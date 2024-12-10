package model

import (
	"encoding/json"
)

func (f *FunctionPythonModel) MarshalJSON() ([]byte, error) {
	type Alias FunctionPythonModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}
