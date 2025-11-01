package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/stretchr/testify/assert"
)

type getTestMockStorage struct {
}

func (g *getTestMockStorage) AddCounter(metricName string, value int64) {
	panic("unimplemented")
}

func (g *getTestMockStorage) GetCounter(metricName string) (int64, error) {
	if metricName == "exists" {
		return 5, nil
	} else {
		return 0, errors.New("")
	}
}

func (g *getTestMockStorage) GetGauge(metricName string) (float64, error) {
	if metricName == "exists" {
		return 5, nil
	} else {
		return 0, errors.New("")
	}
}

func (g *getTestMockStorage) SetGauge(metricName string, value float64) {
	panic("unimplemented")
}

var _ repository.Storage = &getTestMockStorage{}

func TestGetJson(t *testing.T) {
	type testCase struct {
		name      string
		request   string
		wantError bool
		response  string
	}
	testCases := []testCase{
		{
			"Correct counter",
			`{"id":"exists","type":"counter"}`,
			false,
			`{"id":"exists","type":"counter","delta":5}`,
		},
		{
			"Correct gauge",
			`{"id":"exists","type":"gauge"}`,
			false,
			`{"id":"exists","type":"gauge","value":5}`,
		},
		{
			"Gauge not exists",
			`{"id":"n","type":"gauge"}`,
			true,
			"",
		},
		{
			"Counter not exists",
			`{"id":"n","type":"counter"}`,
			true,
			"",
		},
		{
			"Malformed json",
			`:(`,
			true,
			"",
		},
		{
			"Incorrect type",
			`{"id":"n","type":"type"}`,
			true,
			"",
		},
	}
	storage := &getTestMockStorage{}
	handler := GetJson(storage)
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			body := strings.NewReader(tCase.request)
			req := httptest.NewRequest(http.MethodPost, "/", body)
			recorder := httptest.NewRecorder()
			handler(recorder, req)
			res := recorder.Result()
			defer res.Body.Close()

			assert.Equal(t, tCase.wantError, res.StatusCode != http.StatusOK)
			if tCase.wantError {
				return
			}
			responseBytes, _ := io.ReadAll(res.Body)
			got := strings.TrimSuffix(string(responseBytes), "\n")
			assert.Equal(t, tCase.response, got)
		})
	}
}
