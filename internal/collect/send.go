package collect

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
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
			body, err := buildCompressRequestBody(tag, field.Name, value.Interface())
			if err != nil {
				log.Println(err)
				continue
			}
			sendMetricRequest(serverURL, body)
		}
	}
}

func buildCompressRequestBody(mt string, mn string, vl interface{}) ([]byte, error) {
	metric := make(map[string]interface{})

	metric["id"] = mn
	metric["type"] = mt

	switch mt {
	case "gauge":
		if v, ok := vl.(float64); ok {
			metric["value"] = v
		}
	case "counter":
		if v, ok := vl.(int); ok {
			metric["delta"] = v
		}
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	mtb, err := json.Marshal(metric)
	if err != nil {
		return nil, fmt.Errorf("failed marshal data to bytes: %v", err)
	}

	_, err = w.Write(mtb)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}

func sendMetricRequest(url string, body []byte) {
	client := resty.New()

	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type":     "application/json",
			"Content-Encoding": "gzip",
		}).
		SetBody(body).
		Post(fmt.Sprintf("%s/update/", url))

	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode())
	}
}
