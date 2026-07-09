package internal

import (
	"fmt"
	"io"
)

// DisplayPacks displays all scanned packs.
func DisplayPacks(w io.Writer, packs []Pack, showUUID bool) {

	if showUUID {
		fmt.Fprintf(
			w,
			"%-4s %-3s %-8s %-8s %-10s %-36s %s\n",
			"No",
			"Type",
			"Location",
			"Status",
			"Version",
			"UUID",
			"Name",
		)

		fmt.Fprintf(
			w,
			"---- --- -------- -------- ---------- ------------------------------------ ------------------------------\n",
		)

	} else {

		fmt.Fprintf(
			w,
			"%-4s %-3s %-8s %-8s %-10s %s\n",
			"No",
			"Type",
			"Location",
			"Status",
			"Version",
			"Name",
		)

		fmt.Fprintf(
			w,
			"---- --- -------- -------- ---------- ------------------------------\n",
		)
	}

	for i, pack := range packs {

		if showUUID {

			fmt.Fprintf(
				w,
				"%-4d %-3s %-8s %-8s %-10s %-36s %s\n",
				i+1,
				pack.Type.String(),
				pack.Location.String(),
				pack.Status.String(),
				pack.VersionString(),
				pack.UUID,
				pack.DisplayName,
			)

		} else {

			fmt.Fprintf(
				w,
				"%-4d %-3s %-8s %-8s %-10s %s\n",
				i+1,
				pack.Type.String(),
				pack.Location.String(),
				pack.Status.String(),
				pack.VersionString(),
				pack.DisplayName,
			)
		}

		if pack.Error != nil {

			fmt.Fprintf(
				w,
				"      ERROR: %v\n",
				pack.Error,
			)

		}
	}
}

// DisplayChanges displays packs that will be changed.
func DisplayChanges(w io.Writer, before, after []Pack) {

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Changes:")
	fmt.Fprintln(w)

	changed := false

	for i := range before {

		if before[i].Status == after[i].Status {
			continue
		}

		changed = true

		fmt.Fprintf(
			w,
			"[%d] %-30s %s -> %s\n",
			i+1,
			before[i].DisplayName,
			before[i].Status.String(),
			after[i].Status.String(),
		)
	}

	if !changed {
		fmt.Fprintln(w, "No changes.")
	}
}

// DisplaySummary displays summary information.
func DisplaySummary(w io.Writer, packs []Pack) {

	var on int
	var off int
	var system int
	var errors int

	for _, pack := range packs {

		switch pack.Status {

		case StatusOn:
			on++

		case StatusOff:
			off++

		case StatusSystem:
			system++

		case StatusError:
			errors++

		}
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "Enabled : %d\n", on)
	fmt.Fprintf(w, "Disabled: %d\n", off)
	fmt.Fprintf(w, "System  : %d\n", system)
	fmt.Fprintf(w, "Errors  : %d\n", errors)
}

// FindPackByNumber converts the displayed number into a slice index.
func FindPackByNumber(
	packs []Pack,
	number int,
) (*Pack, error) {

	if number < 1 {
		return nil, fmt.Errorf("invalid number")
	}

	if number > len(packs) {
		return nil, fmt.Errorf("invalid number")
	}

	return &packs[number-1], nil
}

// TogglePack switches ON/OFF state.
//
// SYSTEM and ERROR packs are ignored.
func TogglePack(pack *Pack) bool {

	switch pack.Status {

	case StatusOn:
		pack.Status = StatusOff
		return true

	case StatusOff:
		pack.Status = StatusOn
		return true

	default:
		return false
	}
}
