// Compress package provides middleware for reading compressed requests and compressing responses.
package compress

import (
	"net/http"
	"strings"
)

// Middleware is a function that enhances an existing HTTP handler by adding gzip compression.
// It inspects incoming requests and outgoing responses to determine whether to apply gzip
// compression or decompression according to the following rules:
//
// 1. For incoming requests:
//   - If the "Content-Encoding" header specifies "gzip",
//     the middleware will automatically decompress the request body.
//
// 2. For outgoing responses:
//   - If the "Accept-Encoding" header includes "gzip",
//     the middleware will compress the response body.
//
// Errors occurring during compression or decompression are handled by returning
// an Internal Server Error (500).
//
// Arguments:
//   - next: The next handler in the chain that processes the request.
//
// Returns:
//   - A new Handler instance wrapped with gzip compression functionality.
func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			defer cw.Close()
			w = cw
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
