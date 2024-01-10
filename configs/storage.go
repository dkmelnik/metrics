package configs

import (
	"flag"
	"os"
	"strconv"
)

type Storage struct {
	FileStoragePath string
	StoreInterval   int
	Restore         bool
}

func NewStorage() Storage {
	cb := Storage{}

	flag.StringVar(&cb.FileStoragePath, "f", "/tmp/metrics-db.json", "full name of the file where the current values are saved. If empty, /tmp/metrics-db.json is used")
	flag.IntVar(&cb.StoreInterval, "i", 300, "server saved metrics to disk. If empty, 300 is used")
	flag.BoolVar(&cb.Restore, "r", true, "load or not previously saved values. If empty, true is used")
	flag.Parse()

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
