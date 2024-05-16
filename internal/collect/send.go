package collect

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/dkmelnik/metrics/internal/logger"
)

// MetricsCollector представляет собой коллектор метрик.
type MetricsCollector struct {
	ctx       context.Context
	period    *time.Ticker
	payloadCh <-chan *Metrics
	publicKey *rsa.PublicKey
	serverURL string
	signer    Signer
	limit     int
}

type workPayload struct {
	url  string
	body []byte
	hash string
}

func NewMetricsCollector(
	ctx context.Context,
	period *time.Ticker,
	publicKeyPath string,
	payloadCh <-chan *Metrics,
	serverURL string,
	signer Signer,
	limit int,
) (*MetricsCollector, error) {

	mc := &MetricsCollector{
		ctx:       ctx,
		period:    period,
		payloadCh: payloadCh,
		serverURL: serverURL,
		signer:    signer,
		limit:     limit,
	}

	mc.loadPublicKey(publicKeyPath)

	return mc, nil
}

// SendMetricsPeriodically периодически отправляет метрики на сервер.
func (mc *MetricsCollector) SendMetricsPeriodically() {
	go func() {

		for {
			select {
			case <-mc.ctx.Done():
				return
			case <-mc.period.C:
				mc.sendMetrics()
			}
		}
	}()
}

func (mc *MetricsCollector) sendMetrics() {

	md := <-mc.payloadCh
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
			continue
		}

		if mc.publicKey != nil {
			body, err = mc.encrypt(body)
			if err != nil {
				continue
			}
		}

		var hash string
		if mc.signer != nil {
			hash = mc.signer.HashData(body)
		}

		workPayloads = append(workPayloads, workPayload{url: mc.serverURL, body: body, hash: hash})
	}

	mc.processWorkPayloads(workPayloads)
}

func (mc *MetricsCollector) processWorkPayloads(payloads []workPayload) {
	jobs := mc.generator(payloads)

	for w := 0; w < mc.limit; w++ {
		go mc.worker(jobs)
	}
}

func (mc *MetricsCollector) generator(input []workPayload) chan workPayload {
	out := make(chan workPayload, mc.limit)
	go func() {
		defer close(out)
		for _, n := range input {
			out <- n
		}
	}()
	return out
}

func (mc *MetricsCollector) worker(jobs <-chan workPayload) {
	for j := range jobs {
		mc.sendMetricRequest(j.url, j.body, j.hash)
	}
}

func (mc *MetricsCollector) sendMetricRequest(url string, body []byte, hash string) {
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

func (mc *MetricsCollector) loadPublicKey(path string) error {
	keyFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read public key file: %v", err)
	}

	block, _ := pem.Decode(keyFile)
	if block == nil {
		return errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %v", err)
	}

	mc.publicKey = pub

	return nil
}

func (mc *MetricsCollector) encrypt(data []byte) ([]byte, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, mc.publicKey, data)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
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
