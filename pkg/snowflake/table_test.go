package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTableCreate(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	cols := []Column{
		{
			name:     "column1",
			_type:    "OBJECT",
			nullable: true,
		},
		{
			name:     "column2",
			_type:    "VARCHAR",
			nullable: true,
			comment:  "only populated when data is available",
		},
		{
			name:     "column3",
			_type:    "NUMBER(38,0)",
			nullable: false,
			_default: NewColumnDefaultWithSequence(`"test_db"."test_schema"."test_seq"`),
		},
		{
			name:     "column4",
			_type:    "VARCHAR",
			nullable: false,
			_default: NewColumnDefaultWithConstant("test default's"),
		},
		{
			name:     "column5",
			_type:    "TIMESTAMP_NTZ",
			nullable: false,
			_default: NewColumnDefaultWithExpression("CURRENT_TIMESTAMP()"),
		},
	}

	s.WithColumns(Columns(cols))
	r.Equal(s.QualifiedName(), `"test_db"."test_schema"."test_table"`)

	r.Equal(`CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '') DATA_RETENTION_TIME_IN_DAYS = 0 CHANGE_TRACKING = false`, s.Create())

	s.WithComment("Test Comment")
	r.Equal(s.Create(), `CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '') COMMENT = 'Test Comment' DATA_RETENTION_TIME_IN_DAYS = 0 CHANGE_TRACKING = false`)

	s.WithClustering([]string{"column1"})
	r.Equal(s.Create(), `CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '') COMMENT = 'Test Comment' CLUSTER BY LINEAR(column1) DATA_RETENTION_TIME_IN_DAYS = 0 CHANGE_TRACKING = false`)

	s.WithPrimaryKey(PrimaryKey{name: "MY_KEY", keys: []string{"column1"}})
	r.Equal(s.Create(), `CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '' ,CONSTRAINT "MY_KEY" PRIMARY KEY("column1")) COMMENT = 'Test Comment' CLUSTER BY LINEAR(column1) DATA_RETENTION_TIME_IN_DAYS = 0 CHANGE_TRACKING = false`)

	s.WithDataRetentionTimeInDays(10)
	r.Equal(s.Create(), `CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '' ,CONSTRAINT "MY_KEY" PRIMARY KEY("column1")) COMMENT = 'Test Comment' CLUSTER BY LINEAR(column1) DATA_RETENTION_TIME_IN_DAYS = 10 CHANGE_TRACKING = false`)

	s.WithChangeTracking(true)
	r.Equal(s.Create(), `CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '' ,CONSTRAINT "MY_KEY" PRIMARY KEY("column1")) COMMENT = 'Test Comment' CLUSTER BY LINEAR(column1) DATA_RETENTION_TIME_IN_DAYS = 10 CHANGE_TRACKING = true`)
}

func TestTableChangeComment(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ChangeComment("new table comment"), `ALTER TABLE "test_db"."test_schema"."test_table" SET COMMENT = 'new table comment'`)
}

func TestTableRemoveComment(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.RemoveComment(), `ALTER TABLE "test_db"."test_schema"."test_table" UNSET COMMENT`)
}

func TestTableAddColumn(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.AddColumn("new_column", "VARIANT", true, nil, ""), `ALTER TABLE "test_db"."test_schema"."test_table" ADD COLUMN "new_column" VARIANT COMMENT ''`)
}

func TestTableAddColumnWithComment(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.AddColumn("new_column", "VARIANT", true, nil, "some comment"), `ALTER TABLE "test_db"."test_schema"."test_table" ADD COLUMN "new_column" VARIANT COMMENT 'some comment'`)
}

func TestTableAddColumnWithDefault(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.AddColumn("new_column", "NUMBER(38,0)", true, NewColumnDefaultWithConstant("1"), ""), `ALTER TABLE "test_db"."test_schema"."test_table" ADD COLUMN "new_column" NUMBER(38,0) DEFAULT 1 COMMENT ''`)
}

func TestTableDropColumn(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.DropColumn("old_column"), `ALTER TABLE "test_db"."test_schema"."test_table" DROP COLUMN "old_column"`)
}

func TestTableChangeColumnType(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ChangeColumnType("old_column", "BIGINT"), `ALTER TABLE "test_db"."test_schema"."test_table" MODIFY COLUMN "old_column" BIGINT`)
}

func TestTableChangeColumnComment(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ChangeColumnComment("old_column", "some comment"), `ALTER TABLE "test_db"."test_schema"."test_table" MODIFY COLUMN "old_column" COMMENT 'some comment'`)
}

func TestTableDropColumnDefault(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.DropColumnDefault("old_column"), `ALTER TABLE "test_db"."test_schema"."test_table" MODIFY COLUMN "old_column" DROP DEFAULT`)
}

func TestTableChangeClusterBy(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ChangeClusterBy("column2, column3"), `ALTER TABLE "test_db"."test_schema"."test_table" CLUSTER BY LINEAR(column2, column3)`)
}

func TestTableChangeDataRetention(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ChangeDataRetention(5), `ALTER TABLE "test_db"."test_schema"."test_table" SET DATA_RETENTION_TIME_IN_DAYS = 5`)
}

func TestTableChangeChangeTracking(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ChangeChangeTracking(true), `ALTER TABLE "test_db"."test_schema"."test_table" SET CHANGE_TRACKING = true`)
}

func TestTableDropClusterBy(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.DropClustering(), `ALTER TABLE "test_db"."test_schema"."test_table" DROP CLUSTERING KEY`)
}

func TestTableDrop(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.Drop(), `DROP TABLE "test_db"."test_schema"."test_table"`)
}

func TestTableShow(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.Show(), `SHOW TABLES LIKE 'test_table' IN SCHEMA "test_db"."test_schema"`)
}

func TestTableShowPrimaryKeys(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ShowPrimaryKeys(), `SHOW PRIMARY KEYS IN TABLE "test_db"."test_schema"."test_table"`)
}

func TestTableDropPrimaryKeys(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.DropPrimaryKey(), `ALTER TABLE "test_db"."test_schema"."test_table" DROP PRIMARY KEY`)
}

func TestTableChangePrimaryKeysWithConstraintName(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ChangePrimaryKey(PrimaryKey{name: "MY_KEY", keys: []string{"column1", "column2"}}), `ALTER TABLE "test_db"."test_schema"."test_table" ADD CONSTRAINT "MY_KEY" PRIMARY KEY("column1", "column2")`)
}

func TestTableChangePrimaryKeysWithoutConstraintName(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ChangePrimaryKey(PrimaryKey{name: "", keys: []string{"column1", "column2"}}), `ALTER TABLE "test_db"."test_schema"."test_table" ADD PRIMARY KEY("column1", "column2")`)
}
