package encrypt

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plaintext := []byte("SECRET_KEY=super_secret_value")
	passphrase := "my-test-passphrase"

	ciphertext, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := Decrypt(ciphertext, passphrase)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncrypt_DifferentNonceEachTime(t *testing.T) {
	plaintext := []byte("DB_PASSWORD=hunter2")
	passphrase := "passphrase"

	c1, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("first Encrypt failed: %v", err)
	}
	c2, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("second Encrypt failed: %v", err)
	}

	if bytes.Equal(c1, c2) {
		t.Error("two encryptions of the same plaintext should produce different ciphertexts")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	plaintext := []byte("API_TOKEN=abc123")

	ciphertext, err := Encrypt(plaintext, "correct-passphrase")
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt(ciphertext, "wrong-passphrase")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestDecrypt_ShortCiphertext(t *testing.T) {
	_, err := Decrypt([]byte("short"), "passphrase")
	if err == nil {
		t.Fatal("expected error for short ciphertext")
	}
}
