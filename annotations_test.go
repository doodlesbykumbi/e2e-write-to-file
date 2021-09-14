package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type annotationTestCase struct {
	description string
	contents    string
	assert      func(t *testing.T, result map[string]string, err error)
}

func assertGoodAnnotations(expectedResult map[string]string) func (*testing.T, map[string]string, error) {
	return func(t *testing.T, result map[string]string, err error) {
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(
			t,
			expectedResult,
			result,
		)
	}
}

var annotationTestCases = []annotationTestCase{
	{
		description: "valid example",
		contents: `
example1="example1 is a simple value"
example2="example2 is\na more \"complex\" value"
`,
		assert: assertGoodAnnotations(
			map[string]string{
				"example1": "example1 is a simple value",
				"example2": `example2 is
a more "complex" value`,
			},
		),
	},
	{
		description: "malformed line without equals",
		contents: `example1="example1 is a simple value"
this line has no equal sign
example2="example2 is\na more \"complex\" value"
`,
		assert: func(t *testing.T, result map[string]string, err error) {
			assert.Contains(t, err.Error(), "line 2 is malformed")
		},
	},
	{
		description: "malformed line without quoted value",
		contents: `example1="example1 is a simple value"
this=line has no quoted value
example2="example2 is\na more \"complex\" value"
`,
		assert: func(t *testing.T, result map[string]string, err error) {
			assert.Contains(t, err.Error(), "line 2 is malformed")
		},
	},
}

func TestNewAnnotations(t *testing.T) {
	for _, tc := range annotationTestCases {
		t.Run(tc.description, func(t *testing.T) {
			annotations, err := NewAnnotations(strings.NewReader(tc.contents))
			tc.assert(t, annotations, err)
		})
	}
}
