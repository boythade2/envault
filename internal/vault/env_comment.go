package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type CommentStore struct {
	Comments map[string]string `json:"comments"`
}

func commentPath(vaultFile string) string {
	dir := filepath.Dir(vaultFile)
	base := filepath.Base(vaultFile)
	return filepath.Join(dir, "."+base+".comments.json")
}

func LoadComments(vaultFile string) (*CommentStore, error) {
	path := commentPath(vaultFile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &CommentStore{Comments: map[string]string{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read comments: %w", err)
	}
	var cs CommentStore
	if err := json.Unmarshal(data, &cs); err != nil {
		return nil, fmt.Errorf("parse comments: %w", err)
	}
	if cs.Comments == nil {
		cs.Comments = map[string]string{}
	}
	return &cs, nil
}

func saveComments(vaultFile string, cs *CommentStore) error {
	data, err := json.MarshalIndent(cs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal comments: %w", err)
	}
	return os.WriteFile(commentPath(vaultFile), data, 0600)
}

func SetComment(vaultFile, key, comment string) error {
	cs, err := LoadComments(vaultFile)
	if err != nil {
		return err
	}
	cs.Comments[key] = comment
	return saveComments(vaultFile, cs)
}

func RemoveComment(vaultFile, key string) error {
	cs, err := LoadComments(vaultFile)
	if err != nil {
		return err
	}
	delete(cs.Comments, key)
	return saveComments(vaultFile, cs)
}

func GetComment(cs *CommentStore, key string) string {
	return cs.Comments[key]
}
