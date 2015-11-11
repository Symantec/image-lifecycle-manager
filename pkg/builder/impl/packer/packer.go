package packer

import (
	"fmt"
	"io"

	"github.com/Symantec/image-lifecycle-manager/pkg/builder"
	"github.com/Symantec/image-lifecycle-manager/pkg/config"
	"github.com/Symantec/image-lifecycle-manager/pkg/notifier"
)

// Packer image builder implementation
type Packer struct {
	//like /home/illia_svyrydov/packer
	working_directory string
	// like /home/illia_svyrydov/packer/output/output_centos-6_ironic
	// or output/output_centos-6_ironic
	output_directory string
	// like ./packer
	execution_command string
	//relative template path f.e. templates/base-ubuntu-14.04-openstack.json
	template_path string

	notifier  notifier.Notifier
	artifacts []builder.Artifact
}

func (p *Packer) GetImage() io.ReadCloser {
	return nil
}

func (p *Packer) Validate() error {
	//like this [ash2] illia_svyrydov@b0f007ash2012:~/packer$ ./packer validate templates/base-ubuntu-14.04-openstack.json
	//Template validated successfully.
	p.notifier.Notify("packer", "validation", fmt.Sprintf("%s", p))
	artifacts, err := builder.RunCommand(p.working_directory,
		p.execution_command,
		"template_validation",
		[]string{"validate", p.template_path},
	)
	p.artifacts = append(p.artifacts, artifacts...)
	p.notifier.Notify("packer", "validated", fmt.Sprintf("%s", p))
	return err
}

func (p *Packer) Configure(config config.Config, notifier notifier.Notifier) error {
	p.working_directory = config["packer_working_directory"]
	//TODO(illia) probably should be retrieved from template
	p.output_directory = config["packer_output_directory"]
	p.execution_command = config["packer_execution_command"]
	p.template_path = config["packer_template_path"]
	p.notifier = notifier

	if p.working_directory == "" ||
		p.output_directory == "" ||
		p.execution_command == "" ||
		p.template_path == "" {
		return fmt.Errorf("Some of config parameters packer_working_directory, "+
			"packer_output_directory, packer_execution_command, packer_template_path "+
			"are not defined in config %s", config)
	}
	return nil
}

//TODO(illia)implement sending notification during building process with current output
func (p *Packer) BuildImage() error {
	p.notifier.Notify("packer", "building", "")
	artifacts, err := builder.RunCommand(p.working_directory,
		p.execution_command,
		"build_image",
		[]string{"build", p.template_path},
	)
	p.artifacts = append(p.artifacts, artifacts...)
	p.notifier.Notify("packer", "finished", "")
	return err
}

func (p *Packer) GetArtifacts() []builder.Artifact {
	return p.artifacts
}
