package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dkmelnik/metrics/configs"
	"github.com/dkmelnik/metrics/internal/collect"
	"github.com/dkmelnik/metrics/internal/logger"
	"github.com/dkmelnik/metrics/internal/sign"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	if err := run(); err != nil {
		logger.Log.Error("Error starting app", "error", err)
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	ctx, cancel := context.WithCancel(ctx)
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

	cl, err := collect.NewMetricsCollector(ctx, sendPeriod, c.PublicKeyPath, metricsChan, c.Addr, signer, 5)
	if err != nil {
		return err
	}
	cl.SendMetricsPeriodically()

	logger.Log.Info("AGENT RUNNING",
		"ReportInterval", c.ReportInterval,
		"PollInterval", c.PollInterval,
		"buildVersion", buildVersion,
		"buildDate", buildDate,
		"buildCommit", buildCommit,
	)

	<-ctx.Done()

	return nil
}
