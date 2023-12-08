package metrics

import (
	"fmt"
	"github.com/dkmelnik/metrics/internal/models"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

func MakeRequest(url string) (string, error) {
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

func Send(md *models.Metrics) {
	for {
		time.Sleep(time.Second * 10)

		metricType := reflect.TypeOf(*md)
		metricValue := reflect.ValueOf(*md)

		for i := 0; i < metricType.NumField(); i++ {
			field := metricType.Field(i)
			value := metricValue.Field(i)

			tag := field.Tag.Get("metric")
			if tag != "" {
				reqUrl := fmt.Sprintf("http://server:8080/update/%s/%s/%s", tag, field.Name, convertToString(value))
				_, err := MakeRequest(reqUrl)
				if err != nil {
					log.Println(err)
				}
			}
		}

	}
}

func convertToString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Int:
		return strconv.Itoa(int(value.Int()))
	case reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	default:
		return ""
	}
}
