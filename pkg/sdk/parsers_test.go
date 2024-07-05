package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommaSeparatedStringArray(t *testing.T) {
	testCases := []struct {
		Name       string
		Value      string
		TrimQuotes bool
		Result     []string
	}{
		{
			Name:   "empty list",
			Value:  "[]",
			Result: []string{},
		},
		{
			Name:   "empty string",
			Value:  "",
			Result: []string{},
		},
		{
			Name:   "one element in list",
			Value:  "[one]",
			Result: []string{"one"},
		},
		{
			Name:       "one element in list - with quotes",
			Value:      "['one']",
			TrimQuotes: true,
			Result:     []string{"one"},
		},
		{
			Name:       "multiple elements in list - with quotes",
			Value:      "['one', 'two', 'three']",
			TrimQuotes: true,
			Result:     []string{"one", "two", "three"},
		},
		{
			Name:   "multiple elements in list",
			Value:  "[one, two, three]",
			Result: []string{"one", "two", "three"},
		},
		{
			Name:   "multiple elements in list - packed",
			Value:  "[one,two,three]",
			Result: []string{"one", "two", "three"},
		},
		{
			Name:   "multiple elements in list - additional spaces",
			Value:  "[one    ,          two  ,three]",
			Result: []string{"one", "two", "three"},
		},
		{
			Name:   "list without brackets",
			Value:  "one,two,three",
			Result: []string{"one", "two", "three"},
		},
		{
			Name:       "list without brackets - with quotes",
			Value:      "'one','two','three'",
			TrimQuotes: true,
			Result:     []string{"one", "two", "three"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Result, ParseCommaSeparatedStringArray(tc.Value, tc.TrimQuotes))
		})
	}
}
