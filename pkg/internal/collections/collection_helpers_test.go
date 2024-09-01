package collections

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FindOne(t *testing.T) {
	stringSlice := []string{"1", "22", "333", "334"}

	t.Run("basic find", func(t *testing.T) {
		result, resultErr := FindOne(stringSlice, func(s string) bool { return s == "22" })

		require.Equal(t, "22", *result)
		require.Nil(t, resultErr)
	})

	t.Run("two matching, first returned", func(t *testing.T) {
		result, resultErr := FindOne(stringSlice, func(s string) bool { return strings.HasPrefix(s, "33") })

		require.Equal(t, "333", *result)
		require.Nil(t, resultErr)
	})

	t.Run("no item", func(t *testing.T) {
		result, resultErr := FindOne(stringSlice, func(s string) bool { return s == "4444" })

		require.Nil(t, result)
		require.ErrorIs(t, resultErr, ErrObjectNotFound)
	})
}

func Test_Map(t *testing.T) {
	t.Run("basic mapping", func(t *testing.T) {
		stringSlice := []string{"1", "22", "333"}
		stringLenSlice := Map(stringSlice, func(s string) int { return len(s) })
		require.Equal(t, stringLenSlice, []int{1, 2, 3})
	})

	t.Run("validation: empty slice", func(t *testing.T) {
		stringSlice := make([]string, 0)
		stringLenSlice := Map(stringSlice, func(s string) int { return len(s) })
		require.Equal(t, stringLenSlice, []int{})
	})

	t.Run("validation: nil slice", func(t *testing.T) {
		var stringSlice []string = nil
		stringLenSlice := Map(stringSlice, func(s string) int { return len(s) })
		require.Equal(t, stringLenSlice, []int{})
	})

	t.Run("validation: nil mapping function", func(t *testing.T) {
		require.PanicsWithError(t, "runtime error: invalid memory address or nil pointer dereference", func() {
			stringSlice := []string{"1", "22", "333"}
			_ = Map[string, int](stringSlice, nil)
		})
	})
}
