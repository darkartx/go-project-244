package code

import (
	"fmt"
	"slices"
	"strings"
)

type entry struct {
	key        string
	valueLeft  any
	valueRight any
}

func GenDiff(filepathLeft, filepathRight, format string) (string, error) {
	dataLeft, err := parseFile(filepathLeft)
	if err != nil {
		return "", err
	}

	dataRight, err := parseFile(filepathRight)
	if err != nil {
		return "", err
	}

	entries := getEntries(dataLeft, dataRight)

	return genDiff(entries), nil
}

func genDiff(entries []entry) string {
	var sb strings.Builder

	sb.WriteString("{\n")

	for _, entry := range entries {
		if entry.valueLeft == nil {
			fmt.Fprintf(&sb, "  + %s: %v\n", entry.key, entry.valueRight)
			continue
		}

		if entry.valueRight == nil {
			fmt.Fprintf(&sb, "  - %s: %v\n", entry.key, entry.valueLeft)
			continue
		}

		if entry.valueLeft != entry.valueRight {
			fmt.Fprintf(&sb, "  - %s: %v\n", entry.key, entry.valueLeft)
			fmt.Fprintf(&sb, "  + %s: %v\n", entry.key, entry.valueRight)
			continue
		}

		fmt.Fprintf(&sb, "    %s: %v\n", entry.key, entry.valueLeft)
	}

	sb.WriteString("}")

	return sb.String()
}

func getEntries(dataLeft, dataRight map[string]any) []entry {
	var result []entry
	var keys []string

	for k := range dataLeft {
		keys = append(keys, k)
	}

	for k := range dataRight {
		if !slices.Contains(keys, k) {
			keys = append(keys, k)
		}
	}

	slices.Sort(keys)

	for _, key := range keys {
		valueLeft := dataLeft[key]
		valueRight := dataRight[key]
		result = append(result, entry{key, valueLeft, valueRight})
	}

	return result
}
