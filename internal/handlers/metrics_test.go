package handlers

import (
	"github.com/dkmelnik/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
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
				response:    http.StatusText(http.StatusMethodNotAllowed),
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "negative test #2, method does not match pattern update/type/name/value",
			request: "/update/gauge/Alloc",
			method:  http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				response:    http.StatusText(http.StatusNotFound),
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "negative test #3, missing value",
			request: "/update/gauge/Alloc/",
			method:  http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				response:    http.StatusText(http.StatusBadRequest),
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "negative test #4, value not int64 or float64",
			request: "/update/gauge/Alloc/test",
			method:  http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				response:    http.StatusText(http.StatusBadRequest),
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "negative test #5, the specified type was not found",
			request: "/update/test/Alloc/1",
			method:  http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				response:    http.StatusText(http.StatusBadRequest),
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "positive test #6, must return OK",
			request: "/update/gauge/Alloc/1",
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
			request := httptest.NewRequest(tt.method, tt.request, nil)

			w := httptest.NewRecorder()
			s := storage.NewCollection()
			h := http.HandlerFunc(NewHandler(s).Create)
			h(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			//assert.Equal(t, tt.want.response, result.Body)

			res, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.response, strings.ReplaceAll(string(res), "\n", ""))

		})
	}
}
