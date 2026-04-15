package compare_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/compare"
	"github.com/user/envoy-cli/internal/profile"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "compare-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeStore(t *testing.T) *profile.Store {
	t.Helper()
	store, err := profile.LoadStore(filepath.Join(tempDir(t), "profiles"))
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func TestCompare_OnlyInA(t *testing.T) {
	store := makeStore(t)
	store.Add("dev", map[string]string{"FOO": "bar", "ONLY_A": "yes"})
	store.Add("prod", map[string]string{"FOO": "bar"})

	res, err := compare.Profiles(store, "dev", "prod")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.OnlyInA["ONLY_A"]; !ok {
		t.Error("expected ONLY_A in OnlyInA")
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	store := makeStore(t)
	store.Add("dev", map[string]string{"FOO": "bar"})
	store.Add("prod", map[string]string{"FOO": "bar", "ONLY_B": "yes"})

	res, err := compare.Profiles(store, "dev", "prod")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.OnlyInB["ONLY_B"]; !ok {
		t.Error("expected ONLY_B in OnlyInB")
	}
}

func TestCompare_Differ(t *testing.T) {
	store := makeStore(t)
	store.Add("dev", map[string]string{"DB": "localhost"})
	store.Add("prod", map[string]string{"DB": "prod-host"})

	res, err := compare.Profiles(store, "dev", "prod")
	if err != nil {
		t.Fatal(err)
	}
	pair, ok := res.Differ["DB"]
	if !ok {
		t.Fatal("expected DB in Differ")
	}
	if pair[0] != "localhost" || pair[1] != "prod-host" {
		t.Errorf("unexpected values: %v", pair)
	}
}

func TestCompare_Same(t *testing.T) {
	store := makeStore(t)
	store.Add("dev", map[string]string{"SHARED": "value"})
	store.Add("prod", map[string]string{"SHARED": "value"})

	res, err := compare.Profiles(store, "dev", "prod")
	if err != nil {
		t.Fatal(err)
	}
	if res.Same["SHARED"] != "value" {
		t.Error("expected SHARED in Same")
	}
}

func TestCompare_ProfileNotFound(t *testing.T) {
	store := makeStore(t)
	store.Add("dev", map[string]string{"X": "1"})

	_, err := compare.Profiles(store, "dev", "missing")
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestAllKeys_Sorted(t *testing.T) {
	store := makeStore(t)
	store.Add("a", map[string]string{"Z": "1", "A": "2"})
	store.Add("b", map[string]string{"M": "3", "A": "2"})

	res, _ := compare.Profiles(store, "a", "b")
	keys := res.AllKeys()
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}
