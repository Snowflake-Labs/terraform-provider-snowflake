// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExternalTableCreate(t *testing.T) {
	r := require.New(t)
	s := NewExternalTableBuilder("test_table", "test_db", "test_schema")
	s.WithColumns([]map[string]string{{"name": "column1", "type": "OBJECT", "as": "expression1"}, {"name": "column2", "type": "VARCHAR", "as": "expression2"}})
	s.WithLocation("location")
	s.WithPattern("pattern")
	s.WithFileFormat("TYPE = CSV FIELD_DELIMITER = '|'")
	r.Equal(`"test_db"."test_schema"."test_table"`, s.QualifiedName())

	r.Equal(`CREATE EXTERNAL TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT AS expression1, "column2" VARCHAR AS expression2) WITH LOCATION = location REFRESH_ON_CREATE = false AUTO_REFRESH = false PATTERN = 'pattern' FILE_FORMAT = ( TYPE = CSV FIELD_DELIMITER = '|' )`, s.Create())

	s.WithComment("Test Comment")
	r.Equal(`CREATE EXTERNAL TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT AS expression1, "column2" VARCHAR AS expression2) WITH LOCATION = location REFRESH_ON_CREATE = false AUTO_REFRESH = false PATTERN = 'pattern' FILE_FORMAT = ( TYPE = CSV FIELD_DELIMITER = '|' ) COMMENT = 'Test Comment'`, s.Create())
}

func TestExternalTableUpdate(t *testing.T) {
	r := require.New(t)
	s := NewExternalTableBuilder("test_table", "test_db", "test_schema")
	s.WithTags([]TagValue{{Name: "tag1", Value: "value1", Schema: "test_schema", Database: "test_db"}})
	expected := `ALTER EXTERNAL TABLE "test_db"."test_schema"."test_table" TAG "test_db"."test_schema"."tag1" = "value1"`
	r.Equal(expected, s.Update())
}

func TestExternalTableDrop(t *testing.T) {
	r := require.New(t)
	s := NewExternalTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`DROP EXTERNAL TABLE "test_db"."test_schema"."test_table"`, s.Drop())
}

func TestExternalTableShow(t *testing.T) {
	r := require.New(t)
	s := NewExternalTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`SHOW EXTERNAL TABLES LIKE 'test_table' IN SCHEMA "test_db"."test_schema"`, s.Show())
}
