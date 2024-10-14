package collections

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func Test_FindFirst(t *testing.T) {
	stringSlice := []string{"1", "22", "333", "334"}

	t.Run("basic find", func(t *testing.T) {
		result, resultErr := FindFirst(stringSlice, func(s string) bool { return s == "22" })

		require.Equal(t, "22", *result)
		require.Nil(t, resultErr)
	})

	t.Run("two matching, first returned", func(t *testing.T) {
		result, resultErr := FindFirst(stringSlice, func(s string) bool { return strings.HasPrefix(s, "33") })

		require.Equal(t, "333", *result)
		require.Nil(t, resultErr)
	})

	t.Run("no item", func(t *testing.T) {
		result, resultErr := FindFirst(stringSlice, func(s string) bool { return s == "4444" })

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

func Test_MapErr(t *testing.T) {
	t.Run("basic mapping", func(t *testing.T) {
		stringSlice := []string{"1", "22", "333"}
		stringLenSlice, err := MapErr(stringSlice, func(s string) (int, error) { return len(s), nil })
		assert.NoError(t, err)
		assert.Equal(t, stringLenSlice, []int{1, 2, 3})
	})

	t.Run("basic mapping - multiple errors", func(t *testing.T) {
		stringSlice := []string{"1", "22", "333"}
		stringLenSlice, err := MapErr(stringSlice, func(s string) (int, error) {
			if s == "1" {
				return -1, fmt.Errorf("error: 1")
			}
			if s == "22" {
				return -1, fmt.Errorf("error: 22")
			}
			return len(s), nil
		})
		assert.Equal(t, stringLenSlice, []int{-1, -1, 3})
		assert.ErrorContains(t, err, errors.Join(fmt.Errorf("error: 1"), fmt.Errorf("error: 22")).Error())
	})

	t.Run("validation: empty slice", func(t *testing.T) {
		stringSlice := make([]string, 0)
		stringLenSlice, err := MapErr(stringSlice, func(s string) (int, error) { return len(s), nil })
		assert.NoError(t, err)
		assert.Equal(t, stringLenSlice, []int{})
	})

	t.Run("validation: nil slice", func(t *testing.T) {
		var stringSlice []string = nil
		stringLenSlice, err := MapErr(stringSlice, func(s string) (int, error) { return len(s), nil })
		assert.NoError(t, err)
		assert.Equal(t, stringLenSlice, []int{})
	})

	t.Run("validation: nil mapping function", func(t *testing.T) {
		assert.PanicsWithError(t, "runtime error: invalid memory address or nil pointer dereference", func() {
			stringSlice := []string{"1", "22", "333"}
			_, _ = MapErr[string, int](stringSlice, nil)
		})
	})
}
