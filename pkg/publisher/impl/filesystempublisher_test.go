package impl

import (
	"testing"

	"github.com/Symantec/image-lifecycle-manager/pkg/config"
	"github.com/Symantec/image-lifecycle-manager/test"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type FileSystemPublisherSuite struct {
	dir       string
	publisher FileSystemPublisher
	conf      config.Config
}

var _ = Suite(&FileSystemPublisherSuite{})

func (s *FileSystemPublisherSuite) SetUpTest(c *C) {
	s.dir = c.MkDir()
	c.Logf("The test dir is %v", s.dir)
	s.conf = config.Config{
		"file_system_publisher_root_path":   s.dir,
		"file_system_publisher_target_path": "cool",
	}
	c.Logf("Config is %v", s.conf)
	s.publisher = FileSystemPublisher{}
	c.Logf("Publisher %v", s.publisher)
}

func (s *FileSystemPublisherSuite) TestPublishArtifact(c *C) {
	c.Assert(s.publisher.Configure(s.conf, &test.TestNotifier{}), IsNil)
	artifact := test.TestArtifact{}
	c.Logf("artifact is %T and value is %v", artifact, artifact)
	c.Assert(s.publisher.PublishArtifact(&artifact), IsNil)
}

func (s *FileSystemPublisherSuite) TestConfigValidation(c *C) {
	config := config.Config{
		"test": "value",
	}

	notifier := &test.TestNotifier{}
	publisher := FileSystemPublisher{}
	c.Check(publisher.Configure(config, notifier), ErrorMatches, ".*file_system_publisher_root_path is defined.*")

	config["file_system_publisher_root_path"] = "test value"
	c.Check(publisher.Configure(config, notifier), ErrorMatches, ".*file_system_publisher_target_path is defined.*")

	config["file_system_publisher_target_path"] = "test value"
	c.Check(publisher.Configure(config, notifier), ErrorMatches, ".*test value.*no such file or directory")
}
