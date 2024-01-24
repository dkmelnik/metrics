package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dkmelnik/metrics/internal/apperrors"
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/jmoiron/sqlx"
	"time"
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

	uq := `UPDATE metrics SET delta = :delta, value = :value, updated_at = :updated_at WHERE type = :type AND name = :name`
	metric.UpdatedAT = time.Now()
	_, err = r.db.NamedExecContext(ctx, uq, metric)

	return err
}

func (r *RepositoryStorage) SaveOrUpdateMany(metrics []models.Metric) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	ctx := context.Background()

	for _, m := range metrics {
		var existingData models.Metric

		sq := `SELECT * FROM metrics WHERE type = $1 AND name = $2 LIMIT 1`
		err := r.db.GetContext(ctx, &existingData, sq, m.MType, m.Name)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		if errors.Is(err, sql.ErrNoRows) {
			iq := `INSERT INTO metrics (name, type, delta, value) VALUES ($1, $2, $3, $4)`
			_, err := tx.ExecContext(ctx, iq, m.Name, m.MType, m.Delta, m.Value)

			if err != nil {
				return err
			}
		} else {
			uq := `UPDATE metrics SET delta = $1, value = $2, updated_at = $3 WHERE type = $4 AND name = $5`
			m.UpdatedAT = time.Now()
			_, err := tx.ExecContext(ctx, uq, m.Delta, m.Value, m.UpdatedAT, m.MType, m.Name)

			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
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
