package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func FromModelPoc(t *testing.T, model ResourceModel) string {
	t.Helper()

	json, err := NewBasicJsonProvider().JsonFromModel(model)
	require.NoError(t, err)
	t.Logf("Generated json:\n%s", json)

	hcl, err := NewHclV1ConfigProvider().HclFromJson(json)
	require.NoError(t, err)
	t.Logf("Generated config:\n%s", hcl)

	return hcl
}
