package impl

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Symantec/image-lifecycle-manager/pkg/builder"
	"github.com/Symantec/image-lifecycle-manager/pkg/config"
)

const (
	PERMISSION_MASK = 0777
)

type FileSystemPublisher struct {
	Config      config.Config
	root_path   string
	target_path string
}

func (fsp *FileSystemPublisher) Init() error {

	if root_path, ok := fsp.Config["file_system_publisher_root_path"]; ok {
		fsp.root_path = root_path
	} else {
		return fmt.Errorf("No file_system_publisher_root_path is defined in config %v", fsp.Config)
	}

	if target_path, ok := fsp.Config["file_system_publisher_target_path"]; ok {
		fsp.target_path = target_path
	} else {
		return fmt.Errorf("No file_system_publisher_target_path is defined in config %v", fsp.Config)
	}

	log.Printf("root_path '%s'", fsp.root_path)
	log.Printf("target_path '%s'", fsp.target_path)
	_, err := os.Stat(fsp.root_path)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	if err != nil {
		return fmt.Errorf("Some other error happened %v", err)
	}

	err = os.MkdirAll(fsp.root_path+"/"+fsp.target_path, PERMISSION_MASK)
	return err
}

func (fsp *FileSystemPublisher) PublishArtifact(artifact builder.Artifact) error {

	file, err := os.Create(fsp.root_path + "/" + fsp.target_path + "/" + artifact.GetName())
	if err != nil {
		return err
	}

	defer file.Close()

	written, err := io.Copy(file, artifact.GetContent())
	fmt.Printf("Written %v bytes \n", written)
	if err != nil {
		return err
	}

	err = file.Sync()
	return err
}
