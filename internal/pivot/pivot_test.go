package pivot

import (
	"testing"
)

func TestPivot_BasicRows(t *testing.T) {
	profiles := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "PORT": "8080"},
		"prod": {"HOST": "example.com", "PORT": "443"},
	}
	rows := Profiles(profiles, DefaultOptions())
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].Key != "HOST" {
		t.Errorf("expected first key HOST, got %s", rows[0].Key)
	}
}

func TestPivot_MissingKey(t *testing.T) {
	profiles := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "DEBUG": "true"},
		"prod": {"HOST": "example.com"},
	}
	rows := Profiles(profiles, DefaultOptions())
	var debugRow *Row
	for i := range rows {
		if rows[i].Key == "DEBUG" {
			debugRow = &rows[i]
		}
	}
	if debugRow == nil {
		t.Fatal("expected DEBUG row")
	}
	if len(debugRow.Missing) != 1 || debugRow.Missing[0] != "prod" {
		t.Errorf("expected prod missing, got %v", debugRow.Missing)
	}
}

func TestPivot_EmptyProfiles(t *testing.T) {
	rows := Profiles(map[string]map[string]string{}, DefaultOptions())
	if len(rows) != 0 {
		t.Errorf("expected empty rows")
	}
}

func TestPivot_SortKeys(t *testing.T) {
	profiles := map[string]map[string]string{
		"dev": {"ZEBRA": "z", "APPLE": "a", "MANGO": "m"},
	}
	rows := Profiles(profiles, Options{SortKeys: true})
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, r := range rows {
		if r.Key != expected[i] {
			t.Errorf("index %d: expected %s got %s", i, expected[i], r.Key)
		}
	}
}

func TestPivot_ValuesCorrect(t *testing.T) {
	profiles := map[string]map[string]string{
		"dev":  {"URL": "http://dev"},
		"prod": {"URL": "http://prod"},
	}
	rows := Profiles(profiles, DefaultOptions())
	if len(rows) != 1 {
		t.Fatalf("expected 1 row")
	}
	if rows[0].Values["dev"] != "http://dev" {
		t.Errorf("wrong dev value: %s", rows[0].Values["dev"])
	}
	if rows[0].Values["prod"] != "http://prod" {
		t.Errorf("wrong prod value: %s", rows[0].Values["prod"])
	}
}
