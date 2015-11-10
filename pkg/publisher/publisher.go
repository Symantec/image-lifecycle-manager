package publisher

import (
	"github.com/Symantec/image-lifecycle-manager/pkg/builder"
	"github.com/Symantec/image-lifecycle-manager/pkg/config"
)

type Published interface {
	PublishArtifact(artifact builder.Artifact) error
	Configure(config.Config) error
}
