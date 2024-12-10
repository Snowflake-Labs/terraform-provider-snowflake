package model

import (
	"encoding/json"
)

func (f *FunctionJavaModel) MarshalJSON() ([]byte, error) {
	type Alias FunctionJavaModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}
