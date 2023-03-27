package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTableCreate(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
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
		{
			name:          "column6",
			_type:         "VARCHAR",
			nullable:      true,
			maskingPolicy: "TEST_MP",
		},
		{
			name:          "column7",
			_type:         "VARCHAR",
			nullable:      true,
			tags: []TagValue{
					{
						Name:     "columnTag",
						Database: "test_db",
						Schema:   "test_schema",
						Value:    "value",
					},
					{
						Name:     "columnTag2",
						Database: "test_db",
						Schema:   "test_schema",
						Value:    "value2",
					},
		         },
		},
	}

	s.WithColumns(Columns(cols))

	tags := []TagValue{
		{
			Name:     "tag",
			Database: "test_db",
			Schema:   "test_schema",
			Value:    "value",
		},
		{
			Name:     "tag2",
			Database: "test_db",
			Schema:   "test_schema",
			Value:    "value2",
		},
	}

	r.Equal(`"test_db"."test_schema"."test_table"`, s.QualifiedName())

	r.Equal(`CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '', "column6" VARCHAR WITH MASKING POLICY TEST_MP COMMENT '', "column7" VARCHAR WITH TAG ("test_db"."test_schema"."columnTag" = "value", "test_db"."test_schema"."columnTag2" = "value2") COMMENT '') DATA_RETENTION_TIME_IN_DAYS = 0 CHANGE_TRACKING = false`, s.Create())

	s.WithComment("Test Comment")
	r.Equal(`CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '', "column6" VARCHAR WITH MASKING POLICY TEST_MP COMMENT '', "column7" VARCHAR WITH TAG ("test_db"."test_schema"."columnTag" = "value", "test_db"."test_schema"."columnTag2" = "value2") COMMENT '') COMMENT = 'Test Comment' DATA_RETENTION_TIME_IN_DAYS = 0 CHANGE_TRACKING = false`, s.Create())

	s.WithClustering([]string{"column1"})
	r.Equal(`CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '', "column6" VARCHAR WITH MASKING POLICY TEST_MP COMMENT '', "column7" VARCHAR WITH TAG ("test_db"."test_schema"."columnTag" = "value", "test_db"."test_schema"."columnTag2" = "value2") COMMENT '') COMMENT = 'Test Comment' CLUSTER BY LINEAR(column1) DATA_RETENTION_TIME_IN_DAYS = 0 CHANGE_TRACKING = false`, s.Create())

	s.WithPrimaryKey(PrimaryKey{name: "MY_KEY", keys: []string{"column1"}})
	r.Equal(`CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '', "column6" VARCHAR WITH MASKING POLICY TEST_MP COMMENT '', "column7" VARCHAR WITH TAG ("test_db"."test_schema"."columnTag" = "value", "test_db"."test_schema"."columnTag2" = "value2") COMMENT '' ,CONSTRAINT "MY_KEY" PRIMARY KEY("column1")) COMMENT = 'Test Comment' CLUSTER BY LINEAR(column1) DATA_RETENTION_TIME_IN_DAYS = 0 CHANGE_TRACKING = false`, s.Create())

	s.WithDataRetentionTimeInDays(10)
	r.Equal(`CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '', "column6" VARCHAR WITH MASKING POLICY TEST_MP COMMENT '', "column7" VARCHAR WITH TAG ("test_db"."test_schema"."columnTag" = "value", "test_db"."test_schema"."columnTag2" = "value2") COMMENT '' ,CONSTRAINT "MY_KEY" PRIMARY KEY("column1")) COMMENT = 'Test Comment' CLUSTER BY LINEAR(column1) DATA_RETENTION_TIME_IN_DAYS = 10 CHANGE_TRACKING = false`, s.Create())

	s.WithChangeTracking(true)
	r.Equal(`CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '', "column6" VARCHAR WITH MASKING POLICY TEST_MP COMMENT '', "column7" VARCHAR WITH TAG ("test_db"."test_schema"."columnTag" = "value", "test_db"."test_schema"."columnTag2" = "value2") COMMENT '' ,CONSTRAINT "MY_KEY" PRIMARY KEY("column1")) COMMENT = 'Test Comment' CLUSTER BY LINEAR(column1) DATA_RETENTION_TIME_IN_DAYS = 10 CHANGE_TRACKING = true`, s.Create())

	s.WithTags(tags)
	r.Equal(`CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL DEFAULT "test_db"."test_schema"."test_seq".NEXTVAL COMMENT '', "column4" VARCHAR NOT NULL DEFAULT 'test default''s' COMMENT '', "column5" TIMESTAMP_NTZ NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '', "column6" VARCHAR WITH MASKING POLICY TEST_MP COMMENT '', "column7" VARCHAR WITH TAG ("test_db"."test_schema"."columnTag" = "value", "test_db"."test_schema"."columnTag2" = "value2") COMMENT '' ,CONSTRAINT "MY_KEY" PRIMARY KEY("column1")) COMMENT = 'Test Comment' CLUSTER BY LINEAR(column1) DATA_RETENTION_TIME_IN_DAYS = 10 CHANGE_TRACKING = true WITH TAG ("test_db"."test_schema"."tag" = "value", "test_db"."test_schema"."tag2" = "value2")`, s.Create())
}

