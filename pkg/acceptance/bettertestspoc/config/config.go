package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/stretchr/testify/require"
)

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
