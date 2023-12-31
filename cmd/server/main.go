package main

import (
	"log"

	"github.com/dkmelnik/metrics/configs"
	"github.com/dkmelnik/metrics/internal/metrics"
	"github.com/dkmelnik/metrics/internal/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("server is running!")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := configs.CheckUnknownFlags(); err != nil {
		return err
	}

	c := configs.NewServer().Build()

	r := metrics.ConfigureRouter()

	s := server.NewServer(c.Addr, r)

	if err := s.Run(); err != nil {
		return err
	}

	return nil
}
