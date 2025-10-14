package formatters

import (
	"fmt"

	"github.com/darkartx/go-project-244/shared"
)

type UnsupportedFormat struct {
	Format string
}

func (e UnsupportedFormat) Error() string {
	return fmt.Sprintf("unsupported format: %s", e.Format)
}

type Formatter interface {
	Build() string
}

func GetFormater(format string, diff shared.Diff) (Formatter, error) {
	switch format {
	case "stylish":
		return newStylish(diff), nil
	case "plain":
		return newPlain(diff), nil
	default:
		return nil, UnsupportedFormat{Format: format}
	}
}
