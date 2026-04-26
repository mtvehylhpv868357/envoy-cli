package envsign_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/envsign"
)

func TestSign_ProducesEnvelope(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	env, err := envSign.Sign("dev", vars, "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Profile != "dev" {
		t.Errorf("expected profile 'dev', got %q", env.Profile)
	}
	if env.Signature == "" {
		t.Error("expected non-empty signature")
	}
}

func TestVerify_ValidSignature(t *testing.T) {
	vars := map[string]string{"KEY": "value"}
	env, err := envSign.Sign("prod", vars, "mypassphrase")
	if err != nil {
		t.Fatalf("sign error: %v", err)
	}
	if err := envSign.Verify(env, "mypassphrase"); err != nil {
		t.Errorf("expected valid signature, got: %v", err)
	}
}

func TestVerify_WrongPassphrase(t *testing.T) {
	vars := map[string]string{"KEY": "value"}
	env, _ := envSign.Sign("prod", vars, "correct")
	err := envSign.Verify(env, "wrong")
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
	if err != envSign.ErrInvalidSignature {
		t.Errorf("expected ErrInvalidSignature, got: %v", err)
	}
}

func TestVerify_TamperedVars(t *testing.T) {
	vars := map[string]string{"KEY": "value"}
	env, _ := envSign.Sign("dev", vars, "secret")
	env.Vars["KEY"] = "tampered"
	err := envSign.Verify(env, "secret")
	if err != envSign.ErrInvalidSignature {
		t.Errorf("expected ErrInvalidSignature after tampering, got: %v", err)
	}
}

func TestSign_DeterministicSignature(t *testing.T) {
	vars := map[string]string{"B": "2", "A": "1"}
	e1, _ := envSign.Sign("test", vars, "pass")
	e2, _ := envSign.Sign("test", vars, "pass")
	if e1.Signature != e2.Signature {
		t.Error("expected same signature for same input")
	}
}

func TestSign_EmptyPassphrase_Errors(t *testing.T) {
	_, err := envSign.Sign("dev", map[string]string{}, "")
	if err == nil {
		t.Error("expected error for empty passphrase")
	}
}
