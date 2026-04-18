package archive_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-cli/internal/archive"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "archive-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func writeProfile(t *testing.T, dir, name string, vars map[string]string) {
	t.Helper()
	data, _ := json.Marshal(vars)
	if err := os.WriteFile(filepath.Join(dir, name+".json"), data, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestPack_And_Unpack_RoundTrip(t *testing.T) {
	src := tempDir(t)
	writeProfile(t, src, "dev", map[string]string{"FOO": "bar", "PORT": "8080"})
	writeProfile(t, src, "prod", map[string]string{"FOO": "baz"})

	var buf bytes.Buffer
	if err := archive.Pack(src, []string{"dev", "prod"}, &buf); err != nil {
		t.Fatalf("Pack: %v", err)
	}

	dst := tempDir(t)
	names, err := archive.Unpack(&buf, dst, false)
	if err != nil {
		t.Fatalf("Unpack: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(names))
	}

	data, _ := os.ReadFile(filepath.Join(dst, "dev.json"))
	var vars map[string]string
	json.Unmarshal(data, &vars)
	if vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", vars["FOO"])
	}
}

func TestUnpack_NoOverwrite_Errors(t *testing.T) {
	src := tempDir(t)
	writeProfile(t, src, "dev", map[string]string{"A": "1"})

	var buf bytes.Buffer
	archive.Pack(src, []string{"dev"}, &buf)

	dst := tempDir(t)
	writeProfile(t, dst, "dev", map[string]string{"A": "existing"})

	_, err := archive.Unpack(&buf, dst, false)
	if err == nil {
		t.Fatal("expected error when overwrite=false and file exists")
	}
}

func TestUnpack_Overwrite_Succeeds(t *testing.T) {
	src := tempDir(t)
	writeProfile(t, src, "dev", map[string]string{"A": "new"})

	var buf bytes.Buffer
	archive.Pack(src, []string{"dev"}, &buf)

	dst := tempDir(t)
	writeProfile(t, dst, "dev", map[string]string{"A": "old"})

	_, err := archive.Unpack(&buf, dst, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(dst, "dev.json"))
	var vars map[string]string
	json.Unmarshal(data, &vars)
	if vars["A"] != "new" {
		t.Errorf("expected A=new after overwrite, got %q", vars["A"])
	}
}

func TestPack_MissingProfile_Errors(t *testing.T) {
	src := tempDir(t)
	var buf bytes.Buffer
	err := archive.Pack(src, []string{"nonexistent"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}
