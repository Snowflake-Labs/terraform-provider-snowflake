package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaCreate(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.QualifiedName(), `"test"`)

	s.WithDB("db")
	a.Equal(s.QualifiedName(), `"db"."test"`)

	a.Equal(s.Create(), `CREATE SCHEMA "db"."test"`)

	s.Transient()
	a.Equal(s.Create(), `CREATE TRANSIENT SCHEMA "db"."test"`)

	s.Managed()
	a.Equal(s.Create(), `CREATE TRANSIENT SCHEMA "db"."test" WITH MANAGED ACCESS`)

	s.WithDataRetentionDays(7)
	a.Equal(s.Create(), `CREATE TRANSIENT SCHEMA "db"."test" WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 7`)

	s.WithComment("Yeehaw")
	a.Equal(s.Create(), `CREATE TRANSIENT SCHEMA "db"."test" WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 7 COMMENT = 'Yeehaw'`)
}

func TestSchemaRename(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.Rename("bob"), `ALTER SCHEMA "test" RENAME TO "bob"`)
}

func TestSchemaSwap(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.Swap("bob"), `ALTER SCHEMA "test" SWAP WITH "bob"`)
}

func TestSchemaChangeComment(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.ChangeComment("worst schema ever"), `ALTER SCHEMA "test" SET COMMENT = 'worst schema ever'`)
}

func TestSchemaRemoveComment(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.RemoveComment(), `ALTER SCHEMA "test" UNSET COMMENT`)
}

func TestSchemaChangeDataRetentionDays(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.ChangeDataRetentionDays(22), `ALTER SCHEMA "test" SET DATA_RETENTION_TIME_IN_DAYS = 22`)
}

func TestSchemaRemoveDataRetentionDays(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.RemoveDataRetentionDays(), `ALTER SCHEMA "test" UNSET DATA_RETENTION_TIME_IN_DAYS`)
}

func TestSchemaManage(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.Manage(), `ALTER SCHEMA "test" ENABLE MANAGED ACCESS`)
}

func TestSchemaUnmanage(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.Unmanage(), `ALTER SCHEMA "test" DISABLE MANAGED ACCESS`)
}

func TestSchemaDrop(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.Drop(), `DROP SCHEMA "test"`)
}

func TestSchemaUndrop(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.Undrop(), `UNDROP SCHEMA "test"`)
}

func TestSchemaUse(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.Use(), `USE SCHEMA "test"`)
}

func TestSchemaShow(t *testing.T) {
	a := require.New(t)
	s := Schema("test")
	a.Equal(s.Show(), `SHOW SCHEMAS LIKE 'test'`)

	s.WithDB("db")
	a.Equal(s.Show(), `SHOW SCHEMAS LIKE 'test' IN DATABASE "db"`)
}
