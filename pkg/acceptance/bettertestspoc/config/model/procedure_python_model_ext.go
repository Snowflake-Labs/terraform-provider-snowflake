package model

import (
	"encoding/json"
)

func (f *ProcedurePythonModel) MarshalJSON() ([]byte, error) {
	type Alias ProcedurePythonModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}
