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
