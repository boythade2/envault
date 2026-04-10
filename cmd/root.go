package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "envault",
	Short: "A lightweight CLI tool for managing and encrypting environment variable files",
	Long: `envault helps you securely manage .env files across multiple projects.
Encrypt, decrypt, and share environment variables with confidence.`,
	Version: version,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("envault version %s\n", version))
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
