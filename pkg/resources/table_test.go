package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	r := require.New(t)
	err := resources.Table().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestTableCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":    "good_name",
		"comment": "great comment",
		"columns": map[string]interface{}{"column1": "VARCHAR"},
	}
	d := schema.TestResourceDataRaw(t, resources.Table().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE TABLE "good_name" ("column1" VARCHAR) COMMENT = 'great comment'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectTableRead(mock)
		err := resources.CreateTable(d, db)
		r.NoError(err)
	})
}

func expectTableRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"name", "type", "kind", "null?", "primary key", "unique key", "check", "expression", "comment"}).AddRow("good_name", "VARCHAR()", "COLUMN", "Y", "NULL", "N", "N", "NULL", "mock comment")
	mock.ExpectQuery(`SHOW TABLES LIKE 'good_name' IN DATABASE "database_name"`).WillReturnRows(rows)
}

func TestTableRead(t *testing.T) {
	r := require.New(t)

	d := table(t, "database_name|schema_name|good_name", map[string]interface{}{"name": "good_name"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectTableRead(mock)

		err := resources.ReadTable(d, db)
		r.NoError(err)
		r.Equal("good_name", d.Get("name").(string))
		r.Equal("mock comment", d.Get("comment").(string))
	})
}

func TestTableDelete(t *testing.T) {
	r := require.New(t)

	d := table(t, "database_name|schema_name|drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP TABLE "database_name"."schema_name"."drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteTable(d, db)
		r.NoError(err)
	})
}
