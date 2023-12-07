package storage

import "github.com/dkmelnik/metrics/internal/models"

type Collection []models.Metrics

func NewCollection() *Collection {
	return &Collection{}
}

func (c *Collection) Save(m models.Metrics) {
	*c = append(*c, m)
}

func (c *Collection) Find() []models.Metrics {
	return *c
}

func (c *Collection) Clear() {
	*c = []models.Metrics{}
}
