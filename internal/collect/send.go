package collect

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/go-resty/resty/v2"
)

func Send(ctx context.Context, t *time.Ticker, ch <-chan *Metrics, serverURL string) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			loopMetricsAndSend(<-ch, serverURL)
		}
	}
}

func loopMetricsAndSend(md *Metrics, serverURL string) {
	metricType := reflect.TypeOf(*md)
	metricValue := reflect.ValueOf(*md)

	for i := 0; i < metricType.NumField(); i++ {
		field := metricType.Field(i)
		value := metricValue.Field(i)

		tag := field.Tag.Get("metric")
		if tag != "" {
			body := buildRequestBody(tag, field.Name, value.Interface())
			sendMetricRequest(serverURL, body)
		}
	}
}

func buildRequestBody(mt string, mn string, vl interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	m["id"] = mn
	m["type"] = mt

	switch mt {
	case "gauge":
		if v, ok := vl.(float64); ok {
			m["value"] = v
		}
	case "counter":
		if v, ok := vl.(int); ok {
			m["delta"] = v
		}
	}

	return m
}

func sendMetricRequest(url string, body map[string]interface{}) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("%s/update", url))

	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode())
	}
}
