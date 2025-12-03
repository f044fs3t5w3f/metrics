package sign

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
)

type middleware func(next http.Handler) http.Handler

type signWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (sw *signWriter) Write(b []byte) (int, error) {
	return sw.buf.Write(b)
}

func GetSignMiddleware(signFunc signFunc) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hashHeader := r.Header.Get("HashSHA256")
			toCheck := hashHeader != ""
			if !toCheck {
				next.ServeHTTP(w, r)
				return
			}

			reqBody, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "can't read body", http.StatusBadRequest)
				return
			}
			r.Body.Close()

			r.Body = io.NopCloser(bytes.NewReader(reqBody))

			hash := signFunc(reqBody)
			expectedB64 := base64.StdEncoding.EncodeToString(hash)

			if hashHeader != expectedB64 {
				http.Error(w, "Invalid signature", http.StatusBadRequest)
				return
			}

			buf := &bytes.Buffer{}
			sw := &signWriter{ResponseWriter: w, buf: buf}

			next.ServeHTTP(sw, r)

			// считаем подпись ответа
			resp := buf.Bytes()
			respHash := signFunc(resp)
			respHashB64 := base64.StdEncoding.EncodeToString(respHash)

			w.Header().Set("HashSHA256", respHashB64)
			w.Write(resp)
		})
	}
}
