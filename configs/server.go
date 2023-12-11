package configs

import (
	"flag"
	"os"
)

type Server struct {
	Addr string
}

func NewServer() Server {
	return Server{}
}

func (cb Server) Build() Server {
	flag.StringVar(&cb.Addr, "a", "127.0.0.1:8080", "in the form host:port. If empty, 127.0.0.1:8080 is used")
	flag.Parse()

	s, ok := os.LookupEnv("ADDRESS")
	if ok {
		cb.Addr = s
	}

	return cb
}
