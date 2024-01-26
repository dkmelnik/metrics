package db

import (
	"github.com/dkmelnik/metrics/configs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPsqlConnection(c configs.Server) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", c.DBConnectStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
