package apiserver

import "backend/internal/app/store"

// Config
type Config struct {
	BinAddr  string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
	MinIO    struct {
		Endpoint  string `toml:"endpoint"`
		AccessKey string `toml:"access_key"`
		SecretKey string `toml:"secret_key"`
		Bucket    string `toml:"bucket"`
		UseSSL    bool   `toml:"use_ssl"`
	} `toml:"minio"`
}

func NewConfig() *Config {
	return &Config{
		BinAddr:  ":8080",
		LogLevel: "trace",
		Store:    store.NewConfig(),
	}
}
