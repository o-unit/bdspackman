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

	applyWorldOrder(
		packs,
		behaviorEnabled,
		resourceEnabled,
	)

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

	pack.Name = manifest.Header.Name
	pack.DisplayName = ResolveDisplayName(&pack, language)

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

		// Status
		if a.Status != b.Status {
			return a.Status < b.Status
		}

		// ON同士は現在の順番を維持する
		if a.Status == StatusOn {
			return false
		}

		// World → Server
		if a.Location != b.Location {
			return a.Location > b.Location
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

// applyWorldOrder reorders enabled packs to match the order stored in
// world_behavior_packs.json and world_resource_packs.json.
//
// Disabled packs remain after enabled packs.
// The relative order of disabled packs is preserved.
func applyWorldOrder(
	packs []Pack,
	behavior []WorldPack,
	resource []WorldPack,
) {

	var ordered []Pack
	used := make([]bool, len(packs))

	// Behavior Packs
	appendOrderedPacks(
		&ordered,
		used,
		packs,
		PackTypeBehavior,
		behavior,
	)

	// Resource Packs
	appendOrderedPacks(
		&ordered,
		used,
		packs,
		PackTypeResource,
		resource,
	)

	// Remaining packs (OFF / ERROR / SYSTEM)
	for i := range packs {

		if used[i] {
			continue
		}

		ordered = append(ordered, packs[i])
	}

	copy(packs, ordered)
}

func appendOrderedPacks(
	dst *[]Pack,
	used []bool,
	packs []Pack,
	packType PackType,
	world []WorldPack,
) {

	for _, wp := range world {

		// Worldを優先
		for location := PackLocationWorld; location >= PackLocationServer; location-- {

			for i := range packs {

				if used[i] {
					continue
				}

				p := packs[i]

				if p.Type != packType {
					continue
				}

				if p.Location != location {
					continue
				}

				if p.UUID != wp.PackID {
					continue
				}

				*dst = append(*dst, p)
				used[i] = true

				break
			}
		}
	}
}
