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

// GetSignMiddleware returns an HTTP middleware that verifies the request's signature (if present) and signs the response.
//
// The middleware:
//  1. Checks for the "HashSHA256" header in the incoming request.
//  2. If the header is present, it computes request body hash using the provided signFunc,
//     and compares it with the value in the header. If they don't match, it returns a 400 error.
//  3. If the signature is valid or not required, it proceeds to serve the request.
//  4. In case if request HashSHA256 header was provided it signs the response body using the same function.
//
// Parameters:
//   - signFunc: A function that computes the hash/signature of a given byte slice.
//
// Returns:
//   - A middleware that applies request signature verification and response signing.
func GetSignMiddleware(signFunc SignFunc) middleware {
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

			resp := buf.Bytes()
			respHash := signFunc(resp)
			respHashB64 := base64.StdEncoding.EncodeToString(respHash)

			w.Header().Set("HashSHA256", respHashB64)
			w.Write(resp)
		})
	}
}
