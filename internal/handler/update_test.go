package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestHandleUpdate(t *testing.T) {
	type request struct {
		method string
		url    string
	}

	type want struct {
		code int
	}

	type test struct {
		name    string
		request request
		want    want
	}
	tests := []test{
		{
			name:    "Wrong method",
			request: request{method: http.MethodGet, url: "/update/counter/someMetric/527"},
			want:    want{code: http.StatusMethodNotAllowed},
		},
		{
			name:    "Wrong method",
			request: request{method: http.MethodPut, url: "/update/counter/someMetric/527"},
			want:    want{code: http.StatusMethodNotAllowed},
		},
		{
			name:    "Empty url",
			request: request{method: http.MethodPost, url: "/"},
			want:    want{code: http.StatusNotFound},
		},
		{
			name:    "Empty counter",
			request: request{method: http.MethodPost, url: "/update/"},
			want:    want{code: http.StatusNotFound},
		},
		{
			name:    "Bad type",
			request: request{method: http.MethodPost, url: "/update/haha/"},
			want:    want{code: http.StatusNotFound},
		},
		{
			name:    "without metric",
			request: request{method: http.MethodPost, url: "/update/counter/"},
			want:    want{code: http.StatusNotFound},
		},
		{
			name:    "without value",
			request: request{method: http.MethodPost, url: "/update/counter/lol/"},
			want:    want{code: http.StatusBadRequest},
		},
		{
			name:    "bad value",
			request: request{method: http.MethodPost, url: "/update/counter/lol/100f/"},
			want:    want{code: http.StatusBadRequest},
		},
		{
			name:    "Correct counter",
			request: request{method: http.MethodPost, url: "/update/counter/lol/100/"},
			want:    want{code: http.StatusOK},
		},
		{
			name:    "Float counter",
			request: request{method: http.MethodPost, url: "/update/counter/lol/100.1/"},
			want:    want{code: http.StatusBadRequest},
		},
		{
			name:    "Correct gauge",
			request: request{method: http.MethodPost, url: "/update/gauge/lol/100/"},
			want:    want{code: http.StatusOK},
		},
		{
			name:    "Float gauge",
			request: request{method: http.MethodPost, url: "/update/gauge/lol/100.2/"},
			want:    want{code: http.StatusOK},
		},
	}

	storage := models.MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.method, test.request.url, nil)
			w := httptest.NewRecorder()
			Update(storage)(w, request)
			res := w.Result()
			res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
