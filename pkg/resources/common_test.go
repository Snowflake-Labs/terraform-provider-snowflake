package resources

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_suppressIdentifierQuoting(t *testing.T) {
	firstId := "a.b.c"
	firstIdQuoted := "\"a\".b.\"c\""
	secondId := "d.e.f"
	incorrectId := "a.b.c.d.e.f"

	t.Run("old identifier with too many parts", func(t *testing.T) {
		result := suppressIdentifierQuoting("", incorrectId, firstId, nil)
		require.False(t, result)
	})

	t.Run("new identifier with too many parts", func(t *testing.T) {
		result := suppressIdentifierQuoting("", firstId, incorrectId, nil)
		require.False(t, result)
	})

	t.Run("old identifier empty", func(t *testing.T) {
		result := suppressIdentifierQuoting("", "", firstId, nil)
		require.False(t, result)
	})

	t.Run("new identifier empty", func(t *testing.T) {
		result := suppressIdentifierQuoting("", firstId, "", nil)
		require.False(t, result)
	})

	t.Run("identifiers the same", func(t *testing.T) {
		result := suppressIdentifierQuoting("", firstId, firstId, nil)
		require.True(t, result)
	})

	t.Run("identifiers the same (but different quoting)", func(t *testing.T) {
		result := suppressIdentifierQuoting("", firstId, firstIdQuoted, nil)
		require.True(t, result)
	})

	t.Run("identifiers different", func(t *testing.T) {
		result := suppressIdentifierQuoting("", firstId, secondId, nil)
		require.False(t, result)
	})
}

func Test_listValueToSlice(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		trimQuotes bool
		want       []string
	}{
		{
			name: "empty list",
			want: nil,
		},
		{
			name:  "empty list with brackets",
			value: "[]",
			want:  nil,
		},
		{
			name:  "one element in list",
			value: "a",
			want:  []string{"a"},
		},
		{
			name:       "one element in list, wrapped",
			value:      "'a'",
			trimQuotes: true,
			want:       []string{"a"},
		},
		{
			name:  "one element in list with brackets",
			value: "[a]",
			want:  []string{"a"},
		},
		{
			name:       "one element in list wrapped, with brackets",
			value:      "['a']",
			trimQuotes: true,
			want:       []string{"a"},
		},
		{
			name:  "multiple elements in list",
			value: "a, b",
			want:  []string{"a", "b"},
		},
		{
			name:       "multiple elements in list wrapped",
			value:      "'a', 'b'",
			trimQuotes: true,
			want:       []string{"a", "b"},
		},
		{
			name:  "multiple elements in list with brackets",
			value: "[a, b]",
			want:  []string{"a", "b"},
		},
		{
			name:       "multiple elements in list wrapped, with brackets",
			value:      "['a',   'b' ]",
			trimQuotes: true,
			want:       []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listValueToSlice(tt.value, tt.trimQuotes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listValueToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
