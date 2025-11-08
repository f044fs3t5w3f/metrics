package repository

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestoreFromFile(t *testing.T) {
	storage := &memStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
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
	assert.Equal(t, 2, len(storage.gauge), "incorrect number of gauges")
	assert.Equal(t, storage.gauge["gauge1"], float64(100))
	assert.Equal(t, storage.gauge["gauge2"], float64(1000))
	assert.Equal(t, 2, len(storage.counter), "incorrect number of counters")
	assert.Equal(t, storage.counter["counter1"], int64(10000))
	assert.Equal(t, storage.counter["counter2"], int64(100000))
}
func TestRestoreFromFileMalformedJson(t *testing.T) {
	storage := &memStorage{}

	data := `:(`
	reader := bytes.NewReader([]byte(data))
	err := storage.restoreFromFile(reader)
	assert.NotNil(t, err, "malformed json: error is nil")
}
