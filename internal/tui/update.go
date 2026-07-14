package tui

import tea "github.com/charmbracelet/bubbletea"

// Update handles keyboard input.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.Mode {
		case ModeNormal:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "up":
				m.moveUp()
			case "down":
				m.moveDown()
			case " ":
				m.toggleCurrent()
			case "enter":
				m.enterSaveConfirm()
			case "m", "M":
				m.moveCurrentPackLocation()
			case "d", "D":
				m.enterDeleteConfirm()
			case "r", "R":
				m.enterRenameDirMode()
			case "ctrl+up":
				m.moveCurrentUp()
			case "ctrl+down":
				m.moveCurrentDown()
			}

		case ModeConfirmSave:
			switch msg.String() {
			case "y", "Y":
				m.save()
			case "n", "N", "esc":
				m.enterNormalMode()
			case "ctrl+c":
				return m, tea.Quit
			}

		case ModeConfirmDelete:
			switch msg.String() {
			case "y", "Y":
				m.deleteCurrentPack()
			case "n", "N", "esc":
				m.enterNormalMode()
			case "ctrl+c":
				return m, tea.Quit
			}

		case ModeRenameDir:
			switch msg.Type {
			case tea.KeyEnter:
				m.renameCurrentPackDirectory()
			case tea.KeyBackspace, tea.KeyDelete:
				m.renameInputBackspace()
			case tea.KeyEsc:
				m.enterNormalMode()
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyRunes:
				m.renameInputAppend(string(msg.Runes))
			}
		}
	}

	return m, nil
}
