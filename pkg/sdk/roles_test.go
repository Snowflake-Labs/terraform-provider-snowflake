package sdk

import (
	"errors"
	"testing"
)

func TestRolesCreate(t *testing.T) {
	t.Run("if not exists", func(t *testing.T) {
		opts := &CreateRoleOptions{
			name:        NewAccountObjectIdentifier("new_role"),
			IfNotExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE ROLE IF NOT EXISTS "new_role"`)
	})

	t.Run("all options", func(t *testing.T) {
		opts := &CreateRoleOptions{
			name:      NewAccountObjectIdentifier("new_role"),
			OrReplace: Bool(true),
			Tag: []TagAssociation{
				{
					Name:  NewSchemaObjectIdentifier("db1", "schema1", "tag1"),
					Value: "v1",
				},
			},
			Comment: String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE ROLE "new_role" COMMENT = 'comment' TAG ("db1"."schema1"."tag1" = 'v1')`)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := &CreateRoleOptions{
			name: NewAccountObjectIdentifier(""),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: one of OrReplace, IfNotExists", func(t *testing.T) {
		opts := &CreateRoleOptions{
			name:        RandomAccountObjectIdentifier(),
			IfNotExists: Bool(true),
			OrReplace:   Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("OrReplace", "IfNotExists"))
	})
}

func TestRolesDrop(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		opts := &DropRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP ROLE "new_role"`)
	})

	t.Run("if exists", func(t *testing.T) {
		opts := &DropRoleOptions{
			name:     NewAccountObjectIdentifier("new_role"),
			IfExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP ROLE IF EXISTS "new_role"`)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := &DropRoleOptions{
			name: NewAccountObjectIdentifier(""),
		}
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})
}

func TestRolesAlter(t *testing.T) {
	t.Run("rename to", func(t *testing.T) {
		newID := NewAccountObjectIdentifier("new_role123")
		opts := &AlterRoleOptions{
			name:     NewAccountObjectIdentifier("new_role"),
			RenameTo: &newID,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROLE "new_role" RENAME TO "new_role123"`)
	})

	t.Run("set comment", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name:       NewAccountObjectIdentifier("new_role"),
			SetComment: String("some comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROLE "new_role" SET COMMENT = 'some comment'`)
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name:         NewAccountObjectIdentifier("new_role"),
			UnsetComment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROLE "new_role" UNSET COMMENT`)
	})

	t.Run("set tags", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			SetTags: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag-name"),
					Value: "tag-value",
				},
				{
					Name:  NewAccountObjectIdentifier("tag-name2"),
					Value: "tag-value2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROLE "new_role" SET TAG "tag-name" = 'tag-value', "tag-name2" = 'tag-value2'`)
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			UnsetTags: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag-name"),
				NewAccountObjectIdentifier("tag-name2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROLE "new_role" UNSET TAG "tag-name", "tag-name2"`)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name:         NewAccountObjectIdentifier(""),
			UnsetComment: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action specified", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name: RandomAccountObjectIdentifier(),
		}
		assertOptsInvalidJoinedErrors(t, opts, errors.New("no alter action specified"))
	})

	t.Run("validation: more than one alter action specified", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name:         RandomAccountObjectIdentifier(),
			SetComment:   String("comment"),
			UnsetComment: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("RenameTo", "SetComment", "UnsetComment", "SetTags", "UnsetTags"))
	})
}

func TestRolesShow(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		assertOptsValidAndSQLEquals(t, &ShowRoleOptions{}, `SHOW ROLES`)
	})

	t.Run("like", func(t *testing.T) {
		opts := &ShowRoleOptions{
			Like: &Like{
				Pattern: String("new_role"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW ROLES LIKE 'new_role'`)
	})

	t.Run("in class", func(t *testing.T) {
		class := NewAccountObjectIdentifier("some_class")
		opts := &ShowRoleOptions{
			InClass: &RolesInClass{
				Class: &class,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW ROLES IN CLASS "some_class"`)
	})

	t.Run("validation: like with no pattern", func(t *testing.T) {
		opts := &ShowRoleOptions{
			Like: &Like{},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrPatternRequiredForLikeKeyword)
	})

	t.Run("validation: invalid class name", func(t *testing.T) {
		class := NewAccountObjectIdentifier("")
		opts := &ShowRoleOptions{
			InClass: &RolesInClass{
				Class: &class,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
}

func TestRolesGrant(t *testing.T) {
	t.Run("user grant", func(t *testing.T) {
		opts := &GrantRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Grant: GrantRole{
				User: &AccountObjectIdentifier{name: "some_user"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT ROLE "new_role" TO USER "some_user"`)
	})

	t.Run("role grant", func(t *testing.T) {
		opts := &GrantRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Grant: GrantRole{
				Role: &AccountObjectIdentifier{name: "parent_role"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT ROLE "new_role" TO ROLE "parent_role"`)
	})

	t.Run("validation: invalid object identifier and no grant option", func(t *testing.T) {
		opts := &GrantRoleOptions{
			name: NewAccountObjectIdentifier(""),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier, errors.New("only one grant option can be set [TO ROLE or TO USER]"))
	})

	t.Run("validation: invalid object identifier for granted role", func(t *testing.T) {
		id := NewAccountObjectIdentifier("")
		opts := &GrantRoleOptions{
			name: RandomAccountObjectIdentifier(),
			Grant: GrantRole{
				Role: &id,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errors.New("invalid object identifier for granted role"))
	})

	t.Run("validation: invalid object identifier for granted user", func(t *testing.T) {
		id := NewAccountObjectIdentifier("")
		opts := &GrantRoleOptions{
			name: RandomAccountObjectIdentifier(),
			Grant: GrantRole{
				User: &id,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errors.New("invalid object identifier for granted user"))
	})
}

func TestRolesRevoke(t *testing.T) {
	t.Run("revoke user", func(t *testing.T) {
		opts := &RevokeRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Revoke: RevokeRole{
				User: &AccountObjectIdentifier{name: "some_user"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE ROLE "new_role" FROM USER "some_user"`)
	})

	t.Run("revoke role", func(t *testing.T) {
		opts := &RevokeRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Revoke: RevokeRole{
				Role: &AccountObjectIdentifier{name: "parent_role"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE ROLE "new_role" FROM ROLE "parent_role"`)
	})

	t.Run("validation: invalid object identifier and no option set", func(t *testing.T) {
		opts := &RevokeRoleOptions{
			name: NewAccountObjectIdentifier(""),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier, errors.New("only one revoke option can be set [FROM ROLE or FROM USER]"))
	})
}
