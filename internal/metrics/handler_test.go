package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dkmelnik/metrics/internal/metrics/mock"
	"github.com/dkmelnik/metrics/internal/storage"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, body map[string]interface{}) (*http.Response, string) {
	bts, _ := json.Marshal(body)

	req, err := http.NewRequest(method, ts.URL+path, bytes.NewReader(bts))
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func Test_CreateOrUpdateByParams(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "negative test #1, method not allowed",
			request: "/update/gauge/Alloc/242288",
			method:  http.MethodGet,
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "negative test #2, method does not match pattern update/type/name/value",
			request: "/update/gauge/Alloc",
			method:  http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				response:    "404 page not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "negative test #3, value not int64 or float64",
			request: "/update/gauge/Alloc/test",
			method:  http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				response:    "Bad Request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "positive test #4, must return OK",
			request: "/update/gauge/Alloc/" + fmt.Sprint(rand.Intn(1024)),
			method:  http.MethodPost,
			want: want{
				code:        http.StatusOK,
				response:    http.StatusText(http.StatusOK),
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	store, err := storage.NewMemoryStorage("/tmp/metrics-db.json", 10, false)
	if err != nil {
		t.Error(err)
	}

	r, err := ConfigureRouter("", "", nil, store, nil)
	if err != nil {
		t.Error(err)
	}
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, tt.method, tt.request, nil)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, get)
		})
	}
}

