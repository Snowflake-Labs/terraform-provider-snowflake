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

// ResourceFromModel should be used in terraform acceptance tests for Config attribute to get string config from ResourceModel.
// Current implementation is an improved implementation using two steps:
// - .tf.json generation
// - conversion to HCL using hcl v1 lib
// It is still not ideal. HCL v2 should be considered.
func ResourceFromModel(t *testing.T, model ResourceModel) string {
	t.Helper()

	resourceJson, err := DefaultJsonConfigProvider.ResourceJsonFromModel(model)
	require.NoError(t, err)
	t.Logf("Generated json:\n%s", resourceJson)

	hcl, err := DefaultHclConfigProvider.HclFromJson(resourceJson)
	require.NoError(t, err)
	t.Logf("Generated config:\n%s", hcl)

	return hcl
}

// DatasourceFromModel should be used in terraform acceptance tests for Config attribute to get string config from DatasourceModel.
// Current implementation is an improved implementation using two steps:
// - .tf.json generation
// - conversion to HCL using hcl v1 lib
// It is still not ideal. HCL v2 should be considered.
func DatasourceFromModel(t *testing.T, model DatasourceModel) string {
	t.Helper()

	datasourceJson, err := DefaultJsonConfigProvider.DatasourceJsonFromModel(model)
	require.NoError(t, err)
	t.Logf("Generated json:\n%s", datasourceJson)

	hcl, err := DefaultHclConfigProvider.HclFromJson(datasourceJson)
	require.NoError(t, err)
	t.Logf("Generated config:\n%s", hcl)

	return hcl
}

// ProviderFromModel should be used in terraform acceptance tests for Config attribute to get string config from ProviderModel.
// Current implementation is an improved implementation using two steps:
// - .tf.json generation
// - conversion to HCL using hcl v1 lib
// It is still not ideal. HCL v2 should be considered.
func ProviderFromModel(t *testing.T, model ProviderModel) string {
	t.Helper()

	providerJson, err := DefaultJsonConfigProvider.ProviderJsonFromModel(model)
	require.NoError(t, err)
	t.Logf("Generated json:\n%s", providerJson)

	hcl, err := DefaultHclConfigProvider.HclFromJson(providerJson)
	require.NoError(t, err)
	hcl, err = revertEqualSignForMapTypeAttributes(hcl)
	require.NoError(t, err)
	t.Logf("Generated config:\n%s", hcl)

	return hcl
}

// FromModels allows to combine multiple models.
// TODO [SNOW-1501905]: introduce some common interface for all three existing models (ResourceModel, DatasourceModel, and ProviderModel)
func FromModels(t *testing.T, models ...any) string {
	t.Helper()

	var sb strings.Builder
	for i, model := range models {
		switch m := model.(type) {
		case ResourceModel:
			sb.WriteString(ResourceFromModel(t, m))
		case DatasourceModel:
			sb.WriteString(DatasourceFromModel(t, m))
		case ProviderModel:
			sb.WriteString(ProviderFromModel(t, m))
		default:
			t.Fatalf("unknown model: %T", model)
		}
		if i < len(models)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// FromModel should be used in terraform acceptance tests for Config attribute to get string config from ResourceModel.
// Current implementation is really straightforward but it could be improved and tested. It may not handle all cases (like objects, lists, sets) correctly.
// TODO [SNOW-1501905]: use reflection to build config directly from model struct (or some other different way)
// TODO [SNOW-1501905]: add support for config.TestStepConfigFunc (to use as ConfigFile); the naive implementation would be to just create a tmp directory and save file there
// TODO [SNOW-1501905]: add generating MarshalJSON() function
// TODO [SNOW-1501905]: migrate resources to new config generation method (above needed first)
// Use ResourceFromModel, DatasourceFromModel, ProviderFromModel, and FromModels instead.
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

// FromModelsDeprecated allows to combine multiple resource models.
// Use FromModels instead.
func FromModelsDeprecated(t *testing.T, models ...ResourceModel) string {
	t.Helper()
	var sb strings.Builder
	for _, model := range models {
		sb.WriteString(FromModel(t, model) + "\n")
	}
	return sb.String()
}

// ConfigVariablesFromModel constructs config.Variables needed in acceptance tests that are using ConfigVariables in
// combination with ConfigDirectory. It's necessary for cases not supported by FromModel, like lists of objects.
// Use ResourceFromModel, DatasourceFromModel, ProviderFromModel, and FromModels instead.
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
