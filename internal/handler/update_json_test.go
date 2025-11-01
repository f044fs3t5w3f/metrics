package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/stretchr/testify/assert"
)

type updateTestMockStorage struct {
	calls []string
}

var _ repository.Storage = &updateTestMockStorage{}

func (m *updateTestMockStorage) AddCounter(metricName string, value int64) {
	m.calls = append(m.calls, fmt.Sprintf("counter;%s;%d", metricName, value))
}

func (m *updateTestMockStorage) GetCounter(metricName string) (int64, error) {
	panic("unimplemented")
}

func (m *updateTestMockStorage) GetGauge(metricName string) (float64, error) {
	panic("unimplemented")
}

func (m *updateTestMockStorage) SetGauge(metricName string, value float64) {
	m.calls = append(m.calls, fmt.Sprintf("gauge;%s;%f", metricName, value))
}

func TestUpdateJson(t *testing.T) {
	type testCase struct {
		name   string
		body   string
		status int
		calls  []string
	}
	testCases := []testCase{
		{"Good counter", `{"id": "n","type": "counter","delta": 10}`, http.StatusOK, []string{"counter;n;10"}},
		{"Good gauge", `{"id": "n","type": "gauge","value": 10}`, http.StatusOK, []string{"gauge;n;10.000000"}},
		{"Nil counter", `{"id": "n","type": "counter"}`, http.StatusBadRequest, []string{}},
		{"Nil gauge", `{"id": "n","type": "gauge"}`, http.StatusBadRequest, []string{}},
		{"Incorrect type", `{"id": "n","type": "lol"}`, http.StatusBadRequest, []string{}},
		{"Malformed json", `:(`, http.StatusBadRequest, []string{}},
	}

	for _, tCase := range testCases {
		storage := &updateTestMockStorage{calls: []string{}}
		handler := UpdateJSON(storage)
		t.Run(tCase.name, func(t *testing.T) {
			body := strings.NewReader(tCase.body)
			req := httptest.NewRequest(http.MethodPost, "/", body)
			recorder := httptest.NewRecorder()
			handler(recorder, req)
			res := recorder.Result()
			defer res.Body.Close()

			assert.Equal(t, tCase.status, res.StatusCode)
			assert.Equal(t, tCase.calls, storage.calls)
		})
	}

}
