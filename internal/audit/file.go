package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/f044fs3t5w3f/metrics/internal/logger"
)

// Audit interface implementations using file as a storage
type FileAudit struct {
	ch         chan *Event
	closedChan chan struct{}
}

const bufsize = 10

// NewFileAudit creates new FileAudit
// filename is a path to file
// returns *FileAudit  and error if the one has occurred
func NewFileAudit(ctx context.Context, filename string) (*FileAudit, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return nil, fmt.Errorf("NewFileAudit: %w", err)
	}
	ch := make(chan *Event, bufsize)
	closedChan := make(chan struct{})
	go func() {
		encoder := json.NewEncoder(file)
		for ev := range ch {
			encoder.Encode(ev)
		}
		logger.Log.Info("closing audit file")
		file.Close()
		closedChan <- struct{}{}
	}()
	fa := &FileAudit{
		ch:         ch,
		closedChan: closedChan,
	}

	return fa, nil
}

func (f *FileAudit) Notify(ev *Event) {
	f.ch <- ev
}

func (f *FileAudit) Close() {
	close(f.ch)
	<-f.closedChan
}

var _ Audit = (*FileAudit)(nil)
