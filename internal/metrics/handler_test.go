package metrics

import (
	"fmt"
	"github.com/dkmelnik/metrics/internal/metrics/mock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestHandler_Create(t *testing.T) {
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
	ts := httptest.NewServer(ConfigureRouter())
	defer ts.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, tt.method, tt.request)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, get)
		})
	}
}

func TestHandler_Get(t *testing.T) {
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
				response:    "3833856",
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
				response:    "7.708",
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
				response:    "3485734.1",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "positive test #7, valid gauge NextGC metric",
			metricsType: "gauge",
			metricsName: "NextGC",
			want: want{
				code:        http.StatusOK,
				response:    "-3358720",
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

	st := mock.NewStorageMock()
	sr := NewService(st)
	h := NewHandler(sr)

	r.Get("/value/{type}/{name}", h.Get)

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, get := testRequest(t, ts, http.MethodGet, "/value"+"/"+tt.metricsType+"/"+tt.metricsName)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, get)
		})
	}
}

func TestHandler_GetAll(t *testing.T) {
	r := chi.NewRouter()

	st := mock.NewStorageMock()
	sr := NewService(st)
	h := NewHandler(sr)

	r.Get("/", h.GetAll)

	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("positive test #1, return html", func(t *testing.T) {

		resp, _ := testRequest(t, ts, http.MethodGet, "/")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
	})
}
