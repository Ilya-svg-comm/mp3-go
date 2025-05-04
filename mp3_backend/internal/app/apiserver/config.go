package apiserver

import "backend/internal/app/store"

// Config
type Config struct {
	BinAddr  string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		BinAddr:  ":8080",
		LogLevel: "debug",
		Store:    store.NewConfig(),
	}
}
