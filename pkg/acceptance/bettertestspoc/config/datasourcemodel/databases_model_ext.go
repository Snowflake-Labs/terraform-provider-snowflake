package datasourcemodel

import "encoding/json"

// Based on https://medium.com/picus-security-engineering/custom-json-marshaller-in-go-and-common-pitfalls-c43fa774db05.
func (d *DatabasesModel) MarshalJSON() ([]byte, error) {
	type Alias DatabasesModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(d),
		DependsOn: d.DependsOn(),
	})
}

func (d *DatabasesModel) WithDependsOn(values ...string) *DatabasesModel {
	d.SetDependsOn(values...)
	return d
}
