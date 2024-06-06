package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchAllStringsInOrderNonOverlapping(t *testing.T) {
	testCases := map[string]struct {
		parts     []string
		text      string
		wantMatch bool
	}{
		"empty parts and text": {
			parts:     []string{},
			text:      "",
			wantMatch: true,
		},
		"empty parts": {
			parts:     []string{},
			text:      "xyz",
			wantMatch: true,
		},
		"empty text": {
			parts: []string{"a", "b"},
			text:  "",
		},
		"matching non empty": {
			parts:     []string{"a", "b"},
			text:      "xyaxyb",
			wantMatch: true,
		},
		"partial matching": {
			parts: []string{"a", "b"},
			text:  "axyz",
		},
		"not matching": {
			parts: []string{"a", "b"},
			text:  "xyz",
		},
		"wrong order": {
			parts: []string{"a", "b"},
			text:  "ba",
		},
		"overlapping match": {
			parts: []string{"abb", "bba"},
			text:  "abba",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			regex := MatchAllStringsInOrderNonOverlapping(tc.parts)
			require.Equal(t, tc.wantMatch, regex.Match([]byte(tc.text)))
		})
	}
}
