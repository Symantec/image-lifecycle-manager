// Package builder defines basic interfaces of image builder module and contains implementations
package builder

import (
	"io"
	"io/ioutil"
	"os"

	"bytes"
	"fmt"

	"github.com/Symantec/image-lifecycle-manager/pkg/common"
	"github.com/Symantec/image-lifecycle-manager/pkg/config"
	ctx "github.com/Symantec/image-lifecycle-manager/pkg/context"
	"github.com/Symantec/image-lifecycle-manager/pkg/notifier"
)

// ImageBuilder defines interface for implementations
type ImageBuilder interface {
	GetImage() (io.ReadCloser, error)
	Validate() error
	Configure(config.Config, notifier.Notifier) error
	BuildImage() error
	GetArtifacts() []Artifact
	//Cleanup()
}

// Status defined all possible statuses of building image process
type Status string

//TODO(illia) think if status is needed at all
const (
	StatusError Status = "ERROR"
	StatusOk    Status = "OK"
)

// ArtifactType defined all supported artifact types
type ArtifactType string

const (
	ArtifactLog   ArtifactType = "Log"
	ArtifactImage ArtifactType = "Image"
)

// Artifact defines interface of artifact
type Artifact interface {
	GetName() string
	GetType() ArtifactType
	GetContent() io.ReadCloser
}

type FileArtifact struct {
	Name string
	Path string
}

func (f *FileArtifact) GetName() string {
	return f.Name
}

func (f *FileArtifact) GetType() ArtifactType {
	return ArtifactLog
}

func (f *FileArtifact) GetContent() io.ReadCloser {
	file, err := os.Open(f.Path)
	if err == nil {
		return file
	}
	return nil
}

type InMemoryArtifact struct {
	Name         string
	Data         []byte
	ArtifactType ArtifactType
}

func (i *InMemoryArtifact) GetName() string {
	return i.Name
}

func (i *InMemoryArtifact) GetType() ArtifactType {
	return i.ArtifactType
}

func (i *InMemoryArtifact) GetContent() io.ReadCloser {
	return common.NewClosableReader(i.Data)
}

func (i *InMemoryArtifact) String() string {
	cut := i.Data
	if len(i.Data) > 1024 {
		cut = cut[:1024]
	}
	return fmt.Sprintf("name:%s, type:%s, data:%s", i.Name, i.ArtifactType, cut)
}

// ImageArtifact defines artifact for image
type ImageArtifact struct {
	Name            string
	contentResolver func() io.ReadCloser
}

func (i *ImageArtifact) GetName() string {
	return i.Name
}

func (i *ImageArtifact) GetType() ArtifactType {
	return ArtifactImage
}

func (i *ImageArtifact) GetContent() io.ReadCloser {
	return i.contentResolver()
}

// NewImageArtifact creates new artifact for image
func NewImageArtifact(name string, contentResolver func() io.ReadCloser) Artifact {
	return &ImageArtifact{name, contentResolver}
}

// NewTempFileArtifact creates file artifact from bytes.Buffer.
func NewTempFileArtifact(data bytes.Buffer, name string) (Artifact, error) {
	f, err := ioutil.TempFile("", "packer_builder")
	if err != nil {
		return nil, err
	}

	if _, err = f.Write(data.Bytes()); err != nil {
		return nil, err
	}
	a := FileArtifact{Name: name, Path: f.Name()}
	return &a, nil
}

// RunCommand execute command and collect std_err and std_out as artifacts.
// Returns error if any and list of artifacts if any
func RunCommand(dir string, cmd string, artifact_prefix string, args []string) ([]Artifact, error) {
	std_out := bytes.Buffer{}
	std_err := bytes.Buffer{}

	err := ctx.SystemCall(dir, cmd, args, &std_out, &std_err)
	artifacts := []Artifact{}

	if std_err.Len() != 0 {
		artifacts = append(artifacts,
			&InMemoryArtifact{
				Data:         std_err.Bytes(),
				ArtifactType: ArtifactLog,
				Name:         artifact_prefix + "_std_err"})
	}

	if std_out.Len() != 0 {
		artifacts = append(artifacts,
			&InMemoryArtifact{
				Data:         std_out.Bytes(),
				ArtifactType: ArtifactLog,
				Name:         artifact_prefix + "_std_out"})
	}
	return artifacts, err
}
