package metrics

import (
	"fmt"
	"strconv"

	"github.com/dkmelnik/metrics/internal/metrics/interfaces"
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/dkmelnik/metrics/internal/utils"
)

type Service struct {
	metricsRepo interfaces.MetricsRepository
}

func NewService(mr interfaces.MetricsRepository) *Service {
	return &Service{mr}
}

func (s *Service) RecordMetricValue(tp, nm, vl string) error {
	switch tp {
	case string(models.Counter):
		intVal, err := strconv.ParseInt(vl, 10, 64)
		if err != nil {
			return ErrParse
		}
		_, prev, err := s.metricsRepo.FindOneByTypeAndID(tp, nm)
		if err != nil {
			s.metricsRepo.SaveOrUpdate(models.Metric{
				ID:    nm,
				MType: tp,
				Delta: &intVal,
			})
		} else {
			if prev.Delta != nil {
				*prev.Delta += intVal
				s.metricsRepo.SaveOrUpdate(prev)
				return nil
			}
			return ErrParse
		}
		return nil
	case string(models.Gauge):
		flVal, err := strconv.ParseFloat(vl, 64)
		if err != nil {
			return ErrParse
		}
		s.metricsRepo.SaveOrUpdate(models.Metric{
			ID:    nm,
			MType: tp,
			Value: &flVal,
		})
		return nil

	default:
		return ErrTypeNotCorrect
	}
}

// ProcessMetricRequest В теле ответа отправляйте JSON той же структуры с актуальным (изменённым) значением Value.
func (s *Service) ProcessMetricRequest(dto models.Metric) error {
	switch dto.MType {
	case string(models.Gauge):
		s.metricsRepo.SaveOrUpdate(dto)
	case string(models.Counter):
		if _, prev, err := s.metricsRepo.FindOneByTypeAndID(dto.MType, dto.ID); err != nil {
			s.metricsRepo.SaveOrUpdate(dto)
			return nil
		} else {
			if prev.Delta != nil && dto.Delta != nil {
				*prev.Delta += *dto.Delta
				s.metricsRepo.SaveOrUpdate(prev)
				return nil
			}
			return ErrParse
		}
	default:
		return ErrTypeNotCorrect
	}

	return nil
}

// GetMetric теле ответа должен приходить такой же JSON, но с уже заполненными значениями метрик.
func (s *Service) GetMetric(metricType string, metricID string) (models.Metric, error) {
	_, metric, err := s.metricsRepo.FindOneByTypeAndID(metricType, metricID)
	if err != nil {
		return models.Metric{}, err
	}

	return metric, nil
}

func (s *Service) GetMetricValueString(tp, nm string) (string, error) {
	_, metric, err := s.metricsRepo.FindOneByTypeAndID(tp, nm)
	if err != nil {
		return "", err
	}

	if tp == string(models.Gauge) {
		return s.formatToString(tp, *metric.Value)
	} else {
		return s.formatToString(tp, *metric.Delta)
	}
}

func (s *Service) GetAllInHTML() string {
	metrics := s.metricsRepo.GetAllMetrics()

	html := "<html><head><title>Metric Values</title></head><body><h1>Metric Values:</h1><ul>"

	for _, metric := range metrics {
		html += "<li>"
		html += "<strong>" + metric.ID + ": </strong>"
		if metric.MType == "counter" {
			html += fmt.Sprintf("Value: %v", *metric.Delta)
		} else if metric.MType == "gauge" {
			html += fmt.Sprintf("Value: %v", utils.FormatFloat(*metric.Value, 3))
		}
		html += "</li>"
	}

	html += "</ul></body></html>"

	return html
}

func (s *Service) formatToString(tp string, vl interface{}) (string, error) {
	switch tp {
	case string(models.Gauge):
		flVal, err := strconv.ParseFloat(fmt.Sprintf("%v", vl), 64)
		if err != nil {
			return "", ErrParse
		}
		return utils.FormatFloat(flVal, 3), nil
	case string(models.Counter):
		iVal, err := strconv.ParseInt(fmt.Sprintf("%d", vl), 10, 64)
		if err != nil {
			return "", ErrParse
		}
		return fmt.Sprintf("%d", iVal), nil
	default:
		return fmt.Sprintf("%s", vl), nil
	}
}