func TestTableCreateIdentity(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
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
			identity: &ColumnIdentity{2, 5},
		},
		{
			name:          "column4",
			_type:         "VARCHAR",
			nullable:      true,
			maskingPolicy: "TEST_MP",
		},
		{
			name:          "column5",
			_type:         "VARCHAR",
			nullable:      true,
			tags: []TagValue{
					{
						Name:     "columnTag",
						Database: "test_db",
						Schema:   "test_schema",
						Value:    "value",
					},
					{
						Name:     "columnTag2",
						Database: "test_db",
						Schema:   "test_schema",
						Value:    "value2",
					},
		         },
		},
	}

	s.WithColumns(Columns(cols))
	r.Equal(`"test_db"."test_schema"."test_table"`, s.QualifiedName())

	r.Equal(`CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT COMMENT '', "column2" VARCHAR COMMENT 'only populated when data is available', "column3" NUMBER(38,0) NOT NULL IDENTITY(2, 5) COMMENT '', "column4" VARCHAR WITH MASKING POLICY TEST_MP COMMENT '', "column5" VARCHAR WITH TAG ("test_db"."test_schema"."columnTag" = "value", "test_db"."test_schema"."columnTag2" = "value2") COMMENT '') DATA_RETENTION_TIME_IN_DAYS = 0 CHANGE_TRACKING = false`, s.Create())
}

func TestTableChangeComment(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" SET COMMENT = 'new table comment'`, s.ChangeComment("new table comment"))
}

func TestTableRemoveComment(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" UNSET COMMENT`, s.RemoveComment())
}

func TestTableAddColumn(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" ADD COLUMN "new_column" VARIANT COMMENT ''`, s.AddColumn("new_column", "VARIANT", true, nil, nil, "", "", nil))
}

func TestTableAddColumnWithComment(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" ADD COLUMN "new_column" VARIANT COMMENT 'some comment'`, s.AddColumn("new_column", "VARIANT", true, nil, nil, "some comment", "", nil))
}

func TestTableAddColumnWithDefault(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" ADD COLUMN "new_column" NUMBER(38,0) DEFAULT 1 COMMENT ''`, s.AddColumn("new_column", "NUMBER(38,0)", true, NewColumnDefaultWithConstant("1"), nil, "", "", nil))
}

func TestTableAddColumnWithIdentity(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" ADD COLUMN "new_column" NUMBER(38,0) IDENTITY(1, 4) COMMENT ''`, s.AddColumn("new_column", "NUMBER(38,0)", true, nil, &ColumnIdentity{1, 4}, "", "", nil))
}

func TestTableAddColumnWithMaskingPolicy(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" ADD COLUMN "new_column" NUMBER(38,0) IDENTITY(1, 4) WITH MASKING POLICY TEST_MP COMMENT ''`, s.AddColumn("new_column", "NUMBER(38,0)", true, nil, &ColumnIdentity{1, 4}, "", "TEST_MP", nil))
}

func TestTableAddColumnWithTags(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	tags := []TagValue{
		{
			Name:     "columnTag",
			Database: "test_db",
			Schema:   "test_schema",
			Value:    "value",
		},
		{
			Name:     "columnTag2",
			Database: "test_db",
			Schema:   "test_schema",
			Value:    "value2",
		},
	}
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" ADD COLUMN "new_column" NUMBER(38,0) IDENTITY(1, 4) WITH TAG ("test_db"."test_schema"."columnTag" = "value", "test_db"."test_schema"."columnTag2" = "value2") COMMENT ''`, s.AddColumn("new_column", "NUMBER(38,0)", true, nil, &ColumnIdentity{1, 4}, "", "", tags))
}

