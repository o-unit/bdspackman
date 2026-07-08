package internal

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ResolveDisplayName resolves the display name of a pack.
//
// Search order:
//
//  1. texts/<language>.lang
//  2. texts/en_US.lang
//  3. manifest header.name
func ResolveDisplayName(pack *Pack, language string) string {
	if pack == nil {
		return ""
	}

	if pack.Name == "" {
		return ""
	}

	// Not a localization key.
	if !strings.Contains(pack.Name, ".") {
		return pack.Name
	}

	textDir := filepath.Join(pack.Path, "texts")

	for _, lang := range []string{language, "en_US"} {
		langFile := filepath.Join(textDir, lang+".lang")

		if value, ok := lookupLang(langFile, pack.Name); ok {
			return value
		}
	}

	return pack.Name
}

// lookupLang searches a key in a .lang file.
func lookupLang(path, key string) (string, bool) {
	file, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		name, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		if strings.TrimSpace(name) == key {
			return strings.TrimSpace(value), true
		}
	}

	return "", false
}