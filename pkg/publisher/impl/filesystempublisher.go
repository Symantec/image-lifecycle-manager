package impl
import (
	"os"
	"github.com/Symantec/image-lifecycle-manager/pkg/config"
	"github.com/Symantec/image-lifecycle-manager/pkg/builder"
	"fmt"
	"io"
)
const(
	PERMISSION_MASK=0777
)
type FileSystemPublisher struct {
	Config config.Config
	root_path string
	target_path string
}

func (fsp *FileSystemPublisher) Init() error{
	root_path:=fsp.Config["file_system_publisher_root_path"]
	target_path:=fsp.Config["file_system_publisher_target_path"]
	 _, err := os.Stat(root_path)
    if err != nil && os.IsNotExist(err) {
		return err
	}
	if err!=nil{
		return fmt.Errorf("Some other error happened %v", err)
	}

	err = os.MkdirAll(root_path+"/"+target_path, PERMISSION_MASK)
	return err
}

func (fsp *FileSystemPublisher) PublishArtifact(artifact builder.Artifact) error{

	file, err := os.Open(fsp.root_path+"/"+fsp.target_path+"/"+artifact.GetName())
	if err != nil {
		return err
	}

	defer file.Close()

	written, err := io.Copy(file, artifact.GetContent())
	fmt.Printf("Written %v bytes \n", written)
	if err!=nil{
		return err
	}

	err = file.Sync()
	return err
}