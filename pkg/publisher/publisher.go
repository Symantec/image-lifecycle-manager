package publisher

import (
	"github.com/Symantec/image-lifecycle-manager/pkg/builder"
	"github.com/Symantec/image-lifecycle-manager/pkg/config"
	"github.com/Symantec/image-lifecycle-manager/pkg/notifier"
)

type Publisher interface {
	PublishArtifact(artifact builder.Artifact) error
	Configure(config.Config, notifier.Notifier) error
}
