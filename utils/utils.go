package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/olekukonko/tablewriter"
)

func HumanDuration(seconds int64) string {
	createdAt := time.Unix(seconds, 0)

	if createdAt.IsZero() {
		return ""
	}
	// https://github.com/docker/cli/blob/0e70f1b7b831565336006298b9443b015c3c87a5/cli/command/formatter/buildcache.go#L156
	return units.HumanDuration(time.Now().UTC().Sub(createdAt)) + " ago"
}

func HumanSize(size int64) string {
	// https://github.com/docker/cli/blob/0e70f1b7b831565336006298b9443b015c3c87a5/cli/command/formatter/buildcache.go#L148
	return units.HumanSizeWithPrecision(float64(size), 3)
}

func WriteToTable(header []string, rows [][]string) {
	if len(rows) == 0 {
		fmt.Println("You've got nothing")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("   ")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(rows)
	table.Render()
}

func FormatFilter(filter string) (string, string, error) {
	if len(filter) == 0 {
		return "", "", errors.New("empty filter is not allowed")
	}

	if !strings.Contains(filter, "=") {
		return "", "", errors.New(fmt.Sprintf("invalid filter format: `%s`", filter))
	}

	// in above check it is confirm that filter has `=`
	splitter := strings.Split(filter, "=")
	if len(splitter[0]) == 0 {
		return "", "", errors.New(fmt.Sprintf("key must not be empty : `%s`", filter))
	}

	if len(splitter[1]) == 0 {
		return "", "", errors.New(fmt.Sprintf("value must not be empty : `%s`", filter))
	}

	if !isValidKey(splitter[0]) {
		return "", "", errors.New(fmt.Sprintf("invalid key : `%s`\nAllowed keys: dangling, label, before, since, reference", splitter[0]))
	}

	return splitter[0], splitter[1], nil
}

func isValidKey(key string) bool {
	return key == "dangling" || key == "label" || key == "before" || key == "since" || key == "reference"
}
