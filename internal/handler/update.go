package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/f044fs3t5w3f/metrics/internal/models"
)

func Update(storage models.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			w.Write([]byte(":("))
			return
		}
		urlParams := strings.Split(r.URL.Path[1:], "/")

		if len(urlParams) < 4 {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		type_, metricName, merticValueStr := urlParams[1], urlParams[2], urlParams[3]

		switch type_ {
		case models.Gauge:
			merticParsed, err := strconv.ParseFloat(merticValueStr, 64)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			storage.Gauge[metricName] = merticParsed
		case models.Counter:
			merticParsed, err := strconv.ParseInt(merticValueStr, 0, 64)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			if _, ok := storage.Counter[metricName]; !ok {
				storage.Counter[metricName] = 0
			}
			storage.Counter[metricName] += merticParsed
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}
}
