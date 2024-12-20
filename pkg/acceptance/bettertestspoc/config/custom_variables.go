package config

import "encoding/json"

type nullVariable struct{}

// MarshalJSON returns the JSON encoding of nullVariable.
func (v nullVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}

// NullVariable returns nullVariable which implements Variable.
func NullVariable() nullVariable {
	return nullVariable{}
}

type emptyListVariable struct{}

// MarshalJSON returns the JSON encoding of emptyListVariable.
func (v emptyListVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{})
}

// EmptyListVariable returns emptyListVariable which implements Variable.
func EmptyListVariable() emptyListVariable {
	return emptyListVariable{}
}
