package tui

// StatusLevel represents the semantic type of a status message.
type StatusLevel int

const (
	StatusLevelNone StatusLevel = iota
	StatusLevelInfo
	StatusLevelSuccess
	StatusLevelWarning
	StatusLevelError
	StatusLevelConfirm
	StatusLevelSelection
)

// StatusMessage is the common message format used by TUI features.
type StatusMessage struct {
	Level StatusLevel
	Text  string
}

// Empty returns true when the status has no text to display.
func (s StatusMessage) Empty() bool {
	return s.Text == ""
}

// Label returns a short display label for the status level.
func (s StatusMessage) Label() string {
	switch s.Level {
	case StatusLevelInfo:
		return "INFO"
	case StatusLevelSuccess:
		return "OK"
	case StatusLevelWarning:
		return "WARN"
	case StatusLevelError:
		return "ERROR"
	case StatusLevelConfirm:
		return "CONFIRM"
	case StatusLevelSelection:
		return "SELECT"
	default:
		return ""
	}
}

// setStatus updates the status bar message.
func (m *Model) setStatus(level StatusLevel, message string) {
	m.StatusMessage = StatusMessage{
		Level: level,
		Text:  message,
	}
}

// clearStatus clears the active status message.
func (m *Model) clearStatus() {
	m.StatusMessage = StatusMessage{}
}

// setSelectionStatus updates the status message derived from the cursor.
func (m *Model) setSelectionStatus(level StatusLevel, message string) {
	m.SelectionMessage = StatusMessage{
		Level: level,
		Text:  message,
	}
}

// currentStatusMessage returns the active feature message, or the selection
// message when no feature message is active.
func (m Model) currentStatusMessage() StatusMessage {
	if !m.StatusMessage.Empty() {
		return m.StatusMessage
	}

	return m.SelectionMessage
}
