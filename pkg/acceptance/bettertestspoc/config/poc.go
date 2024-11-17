package config

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

// SnowflakeProviderModel contains our provider's configuration
// TODO: generate model + builders
type SnowflakeProviderModel struct {
	Authenticator                  sdk.AuthenticationType `json:"authenticator,omitempty"`
	ClientStoreTemporaryCredential string                 `json:"client_store_temporary_credential,omitempty"`
	LoginTimeout                   string                 `json:"login_timeout,omitempty"`
	OktaUrl                        string                 `json:"okta_url,omitempty"`
	Port                           int                    `json:"port,omitempty"`
	Profile                        string                 `json:"profile,omitempty"`
	Protocol                       string                 `json:"protocol,omitempty"`
	Role                           string                 `json:"role,omitempty"`
	ValidateDefaultParameters      string                 `json:"validate_default_parameters,omitempty"`
	Warehouse                      string                 `json:"warehouse,omitempty"`

	*ProviderModelMeta
}

func SnowflakeProvider() SnowflakeProviderModel {
	return SnowflakeProviderModel{ProviderModelMeta: &ProviderModelMeta{name: "snowflake"}}
}

func SnowflakeProviderAlias(alias string) SnowflakeProviderModel {
	return SnowflakeProviderModel{ProviderModelMeta: &ProviderModelMeta{name: "snowflake", alias: alias}}
}

func (m *SnowflakeProviderModel) ProviderName() string {
	return m.name
}

func (m *SnowflakeProviderModel) Alias() string {
	return m.alias
}
