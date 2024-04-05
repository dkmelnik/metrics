package main

import (
	"context"
	"os"
	"time"

	"github.com/dkmelnik/metrics/configs"
	"github.com/dkmelnik/metrics/internal/collect"
	"github.com/dkmelnik/metrics/internal/logger"
	"github.com/dkmelnik/metrics/internal/sign"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	c := configs.NewAgent()

	if err := logger.Setup(c, os.Stdout); err != nil {
		return err
	}

	if err := configs.CheckUnknownFlags(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collectPeriod := time.NewTicker(time.Second * time.Duration(c.PollInterval))
	defer collectPeriod.Stop()

	metricsChan := collect.MetricsGenerator(ctx, collectPeriod)

	// SENDER
	sendPeriod := time.NewTicker(time.Second * time.Duration(c.ReportInterval))
	defer sendPeriod.Stop()

	var signer collect.Signer
	if c.Key != "" {
		signer = sign.NewSign(c.Key)
	}

	go collect.Send(ctx, sendPeriod, metricsChan, c.Addr, signer)

	logger.Log.Info("AGENT RUNNING", "ReportInterval", c.ReportInterval, "PollInterval", c.PollInterval)

	done := make(chan struct{})
	<-done

	return nil
}
