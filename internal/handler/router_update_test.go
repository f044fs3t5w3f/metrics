package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/f044fs3t5w3f/metrics/internal/audit"
	"github.com/f044fs3t5w3f/metrics/internal/repository/memory"
	"github.com/f044fs3t5w3f/metrics/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type request struct {
	method string
	url    string
}

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

func TestRouter(t *testing.T) {

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
			want:    want{code: http.StatusMethodNotAllowed},
		},
		{
			name:    "Empty url",
			request: request{method: http.MethodGet, url: "/"},
			want:    want{code: http.StatusOK},
		},
		{
			name:    "Empty type",
			request: request{method: http.MethodPost, url: "/update/haha/"},
			want:    want{code: http.StatusNotFound},
		},
		{
			name:    "bad type",
			request: request{method: http.MethodPost, url: "/update/cntr/lol/100"},
			want:    want{code: http.StatusBadRequest},
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
			name:    "Bad gauge",
			request: request{method: http.MethodPost, url: "/update/gauge/lol/100j"},
			want:    want{code: http.StatusBadRequest},
		},
		{
			name:    "Float gauge",
			request: request{method: http.MethodPost, url: "/update/gauge/lol/100.2"},
			want:    want{code: http.StatusOK},
		},
	}

	storage := memory.NewMemStorage()
	service := service.NewService(storage, audit.Dummy{})
	router := GetRouter(storage, service, "")
	ts := httptest.NewServer(router)
	defer ts.Close()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.request.method, test.request.url)
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
		})
	}
}

func TestSequense(t *testing.T) {
	type response struct {
		statusCode int
		response   string
	}
	type testCase struct {
		name     string
		requests []request
		want     response
	}
	testTable := []testCase{
		{
			"Empty gauge",
			[]request{{http.MethodGet, "/value/gauge/lol"}},
			response{http.StatusNotFound, ""},
		},
		{
			"Set gauge",
			[]request{
				{http.MethodPost, "/update/gauge/lol/4.5"},
				{http.MethodGet, "/value/gauge/lol"},
			},
			response{http.StatusOK, "4.5"},
		},
		{
			"Empty counter",
			[]request{{http.MethodGet, "/value/counter/lol"}},
			response{http.StatusNotFound, ""},
		},
		{
			"Set counter counter",
			[]request{
				{http.MethodPost, "/update/counter/lol/4"},
				{http.MethodGet, "/value/counter/lol"},
			},
			response{http.StatusOK, "4"},
		},
	}

	storage := memory.NewMemStorage()
	service := service.NewService(storage, audit.Dummy{})
	router := GetRouter(storage, service, "")
	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			var response *http.Response
			var content string
			for _, requst := range testCase.requests {
				response, content = testRequest(t, ts, requst.method, requst.url)
				response.Body.Close()
			}
			assert.Equal(t, response.StatusCode, testCase.want.statusCode)
			if testCase.want.response != "" {
				assert.Equal(t, testCase.want.response, content)
			}
		})
	}

}
