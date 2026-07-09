package internal

import "fmt"

// ValidatePacks validates all scanned packs.
//
// Returns true if one or more validation errors were found.
func ValidatePacks(packs []Pack) bool {
	validateDuplicateUUIDs(packs)

	for _, pack := range packs {
		if pack.Status == StatusError {
			return true
		}
	}

	return false
}

// validateDuplicateUUIDs marks duplicated UUIDs as StatusError.
//
// Duplicate UUIDs are checked independently for each:
//
//   - Behavior Packs (Server)
//   - Behavior Packs (World)
//   - Resource Packs (Server)
//   - Resource Packs (World)
func validateDuplicateUUIDs(packs []Pack) {
	type key struct {
		PackType
		PackLocation
	}

	uuidMap := make(map[key]map[string]int)

	for i := range packs {

		// Skip packs that already have an error.
		if packs[i].Status == StatusError {
			continue
		}

		// System packs are ignored.
		if packs[i].Status == StatusSystem {
			continue
		}

		k := key{
			PackType:     packs[i].Type,
			PackLocation: packs[i].Location,
		}

		if _, ok := uuidMap[k]; !ok {
			uuidMap[k] = make(map[string]int)
		}

		if first, exists := uuidMap[k][packs[i].UUID]; exists {

			// Mark both packs as invalid.

			packs[first].Status = StatusError
			packs[first].Error = fmt.Errorf("duplicate UUID")

			packs[i].Status = StatusError
			packs[i].Error = fmt.Errorf("duplicate UUID")

			continue
		}

		uuidMap[k][packs[i].UUID] = i
	}
}
