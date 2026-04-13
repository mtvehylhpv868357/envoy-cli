package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// EventType represents the kind of audit event.
type EventType string

const (
	EventProfileSwitch EventType = "profile.switch"
	EventProfileCreate EventType = "profile.create"
	EventProfileDelete EventType = "profile.delete"
	EventVaultEncrypt  EventType = "vault.encrypt"
	EventVaultDecrypt  EventType = "vault.decrypt"
	EventSnapshotSave  EventType = "snapshot.save"
	EventSnapshotLoad  EventType = "snapshot.load"
)

// Entry is a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	Project   string    `json:"project,omitempty"`
	Detail    string    `json:"detail,omitempty"`
}

// Logger writes audit entries to a log file.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the file at path.
func NewLogger(dir string) (*Logger, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("audit: create dir: %w", err)
	}
	return &Logger{path: filepath.Join(dir, "audit.log")}, nil
}

// Log appends an entry to the audit log.
func (l *Logger) Log(event EventType, project, detail string) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Project:   project,
		Detail:    detail,
	}
	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}

// ReadAll returns all entries from the audit log.
func (l *Logger) ReadAll() ([]Entry, error) {
	data, err := os.ReadFile(l.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read log: %w", err)
	}

	var entries []Entry
	for _, raw := range splitLines(data) {
		if len(raw) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(raw, &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
