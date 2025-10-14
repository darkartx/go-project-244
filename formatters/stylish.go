package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/darkartx/go-project-244/shared"
)

type stylish struct {
	rootDiff shared.Diff
	builder  strings.Builder
	indent   uint8
}

func newStylish(diff shared.Diff) *stylish {
	result := &stylish{rootDiff: diff}
	return result
}

func (f *stylish) Build() string {
	f.addDiff(f.rootDiff.Key, f.rootDiff)
	return f.builder.String()
}

func (f *stylish) addDiff(key string, diff shared.Diff) {
	if key == "" {
		f.addIndeted(' ', "{\n")
	} else {
		f.addIndeted(' ', "%s: {\n", key)
	}

	keys := slices.Sorted(maps.Keys(diff.Child))

	f.indent += 1
	for _, key := range keys {
		diffItem := diff.Child[key]

		switch diffItem.Change {
		case "unchanged":
			f.addUnchanged(key, diffItem)
		case "added":
			f.addAdded(key, diffItem)
		case "removed":
			f.addRemoved(key, diffItem)
		case "value_changed":
			f.addValueChanged(key, diffItem)
		case "diff":
			f.addDiff(key, diffItem)
		}
	}
	f.indent -= 1

	f.addIndeted(' ', "}\n")
}

func (f *stylish) addUnchanged(key string, diff shared.Diff) {
	f.addIndeted(' ', "%s: ", key)
	f.addValue(diff.LeftValue)
}

func (f *stylish) addAdded(key string, diff shared.Diff) {
	f.addIndeted('+', "%s: ", key)
	f.addValue(diff.RightValue)
}

func (f *stylish) addRemoved(key string, diff shared.Diff) {
	f.addIndeted('-', "%s: ", key)
	f.addValue(diff.LeftValue)
}

func (f *stylish) addValueChanged(key string, diff shared.Diff) {
	f.addRemoved(key, diff)
	f.addAdded(key, diff)
}

func (f *stylish) addValue(value any) {
	switch value := value.(type) {
	case map[string]any:
		f.addMap(value)
	case nil:
		fmt.Fprint(&f.builder, "null\n")
	default:
		fmt.Fprintf(&f.builder, "%v\n", value)
	}
}

func (f *stylish) addIndeted(sym rune, format string, args ...any) {
	if f.indent > 0 {
		fmt.Fprintf(
			&f.builder,
			"%s  %c ",
			strings.Repeat(" ", (int(f.indent)-1)*4),
			sym,
		)
	}
	fmt.Fprintf(&f.builder, format, args...)
}

func (f *stylish) addMap(value map[string]any) {
	fmt.Fprintf(&f.builder, "{\n")

	keys := slices.Sorted(maps.Keys(value))

	f.indent += 1
	for _, key := range keys {
		item := value[key]

		f.addIndeted(' ', "%s: ", key)
		f.addValue(item)
	}
	f.indent -= 1

	f.addIndeted(' ', "}\n")
}
