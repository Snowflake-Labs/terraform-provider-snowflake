package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
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
		"name":     "good_name",
		"database": "database_name",
		"schema":   "schema_name",
		"comment":  "great comment",
		"column": []interface{}{
			map[string]interface{}{
				"name": "column1",
				"type": "OBJECT",
			},
			map[string]interface{}{
				"name":     "column2",
				"type":     "VARCHAR",
				"nullable": false,
			},
			map[string]interface{}{
				"name":    "column3",
				"type":    "NUMBER(38,0)",
				"comment": "some comment",
			},
			map[string]interface{}{
				"name":           "column4",
				"type":           "VARCHAR",
				"nullable":       true,
				"masking_policy": "TEST_MP",
			},
			map[string]interface{}{
				"name":     "column5",
				"type":     "VARCHAR",
				"nullable": false,
				"default": []interface{}{
					map[string]interface{}{
						"constant": "hello",
					},
				},
			},
			map[string]interface{}{
				"name":           "column6",
				"type":           "VARCHAR",
				"nullable":       true,
				"tag":            []snowflake.TagValue{
					{
						Name:     "columnTag",
						Database: "database_name",
						Schema:   "schema_name",
						Value:    "value",
					},
					{
						Name:     "columnTag2",
						Database: "database_name",
						Schema:   "schema_name",
						Value:    "value2",
					},
				},
			},
		},
		"primary_key": []interface{}{map[string]interface{}{"name": "MY_KEY", "keys": []interface{}{"column1"}}},
	}
	d := table(t, "database_name|schema_name|good_name", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE TABLE "database_name"."schema_name"."good_name" \("column1" OBJECT COMMENT '', "column2" VARCHAR NOT NULL COMMENT '', "column3" NUMBER\(38,0\) COMMENT 'some comment', "column4" VARCHAR WITH MASKING POLICY TEST_MP COMMENT '', "column5" VARCHAR NOT NULL DEFAULT 'hello' COMMENT '', "column6" VARCHAR WITH TAG \("test_db"."test_schema"."columnTag" = "value", "test_db"."test_schema"."columnTag2" = "value2"\) COMMENT '' ,CONSTRAINT "MY_KEY" PRIMARY KEY\("column1"\)\) COMMENT = 'great comment' DATA_RETENTION_TIME_IN_DAYS = 1 CHANGE_TRACKING = false`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectTableRead(mock)
		err := resources.CreateTable(d, db)
		r.NoError(err)
		r.Equal("good_name", d.Get("name").(string))
		columns := d.Get("column").([]interface{})
		r.Equal(6, len(columns))
		col1 := columns[0].(map[string]interface{})
		r.Equal("column1", col1["name"].(string))
		r.Equal("OBJECT", col1["type"].(string))
		r.Equal(true, col1["nullable"].(bool))
		col2 := columns[1].(map[string]interface{})
		r.Equal("column2", col2["name"].(string))
		r.Equal("VARCHAR", col2["type"].(string))
		r.Equal(false, col2["nullable"].(bool))
		col3 := columns[2].(map[string]interface{})
		r.Equal("column3", col3["name"].(string))
		r.Equal("NUMBER(38,0)", col3["type"].(string))
		r.Equal(true, col3["nullable"].(bool))
		r.Equal("some comment", col3["comment"].(string))
		col4 := columns[3].(map[string]interface{})
		r.Equal("column4", col4["name"].(string))
		r.Equal("VARCHAR", col4["type"].(string))
		r.Equal(true, col4["nullable"].(bool))
		r.Equal("TEST_MP", col4["masking_policy"].(string))
		col5 := columns[4].(map[string]interface{})
		r.Equal("column5", col5["name"].(string))
		r.Equal("VARCHAR", col5["type"].(string))
		r.NotNil(col5["default"])
		col5Default := col5["default"].([]interface{})
		r.Equal(1, len(col5Default))
		col5DefaultParams := col5Default[0].(map[string]interface{})
		r.Equal("hello", col5DefaultParams["constant"].(string))
		col6 := columns[5].(map[string]interface{})
		r.Equal("column6", col6["name"].(string))
		r.Equal("VARCHAR", col6["type"].(string))
		r.Equal(true, col6["nullable"].(bool))
		r.Equal([]snowflake.TagValue{
			{
				Name:     "columnTag",
				Database: "database_name",
				Schema:   "schema_name",
				Value:    "value",
			},
			{
				Name:     "columnTag2",
				Database: "database_name",
				Schema:   "schema_name",
				Value:    "value2",
			},
		}, col6["tag"].([]snowflake.TagValue))
	})
}

func expectTableRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"name", "type", "kind", "null?", "default", "primary key", "unique key", "check", "expression", "comment"}).AddRow("good_name", "VARCHAR()", "COLUMN", "Y", "NULL", "NULL", "N", "N", "NULL", "mock comment")
	mock.ExpectQuery(`SHOW TABLES LIKE 'good_name' IN SCHEMA "database_name"."schema_name"`).WillReturnRows(rows)

	describeRows := sqlmock.NewRows([]string{"name", "type", "kind", "null?", "default", "policy name", "comment"}).
		AddRow("column1", "OBJECT", "COLUMN", "Y", nil, nil, nil).
		AddRow("column2", "VARCHAR", "COLUMN", "N", nil, nil, nil).
		AddRow("column3", "NUMBER(38,0)", "COLUMN", "Y", nil, nil, "some comment").
		AddRow("column4", "VARCHAR", "COLUMN", "Y", nil, "TEST_MP", nil).
		AddRow("column5", "VARCHAR", "COLUMN", "N", "'hello'", nil, nil)

	mock.ExpectQuery(`DESC TABLE "database_name"."schema_name"."good_name"`).WillReturnRows(describeRows)
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
		q := snowflake.NewTableBuilder("good_name", "database_name", "schema_name").Show()
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
