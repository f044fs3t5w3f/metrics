package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			request: request{method: http.MethodPost, url: "/update/counter/lol"},
			want:    want{code: http.StatusNotFound},
		},
		{
			name:    "bad value",
			request: request{method: http.MethodPost, url: "/update/counter/lol/100f"},
			want:    want{code: http.StatusBadRequest},
		},
		{
			name:    "Correct counter",
			request: request{method: http.MethodPost, url: "/update/counter/lol/100"},
			want:    want{code: http.StatusOK},
		},
		{
			name:    "Float counter",
			request: request{method: http.MethodPost, url: "/update/counter/lol/100.1"},
			want:    want{code: http.StatusBadRequest},
		},
		{
			name:    "Correct gauge",
			request: request{method: http.MethodPost, url: "/update/gauge/lol/100"},
			want:    want{code: http.StatusOK},
		},
		{
			name:    "Float gauge",
			request: request{method: http.MethodPost, url: "/update/gauge/lol/100.2"},
			want:    want{code: http.StatusOK},
		},
	}

	storage := repository.NewMemStorage()
	router := GetRouter(storage)
	ts := httptest.NewServer(router)
	defer ts.Close()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.request.method, test.request.url)
			assert.Equal(t, test.want.code, resp.StatusCode)
		})
	}
}
