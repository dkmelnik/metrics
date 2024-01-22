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
	c := configs.NewServer()

	if err := logger.Setup(c, os.Stdout); err != nil {
		return err
	}

	if err := configs.CheckUnknownFlags(); err != nil {
		return err
	}

	// TODO: если обработать тесты упадут
	conn, _ := db.NewPsqlConnection(c)
	//conn = nil
	r, err := metrics.ConfigureRouter(conn, c)
	if err != nil {
		return err
	}

	s := server.NewServer(c.Addr, r)

	logger.Log.Info("SERVER LISTEN AND SERVE", "addr", c.Addr, "DBConnected", conn != nil)
	if err = s.Run(); err != nil {
		return err
	}

	return nil
}
