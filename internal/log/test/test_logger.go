package test

import (
	"github.com/epistax1s/gomer/internal/config"
	"github.com/epistax1s/gomer/internal/log"
)

// InitTestLogger initializes a simple logger for testing purposes
func InitTestLogger() {
	log.InitLogger(&config.LogConfig{
		Level:  "debug",
		Stdout: true,
	})
}
