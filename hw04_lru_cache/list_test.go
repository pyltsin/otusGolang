package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("len change only when len list change", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())

		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
		require.Equal(t, 0, l.Len())

		l.PushFront(1)
		require.Equal(t, 1, l.Len())

		movedItem := l.PushBack(2)
		require.Equal(t, 2, l.Len())

		l.MoveToFront(movedItem)
		require.Equal(t, 2, l.Len())

		l.Remove(movedItem)
		require.Equal(t, 1, l.Len())
	})

	t.Run("add to front", func(t *testing.T) {
		l := NewList()

		l.PushFront(3)
		l.PushFront(2)
		l.PushFront(1)

		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 2, l.Front().Next.Value)
		require.Equal(t, 3, l.Front().Next.Next.Value)
		require.Nil(t, nil, l.Front().Next.Next.Next)

		require.Equal(t, 3, l.Back().Value)
		require.Equal(t, 2, l.Back().Prev.Value)
		require.Equal(t, 1, l.Back().Prev.Prev.Value)
		require.Nil(t, nil, l.Back().Prev.Prev.Prev)
	})

	t.Run("add to back", func(t *testing.T) {
		l := NewList()

		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)

		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 2, l.Front().Next.Value)
		require.Equal(t, 3, l.Front().Next.Next.Value)
		require.Nil(t, nil, l.Front().Next.Next.Next)

		require.Equal(t, 3, l.Back().Value)
		require.Equal(t, 2, l.Back().Prev.Value)
		require.Equal(t, 1, l.Back().Prev.Prev.Value)
		require.Nil(t, nil, l.Back().Prev.Prev.Prev)
	})

	//move to front from middle
	//move to front from front

	t.Run("move to front from back", func(t *testing.T) {
		l := NewList()

		l.PushBack(1)
		l.PushBack(2)
		back := l.PushBack(3)

		l.MoveToFront(back)

		require.Equal(t, 3, l.Front().Value)
		require.Equal(t, 1, l.Front().Next.Value)
		require.Equal(t, 2, l.Front().Next.Next.Value)
		require.Nil(t, nil, l.Front().Next.Next.Next)

		require.Equal(t, 2, l.Back().Value)
		require.Equal(t, 1, l.Back().Prev.Value)
		require.Equal(t, 3, l.Back().Prev.Prev.Value)
		require.Nil(t, nil, l.Back().Prev.Prev.Prev)
	})

	t.Run("move to front from front", func(t *testing.T) {
		l := NewList()

		front := l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)

		l.MoveToFront(front)

		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 2, l.Front().Next.Value)
		require.Equal(t, 3, l.Front().Next.Next.Value)
		require.Nil(t, nil, l.Front().Next.Next.Next)

		require.Equal(t, 3, l.Back().Value)
		require.Equal(t, 2, l.Back().Prev.Value)
		require.Equal(t, 1, l.Back().Prev.Prev.Value)
		require.Nil(t, nil, l.Back().Prev.Prev.Prev)
	})

	t.Run("move to front from middle", func(t *testing.T) {
		l := NewList()

		l.PushBack(1)
		middle := l.PushBack(2)
		l.PushBack(3)

		l.MoveToFront(middle)

		require.Equal(t, 2, l.Front().Value)
		require.Equal(t, 1, l.Front().Next.Value)
		require.Equal(t, 3, l.Front().Next.Next.Value)
		require.Nil(t, nil, l.Front().Next.Next.Next)

		require.Equal(t, 3, l.Back().Value)
		require.Equal(t, 1, l.Back().Prev.Value)
		require.Equal(t, 2, l.Back().Prev.Prev.Value)
		require.Nil(t, nil, l.Back().Prev.Prev.Prev)
	})

	t.Run("remove from front", func(t *testing.T) {
		l := NewList()

		front := l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)

		l.Remove(front)

		require.Equal(t, 2, l.Front().Value)
		require.Equal(t, 3, l.Front().Next.Value)
		require.Nil(t, nil, l.Front().Next.Next)

		require.Equal(t, 3, l.Back().Value)
		require.Equal(t, 2, l.Back().Prev.Value)
		require.Nil(t, nil, l.Back().Prev.Prev)
	})

	t.Run("remove from back", func(t *testing.T) {
		l := NewList()

		l.PushBack(1)
		l.PushBack(2)
		back := l.PushBack(3)

		l.Remove(back)

		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 2, l.Front().Next.Value)
		require.Nil(t, nil, l.Front().Next.Next)

		require.Equal(t, 2, l.Back().Value)
		require.Equal(t, 1, l.Back().Prev.Value)
		require.Nil(t, nil, l.Back().Prev.Prev)
	})

	t.Run("remove from middle", func(t *testing.T) {
		l := NewList()

		l.PushBack(1)
		middle := l.PushBack(2)
		l.PushBack(3)

		l.Remove(middle)

		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 3, l.Front().Next.Value)
		require.Nil(t, nil, l.Front().Next.Next)

		require.Equal(t, 3, l.Back().Value)
		require.Equal(t, 1, l.Back().Prev.Value)
		require.Nil(t, nil, l.Back().Prev.Prev)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
