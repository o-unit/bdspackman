package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/o-unit/bdspackman/internal"
)

func main() {

	cfg, err := internal.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	packs, err := internal.ScanPacks(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	internal.DisplayPacks(
		os.Stdout,
		packs,
		cfg.ShowUUID,
	)

	internal.DisplaySummary(
		os.Stdout,
		packs,
	)

	before := clonePacks(packs)

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("\nNumber (toggle), s=save, q=quit > ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		switch strings.ToLower(line) {

		case "q":

			fmt.Println("Canceled.")

			return

		case "s":

			goto SAVE
		}

		number, err := strconv.Atoi(line)
		if err != nil {

			fmt.Println("Invalid number.")

			continue
		}

		pack, err := internal.FindPackByNumber(
			packs,
			number,
		)

		if err != nil {

			fmt.Println(err)

			continue
		}

		if !internal.TogglePack(pack) {

			fmt.Println("This pack cannot be changed.")

			continue
		}

		fmt.Println()

		internal.DisplayPacks(
			os.Stdout,
			packs,
			cfg.ShowUUID,
		)

		internal.DisplaySummary(
			os.Stdout,
			packs,
		)
	}

SAVE:

	fmt.Println()

	internal.DisplayChanges(
		os.Stdout,
		before,
		packs,
	)

	if !confirmSave(reader) {

		fmt.Println("Canceled.")

		return
	}

	err = internal.SaveWorldPackFiles(
		cfg,
		packs,
		internal.SaveOptions{
			DryRun: cfg.DryRun,
		},
	)

	if err != nil {

		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	fmt.Println("Saved.")
}

func confirmSave(reader *bufio.Reader) bool {

	for {

		fmt.Print("Save changes? [y/N]: ")

		line, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		line = strings.TrimSpace(strings.ToLower(line))

		switch line {

		case "y", "yes":
			return true

		case "", "n", "no":
			return false

		default:
			fmt.Println("Please enter y or n.")
		}
	}
}

// clonePacks returns a copy of the pack slice.
//
// This is used to compare the state before and after editing.
func clonePacks(src []internal.Pack) []internal.Pack {

	dst := make([]internal.Pack, len(src))
	copy(dst, src)

	return dst
}
