package resources_test

import (
	"database/sql"
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestWarehouse(t *testing.T) {
	resources.Warehouse().InternalValidate(provider.Provider().Schema, false)
}

func TestValiateWarehouseName(t *testing.T) {
	a := assert.New(t)

	warns, errs := resources.ValidateWarehouseName("foo", "name")
	a.Len(warns, 0)
	a.Len(errs, 0)
}

func TestWarehouseCreate(t *testing.T) {
	a := assert.New(t)
	w := resources.NewResourceWarehouse()
	a.NotNil(w)

	in := map[string]interface{}{
		"name":    "good_name",
		"comment": "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Warehouse().Schema, in)
	a.NotNil(d)

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {

		mock.ExpectExec("CREATE WAREHOUSE good_name COMMENT='great comment").WillReturnResult(sqlmock.NewResult(1, 1))
		err := w.Create(d, db)
		a.NoError(err)
	})
}

func resource(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	a := assert.New(t)
	d := schema.TestResourceDataRaw(t, resources.Warehouse().Schema, params)
	a.NotNil(d)
	d.SetId(id)
	return d
}

func TestWarehouseRead(t *testing.T) {
	a := assert.New(t)
	w := resources.NewResourceWarehouse()

	d := resource(t, "good_name", map[string]interface{}{"name": "good_name"})

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("SHOW WAREHOUSES LIKE 'good_name'").WillReturnResult(sqlmock.NewResult(1, 1))
		rows := sqlmock.NewRows([]string{"name", "comment"}).AddRow("good_name", "mock comment")
		mock.ExpectQuery(`select "name", "comment" from table\(result_scan\(last_query_id\(\)\)\)`).WillReturnRows(rows)
		err := w.Read(d, db)
		a.NoError(err)
		a.Equal("mock comment", d.Get("comment").(string))
	})
}

func TestWarehouseDelete(t *testing.T) {
	a := assert.New(t)
	w := resources.NewResourceWarehouse()
	a.NotNil(w)

	d := resource(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("DROP WAREHOUSE drop_it").WillReturnResult(sqlmock.NewResult(1, 1))
		err := w.Delete(d, db)
		a.NoError(err)
	})
}

func withMockDb(t *testing.T, f func(*sql.DB, sqlmock.Sqlmock)) {
	a := assert.New(t)
	db, mock, err := sqlmock.New()
	defer db.Close()
	a.NoError(err)

	f(db, mock)
}
