package collections

type Queue[T any] struct {
	data []T
}

func (s *Queue[T]) Head() *T {
	if len(s.data) == 0 {
		return nil
	}
	return &s.data[0]
}

func (s *Queue[T]) Pop() *T {
	elem := s.Head()
	if elem != nil {
		s.data = s.data[1:]
	}
	return elem
}

func (s *Queue[T]) Push(elem T) {
	s.data = append(s.data, elem)
}

func NewQueue[T any]() Queue[T] {
	return Queue[T]{
		data: make([]T, 0),
	}
}
