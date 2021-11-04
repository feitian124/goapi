package config

type Config struct {
	DSN DSN `yaml:"dsn"`
}

type DSN struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

// New return Config
func New() (*Config, error) {
	c := Config{}
	return &c, nil
}
