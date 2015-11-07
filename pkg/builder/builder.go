// Package builder defines basic interfaces of image builder module and contains implementations
package builder

import (
	"io"
	"github.com/Symantec/image-lifecycle-manager/pkg/config"
)

// ImageBuilder defines interface for implementations
type ImageBuilder interface{
	GetImage() io.Reader
	GetStatus() Status
	Validate() error
	Configure(config.Config) error
	GetMetadata(map[string]string)
	BuildImage() error
	GetArtifacts([]Artifact)
}
// Status defined all possible statuses of building image process
type Status string

const(
	StatusError Status = "ERROR"
	StatusOk Status = "OK"
)

// ArtifactType defined all supported artifact types
type ArtifactType string

const(
	ArtifactLog ArtifactType = "Log"
	ArtifactImage ArtifactType = "Image"
)

// Artifact defines interface of artifact
type Artifact interface{
	GetName() string
	GetType() ArtifactType
	GetContent() io.Reader
}
