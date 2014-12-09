package gohaml_test

import (
	"encoding/json"
	"errors"
	"github.com/realistschuckle/gohaml"
	"os"
	"strings"
	"testing"
)

type test struct {
	Haml     string
	Html     string
	Optional bool
	Config   map[string]string
	Locals   map[string]interface{}
}

type results struct {
	passed int
	failed int
}

func replaceNewlines(val string) (new string) {
	new = strings.Replace(val, "\n", "\\n", -1)
	new = strings.Replace(new, "\r", "\\r", -1)
	return
}

func str2bool(s string) (b bool, err error) {
	switch s {
	case "true":
		b = true
	case "false":
		b = false
	default:
		err = errors.New("Unknown bool value: " + s)
	}
	return
}

func mergeConfig(config map[string]string) (options gohaml.EngineOptions, err error) {
	options = gohaml.DefaultEngineOptions()
	if config != nil {
		for name, val := range config {
			switch name {
			case "escape_html":
				options.EscapeHtml, err = str2bool(val)
			case "format":
				options.Format = val
			}
		}
	}
	return
}

func TestSpecifications(t *testing.T) {
	var file *os.File
	var err error

	if file, err = os.Open("tests.json"); err != nil {
		t.Fatalf(err.Error())
	}

	var tests map[string]map[string]test

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&tests); err != nil {
		t.Fatalf(err.Error())
	}

	res := results{}
	for categoryName := range tests {
		failed := false
		for testName, test := range tests[categoryName] {
			var opts gohaml.EngineOptions
			var err error
			if opts, err = mergeConfig(test.Config); err != nil {
				t.Fatalf("\t%v\n", err.Error())
			}
			if engine, err := gohaml.NewEngine(test.Haml, &opts); err != nil {
				t.Errorf("  ERROR IN '%s': %s\n", testName, err.Error())
				res.failed += 1
			} else {
				output, err := engine.Render(test.Locals)
				switch {
				case err != nil:
					t.Errorf("  ERROR IN '%s': %s\n", testName, err.Error())
					res.failed += 1
				case output != test.Html:
					if !failed {
						failed = true
						t.Log(categoryName)
					}
					t.Errorf("  %s\n", testName)
					t.Errorf("    input   : %s\n", replaceNewlines(test.Haml))
					t.Errorf("    expected: %s\n", replaceNewlines(test.Html))
					t.Errorf("    got     : %s\n", replaceNewlines(output))
					res.failed += 1
				default:
					res.passed += 1
				}
			}
		}
	}
	t.Logf("%d tests, %d failures", res.passed+res.failed, res.failed)
}
