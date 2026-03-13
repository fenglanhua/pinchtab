package main

import "github.com/pinchtab/pinchtab/internal/config"

func loadConfig() *config.RuntimeConfig {
	return config.Load()
}
