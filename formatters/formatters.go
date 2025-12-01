package formatters

import (
	"fmt"

	"github.com/darkartx/go-project-244/internal"
)

type UnsupportedFormat struct {
	Format string
}

func (e UnsupportedFormat) Error() string {
	return fmt.Sprintf("unsupported format: %s", e.Format)
}

type Formatter interface {
	Build() (string, error)
}

func GetFormater(format string, diff internal.Diff) (Formatter, error) {
	switch format {
	case "stylish":
		return NewStylish(diff), nil
	case "plain":
		return NewPlain(diff), nil
	case "json":
		return NewJson(diff), nil
	default:
		return nil, UnsupportedFormat{Format: format}
	}
}
