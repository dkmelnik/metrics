package configs

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
)

// Agent stores properties that configure the agent service.
// Properties can be taken from environment variables or flags.
type Agent struct {
	Addr           string `json:"addr"`
	Mode           string `json:"mode"`
	Level          string `json:"level"`
	Key            string `json:"key"`
	PublicKeyPath  string `json:"public_key_path"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
	RateLimit      int    `json:"rate_limit"`
	configPath     string
}

// NewAgent initializes a new Agent with default values and parses flags and environment variables.
func NewAgent() Agent {
	ca := Agent{}

	return ca.setFlagValues().
		setEnvValues().
		setFileValues()
}

// setFlagValues sets configuration values from command line flags.
func (c Agent) setFlagValues() Agent {
	flag.StringVar(&c.configPath, "c", c.configPath, "path to configuration file")
	flag.StringVar(&c.Addr, "a", "http://localhost:8080", "server by collected metric address")
	flag.StringVar(&c.Mode, "m", "production", "app mode. If empty, production is used")
	flag.StringVar(&c.Level, "la", "info", "logging level. If empty, warn is used")
	flag.StringVar(&c.Key, "k", "", "signature key")
	flag.IntVar(&c.ReportInterval, "r", 10, "period for sending metrics to the server")
	flag.IntVar(&c.PollInterval, "p", 2, "metrics collection period")
	flag.IntVar(&c.RateLimit, "l", 5, "req rate limit. If empty, 5 is used")
	flag.StringVar(&c.PublicKeyPath, "crypto-key", "", "public key address for asymmetric encryption")
	flag.Parse()

	return c
}

// setEnvValues sets configuration values from environment variables.
func (c Agent) setEnvValues() Agent {
	setEnvString(&c.Addr, "ADDRESS")
	setEnvString(&c.Mode, "APP_MODE")
	setEnvString(&c.Level, "LOG_LEVEL")
	setEnvString(&c.Key, "KEY")
	setEnvString(&c.PublicKeyPath, "CRYPTO_KEY")
	setEnvInt(&c.ReportInterval, "REPORT_INTERVAL")
	setEnvInt(&c.PollInterval, "POLL_INTERVAL")
	setEnvInt(&c.RateLimit, "RATE_LIMIT")

	return c
}

// setFileValues sets configuration values from a JSON file.
func (c Agent) setFileValues() Agent {
	configPath := c.configPath

	envConfigFile := os.Getenv("CONFIG")
	if envConfigFile != "" {
		configPath = envConfigFile
	}

	config := Agent{}

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

	if c.Mode == "" {
		c.Mode = config.Mode
	}

	if c.Level == "" {
		c.Level = config.Level
	}

	if c.Key == "" {
		c.Key = config.Key
	}

	if c.PublicKeyPath == "" {
		c.PublicKeyPath = config.PublicKeyPath
	}

	if c.ReportInterval == 0 {
		c.ReportInterval = config.ReportInterval
	}

	if c.PollInterval == 0 {
		c.PollInterval = config.PollInterval
	}

	if c.RateLimit == 0 {
		c.RateLimit = config.RateLimit
	}

	return c
}

// setEnvString sets string value from environment variable if available.
func setEnvString(value *string, envName string) {
	if v, ok := os.LookupEnv(envName); ok {
		*value = v
	}
}

// setEnvInt sets integer value from environment variable if available.
func setEnvInt(value *int, envName string) {
	if v, ok := os.LookupEnv(envName); ok {
		if i, err := strconv.Atoi(v); err == nil {
			*value = i
		}
	}
}
