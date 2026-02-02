package audit

// Dummy implementation using for testing
type Dummy struct{}

func (d Dummy) Notify(*Event) {
}
