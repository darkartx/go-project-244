package internal

import (
	"reflect"
)

const (
	DIFF_CHANGE_UNCAHGNED = iota
	DIFF_CHANGE_VALUE_CHANGED
	DIFF_CHANGE_VALUE_ADDED
	DIFF_CHANGE_VALUE_REMOVED
	DIFF_CHANGE_DIFF
)

type Empty struct{}

type Diff struct {
	Key        string
	Change     int
	LeftValue  any
	RightValue any
	Child      map[string]Diff
}

func NewDiff(key string, change int, leftValue, rightValue any, child map[string]Diff) Diff {
	return Diff{
		key, change, leftValue, rightValue, child,
	}
}

func BuildDiff(dataLeft, dataRight map[string]any) Diff {
	result := make(map[string]Diff)
	keys := make(map[string]bool)

	for k := range dataLeft {
		keys[k] = true
	}

	for k := range dataRight {
		if !keys[k] {
			keys[k] = true
		}
	}

	for key := range keys {
		valueLeft, existsLeft := dataLeft[key]
		if !existsLeft {
			valueLeft = Empty{}
		}

		valueRight, existsRight := dataRight[key]
		if !existsRight {
			valueRight = Empty{}
		}

		result[key] = MakeDiff(key, valueLeft, valueRight)
	}

	return NewDiff("", DIFF_CHANGE_DIFF, dataLeft, dataRight, result)
}

func MakeDiff(key string, valueLeft, valueRight any) Diff {
	typeLeft := reflect.TypeOf(valueLeft)
	typeRight := reflect.TypeOf(valueRight)
	change := DIFF_CHANGE_UNCAHGNED
	var child map[string]Diff

	switch {
	case typeLeft == nil || typeRight == nil:
		if typeLeft != typeRight {
			change = DIFF_CHANGE_VALUE_CHANGED
		}
	case typeLeft == reflect.TypeOf(Empty{}):
		change = DIFF_CHANGE_VALUE_ADDED
	case typeRight == reflect.TypeOf(Empty{}):
		change = DIFF_CHANGE_VALUE_REMOVED
	case typeLeft.Kind() == reflect.Map && typeRight.Kind() == reflect.Map:
		result := BuildDiff(valueLeft.(map[string]any), valueRight.(map[string]any))
		result.Key = key
		return result
	case typeLeft.Kind() == reflect.Map && typeRight.Kind() != reflect.Map:
		fallthrough
	case typeLeft.Kind() != reflect.Map && typeRight.Kind() == reflect.Map:
		change = DIFF_CHANGE_VALUE_CHANGED
	default:
		if valueLeft != valueRight {
			change = DIFF_CHANGE_VALUE_CHANGED
		}
	}

	return NewDiff(key, change, valueLeft, valueRight, child)
}
