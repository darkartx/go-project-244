package code

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenDiff_JsonStylish(t *testing.T) {
	cases := []struct {
		leftFile, rightFile string
		expected            string
	}{
		// Flat json
		{
			"file1.json",
			"file2.json",
			`{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`,
		},
		// Recursive json
		{
			"file1_recursive.json",
			"file2_recursive.json",
			`{
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
}`,
		},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%s_%s", c.leftFile, c.rightFile)

		t.Run(name, func(t *testing.T) {
			leftFile := filepath.Join("testdata", "fixture", c.leftFile)
			rightFile := filepath.Join("testdata", "fixture", c.rightFile)
			actual, err := GenDiff(leftFile, rightFile, "stylish")
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestGenDiff_JsonPlain(t *testing.T) {
	cases := []struct {
		leftFile, rightFile string
		expected            string
	}{
		// Flat json
		{
			"file1.json",
			"file2.json",
			`Property 'follow' was removed
Property 'proxy' was removed
Property 'timeout' was updated. From 50 to 20
Property 'verbose' was added with value: true`,
		},
		// Recursive json
		{
			"file1_recursive.json",
			"file2_recursive.json",
			`Property 'common.follow' was added with value: false
Property 'common.setting2' was removed
Property 'common.setting3' was updated. From true to null
Property 'common.setting4' was added with value: 'blah blah'
Property 'common.setting5' was added with value: [complex value]
Property 'common.setting6.doge.wow' was updated. From '' to 'so much'
Property 'common.setting6.ops' was added with value: 'vops'
Property 'group1.baz' was updated. From 'bas' to 'bars'
Property 'group1.nest' was updated. From [complex value] to 'str'
Property 'group2' was removed
Property 'group3' was added with value: [complex value]`,
		},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%s_%s", c.leftFile, c.rightFile)

		t.Run(name, func(t *testing.T) {
			leftFile := filepath.Join("testdata", "fixture", c.leftFile)
			rightFile := filepath.Join("testdata", "fixture", c.rightFile)
			actual, err := GenDiff(leftFile, rightFile, "plain")
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestGenDiff_JsonJson(t *testing.T) {
	cases := []struct {
		leftFile, rightFile string
		expected            string
	}{
		// Flat json
		{
			"file1.json",
			"file2.json",
			`{
  "follow": {
    "change": "removed",
    "left_value": false
  },
  "host": {
    "change": "unchanged",
    "left_value": "hexlet.io"
  },
  "proxy": {
    "change": "removed",
    "left_value": "123.234.53.22"
  },
  "timeout": {
    "change": "value_changed",
    "left_value": 50,
    "right_value": 20
  },
  "verbose": {
    "change": "added",
    "right_value": true
  }
}`,
		},
		// Recursive json
		{
			"file1_recursive.json",
			"file2_recursive.json",
			`{
  "common": {
    "change": "diff",
    "diff": {
      "follow": {
        "change": "added",
        "right_value": false
      },
      "setting1": {
        "change": "unchanged",
        "left_value": "Value 1"
      },
      "setting2": {
        "change": "removed",
        "left_value": 200
      },
      "setting3": {
        "change": "value_changed",
        "left_value": true,
        "right_value": null
      },
      "setting4": {
        "change": "added",
        "right_value": "blah blah"
      },
      "setting5": {
        "change": "added",
        "right_value": {
          "key5": "value5"
        }
      },
      "setting6": {
        "change": "diff",
        "diff": {
          "doge": {
            "change": "diff",
            "diff": {
              "wow": {
                "change": "value_changed",
                "left_value": "",
                "right_value": "so much"
              }
            }
          },
          "key": {
            "change": "unchanged",
            "left_value": "value"
          },
          "ops": {
            "change": "added",
            "right_value": "vops"
          }
        }
      }
    }
  },
  "group1": {
    "change": "diff",
    "diff": {
      "baz": {
        "change": "value_changed",
        "left_value": "bas",
        "right_value": "bars"
      },
      "foo": {
        "change": "unchanged",
        "left_value": "bar"
      },
      "nest": {
        "change": "value_changed",
        "left_value": {
          "key": "value"
        },
        "right_value": "str"
      }
    }
  },
  "group2": {
    "change": "removed",
    "left_value": {
      "abc": 12345,
      "deep": {
        "id": 45
      }
    }
  },
  "group3": {
    "change": "added",
    "right_value": {
      "deep": {
        "id": {
          "number": 45
        }
      },
      "fee": 100500
    }
  }
}`,
		},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%s_%s", c.leftFile, c.rightFile)

		t.Run(name, func(t *testing.T) {
			leftFile := filepath.Join("testdata", "fixture", c.leftFile)
			rightFile := filepath.Join("testdata", "fixture", c.rightFile)
			actual, err := GenDiff(leftFile, rightFile, "json")
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestGenDiff_ErrorBadJson(t *testing.T) {
	_, err := GenDiff("testdata/fixture/file1.json", "testdata/fixture/bad_json.json", "stylish")

	assert.Error(t, err)
}

func TestGenDiff_YamlStylish(t *testing.T) {
	cases := []struct {
		leftFile, rightFile string
		expected            string
	}{
		// Flat yaml
		{
			"file1.yml",
			"file2.yml",
			`{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`,
		},
		// Recursive yaml
		{
			"file1_recursive.yml",
			"file2_recursive.yml",
			`{
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
}`,
		},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%s_%s", c.leftFile, c.rightFile)

		t.Run(name, func(t *testing.T) {
			leftFile := filepath.Join("testdata", "fixture", c.leftFile)
			rightFile := filepath.Join("testdata", "fixture", c.rightFile)
			actual, err := GenDiff(leftFile, rightFile, "stylish")
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestGenDiff_YamlPlain(t *testing.T) {
	cases := []struct {
		leftFile, rightFile string
		expected            string
	}{
		// Flat yaml
		{
			"file1.yml",
			"file2.yml",
			`Property 'follow' was removed
Property 'proxy' was removed
Property 'timeout' was updated. From 50 to 20
Property 'verbose' was added with value: true`,
		},
		// Recursive yaml
		{
			"file1_recursive.yml",
			"file2_recursive.yml",
			`Property 'common.follow' was added with value: false
Property 'common.setting2' was removed
Property 'common.setting3' was updated. From true to null
Property 'common.setting4' was added with value: 'blah blah'
Property 'common.setting5' was added with value: [complex value]
Property 'common.setting6.doge.wow' was updated. From '' to 'so much'
Property 'common.setting6.ops' was added with value: 'vops'
Property 'group1.baz' was updated. From 'bas' to 'bars'
Property 'group1.nest' was updated. From [complex value] to 'str'
Property 'group2' was removed
Property 'group3' was added with value: [complex value]`,
		},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%s_%s", c.leftFile, c.rightFile)

		t.Run(name, func(t *testing.T) {
			leftFile := filepath.Join("testdata", "fixture", c.leftFile)
			rightFile := filepath.Join("testdata", "fixture", c.rightFile)
			actual, err := GenDiff(leftFile, rightFile, "plain")
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestGenDiff_YamlJson(t *testing.T) {
	cases := []struct {
		leftFile, rightFile string
		expected            string
	}{
		// Flat yaml
		{
			"file1.yml",
			"file2.yml",
			`{
  "follow": {
    "change": "removed",
    "left_value": false
  },
  "host": {
    "change": "unchanged",
    "left_value": "hexlet.io"
  },
  "proxy": {
    "change": "removed",
    "left_value": "123.234.53.22"
  },
  "timeout": {
    "change": "value_changed",
    "left_value": 50,
    "right_value": 20
  },
  "verbose": {
    "change": "added",
    "right_value": true
  }
}`,
		},
		// Recursive yaml
		{
			"file1_recursive.yml",
			"file2_recursive.yml",
			`{
  "common": {
    "change": "diff",
    "diff": {
      "follow": {
        "change": "added",
        "right_value": false
      },
      "setting1": {
        "change": "unchanged",
        "left_value": "Value 1"
      },
      "setting2": {
        "change": "removed",
        "left_value": 200
      },
      "setting3": {
        "change": "value_changed",
        "left_value": true,
        "right_value": null
      },
      "setting4": {
        "change": "added",
        "right_value": "blah blah"
      },
      "setting5": {
        "change": "added",
        "right_value": {
          "key5": "value5"
        }
      },
      "setting6": {
        "change": "diff",
        "diff": {
          "doge": {
            "change": "diff",
            "diff": {
              "wow": {
                "change": "value_changed",
                "left_value": "",
                "right_value": "so much"
              }
            }
          },
          "key": {
            "change": "unchanged",
            "left_value": "value"
          },
          "ops": {
            "change": "added",
            "right_value": "vops"
          }
        }
      }
    }
  },
  "group1": {
    "change": "diff",
    "diff": {
      "baz": {
        "change": "value_changed",
        "left_value": "bas",
        "right_value": "bars"
      },
      "foo": {
        "change": "unchanged",
        "left_value": "bar"
      },
      "nest": {
        "change": "value_changed",
        "left_value": {
          "key": "value"
        },
        "right_value": "str"
      }
    }
  },
  "group2": {
    "change": "removed",
    "left_value": {
      "abc": 12345,
      "deep": {
        "id": 45
      }
    }
  },
  "group3": {
    "change": "added",
    "right_value": {
      "deep": {
        "id": {
          "number": 45
        }
      },
      "fee": 100500
    }
  }
}`,
		},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%s_%s", c.leftFile, c.rightFile)

		t.Run(name, func(t *testing.T) {
			leftFile := filepath.Join("testdata", "fixture", c.leftFile)
			rightFile := filepath.Join("testdata", "fixture", c.rightFile)
			actual, err := GenDiff(leftFile, rightFile, "json")
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.expected, actual)
		})
	}
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
