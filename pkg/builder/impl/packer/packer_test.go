package packer

import (
	"testing"

	"github.com/Symantec/image-lifecycle-manager/pkg/config"

	"os"

	. "github.com/Symantec/image-lifecycle-manager/test"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type PackerSuite struct {
	working_dir string
	notifier    *TestNotifier
}

var _ = Suite(&PackerSuite{})

func (s *PackerSuite) SetUpTest(c *C) {
	if _, ok := os.LookupEnv("TRAVIS_CI"); ok {
		c.Skip("Can't run test on Travis CI. We need configured and installed packer for this")
	}

	s.working_dir = c.MkDir()
	s.notifier = &TestNotifier{}
	c.Logf("The test dir is %v", s.working_dir)
	conf := config.Config{
		"file_system_publisher_target_path": "cool",
	}
	c.Logf("Config is %v", conf)
}

func (s *PackerSuite) TestConfig(c *C) {
	p := Packer{}
	config := config.Config{}
	c.Assert(p.Configure(config, s.notifier), ErrorMatches, ".*are not defined in config.*")

	config["packer_working_directory"] = "/Users/isviridov/projects/symantec/own-packer"
	config["packer_output_directory"] = "output/vbox_debug_output_ubuntu-14.04_ironic"
	config["packer_execution_command"] = "ls -la"
	config["packer_template_path"] = "imr_template.json"

	c.Assert(p.Configure(config, s.notifier), IsNil)

}

//TODO(illia) decouple test from real packer installation
func (s *PackerSuite) TestValidation(c *C) {
	p := Packer{}

	config := config.Config{}
	config["packer_working_directory"] = "/Users/isviridov/projects/symantec/own-packer"
	config["packer_output_directory"] = "output/vbox_debug_output_ubuntu-14.04_ironic"
	config["packer_execution_command"] = "/usr/local/bin/packer"
	config["packer_template_path"] = "imr_template.json"

	s.notifier.Clean()
	c.Check(p.Configure(config, s.notifier), IsNil)
	c.Check(p.Validate(), IsNil)

	c.Logf("Artifacts %s", p.GetArtifacts())
	c.Check(len(p.GetArtifacts()), Equals, 1)
	c.Check(len(s.notifier.GetNotifications()), Equals, 2,
		Commentf("Expected 2 notifications, but were %s", s.notifier.GetNotifications()))

	p = Packer{}
	config["packer_template_path"] = "no_existing.json"
	s.notifier.Clean()
	c.Check(p.Configure(config, s.notifier), IsNil)
	c.Check(p.Validate(), ErrorMatches, "exit status 1")
	c.Check(len(p.GetArtifacts()), Equals, 1)

	c.Check(p.GetArtifacts()[0], ArtifactContentChecker,
		".*no such file or directory.*",
		Commentf("Content doesn't match"))

	c.Check(len(s.notifier.GetNotifications()), Equals, 2,
		Commentf("Expected 2 notifications, but were %s", s.notifier.GetNotifications()))
}

func (s *PackerSuite) TestBuild(c *C) {
	//TODO(illia) Depends on real packer installation
}
