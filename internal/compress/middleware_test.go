package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	var _1000a = make([]byte, 1000)
	for i := range 1000 {
		_1000a[i] = 'a'
	}
	statusOKHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(_1000a)
		w.WriteHeader(http.StatusOK)
	})
	emptyContentTypeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("!"))
		w.WriteHeader(http.StatusOK)
	})
	internalServerErrorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})

	type testCase struct {
		name                 string
		handler              http.HandlerFunc
		acceptGzip           bool
		expectedResponseBody []byte
		expectCompress       bool
	}
	testCases := []testCase{
		{"Supports gzip", statusOKHandler, true, _1000a, true},
		{"Empty content type", emptyContentTypeHandler, true, []byte("!"), false},
		{"Doesn't support gzip", statusOKHandler, false, _1000a, false},
		{"Error handler", internalServerErrorHandler, true, []byte("Internal Server Error\n"), false},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", nil)
			if tCase.acceptGzip {
				request.Header.Add("Accept-Encoding", "gzip")
			}
			recorder := httptest.NewRecorder()
			Middleware(tCase.handler).ServeHTTP(recorder, request)
			encoding := recorder.Header().Get("Content-Encoding")
			gziped := encoding == "gzip"
			assert.Equal(t, tCase.expectCompress, gziped)
			body, _ := io.ReadAll(recorder.Body)
			if tCase.expectCompress {
				gzipReader, err := gzip.NewReader(bytes.NewReader(body))
				assert.NoError(t, err, "Error while creating gzip decoder")
				decoded, err := io.ReadAll(gzipReader)
				assert.NoError(t, err, "Error while decoding")
				assert.Equal(t, tCase.expectedResponseBody, decoded)
			} else {
				assert.Equal(t, tCase.expectedResponseBody, body)

			}
		})
	}
}
