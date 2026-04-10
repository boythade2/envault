package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	plaintext := []byte("SECRET_KEY=supersecret\nDB_PASS=hunter2")
	passphrase := "my-strong-passphrase"

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
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncryptProducesUniqueOutputs(t *testing.T) {
	plaintext := []byte("API_KEY=abc123")
	passphrase := "passphrase"

	first, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("first Encrypt failed: %v", err)
	}

	second, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("second Encrypt failed: %v", err)
	}

	if bytes.Equal(first, second) {
		t.Fatal("two encryptions of the same plaintext should produce different ciphertexts (random nonce)")
	}
}

func TestDecryptWithWrongPassphrase(t *testing.T) {
	plaintext := []byte("SECRET=value")

	ciphertext, err := Encrypt(plaintext, "correct-passphrase")
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt(ciphertext, "wrong-passphrase")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestDecryptShortCiphertext(t *testing.T) {
	_, err := Decrypt([]byte("short"), "passphrase")
	if err == nil {
		t.Fatal("expected error for short ciphertext")
	}
}
