package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRolesCreate(t *testing.T) {
	t.Run("if not exists", func(t *testing.T) {
		opts := &RoleCreateOptions{
			name:        NewAccountObjectIdentifier("new_role"),
			IfNotExists: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE ROLE IF NOT EXISTS "new_role"`
		assert.Equal(t, expected, actual)
	})

	t.Run("all options", func(t *testing.T) {
		opts := &RoleCreateOptions{
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
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE ROLE "new_role" COMMENT = 'comment' TAG ("db1"."schema1"."tag1" = 'v1')`
		assert.Equal(t, expected, actual)
	})
}

func TestRolesDrop(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		opts := &RoleDropOptions{
			name: NewAccountObjectIdentifier("new_role"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP ROLE "new_role"`
		assert.Equal(t, expected, actual)
	})

	t.Run("if exists", func(t *testing.T) {
		opts := &RoleDropOptions{
			name:     NewAccountObjectIdentifier("new_role"),
			IfExists: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP ROLE IF EXISTS "new_role"`
		assert.Equal(t, expected, actual)
	})
}

func TestRolesAlter(t *testing.T) {
	t.Run("rename to", func(t *testing.T) {
		opts := &RoleAlterOptions{
			name:     NewAccountObjectIdentifier("new_role"),
			RenameTo: NewAccountObjectIdentifier("new_role123"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ROLE "new_role" RENAME TO "new_role123"`
		assert.Equal(t, expected, actual)
	})

	t.Run("set comment", func(t *testing.T) {
		opts := &RoleAlterOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Set: &RoleSet{
				Comment: String("some comment"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ROLE "new_role" SET COMMENT = 'some comment'`
		assert.Equal(t, expected, actual)
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := &RoleAlterOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Unset: &RoleUnset{
				Comment: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ROLE "new_role" UNSET COMMENT`
		assert.Equal(t, expected, actual)
	})

	t.Run("set tags", func(t *testing.T) {
		opts := &RoleAlterOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Set: &RoleSet{
				Tag: []TagAssociation{
					{
						Name:  NewAccountObjectIdentifier("tagname"),
						Value: "tagvalue",
					},
					{
						Name:  NewAccountObjectIdentifier("tagname2"),
						Value: "tagvalue2",
					},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ROLE "new_role" SET TAG "tagname" = 'tagvalue', "tagname2" = 'tagvalue2'`
		assert.Equal(t, expected, actual)
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := &RoleAlterOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Unset: &RoleUnset{
				Tag: []ObjectIdentifier{
					NewAccountObjectIdentifier("tagname"),
					NewAccountObjectIdentifier("tagname2"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER ROLE "new_role" UNSET TAG "tagname", "tagname2"`
		assert.Equal(t, expected, actual)
	})
}

func TestRolesShow(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		opts := &RoleShowOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW ROLES`
		assert.Equal(t, expected, actual)
	})

	t.Run("like", func(t *testing.T) {
		opts := &RoleShowOptions{
			Like: &Like{
				Pattern: String("new_role"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW ROLES LIKE 'new_role'`
		assert.Equal(t, expected, actual)
	})
}

func TestTest(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := []int{4, 5, 6}
	c := make([]int, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	fmt.Println(c)
}

func TestRolesUse(t *testing.T) {
	opts := &RoleUseOptions{
		name: NewAccountObjectIdentifier("new_role"),
	}
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	expected := `USE ROLE "new_role"`
	assert.Equal(t, expected, actual)
}

func TestRolesGrant(t *testing.T) {
	t.Run("user grant", func(t *testing.T) {
		opts := &RoleGrantOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Grant: GrantRole{
				User: &AccountObjectIdentifier{name: "some_user"},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `GRANT ROLE "new_role" TO USER "some_user"`
		assert.Equal(t, expected, actual)
	})

	t.Run("role grant", func(t *testing.T) {
		opts := &RoleGrantOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Grant: GrantRole{
				Role: &AccountObjectIdentifier{name: "parent_role"},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `GRANT ROLE "new_role" TO ROLE "parent_role"`
		assert.Equal(t, expected, actual)
	})
}

func TestRolesRevoke(t *testing.T) {
	t.Run("revoke user", func(t *testing.T) {
		opts := &RoleRevokeOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Revoke: RevokeRole{
				User: &AccountObjectIdentifier{name: "some_user"},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `REVOKE ROLE "new_role" FROM USER "some_user"`
		assert.Equal(t, expected, actual)
	})

	t.Run("revoke role", func(t *testing.T) {
		opts := &RoleRevokeOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Revoke: RevokeRole{
				Role: &AccountObjectIdentifier{name: "parent_role"},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `REVOKE ROLE "new_role" FROM ROLE "parent_role"`
		assert.Equal(t, expected, actual)
	})
}

func TestRolesUseSecondaryRoles(t *testing.T) {
	t.Run("use 'ALL' secondary functions", func(t *testing.T) {
		opts := &RoleUseSecondaryOptions{
			SecondaryRoleOption: AllSecondaryRoles,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `USE SECONDARY ROLES ALL`
		assert.Equal(t, expected, actual)
	})

	t.Run("use 'NONE' secondary functions", func(t *testing.T) {
		opts := &RoleUseSecondaryOptions{
			SecondaryRoleOption: NoneSecondaryRoles,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `USE SECONDARY ROLES NONE`
		assert.Equal(t, expected, actual)
	})
}
