package configs

import (
	"flag"
	"os"
	"strconv"
)

type Server struct {
	Addr, LogLevel, FileStoragePath string
	StoreInterval                   int
	Restore                         bool
}

func NewServer() Server {
	cb := Server{}

	flag.StringVar(&cb.Addr, "a", "0.0.0.0:8080", "in the form host:port. If empty, 0.0.0.0:8080 is used")
	flag.StringVar(&cb.LogLevel, "l", "warn", "logging level. If empty, warn is used")
	flag.StringVar(&cb.FileStoragePath, "f", "/tmp/metrics-db.json", "full name of the file where the current values are saved. If empty, /tmp/metrics-db.json is used")
	flag.IntVar(&cb.StoreInterval, "i", 10, "server saved metrics to disk. If empty, 300 is used")
	flag.BoolVar(&cb.Restore, "r", true, "load or not previously saved values. If empty, true is used")
	flag.Parse()

	l, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		cb.LogLevel = l
	}

	s, ok := os.LookupEnv("ADDRESS")
	if ok {
		cb.Addr = s
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

	return cb
}
