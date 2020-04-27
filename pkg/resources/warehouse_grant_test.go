package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestWarehouseGrant(t *testing.T) {
	r := require.New(t)
	err := resources.WarehouseGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestWarehouseGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"warehouse_name": "test-warehouse",
		"privilege":      "USAGE",
		"roles":          []interface{}{"test-role-1", "test-role-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.WarehouseGrant().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT USAGE ON WAREHOUSE "test-warehouse" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON WAREHOUSE "test-warehouse" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadWarehouseGrant(mock)
		err := resources.CreateWarehouseGrant(d, db)
		r.NoError(err)
	})
}

func expectReadWarehouseGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "WAREHOUSE", "test-warehouse", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "WAREHOUSE", "test-warehouse", "ROLE", "test-role-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON WAREHOUSE "test-warehouse"$`).WillReturnRows(rows)
}
