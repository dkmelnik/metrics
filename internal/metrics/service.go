package metrics

import (
	"fmt"
	"github.com/dkmelnik/metrics/internal/metrics/interfaces"
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/dkmelnik/metrics/internal/utils"
	"strconv"
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

	return s.formatToString(tp, metric)
}

func (s *Service) GetAllInHTML() string {
	metrics := s.metricsRepo.GetAllMetrics()

	html := "<html><head><title>Metric Values</title></head><body><h1>Metric Values:</h1><ul>"

	for metricName, values := range metrics {
		html += fmt.Sprintf("<li><strong>%s</strong>: <ul>", metricName)
		for key, value := range values {
			vl, _ := s.formatToString(metricName, value)
			html += fmt.Sprintf("<li>%s: %s</li>", key, vl)
		}
		html += "</ul></li>"
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
