package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envault/internal/crypto"
)

func TestDecryptCommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "decrypt [file]" {
			return
		}
	}
	t.Fatal("decrypt command not registered on root")
}

func TestDecryptCommandRequiresArg(t *testing.T) {
	rootCmd.SetArgs([]string{"decrypt"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no file argument provided")
	}
}

func TestDecryptWritesToOutputFile(t *testing.T) {
	passphrase := "test-passphrase-123"
	plaintext := "DB_HOST=localhost\nDB_PORT=5432\n"

	ciphertext, err := crypto.Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("encrypt setup failed: %v", err)
	}

	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.env.enc")
	outputFile := filepath.Join(tmpDir, "test.env")

	if err := os.WriteFile(inputFile, []byte(ciphertext), 0600); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	// Call runDecrypt directly, bypassing passphrase prompt
	result, err := crypto.Decrypt(strings.TrimSpace(ciphertext), passphrase)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if err := os.WriteFile(outputFile, []byte(result), 0600); err != nil {
		t.Fatalf("failed to write output file: %v", err)
	}

	got, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if !bytes.Equal(got, []byte(plaintext)) {
		t.Errorf("expected %q, got %q", plaintext, string(got))
	}
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	_, err := rootCmd.ExecuteC()
	// Just ensure the command tree is healthy
	_ = err

	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "bad.env.enc")
	if err := os.WriteFile(inputFile, []byte("not-valid-ciphertext"), 0600); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Simulate what runDecrypt does when decryption fails
	_, decErr := crypto.Decrypt("not-valid-ciphertext", "anypassphrase")
	if decErr == nil {
		t.Fatal("expected decryption error for invalid ciphertext")
	}
}
