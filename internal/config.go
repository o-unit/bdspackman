package internal

import "path/filepath"

// Config stores the runtime configuration.
type Config struct {
	// ServerDir is the root directory of the Bedrock Dedicated Server.
	ServerDir string

	// World is the world name specified by --world.
	World string

	// Language used when resolving localized pack names.
	Language string

	// ShowUUID controls whether UUIDs are displayed.
	ShowUUID bool

	// DryRun suppresses writing world json files.
	DryRun bool

	// ShowSystemPacks controls whether system packs are displayed.
	ShowSystemPacks bool
}

// WorldDir returns the absolute path to the selected world.
func (c Config) WorldDir() string {
	return filepath.Join(c.ServerDir, "worlds", c.World)
}

// ServerBehaviorPackDir returns the server behavior_packs directory.
func (c Config) ServerBehaviorPackDir() string {
	return filepath.Join(c.ServerDir, "behavior_packs")
}

// ServerResourcePackDir returns the server resource_packs directory.
func (c Config) ServerResourcePackDir() string {
	return filepath.Join(c.ServerDir, "resource_packs")
}

// WorldBehaviorPackDir returns the world's behavior_packs directory.
func (c Config) WorldBehaviorPackDir() string {
	return filepath.Join(c.WorldDir(), "behavior_packs")
}

// WorldResourcePackDir returns the world's resource_packs directory.
func (c Config) WorldResourcePackDir() string {
	return filepath.Join(c.WorldDir(), "resource_packs")
}

// WorldBehaviorJSON returns world_behavior_packs.json.
func (c Config) WorldBehaviorJSON() string {
	return filepath.Join(c.WorldDir(), "world_behavior_packs.json")
}

// WorldResourceJSON returns world_resource_packs.json.
func (c Config) WorldResourceJSON() string {
	return filepath.Join(c.WorldDir(), "world_resource_packs.json")
}
