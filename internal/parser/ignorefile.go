package parser

import (
	"bufio"
	"os"
	"strings"
)

// ParseIgnoreFile reads a file containing one key per line (comments with #
// and blank lines are skipped) and returns the list of keys to ignore.
func ParseIgnoreFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var keys []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		keys = append(keys, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}
