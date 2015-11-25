package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Symantec/image-lifecycle-manager/pkg/builder"
	"github.com/Symantec/image-lifecycle-manager/pkg/builder/impl/packer"
	"github.com/Symantec/image-lifecycle-manager/pkg/common"
	"github.com/Symantec/image-lifecycle-manager/pkg/config"
	"github.com/Symantec/image-lifecycle-manager/pkg/notifier"
	"github.com/Symantec/image-lifecycle-manager/pkg/publisher"
	"github.com/Symantec/image-lifecycle-manager/pkg/publisher/impl"
)

var (
	templateParam = flag.String("template", "",
		"Template file of image to build. Example: ubuntu-14.04.json")
	helpParam = flag.Bool("help", false,
		"Prints out usage information")
	builderParam = flag.String("builder", "packer",
		"The builder used for building image. Allowed values [packer]. Default: packer")
	publisherParam = flag.String("publisher", "filesystem",
		"The publisher used for publish build artifacts. Allowed values [filesystem]. Default: filesystem")
	notifierParam = flag.String("notifier", "log",
		"The notifier used for sending notifications. Allowed values [log]. Default: log")
	configSourceParam = flag.String("configSource", "env",
		"Source of configurations of submodules like builder, publisher, notifier. Allowed values [env]. Default: env")
	paramsParam map[string]string = map[string]string{}
	paramList   common.ParamList  = common.ParamList{paramsParam}
	arts                          = []builder.Artifact{}

	exit = os.Exit
)

func init() {
	flag.Var(&paramList, "param", "Parameters for building process. "+
		"Paramaters are set here will overwrite any other previously defined. "+
		"Example: -param template_path=/var/test.json -param working_dir=/var")
}

func main() {
	flag.Parse()
	if *helpParam {
		flag.Usage()
		exit(0)
	}

	//building config
	config := getConfig(*configSourceParam)
	mergeConfigs(config, *templateParam, paramsParam)

	//building image builder and dependencies
	b := getImageBuilder(*builderParam)
	notifier := getNotifier(*notifierParam)
	p := getPublisher(*publisherParam)

	err := b.Configure(config, notifier)
	if err != nil {
		log.Printf("Error during builder configuration: %s", err)
		exit(1)
	}

	err = p.Configure(config, notifier)
	if err != nil {
		log.Printf("Error during publisher configuration: %s", err)
		exit(1)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovering from panic : %s", r)
			publishArtifacts(p, arts)
			exit(1)
		}
	}()

	//validating image and environment: utilities, folders permissions so on.
	err = b.Validate()
	arts = append(arts, b.GetArtifacts()...)
	if err != nil {
		panic(fmt.Sprintf("Error during validation of template %s", err))
	}

	//building image
	err = b.BuildImage()
	arts = append(arts, b.GetArtifacts()...)
	if err != nil {
		panic(fmt.Sprintf("Error during building image: %s", err))
	}

	//publishing image artifact
	imageArtifact := builder.NewImageArtifact("image", func() io.ReadCloser {
		if content, err := b.GetImage(); err != nil {
			return nil
		} else {
			return content
		}
	})

	arts = append(arts, imageArtifact)
	publishArtifacts(p, arts)
}

func publishArtifacts(p publisher.Publisher, arts []builder.Artifact) {
	//publishing artifacts
	for _, a := range arts {
		if err := p.PublishArtifact(a); err != nil {
			log.Printf("Error during publishing artifact %s. err: %s", a, err)
		}
	}
}

// Merges config with command line parameters
func mergeConfigs(config config.Config, templateName string, params map[string]string) {
	for key, value := range params {
		config[key] = value
	}
	config["packer_template_path"] = templateName
}

func getNotifier(name string) notifier.Notifier {
	switch name {
	case "log", "":
		return &notifier.LogNotifier{}
	default:
		log.Printf("Unknown notifier type %s.  Using default one", name)
		return &notifier.LogNotifier{}
	}
}

func getPublisher(name string) publisher.Publisher {
	switch name {
	case "filesystem", "":
		return &impl.FileSystemPublisher{}
	default:
		log.Printf("Unknown publisher type %s.  Using default one", name)
		return &impl.FileSystemPublisher{}
	}
}

func getImageBuilder(name string) builder.ImageBuilder {
	switch name {
	case "packer", "":
		return &packer.Packer{}
	default:
		log.Printf("Unknown image builder type %s.  Using default one", name)
		return &packer.Packer{}
	}
}

func getConfig(source string) config.Config {
	switch source {
	case "env", "":
		return config.BuildEnvConfig()
	default:
		log.Printf("Unknown config source %s.  Using default one", source)
		return config.BuildEnvConfig()
	}
}
