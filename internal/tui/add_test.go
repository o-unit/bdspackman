package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-unit/bdspackman/internal"
)

func TestAddInputCanBeCancelled(t *testing.T) {
	model := NewModel(internal.Config{}, nil)

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("A")})
	addModel := updated.(Model)
	if addModel.Mode != ModeAddInput || addModel.AddToWorld {
		t.Fatalf("add mode = %v, toWorld = %v", addModel.Mode, addModel.AddToWorld)
	}

	updated, _ = addModel.Update(tea.KeyMsg{Type: tea.KeyEsc})
	normalModel := updated.(Model)
	if normalModel.Mode != ModeNormal {
		t.Fatalf("mode = %v, want normal", normalModel.Mode)
	}
	if normalModel.StatusMessage.Text != "Add cancelled." {
		t.Fatalf("status = %q", normalModel.StatusMessage.Text)
	}
}
