package envcast_test

import (
	"testing"
	"time"

	"github.com/your-org/envoy-cli/internal/envcast"
)

func TestToBool_TrueValues(t *testing.T) {
	for _, v := range []string{"true", "1", "yes", "TRUE", "YES"} {
		got, err := envcast.ToBool(v)
		if err != nil || !got {
			t.Errorf("ToBool(%q) = %v, %v; want true, nil", v, got, err)
		}
	}
}

func TestToBool_FalseValues(t *testing.T) {
	for _, v := range []string{"false", "0", "no", "FALSE", "NO"} {
		got, err := envcast.ToBool(v)
		if err != nil || got {
			t.Errorf("ToBool(%q) = %v, %v; want false, nil", v, got, err)
		}
	}
}

func TestToBool_Invalid(t *testing.T) {
	_, err := envcast.ToBool("maybe")
	if err == nil {
		t.Fatal("expected error for invalid bool value")
	}
}

func TestToInt_Valid(t *testing.T) {
	got, err := envcast.ToInt("42")
	if err != nil || got != 42 {
		t.Fatalf("ToInt(\"42\") = %d, %v; want 42, nil", got, err)
	}
}

func TestToInt_Invalid(t *testing.T) {
	_, err := envcast.ToInt("abc")
	if err == nil {
		t.Fatal("expected error for non-integer string")
	}
}

func TestToFloat64_Valid(t *testing.T) {
	got, err := envcast.ToFloat64("3.14")
	if err != nil || got != 3.14 {
		t.Fatalf("ToFloat64 = %f, %v; want 3.14, nil", got, err)
	}
}

func TestToFloat64_Invalid(t *testing.T) {
	_, err := envcast.ToFloat64("not-a-float")
	if err == nil {
		t.Fatal("expected error for non-float string")
	}
}

func TestToDuration_Valid(t *testing.T) {
	got, err := envcast.ToDuration("2m30s")
	if err != nil || got != 2*time.Minute+30*time.Second {
		t.Fatalf("ToDuration = %v, %v; want 2m30s, nil", got, err)
	}
}

func TestToDuration_Invalid(t *testing.T) {
	_, err := envcast.ToDuration("forever")
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}

func TestMap_AllValid(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2", "C": "3"}
	out, err := envcast.Map(vars, envcast.ToInt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["B"] != 2 {
		t.Errorf("expected out[B]=2, got %d", out["B"])
	}
}

func TestMap_PartialError(t *testing.T) {
	vars := map[string]string{"X": "10", "Y": "bad"}
	_, err := envcast.Map(vars, envcast.ToInt)
	if err == nil {
		t.Fatal("expected error for partial conversion failure")
	}
}
