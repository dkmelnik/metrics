package collect

import (
	"context"
	"fmt"
	"github.com/dkmelnik/metrics/internal/models"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"
)

func Send(ctx context.Context, t *time.Ticker, md *models.Metrics, serverURL string) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if md == nil {
				continue
			}
			loopMetricsAndSend(md, serverURL)
		}
	}
}

func loopMetricsAndSend(md *models.Metrics, serverURL string) {
	metricType := reflect.TypeOf(*md)
	metricValue := reflect.ValueOf(*md)

	for i := 0; i < metricType.NumField(); i++ {
		field := metricType.Field(i)
		value := metricValue.Field(i)

		tag := field.Tag.Get("metric")
		if tag != "" {
			reqURL := buildRequestURL(serverURL, tag, field.Name, convertToString(value))
			_, err := makeRequest(reqURL)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func makeRequest(url string) (string, error) {
	client := http.Client{
		Timeout: 40 * time.Second,
	}

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func buildRequestURL(serverURL, tag, fieldName, value string) string {
	return fmt.Sprintf("%s/update/%s/%s/%s", serverURL, tag, fieldName, value)
}

func convertToString(value reflect.Value) string {
	return fmt.Sprintf("%v", value.Interface())
}
