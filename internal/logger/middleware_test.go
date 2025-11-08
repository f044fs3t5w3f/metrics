package logger

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const sleepDuration = 1.0

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	time.Sleep(sleepDuration * time.Second)
	w.Write([]byte("Hello World!"))
}

func TestRequestLogger(t *testing.T) {
	const path = "/path/"
	request := httptest.NewRequest(http.MethodGet, path, nil)
	recorder := httptest.NewRecorder()
	mycore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stderr,
		zapcore.InfoLevel,
	)
	observed, logs := observer.New(zapcore.InfoLevel)
	Log = zap.New(zapcore.NewTee(mycore, observed))
	var h http.HandlerFunc = helloWorldHandler
	RequestLogger(h).ServeHTTP(recorder, request)
	entries := logs.All()
	if len(entries) != 1 {
		t.Error("Incorrect number of logs")
	}
	entry := entries[0]
	records := entry.ContextMap()
	if methodGot := records["method"]; methodGot != http.MethodGet {
		t.Errorf("Incorrect logged method, want: %s, got: %s", http.MethodGet, methodGot)
	}
	if pathGot := records["path"]; pathGot != path {
		t.Errorf("Incorrect logged method, want: %s, got: %s", path, pathGot)
	}
	duration := records["duration"]
	durationSeconds, ok := duration.(float64)
	if !ok {
		t.Errorf("Failed to convert duration")
	}
	if durationSeconds < sleepDuration {
		t.Errorf("Incorrect logged method, want at least %f, got: %f", sleepDuration, durationSeconds)
	}
}
