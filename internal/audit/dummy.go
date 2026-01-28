package audit

type Dummy struct{}

func (d Dummy) Notify(*Event) {
}
