package config

import (
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
