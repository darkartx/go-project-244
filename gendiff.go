package code

import (
	"reflect"
	"slices"
)

type empty struct{}

type diff struct {
	key        string
	change     string
	valueLeft  any
	valueRight any
	child      map[string]diff
}

func GenDiff(filepathLeft, filepathRight, format string) (string, error) {
	dataLeft, err := parseFile(filepathLeft)
	if err != nil {
		return "", err
	}

	dataRight, err := parseFile(filepathRight)
	if err != nil {
		return "", err
	}

	diff := genDiff(dataLeft, dataRight)

	formater, err := getFormater(format, diff)
	if err != nil {
		return "", err
	}

	return formater.build(), nil
}

func genDiff(dataLeft, dataRight map[string]any) diff {
	result := make(map[string]diff)
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

	return diff{
		key:        "",
		change:     "diff",
		valueLeft:  dataLeft,
		valueRight: dataRight,
		child:      result,
	}
}

func makeDiff(key string, valueLeft, valueRight any) diff {
	typeLeft := reflect.TypeOf(valueLeft)
	typeRight := reflect.TypeOf(valueRight)
	change := "unchanged"
	var child map[string]diff

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
		child = genDiff(valueLeft.(map[string]any), valueRight.(map[string]any)).child
	case typeLeft.Kind() == reflect.Map && typeRight.Kind() != reflect.Map:
		fallthrough
	case typeLeft.Kind() != reflect.Map && typeRight.Kind() == reflect.Map:
		change = "value_changed"
	default:
		if valueLeft != valueRight {
			change = "value_changed"
		}
	}

	return diff{key, change, valueLeft, valueRight, child}
}
