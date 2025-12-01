package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type UnsupportedFileFormatError struct {
	Format string
}

func (e UnsupportedFileFormatError) Error() string {
	return fmt.Sprintf("unsupported file format %s", e.Format)
}

type ParseError struct {
	FileName string
	Previous error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("parsing %s: %s", e.FileName, e.Previous)
}

func ParseFile(fp string) (map[string]any, error) {
	format := filepath.Ext(fp)
	name := filepath.Base(fp)

	if len(format) > 0 {
		format = format[1:]
	}

	data, err := os.ReadFile(fp)

	if err != nil {
		return nil, ParseError{FileName: name, Previous: err}
	}

	result, err := parse(data, format)
	if err != nil {
		return nil, ParseError{FileName: name, Previous: err}
	}

	return result, nil
}

func parse(data []byte, format string) (map[string]any, error) {
	switch format {
	case "json":
		return parseJson(data)
	case "yaml", "yml":
		return parseYaml(data)
	default:
		return nil, UnsupportedFileFormatError{Format: format}
	}
}

func parseJson(data []byte) (map[string]any, error) {
	var result map[string]any

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func parseYaml(data []byte) (map[string]any, error) {
	var result map[string]any

	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}
