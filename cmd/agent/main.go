package main

import (
	"context"
	"github.com/dkmelnik/metrics/internal/metrics"
	"github.com/dkmelnik/metrics/internal/models"
	"log"
	"time"
)

func main() {
	log.Println("agent is running!")

	md := &models.Metrics{}

	collectPeriod := time.NewTicker(time.Second * 2)
	defer collectPeriod.Stop()
	go metrics.Collect(context.Background(), collectPeriod, md)

	sendPeriod := time.NewTicker(time.Second * 10)
	defer sendPeriod.Stop()
	go metrics.Send(context.Background(), sendPeriod, md, "http://server:8080")

	done := make(chan struct{})
	<-done

}
