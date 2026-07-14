package internal

import "fmt"

// DeletePack removes a pack directory from disk.
func DeletePack(pack Pack) error {
	if pack.Path == "" {
		return fmt.Errorf("pack path is empty")
	}

	return RemovePath(pack.Path)
}
