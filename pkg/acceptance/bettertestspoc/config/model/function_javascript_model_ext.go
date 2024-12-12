package model

import (
	"encoding/json"
)

func (f *FunctionJavascriptModel) MarshalJSON() ([]byte, error) {
	type Alias FunctionJavascriptModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}
