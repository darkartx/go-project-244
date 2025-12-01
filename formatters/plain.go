package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/darkartx/go-project-244/internal"
)

type Plain struct {
	rootDiff internal.Diff
	builder  strings.Builder
	scope    []string
}

func NewPlain(diff internal.Diff) *Plain {
	result := &Plain{rootDiff: diff}
	return result
}

func (f *Plain) Build() (string, error) {
	f.addDiff(&f.rootDiff)
	return f.builder.String()[1:], nil
}

func (f *Plain) addDiff(diff *internal.Diff) {
	keys := slices.Sorted(maps.Keys(diff.Child))

	for _, key := range keys {
		diffItem := diff.Child[key]
		f.scope = append(f.scope, key)

		switch diffItem.Change {
		case internal.DIFF_CHANGE_VALUE_ADDED:
			fmt.Fprint(&f.builder, "\n")
			f.addAdded(&diffItem)
		case internal.DIFF_CHANGE_VALUE_REMOVED:
			fmt.Fprint(&f.builder, "\n")
			f.addRemoved()
		case internal.DIFF_CHANGE_VALUE_CHANGED:
			fmt.Fprint(&f.builder, "\n")
			f.addValueChanged(&diffItem)
		case internal.DIFF_CHANGE_DIFF:
			f.addDiff(&diffItem)
		}

		f.scope = f.scope[:len(f.scope)-1]
	}
}

func (f *Plain) addAdded(diff *internal.Diff) {
	key := strings.Join(f.scope, ".")
	value := plainValue(diff.RightValue)

	fmt.Fprintf(&f.builder, "Property '%s' was added with value: %s", key, value)
}

func (f *Plain) addRemoved() {
	key := strings.Join(f.scope, ".")

	fmt.Fprintf(&f.builder, "Property '%s' was removed", key)
}

func (f *Plain) addValueChanged(diff *internal.Diff) {
	key := strings.Join(f.scope, ".")
	valueLeft := plainValue(diff.LeftValue)
	valueRight := plainValue(diff.RightValue)

	fmt.Fprintf(&f.builder, "Property '%s' was updated. From %s to %s", key, valueLeft, valueRight)
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
