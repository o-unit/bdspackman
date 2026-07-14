package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/o-unit/bdspackman/internal"
	"github.com/o-unit/bdspackman/internal/tui"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {

	cfg, err := internal.LoadFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if cfg.Version {
		fmt.Printf(
			"bdspackman %s (commit: %s)\n",
			version,
			commit,
		)
		return
	}

	cfg, err = internal.CompleteConfig(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	packs, err := internal.ScanPacks(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	model := tui.NewModel(cfg, packs)

	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
