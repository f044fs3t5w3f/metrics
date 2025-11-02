package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  io.Writer
	ow http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  nil,
		ow: w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.ow.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	if c.w == nil {
		if c.isJSONContentType() {
			c.ow.Header().Set("Content-Encoding", "gzip")
			c.w = c.zw
		} else {
			c.w = c.ow
		}
	}
	return c.w.Write(p)
}

func (c *compressWriter) isJSONContentType() bool {
	contentType := c.Header().Get("Content-Type")
	return strings.Contains(contentType, "application/json")

}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode >= 300 || !c.isJSONContentType() {
		c.w = c.ow
	} else {
		c.ow.Header().Set("Content-Encoding", "gzip")
		c.w = c.zw
	}
	c.ow.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	if c.w == c.zw {
		return c.zw.Close()
	}
	return nil
}
