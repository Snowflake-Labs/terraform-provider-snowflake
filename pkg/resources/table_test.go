package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
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
		"name":        "good_name",
		"database":    "database_name",
		"schema":      "schema_name",
		"comment":     "great comment",
		"column":      []interface{}{map[string]interface{}{"name": "column1", "type": "OBJECT"}, map[string]interface{}{"name": "column2", "type": "VARCHAR", "nullable": false}},
		"primary_key": []interface{}{map[string]interface{}{"name": "MY_KEY", "keys": []interface{}{"column1"}}},
	}
	d := table(t, "database_name|schema_name|good_name", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE TABLE "database_name"."schema_name"."good_name" \("column1" OBJECT, "column2" VARCHAR NOT NULL ,CONSTRAINT "MY_KEY" PRIMARY KEY\("column1"\)\) COMMENT = 'great comment'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectTableRead(mock)
		err := resources.CreateTable(d, db)
		r.NoError(err)
		r.Equal("good_name", d.Get("name").(string))
	})
}

func expectTableRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"name", "type", "kind", "null?", "default", "primary key", "unique key", "check", "expression", "comment"}).AddRow("good_name", "VARCHAR()", "COLUMN", "Y", "NULL", "NULL", "N", "N", "NULL", "mock comment")
	mock.ExpectQuery(`SHOW TABLES LIKE 'good_name' IN SCHEMA "database_name"."schema_name"`).WillReturnRows(rows)

	describeRows := sqlmock.NewRows([]string{"name", "type", "kind", "null?"}).
		AddRow("column1", "OBJECT", "COLUMN", "Y").
		AddRow("column2", "VARCHAR", "COLUMN", "N")

	mock.ExpectQuery(`DESC TABLE "database_name"."schema_name"."good_name"`).WillReturnRows(describeRows)

	pkRows := sqlmock.NewRows([]string{"column_name", "key_sequence", "constraint_name"}).AddRow("column1", "1", "MY_PK")

	mock.ExpectQuery(`SHOW PRIMARY KEYS IN TABLE "database_name"."schema_name"."good_name"`).WillReturnRows(pkRows)

}

func TestTableRead(t *testing.T) {
	r := require.New(t)

	d := table(t, "database_name|schema_name|good_name", map[string]interface{}{"name": "good_name", "comment": "mock comment"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectTableRead(mock)

		err := resources.ReadTable(d, db)
		r.NoError(err)
		r.Equal("good_name", d.Get("name").(string))
		r.Equal("mock comment", d.Get("comment").(string))

		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.Table("good_name", "database_name", "schema_name").Show()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err2 := resources.ReadTable(d, db)
		r.Empty(d.State())
		r.Nil(err2)
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
