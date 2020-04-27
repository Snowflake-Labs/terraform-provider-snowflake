package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaCreate(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.QualifiedName(), `"test"`)

	s.WithDB("db")
	r.Equal(s.QualifiedName(), `"db"."test"`)

	r.Equal(s.Create(), `CREATE SCHEMA "db"."test"`)

	s.Transient()
	r.Equal(s.Create(), `CREATE TRANSIENT SCHEMA "db"."test"`)

	s.Managed()
	r.Equal(s.Create(), `CREATE TRANSIENT SCHEMA "db"."test" WITH MANAGED ACCESS`)

	s.WithDataRetentionDays(7)
	r.Equal(s.Create(), `CREATE TRANSIENT SCHEMA "db"."test" WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 7`)

	s.WithComment("Yeehaw")
	r.Equal(s.Create(), `CREATE TRANSIENT SCHEMA "db"."test" WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 7 COMMENT = 'Yeehaw'`)
}

func TestSchemaRename(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.Rename("bob"), `ALTER SCHEMA "test" RENAME TO "bob"`)
}

func TestSchemaSwap(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.Swap("bob"), `ALTER SCHEMA "test" SWAP WITH "bob"`)
}

func TestSchemaChangeComment(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.ChangeComment("worst schema ever"), `ALTER SCHEMA "test" SET COMMENT = 'worst schema ever'`)
}

func TestSchemaRemoveComment(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.RemoveComment(), `ALTER SCHEMA "test" UNSET COMMENT`)
}

func TestSchemaChangeDataRetentionDays(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.ChangeDataRetentionDays(22), `ALTER SCHEMA "test" SET DATA_RETENTION_TIME_IN_DAYS = 22`)
}

func TestSchemaRemoveDataRetentionDays(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.RemoveDataRetentionDays(), `ALTER SCHEMA "test" UNSET DATA_RETENTION_TIME_IN_DAYS`)
}

func TestSchemaManage(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.Manage(), `ALTER SCHEMA "test" ENABLE MANAGED ACCESS`)
}

func TestSchemaUnmanage(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.Unmanage(), `ALTER SCHEMA "test" DISABLE MANAGED ACCESS`)
}

func TestSchemaDrop(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.Drop(), `DROP SCHEMA "test"`)
}

func TestSchemaUndrop(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.Undrop(), `UNDROP SCHEMA "test"`)
}

func TestSchemaUse(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.Use(), `USE SCHEMA "test"`)
}

func TestSchemaShow(t *testing.T) {
	r := require.New(t)
	s := Schema("test")
	r.Equal(s.Show(), `SHOW SCHEMAS LIKE 'test'`)

	s.WithDB("db")
	r.Equal(s.Show(), `SHOW SCHEMAS LIKE 'test' IN DATABASE "db"`)
}
