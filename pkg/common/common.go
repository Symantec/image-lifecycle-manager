package common

import (
	"bytes"
	"fmt"
	"io"
	"strings"
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

// ParamList parses multi-value string parameters
type ParamList struct {
	Params map[string]string
}

func (p *ParamList) String() string {
	return fmt.Sprint("%v", p.Params)
}

func (p *ParamList) Set(value string) error {
	if value == "" {
		return fmt.Errorf("Parameter has empty value")
	}

	pair := strings.Split(value, "=")
	if len(pair) != 2 {
		return fmt.Errorf("Wrong paramater format. Expected 'name=value', but was %s", value)
	}
	p.Params[pair[0]] = pair[1]
	return nil
}

type SystemCaller interface {
}
