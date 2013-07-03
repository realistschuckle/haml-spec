package gohaml_test

import (
	"encoding/json"
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
	emptyScope := make(map[string]interface{})
	for categoryName := range tests {
		t.Log(categoryName)
		for testName, test := range tests[categoryName] {
			if engine, err := gohaml.NewEngine(test.Haml); err != nil {
				t.Error(err.Error())
				res.failed += 1
			} else {
				scope := emptyScope
				if test.Locals != nil {
					scope = test.Locals
				}
				if output := engine.Render(scope); output != test.Html {
					t.Errorf("  %s\n", testName)
					t.Errorf("    input   : %s\n", replaceNewlines(test.Haml))
					t.Errorf("    expected: %s\n", replaceNewlines(test.Html))
					t.Errorf("    got     : %s\n", replaceNewlines(output))
					res.failed += 1
				} else {
					res.passed += 1
				}
			}
		}
	}
	t.Logf("%d tests, %d failures", res.passed+res.failed, res.failed)
}
