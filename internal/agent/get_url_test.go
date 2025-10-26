package agent

import (
	"fmt"
	"testing"

	"github.com/f044fs3t5w3f/metrics/internal/models"
)

func int64Ptr(val int64) *int64 {
	return &val
}
func float64Ptr(val float64) *float64 {
	return &val
}

func TestGetURL(t *testing.T) {
	hostName := "example.com"
	type testCase struct {
		metric *models.Metrics
		want   string
	}

	cases := []testCase{
		{
			metric: &models.Metrics{
				ID:    "M1",
				MType: models.Counter,
				Delta: int64Ptr(42),
			},
			want: "http://example.com/update/counter/M1/42",
		},
		{
			metric: &models.Metrics{
				ID:    "M2",
				MType: models.Gauge,
				Value: float64Ptr(42.1),
			},
			want: "http://example.com/update/gauge/M2/42.1",
		},
	}
	for _, testCase := range cases {
		url, err := getURL(hostName, testCase.metric)
		if err != nil {
			t.Error("GetURL: Expected no error but got one")
		}
		if url != testCase.want {
			t.Errorf("GetURL: Expected \"%s\" but got \"%s\"", testCase.want, url)
		}
	}
}

func TestGetURLError(t *testing.T) {
	hostName := "example.com"
	cases := []*models.Metrics{
		{
			ID:    "M1",
			MType: models.Counter,
		},
		{
			ID:    "M2",
			MType: models.Gauge,
		},
	}
	for _, metric := range cases {
		url, err := getURL(hostName, metric)
		if err == nil {
			fmt.Println(url)
			t.Error("Expected an error, but got nil")
		}
	}
}
