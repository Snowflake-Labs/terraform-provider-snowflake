package model

import (
	"encoding/json"
)

func (f *FunctionScalaModel) MarshalJSON() ([]byte, error) {
	type Alias FunctionScalaModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}
