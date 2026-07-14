package tui

import "testing"

func TestCurrentStatusMessagePrefersActiveStatus(t *testing.T) {
	model := Model{
		StatusMessage: StatusMessage{
			Level: StatusLevelWarning,
			Text:  "active",
		},
		SelectionMessage: StatusMessage{
			Level: StatusLevelSelection,
			Text:  "selection",
		},
	}

	status := model.currentStatusMessage()

	if status.Text != "active" {
		t.Fatalf("status text = %q, want %q", status.Text, "active")
	}

	if status.Level != StatusLevelWarning {
		t.Fatalf("status level = %v, want %v", status.Level, StatusLevelWarning)
	}
}

func TestCurrentStatusMessageFallsBackToSelectionStatus(t *testing.T) {
	model := Model{
		SelectionMessage: StatusMessage{
			Level: StatusLevelSelection,
			Text:  "selection",
		},
	}

	status := model.currentStatusMessage()

	if status.Text != "selection" {
		t.Fatalf("status text = %q, want %q", status.Text, "selection")
	}

	if status.Level != StatusLevelSelection {
		t.Fatalf("status level = %v, want %v", status.Level, StatusLevelSelection)
	}
}
