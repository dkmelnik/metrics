package configs

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Server struct {
	Addr            string `json:"addr"`
	DBConnectStr    string `json:"db_connect_str"`
	Mode            string `json:"mode"`
	Level           string `json:"level"`
	FileStoragePath string `json:"file_storage_path"`
	StoreInterval   int    `json:"store_interval"`
	Restore         bool   `json:"restore"`
	PrivateKeyPath  string `json:"private_key_path"`
	Key             string `json:"key"`
	configPath      string
}

func NewServer() Server {
	c := Server{}

	return c.setFlagValues().
		setEnvValues().
		setFileValues()
}

func (c Server) setFlagValues() Server {
	flag.StringVar(&c.configPath, "c", c.configPath, "path to configuration file")
	flag.StringVar(&c.Addr, "a", "0.0.0.0:8080", "address in the form host:port. If empty, 0.0.0.0:8080 is used")
	flag.StringVar(&c.DBConnectStr, "d", c.DBConnectStr, "Database Connection String")
	flag.StringVar(&c.Mode, "m", "production", "app mode. If empty, production is used")
	flag.StringVar(&c.Level, "l", "info", "logging level. If empty, warn is used")
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "full name of the file where the current values are saved. If empty, /tmp/metrics-db.json is used")
	flag.IntVar(&c.StoreInterval, "i", 300, "server saved metrics to disk. If empty, 300 is used")
	flag.BoolVar(&c.Restore, "r", true, "load or not previously saved values. If empty, true is used")
	flag.StringVar(&c.PrivateKeyPath, "crypto-key", "", "private key address for asymmetric encryption")
	flag.StringVar(&c.Key, "k", "", "signature key")
	flag.Parse()

	return c
}

func (c Server) setEnvValues() Server {
	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		c.Addr = envAddr
	}
	if envDBConnectStr := os.Getenv("DATABASE_DSN"); envDBConnectStr != "" {
		c.DBConnectStr = envDBConnectStr
	}
	if envMode := os.Getenv("APP_MODE"); envMode != "" {
		c.Mode = envMode
	}
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		c.Level = envLevel
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		c.FileStoragePath = envFileStoragePath
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		fmt.Sscanf(envStoreInterval, "%d", &c.StoreInterval)
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		c.Restore = envRestore == "true"
	}
	if envPrivateKeyPath := os.Getenv("CRYPTO_KEY"); envPrivateKeyPath != "" {
		c.PrivateKeyPath = envPrivateKeyPath
	}
	if envKey := os.Getenv("KEY"); envKey != "" {
		c.Key = envKey
	}

	return c
}

func (c Server) setFileValues() Server {
	configPath := c.configPath

	envConfigFile := os.Getenv("CONFIG")
	if envConfigFile != "" {
		configPath = envConfigFile
	}

	config := Server{}

	if configPath == "" {
		return c
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return c
	}

	if err = json.Unmarshal(file, &config); err != nil {
		return c
	}

	if c.Addr == "" {
		c.Addr = config.Addr
	}

	if c.DBConnectStr == "" {
		c.DBConnectStr = config.DBConnectStr
	}

	if c.Mode == "" {
		c.Mode = config.Mode
	}

	if c.Level == "" {
		c.Level = config.Level
	}

	if c.FileStoragePath == "" {
		c.FileStoragePath = config.FileStoragePath
	}

	if c.StoreInterval == 0 {
		c.StoreInterval = config.StoreInterval
	}

	if c.Restore == false {
		c.Restore = config.Restore
	}

	if c.PrivateKeyPath == "" {
		c.PrivateKeyPath = config.PrivateKeyPath
	}

	if c.Key == "" {
		c.Key = config.Key
	}

	return c
}

func (c Server) GetLevel() string {
	return c.Level
}

func (c Server) GetMode() string {
	return c.Mode
}
