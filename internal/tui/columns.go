package tui

import (
	"fmt"

	"github.com/o-unit/bdspackman/internal"
)

type valueFunc func(internal.Pack, int) string

type column struct {
	Header string
	Width  int
	Value  valueFunc
}

// columns returns the columns displayed in the pack list.
func (m Model) columns() []column {

	columns := []column{
		{
			Header: "No",
			Width:  4,
			Value: func(_ internal.Pack, index int) string {
				return fmt.Sprintf("%d", index+1)
			},
		},
		{
			Header: "Type",
			Width:  4,
			Value: func(pack internal.Pack, _ int) string {
				return pack.Type.String()
			},
		},
		{
			Header: "Location",
			Width:  8,
			Value: func(pack internal.Pack, _ int) string {
				return pack.Location.String()
			},
		},
		{
			Header: "Status",
			Width:  8,
			Value: func(pack internal.Pack, _ int) string {
				return pack.Status.String()
			},
		},
		{
			Header: "Version",
			Width:  10,
			Value: func(pack internal.Pack, _ int) string {
				return pack.VersionString()
			},
		},
	}

	if m.Config.ShowUUID {
		columns = append(columns, column{
			Header: "UUID",
			Width:  36,
			Value: func(pack internal.Pack, _ int) string {
				return pack.UUID
			},
		})
	}

	if m.Config.ShowDirName {
		columns = append(columns, column{
			Header: "DirName",
			Width:  30,
			Value: func(pack internal.Pack, _ int) string {
				return pack.FolderName
			},
		})
	}

	columns = append(columns, column{
		Header: "Name",
		Width:  0, // 0 = remaining width
		Value: func(pack internal.Pack, _ int) string {
			return pack.DisplayName
		},
	})

	return columns
}
