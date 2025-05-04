package store

type Config struct {
	DataBaseURL string `toml:"database_url"`
}

// NewConfig...
func NewConfig() *Config {
	return &Config{}
}
