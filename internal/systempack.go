package internal

import "strings"

// IsSystemPack reports whether the specified pack is a Bedrock system pack.
//
// Only packs stored in the server pack directories can be treated as
// system packs.
func IsSystemPack(location PackLocation, folderName string) bool {
	if location != PackLocationServer {
		return false
	}

	name := strings.ToLower(strings.TrimSpace(folderName))

	switch {
	case strings.HasPrefix(name, "chemistry"):
		return true

	case name == "editor":
		return true

	case strings.HasPrefix(name, "experimental_"):
		return true

	case strings.HasPrefix(name, "server") &&
		strings.HasSuffix(name, "library"):
		return true

	case name == "vanilla":
		return true

	case strings.HasPrefix(name, "vanilla_"):
		return true
	}

	return false
}
