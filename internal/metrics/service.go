package metrics

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/dkmelnik/metrics/internal/metrics/dto"

	"github.com/dkmelnik/metrics/internal/apperrors"
	"github.com/dkmelnik/metrics/internal/models"
)

// Service represents the business logic layer.
type Service struct {
	metricsRepo IRepository
}

// NewService creates a new instance of Service.
func NewService(mr IRepository) *Service {
	return &Service{mr}
}

// CreateOrUpdateByParams is a method of Service to handle creating or updating metrics by params.
func (s *Service) CreateOrUpdateByParams(tp, nm, vl string) error {
	metric := models.Metric{
		Name:  nm,
		MType: tp,
	}
	switch tp {
	case string(models.Counter):
		iVal, err := strconv.ParseInt(vl, 10, 64)
		if err != nil {
			return apperrors.ErrParse
		}
		metric.SetDelta(iVal)
	case string(models.Gauge):
		iVal, err := strconv.ParseFloat(vl, 64)
		if err != nil {
			return apperrors.ErrParse
		}
		metric.SetValue(iVal)
	default:
		return apperrors.ErrTypeNotCorrect
	}

	return s.CreateOrUpdate(metric)
}

// CreateOrUpdate is a method of Service to handle creating or updating metrics by dto models.Metric.
func (s *Service) CreateOrUpdate(dto models.Metric) error {
	ctx := context.Background()

	if string(models.Counter) == dto.MType {
		prev, err := s.metricsRepo.FindOneByTypeAndName(ctx, dto.MType, dto.Name)
		if nil != err && !errors.Is(err, apperrors.ErrNotFound) {
			return err
		}
		dto.UpdateDelta(prev.Delta.Int64)
	}
	return s.metricsRepo.SaveOrUpdate(ctx, dto)
}

// CreateOrUpdateMany is a method of Service to handle creating or updating slice of metrics.
func (s *Service) CreateOrUpdateMany(dtos []models.Metric) error {
	for _, dto := range dtos {
		if err := s.CreateOrUpdate(dto); err != nil {
			return err
		}
	}
	return nil
}

// GetMetric is a method of Service to handle getting metric details dto.Details.
func (s *Service) GetMetric(tp, nm string) (dto.Details, error) {
	ctx := context.Background()

	m, err := s.metricsRepo.FindOneByTypeAndName(ctx, tp, nm)
	if err != nil {
		return dto.Details{}, err
	}
	var out dto.Details
	out.FillFromModel(m)
	return out, nil
}

// GetMetricValue is a method of Service to handle getting metric value by type and name.
func (s *Service) GetMetricValue(tp, nm string) (interface{}, error) {
	ctx := context.Background()

	m, err := s.metricsRepo.FindOneByTypeAndName(ctx, tp, nm)
	if err != nil {
		return nil, err
	}
	return m.GetValueByType(), nil
}

// GetAllInHTML is a method of Service to handle getting metrics in html format.
func (s *Service) GetAllInHTML() (string, error) {
	ctx := context.Background()

	metrics, err := s.metricsRepo.Find(ctx)
	if err != nil {
		return "", err
	}
	html := "<html><head><title>Metric Values</title></head><body><h1>Metric Values:</h1><ul>"

	for _, metric := range metrics {
		html += "<li>"
		html += "<strong>" + metric.Name + ": </strong>"
		html += fmt.Sprintf("Guid: %v\t", metric.ID)
		html += fmt.Sprintf("Value: %v", metric.GetValueByType())
		html += "</li>"
	}

	html += "</ul></body></html>"

	return html, nil
}
