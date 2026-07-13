package tui

import "github.com/o-unit/bdspackman/internal"

// setStatus updates the status bar message.
func (m *Model) setStatus(message string) {
	m.StatusMessage = message
}

func (m *Model) setHelp(message string) {
	m.HelpMessage = message
}

// updateHelp updates the help message according to the current mode.
func (m *Model) updateHelp() {

	switch m.Mode {

	case ModeNormal:
		m.HelpMessage =
			"↑↓ Move  Ctrl+↑↓ Reorder Space Toggle  Esc Quit"

	case ModeConfirmSave:
		m.HelpMessage =
			"Y Save  N Cancel"
	default:
		m.HelpMessage = ""
	}
}

// updateStatusFromSelection updates the message for the selected pack.
func (m *Model) updateStatusFromSelection() {

	if len(m.Packs) == 0 {
		m.SelectionMessage = ""
		return
	}

	pack := m.Packs[m.Cursor]

	switch pack.Status {

	case internal.StatusError:

		if pack.Error != nil {
			m.SelectionMessage = pack.Error.Error()
		} else {
			m.SelectionMessage = "Unknown error."
		}

	case internal.StatusSystem:
		m.SelectionMessage = "System pack"

	case internal.StatusOn:
		m.SelectionMessage = "Enabled"

	case internal.StatusOff:
		m.SelectionMessage = "Disabled"

	default:
		m.SelectionMessage = ""
	}
}

// enterSaveConfirm switches to save confirmation mode.
func (m *Model) enterSaveConfirm() {

	m.Mode = ModeConfirmSave

	m.setStatus("Save changes? (Y=yes / N=no)")
	m.updateHelp()
}

// enterNormalMode switches back to normal mode.
func (m *Model) enterNormalMode() {

	m.Mode = ModeNormal
	m.updateHelp()
}

// save writes the current pack configuration.
func (m *Model) save() {

	message, err := internal.SaveWorldPackFiles(
		m.Config,
		m.Packs,
		internal.SaveOptions{
			ExportPrefix: m.Config.ExportPrefix,
			ExportDir:    m.Config.ExportDir,
		},
	)

	m.enterNormalMode()

	if err != nil {
		m.setStatus("Save failed: " + err.Error())
		return
	}

	m.setStatus(message)
}

// moveUp moves the cursor up by one line.
func (m *Model) moveUp() {

	if m.Cursor > 0 {
		m.Cursor--
	}

	m.updateStatusFromSelection()
}

// moveDown moves the cursor down by one line.
func (m *Model) moveDown() {

	if m.Cursor < len(m.Packs)-1 {
		m.Cursor++
	}

	m.updateStatusFromSelection()
}

// toggleCurrent toggles the currently selected pack.
func (m *Model) toggleCurrent() {

	if len(m.Packs) == 0 {
		return
	}

	pack := &m.Packs[m.Cursor]

	if pack.Status == internal.StatusDuplicate {

		m.setStatus(
			"Duplicate pack. Toggle the World pack instead.",
		)

		return
	}

	internal.TogglePack(pack)
	m.updateStatusFromSelection()
}

func (m *Model) moveCurrentUp() {

	if m.Cursor <= 0 {
		return
	}

	m.Packs[m.Cursor], m.Packs[m.Cursor-1] =
		m.Packs[m.Cursor-1], m.Packs[m.Cursor]

	m.Cursor--

	m.updateStatusFromSelection()
}

func (m *Model) moveCurrentDown() {

	if m.Cursor >= len(m.Packs)-1 {
		return
	}

	m.Packs[m.Cursor], m.Packs[m.Cursor+1] =
		m.Packs[m.Cursor+1], m.Packs[m.Cursor]

	m.Cursor++

	m.updateStatusFromSelection()
}
