package test

import (
	"bytes"
	"io"
	"log"

	"github.com/Symantec/image-lifecycle-manager/pkg/builder"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type closableReader struct {
	*bytes.Reader
}

func (*closableReader) Close() error {
	return nil
}

func newClosableReader(b []byte) io.ReadCloser {
	return &closableReader{bytes.NewReader(b)}
}

type TestArtifact struct{}

func (ta *TestArtifact) GetName() string {
	return "TestArtifactName"
}

func (ta *TestArtifact) GetType() builder.ArtifactType {
	return builder.ArtifactLog
}

func (ta *TestArtifact) GetContent() io.ReadCloser {
	return newClosableReader([]byte{1, 2, 3})
}
