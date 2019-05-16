package resources_test

import (
	"database/sql"
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestWarehouse(t *testing.T) {
	t.Parallel()
	r := require.New(t)
	err := resources.Warehouse().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestWarehouseCreate(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	in := map[string]interface{}{
		"name":    "good_name",
		"comment": "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Warehouse().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE WAREHOUSE "good_name" COMMENT='great comment`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadWarehouse(mock)
		err := resources.CreateWarehouse(d, db)
		a.NoError(err)
	})
}

func expectReadWarehouse(mock sqlmock.Sqlmock) {
	mock.ExpectExec(`USE WAREHOUSE "good_name"`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("SHOW WAREHOUSES LIKE 'good_name'").WillReturnResult(sqlmock.NewResult(1, 1))
	rows := sqlmock.NewRows([]string{"name", "comment", "size"}).AddRow("good_name", "mock comment", "SMALL")
	mock.ExpectQuery(`select "name", "comment", "size" from table\(result_scan\(last_query_id\(\)\)\)`).WillReturnRows(rows)
}

func TestWarehouseRead(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	d := warehouse(t, "good_name", map[string]interface{}{"name": "good_name"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadWarehouse(mock)
		err := resources.ReadWarehouse(d, db)
		a.NoError(err)
		a.Equal("mock comment", d.Get("comment").(string))
	})
}

func TestWarehouseDelete(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	d := warehouse(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP WAREHOUSE "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteWarehouse(d, db)
		a.NoError(err)
	})
}
