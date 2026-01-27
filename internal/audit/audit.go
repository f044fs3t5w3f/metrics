package audit

type Event struct {
	Timestamp int64
	Metrics   []string
	IP        string
}

type Audit interface {
	Notify(*Event)
}
