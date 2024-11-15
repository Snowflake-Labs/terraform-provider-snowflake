package config

import (
	"encoding/json"
	"strings"
	"testing"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/stretchr/testify/require"
)

type Item struct {
	IntField    int
	StringField string
}

type SomeModel struct {
	Comment tfconfig.Variable `json:"comment,omitempty"`
	Name    tfconfig.Variable `json:"name,omitempty"`

	StringList tfconfig.Variable `json:"string_list,omitempty"`
	StringSet  tfconfig.Variable `json:"string_set,omitempty"`
	ObjectList tfconfig.Variable `json:"object_list,omitempty"`

	*ResourceModelMeta
}

func Some(
	resourceName string,
	name string,
) *SomeModel {
	d := &SomeModel{ResourceModelMeta: Meta(resourceName, resources.Share)}
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

// https://medium.com/picus-security-engineering/custom-json-marshaller-in-go-and-common-pitfalls-c43fa774db05
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

func Test_FromModelPoc(t *testing.T) {

	t.Run("test basic", func(t *testing.T) {
		someModel := Some("test", "Some Name")
		expectedOutput := `
"resource" "snowflake_share" "test" {
  "name" = "Some Name"
}
`
		result := FromModelPoc(t, someModel)

		require.Equal(t, strings.TrimPrefix(expectedOutput, "\n"), result)
	})

	t.Run("test full", func(t *testing.T) {
		someModel := Some("test", "Some Name").
			WithComment("Some Comment").
			WithStringList("a", "b", "a").
			WithStringSet("a", "b", "c").
			WithObjectList(
				Item{IntField: 1, StringField: "first item"},
				Item{IntField: 2, StringField: "second item"},
			).WithDependsOn("some_other_resource.some_name", "other_resource.some_other_name", "third_resource.third_name")
		expectedOutput := `
"resource" "snowflake_share" "test" {
  "comment" = "Some Comment"
  "name" = "Some Name"
  "string_list" = ["a", "b", "a"]
  "string_set" = ["a", "b", "c"]
  "object_list" = {
    "int_field" = 1
    "string_field" = "first item"
  }
  "object_list" = {
    "int_field" = 2
    "string_field" = "second item"
  }
  "depends_on" = [some_other_resource.some_name, other_resource.some_other_name, third_resource.third_name]
}
`

		result := FromModelPoc(t, someModel)

		require.Equal(t, strings.TrimPrefix(expectedOutput, "\n"), result)
	})
}
