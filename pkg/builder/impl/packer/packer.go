package packer

import (
	"fmt"
	"io"

	"os"

	"github.com/Symantec/image-lifecycle-manager/pkg/builder"
	"github.com/Symantec/image-lifecycle-manager/pkg/config"
	"github.com/Symantec/image-lifecycle-manager/pkg/notifier"

	"log"
	"path"

	"github.com/Symantec/image-lifecycle-manager/pkg/context"
)

/*
	{
  	"variables": {
		"output_directory": "output",
		"vm_name": "vm",
		"http_port": "10092",
		"http_host": "172.22.162.153"
 		 },
 	...
*/

const (
	// ParamNameVMName is a name of parameter in template for vm_name -> output image file name
	ParamNameVMName = "vm_name"
	// ParamNameOutDir is a name of parameter in template for output directory
	ParamNameOutDir = "output_directory"
	// ParamNameHTTPPort is a name of parameter in template for http_port
	ParamNameHTTPPort = "http_port"
	// ParamNameHTTPHost is a name of parameter in template for http_host
	ParamNameHTTPHost = "http_host"
	// ParamNameFormat is a name of parameter in template for format
	ParamNameFormat = "format"

	//ParamValueVMName defines the name for virtual machine during build -> hostname, output file name
	ParamValueVMName = "packer-vm"
	// ParamValueFormat is a name of parameter in template for format
	ParamValueFormat = "qcow2"
)

var (
	// TemplateParameters enumerates list of parameters defined in packer template
	TemplateParameters = [...]string{ParamNameVMName, ParamNameOutDir, ParamNameHTTPPort, ParamNameHTTPHost}
)

// Packer image builder implementation
type Packer struct {
	//like /home/illia_svyrydov/packer
	workingDirectory string

	//relative to workingDirectory like: output
	outputDirectory string

	//image file format like: qcow2, ova, ovf. Depends on chosen builder
	outputFileFormat string

	builderHostName string

	//like ./packer
	executionCommand string

	//relative template path f.e. templates/base-ubuntu-14.04-openstack.json
	templatePath string

	notifier  notifier.Notifier
	artifacts []builder.Artifact
}

// GetImage required by builder.ImageBuilder.
func (p *Packer) GetImage() (io.ReadCloser, error) {
	imagePath := path.Join(p.workingDirectory, p.outputDirectory, ParamValueVMName) + "." + p.outputFileFormat
	log.Printf("Looking for image at %s", imagePath)
	file, err := os.Open(imagePath)
	return file, err
}

// Validate required by builder.ImageBuilder.
func (p *Packer) Validate() error {
	//like this [ash2] illia_svyrydov@b0f007ash2012:~/packer$ ./packer validate templates/base-ubuntu-14.04-openstack.json
	//Template validated successfully.
	p.notifier.Notify("packer", "validation", fmt.Sprintf("%s", p))
	artifacts, err := builder.RunCommand(p.workingDirectory,
		p.executionCommand,
		"template_validation",
		[]string{"validate", p.templatePath},
	)
	p.artifacts = append(p.artifacts, artifacts...)
	p.notifier.Notify("packer", "validated", fmt.Sprintf("%s", p))
	return err
}

// Configure required by builder.ImageBuilder.
func (p *Packer) Configure(config config.Config, notifier notifier.Notifier) error {
	p.workingDirectory = config["packer_working_directory"]
	p.outputDirectory = config["packer_output_directory"]
	p.outputFileFormat = config["packer_output_format"]

	p.builderHostName = config["packer_builder_hostname"]
	p.executionCommand = config["packer_execution_command"]
	p.templatePath = config["packer_template_path"]
	p.notifier = notifier

	if p.workingDirectory == "" ||
		p.outputDirectory == "" ||
		p.executionCommand == "" ||
		p.templatePath == "" ||
		p.outputFileFormat == "" {
		return fmt.Errorf("Some of config parameters packer_working_directory, "+
			"packer_output_directory, packer_output_format, packer_execution_command, packer_template_path, packer_builder_hostname "+
			"are not defined in config %s", config)
	}
	return nil
}

// BuildImage required by builder.ImageBuilder.
//TODO(illia)implement sending notification during building process with current output
func (p *Packer) BuildImage() error {
	p.notifier.Notify("packer", "building", "")
	artifacts, err := builder.RunCommand(p.workingDirectory,
		p.executionCommand,
		"build_image",
		p.buildPackerVars("build", p.templatePath),
	)
	p.artifacts = append(p.artifacts, artifacts...)
	p.notifier.Notify("packer", "finished", "")
	return err
}

// GetArtifacts required by builder.ImageBuilder.
func (p *Packer) GetArtifacts() []builder.Artifact {
	return p.artifacts
}

func (p *Packer) buildPackerVars(action, template string) []string {
	varArgs := []string{
		action,
		"-var", fmt.Sprintf("%s=%v", ParamNameVMName, ParamValueVMName),
		"-var", fmt.Sprintf("%s=%v", ParamNameFormat, p.outputFileFormat),
		"-var", fmt.Sprintf("%s=%v", ParamNameOutDir, p.outputDirectory),
		"-var", fmt.Sprintf("%s=%v", ParamNameHTTPHost, p.builderHostName),
		"-var", fmt.Sprintf("%s=%v", ParamNameHTTPPort, context.GetFreePort()),
		template,
	}
	return varArgs
}
