package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchAllStringsInOrder(t *testing.T) {
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
		"not matching": {
			parts: []string{"a", "b"},
			text:  "axyz",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			regex := MatchAllStringsInOrder(tc.parts)
			require.Equal(t, tc.wantMatch, regex.Match([]byte(tc.text)))
		})
	}
}
