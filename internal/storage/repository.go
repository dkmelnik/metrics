package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dkmelnik/metrics/internal/apperrors"
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/jmoiron/sqlx"
)

type RepositoryStorage struct {
	db *sqlx.DB
}

func NewRepositoryStorage(db *sqlx.DB) (*RepositoryStorage, error) {
	if db == nil {
		return nil, errors.New("db not instance")
	}
	return &RepositoryStorage{db}, nil
}

func (r *RepositoryStorage) SaveOrUpdate(metric models.Metric) error {
	var existingData models.Metric
	ctx := context.Background()

	sq := `SELECT * FROM metrics WHERE type = $1 AND name = $2 LIMIT 1`
	err := r.db.GetContext(ctx, &existingData, sq, metric.MType, metric.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		iq := `INSERT INTO metrics (name, type, delta, value) VALUES (:name, :type, :delta, :value)`
		_, err = r.db.NamedExecContext(ctx, iq, metric)
		return err
	}

	uq := `UPDATE metrics SET delta = :delta, value = :value WHERE type = :type AND name = :name`
	_, err = r.db.NamedExecContext(ctx, uq, metric)

	return err
}

func (r *RepositoryStorage) FindOneByTypeAndName(mType, mName string) (models.Metric, error) {
	var existingData models.Metric
	ctx := context.Background()
	q := `SELECT * FROM metrics WHERE type = $1 AND name = $2 LIMIT 1`
	err := r.db.GetContext(ctx, &existingData, q, mType, mName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return existingData, apperrors.ErrNotFound
		}
		return existingData, err
	}

	return existingData, nil
}

func (r *RepositoryStorage) Find() ([]models.Metric, error) {
	var existingData []models.Metric
	ctx := context.Background()
	err := r.db.SelectContext(ctx, &existingData, "SELECT * FROM metrics")

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return existingData, nil
}
