package config

type Config struct {
	DSN DSN `yaml:"dsn"`
}

type DSN struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

func New() *Config {
	dsn := DSN{
		URL: "my://root:mypass@localhost:33308/testdb",
	}
	c := Config{
		dsn,
	}
	return &c
}
