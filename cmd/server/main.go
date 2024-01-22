package main

import (
	"database/sql"
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/dkmelnik/metrics/configs"
	"github.com/dkmelnik/metrics/internal/db"
	"github.com/dkmelnik/metrics/internal/logger"
	"github.com/dkmelnik/metrics/internal/metrics"
	"github.com/dkmelnik/metrics/internal/server"
)

func main() {
	if err := run(); err != nil {
		logger.Log.Error("Error starting app", "error", err)
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

	if conn != nil {
		if err := migrateDB(c.DBConnectStr); err != nil {
			return err
		}
	}

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

func migrateDB(connStr string) error {

	dbinst, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer dbinst.Close()

	driver, err := postgres.WithInstance(dbinst, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
