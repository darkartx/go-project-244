package code

import (
	"code/formatters"
	"code/internal"
)

func GenDiff(filepathLeft, filepathRight, format string) (string, error) {
	if format == "" {
		format = "stylish"
	}

	dataLeft, err := internal.ParseFile(filepathLeft)
	if err != nil {
		return "", err
	}

	dataRight, err := internal.ParseFile(filepathRight)
	if err != nil {
		return "", err
	}

	diff := internal.BuildDiff(dataLeft, dataRight)

	formater, err := formatters.GetFormater(format, diff)
	if err != nil {
		return "", err
	}

	return formater.Build()
}
