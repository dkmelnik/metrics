package configs

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

// Agent stores properties that configure agent service.
// Properties can be taken from environment variables or flags.
type Agent struct {
	Addr, Mode, Level, Key                  string
	ReportInterval, PollInterval, RateLimit int
}

func NewAgent() Agent {
	cb := Agent{}

	flag.StringVar(&cb.Addr, "a", "http://localhost:8080", "server by collected metric address")
	flag.IntVar(&cb.ReportInterval, "r", 10, "period for sending metrics to the server")
	flag.IntVar(&cb.PollInterval, "p", 2, "metrics collection period")
	flag.StringVar(&cb.Mode, "m", "production", "app mode. If empty, production is used")
	flag.StringVar(&cb.Level, "la", "info", "logging level. If empty, warn is used")
	flag.StringVar(&cb.Key, "k", "", "signature key")
	flag.IntVar(&cb.RateLimit, "l", 5, "req rate limit. If empty, 5 is used")
	flag.Parse()

	k, ok := os.LookupEnv("KEY")
	if ok {
		cb.Key = k
	}

	l, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		cb.Level = l
	}

	m, ok := os.LookupEnv("APP_MODE")
	if ok {
		cb.Mode = m
	}

	sa, ok := os.LookupEnv("ADDRESS")
	if ok {
		cb.Addr = sa
	}

	s, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		i, err := strconv.Atoi(s)
		if err == nil {
			cb.ReportInterval = i
		}
	}

	r, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		i, err := strconv.Atoi(r)
		if err == nil {
			cb.RateLimit = i
		}
	}

	s, ok = os.LookupEnv("POLL_INTERVAL")
	if ok {
		i, err := strconv.Atoi(s)
		if err == nil {
			cb.PollInterval = i
		}
	}
	if !strings.HasPrefix(cb.Addr, "http://") && !strings.HasPrefix(cb.Addr, "https://") {
		cb.Addr = "http://" + cb.Addr
	}
	return cb
}

func (c Agent) GetLevel() string {
	return c.Level
}

func (c Agent) GetMode() string {
	return c.Mode
}
