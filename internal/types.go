package internal

// PackType represents the type of a pack.
type PackType int

const (
	PackTypeBehavior PackType = iota
	PackTypeResource
)

// PackLocation represents where a pack is stored.
type PackLocation int

const (
	PackLocationServer PackLocation = iota
	PackLocationWorld
)

// PackStatus represents the current state of a pack.
type PackStatus int

const (
	StatusOff PackStatus = iota
	StatusOn
	StatusError
	StatusSystem
)

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
