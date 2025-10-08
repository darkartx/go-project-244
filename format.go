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
	unchanged(key string, value any)
	added(key string, value any)
	removed(key string, value any)
	valueChanged(key string, left any, right any)
	diffStart(scope string, value diff)
	diffEnd()
	build() string
}

type stylish struct {
	diff    map[string]diff
	scope   uint8
	builder strings.Builder
}

func newStylish(diff map[string]diff) *stylish {
	result := &stylish{diff: diff}
	fmt.Fprint(&result.builder, "{\n")
	return result
}

func (f *stylish) unchanged(key string, value any) {
	f.addLine("    %s: %v", key, value)
}

func (f *stylish) added(key string, value any) {
	switch value := value.(type) {
	case map[string]any:
		f.addLine("  + %s: {", key)
		// f.addMap(0, value)
		// f.addLine("}")
	default:
		f.addLine("  + %s: %v", key, value)
	}
}

func (f *stylish) removed(key string, value any) {
	switch value := value.(type) {
	case map[string]any:
		f.addLine("  - %s: {", key)
		f.addMap(1, value)
		f.addLine("}")
	default:
		f.addLine("  - %s: %v", key, value)
	}
}

func (f *stylish) valueChanged(key string, left any, right any) {
	f.removed(key, left)
	f.added(key, right)
}

func (f *stylish) diffStart(scope string, value diff) {

}

func (f *stylish) diffEnd() {

}

func (f *stylish) build() string {
	fmt.Fprint(&f.builder, "}")
	return f.builder.String()
}

func (f *stylish) addLine(format string, values ...any) {
	fmt.Fprintf(&f.builder, "%*s", f.scope*4, " ")
	fmt.Fprintf(&f.builder, format, values...)
	fmt.Fprint(&f.builder, "\n")
}

func (f *stylish) addMap(indent uint8, value map[string]any) {
	keys := slices.Sorted(maps.Keys(value))
	indent += 1

	for _, key := range keys {
		item := value[key]
		switch item := item.(type) {
		case map[string]any:
			f.addLine("%*s%s: {", indent*4, " ", key)
			f.addMap(indent, item)
			f.addLine("%*s}", indent*4, " ")
		default:
			f.addLine("%*s: %v", indent*4, key, item)
		}
	}
}

func getFormater(format string, diff map[string]diff) (formater, error) {
	switch format {
	case "stylish":
		return newStylish(diff), nil
	default:
		return nil, UnsupportedFormat{Format: format}
	}
}
