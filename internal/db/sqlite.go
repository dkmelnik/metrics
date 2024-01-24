package db

import (
	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
)

func NewSQLITEConnection(dbName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
