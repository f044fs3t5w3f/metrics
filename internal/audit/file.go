package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Audit interface implementations using file as a storage
type FileAudit struct {
	file *os.File
	lock sync.Mutex
}

// NewFileAudit creates new FileAudit
// filename is a path to file
// returns *FileAudit  and error if the one has occurred
func NewFileAudit(filename string) (*FileAudit, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return nil, fmt.Errorf("NewFileAudit: %w", err)
	}
	return &FileAudit{
		file: file,
	}, nil
}

func (f *FileAudit) Notify(ev *Event) {
	f.lock.Lock()
	go func() {
		defer f.lock.Unlock()
		encoder := json.NewEncoder(f.file)
		encoder.Encode(ev)
	}()
}

var _ Audit = (*FileAudit)(nil)
