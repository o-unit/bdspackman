package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadServerProperties reads server.properties and returns all key/value pairs.
func LoadServerProperties(serverDir string) (map[string]string, error) {

	path := filepath.Join(serverDir, "server.properties")

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	props := make(map[string]string)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		props[strings.TrimSpace(key)] =
			strings.TrimSpace(value)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf(
			"read server.properties: %w",
			err,
		)
	}

	return props, nil
}
