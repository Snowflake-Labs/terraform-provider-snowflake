package model

import (
	"encoding/json"
)

func (f *ProcedureJavascriptModel) MarshalJSON() ([]byte, error) {
	type Alias ProcedureJavascriptModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}
