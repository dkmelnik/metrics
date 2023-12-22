package main

import (
	"context"
	"log"
	"time"

	"github.com/dkmelnik/metrics/configs"
	"github.com/dkmelnik/metrics/internal/collect"
	"github.com/dkmelnik/metrics/internal/models"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("agent is running!")
	return
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := configs.CheckUnknownFlags(); err != nil {
		return err
	}

	c := configs.NewAgent().Build()

	md := &models.Metrics{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collectPeriod := time.NewTicker(time.Second * time.Duration(c.PollInterval))
	defer collectPeriod.Stop()
	go collect.Collect(ctx, collectPeriod, md)

	sendPeriod := time.NewTicker(time.Second * time.Duration(c.ReportInterval))
	defer sendPeriod.Stop()
	go collect.Send(ctx, sendPeriod, md, c.Addr)

	done := make(chan struct{})
	<-done

	return nil
}
