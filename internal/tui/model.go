package tui

import (
	"github.com/o-unit/bdspackman/internal"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Mode represents the current screen state.
type Mode int

const (
	// ModeNormal is the default pack list screen.
	ModeNormal Mode = iota
	ModeConfirmSave
	ModeConfirmDelete
	ModeRenameDir
	ModeAddInput

	// (Future)
	// ModeSearch
	// ModeHelp
)

type Model struct {
	Config internal.Config
	Packs  []internal.Pack

	// Cursor position.
	Cursor int

	// Viewport width.
	ViewportWidth int

	// Current UI mode.
	Mode Mode

	// Status message shown above the help line.
	StatusMessage StatusMessage

	// Help message shown at the bottom.
	HelpMessage string

	// Message for the currently selected pack.
	SelectionMessage StatusMessage

	// Current directory name input for rename mode.
	RenameInput string

	// ----------------------------------------------------------------
	// Add Pack
	// ----------------------------------------------------------------

	// true: Worldへ追加
	// false: Serverへ追加
	AddToWorld bool

	// Path input
	AddInput textinput.Model

	// Current completion candidates.
	AddCompletions []string

	// Current completion index.
	AddCompletionIndex int

	// Input used to build the current completion list.
	LastCompletionInput string
}

// NewModel creates a new TUI model.
// NewModel creates a new TUI model.
func NewModel(cfg internal.Config, packs []internal.Pack) Model {

	model := Model{
		Config: cfg,
		Packs:  packs,

		Cursor: 0,
		Mode:   ModeNormal,

		StatusMessage:       StatusMessage{},
		HelpMessage:         "",
		SelectionMessage:    StatusMessage{},
		RenameInput:         "",
		AddToWorld:          false,
		AddCompletions:      nil,
		AddCompletionIndex:  0,
		LastCompletionInput: "",
	}

	model.updateHelp()
	model.updateStatusFromSelection()

	ti := textinput.New()
	ti.Placeholder = "Path (.mcaddon/.mcpack/.zip or directory)"
	ti.CharLimit = 4096
	ti.Width = 40

	model.AddInput = ti

	return model
}

// Init is called when the program starts.
func (m Model) Init() tea.Cmd {
	return nil
}
