package gen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToSnakeCase(t *testing.T) {

	type test struct {
		input    string
		expected string
	}

	testCases := []test{
		{input: "CamelCase", expected: "camel_case"},
		{input: "ACamelCase", expected: "a_camel_case"},
		{input: "URLParser", expected: "url_parser"},
		{input: "Camel1Case", expected: "camel1_case"},
		{input: "camelCase", expected: "camel_case"},
		{input: "camelURL", expected: "camel_url"},
		{input: "camelURLSomething", expected: "camel_url_something"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s=>%s", tc.input, tc.expected), func(t *testing.T) {
			result := ToSnakeCase(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
