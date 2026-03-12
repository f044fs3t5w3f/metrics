package net

import (
	"net"
	"net/http"
)

func GetCheckSubnetMiddleware(network *net.IPNet) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if network == nil {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			realIP := r.Header.Get("X-Real-IP")
			if realIP == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			ip := net.ParseIP(realIP)
			if !network.Contains(ip) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
