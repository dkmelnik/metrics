package metrics

import (
	"errors"
	"fmt"
	"github.com/dkmelnik/metrics/internal/logger"
	"strconv"

	"github.com/dkmelnik/metrics/internal/apperrors"
	"github.com/dkmelnik/metrics/internal/metrics/interfaces"
	"github.com/dkmelnik/metrics/internal/models"
)

type Service struct {
	metricsRepo interfaces.MetricsRepository
}

func NewService(mr interfaces.MetricsRepository) *Service {
	return &Service{mr}
}

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

func (s *Service) CreateOrUpdate(dto models.Metric) error {
	if string(models.Counter) == dto.MType {
		prev, err := s.metricsRepo.FindOneByTypeAndName(dto.MType, dto.Name)
		if nil != err && !errors.Is(err, apperrors.ErrNotFound) {
			logger.Log.Error("Error while getting metric", "error", err)
			return err
		}
		dto.UpdateDelta(prev.Delta.Int64)
	}
	return s.metricsRepo.SaveOrUpdate(dto)
}

func (s *Service) GetMetric(tp, nm string) (models.Metric, error) {
	return s.metricsRepo.FindOneByTypeAndName(tp, nm)
}

func (s *Service) GetMetricValue(tp, nm string) (interface{}, error) {
	m, err := s.GetMetric(tp, nm)
	if err != nil {
		return nil, err
	}
	return m.GetValue(), nil
}

func (s *Service) GetAllInHTML() (string, error) {
	metrics, err := s.metricsRepo.Find()
	if err != nil {
		return "", err
	}
	html := "<html><head><title>Metric Values</title></head><body><h1>Metric Values:</h1><ul>"

	for _, metric := range metrics {
		html += "<li>"
		html += "<strong>" + metric.Name + ": </strong>"
		html += fmt.Sprintf("Guid: %v\t", metric.ID)
		html += fmt.Sprintf("Value: %v", metric.GetValue())
		html += "</li>"
	}

	html += "</ul></body></html>"

	return html, nil
}
