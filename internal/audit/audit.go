// Audit package provides an interface for auditing metrics updates
// and it's implementations
package audit

// Event is type representing an metric update
type Event struct {
	Timestamp int64
	Metrics   []string
	IP        string
}

// Audit is an interface for auditing metrics updates
type Audit interface {
	Notify(*Event)
}
