package collect

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/dkmelnik/metrics/internal/logger"
)

// Send periodically sends metrics to a server using the provided time ticker, metrics channel, server URL, and signer.
//
// The function continuously listens for incoming metrics from the provided channel and sends them to the specified server URL.
// It utilizes the provided context to allow for cancellation of the sending routine.
// Additionally, it expects a time.Ticker to determine the interval between metric sending attempts.
// The signer is used for signing the data before sending it to the server.
func Send(ctx context.Context, t *time.Ticker, ch <-chan *Metrics, serverURL string, signer Signer) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			loopMetricsAndSend(<-ch, serverURL, signer)
		}
	}
}

func loopMetricsAndSend(md *Metrics, serverURL string, signer Signer) {
	metricType := reflect.TypeOf(*md)
	metricValue := reflect.ValueOf(*md)

	workPayloads := make([]workPayload, 0)

	for i := 0; i < metricType.NumField(); i++ {
		field := metricType.Field(i)
		value := metricValue.Field(i)

		tag := field.Tag.Get("metric")
		if tag == "" {
			continue
		}
		body, err := buildCompressRequestBody(tag, field.Name, value.Interface())
		if err != nil {
			logger.Log.ErrorWithContext(context.Background(), err)
			continue
		}
		var hash string
		if signer != nil {
			hash = signer.HashData(body)
		}
		workPayloads = append(workPayloads, workPayload{url: serverURL, body: body, hash: hash})
	}

	limit := 5

	jobs := generator(workPayloads, limit)

	for w := 0; w <= limit; w++ {
		go worker(jobs)
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

func sendMetricRequest(url string, body []byte, hash string) {
	client := resty.New()

	header := map[string]string{
		"Content-Type":     "application/json",
		"Content-Encoding": "gzip",
	}

	if hash != "" {
		header["HashSHA256"] = hash
	}

	resp, err := client.R().
		SetHeaders(header).
		SetBody(body).
		Post(fmt.Sprintf("%s/update/", url))

	if err != nil {
		logger.Log.Error("sendMetricRequest", "err", err.Error(), "body", string(body))
		return
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Error(
			"sendMetricRequest",
			"err", "status not ok",
			"body", string(body),
			"resp", string(resp.Body()),
			"code", resp.StatusCode(),
		)
	}
}

func generator(input []workPayload, limit int) chan workPayload {
	out := make(chan workPayload, limit)
	go func() {
		defer close(out)
		for _, n := range input {
			out <- n
		}
	}()
	return out
}

type workPayload struct {
	url  string
	body []byte
	hash string
}

func worker(jobs <-chan workPayload) {
	for j := range jobs {
		sendMetricRequest(j.url, j.body, j.hash)
	}
}
