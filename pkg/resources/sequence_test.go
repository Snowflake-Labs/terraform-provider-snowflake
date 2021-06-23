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

func TestSequence(t *testing.T) {
	r := require.New(t)
	err := resources.Sequence().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestSequenceCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "good_name",
		"schema":   "schema",
		"database": "database",
		"comment":  "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Sequence().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE SEQUENCE "database"."schema"."good_name" COMMENT = 'great comment`).WillReturnResult(sqlmock.NewResult(1, 1))

		rows := sqlmock.NewRows([]string{
			"name",
			"database_name",
			"schema_name",
			"next_value",
			"interval",
			"created_on",
			"owner",
			"comment",
		}).AddRow(
			"good_name",
			"database",
			"schema",
			"25",
			"1",
			"created_on",
			"owner",
			"mock comment",
		)
		mock.ExpectQuery(`SHOW SEQUENCES LIKE 'good_name' IN SCHEMA "database"."schema"`).WillReturnRows(rows)
		err := resources.CreateSequence(d, db)
		r.NoError(err)
		r.Equal("database|schema|good_name", d.Id())
	})
}

func TestSequenceRead(t *testing.T) {
	r := require.New(t)
	in := map[string]interface{}{
		"name":     "good_name",
		"schema":   "schema",
		"database": "database",
	}

	d := sequence(t, "good_name", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{
			"name",
			"database_name",
			"schema_name",
			"next_value",
			"interval",
			"created_on",
			"owner",
			"comment",
		}).AddRow(
			"good_name",
			"database",
			"schema",
			"5",
			"25",
			"created_on",
			"owner",
			"mock comment",
		)
		mock.ExpectQuery(`SHOW SEQUENCES LIKE 'good_name' IN SCHEMA "database"."schema"`).WillReturnRows(rows)
		err := resources.ReadSequence(d, db)
		r.NoError(err)
		r.Equal("good_name", d.Get("name").(string))
		r.Equal("schema", d.Get("schema").(string))
		r.Equal("database", d.Get("database").(string))
		r.Equal("mock comment", d.Get("comment").(string))
		r.Equal(25, d.Get("increment").(int))
		r.Equal(5, d.Get("next_value").(int))
		r.Equal("database|schema|good_name", d.Id())
	})
}

func TestSequenceDelete(t *testing.T) {
	r := require.New(t)
	in := map[string]interface{}{
		"name":     "drop_it",
		"schema":   "schema",
		"database": "database",
	}

	d := sequence(t, "drop_it", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP SEQUENCE "database"."schema"."drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteSequence(d, db)
		r.NoError(err)
		r.Equal("", d.Id())
	})
}
