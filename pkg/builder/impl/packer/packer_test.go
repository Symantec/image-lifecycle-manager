package packer

import (
	"testing"

	"github.com/Symantec/image-lifecycle-manager/pkg/config"

	"fmt"

	. "github.com/Symantec/image-lifecycle-manager/test"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type PackerBaseSuite struct {
	workingDir string
	notifier   *TestNotifier
	config     config.Config
	packer     Packer
}

type PackerUnitSuite struct {
	PackerBaseSuite
}

var (
	_ = Suite(&PackerUnitSuite{})
)

func (s *PackerUnitSuite) SetUpTest(c *C) {
	s.workingDir = c.MkDir()
	s.notifier = &TestNotifier{}

	s.config = config.Config{}
	s.config["packer_working_directory"] = "/Users/isviridov/projects/symantec/own-packer"
	s.config["packer_output_directory"] = "output/vbox_debug_output_ubuntu-14.04_ironic"
	s.config["packer_output_format"] = "ovf"
	s.config["packer_execution_command"] = "ls -la"
	s.config["packer_template_path"] = "imr_template.json"
	s.config["packer_builder_hostname"] = "host"

	s.packer = Packer{}

}

func (s *PackerUnitSuite) TestConfig(c *C) {
	delete(s.config, "packer_template_path")

	c.Assert(s.packer.Configure(s.config, s.notifier), ErrorMatches, ".*are not defined in config.*")
	s.config["packer_template_path"] = "imr_template.json"
	c.Assert(s.packer.Configure(s.config, s.notifier), IsNil)

}

func (s *PackerUnitSuite) TestValidation(c *C) {
	c.Check(s.packer.Configure(s.config, s.notifier), IsNil)

	MockSystemCall("Validation passed", "", nil)
	c.Check(s.packer.Validate(), IsNil)

	c.Logf("Artifacts %s", s.packer.GetArtifacts())
	c.Check(len(s.packer.GetArtifacts()), Equals, 1)
	c.Check(len(s.notifier.GetNotifications()), Equals, 2,
		Commentf("Expected 2 notifications, but were %s", s.notifier.GetNotifications()))

}

func (s *PackerUnitSuite) TestValidationNoExistingTemplate(c *C) {

	MockSystemCall("", "no such file or directory", fmt.Errorf("exit status 1"))
	c.Check(s.packer.Configure(s.config, s.notifier), IsNil)
	c.Check(s.packer.Validate(), ErrorMatches, "exit status 1")
	c.Check(len(s.packer.GetArtifacts()), Equals, 1)

	c.Check(s.packer.GetArtifacts()[0], ArtifactContentChecker,
		".*no such file or directory.*",
		Commentf("Content doesn't match"))

	c.Check(len(s.notifier.GetNotifications()), Equals, 2,
		Commentf("Expected 2 notifications, but were %s", s.notifier.GetNotifications()))
}

func (s *PackerUnitSuite) TestBuild(c *C) {
	//TODO(illia) cover with test
}
