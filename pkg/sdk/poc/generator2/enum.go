package generator2

type Enum[T any] struct {
	Name   string
	Values []EnumValue[T]
}

type EnumValue[T any] struct {
	Name  string
	Value T
}

func NewEnum[T any](name string) *Enum[T] {
	return &Enum[T]{
		Name:   name,
		Values: make([]EnumValue[T], 0),
	}
}

func (e *Enum[T]) WithValue(variableName string, value T) *Enum[T] {
	e.Values = append(e.Values, EnumValue[T]{
		Name:  variableName,
		Value: value,
	})
	return e
}

func (e *Enum[T]) WithValue(variableName string, value T) *Enum[T] {

}
