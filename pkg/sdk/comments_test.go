package sdk

import (
	"testing"
)

func TestComments(t *testing.T) {
	t.Run("set on schema", func(t *testing.T) {
		id := NewDatabaseObjectIdentifier("db1", "schema2")
		opts := &SetCommentOptions{
			ObjectType: ObjectTypeSchema,
			ObjectName: &id,
			Value:      String("mycomment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `COMMENT ON SCHEMA "db1"."schema2" IS 'mycomment'`)
	})

	t.Run("set if exists", func(t *testing.T) {
		id := NewAccountObjectIdentifier("maskpol")
		opts := &SetCommentOptions{
			IfExists:   Bool(true),
			ObjectType: ObjectTypeMaskingPolicy,
			ObjectName: &id,
			Value:      String("mycomment2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `COMMENT IF EXISTS ON MASKING POLICY "maskpol" IS 'mycomment2'`)
	})

	t.Run("set column comment", func(t *testing.T) {
		opts := &SetColumnCommentOptions{
			Column: NewDatabaseObjectIdentifier("table3", "column4"),
			Value:  String("mycomment3"),
		}
		assertOptsValidAndSQLEquals(t, opts, `COMMENT ON COLUMN "table3"."column4" IS 'mycomment3'`)
	})
}
