package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func ResourceFromModelPoc(t *testing.T, model ResourceModel) string {
	t.Helper()

	json, err := DefaultJsonProvider.ResourceJsonFromModel(model)
	require.NoError(t, err)
	t.Logf("Generated json:\n%s", json)

	hcl, err := DefaultHclProvider.HclFromJson(json)
	require.NoError(t, err)
	t.Logf("Generated config:\n%s", hcl)

	return hcl
}

func DatasourceFromModelPoc(t *testing.T, model DatasourceModel) string {
	t.Helper()

	json, err := DefaultJsonProvider.DatasourceJsonFromModel(model)
	require.NoError(t, err)
	t.Logf("Generated json:\n%s", json)

	hcl, err := DefaultHclProvider.HclFromJson(json)
	require.NoError(t, err)
	t.Logf("Generated config:\n%s", hcl)

	return hcl
}

func ProviderFromModelPoc(t *testing.T, model ProviderModel) string {
	t.Helper()

	json, err := DefaultJsonProvider.ProviderJsonFromModel(model)
	require.NoError(t, err)
	t.Logf("Generated json:\n%s", json)

	hcl, err := DefaultHclProvider.HclFromJson(json)
	require.NoError(t, err)
	t.Logf("Generated config:\n%s", hcl)

	return hcl
}

// TODO: have a common interface for all models
func ConfigFromModelsPoc(t *testing.T, models ...any) string {
	t.Helper()

	var sb strings.Builder
	for i, model := range models {
		switch m := model.(type) {
		case ResourceModel:
			sb.WriteString(ResourceFromModelPoc(t, m))
		case DatasourceModel:
			sb.WriteString(DatasourceFromModelPoc(t, m))
		case ProviderModel:
			sb.WriteString(ProviderFromModelPoc(t, m))
		default:
			t.Fatalf("unknown model: %T", model)
		}
		if i < len(models)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
