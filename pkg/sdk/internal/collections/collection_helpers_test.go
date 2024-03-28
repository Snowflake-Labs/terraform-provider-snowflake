package collections

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
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
