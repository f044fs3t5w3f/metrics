package repository

type Storage interface {
	GetCounter(metricName string) (int64, error)
	GetGauge(metricName string) (float64, error)
	AddCounter(metricName string, value int64)
	SetGauge(metricName string, value float64)
}
