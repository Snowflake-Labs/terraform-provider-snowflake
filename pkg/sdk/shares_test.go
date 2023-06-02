package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSharesCreate(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &CreateShareOptions{
			name: NewAccountObjectIdentifier("myshare"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE SHARE "myshare"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with complete options", func(t *testing.T) {
		comment := randomComment(t)
		opts := &CreateShareOptions{
			OrReplace: Bool(true),
			name:      NewAccountObjectIdentifier("complete_share"),
			Comment:   String(comment),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE SHARE "complete_share" COMMENT = '` + comment + `'`
		assert.Equal(t, expected, actual)
	})
}

func TestShareAlter(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &AlterShareOptions{
			name: NewAccountObjectIdentifier("myshare"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SHARE "myshare"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with add", func(t *testing.T) {
		accounts := []AccountIdentifier{NewAccountIdentifier("my-org", "myaccount")}
		opts := &AlterShareOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Add: &ShareAdd{
				Accounts:          accounts,
				ShareRestrictions: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SHARE IF EXISTS "myshare" ADD ACCOUNTS = "my-org.myaccount" SHARE_RESTRICTIONS = true`
		assert.Equal(t, expected, actual)
	})

	t.Run("with remove", func(t *testing.T) {
		accounts := []AccountIdentifier{NewAccountIdentifier("my-org", "myaccount"), NewAccountIdentifier("my-org", "myaccount2")}
		opts := &AlterShareOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Remove: &ShareRemove{
				Accounts: accounts,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SHARE IF EXISTS "myshare" REMOVE ACCOUNTS = "my-org.myaccount", "my-org.myaccount2"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set", func(t *testing.T) {
		accounts := []AccountIdentifier{NewAccountIdentifier("my-org", "myaccount")}
		comment := randomComment(t)
		opts := &AlterShareOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Set: &ShareSet{
				Accounts: accounts,
				Comment:  &comment,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SHARE IF EXISTS "myshare" SET ACCOUNTS = "my-org.myaccount" COMMENT = '` + comment + `'`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set tag", func(t *testing.T) {
		opts := &AlterShareOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Set: &ShareSet{
				Tag: []TagAssociation{
					{
						Name:  NewSchemaObjectIdentifier("db", "schema", "tag"),
						Value: "v1",
					},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SHARE IF EXISTS "myshare" SET TAG "db"."schema"."tag" = 'v1'`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &AlterShareOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Unset: &ShareUnset{
				Comment: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SHARE IF EXISTS "myshare" UNSET COMMENT`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset tag", func(t *testing.T) {
		opts := &AlterShareOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Unset: &ShareUnset{
				Tag: []ObjectIdentifier{
					NewSchemaObjectIdentifier("db", "schema", "tag"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER SHARE IF EXISTS "myshare" UNSET TAG "db"."schema"."tag"`
		assert.Equal(t, expected, actual)
	})
}

func TestShareShow(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		opts := &ShowShareOptions{
			Like: &Like{
				Pattern: String("myshare"),
			},
			StartsWith: String("my"),
			Limit: &LimitFrom{
				Rows: Int(10),
				From: String("my_other_share"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW SHARES LIKE 'myshare' STARTS WITH 'my' LIMIT 10 FROM 'my_other_share'`
		assert.Equal(t, expected, actual)
	})
}

func TestShareDrop(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &shareDropOptions{
			name: NewAccountObjectIdentifier("myshare"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP SHARE "myshare"`
		assert.Equal(t, expected, actual)
	})
}

func TestShareDescribe(t *testing.T) {
	t.Run("describe provider", func(t *testing.T) {
		opts := &shareDescribeOptions{
			name: NewAccountObjectIdentifier("myprovider"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DESCRIBE SHARE "myprovider"`
		assert.Equal(t, expected, actual)
	})
	t.Run("describe consumer", func(t *testing.T) {
		opts := &shareDescribeOptions{
			name: NewAccountObjectIdentifier("myconsumer"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DESCRIBE SHARE "myconsumer"`
		assert.Equal(t, expected, actual)
	})
}
