package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
)

func TestSharesCreate(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &CreateShareOptions{
			name: NewAccountObjectIdentifier("myshare"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE SHARE "myshare"`)
	})

	t.Run("with complete options", func(t *testing.T) {
		comment := random.Comment()
		opts := &CreateShareOptions{
			OrReplace: Bool(true),
			name:      NewAccountObjectIdentifier("complete_share"),
			Comment:   String(comment),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SHARE "complete_share" COMMENT = '%s'`, comment)
	})
}

func TestShareAlter(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &AlterShareOptions{
			name: NewAccountObjectIdentifier("myshare"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterShareOptions", "Add", "Remove", "Set", "Unset", "SetTag", "UnsetTag"))
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER SHARE IF EXISTS "myshare" ADD ACCOUNTS = "my-org"."myaccount" SHARE_RESTRICTIONS = true`)
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER SHARE IF EXISTS "myshare" REMOVE ACCOUNTS = "my-org"."myaccount", "my-org"."myaccount2"`)
	})

	t.Run("with set", func(t *testing.T) {
		accounts := []AccountIdentifier{NewAccountIdentifier("my-org", "myaccount")}
		comment := random.Comment()
		opts := &AlterShareOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Set: &ShareSet{
				Accounts: accounts,
				Comment:  &comment,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SHARE IF EXISTS "myshare" SET ACCOUNTS = "my-org"."myaccount" COMMENT = '%s'`, comment)
	})

	t.Run("with set tag", func(t *testing.T) {
		opts := &AlterShareOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			SetTag: []TagAssociation{
				{
					Name:  NewSchemaObjectIdentifier("db", "schema", "tag"),
					Value: "v1",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SHARE IF EXISTS "myshare" SET TAG "db"."schema"."tag" = 'v1'`)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &AlterShareOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			Unset: &ShareUnset{
				Comment: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SHARE IF EXISTS "myshare" UNSET COMMENT`)
	})

	t.Run("with unset tag", func(t *testing.T) {
		opts := &AlterShareOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("myshare"),
			UnsetTag: []ObjectIdentifier{
				NewSchemaObjectIdentifier("db", "schema", "tag"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SHARE IF EXISTS "myshare" UNSET TAG "db"."schema"."tag"`)
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
		assertOptsValidAndSQLEquals(t, opts, `SHOW SHARES LIKE 'myshare' STARTS WITH 'my' LIMIT 10 FROM 'my_other_share'`)
	})
}

func TestShareDrop(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &dropShareOptions{
			name: NewAccountObjectIdentifier("myshare"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP SHARE "myshare"`)
	})
}

func TestShareDescribe(t *testing.T) {
	t.Run("describe provider", func(t *testing.T) {
		opts := &describeShareOptions{
			name: NewAccountObjectIdentifier("myprovider"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE SHARE "myprovider"`)
	})

	t.Run("describe consumer", func(t *testing.T) {
		opts := &describeShareOptions{
			name: NewAccountObjectIdentifier("myconsumer"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE SHARE "myconsumer"`)
	})
}
