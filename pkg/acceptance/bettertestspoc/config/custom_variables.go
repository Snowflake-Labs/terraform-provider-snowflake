package config

import "encoding/json"

type emptyListVariable struct{}

// MarshalJSON returns the JSON encoding of emptyListVariable.
func (v emptyListVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{})
}

// EmptyListVariable returns Variable representing an empty list. This is because the current hcl parser handles empty SetVariable incorrectly.
func EmptyListVariable() emptyListVariable {
	return emptyListVariable{}
}

type replacementPlaceholderVariable struct {
	placeholder ReplacementPlaceholder
}

// MarshalJSON returns the JSON encoding of replacementPlaceholderVariable.
func (v replacementPlaceholderVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.placeholder)
}

// ReplacementPlaceholderVariable returns Variable containing one of the ReplacementPlaceholder which is later replaced by HclFormatter.
func ReplacementPlaceholderVariable(placeholder ReplacementPlaceholder) replacementPlaceholderVariable {
	return replacementPlaceholderVariable{placeholder}
}
