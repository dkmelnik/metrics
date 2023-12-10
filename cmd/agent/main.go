package main

import (
	"context"
	"errors"
	"flag"
	"github.com/dkmelnik/metrics/internal/metrics"
	"github.com/dkmelnik/metrics/internal/models"
	"log"
	"time"
)

func main() {
	log.Println("agent is running!")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	addr                         string
	reportInterval, pollInterval int
}

func flags() (Config, error) {
	c := Config{}
	flag.StringVar(&c.addr, "a", "server:8080", "server by collected metric address ")

	flag.IntVar(&c.reportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.IntVar(&c.pollInterval, "p", 2, "metrics polling frequency")

	flag.Parse()

	if len(flag.Args()) > 0 {
		return c, errors.New("unknown flags or parameters")
	}
	return c, nil
}

func run() error {
	c, err := flags()
	if err != nil {
		return err
	}

	md := &models.Metrics{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collectPeriod := time.NewTicker(time.Second * time.Duration(c.pollInterval))
	defer collectPeriod.Stop()
	go metrics.Collect(ctx, collectPeriod, md)

	sendPeriod := time.NewTicker(time.Second * time.Duration(c.reportInterval))
	defer sendPeriod.Stop()
	go metrics.Send(ctx, sendPeriod, md, c.addr)

	done := make(chan struct{})
	<-done

	return nil
}
