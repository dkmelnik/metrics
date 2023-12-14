package metrics

import (
	"fmt"
	"math"
	"strconv"

	"github.com/dkmelnik/metrics/internal/metrics/interfaces"
	"github.com/dkmelnik/metrics/internal/models"
)

type Service struct {
	metricsRepo interfaces.MetricsRepository
}

func NewService(mr interfaces.MetricsRepository) *Service {
	return &Service{mr}
}

func (s *Service) SaveMetricData(tp, nm, vl string) error {
	switch tp {
	case string(models.Counter):
		intVal, err := strconv.ParseInt(vl, 10, 64)
		if err != nil {
			return ErrParse
		}
		prev, err := s.metricsRepo.FindOneByTypeName(tp, nm)
		if err != nil {
			s.metricsRepo.Save(tp, nm, intVal)
		} else {
			i := prev.(int64)
			s.metricsRepo.Save(tp, nm, intVal+i)
		}
		return nil
	case string(models.Gauge):
		flVal, err := strconv.ParseFloat(vl, 64)
		if err != nil {
			return ErrParse
		}
		s.metricsRepo.Save(tp, nm, flVal)
		return nil

	default:
		return ErrTypeNotCorrect
	}
}

func (s *Service) GetMetricData(tp, nm string) (string, error) {
	metric, err := s.metricsRepo.FindOneByTypeName(tp, nm)
	if err != nil {
		return "", err
	}

	var out string
	if tp == string(models.Gauge) {
		flVal, err := strconv.ParseFloat(fmt.Sprintf("%v", metric), 64)
		if err != nil {
			return "", ErrParse
		}
		out = fmt.Sprintf("%v", math.Round(flVal*1000)/1000)
	} else {
		out = fmt.Sprintf("%v", metric)
	}

	return out, nil
}

func (s *Service) GetAllInHTML() string {
	metrics := s.metricsRepo.GetAllMetrics()

	html := "<html><head><title>Metric Values</title></head><body><h1>Metric Values:</h1><ul>"

	for metricName, values := range metrics {
		html += fmt.Sprintf("<li><strong>%s</strong>: <ul>", metricName)
		for key, value := range values {
			html += fmt.Sprintf("<li>%s: %v</li>", key, value)
		}
		html += "</ul></li>"
	}

	html += "</ul></body></html>"

	return html
}
