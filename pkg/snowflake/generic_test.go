package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatStringList(t *testing.T) {
	r := require.New(t)

	in := []string{"this", "is", "just", "a", "test"}
	out := formatStringList(in)

	r.Equal("('this', 'is', 'just', 'a', 'test')", out)
}

func TestFormatStringListWithEscape(t *testing.T) {
	r := require.New(t)

	in := []string{"th'is", "is", "just", "a", "test"}
	out := formatStringList(in)

	r.Equal("('th\\'is', 'is', 'just', 'a', 'test')", out)
}
