package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestExternalOauthIntegration(t *testing.T) {
	r := require.New(t)
	err := resources.ExternalOauthIntegration().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestExternalOauthIntegrationCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                             "test_external_oauth_integration",
		"type":                             "AZURE",
		"enabled":                          true,
		"issuer":                           "https://sts.windows.net/00000000-0000-0000-0000-000000000000",
		"snowflake_user_mapping_attribute": "LOGIN_NAME",
		"token_user_mapping_claims":        []interface{}{"upn"},
	}
	d := schema.TestResourceDataRaw(t, resources.ExternalOauthIntegration().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE SECURITY INTEGRATION "test_external_oauth_integration" TYPE=EXTERNAL_OAUTH EXTERNAL_OAUTH_ANY_ROLE_MODE='DISABLE' EXTERNAL_OAUTH_ISSUER='https://sts.windows.net/00000000-0000-0000-0000-000000000000' EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE='LOGIN_NAME' EXTERNAL_OAUTH_TYPE='AZURE' EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM=\('upn'\) ENABLED=true$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadExternalOauthIntegration(mock)

		err := resources.CreateExternalOauthIntegration(d, db)
		r.NoError(err)
	})
}

func TestExternalOauthIntegrationRead(t *testing.T) {
	r := require.New(t)

	d := externalOauthIntegration(t, "test_external_oauth_integration", map[string]interface{}{"name": "test_external_oauth_integration"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadExternalOauthIntegration(mock)

		err := resources.ReadExternalOauthIntegration(d, db)
		r.NoError(err)
	})
}

func TestExternalOauthIntegrationDelete(t *testing.T) {
	r := require.New(t)

	d := externalOauthIntegration(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP SECURITY INTEGRATION "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteExternalOauthIntegration(d, db)
		r.NoError(err)
	})
}

func expectReadExternalOauthIntegration(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "comment", "created_on"},
	).AddRow("test_external_oauth_integration", "EXTERNAL_OAUTH - AZURE", "SECURITY", true, nil, "now")
	mock.ExpectQuery(`^SHOW SECURITY INTEGRATIONS LIKE 'test_external_oauth_integration'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).AddRow("EXTERNAL_OAUTH_ISSUER", "String", "https://sts.windows.net/00000000-0000-0000-0000-000000000000", nil).
		AddRow("EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM", "List", "['upn']", nil).
		AddRow("EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE", "String", "upn", nil).
		AddRow("EXTERNAL_OAUTH_RSA_PUBLIC_KEY", "String", "", nil).
		AddRow("EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2", "String", "", nil).
		AddRow("EXTERNAL_OAUTH_BLOCKED_ROLES_LIST", "List", "ACCOUNTADMIN,SECURITYADMIN", nil).
		AddRow("EXTERNAL_OAUTH_JWS_KEYS_URL", "List", "", nil).
		AddRow("EXTERNAL_OAUTH_ALLOWED_ROLES_LIST", "List", "", nil).
		AddRow("EXTERNAL_OAUTH_AUDIENCE_LIST", "List", "", nil).
		AddRow("EXTERNAL_OAUTH_ANY_ROLE_MODE", "String", "", nil)

	mock.ExpectQuery(`DESCRIBE SECURITY INTEGRATION "test_external_oauth_integration"$`).WillReturnRows(descRows)
}
