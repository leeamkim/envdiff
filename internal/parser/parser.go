package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a parsed .env file as a map of key-value pairs.
type EnvMap map[string]string

// ParseFile reads a .env file and returns an EnvMap.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid syntax at line %d: %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, fmt.Errorf("empty key at line %d: %q", lineNum, line)
		}

		value := strings.TrimSpace(parts[1])
		value = stripQuotes(value)

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file %q: %w", path, err)
	}

	return env, nil
}

// stripQuotes removes a matching pair of surrounding single or double quotes
// from a value string, if present.
func stripQuotes(value string) string {
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			return value[1 : len(value)-1]
		}
	}
	return value
}
