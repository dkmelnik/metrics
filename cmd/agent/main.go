package main

import (
	"context"
	"log"
	"time"

	"github.com/dkmelnik/metrics/configs"
	"github.com/dkmelnik/metrics/internal/collect"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("agent is running!")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
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

	done := make(chan struct{})
	<-done

	return nil
}
