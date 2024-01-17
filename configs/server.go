package configs

import (
	"flag"
	"os"
)

type Server struct {
	Addr, DBConnectStr string
}

func NewServer() Server {
	cb := Server{}

	flag.StringVar(&cb.Addr, "a", "0.0.0.0:8080", "in the form host:port. If empty, 0.0.0.0:8080 is used")
	flag.StringVar(&cb.DBConnectStr, "d", "", "string for db connect")
	flag.Parse()

	s, ok := os.LookupEnv("ADDRESS")
	if ok {
		cb.Addr = s
	}

	db, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		cb.DBConnectStr = db
	}

	return cb
}
