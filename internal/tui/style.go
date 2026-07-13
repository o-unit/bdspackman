package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/o-unit/bdspackman/internal"
)

var (
	normalStyle = lipgloss.NewStyle()

	selectedStyle = lipgloss.NewStyle().
			Bold(true)

	statusOnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	statusOffStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	statusErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("9")).
				Bold(true)

	statusDuplicateStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("11"))

	statusSystemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("12"))

	headerStyle = lipgloss.NewStyle().
			Bold(true)
)

// RenderHeader applies the style for a column header.
func RenderHeader(text string) string {
	return headerStyle.Render(text)
}

func RenderStatus(value string, status internal.PackStatus) string {
	switch status {
	case internal.StatusOn:
		return statusOnStyle.Render(value)

	case internal.StatusOff:
		return statusOffStyle.Render(value)

	case internal.StatusError:
		return statusErrorStyle.Render(value)

	case internal.StatusDuplicate:
		return statusDuplicateStyle.Render(value)

	case internal.StatusSystem:
		return statusSystemStyle.Render(value)

	default:
		return value
	}
}
