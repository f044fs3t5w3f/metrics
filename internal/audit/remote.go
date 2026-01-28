package audit

import "github.com/f044fs3t5w3f/metrics/pkg/net"

type RemoteAudit struct {
	url string
}

// TODO: use client

func (r *RemoteAudit) Notify(ev *Event) {
	go func() {
		net.SendZippedJSON(r.url, ev)
	}()
}

func NewRemoteAudit(url string) *RemoteAudit {
	return &RemoteAudit{
		url: url,
	}
}

var _ Audit = (*RemoteAudit)(nil)
