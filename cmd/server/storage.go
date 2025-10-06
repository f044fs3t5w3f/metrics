package main

type memStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}
