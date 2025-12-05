package formatters

import (
	"encoding/json"
	"maps"
	"slices"

	"code/internal"
)

var DiffChangeMap = map[int]string{
	internal.DIFF_CHANGE_UNCAHGNED:     "unchanged",
	internal.DIFF_CHANGE_VALUE_ADDED:   "added",
	internal.DIFF_CHANGE_VALUE_CHANGED: "value_changed",
	internal.DIFF_CHANGE_VALUE_REMOVED: "removed",
	internal.DIFF_CHANGE_DIFF:          "diff",
}

type Json struct {
	rootDiff internal.Diff
	result   map[string]any
	current  *map[string]any
}

func NewJson(diff internal.Diff) *Json {
	result := make(map[string]any)
	current := &result
	return &Json{rootDiff: diff, result: result, current: current}
}

func (f *Json) Build() (string, error) {
	f.addDiff(&f.rootDiff)

	byteJson, err := json.MarshalIndent(f.result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(byteJson), nil
}

func (f *Json) addDiff(diff *internal.Diff) {
	keys := slices.Sorted(maps.Keys(diff.Child))

	for _, key := range keys {
		diffItem := diff.Child[key]
		f.addItem(&diffItem)
	}
}

func (f *Json) addItem(diff *internal.Diff) {
	value := make(map[string]any)

	value["change"] = DiffChangeMap[diff.Change]

	if diff.Change == internal.DIFF_CHANGE_VALUE_REMOVED ||
		diff.Change == internal.DIFF_CHANGE_VALUE_CHANGED ||
		diff.Change == internal.DIFF_CHANGE_UNCAHGNED {
		value["left_value"] = diff.LeftValue
	}

	if diff.Change == internal.DIFF_CHANGE_VALUE_ADDED ||
		diff.Change == internal.DIFF_CHANGE_VALUE_CHANGED {
		value["right_value"] = diff.RightValue
	}

	if diff.Change == internal.DIFF_CHANGE_DIFF {
		current := make(map[string]any)
		prev := f.current
		f.current = &current

		f.addDiff(diff)
		value["diff"] = current

		f.current = prev
	}

	(*f.current)[diff.Key] = value
}
