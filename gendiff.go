package code

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type diff struct {
	key        string
	change     string
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

	diff := genDiff(dataLeft, dataRight)

	return buildDiff(diff), nil
}

func buildDiff(diff map[string]diff) string {
	var sb strings.Builder

	keys := slices.Sorted(maps.Keys(diff))

	sb.WriteString("{\n")

	for _, key := range keys {
		diffItem := diff[key]

		if diffItem.change == "added" {
			fmt.Fprintf(&sb, "  + %s: %v\n", key, diffItem.valueRight)
			continue
		}

		if diffItem.change == "removed" {
			fmt.Fprintf(&sb, "  - %s: %v\n", key, diffItem.valueLeft)
			continue
		}

		if diffItem.change == "value_changed" {
			fmt.Fprintf(&sb, "  - %s: %v\n", key, diffItem.valueLeft)
			fmt.Fprintf(&sb, "  + %s: %v\n", key, diffItem.valueRight)
			continue
		}

		fmt.Fprintf(&sb, "    %s: %v\n", key, diffItem.valueLeft)
	}

	sb.WriteString("}")

	return sb.String()
}

func genDiff(dataLeft, dataRight map[string]any) map[string]diff {
	result := make(map[string]diff)
	var keys []string

	for k := range dataLeft {
		keys = append(keys, k)
	}

	for k := range dataRight {
		if !slices.Contains(keys, k) {
			keys = append(keys, k)
		}
	}

	for _, key := range keys {
		valueLeft, existsLeft := dataLeft[key]

		valueRight, existsRight := dataRight[key]

		change := "unchanged"

		if !existsLeft {
			change = "added"
		} else if !existsRight {
			change = "removed"
		} else if valueLeft != valueRight {
			change = "value_changed"
		}

		result[key] = diff{key, change, valueLeft, valueRight}
	}

	return result
}
