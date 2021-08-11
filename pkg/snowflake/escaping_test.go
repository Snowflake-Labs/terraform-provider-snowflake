package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestEscapeString(t *testing.T) {
	r := require.New(t)

	r.Equal(`\'`, snowflake.EscapeString(`'`))
	r.Equal(`\\\'`, snowflake.EscapeString(`\'`))
}

func TestEscapeSnowflakeString(t *testing.T) {
	r := require.New(t)
	r.Equal(`'table''s quoted'`, snowflake.EscapeSnowflakeString(`table's quoted`))
}

func TestUnescapeSnowflakeString(t *testing.T) {
	r := require.New(t)
	r.Equal(`table's quoted`, snowflake.UnescapeSnowflakeString(`'table''s quoted'`))
}

func TestAddressEscape(t *testing.T) {
	testCases := []struct {
		id       string
		name     []string
		expected string
	}{
		{
			id:       "single no escape",
			name:     []string{"HELLO"},
			expected: "HELLO",
		},
		{
			id:       "multiple no escape",
			name:     []string{"HELLO", "WORLD"},
			expected: "HELLO.WORLD",
		},
		{
			id:       "single escape",
			name:     []string{"hello"},
			expected: `"hello"`,
		},
		{
			id:       "multiple escape",
			name:     []string{"hello", "world"},
			expected: `"hello"."world"`,
		},
		{
			id:       "mixed escape",
			name:     []string{"hello", "world", "NOTHERE"},
			expected: `"hello"."world".NOTHERE`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			r := require.New(t)
			r.Equal(testCase.expected, snowflake.AddressEscape(testCase.name...))
		})
	}
}
