package code

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type UnsupportedFormat struct {
	Format string
}

func (e UnsupportedFormat) Error() string {
	return fmt.Sprintf("unsupported format: %s", e.Format)
}

type formater interface {
	build() string
}

type stylish struct {
	rootDiff diff
	builder  strings.Builder
	indent   uint8
}

func newStylish(diff diff) *stylish {
	result := &stylish{rootDiff: diff}
	return result
}

func (f *stylish) build() string {
	f.addDiff(f.rootDiff.key, f.rootDiff)
	return f.builder.String()
}

func (f *stylish) addDiff(key string, diff diff) {
	if key == "" {
		f.addIndeted(' ', "{\n")
	} else {
		f.addIndeted(' ', "%s: {\n", key)
	}

	keys := slices.Sorted(maps.Keys(diff.child))

	f.indent += 1
	for _, key := range keys {
		diffItem := diff.child[key]

		switch diffItem.change {
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

func (f *stylish) addUnchanged(key string, diff diff) {
	f.addIndeted(' ', "%s: ", key)
	f.addValue(diff.valueLeft)
}

func (f *stylish) addAdded(key string, diff diff) {
	f.addIndeted('+', "%s: ", key)
	f.addValue(diff.valueRight)
}

func (f *stylish) addRemoved(key string, diff diff) {
	f.addIndeted('-', "%s: ", key)
	f.addValue(diff.valueLeft)
}

func (f *stylish) addValueChanged(key string, diff diff) {
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

type plain struct {
	rootDiff diff
	builder  strings.Builder
	scope    []string
}

func newPlain(diff diff) *plain {
	result := &plain{rootDiff: diff}
	return result
}

func (f *plain) build() string {
	f.addDiff(f.rootDiff.key, f.rootDiff)
	return f.builder.String()
}

func (f *plain) addDiff(key string, diff diff) {
	keys := slices.Sorted(maps.Keys(diff.child))

	for _, key := range keys {
		diffItem := diff.child[key]
		f.scope = append(f.scope, key)

		switch diffItem.change {
		case "added":
			f.addAdded(diffItem)
		case "removed":
			f.addRemoved()
		case "value_changed":
			f.addValueChanged(diffItem)
		case "diff":
			f.addDiff(key, diffItem)
		}

		f.scope = f.scope[:len(f.scope)-1]
	}
}

func (f *plain) addAdded(diff diff) {
	key := strings.Join(f.scope, ".")
	value := plainValue(diff.valueRight)

	fmt.Fprintf(&f.builder, "Property '%s' was added with value: %s\n", key, value)
}

func (f *plain) addRemoved() {
	key := strings.Join(f.scope, ".")

	fmt.Fprintf(&f.builder, "Property '%s' was removed\n", key)
}

func (f *plain) addValueChanged(diff diff) {
	key := strings.Join(f.scope, ".")
	valueLeft := plainValue(diff.valueLeft)
	valueRight := plainValue(diff.valueRight)

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

func getFormater(format string, diff diff) (formater, error) {
	switch format {
	case "stylish":
		return newStylish(diff), nil
	case "plain":
		return newPlain(diff), nil
	default:
		return nil, UnsupportedFormat{Format: format}
	}
}
