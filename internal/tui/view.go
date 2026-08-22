package tui

import (
	"fmt"
	"strings"

	"github.com/o-unit/bdspackman/internal"
)

// View renders the application.
func (m Model) View() string {

	var b strings.Builder

	// ----------------------------------------------------------------
	// Add Pack Input
	// ----------------------------------------------------------------

	if m.Mode == ModeAddInput {

		target := "Server"
		if m.AddToWorld {
			target = "World"
		}

		fmt.Fprint(&b, "bdspackman")
		fmt.Fprintln(&b)
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "Add Pack (%s)\n\n", target)

		fmt.Fprintf(&b, "Path : %s\n", m.AddInput.View())

		status := m.currentStatusMessage()
		if !status.Empty() {
			fmt.Fprintln(&b)
			fmt.Fprintln(&b, strings.Repeat("-", 40))
			fmt.Fprintln(&b, RenderStatusMessage(status))
		}

		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "Enter: Add    Esc: Cancel")

		return b.String()
	}

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

		header := c.Header

		// 幅を揃える
		if c.Width > 0 {
			header = fmt.Sprintf("%-*s", c.Width, header)
		}

		// 必要なら色付け
		header = RenderHeader(header)

		// 出力
		if c.Width == 0 {
			b.WriteString(header)
		} else {
			b.WriteString(header)
			b.WriteByte(' ')
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
			cursor = selectedStyle.Render("▶")
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

			// まず幅を揃える
			if c.Width > 0 {
				value = fmt.Sprintf("%-*s", c.Width, value)
			}

			// Status列だけ色を付ける
			if c.Header == "Status" {
				value = RenderStatus(value, pack.Status)
			}

			// 出力
			if c.Width == 0 {
				b.WriteString(value)
			} else {
				b.WriteString(value)
				b.WriteByte(' ')
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
		dup    int
		system int
		errors int
	)

	for _, pack := range m.Packs {

		switch pack.Status {

		case internal.StatusOn:
			on++

		case internal.StatusOff:
			off++

		case internal.StatusDuplicate:
			dup++

		case internal.StatusSystem:
			system++

		case internal.StatusError:
			errors++
		}
	}

	fmt.Fprintln(&b)
	fmt.Fprintf(
		&b,
		"%s : %d   %s : %d   %s : %d   %s : %d   %s : %d\n",
		statusOnStyle.Render("ON"),
		on,
		statusOffStyle.Render("OFF"),
		off,
		statusDuplicateStyle.Render("DUP"),
		dup,
		statusSystemStyle.Render("System"),
		system,
		statusErrorStyle.Render("Errors"),
		errors,
	)

	// ----------------------------------------------------------------
	// Status
	// ----------------------------------------------------------------

	status := m.currentStatusMessage()
	if !status.Empty() {

		fmt.Fprintln(&b)
		fmt.Fprintln(&b, strings.Repeat("-", 40))
		fmt.Fprintln(&b, RenderStatusMessage(status))
	}

	// ----------------------------------------------------------------
	// Help
	// ----------------------------------------------------------------

	fmt.Fprintln(&b)
	fmt.Fprintln(&b, m.HelpMessage)

	return b.String()
}
