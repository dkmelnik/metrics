package main

import (
	"context"
	"os"
	"time"

	"github.com/dkmelnik/metrics/configs"
	"github.com/dkmelnik/metrics/internal/collect"
	"github.com/dkmelnik/metrics/internal/logger"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Setup(configs.NewLogger(), os.Stdout); err != nil {
		return err
	}

	if err := configs.CheckUnknownFlags(); err != nil {
		return err
	}

	c := configs.NewAgent()

	metricsChan := make(chan *collect.Metrics)

	// CTX for stopping sender and collector
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// COLLECTOR
	collectPeriod := time.NewTicker(time.Second * time.Duration(c.PollInterval))
	defer collectPeriod.Stop()
	go collect.Collect(ctx, collectPeriod, metricsChan)

	// SENDER
	sendPeriod := time.NewTicker(time.Second * time.Duration(c.ReportInterval))
	defer sendPeriod.Stop()
	go collect.Send(ctx, sendPeriod, metricsChan, c.Addr)

	logger.Log.Info("AGENT RUNNING", "ReportInterval", c.ReportInterval, "PollInterval", c.PollInterval)

	done := make(chan struct{})
	<-done

	return nil
}
