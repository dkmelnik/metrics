package configs

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

type Agent struct {
	Addr                         string
	ReportInterval, PollInterval int
}

func NewAgent() Agent {
	return Agent{}
}

func (cb Agent) Build() Agent {
	flag.StringVar(&cb.Addr, "a", "http://localhost:8080", "server by collected metric address ")
	flag.IntVar(&cb.ReportInterval, "r", 10, "frequency of sending collect to the server")
	flag.IntVar(&cb.PollInterval, "p", 2, "frequency of sending collect to the server")
	flag.Parse()

	s, ok := os.LookupEnv("ADDRESS")
	if ok {
		cb.Addr = s
	}

	s, ok = os.LookupEnv("REPORT_INTERVAL")
	if ok {
		i, err := strconv.Atoi(s)
		if err == nil {
			cb.ReportInterval = i
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
