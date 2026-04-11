package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuditEvent represents a single recorded action on a vault.
type AuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Key       string    `json:"key,omitempty"`
	VaultPath string    `json:"vault_path"`
	Details   string    `json:"details,omitempty"`
}

// AuditLog is a collection of audit events.
type AuditLog struct {
	Events []AuditEvent `json:"events"`
}

func auditPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".audit.json")
}

// LoadAuditLog loads the audit log for the given vault file.
// Returns an empty log if the file does not exist.
func LoadAuditLog(vaultPath string) (*AuditLog, error) {
	path := auditPath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &AuditLog{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading audit log: %w", err)
	}
	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("parsing audit log: %w", err)
	}
	return &log, nil
}

// RecordEvent appends a new event to the audit log for the given vault.
func RecordEvent(vaultPath, action, key, details string) error {
	log, err := LoadAuditLog(vaultPath)
	if err != nil {
		return err
	}
	log.Events = append(log.Events, AuditEvent{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Key:       key,
		VaultPath: vaultPath,
		Details:   details,
	})
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding audit log: %w", err)
	}
	path := auditPath(vaultPath)
	return os.WriteFile(path, data, 0600)
}
