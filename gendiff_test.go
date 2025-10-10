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
}
`

	assert.Equal(t, expected, actual)
}

func TestGenDiff_JsonRecurcive(t *testing.T) {
	actual, err := GenDiff("testdata/fixture/file1_recursive.json", "testdata/fixture/file2_recursive.json", "stylish")
	if err != nil {
		t.Fatal(err)
	}

	expected := `{
    common: {
      + follow: false
        setting1: Value 1
      - setting2: 200
      - setting3: true
      + setting3: null
      + setting4: blah blah
      + setting5: {
            key5: value5
        }
        setting6: {
            doge: {
              - wow: 
              + wow: so much
            }
            key: value
          + ops: vops
        }
    }
    group1: {
      - baz: bas
      + baz: bars
        foo: bar
      - nest: {
            key: value
        }
      + nest: str
    }
  - group2: {
        abc: 12345
        deep: {
            id: 45
        }
    }
  + group3: {
        deep: {
            id: {
                number: 45
            }
        }
        fee: 100500
    }
}
`

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
}
`

	assert.Equal(t, expected, actual)
}

func TestGenDiff_YamlRecurcive(t *testing.T) {
	actual, err := GenDiff("testdata/fixture/file1_recursive.yml", "testdata/fixture/file2_recursive.yml", "stylish")
	if err != nil {
		t.Fatal(err)
	}

	expected := `{
    common: {
      + follow: false
        setting1: Value 1
      - setting2: 200
      - setting3: true
      + setting3: null
      + setting4: blah blah
      + setting5: {
            key5: value5
        }
        setting6: {
            doge: {
              - wow: 
              + wow: so much
            }
            key: value
          + ops: vops
        }
    }
    group1: {
      - baz: bas
      + baz: bars
        foo: bar
      - nest: {
            key: value
        }
      + nest: str
    }
  - group2: {
        abc: 12345
        deep: {
            id: 45
        }
    }
  + group3: {
        deep: {
            id: {
                number: 45
            }
        }
        fee: 100500
    }
}
`

	assert.Equal(t, expected, actual)
}

func TestGenDiff_ErrorBadYaml(t *testing.T) {
	_, err := GenDiff("testdata/fixture/file1.yml", "testdata/fixture/bad_yml.yml", "stylish")

	assert.Error(t, err)
}

func TestGenDiff_ErrorUnsupportedFileFormat(t *testing.T) {
	_, err := GenDiff("testdata/fixture/file1.toml", "testdata/fixture/file2.json", "stylish")

	assert.EqualError(t, err, "unsupported file format: toml")
}

func TestGenDiff_ErrorNoFile(t *testing.T) {
	_, err := GenDiff("testdata/fixture/no_file.json", "testdata/fixture/file2.json", "stylish")

	assert.Error(t, err)
}
