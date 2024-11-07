package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1501905]: add possibility to have reference to another object (e.g. WithResourceMonitorReference); new config.Variable impl?
// TODO [SNOW-1501905]: generate With/SetDependsOn for the resources to preserve builder pattern
// TODO [SNOW-1501905]: add a convenience method to use multiple configs from multiple models

// ResourceModel is the base interface all of our config models will implement.
// To allow easy implementation, ResourceModelMeta can be embedded inside the struct (and the struct will automatically implement it).
type ResourceModel interface {
	Resource() resources.Resource
	ResourceName() string
	SetResourceName(name string)
	ResourceReference() string
	DependsOn() []string
	SetDependsOn(values ...string)
}

type ResourceModelMeta struct {
	name      string
	resource  resources.Resource
	dependsOn []string
}

func (m *ResourceModelMeta) Resource() resources.Resource {
	return m.resource
}

func (m *ResourceModelMeta) ResourceName() string {
	return m.name
}

func (m *ResourceModelMeta) SetResourceName(name string) {
	m.name = name
}

func (m *ResourceModelMeta) ResourceReference() string {
	return fmt.Sprintf(`%s.%s`, m.resource, m.name)
}

func (m *ResourceModelMeta) DependsOn() []string {
	return m.dependsOn
}

func (m *ResourceModelMeta) SetDependsOn(values ...string) {
	m.dependsOn = values
}

// DefaultResourceName is exported to allow assertions against the resources using the default name.
const DefaultResourceName = "test"

func DefaultMeta(resource resources.Resource) *ResourceModelMeta {
	return &ResourceModelMeta{name: DefaultResourceName, resource: resource}
}

func Meta(resourceName string, resource resources.Resource) *ResourceModelMeta {
	return &ResourceModelMeta{name: resourceName, resource: resource}
}

// FromModel should be used in terraform acceptance tests for Config attribute to get string config from ResourceModel.
// Current implementation is really straightforward but it could be improved and tested. It may not handle all cases (like objects, lists, sets) correctly.
// TODO [SNOW-1501905]: use reflection to build config directly from model struct (or some other different way)
// TODO [SNOW-1501905]: add support for config.TestStepConfigFunc (to use as ConfigFile); the naive implementation would be to just create a tmp directory and save file there
func FromModel(t *testing.T, model ResourceModel) string {
	t.Helper()

	b, err := json.Marshal(model)
	require.NoError(t, err)

	var objMap map[string]json.RawMessage
	err = json.Unmarshal(b, &objMap)
	require.NoError(t, err)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`resource "%s" "%s" {`, model.Resource(), model.ResourceName()))
	sb.WriteRune('\n')
	for k, v := range objMap {
		sb.WriteString(fmt.Sprintf("\t%s = %s\n", k, v))
	}
	if len(model.DependsOn()) > 0 {
		sb.WriteString(fmt.Sprintf("\tdepends_on = [%s]\n", strings.Join(model.DependsOn(), ", ")))
	}
	sb.WriteString(`}`)
	sb.WriteRune('\n')
	s := sb.String()
	t.Logf("Generated config:\n%s", s)
	return s
}

func FromModels(t *testing.T, models ...ResourceModel) string {
	t.Helper()
	var sb strings.Builder
	for _, model := range models {
		sb.WriteString(FromModel(t, model) + "\n")
	}
	return sb.String()
}

// ConfigVariablesFromModel constructs config.Variables needed in acceptance tests that are using ConfigVariables in
// combination with ConfigDirectory. It's necessary for cases not supported by FromModel, like lists of objects.
func ConfigVariablesFromModel(t *testing.T, model ResourceModel) tfconfig.Variables {
	t.Helper()
	variables := make(tfconfig.Variables)
	rType := reflect.TypeOf(model).Elem()
	rValue := reflect.ValueOf(model).Elem()
	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		if jsonTag, ok := field.Tag.Lookup("json"); ok {
			name := strings.Split(jsonTag, ",")[0]
			if fieldValue, ok := rValue.Field(i).Interface().(tfconfig.Variable); ok {
				variables[name] = fieldValue
			}
		}
	}
	return variables
}

// ConfigVariablesFromModels can be used to create a list of objects that are referring to the same resource model.
// It's useful when there's a need to create associations between objects of the same type in Snowflake.
func ConfigVariablesFromModels(t *testing.T, variableName string, models ...ResourceModel) tfconfig.Variables {
	t.Helper()
	allVariables := make([]tfconfig.Variable, 0)
	for _, model := range models {
		rType := reflect.TypeOf(model).Elem()
		rValue := reflect.ValueOf(model).Elem()
		variables := make(tfconfig.Variables)
		for i := 0; i < rType.NumField(); i++ {
			field := rType.Field(i)
			if jsonTag, ok := field.Tag.Lookup("json"); ok {
				name := strings.Split(jsonTag, ",")[0]
				if fieldValue, ok := rValue.Field(i).Interface().(tfconfig.Variable); ok {
					variables[name] = fieldValue
				}
			}
		}
		allVariables = append(allVariables, tfconfig.ObjectVariable(variables))
	}
	return tfconfig.Variables{
		variableName: tfconfig.ListVariable(allVariables...),
	}
}

type nullVariable struct{}

// MarshalJSON returns the JSON encoding of nullVariable.
func (v nullVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}

// NullVariable returns nullVariable which implements Variable.
func NullVariable() nullVariable {
	return nullVariable{}
}
