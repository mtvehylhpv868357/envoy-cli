package importer

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestFromBytes_DotEnv(t *testing.T) {
	data := []byte("APP_ENV=production\nDEBUG=false\n# comment\n")
	res, err := FromBytes(data, FormatDotEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", res.Vars["APP_ENV"])
	}
	if res.Vars["DEBUG"] != "false" {
		t.Errorf("expected DEBUG=false, got %q", res.Vars["DEBUG"])
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected no skipped lines, got %v", res.Skipped)
	}
}

func TestFromBytes_ExportFormat(t *testing.T) {
	data := []byte("export HOST=localhost\nexport PORT=8080\n")
	res, err := FromBytes(data, FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", res.Vars["HOST"])
	}
	if res.Format != FormatExport {
		t.Errorf("expected format=export, got %q", res.Format)
	}
}

func TestFromBytes_AutoDetect_Export(t *testing.T) {
	data := []byte("export KEY=value\n")
	res, err := FromBytes(data, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", res.Vars["KEY"])
	}
}

func TestFromBytes_SkipsInvalidLines(t *testing.T) {
	data := []byte("VALID=yes\nNOEQUALSSIGN\n")
	res, err := FromBytes(data, FormatDotEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "NOEQUALSIGN" {
		t.Errorf("expected 1 skipped line, got %v", res.Skipped)
	}
}

func TestFromFile_ReadsFile(t *testing.T) {
	path := writeTempFile(t, "FROM_FILE=yes\n")
	res, err := FromFile(path, FormatDotEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["FROM_FILE"] != "yes" {
		t.Errorf("expected FROM_FILE=yes, got %q", res.Vars["FROM_FILE"])
	}
}

func TestFromFile_NotFound(t *testing.T) {
	_, err := FromFile("/nonexistent/.env", FormatDotEnv)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestFromBytes_UnsupportedFormat(t *testing.T) {
	_, err := FromBytes([]byte("K=V"), FormatJSON)
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}
