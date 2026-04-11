package vault

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ImportDotenv reads a .env file and returns a map of key-value pairs.
func ImportDotenv(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		result[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file: %w", err)
	}
	return result, nil
}

// ImportJSON reads a flat JSON object file and returns a map of key-value pairs.
func ImportJSON(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	result := make(map[string]string, len(raw))
	for k, v := range raw {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
