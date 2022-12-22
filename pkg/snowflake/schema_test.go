package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaCreate(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`"test"`, s.QualifiedName())

	s.WithDB("db")
	r.Equal(`"db"."test"`, s.QualifiedName())

	r.Equal(`CREATE SCHEMA "db"."test"`, s.Create())

	s.Transient()
	r.Equal(`CREATE TRANSIENT SCHEMA "db"."test"`, s.Create())

	s.Managed()
	r.Equal(`CREATE TRANSIENT SCHEMA "db"."test" WITH MANAGED ACCESS`, s.Create())

	s.WithDataRetentionDays(7)
	r.Equal(`CREATE TRANSIENT SCHEMA "db"."test" WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 7`, s.Create())

	s.WithComment("Yee'haw")
	r.Equal(`CREATE TRANSIENT SCHEMA "db"."test" WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 7 COMMENT = 'Yee\'haw'`, s.Create())
}

func TestSchemaRename(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`ALTER SCHEMA "test" RENAME TO "bob"`, s.Rename("bob"))
}

func TestSchemaSwap(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`ALTER SCHEMA "test" SWAP WITH "bob"`, s.Swap("bob"))
}

func TestSchemaChangeComment(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`ALTER SCHEMA "test" SET COMMENT = 'worst\' schema ever'`, s.ChangeComment("worst' schema ever"))
}

func TestSchemaRemoveComment(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`ALTER SCHEMA "test" UNSET COMMENT`, s.RemoveComment())
}

func TestSchemaChangeDataRetentionDays(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`ALTER SCHEMA "test" SET DATA_RETENTION_TIME_IN_DAYS = 22`, s.ChangeDataRetentionDays(22))
}

func TestSchemaRemoveDataRetentionDays(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`ALTER SCHEMA "test" UNSET DATA_RETENTION_TIME_IN_DAYS`, s.RemoveDataRetentionDays())
}

func TestSchemaManage(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`ALTER SCHEMA "test" ENABLE MANAGED ACCESS`, s.Manage())
}

func TestSchemaUnmanage(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`ALTER SCHEMA "test" DISABLE MANAGED ACCESS`, s.Unmanage())
}

func TestSchemaDrop(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`DROP SCHEMA "test"`, s.Drop())
}

func TestSchemaUndrop(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`UNDROP SCHEMA "test"`, s.Undrop())
}

func TestSchemaUse(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`USE SCHEMA "test"`, s.Use())
}

func TestSchemaShow(t *testing.T) {
	r := require.New(t)
	s := NewSchemaBuilder("test")
	r.Equal(`SHOW SCHEMAS LIKE 'test'`, s.Show())

	s.WithDB("db")
	r.Equal(`SHOW SCHEMAS LIKE 'test' IN DATABASE "db"`, s.Show())
}
