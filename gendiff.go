package code

import (
	"reflect"
	"slices"

	"github.com/darkartx/go-project-244/formatters"
	"github.com/darkartx/go-project-244/shared"
)

type empty struct{}

func GenDiff(filepathLeft, filepathRight, format string) (string, error) {
	if format == "" {
		format = "stylish"
	}

	dataLeft, err := parseFile(filepathLeft)
	if err != nil {
		return "", err
	}

	dataRight, err := parseFile(filepathRight)
	if err != nil {
		return "", err
	}

	diff := genDiff(dataLeft, dataRight)

	formater, err := formatters.GetFormater(format, diff)
	if err != nil {
		return "", err
	}

	return formater.Build(), nil
}

func genDiff(dataLeft, dataRight map[string]any) shared.Diff {
	result := make(map[string]shared.Diff)
	var keys []string

	for k := range dataLeft {
		keys = append(keys, k)
	}

	for k := range dataRight {
		if !slices.Contains(keys, k) {
			keys = append(keys, k)
		}
	}

	for _, key := range keys {
		valueLeft, existsLeft := dataLeft[key]
		if !existsLeft {
			valueLeft = empty{}
		}

		valueRight, existsRight := dataRight[key]
		if !existsRight {
			valueRight = empty{}
		}

		result[key] = makeDiff(key, valueLeft, valueRight)
	}

	return shared.NewDiff("", "diff", dataLeft, dataRight, result)
}

func makeDiff(key string, valueLeft, valueRight any) shared.Diff {
	typeLeft := reflect.TypeOf(valueLeft)
	typeRight := reflect.TypeOf(valueRight)
	change := "unchanged"
	var child map[string]shared.Diff

	switch {
	case typeLeft == nil || typeRight == nil:
		if typeLeft != typeRight {
			change = "value_changed"
		}
	case typeLeft == reflect.TypeOf(empty{}):
		change = "added"
	case typeRight == reflect.TypeOf(empty{}):
		change = "removed"
	case typeLeft.Kind() == reflect.Map && typeRight.Kind() == reflect.Map:
		change = "diff"
		child = genDiff(valueLeft.(map[string]any), valueRight.(map[string]any)).Child
	case typeLeft.Kind() == reflect.Map && typeRight.Kind() != reflect.Map:
		fallthrough
	case typeLeft.Kind() != reflect.Map && typeRight.Kind() == reflect.Map:
		change = "value_changed"
	default:
		if valueLeft != valueRight {
			change = "value_changed"
		}
	}

	return shared.NewDiff(key, change, valueLeft, valueRight, child)
}
