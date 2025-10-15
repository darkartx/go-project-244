package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/darkartx/go-project-244/shared"
)

type plain struct {
	rootDiff shared.Diff
	builder  strings.Builder
	scope    []string
}

func newPlain(diff shared.Diff) *plain {
	result := &plain{rootDiff: diff}
	return result
}

func (f *plain) Build() string {
	f.addDiff(f.rootDiff)
	return f.builder.String()
}

func (f *plain) addDiff(diff shared.Diff) {
	keys := slices.Sorted(maps.Keys(diff.Child))

	for _, key := range keys {
		diffItem := diff.Child[key]
		f.scope = append(f.scope, key)

		switch diffItem.Change {
		case "added":
			f.addAdded(diffItem)
		case "removed":
			f.addRemoved()
		case "value_changed":
			f.addValueChanged(diffItem)
		case "diff":
			f.addDiff(diffItem)
		}

		f.scope = f.scope[:len(f.scope)-1]
	}
}

func (f *plain) addAdded(diff shared.Diff) {
	key := strings.Join(f.scope, ".")
	value := plainValue(diff.RightValue)

	fmt.Fprintf(&f.builder, "Property '%s' was added with value: %s\n", key, value)
}

func (f *plain) addRemoved() {
	key := strings.Join(f.scope, ".")

	fmt.Fprintf(&f.builder, "Property '%s' was removed\n", key)
}

func (f *plain) addValueChanged(diff shared.Diff) {
	key := strings.Join(f.scope, ".")
	valueLeft := plainValue(diff.LeftValue)
	valueRight := plainValue(diff.RightValue)

	fmt.Fprintf(&f.builder, "Property '%s' was updated. From %s to %s\n", key, valueLeft, valueRight)
}

func plainValue(value any) string {
	var result string

	switch value := value.(type) {
	case nil:
		result = "null"
	case string:
		result = fmt.Sprintf("'%s'", value)
	case map[string]any:
		result = "[complex value]"
	default:
		result = fmt.Sprintf("%v", value)
	}

	return result
}
