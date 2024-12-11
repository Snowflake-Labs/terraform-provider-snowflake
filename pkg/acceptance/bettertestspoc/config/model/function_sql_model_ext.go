package model

import (
	"encoding/json"
)

func (f *FunctionSqlModel) MarshalJSON() ([]byte, error) {
	type Alias FunctionSqlModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}
