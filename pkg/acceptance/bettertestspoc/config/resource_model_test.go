package config_test

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type Item struct {
	IntField    int
	StringField string
}

// SomeModel is an example model struct similar to the ones being generated for our resources.
// It contains all the interesting types of variables.
type SomeModel struct {
	Comment tfconfig.Variable `json:"comment,omitempty"`
	Name    tfconfig.Variable `json:"name,omitempty"`

	StringList tfconfig.Variable `json:"string_list,omitempty"`
	StringSet  tfconfig.Variable `json:"string_set,omitempty"`
	// contains list of Item
	ObjectList tfconfig.Variable `json:"object_list,omitempty"`

	*config.ResourceModelMeta
}

func Some(
	resourceName string,
	name string,
) *SomeModel {
	// resources enum is closed so using one of the existing ones
	d := &SomeModel{ResourceModelMeta: config.Meta(resourceName, resources.Share)}
	d.WithName(name)
	return d
}

func (m *SomeModel) WithComment(comment string) *SomeModel {
	m.Comment = tfconfig.StringVariable(comment)
	return m
}

func (m *SomeModel) WithName(name string) *SomeModel {
	m.Name = tfconfig.StringVariable(name)
	return m
}

func (m *SomeModel) WithStringList(items ...string) *SomeModel {
	variables := make([]tfconfig.Variable, 0)
	for _, i := range items {
		variables = append(variables, tfconfig.StringVariable(i))
	}
	m.StringList = tfconfig.ListVariable(variables...)
	return m
}

func (m *SomeModel) WithStringSet(items ...string) *SomeModel {
	variables := make([]tfconfig.Variable, 0)
	for _, i := range items {
		variables = append(variables, tfconfig.StringVariable(i))
	}
	m.StringSet = tfconfig.SetVariable(variables...)
	return m
}

func (m *SomeModel) WithObjectList(items ...Item) *SomeModel {
	variables := make([]tfconfig.Variable, 0)
	for _, i := range items {
		variables = append(variables, tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"int_field":    tfconfig.IntegerVariable(i.IntField),
				"string_field": tfconfig.StringVariable(i.StringField),
			},
		))
	}
	m.ObjectList = tfconfig.ListVariable(variables...)
	return m
}

func (m *SomeModel) WithDependsOn(values ...string) *SomeModel {
	m.SetDependsOn(values...)
	return m
}

// Based on https://medium.com/picus-security-engineering/custom-json-marshaller-in-go-and-common-pitfalls-c43fa774db05.
func (m *SomeModel) MarshalJSON() ([]byte, error) {
	type Alias SomeModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(m),
		DependsOn: m.DependsOn(),
	})
}
