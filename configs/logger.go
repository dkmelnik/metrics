package configs

import (
	"flag"
	"os"
)

type Logger struct {
	Mode, Level string
}

func NewLogger() Logger {
	cb := Logger{}

	flag.StringVar(&cb.Mode, "m", "production", "app mode. If empty, production is used")
	flag.StringVar(&cb.Level, "l", "info", "logging level. If empty, warn is used")
	flag.Parse()

	l, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		cb.Level = l
	}

	s, ok := os.LookupEnv("APP_MODE")
	if ok {
		cb.Mode = s
	}

	return cb
}
