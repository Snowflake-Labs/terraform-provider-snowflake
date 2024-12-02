package datasourcemodel

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

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

func (d *DatabasesModel) WithLimit(rows int) *DatabasesModel {
	return d.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}