func TestTableDropColumn(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" DROP COLUMN "old_column"`, s.DropColumn("old_column"))
}

func TestTableChangeColumnType(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" MODIFY COLUMN "old_column" BIGINT`, s.ChangeColumnType("old_column", "BIGINT"))
}

func TestTableChangeColumnComment(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" MODIFY COLUMN "old_column" COMMENT 'some comment'`, s.ChangeColumnComment("old_column", "some comment"))
}

func TestTableChangeColumnMaskingPolicy(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" MODIFY COLUMN "old_column" SET MASKING POLICY TEST_MP`, s.ChangeColumnMaskingPolicy("old_column", "TEST_MP"))
}

func TestTableChangeColumnTags(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	tags := []TagValue{
		{
			Name:     "columnTag",
			Database: "test_db",
			Schema:   "test_schema",
			Value:    "value",
		},
		{
			Name:     "columnTag2",
			Database: "test_db",
			Schema:   "test_schema",
			Value:    "value2",
		},
	}
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" MODIFY COLUMN "old_column" SET TAG "test_db"."test_schema"."columnTag" = "value", "test_db"."test_schema"."columnTag2" = "value2"`, s.ChangeColumnTags("old_column", tags))
}
func TestTableDropColumnDefault(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" MODIFY COLUMN "old_column" DROP DEFAULT`, s.DropColumnDefault("old_column"))
}

func TestTableChangeClusterBy(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" CLUSTER BY LINEAR(column2, column3)`, s.ChangeClusterBy("column2, column3"))
}

func TestTableChangeDataRetention(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" SET DATA_RETENTION_TIME_IN_DAYS = 5`, s.ChangeDataRetention(5))
}

func TestTableChangeChangeTracking(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" SET CHANGE_TRACKING = true`, s.ChangeChangeTracking(true))
}

func TestTableDropClusterBy(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" DROP CLUSTERING KEY`, s.DropClustering())
}

func TestTableDrop(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`DROP TABLE "test_db"."test_schema"."test_table"`, s.Drop())
}

func TestTableShow(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`SHOW TABLES LIKE 'test_table' IN SCHEMA "test_db"."test_schema"`, s.Show())
}

func TestTableShowPrimaryKeys(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`SHOW PRIMARY KEYS IN TABLE "test_db"."test_schema"."test_table"`, s.ShowPrimaryKeys())
}

func TestTableDropPrimaryKeys(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" DROP PRIMARY KEY`, s.DropPrimaryKey())
}

func TestTableChangePrimaryKeysWithConstraintName(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" ADD CONSTRAINT "MY_KEY" PRIMARY KEY("column1", "column2")`, s.ChangePrimaryKey(PrimaryKey{name: "MY_KEY", keys: []string{"column1", "column2"}}))
}

func TestTableChangePrimaryKeysWithoutConstraintName(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" ADD PRIMARY KEY("column1", "column2")`, s.ChangePrimaryKey(PrimaryKey{name: "", keys: []string{"column1", "column2"}}))
}

func TestTableAddTag(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" SET TAG "test_db"."test_schema"."tag" = "value"`, s.AddTag(TagValue{Name: "tag", Schema: "test_schema", Database: "test_db", Value: "value"}))
}

func TestTableChangeTag(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" SET TAG "test_db"."test_schema"."tag" = "value"`, s.ChangeTag(TagValue{Name: "tag", Schema: "test_schema", Database: "test_db", Value: "value"}))
}

func TestTableUnsetTag(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table" UNSET TAG "test_db"."test_schema"."tag"`, s.UnsetTag(TagValue{Name: "tag", Schema: "test_schema", Database: "test_db"}))
}

func TestTableRename(t *testing.T) {
	r := require.New(t)
	s := NewTableBuilder("test_table1", "test_db", "test_schema")
	r.Equal(`ALTER TABLE "test_db"."test_schema"."test_table1" RENAME TO "test_db"."test_schema"."test_table2"`, s.Rename("test_table2"))
}
