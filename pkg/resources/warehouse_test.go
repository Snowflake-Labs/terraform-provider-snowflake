package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestWarehouse(t *testing.T) {
	r := require.New(t)
	err := resources.Warehouse().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestWarehouseCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":    "good_name",
		"comment": "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Warehouse().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE WAREHOUSE "good_name" COMMENT='great comment`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadWarehouse(mock)
		err := resources.CreateWarehouse(d, db)
		r.NoError(err)
	})
}

func expectReadWarehouse(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"name", "comment", "size"}).AddRow("good_name", "mock comment", "SMALL")
	mock.ExpectQuery("SHOW WAREHOUSES LIKE 'good_name'").WillReturnRows(rows)

	rows = sqlmock.NewRows(
		[]string{"key", "value", "default", "level", "description", "type"},
	).AddRow("MAX_CONCURRENCY_LEVEL", 8, 8, "WAREHOUSE", "", "NUMBER")
	mock.ExpectQuery("SHOW PARAMETERS IN WAREHOUSE good_name").WillReturnRows(rows)
}

func TestWarehouseRead(t *testing.T) {
	r := require.New(t)

	d := warehouse(t, "good_name", map[string]interface{}{"name": "good_name"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadWarehouse(mock)
		err := resources.ReadWarehouse(d, db)
		r.NoError(err)
		r.Equal("mock comment", d.Get("comment").(string))

		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.Warehouse(d.Id()).Show()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err2 := resources.ReadWarehouse(d, db)
		r.Empty(d.State())
		r.Nil(err2)
	})
}

func TestWarehouseDelete(t *testing.T) {
	r := require.New(t)

	d := warehouse(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP WAREHOUSE "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteWarehouse(d, db)
		r.NoError(err)
	})
}
