package main

import (
	"github.com/dkmelnik/metrics/configs"
	"github.com/dkmelnik/metrics/internal/logger"
	"github.com/dkmelnik/metrics/internal/metrics"
	"github.com/dkmelnik/metrics/internal/server"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := logger.Setup(configs.NewLogger(), os.Stdout); err != nil {
		return err
	}

	if err := configs.CheckUnknownFlags(); err != nil {
		return err
	}

	c := configs.NewServer()

	r, err := metrics.ConfigureRouter(configs.NewStorage())
	if err != nil {
		return err
	}

	s := server.NewServer(c.Addr, r)

	logger.Log.Info("LISTEN AND SERVE", "addr", c.Addr)
	if err = s.Run(); err != nil {
		return err
	}

	return nil
}
