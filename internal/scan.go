package internal

import (
	"os"
	"path/filepath"
	"sort"
)

// ScanPacks scans every pack directory and returns all discovered packs.
func ScanPacks(cfg Config) ([]Pack, error) {

	behaviorEnabled, err := LoadWorldPacks(cfg.WorldBehaviorJSON())
	if err != nil {
		return nil, err
	}

	resourceEnabled, err := LoadWorldPacks(cfg.WorldResourceJSON())
	if err != nil {
		return nil, err
	}

	behaviorSet := WorldPackSet(behaviorEnabled)
	resourceSet := WorldPackSet(resourceEnabled)

	var packs []Pack

	type target struct {
		Path     string
		Type     PackType
		Location PackLocation
		Enabled  map[string]struct{}
	}

	targets := []target{
		{
			Path:     cfg.ServerBehaviorPackDir(),
			Type:     PackTypeBehavior,
			Location: PackLocationServer,
			Enabled:  behaviorSet,
		},
		{
			Path:     cfg.ServerResourcePackDir(),
			Type:     PackTypeResource,
			Location: PackLocationServer,
			Enabled:  resourceSet,
		},
		{
			Path:     cfg.WorldBehaviorPackDir(),
			Type:     PackTypeBehavior,
			Location: PackLocationWorld,
			Enabled:  behaviorSet,
		},
		{
			Path:     cfg.WorldResourcePackDir(),
			Type:     PackTypeResource,
			Location: PackLocationWorld,
			Enabled:  resourceSet,
		},
	}

	for _, target := range targets {

		list, err := scanDirectory(
			target.Path,
			target.Type,
			target.Location,
			target.Enabled,
			cfg.Language,
			cfg.ShowSystemPacks,
		)

		if err != nil {
			return nil, err
		}

		packs = append(packs, list...)
	}

	ValidatePacks(packs)

	sortPacks(packs)

	return packs, nil
}

func scanDirectory(
	dir string,
	packType PackType,
	location PackLocation,
	enabled map[string]struct{},
	language string,
	showSystem bool,
) ([]Pack, error) {

	entries, err := os.ReadDir(dir)
	if err != nil {

		if os.IsNotExist(err) {
			return []Pack{}, nil
		}

		return nil, err
	}

	var packs []Pack

	for _, entry := range entries {

		if !entry.IsDir() {
			continue
		}

		pack, ok := buildPack(
			filepath.Join(dir, entry.Name()),
			entry.Name(),
			packType,
			location,
			enabled,
			language,
			showSystem,
		)

		if !ok {
			continue
		}

		packs = append(packs, pack)
	}

	return packs, nil
}

func buildPack(
	path string,
	folderName string,
	packType PackType,
	location PackLocation,
	enabled map[string]struct{},
	language string,
	showSystem bool,
) (Pack, bool) {

	pack := Pack{
		Type:       packType,
		Location:   location,
		FolderName: folderName,
		Path:       path,
		Status:     StatusOff,
	}

	if IsSystemPack(location, folderName) {

		if !showSystem {
			return pack, false
		}

		pack.Status = StatusSystem
	}

	manifest, err := LoadManifest(filepath.Join(path, "manifest.json"))
	if err != nil {

		pack.Status = StatusError
		pack.Error = err

		return pack, true
	}

	pack.UUID = manifest.Header.UUID
	pack.Version = manifest.Header.Version

	name, err := ResolveDisplayName(
		path,
		language,
		manifest.Header.Name,
	)

	if err != nil {
		pack.DisplayName = manifest.Header.Name
	} else {
		pack.DisplayName = name
	}

	if pack.Status != StatusSystem {

		if _, ok := enabled[pack.UUID]; ok {
			pack.Status = StatusOn
		} else {
			pack.Status = StatusOff
		}
	}

	return pack, true
}

