package main

import (
	"github.com/dkmelnik/metrics/internal/metrics"
	"github.com/dkmelnik/metrics/internal/models"
	"log"
)

func main() {
	log.Println("agent is running!")

	md := &models.Metrics{}

	go metrics.Collect(md)
	go metrics.Send(md)

	done := make(chan struct{})
	<-done

}
