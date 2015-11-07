package config

type Config map[string]string

func (*Config) String() string {
	return "Condig string representation"
}