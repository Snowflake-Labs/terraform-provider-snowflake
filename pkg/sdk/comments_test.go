package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestComments(t *testing.T) {
	t.Run("set on schema", func(t *testing.T) {
		id := NewSchemaIdentifier("db1", "schema2")
		opts := &SetCommentOptions{
			ObjectType: ObjectTypeSchema,
			ObjectName: &id,
			Value:      String("mycomment"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `COMMENT ON SCHEMA "db1"."schema2" IS 'mycomment'`
		assert.Equal(t, expected, actual)
	})

	t.Run("set if exists", func(t *testing.T) {
		id := NewAccountObjectIdentifier("maskpol")
		opts := &SetCommentOptions{
			IfExists:   Bool(true),
			ObjectType: ObjectTypeMaskingPolicy,
			ObjectName: &id,
			Value:      String("mycomment2"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `COMMENT IF EXISTS ON MASKING POLICY "maskpol" IS 'mycomment2'`
		assert.Equal(t, expected, actual)
	})

	t.Run("set column comment", func(t *testing.T) {
		opts := &SetColumnCommentOptions{
			Column: NewSchemaIdentifier("table3", "column4"),
			Value:  String("mycomment3"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `COMMENT ON COLUMN "table3"."column4" IS 'mycomment3'`
		assert.Equal(t, expected, actual)
	})
}
