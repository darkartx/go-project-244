package code

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

type UnsupportedFileFormat struct {
	Format string
}

func (e UnsupportedFileFormat) Error() string {
	return fmt.Sprintf("unsupported file format: %s", e.Format)
}

func parseFile(fp string) (map[string]any, error) {
	format := filepath.Ext(fp)

	if len(format) > 0 {
		format = format[1:]
	}

	data, err := os.ReadFile(fp)

	if err != nil {
		return nil, err
	}

	return parse(data, format)
}

func parse(data []byte, format string) (map[string]any, error) {
	switch format {
	case "json":
		return parseJson(data)
	case "yaml", "yml":
		return parseYaml(data)
	default:
		return nil, UnsupportedFileFormat{Format: format}
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
	var result map[any]any

	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	normalizedResult := normalizeData(result)

	return normalizedResult, nil
}

func normalizeData(data map[any]any) map[string]any {
	result := make(map[string]any)

	for k, v := range data {
		keyStr := fmt.Sprintf("%v", k)

		switch val := v.(type) {
		case map[any]any:
			result[keyStr] = normalizeData(val)
		default:
			result[keyStr] = val
		}
	}

	return result
}
