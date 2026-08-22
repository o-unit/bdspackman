package tui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CompletePath returns filesystem completion candidates for the given path.
func CompletePath(input string) ([]string, error) {

	// ~ はホームディレクトリへ展開
	if strings.HasPrefix(input, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			input = filepath.Join(home, strings.TrimPrefix(input, "~"))
		}
	}

	dir := input
	prefix := ""

	info, err := os.Stat(input)
	if err == nil && info.IsDir() {
		dir = input
	} else {
		dir = filepath.Dir(input)
		prefix = filepath.Base(input)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []string

	for _, entry := range entries {

		name := entry.Name()

		if !strings.HasPrefix(
			strings.ToLower(name),
			strings.ToLower(prefix),
		) {
			continue
		}

		full := filepath.Join(dir, name)

		if entry.IsDir() {
			full += string(os.PathSeparator)
		}

		result = append(result, full)
	}

	sort.Strings(result)

	return result, nil
}
