package main

import (
	"os"

	"github.com/dkmelnik/metrics/configs"
	"github.com/dkmelnik/metrics/internal/db"
	"github.com/dkmelnik/metrics/internal/logger"
	"github.com/dkmelnik/metrics/internal/metrics"
	"github.com/dkmelnik/metrics/internal/server"
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

	c := configs.NewServer()

	conn, err := db.NewPsqlConnection(c)
	if err != nil {
		return err
	}

	r, err := metrics.ConfigureRouter(conn, configs.NewStorage())
	if err != nil {
		return err
	}

	s := server.NewServer(c.Addr, r)

	logger.Log.Info("SERVER LISTEN AND SERVE", "addr", c.Addr, "db", "connected")
	if err = s.Run(); err != nil {
		return err
	}

	return nil
}
