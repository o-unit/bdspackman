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
			case "ctrl+up":
				m.moveCurrentUp()
			case "ctrl+down":
				m.moveCurrentDown()
			}

		case ModeConfirmSave:
			switch msg.String() {
			case "y":
				m.save()
			case "n", "esc":
				m.enterNormalMode()
			case "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}
