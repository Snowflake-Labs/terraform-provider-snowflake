package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSharesCreate(t *testing.T) {
	builder := testBuilder(t)

	t.Run("only name", func(t *testing.T) {
		opts := &ShareCreateOptions{
			name: NewAccountObjectIdentifier("myshare"),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`CREATE SHARE "myshare"`,
			builder.sql(clauses...),
		)
	})

	t.Run("with complete options", func(t *testing.T) {
		comment := randomComment(t)
		opts := &ShareCreateOptions{
			OrReplace: Bool(true),
			name:      NewAccountObjectIdentifier("complete_share"),
			Comment:   String(comment),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`CREATE OR REPLACE SHARE "complete_share" COMMENT = '`+comment+`'`,
			builder.sql(clauses...),
		)
	})
}

func TestShareAlter(t *testing.T) {
	builder := testBuilder(t)

	t.Run("only name", func(t *testing.T) {
		opts := &ShareAlterOptions{
			name: NewAccountObjectIdentifier("myshare"),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`ALTER SHARE "myshare"`,
			builder.sql(clauses...),
		)
	})

	t.Run("with add", func(t *testing.T) {
		accounts := []AccountIdentifier{NewAccountIdentifier("my-org", "myaccount")}
		opts := &ShareAlterOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Add: &ShareAdd{
				Accounts:          accounts,
				ShareRestrictions: Bool(true),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`ALTER SHARE IF EXISTS "myshare" ADD ACCOUNTS = "my-org.myaccount" SHARE_RESTRICTIONS = true`,
			builder.sql(clauses...),
		)
	})

	t.Run("with remove", func(t *testing.T) {
		accounts := []AccountIdentifier{NewAccountIdentifier("my-org", "myaccount"), NewAccountIdentifier("my-org", "myaccount2")}
		opts := &ShareAlterOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Remove: &ShareRemove{
				Accounts: accounts,
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`ALTER SHARE IF EXISTS "myshare" REMOVE ACCOUNTS = "my-org.myaccount","my-org.myaccount2"`,
			builder.sql(clauses...),
		)
	})

	t.Run("with set", func(t *testing.T) {
		accounts := []AccountIdentifier{NewAccountIdentifier("my-org", "myaccount")}
		comment := randomComment(t)
		opts := &ShareAlterOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Set: &ShareSet{
				Accounts: accounts,
				Comment:  &comment,
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`ALTER SHARE IF EXISTS "myshare" SET ACCOUNTS = "my-org.myaccount" COMMENT = '`+comment+`'`,
			builder.sql(clauses...),
		)
	})

	t.Run("with set tag", func(t *testing.T) {
		opts := &ShareAlterOptions{
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`ALTER SHARE IF EXISTS "myshare" SET TAG "db"."schema"."tag" = 'v1'`,
			builder.sql(clauses...),
		)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &ShareAlterOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Unset: &ShareUnset{
				Comment: Bool(true),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`ALTER SHARE IF EXISTS "myshare" UNSET COMMENT`,
			builder.sql(clauses...),
		)
	})

	t.Run("with unset tag", func(t *testing.T) {
		opts := &ShareAlterOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Unset: &ShareUnset{
				Tag: []ObjectIdentifier{
					NewSchemaObjectIdentifier("db", "schema", "tag"),
				},
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`ALTER SHARE IF EXISTS "myshare" UNSET TAG "db"."schema"."tag"`,
			builder.sql(clauses...),
		)
	})
}

func TestShareShow(t *testing.T) {
	builder := testBuilder(t)

	t.Run("complete", func(t *testing.T) {
		opts := &ShareShowOptions{
			Like: &Like{
				Pattern: String("myshare"),
			},
			StartsWith: String("my"),
			Limit: &LimitFrom{
				Rows: Int(10),
				From: String("my_other_share"),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`SHOW SHARES LIKE 'myshare' STARTS WITH 'my' LIMIT 10 FROM 'my_other_share'`,
			builder.sql(clauses...),
		)
	})
}
