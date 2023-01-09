package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagCreate(t *testing.T) {
	r := require.New(t)
	o := NewTagBuilder("test")
	r.Equal(`"test"`, o.QualifiedName())

	o.WithDB("db")
	r.Equal(`"db"."test"`, o.QualifiedName())

	o.WithSchema("schema")
	r.Equal(`"db"."schema"."test"`, o.QualifiedName())

	r.Equal(`CREATE TAG "db"."schema"."test"`, o.Create())

	allowedValues := []string{"marketing", "finance"}
	o.WithAllowedValues(allowedValues)

	o.WithComment("Yee'haw")
	r.Equal(`CREATE TAG "db"."schema"."test" ALLOWED_VALUES 'marketing', 'finance' COMMENT = 'Yee\'haw'`, o.Create())
}

func TestTagRename(t *testing.T) {
	r := require.New(t)
	o := NewTagBuilder("test")
	r.Equal(`ALTER TAG "test" RENAME TO "bob"`, o.Rename("bob"))
}

func TestTagChangeComment(t *testing.T) {
	r := require.New(t)
	o := NewTagBuilder("test")
	r.Equal(`ALTER TAG "test" SET COMMENT = 'worst\' tag ever'`, o.ChangeComment("worst' tag ever"))
}

func TestTagRemoveComment(t *testing.T) {
	r := require.New(t)
	o := NewTagBuilder("test")
	r.Equal(`ALTER TAG "test" UNSET COMMENT`, o.RemoveComment())
}

func TestTagAddAllowedValues(t *testing.T) {
	r := require.New(t)
	o := NewTagBuilder("test")
	avs := []string{"foo", "bar"}
	r.Equal(`ALTER TAG "test" ADD ALLOWED_VALUES 'foo', 'bar'`, o.AddAllowedValues(avs))
}

func TestTagDropAllowedValues(t *testing.T) {
	r := require.New(t)
	o := NewTagBuilder("test")
	avs := []string{"foo"}
	r.Equal(`ALTER TAG "test" DROP ALLOWED_VALUES 'foo'`, o.DropAllowedValues(avs))
}

func TestTagRemoveAllowedValues(t *testing.T) {
	r := require.New(t)
	o := NewTagBuilder("test")
	r.Equal(`ALTER TAG "test" UNSET ALLOWED_VALUES`, o.RemoveAllowedValues())
}

func TestTagDrop(t *testing.T) {
	r := require.New(t)
	o := NewTagBuilder("test")
	r.Equal(`DROP TAG "test"`, o.Drop())
}

func TestTagUndrop(t *testing.T) {
	r := require.New(t)
	o := NewTagBuilder("test")
	r.Equal(`UNDROP TAG "test"`, o.Undrop())
}

func TestTagShow(t *testing.T) {
	r := require.New(t)
	o := NewTagBuilder("test")
	r.Equal(`SHOW TAGS LIKE 'test'`, o.Show())

	o.WithDB("db")
	r.Equal(`SHOW TAGS LIKE 'test' IN DATABASE "db"`, o.Show())

	o.WithSchema("schema")
	r.Equal(`SHOW TAGS LIKE 'test' IN SCHEMA "db"."schema"`, o.Show())
}

func TestTagShowAttachedPolicy(t *testing.T) {
	r := require.New(t)
	o := NewTagBuilder("test")
	r.Equal(`SHOW TAGS LIKE 'test'`, o.Show())

	o.WithDB("db")
	r.Equal(`SHOW TAGS LIKE 'test' IN DATABASE "db"`, o.Show())

	o.WithSchema("schema")
	r.Equal(`SHOW TAGS LIKE 'test' IN SCHEMA "db"."schema"`, o.Show())

	mP := MaskingPolicy("policy", "db2", "schema2")
	o.WithMaskingPolicy(mP)
	r.Equal(`SELECT * from table ("db".information_schema.policy_references(ref_entity_name => '"db"."schema"."test"', ref_entity_domain => 'TAG')) where policy_db='db2' and policy_schema='schema2' and policy_name='policy'`, o.ShowAttachedPolicy())
}
