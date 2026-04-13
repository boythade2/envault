package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommandHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})

	// Execute should not return an error for --help
	_ = rootCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "envault") {
		t.Errorf("expected help output to contain 'envault', got: %s", output)
	}
}

func TestRootCommandVersion(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--version"})

	_ = rootCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, version) {
		t.Errorf("expected version output to contain %q, got: %s", version, output)
	}
}

func TestExecuteReturnsNoError(t *testing.T) {
	rootCmd.SetArgs([]string{})
	if err := Execute(); err != nil {
		t.Errorf("Execute() returned unexpected error: %v", err)
	}
}

func TestRootCommandUnknownFlag(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--unknown-flag"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected an error for unknown flag, got nil")
	}
}
