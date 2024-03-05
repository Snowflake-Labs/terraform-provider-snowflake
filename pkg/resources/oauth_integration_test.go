package resources_test

import (
	"database/sql"
	"testing"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestOAuthIntegration(t *testing.T) {
	r := require.New(t)
	err := resources.OAuthIntegration().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestOAuthIntegrationCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":         "test_oauth_integration",
		"oauth_client": "TABLEAU_DESKTOP",
	}
	d := schema.TestResourceDataRaw(t, resources.OAuthIntegration().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE SECURITY INTEGRATION "test_oauth_integration" TYPE=OAUTH OAUTH_CLIENT='TABLEAU_DESKTOP' OAUTH_USE_SECONDARY_ROLES='NONE'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadOAuthIntegration(mock)

		err := resources.CreateOAuthIntegration(d, &internalprovider.Context{
			Client: sdk.NewClientFromDB(db),
		})
		r.NoError(err)
	})
}

func TestOAuthIntegrationRead(t *testing.T) {
	r := require.New(t)

	d := oauthIntegration(t, "test_oauth_integration", map[string]interface{}{"name": "test_oauth_integration"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadOAuthIntegration(mock)

		err := resources.ReadOAuthIntegration(d, &internalprovider.Context{
			Client: sdk.NewClientFromDB(db),
		})
		r.NoError(err)
	})
}

func TestOAuthIntegrationDelete(t *testing.T) {
	r := require.New(t)

	d := oauthIntegration(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP SECURITY INTEGRATION "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteOAuthIntegration(d, &internalprovider.Context{
			Client: sdk.NewClientFromDB(db),
		})
		r.NoError(err)
	})
}

func expectReadOAuthIntegration(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "comment", "created_on",
	},
	).AddRow("test_oauth_integration", "OAUTH - TABLEAU_DESKTOP", "SECURITY", true, nil, "now")
	mock.ExpectQuery(`^SHOW SECURITY INTEGRATIONS LIKE 'test_oauth_integration'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).AddRow("OAUTH_ISSUE_REFRESH_TOKENS", "Boolean", "true", "true").
		AddRow("OAUTH_REFRESH_TOKEN_VALIDITY", "Integer", "86400", "7776000").
		AddRow("BLOCKED_ROLES_LIST", "List", "ACCOUNTADMIN,SECURITYADMIN", nil)

	mock.ExpectQuery(`DESCRIBE SECURITY INTEGRATION "test_oauth_integration"$`).WillReturnRows(descRows)
}
