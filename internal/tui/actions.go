package tui

import "github.com/o-unit/bdspackman/internal"

func (m *Model) setHelp(message string) {
	m.HelpMessage = message
}

// updateHelp updates the help message according to the current mode.
func (m *Model) updateHelp() {

	switch m.Mode {

	case ModeNormal:
		m.HelpMessage =
			"↑↓ Move  Ctrl+↑↓ Reorder Space Toggle  M MovePack  D Delete  R Rename  Esc Quit"

	case ModeConfirmSave:
		m.HelpMessage =
			"Y Save  N Cancel"
	case ModeConfirmDelete:
		m.HelpMessage =
			"Y Delete  N Cancel"
	case ModeRenameDir:
		m.HelpMessage =
			"Enter Rename  Backspace DeleteChar  Esc Cancel"
	default:
		m.HelpMessage = ""
	}
}

// updateStatusFromSelection updates the message for the selected pack.
func (m *Model) updateStatusFromSelection() {

	if len(m.Packs) == 0 {
		m.SelectionMessage = StatusMessage{}
		return
	}

	pack := m.Packs[m.Cursor]

	switch pack.Status {

	case internal.StatusError:

		if pack.Error != nil {
			m.setSelectionStatus(StatusLevelError, pack.Error.Error())
		} else {
			m.setSelectionStatus(StatusLevelError, "Unknown error.")
		}

	case internal.StatusSystem:
		m.setSelectionStatus(StatusLevelSelection, "System pack")

	case internal.StatusOn:
		m.setSelectionStatus(StatusLevelSelection, "Enabled")

	case internal.StatusOff:
		m.setSelectionStatus(StatusLevelSelection, "Disabled")

	default:
		m.SelectionMessage = StatusMessage{}
	}
}

// enterSaveConfirm switches to save confirmation mode.
func (m *Model) enterSaveConfirm() {

	m.Mode = ModeConfirmSave

	m.setStatus(StatusLevelConfirm, "Save changes? (Y=yes / N=no)")
	m.updateHelp()
}

func (m *Model) enterDeleteConfirm() {
	if len(m.Packs) == 0 {
		return
	}

	pack := m.Packs[m.Cursor]
	if pack.Status == internal.StatusSystem {
		m.setStatus(StatusLevelWarning, "System pack cannot be deleted.")
		return
	}

	m.Mode = ModeConfirmDelete
	m.setStatus(StatusLevelConfirm, "Delete "+pack.FolderName+"? (Y=yes / N=no)")
	m.updateHelp()
}

func (m *Model) enterRenameDirMode() {
	if len(m.Packs) == 0 {
		return
	}

	pack := m.Packs[m.Cursor]
	if pack.Status == internal.StatusSystem {
		m.setStatus(StatusLevelWarning, "System pack cannot be renamed.")
		return
	}

	m.Mode = ModeRenameDir
	m.RenameInput = pack.FolderName
	m.updateRenameStatus()
	m.updateHelp()
}

// enterNormalMode switches back to normal mode.
func (m *Model) enterNormalMode() {

	m.Mode = ModeNormal
	m.RenameInput = ""
	m.clearStatus()
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
		m.setStatus(StatusLevelError, "Save failed: "+err.Error())
		return
	}

	m.setStatus(StatusLevelSuccess, message)
}

func (m *Model) deleteCurrentPack() {
	if len(m.Packs) == 0 {
		m.enterNormalMode()
		return
	}

	pack := m.Packs[m.Cursor]

	if err := internal.DeletePack(pack); err != nil {
		m.enterNormalMode()
		m.setStatus(StatusLevelError, "Delete failed: "+err.Error())
		return
	}

	m.Packs = append(m.Packs[:m.Cursor], m.Packs[m.Cursor+1:]...)
	if m.Cursor >= len(m.Packs) && m.Cursor > 0 {
		m.Cursor--
	}

	m.enterNormalMode()
	m.setStatus(StatusLevelSuccess, "Deleted "+pack.FolderName+".")
	m.updateStatusFromSelection()
}

func (m *Model) updateRenameStatus() {
	m.setStatus(StatusLevelConfirm, "Directory name: "+m.RenameInput)
}

