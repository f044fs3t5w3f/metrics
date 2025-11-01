package compress

import (
	"net/http"
	"strings"
)

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			defer cw.Close()
			w = cw
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
