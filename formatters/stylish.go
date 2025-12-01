package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/darkartx/go-project-244/internal"
)

type Stylish struct {
	rootDiff internal.Diff
	builder  strings.Builder
	indent   uint8
}

func NewStylish(diff internal.Diff) *Stylish {
	result := &Stylish{rootDiff: diff}
	return result
}

func (f *Stylish) Build() (string, error) {
	f.addDiff(&f.rootDiff)
	return f.builder.String(), nil
}

func (f *Stylish) addDiff(diff *internal.Diff) {
	if diff.Key == "" {
		f.addIndeted(' ', "{\n")
	} else {
		f.addIndeted(' ', "%s: {\n", diff.Key)
	}

	keys := slices.Sorted(maps.Keys(diff.Child))

	f.indent += 1
	for _, key := range keys {
		diffItem := diff.Child[key]

		switch diffItem.Change {
		case internal.DIFF_CHANGE_UNCAHGNED:
			f.addUnchanged(&diffItem)
		case internal.DIFF_CHANGE_VALUE_ADDED:
			f.addAdded(&diffItem)
		case internal.DIFF_CHANGE_VALUE_REMOVED:
			f.addRemoved(&diffItem)
		case internal.DIFF_CHANGE_VALUE_CHANGED:
			f.addValueChanged(&diffItem)
		case internal.DIFF_CHANGE_DIFF:
			f.addDiff(&diffItem)
		}

		fmt.Fprint(&f.builder, "\n")
	}
	f.indent -= 1

	f.addIndeted(' ', "}")
}

func (f *Stylish) addUnchanged(diff *internal.Diff) {
	f.addIndeted(' ', "%s: ", diff.Key)
	f.addValue(diff.LeftValue)
}

func (f *Stylish) addAdded(diff *internal.Diff) {
	f.addIndeted('+', "%s: ", diff.Key)
	f.addValue(diff.RightValue)
}

func (f *Stylish) addRemoved(diff *internal.Diff) {
	f.addIndeted('-', "%s: ", diff.Key)
	f.addValue(diff.LeftValue)
}

func (f *Stylish) addValueChanged(diff *internal.Diff) {
	f.addRemoved(diff)
	fmt.Fprint(&f.builder, "\n")
	f.addAdded(diff)
}

func (f *Stylish) addValue(value any) {
	switch value := value.(type) {
	case map[string]any:
		f.addMap(value)
	case nil:
		fmt.Fprint(&f.builder, "null")
	default:
		fmt.Fprintf(&f.builder, "%v", value)
	}
}

func (f *Stylish) addIndeted(sym rune, format string, args ...any) {
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

func (f *Stylish) addMap(value map[string]any) {
	fmt.Fprintf(&f.builder, "{\n")

	keys := slices.Sorted(maps.Keys(value))

	f.indent += 1
	for _, key := range keys {
		item := value[key]

		f.addIndeted(' ', "%s: ", key)
		f.addValue(item)
		fmt.Fprint(&f.builder, "\n")
	}
	f.indent -= 1

	f.addIndeted(' ', "}")
}