func (m *Model) renameInputAppend(text string) {
	m.RenameInput += text
	m.updateRenameStatus()
}

func (m *Model) renameInputBackspace() {
	if m.RenameInput == "" {
		m.updateRenameStatus()
		return
	}

	runes := []rune(m.RenameInput)
	m.RenameInput = string(runes[:len(runes)-1])
	m.updateRenameStatus()
}

func (m *Model) renameCurrentPackDirectory() {
	if len(m.Packs) == 0 {
		m.enterNormalMode()
		return
	}

	pack := &m.Packs[m.Cursor]
	oldName := pack.FolderName

	if err := internal.RenamePackDirectory(pack, m.RenameInput); err != nil {
		m.setStatus(StatusLevelError, "Rename failed: "+err.Error())
		return
	}

	m.enterNormalMode()

	if oldName == pack.FolderName {
		m.setStatus(StatusLevelWarning, "Directory name was not changed.")
	} else {
		m.setStatus(StatusLevelSuccess, "Renamed "+oldName+" to "+pack.FolderName+".")
	}

	m.updateStatusFromSelection()
}

// moveUp moves the cursor up by one line.
func (m *Model) moveUp() {

	if m.Cursor > 0 {
		m.Cursor--
		m.clearStatus()
	}

	m.updateStatusFromSelection()
}

// moveDown moves the cursor down by one line.
func (m *Model) moveDown() {

	if m.Cursor < len(m.Packs)-1 {
		m.Cursor++
		m.clearStatus()
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
			StatusLevelWarning,
			"Duplicate pack. Toggle the World pack instead.",
		)

		return
	}

	internal.TogglePack(pack)
	m.clearStatus()
	m.updateStatusFromSelection()
}

func (m *Model) moveCurrentPackLocation() {
	if len(m.Packs) == 0 {
		return
	}

	pack := &m.Packs[m.Cursor]

	if pack.Status == internal.StatusSystem {
		m.setStatus(StatusLevelWarning, "System pack cannot be moved.")
		return
	}

	if pack.Status == internal.StatusError {
		m.setStatus(StatusLevelWarning, "Pack with errors cannot be moved.")
		return
	}

	if pack.Status == internal.StatusDuplicate {
		m.setStatus(StatusLevelWarning, "Duplicate pack cannot be moved.")
		return
	}

	targetLocation, _ := internal.PackMoveTarget(m.Config, *pack)
	if m.hasPackConflict(*pack, targetLocation) {
		m.setStatus(StatusLevelWarning, "A pack with the same UUID already exists in the destination.")
		return
	}

	if err := internal.MovePack(m.Config, pack); err != nil {
		m.setStatus(StatusLevelError, "Move failed: "+err.Error())
		return
	}

	m.setStatus(
		StatusLevelSuccess,
		"Moved pack to "+pack.Location.String()+".",
	)
	m.updateStatusFromSelection()
}

func (m *Model) moveCurrentUp() {

	if m.Cursor <= 0 {
		return
	}

	m.Packs[m.Cursor], m.Packs[m.Cursor-1] =
		m.Packs[m.Cursor-1], m.Packs[m.Cursor]

	m.Cursor--

	m.clearStatus()
	m.updateStatusFromSelection()
}

func (m Model) hasPackConflict(pack internal.Pack, targetLocation internal.PackLocation) bool {
	if pack.UUID == "" {
		return false
	}

	for i := range m.Packs {
		other := m.Packs[i]
		if i == m.Cursor {
			continue
		}

		if other.Type != pack.Type {
			continue
		}

		if other.Location != targetLocation {
			continue
		}

		if other.UUID == pack.UUID {
			return true
		}
	}

	return false
}

func (m *Model) moveCurrentDown() {

	if m.Cursor >= len(m.Packs)-1 {
		return
	}

	m.Packs[m.Cursor], m.Packs[m.Cursor+1] =
		m.Packs[m.Cursor+1], m.Packs[m.Cursor]

	m.Cursor++

	m.clearStatus()
	m.updateStatusFromSelection()
}
