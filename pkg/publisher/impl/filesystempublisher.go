package impl

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Symantec/image-lifecycle-manager/pkg/builder"
	"github.com/Symantec/image-lifecycle-manager/pkg/config"
	"github.com/Symantec/image-lifecycle-manager/pkg/notifier"
)

var (
	// FolderPermissionMask is default permissions for newle created folder
	FolderPermissionMask os.FileMode = 0777
)

// FileSystemPublisher publishes artifact on filesystem
type FileSystemPublisher struct {
	config      config.Config
	root_path   string
	target_path string
}

func (fsp *FileSystemPublisher) Configure(conf config.Config, notifier notifier.Notifier) error {
	fsp.config = conf

	if root_path, ok := fsp.config["file_system_publisher_root_path"]; ok {
		fsp.root_path = root_path
	} else {
		return fmt.Errorf("No file_system_publisher_root_path is defined in config %v", fsp.config)
	}

	if target_path, ok := fsp.config["file_system_publisher_target_path"]; ok {
		fsp.target_path = target_path
	} else {
		return fmt.Errorf("No file_system_publisher_target_path is defined in config %v", fsp.config)
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

	err = os.MkdirAll(fsp.root_path+"/"+fsp.target_path, FolderPermissionMask)
	return err
}

func (fsp *FileSystemPublisher) PublishArtifact(artifact builder.Artifact) error {
	if fsp.config == nil {
		return fmt.Errorf("FileSystemPublisher is not initialized")
	}
	file, err := os.Create(fsp.root_path + "/" + fsp.target_path + "/" + artifact.GetName())
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, artifact.GetContent())
	if err != nil {
		return err
	}

	err = file.Sync()
	return err
}
