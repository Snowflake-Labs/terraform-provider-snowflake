package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestWarehouseGrant(t *testing.T) {
	r := require.New(t)
	err := resources.WarehouseGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestWarehouseGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"warehouse_name": "test-warehouse",
		"privilege":      "USAGE",
		"roles":          []interface{}{"test-role-1", "test-role-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.WarehouseGrant().Resource.Schema, in)
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

func TestParseWarehouseGrantID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseWarehouseGrantID("test-db|USAGE|false|role1,role2")
	r.NoError(err)
	r.Equal("test-db", grantID.ObjectName)
	r.Equal("USAGE", grantID.Privilege)
	r.Equal(false, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}

func TestParseWarehouseGrantOldID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseWarehouseGrantID("test-db|||USAGE|role1,role2|false")
	r.NoError(err)
	r.Equal("test-db", grantID.ObjectName)
	r.Equal("USAGE", grantID.Privilege)
	r.Equal(false, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}

func TestParseWarehouseGrantReallyOldID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseWarehouseGrantID("test-db|||USAGE|false")
	r.NoError(err)
	r.Equal("test-db", grantID.ObjectName)
	r.Equal("USAGE", grantID.Privilege)
	r.Equal(false, grantID.WithGrantOption)
	r.Equal(0, len(grantID.Roles))
}
