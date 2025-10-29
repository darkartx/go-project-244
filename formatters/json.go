package formatters

import (
	enJson "encoding/json"
	"maps"
	"slices"

	"github.com/darkartx/go-project-244/shared"
)

type json struct {
	rootDiff shared.Diff
	result   map[string]any
	current  *map[string]any
}

func newJson(diff shared.Diff) *json {
	result := make(map[string]any)
	current := &result
	return &json{rootDiff: diff, result: result, current: current}
}

func (f *json) Build() (string, error) {
	f.addDiff(&f.rootDiff)

	byteJson, err := enJson.MarshalIndent(f.result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(byteJson), nil
}

func (f *json) addDiff(diff *shared.Diff) {
	keys := slices.Sorted(maps.Keys(diff.Child))

	for _, key := range keys {
		diffItem := diff.Child[key]
		f.addItem(&diffItem)
	}
}

func (f *json) addItem(diff *shared.Diff) {
	value := make(map[string]any)

	value["change"] = diff.Change

	if diff.Change == "removed" || diff.Change == "value_changed" || diff.Change == "unchanged" {
		value["left_value"] = diff.LeftValue
	}

	if diff.Change == "added" || diff.Change == "value_changed" {
		value["right_value"] = diff.RightValue
	}

	if diff.Change == "diff" {
		current := make(map[string]any)
		prev := f.current
		f.current = &current

		f.addDiff(diff)
		value["diff"] = current

		f.current = prev
	}

	(*f.current)[diff.Key] = value
}
