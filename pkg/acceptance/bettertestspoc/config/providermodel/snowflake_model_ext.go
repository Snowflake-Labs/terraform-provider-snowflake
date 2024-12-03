package providermodel

import (
	"encoding/json"
	"fmt"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// Based on https://medium.com/picus-security-engineering/custom-json-marshaller-in-go-and-common-pitfalls-c43fa774db05.
func (m *SnowflakeModel) MarshalJSON() ([]byte, error) {
	type AliasModelType SnowflakeModel
	return json.Marshal(&struct {
		*AliasModelType
		Alias string `json:"alias,omitempty"`
	}{
		AliasModelType: (*AliasModelType)(m),
		Alias:          m.Alias(),
	})
}

func (m *SnowflakeModel) WithUserId(userId sdk.AccountObjectIdentifier) *SnowflakeModel {
	return m.WithUser(userId.Name())
}

func (m *SnowflakeModel) WithRoleId(roleId sdk.AccountObjectIdentifier) *SnowflakeModel {
	return m.WithRole(roleId.Name())
}

func (m *SnowflakeModel) WithWarehouseId(warehouseId sdk.AccountObjectIdentifier) *SnowflakeModel {
	return m.WithWarehouse(warehouseId.Name())
}

func (m *SnowflakeModel) WithAuthenticatorType(authenticationType sdk.AuthenticationType) *SnowflakeModel {
	return m.WithAuthenticator(string(authenticationType))
}

func (m *SnowflakeModel) WithPrivateKeyMultiline(privateKey string) *SnowflakeModel {
	return m.WithPrivateKey(fmt.Sprintf(`%[1]s%[2]s%[1]s`, config.SnowflakeProviderConfigPrivateKey, privateKey))
}

func (m *SnowflakeModel) WithClientStoreTemporaryCredentialBool(clientStoreTemporaryCredential bool) *SnowflakeModel {
	m.ClientStoreTemporaryCredential = tfconfig.BoolVariable(clientStoreTemporaryCredential)
	return m
}

func (m *SnowflakeModel) AllFields(tmpConfig *helpers.TmpTomlConfig, tmpUser *helpers.TmpServiceUser) *SnowflakeModel {
	return SnowflakeProvider().
		WithProfile(tmpConfig.Profile).
		WithOrganizationName(tmpUser.AccountId.OrganizationName()).
		WithAccountName(tmpUser.AccountId.AccountName()).
		WithUserId(tmpUser.UserId).
		WithPrivateKeyMultiline(tmpUser.PrivateKey).
		WithWarehouseId(tmpUser.WarehouseId).
		WithProtocol("https").
		WithPort(443).
		WithRoleId(tmpUser.RoleId).
		WithValidateDefaultParameters("true").
		WithClientIp("3.3.3.3").
		WithAuthenticatorType(sdk.AuthenticationTypeJwt).
		WithOktaUrl("https://example-tf.com").
		WithLoginTimeout(101).
		WithRequestTimeout(201).
		WithJwtExpireTimeout(301).
		WithClientTimeout(401).
		WithJwtClientTimeout(501).
		WithExternalBrowserTimeout(601).
		WithInsecureMode(true).
		WithOcspFailOpen("true").
		WithKeepSessionAlive(true).
		WithDisableTelemetry(true).
		WithClientRequestMfaToken("true").
		WithClientStoreTemporaryCredential("true").
		WithDisableQueryContextCache(true).
		WithIncludeRetryReason("true").
		WithMaxRetryCount(3).
		WithDriverTracing("warning").
		WithTmpDirectoryPath("../../").
		WithDisableConsoleLogin("true").
		WithParamsValue(
			tfconfig.ObjectVariable(
				map[string]tfconfig.Variable{
					"foo": tfconfig.StringVariable("piyo"),
				},
			),
		)
}
