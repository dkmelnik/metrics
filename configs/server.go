package configs

import (
	"flag"
	"os"
	"strconv"
)

// Server stores properties that configure server service.
// Properties can be taken from environment variables or flags.
type Server struct {
	Addr, DBConnectStr, Mode, Level, FileStoragePath, Key string
	StoreInterval                                         int
	Restore                                               bool
}

func NewServer() Server {
	cb := Server{}

	flag.StringVar(&cb.Addr, "a", "0.0.0.0:8080", "in the form host:port. If empty, 0.0.0.0:8080 is used")
	flag.StringVar(&cb.DBConnectStr, "d", "", "string for db connect")
	flag.StringVar(&cb.Mode, "m", "production", "app mode. If empty, production is used")
	flag.StringVar(&cb.Level, "l", "info", "logging level. If empty, warn is used")
	flag.StringVar(&cb.FileStoragePath, "f", "/tmp/metrics-db.json", "full name of the file where the current values are saved. If empty, /tmp/metrics-db.json is used")
	flag.IntVar(&cb.StoreInterval, "i", 300, "server saved metrics to disk. If empty, 300 is used")
	flag.BoolVar(&cb.Restore, "r", true, "load or not previously saved values. If empty, true is used")
	flag.StringVar(&cb.Key, "k", "", "signature key")
	flag.Parse()

	k, ok := os.LookupEnv("KEY")
	if ok {
		cb.Key = k
	}

	f, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		cb.FileStoragePath = f
	}

	r, ok := os.LookupEnv("RESTORE")
	if ok {
		toBool, err := strconv.ParseBool(r)
		if err == nil {
			cb.Restore = toBool
		}
	}

	i, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		toInt, err := strconv.Atoi(i)
		if err == nil {
			cb.StoreInterval = toInt
		}
	}

	l, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		cb.Level = l
	}

	s, ok := os.LookupEnv("APP_MODE")
	if ok {
		cb.Mode = s
	}

	sa, ok := os.LookupEnv("ADDRESS")
	if ok {
		cb.Addr = sa
	}

	db, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		cb.DBConnectStr = db
	}

	return cb
}

func (c Server) GetLevel() string {
	return c.Level
}

func (c Server) GetMode() string {
	return c.Mode
}
