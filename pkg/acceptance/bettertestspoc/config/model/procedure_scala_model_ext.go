package model

import (
	"encoding/json"
)

func (f *ProcedureScalaModel) MarshalJSON() ([]byte, error) {
	type Alias ProcedureScalaModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}
