package collections

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Queue(t *testing.T) {
	t.Run("empty queue initialization", func(t *testing.T) {
		queue := NewQueue[int]()

		require.Nil(t, queue.Head())
		require.Nil(t, queue.Pop())
	})

	t.Run("returns head multiple times", func(t *testing.T) {
		queue := NewQueue[int]()

		queue.Push(1)

		require.Equal(t, 1, *queue.Head())
		require.Equal(t, 1, *queue.Head())
	})

	t.Run("returns empty head after pop", func(t *testing.T) {
		queue := NewQueue[int]()

		queue.Pop()

		require.Nil(t, queue.Head())
	})

	t.Run("multiple operations", func(t *testing.T) {
		queue := NewQueue[int]()

		queue.Push(1)
		require.Equal(t, 1, *queue.Head())

		queue.Push(2)
		require.Equal(t, 1, *queue.Head())

		elem := queue.Pop()
		require.Equal(t, 1, *elem)
		require.Equal(t, 2, *queue.Head())

		elem = queue.Pop()
		require.Equal(t, 2, *elem)
		require.Nil(t, queue.Head())
	})
}
