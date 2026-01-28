package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type FileAudit struct {
	file *os.File
	lock sync.Mutex
}

func (f *FileAudit) Notify(ev *Event) {
	f.lock.Lock()
	go func() {
		defer f.lock.Unlock()
		encoder := json.NewEncoder(f.file)
		encoder.Encode(ev)
	}()
}

func NewFileAudit(filename string) (*FileAudit, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return nil, fmt.Errorf("NewFileAudit: %w", err)
	}
	return &FileAudit{
		file: file,
	}, nil
}

var _ Audit = (*FileAudit)(nil)
