package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenDiff_Json(t *testing.T) {
	actual, err := GenDiff("testdata/fixture/file1.json", "testdata/fixture/file2.json", "stylish")
	if err != nil {
		t.Fatal(err)
	}

	expected := `{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`

	assert.Equal(t, expected, actual)
}

func TestGenDiff_ErrorBadJson(t *testing.T) {
	_, err := GenDiff("testdata/fixture/file1.json", "testdata/fixture/bad_json.json", "stylish")

	assert.Error(t, err)
}

func TestGenDiff_Yaml(t *testing.T) {
	actual, err := GenDiff("testdata/fixture/file1.yml", "testdata/fixture/file2.yml", "stylish")
	if err != nil {
		t.Fatal(err)
	}

	expected := `{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`

	assert.Equal(t, expected, actual)
}

func TestGenDiff_ErrorBadYaml(t *testing.T) {
	_, err := GenDiff("testdata/fixture/file1.yml", "testdata/fixture/bad_yml.yml", "stylish")

	assert.Error(t, err)
}

func TestGenDiff_ErrorUnsupportedFormat(t *testing.T) {
	_, err := GenDiff("testdata/fixture/file1.toml", "testdata/fixture/file2.json", "stylish")

	assert.EqualError(t, err, "unsupported format: toml")
}

func TestGenDiff_ErrorNoFile(t *testing.T) {
	_, err := GenDiff("testdata/fixture/no_file.json", "testdata/fixture/file2.json", "stylish")

	assert.Error(t, err)
}