func sortPacks(packs []Pack) {
	sort.SliceStable(packs, func(i, j int) bool {

		a := packs[i]
		b := packs[j]

		// Server → World
		if a.Location != b.Location {
			return a.Location < b.Location
		}

		// BP → RP
		if a.Type != b.Type {
			return a.Type < b.Type
		}

		// Display Name
		if a.DisplayName != b.DisplayName {
			return a.DisplayName < b.DisplayName
		}

		// Version
		if cmp := compareVersion(a.Version, b.Version); cmp != 0 {
			return cmp < 0
		}

		// Folder name
		if a.FolderName != b.FolderName {
			return a.FolderName < b.FolderName
		}

		return a.UUID < b.UUID
	})
}

func compareVersion(a, b []uint32) int {

	max := len(a)

	if len(b) > max {
		max = len(b)
	}

	for i := 0; i < max; i++ {

		var av uint32
		var bv uint32

		if i < len(a) {
			av = a[i]
		}

		if i < len(b) {
			bv = b[i]
		}

		if av < bv {
			return -1
		}

		if av > bv {
			return 1
		}
	}

	return 0
}

// findPackByUUID returns the first pack that matches the UUID.
// Returns nil if not found.
func findPackByUUID(
	packs []Pack,
	packType PackType,
	location PackLocation,
	uuid string,
) *Pack {

	for i := range packs {

		if packs[i].Type != packType {
			continue
		}

		if packs[i].Location != location {
			continue
		}

		if packs[i].UUID != uuid {
			continue
		}

		return &packs[i]
	}

	return nil
}

// updateEnabledStatus updates every pack according to the
// contents of world_behavior_packs.json and world_resource_packs.json.
func updateEnabledStatus(
	packs []Pack,
	behavior map[string]struct{},
	resource map[string]struct{},
) {

	for i := range packs {

		if packs[i].Status == StatusSystem ||
			packs[i].Status == StatusError {
			continue
		}

		switch packs[i].Type {

		case PackTypeBehavior:

			if _, ok := behavior[packs[i].UUID]; ok {
				packs[i].Status = StatusOn
			} else {
				packs[i].Status = StatusOff
			}

		case PackTypeResource:

			if _, ok := resource[packs[i].UUID]; ok {
				packs[i].Status = StatusOn
			} else {
				packs[i].Status = StatusOff
			}
		}
	}
}

// hasPackEnabled reports whether the UUID exists in the enabled pack set.
func hasPackEnabled(enabled map[string]struct{}, uuid string) bool {
	_, ok := enabled[uuid]
	return ok
}

// filterVisiblePacks removes packs that should not be displayed.
// At present this simply returns the original slice because system packs
// are filtered during scanning. This function exists to keep future
// filtering logic localized.
func filterVisiblePacks(packs []Pack) []Pack {
	return packs
}

// scanTargets returns the scan targets used by ScanPacks.
//
// This helper exists so additional scan locations can be added without
// changing the ScanPacks implementation.
func scanTargets(cfg Config, behaviorSet, resourceSet map[string]struct{}) []struct {
	Path     string
	Type     PackType
	Location PackLocation
	Enabled  map[string]struct{}
} {
	return []struct {
		Path     string
		Type     PackType
		Location PackLocation
		Enabled  map[string]struct{}
	}{
		{
			Path:     cfg.ServerBehaviorPackDir(),
			Type:     PackTypeBehavior,
			Location: PackLocationServer,
			Enabled:  behaviorSet,
		},
		{
			Path:     cfg.ServerResourcePackDir(),
			Type:     PackTypeResource,
			Location: PackLocationServer,
			Enabled:  resourceSet,
		},
		{
			Path:     cfg.WorldBehaviorPackDir(),
			Type:     PackTypeBehavior,
			Location: PackLocationWorld,
			Enabled:  behaviorSet,
		},
		{
			Path:     cfg.WorldResourcePackDir(),
			Type:     PackTypeResource,
			Location: PackLocationWorld,
			Enabled:  resourceSet,
		},
	}
}
