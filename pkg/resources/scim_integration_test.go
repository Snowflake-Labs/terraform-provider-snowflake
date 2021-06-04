package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestSCIMIntegration(t *testing.T) {
	r := require.New(t)
	err := resources.SCIMIntegration().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestSCIMIntegrationCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":             "test_scim_integration",
		"scim_client":      "AZURE",
		"provisioner_role": "AAD_PROVISIONER",
		"network_policy":   "AAD_NETWORK_POLICY",
	}
	d := schema.TestResourceDataRaw(t, resources.SCIMIntegration().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE SECURITY INTEGRATION "test_scim_integration" TYPE=SCIM NETWORK_POLICY='AAD_NETWORK_POLICY' RUN_AS_ROLE='AAD_PROVISIONER' SCIM_CLIENT='AZURE'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadSCIMIntegration(mock)

		err := resources.CreateSCIMIntegration(d, db)
		r.NoError(err)
	})
}

func TestSCIMIntegrationRead(t *testing.T) {
	r := require.New(t)

	d := scimIntegration(t, "test_scim_integration", map[string]interface{}{"name": "test_scim_integration"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadSCIMIntegration(mock)

		err := resources.ReadSCIMIntegration(d, db)
		r.NoError(err)
	})
}

func TestSCIMIntegrationDelete(t *testing.T) {
	r := require.New(t)

	d := scimIntegration(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP SECURITY INTEGRATION "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteSCIMIntegration(d, db)
		r.NoError(err)
	})
}

func expectReadSCIMIntegration(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"name", "type", "category", "created_on"},
	).AddRow("test_scim_integration", "SCIM - AZURE", "SECURITY", "now")
	mock.ExpectQuery(`^SHOW SECURITY INTEGRATIONS LIKE 'test_scim_integration'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).AddRow("NETWORK_POLICY", "String", "AAD_NETWORK_POLICY", nil).
		AddRow("RUN_AS_ROLE", "String", "AAD_PROVISIONER", nil)

	mock.ExpectQuery(`DESCRIBE SECURITY INTEGRATION "test_scim_integration"$`).WillReturnRows(descRows)
}
