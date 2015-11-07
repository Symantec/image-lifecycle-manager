package impl
import (
	"github.com/Symantec/image-lifecycle-manager/pkg/builder"
	"github.com/Symantec/image-lifecycle-manager/pkg/config"
	"github.com/hypersleep/easyssh"
)

type SshPublisher struct {
	Config config.Config
	connection easyssh.MakeConfig
}

func (sp *SshPublisher) init() error{
	user := sp.Config["ssh_publisher_user"]
	server := sp.Config["ssh_publisher_server"]
	key_file := sp.Config["ssh_publisher_key_file"]
	port := sp.Config["ssh_publisher_port"]
	password := sp.Config["ssh_publisher_password"]
	target_path:= sp.Config["ssh_published_target_path"]
	sp.connection = easyssh.MakeConfig{User:user, Server:server, Key:key_file, Port:port, Password:password}

	return sp.connection.Scp(target_path)
}

func (sp *SshPublisher) Publish(artifact builder.Artifact) error{
	//TODO implement or even try to find another lib
	return nil
}
