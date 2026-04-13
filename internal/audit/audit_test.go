package audit

import (
	"os"
	"testing"
	"time"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "audit-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNewLogger(t *testing.T) {
	dir := tempDir(t)
	l, err := NewLogger(dir)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLog_CreatesEntry(t *testing.T) {
	dir := tempDir(t)
	l, _ := NewLogger(dir)

	if err := l.Log(EventProfileSwitch, "myproject", "dev"); err != nil {
		t.Fatalf("Log: %v", err)
	}

	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Event != EventProfileSwitch {
		t.Errorf("event: got %q, want %q", e.Event, EventProfileSwitch)
	}
	if e.Project != "myproject" {
		t.Errorf("project: got %q, want %q", e.Project, "myproject")
	}
	if e.Detail != "dev" {
		t.Errorf("detail: got %q, want %q", e.Detail, "dev")
	}
}

func TestLog_MultipleEntries(t *testing.T) {
	dir := tempDir(t)
	l, _ := NewLogger(dir)

	events := []EventType{EventProfileCreate, EventVaultEncrypt, EventSnapshotSave}
	for _, ev := range events {
		if err := l.Log(ev, "proj", "detail"); err != nil {
			t.Fatalf("Log(%s): %v", ev, err)
		}
	}

	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != len(events) {
		t.Fatalf("expected %d entries, got %d", len(events), len(entries))
	}
}

func TestReadAll_EmptyLog(t *testing.T) {
	dir := tempDir(t)
	l, _ := NewLogger(dir)

	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestLog_TimestampIsUTC(t *testing.T) {
	dir := tempDir(t)
	l, _ := NewLogger(dir)
	before := time.Now().UTC()
	_ = l.Log(EventProfileDelete, "", "")
	after := time.Now().UTC()

	entries, _ := l.ReadAll()
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}
