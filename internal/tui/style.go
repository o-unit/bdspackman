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

	statusInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12"))

	statusSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("10"))

	statusWarningStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("11"))

	statusConfirmStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("13")).
				Bold(true)

	headerStyle = lipgloss.NewStyle().
			Bold(true)
)

// RenderHeader applies the style for a column header.
func RenderHeader(text string) string {
	return headerStyle.Render(text)
}

func RenderStatusMessage(status StatusMessage) string {
	if status.Empty() {
		return ""
	}

	label := status.Label()
	if label != "" {
		label += ": "
	}

	value := label + status.Text

	switch status.Level {
	case StatusLevelInfo, StatusLevelSelection:
		return statusInfoStyle.Render(value)

	case StatusLevelSuccess:
		return statusSuccessStyle.Render(value)

	case StatusLevelWarning:
		return statusWarningStyle.Render(value)

	case StatusLevelError:
		return statusErrorStyle.Render(value)

	case StatusLevelConfirm:
		return statusConfirmStyle.Render(value)

	default:
		return value
	}
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