func Test_CreateOrUpdateByJSON(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		body    map[string]interface{}
		method  string
		want    want
		wantErr bool
	}{
		{
			name:    "negative test #1, empty body",
			body:    nil,
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusBadRequest,
				response:    http.StatusText(http.StatusBadRequest) + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #2, empty type",
			body: map[string]interface{}{
				"id":    "testCounter",
				"delta": 1,
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusBadRequest,
				response:    http.StatusText(http.StatusBadRequest) + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #3, incorrect type",
			body: map[string]interface{}{
				"id":    "testCounter",
				"type":  "none",
				"delta": 1,
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusBadRequest,
				response:    http.StatusText(http.StatusBadRequest) + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #4, counter type incorrect delta",
			body: map[string]interface{}{
				"id":    "testCounter",
				"type":  "counter",
				"delta": "none",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusBadRequest,
				response:    http.StatusText(http.StatusBadRequest) + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #4, gauge type incorrect value",
			body: map[string]interface{}{
				"id":    "LastGC",
				"type":  "gauge",
				"value": "none",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusBadRequest,
				response:    http.StatusText(http.StatusBadRequest) + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		}, {
			name: "negative test #5, gauge type empty value",
			body: map[string]interface{}{
				"id":   "LastGC",
				"type": "gauge",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusBadRequest,
				response:    http.StatusText(http.StatusBadRequest) + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		}, {
			name: "negative test #6, counter type empty delta",
			body: map[string]interface{}{
				"id":   "TestCounter",
				"type": "counter",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusBadRequest,
				response:    http.StatusText(http.StatusBadRequest) + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		}, {
			name: "positive test 7, gauge type return json",
			body: map[string]interface{}{
				"id":    "LastGC",
				"type":  "gauge",
				"value": 10.123123,
			},
			method:  http.MethodPost,
			wantErr: false,
			want: want{
				code: http.StatusOK,
				response: `{
					"id": "LastGC",
					"type": "gauge",
					"value": 10.123123
				}`,
				contentType: "application/json",
			},
		}, {
			name: "positive test #8, counter type return json",
			body: map[string]interface{}{
				"id":    "TestCounter",
				"type":  "counter",
				"delta": 1,
			},
			method:  http.MethodPost,
			wantErr: false,
			want: want{
				code: http.StatusOK,
				response: `{
					"id": "TestCounter",
					"type": "counter",
					"delta": 1
				}`,
				contentType: "application/json",
			},
		},
	}
	store, err := storage.NewMemoryStorage("/tmp/metrics-db.json", 10, false)
	if err != nil {
		t.Error(err)
	}

	r, err := ConfigureRouter("", "", nil, store, nil)
	if err != nil {
		t.Error(err)
	}
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, tt.method, "/update/", tt.body)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			if tt.wantErr {
				assert.Equal(t, tt.want.response, get)
			} else {
				assert.JSONEq(t, tt.want.response, get)
			}
		})
	}
}

func Test_GetMetricValue(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name        string
		metricsType string
		metricsName string
		want        want
	}{
		{
			name:        "positive test #1, valid gauge HeapAlloc metric",
			metricsType: "gauge",
			metricsName: "HeapAlloc",
			want: want{
				code:        http.StatusOK,
				response:    "150112",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "positive test #2, valid gauge HeapSys metric",
			metricsType: "gauge",
			metricsName: "HeapSys",
			want: want{
				code:        http.StatusOK,
				response:    "3.833856e+06",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "positive test #3, valid gauge MCacheSys metric",
			metricsType: "gauge",
			metricsName: "MCacheSys",
			want: want{
				code:        http.StatusOK,
				response:    "13",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "positive test #4, valid gauge TotalAlloc metric",
			metricsType: "gauge",
			metricsName: "TotalAlloc",
			want: want{
				code:        http.StatusOK,
				response:    "7.70766",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "positive test #5, valid gauge Mallocs metric",
			metricsType: "gauge",
			metricsName: "Mallocs",
			want: want{
				code:        http.StatusOK,
				response:    "282",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "positive test #6, valid gauge OtherSys metric",
			metricsType: "gauge",
			metricsName: "OtherSys",
			want: want{
				code:        http.StatusOK,
				response:    "3.4857341e+06",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "positive test #7, valid gauge NextGC metric",
			metricsType: "gauge",
			metricsName: "NextGC",
			want: want{
				code:        http.StatusOK,
				response:    "-3.35872e+06",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "positive test #8, valid gauge LastGC metric",
			metricsType: "gauge",
			metricsName: "LastGC",
			want: want{
				code:        http.StatusOK,
				response:    "0",
				contentType: "text/plain; charset=utf-8",
			},
		}, {
			name:        "positive test #9, valid counter PollCount metric",
			metricsType: "counter",
			metricsName: "PollCount",
			want: want{
				code:        http.StatusOK,
				response:    "14123413542",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	r := chi.NewRouter()

	st, err := mock.NewStorageMock()
	if err != nil {
		t.Error(err)
	}
	sr := NewService(st)
	h := NewHandler(nil, sr)

	r.Get("/value/{type}/{name}", h.GetMetricValue)

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, get := testRequest(t, ts, http.MethodGet, "/value"+"/"+tt.metricsType+"/"+tt.metricsName, nil)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, get)
		})
	}
}

func Test_GetMetric(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		body    map[string]interface{}
		method  string
		want    want
		wantErr bool
	}{
		{
			name: "negative test #1, incorrect body",
			body: map[string]interface{}{
				"id":   1223,
				"type": 1435,
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusBadRequest,
				response:    http.StatusText(http.StatusBadRequest) + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #2, empty type, record not found",
			body: map[string]interface{}{
				"id": "testCounter",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusNotFound,
				response:    http.StatusText(http.StatusNotFound) + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #3, incorrect type, record not found",
			body: map[string]interface{}{
				"id":   "testCounter",
				"type": "none",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusNotFound,
				response:    http.StatusText(http.StatusNotFound) + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #4, empty id, record not found",
			body: map[string]interface{}{
				"type": "counter",
			},
			method:  http.MethodPost,
			wantErr: true,
			want: want{
				code:        http.StatusNotFound,
				response:    http.StatusText(http.StatusNotFound) + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		}, {
			name: "positive test #5, counter type return json",
			body: map[string]interface{}{
				"id":   "PollCount",
				"type": "counter",
			},
			method:  http.MethodPost,
			wantErr: false,
			want: want{
				code: http.StatusOK,
				response: `{
					"id": "PollCount",
					"type": "counter",
					"delta": 14123413542
				}`,
				contentType: "application/json",
			},
		},
	}
	r := chi.NewRouter()

	st, err := mock.NewStorageMock()
	if err != nil {
		t.Error(err)
	}
	sr := NewService(st)
	h := NewHandler(nil, sr)

	r.Post("/value/", h.GetMetric)

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, tt.method, "/value/", tt.body)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			if tt.wantErr {
				assert.Equal(t, tt.want.response, get)
			} else {
				assert.JSONEq(t, tt.want.response, get)
			}
		})
	}
}

func Test_GetAllMetrics(t *testing.T) {
	r := chi.NewRouter()

	st, err := mock.NewStorageMock()
	if err != nil {
		t.Error(err)
	}

	sr := NewService(st)
	h := NewHandler(nil, sr)

	r.Get("/", h.GetAllMetrics)

	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("positive test #1, return html", func(t *testing.T) {

		resp, _ := testRequest(t, ts, http.MethodGet, "/", nil)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
	})
}

func ExampleHandler_CreateOrUpdateByParams() {
	store, _ := storage.NewMemoryStorage("", 0, false) // Instantiate your storage.
	service := NewService(store)                       // Instantiate your service.
	handler := NewHandler(nil, service)                // Instantiate your handler.

	// Create a new HTTP request with query params.
	req := httptest.NewRequest("GET", "/metrics?type=cpu&name=usage&value=50", nil)
	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()

	// Call the handler function directly, passing in the ResponseRecorder and Request.
	handler.CreateOrUpdateByParams(rr, req)

	// Check the status code and response body to verify the result.
	if status := rr.Code; status != http.StatusOK {
		fmt.Printf("Example 1: Unexpected status code: %d\n", status)
	}
}

func ExampleHandler_CreateOrUpdateByJSON() {
	store, _ := storage.NewMemoryStorage("", 0, false) // Instantiate your storage.
	service := NewService(store)                       // Instantiate your service.
	handler := NewHandler(nil, service)                // Instantiate your handler.

	// Example 1: Mocking a successful JSON request
	// Define a sample JSON request body.
	requestBody := `{"ID": "123", "MType": "cpu", "Delta": 10}`
	// Create a new HTTP request with the defined JSON body.
	req := httptest.NewRequest("POST", "/create", strings.NewReader(requestBody))
	// Set the request Content-Type header to indicate JSON format.
	req.Header.Set("Content-Type", "application/json")
	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	// Call the handler function directly, passing in the ResponseRecorder and Request.
	handler.CreateOrUpdateByJSON(rr, req)
	// Check the status code and response body to verify the result.
	if status := rr.Code; status != http.StatusOK {
		fmt.Printf("Example 1: Unexpected status code: %d\n", status)
	}
}

func ExampleHandler_CreateOrUpdateMany() {
	store, _ := storage.NewMemoryStorage("", 0, false) // Instantiate your storage.
	service := NewService(store)                       // Instantiate your service.
	handler := NewHandler(nil, service)                // Instantiate your handler.

	// Example 1: Mocking a successful JSON request with multiple metrics
	// Define a sample JSON request body with multiple metrics.
	requestBody := `[{"ID": "123", "MType": "cpu", "Delta": 10}, {"ID": "456", "MType": "memory", "Value": 80}]`
	// Create a new HTTP request with the defined JSON body.
	req := httptest.NewRequest("POST", "/create-many", strings.NewReader(requestBody))
	// Set the request Content-Type header to indicate JSON format.
	req.Header.Set("Content-Type", "application/json")
	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	// Call the handler function directly, passing in the ResponseRecorder and Request.
	handler.CreateOrUpdateMany(rr, req)
	// Check the status code and response body to verify the result.
	if status := rr.Code; status != http.StatusOK {
		fmt.Printf("Example 1: Unexpected status code: %d\n", status)
	}
}

func ExampleHandler_GetMetric() {
	store, _ := storage.NewMemoryStorage("", 0, false) // Instantiate your storage.
	service := NewService(store)                       // Instantiate your service.
	handler := NewHandler(nil, service)                // Instantiate your handler.

	// Example 1: Mocking a successful request to get a metric
	// Define a sample JSON request body.
	requestBody := `{"MType": "cpu", "ID": "123"}`
	// Create a new HTTP request with the defined JSON body.
	req := httptest.NewRequest("GET", "/metric", strings.NewReader(requestBody))
	// Set the request Content-Type header to indicate JSON format.
	req.Header.Set("Content-Type", "application/json")
	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	// Call the handler function directly, passing in the ResponseRecorder and Request.
	handler.GetMetric(rr, req)
	// Check the status code and response body to verify the result.
	if status := rr.Code; status != http.StatusOK {
		fmt.Printf("Example 1: Unexpected status code: %d\n", status)
	}
}

func ExampleHandler_GetMetricValue() {
	store, _ := storage.NewMemoryStorage("", 0, false) // Instantiate your storage.
	service := NewService(store)                       // Instantiate your service.
	handler := NewHandler(nil, service)                // Instantiate your handler.

	// Example 1: Mocking a successful request to get a metric value
	// Create a new HTTP request with desired URL parameters.
	req := httptest.NewRequest("GET", "/metric?type=cpu&name=usage", nil)
	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	// Call the handler function directly, passing in the ResponseRecorder and Request.
	handler.GetMetricValue(rr, req)
	// Check the status code and response body to verify the result.
	if status := rr.Code; status != http.StatusOK {
		fmt.Printf("Example 1: Unexpected status code: %d\n", status)
	}
}

func ExampleHandler_GetAllMetrics() {
	store, _ := storage.NewMemoryStorage("", 0, false) // Instantiate your storage.
	service := NewService(store)                       // Instantiate your service.
	handler := NewHandler(nil, service)                // Instantiate your handler.

	// Example 1: Mocking a successful request to get all metrics in HTML format
	// Create a new HTTP request.
	req := httptest.NewRequest("GET", "/metrics", nil)
	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	// Call the handler function directly, passing in the ResponseRecorder and Request.
	handler.GetAllMetrics(rr, req)
	// Check the status code and response body to verify the result.
	if status := rr.Code; status != http.StatusOK {
		fmt.Printf("Example 1: Unexpected status code: %d\n", status)
	}
}
