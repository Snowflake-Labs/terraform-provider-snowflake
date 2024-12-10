package model

import (
	"encoding/json"
)

func (f *ProcedureJavaModel) MarshalJSON() ([]byte, error) {
	type Alias ProcedureJavaModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}
