package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestIntegrationGrant(t *testing.T) {
	r := require.New(t)
	err := resources.IntegrationGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestIntegrationGrantCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"integration_name": "test-integration",
		"privilege":        "USAGE",
		"roles":            []interface{}{"test-role-1", "test-role-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.IntegrationGrant().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT USAGE ON INTEGRATION "test-integration" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON INTEGRATION "test-integration" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadIntegrationGrant(mock)
		err := resources.CreateIntegrationGrant(d, db)
		a.NoError(err)
	})
}

func TestIntegrationGrantRead(t *testing.T) {
	a := assert.New(t)

	d := integrationGrant(t, "test-integration|||IMPORTED PRIVILIGES", map[string]interface{}{
		"integration_name": "test-integration",
		"privilege":        "IMPORTED PRIVILIGES",
		"roles":            []interface{}{"test-role-1", "test-role-2"},
	})

	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadIntegrationGrant(mock)
		err := resources.ReadIntegrationGrant(d, db)
		a.NoError(err)
	})
}

func expectReadIntegrationGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "INTEGRATION", "test-integration", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "INTEGRATION", "test-integration", "ROLE", "test-role-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON INTEGRATION "test-integration"$`).WillReturnRows(rows)
}
