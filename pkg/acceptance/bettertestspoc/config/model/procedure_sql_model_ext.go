package model

import (
	"encoding/json"
)

func (f *ProcedureSqlModel) MarshalJSON() ([]byte, error) {
	type Alias ProcedureSqlModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}
