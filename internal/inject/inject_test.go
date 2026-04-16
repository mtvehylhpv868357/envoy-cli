package inject_test

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/inject"
)

func TestIntoMap_Overwrite(t *testing.T) {
	dst := map[string]string{"FOO": "old", "BAR": "keep"}
	src := map[string]string{"FOO": "new", "BAZ": "added"}
	opts := inject.DefaultOptions()
	result := inject.IntoMap(dst, src, opts)
	if result["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %s", result["FOO"])
	}
	if result["BAR"] != "keep" {
		t.Errorf("expected BAR=keep, got %s", result["BAR"])
	}
	if result["BAZ"] != "added" {
		t.Errorf("expected BAZ=added, got %s", result["BAZ"])
	}
}

func TestIntoMap_NoOverwrite(t *testing.T) {
	dst := map[string]string{"FOO": "original"}
	src := map[string]string{"FOO": "new"}
	opts := inject.Options{Overwrite: false}
	inject.IntoMap(dst, src, opts)
	if dst["FOO"] != "original" {
		t.Errorf("expected FOO to remain original, got %s", dst["FOO"])
	}
}

func TestIntoMap_Prefix(t *testing.T) {
	dst := map[string]string{}
	src := map[string]string{"KEY": "val"}
	opts := inject.Options{Overwrite: true, Prefix: "APP_"}
	inject.IntoMap(dst, src, opts)
	if dst["APP_KEY"] != "val" {
		t.Errorf("expected APP_KEY=val, got %v", dst)
	}
}

func TestIntoEnviron_Overwrite(t *testing.T) {
	environ := []string{"FOO=old", "BAR=keep"}
	src := map[string]string{"FOO": "new", "BAZ": "added"}
	opts := inject.DefaultOptions()
	result := inject.IntoEnviron(environ, src, opts)
	check := func(key, want string) {
		for _, e := range result {
			if strings.HasPrefix(e, key+"=") {
				got := strings.SplitN(e, "=", 2)[1]
				if got != want {
					t.Errorf("%s: want %s got %s", key, want, got)
				}
				return
			}
		}
		t.Errorf("%s not found in environ", key)
	}
	check("FOO", "new")
	check("BAR", "keep")
	check("BAZ", "added")
}

func TestIntoEnviron_NoOverwrite(t *testing.T) {
	environ := []string{"FOO=original"}
	src := map[string]string{"FOO": "new"}
	opts := inject.Options{Overwrite: false}
	result := inject.IntoEnviron(environ, src, opts)
	if result[0] != "FOO=original" {
		t.Errorf("expected FOO=original, got %s", result[0])
	}
}
