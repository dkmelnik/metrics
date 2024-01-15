package configs

import (
	"flag"
	"os"
)

type Server struct {
	Addr, LogLevel string
}

func NewServer() Server {
	cb := Server{}

	flag.StringVar(&cb.Addr, "a", "0.0.0.0:8080", "in the form host:port. If empty, 0.0.0.0:8080 is used")
	flag.StringVar(&cb.LogLevel, "l", "warn", "logging level. If empty, warn is used")
	flag.Parse()

	l, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		cb.LogLevel = l
	}

	s, ok := os.LookupEnv("ADDRESS")
	if ok {
		cb.Addr = s
	}

	return cb
}
