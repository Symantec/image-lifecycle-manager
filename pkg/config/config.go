package config

import (
	"os"
	"strings"
)

type Config map[string]string

func (*Config) String() string {
	return "Condig string representation"
}

// BuildEnvConfig builds config object from os.Environ()
func BuildEnvConfig() Config {
	config := Config{}
	for _, v := range os.Environ() {
		pair := strings.Split(v, "=")
		config[pair[0]] = pair[1]
	}
	return config
}
