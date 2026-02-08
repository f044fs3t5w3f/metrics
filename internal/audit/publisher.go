package audit

// Publisher is struct using as linker of different audit implementations
type Publisher struct {
	subscribers []Audit
}

func (p *Publisher) Notify(event *Event) {
	for _, subscriber := range p.subscribers {
		subscriber.Notify(event)
	}
}

func NewAuditPublisher(subscribers []Audit) *Publisher {
	return &Publisher{
		subscribers: subscribers,
	}
}

func (p *Publisher) AddSubscriber(subscriber Audit) {
	for _, s := range p.subscribers {
		if s == subscriber {
			return
		}
	}
	p.subscribers = append(p.subscribers, subscriber)
}

func (p *Publisher) RemoveSubscriber(subscriber Audit) {
	for i, s := range p.subscribers {
		if s == subscriber {
			p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
			return
		}
	}
}

var _ Audit = (*Publisher)(nil)
