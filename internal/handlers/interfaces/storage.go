package interfaces

import "github.com/dkmelnik/metrics/internal/models"

type Storage interface {
	Save(m models.Metrics)
	Find() []models.Metrics
}
