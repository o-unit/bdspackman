package tui

import (
	"fmt"
	"strings"

	"github.com/o-unit/bdspackman/internal"
)

// View renders the application.
func (m Model) View() string {

	var b strings.Builder

	fmt.Fprint(&b, "bdspackman")

	if m.Config.ExportPrefix != "" {
		fmt.Fprintf(
			&b,
			"Export : %s*\n\n",
			m.Config.ExportPrefix,
		)
	}

	fmt.Fprintln(&b)
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "Server : %s\n", m.Config.ServerDir)
	fmt.Fprintf(&b, "World  : %s\n\n", m.Config.WorldName)

	if len(m.Packs) == 0 {
		fmt.Fprintln(&b, "No packs found.")
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, m.HelpMessage)
		return b.String()
	}

	columns := m.columns()

	// ----------------------------------------------------------------
	// Header
	// ----------------------------------------------------------------

	for _, c := range columns {

		if c.Width == 0 {
			fmt.Fprintf(&b, "%s", c.Header)
		} else {
			fmt.Fprintf(&b, "%-*s ", c.Width, c.Header)
		}
	}

	b.WriteByte('\n')

	// ----------------------------------------------------------------
	// Separator
	// ----------------------------------------------------------------

	for _, c := range columns {

		width := c.Width
		if width == 0 {
			width = 30
		}

		b.WriteString(strings.Repeat("-", width))
		b.WriteByte(' ')
	}

	b.WriteByte('\n')

	// ----------------------------------------------------------------
	// Rows
	// ----------------------------------------------------------------

	for i, pack := range m.Packs {

		cursor := " "
		if i == m.Cursor {
			cursor = "▶"
		}

		for col, c := range columns {

			value := c.Value(pack, i)

			// No列だけカーソル込みで描画
			if col == 0 {

				fmt.Fprintf(
					&b,
					"%s%3s ",
					cursor,
					value,
				)

				continue
			}

			if c.Width == 0 {
				b.WriteString(value)
			} else {
				fmt.Fprintf(
					&b,
					"%-*s ",
					c.Width,
					value,
				)
			}
		}

		b.WriteByte('\n')
	}

	// ----------------------------------------------------------------
	// Summary
	// ----------------------------------------------------------------

	var (
		on     int
		off    int
		system int
		errors int
	)

	for _, pack := range m.Packs {

		switch pack.Status {

		case internal.StatusOn:
			on++

		case internal.StatusOff:
			off++

		case internal.StatusSystem:
			system++

		case internal.StatusError:
			errors++
		}
	}

	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "ON     : %d\n", on)
	fmt.Fprintf(&b, "OFF    : %d\n", off)
	fmt.Fprintf(&b, "System : %d\n", system)
	fmt.Fprintf(&b, "Errors : %d\n", errors)

	// ----------------------------------------------------------------
	// Status
	// ----------------------------------------------------------------

	if m.StatusMessage != "" {

		fmt.Fprintln(&b)
		fmt.Fprintln(&b, strings.Repeat("-", 40))
		fmt.Fprintln(&b, m.StatusMessage)
	}

	// ----------------------------------------------------------------
	// Help
	// ----------------------------------------------------------------

	fmt.Fprintln(&b)
	fmt.Fprintln(&b, m.HelpMessage)

	return b.String()
}
