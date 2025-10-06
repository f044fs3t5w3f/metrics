package main

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

var memStorage MemStorage

func init() {
	memStorage = MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		w.Write([]byte(":("))
		return
	}
	urlParams := strings.Split(r.URL.Path[1:], "/")

	if len(urlParams) < 4 || urlParams[0] != "update" {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	type_, metricName, merticValueStr := urlParams[1], urlParams[2], urlParams[3]
	if type_ != Gauge && type_ != Counter {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	switch type_ {
	case Gauge:
		merticParsed, err := strconv.ParseFloat(merticValueStr, 64)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		memStorage.gauge[metricName] = merticParsed
	case Counter:
		merticParsed, err := strconv.ParseInt(merticValueStr, 0, 64)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if _, ok := memStorage.counter[metricName]; !ok {
			memStorage.counter[metricName] = 0
		}
		memStorage.counter[metricName] += merticParsed
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleMain)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(":(")
	}

}
