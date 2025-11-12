package file

import (
	"bytes"
	"testing"

	"github.com/f044fs3t5w3f/metrics/internal/repository/memory"
	"github.com/stretchr/testify/assert"
)

func TestRestoreFromFile(t *testing.T) {
	storage := &fileStorage{
		Storage: memory.NewMemStorage(),
	}

	data := `[
		{"id":"gauge1","type":"gauge","value":100},
		{"id":"gauge2","type":"gauge","value":1000},
		{"id":"counter1","type":"counter","delta":10000},
		{"id":"counter2","type":"counter","delta":100000}
]`
	reader := bytes.NewReader([]byte(data))
	err := storage.restoreFromFile(reader)
	assert.Nil(t, err, "error while reading correct json")
	// assert.Equal(t, 2, len(storage.gauge), "incorrect number of gauges")
	value1, _ := storage.GetGauge("gauge1")
	assert.Equal(t, value1, float64(100))
	value2, _ := storage.GetGauge("gauge2")
	assert.Equal(t, value2, float64(1000))

	value3, _ := storage.GetCounter("counter1")
	assert.Equal(t, value3, int64(10000))

	value4, _ := storage.GetCounter("counter2")
	assert.Equal(t, value4, int64(100000))

}
func TestRestoreFromFileMalformedJson(t *testing.T) {
	storage := &fileStorage{
		Storage: memory.NewMemStorage(),
	}

	data := `:(`
	reader := bytes.NewReader([]byte(data))
	err := storage.restoreFromFile(reader)
	assert.NotNil(t, err, "malformed json: error is nil")
}
