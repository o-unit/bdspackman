package tui

import (
	"github.com/o-unit/bdspackman/internal"

	tea "github.com/charmbracelet/bubbletea"
)

// Mode represents the current screen state.
type Mode int

const (
	// ModeNormal is the default pack list screen.
	ModeNormal Mode = iota
	ModeConfirmSave
	ModeConfirmDelete

	// (Future)
	// ModeSearch
	// ModeHelp
)

type Model struct {
	Config internal.Config
	Packs  []internal.Pack

	// Cursor position.
	Cursor int

	// Current UI mode.
	Mode Mode

	// Status message shown above the help line.
	StatusMessage StatusMessage

	// Help message shown at the bottom.
	HelpMessage string

	// Message for the currently selected pack.
	SelectionMessage StatusMessage
}

// NewModel creates a new TUI model.
func NewModel(cfg internal.Config, packs []internal.Pack) Model {

	model := Model{
		Config: cfg,
		Packs:  packs,

		Cursor: 0,
		Mode:   ModeNormal,

		StatusMessage:    StatusMessage{},
		HelpMessage:      "",
		SelectionMessage: StatusMessage{},
	}

	model.updateHelp()
	model.updateStatusFromSelection()
	return model
}

// Init is called when the program starts.
func (m Model) Init() tea.Cmd {
	return nil
}
