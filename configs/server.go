package configs

import (
	"flag"
	"os"
)

type Server struct {
	Addr string
}

func NewServer() Server {
	cb := Server{}

	flag.StringVar(&cb.Addr, "a", "0.0.0.0:8080", "in the form host:port. If empty, 0.0.0.0:8080 is used")
	flag.Parse()

	s, ok := os.LookupEnv("ADDRESS")
	if ok {
		cb.Addr = s
	}

	return cb
}
