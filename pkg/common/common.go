package common

import (
	"bytes"
	"io"
)

type closableReader struct {
	*bytes.Reader
}

func (*closableReader) Close() error {
	return nil
}

func NewClosableReader(b []byte) io.ReadCloser {
	return &closableReader{bytes.NewReader(b)}
}
