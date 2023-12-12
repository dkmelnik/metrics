package metrics

import (
	"fmt"
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(ConfigureRouter())
			defer ts.Close()

			resp, get := testRequest(t, ts, tt.method, tt.request)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, get)
		})
	}
}
