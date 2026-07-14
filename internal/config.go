package internal

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config stores program configuration.
type Config struct {
	Version         bool
	ServerDir       string
	WorldName       string
	Language        string
	ShowUUID        bool
	ShowDirName     bool
	ShowSystemPacks bool
	ExportPrefix    string
	ExportDir       string
	ForceDelete     bool
}

// LoadConfig parses command line arguments.
func LoadFlags() (Config, error) {

	cfg := Config{}

	flag.BoolVar(
		&cfg.Version,
		"version",
		false,
		"show version",
	)

	flag.StringVar(
		&cfg.ServerDir,
		"serverdir",
		".",
		"Bedrock Dedicated Server directory",
	)

	flag.StringVar(
		&cfg.WorldName,
		"world",
		"",
		"World name. If omitted, level-name in server.properties is used.",
	)

	flag.StringVar(
		&cfg.Language,
		"lang",
		"ja_JP",
		"Language code",
	)

	flag.BoolVar(
		&cfg.ShowUUID,
		"uuid",
		false,
		"Show UUID column",
	)

	flag.BoolVar(
		&cfg.ShowDirName,
		"dirname",
		false,
		"Show directory name column",
	)

	flag.BoolVar(
		&cfg.ShowSystemPacks,
		"systempack",
		false,
		"Show system packs",
	)

	flag.StringVar(
		&cfg.ExportPrefix,
		"export-prefix",
		"",
		"Prefix added to exported filenames",
	)

	flag.StringVar(
		&cfg.ExportDir,
		"export-dir",
		"",
		"Directory to export world pack json files",
	)

	flag.BoolVar(
		&cfg.ForceDelete,
		"force-delete",
		false,
		"Delete packs permanently instead of moving them to delpacks.",
	)

	flag.Parse()

	return cfg, nil
}

func CompleteConfig(cfg Config) (Config, error) {
	if cfg.WorldName == "" {

		props, err := LoadServerProperties(cfg.ServerDir)
		if err != nil {
			return cfg, fmt.Errorf(
				"cannot determine world name: %w",
				err,
			)
		}

		world := props["level-name"]
		if world == "" {
			return cfg, fmt.Errorf(
				"cannot determine world name: level-name not found",
			)
		}

		cfg.WorldName = world
	}

	serverDir, err := filepath.Abs(cfg.ServerDir)
	if err != nil {
		return cfg, err
	}

	cfg.ServerDir = serverDir

	if _, err := os.Stat(cfg.ServerDir); err != nil {
		return cfg, fmt.Errorf("server directory not found: %s", cfg.ServerDir)
	}

	if _, err := os.Stat(cfg.WorldDir()); err != nil {
		return cfg, fmt.Errorf("world directory not found: %s", cfg.WorldDir())
	}

	return cfg, nil
}

// WorldDir returns the world directory.
func (c Config) WorldDir() string {
	return filepath.Join(
		c.ServerDir,
		"worlds",
		c.WorldName,
	)
}

// ServerBehaviorPackDir returns the server behavior pack directory.
func (c Config) ServerBehaviorPackDir() string {
	return filepath.Join(
		c.ServerDir,
		"behavior_packs",
	)
}

// ServerResourcePackDir returns the server resource pack directory.
func (c Config) ServerResourcePackDir() string {
	return filepath.Join(
		c.ServerDir,
		"resource_packs",
	)
}

// WorldBehaviorPackDir returns the world behavior pack directory.
func (c Config) WorldBehaviorPackDir() string {
	return filepath.Join(
		c.WorldDir(),
		"behavior_packs",
	)
}

// WorldResourcePackDir returns the world resource pack directory.
func (c Config) WorldResourcePackDir() string {
	return filepath.Join(
		c.WorldDir(),
		"resource_packs",
	)
}

// WorldBehaviorJSON returns world_behavior_packs.json.
func (c Config) WorldBehaviorJSON() string {
	return filepath.Join(
		c.WorldDir(),
		"world_behavior_packs.json",
	)
}

// WorldResourceJSON returns world_resource_packs.json.
func (c Config) WorldResourceJSON() string {
	return filepath.Join(
		c.WorldDir(),
		"world_resource_packs.json",
	)
}
