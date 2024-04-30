package main

import (
	"database/sql"
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/dkmelnik/metrics/configs"
	"github.com/dkmelnik/metrics/internal/db"
	"github.com/dkmelnik/metrics/internal/logger"
	"github.com/dkmelnik/metrics/internal/metrics"
	"github.com/dkmelnik/metrics/internal/server"
	"github.com/dkmelnik/metrics/internal/sign"
	"github.com/dkmelnik/metrics/internal/storage"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
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
	connPG, _ := db.NewPsqlConnection(c)
	if connPG != nil {
		defer connPG.Close()
		if err := migrateDatabase(
			"postgres",
			c.DBConnectStr,
			"file://migrations/pg",
		); err != nil {
			return err
		}
	}

	dbName := "file:db.sqlite3?cache=shared"
	connSQLITE, _ := db.NewSQLITEConnection(dbName)
	if connSQLITE != nil {
		defer connSQLITE.Close()
		if err := migrateDatabase(
			"sqlite3",
			dbName,
			"file://migrations/sqlite",
		); err != nil {
			return err
		}
	}

	var store metrics.IRepository
	if connPG != nil {
		store, _ = storage.NewRepositoryStorage(connPG)
	} else if connSQLITE != nil {
		store, _ = storage.NewMemoryStorage(c.FileStoragePath, c.StoreInterval, c.Restore)
	}

	signer := sign.NewSign(c.Key)

	r, err := metrics.ConfigureRouter(c, connPG, store, signer)
	if err != nil {
		return err
	}

	s := server.NewServer(c.Addr, r)

	logger.Log.Info("SERVER LISTEN AND SERVE",
		"addr", c.Addr,
		"DBConnected", connPG != nil,
		"buildVersion", buildVersion,
		"buildDate", buildDate,
		"buildCommit", buildCommit,
	)
	if err = s.Run(); err != nil {
		return err
	}

	return nil
}

func migrateDatabase(driverName, connStr, migrationPath string) error {
	dbinst, err := sql.Open(driverName, connStr)
	if err != nil {
		return err
	}
	defer dbinst.Close()

	var driver database.Driver
	switch driverName {
	case "postgres":
		driver, err = postgres.WithInstance(dbinst, &postgres.Config{})
	case "sqlite3":
		driver, err = sqlite3.WithInstance(dbinst, &sqlite3.Config{})
	default:
		return errors.New("unsupported database driver")
	}

	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"db", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
