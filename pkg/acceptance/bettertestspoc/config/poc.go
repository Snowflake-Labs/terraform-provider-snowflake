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

type ProviderModel interface {
	ProviderName() string
	Alias() string
}

type ProviderModelMeta struct {
	name  string
	alias string
}

func DefaultProviderMeta(name string) *ProviderModelMeta {
	return &ProviderModelMeta{name: name}
}

func ProviderMeta(name string, alias string) *ProviderModelMeta {
	return &ProviderModelMeta{name: name, alias: alias}
}

func (m *ProviderModelMeta) ProviderName() string {
	return m.name
}

func (m *ProviderModelMeta) Alias() string {
	return m.alias
}
