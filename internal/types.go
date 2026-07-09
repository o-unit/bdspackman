package internal

import "fmt"

// PackType represents the type of a pack.
type PackType int

const (
	PackTypeBehavior PackType = iota
	PackTypeResource
)

func (t PackType) String() string {
	switch t {
	case PackTypeBehavior:
		return "BP"
	case PackTypeResource:
		return "RP"
	default:
		return "UNKNOWN"
	}
}

// PackLocation represents where a pack is stored.
type PackLocation int

const (
	PackLocationServer PackLocation = iota
	PackLocationWorld
)

func (l PackLocation) String() string {
	switch l {
	case PackLocationServer:
		return "Server"
	case PackLocationWorld:
		return "World"
	default:
		return "Unknown"
	}
}

// PackStatus represents the current state of a pack.
type PackStatus int

const (
	StatusOff PackStatus = iota
	StatusOn
	StatusError
	StatusSystem
)

func (s PackStatus) String() string {
	switch s {
	case StatusOff:
		return "OFF"
	case StatusOn:
		return "ON"
	case StatusError:
		return "ERROR"
	case StatusSystem:
		return "SYSTEM"
	default:
		return "UNKNOWN"
	}
}

// Pack contains all information about a behavior pack or resource pack.
type Pack struct {
	// Display information
	Name        string // header.name
	DisplayName string // Localized display name
	FolderName  string // Directory name

	// Manifest information
	UUID    string
	Version []uint32

	// Classification
	Type     PackType
	Location PackLocation
	Status   PackStatus

	// File system
	Path         string
	ManifestPath string

	// Error information
	Error error
}

func (p Pack) VersionString() string {
	if len(p.Version) != 3 {
		return ""
	}

	return fmt.Sprintf(
		"%d.%d.%d",
		p.Version[0],
		p.Version[1],
		p.Version[2],
	)
}

// Manifest represents manifest.json.
type Manifest struct {
	Header ManifestHeader `json:"header"`
}

// ManifestHeader represents manifest.header.
type ManifestHeader struct {
	Name    string   `json:"name"`
	UUID    string   `json:"uuid"`
	Version []uint32 `json:"version"`
}

// WorldPack represents one entry in
// world_behavior_packs.json or world_resource_packs.json.
type WorldPack struct {
	PackID  string   `json:"pack_id"`
	Version []uint32 `json:"version"`
}
