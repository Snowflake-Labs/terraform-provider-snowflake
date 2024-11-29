package datasourcemodel

import (
	"encoding/json"
)

// Based on https://medium.com/picus-security-engineering/custom-json-marshaller-in-go-and-common-pitfalls-c43fa774db05.
func (d *DatabaseModel) MarshalJSON() ([]byte, error) {
	type Alias DatabaseModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(d),
		DependsOn: d.DependsOn(),
	})
}
