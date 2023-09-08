package example

import "testing"

func TestDatabaseRoles_Create(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	defaultOpts := func() *CreateDatabaseRoleOptions {
		return &CreateDatabaseRoleOptions{
			name: id,
		}
	}
	// TODO: remove me
	_ = defaultOpts()

	// TODO: fill me

	// TODO: validate valid identifier for [opts.name]
	// TODO: validate conflicting fields for [opts.OrReplace opts.IfNotExists]

}

func TestDatabaseRoles_Alter(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	defaultOpts := func() *AlterDatabaseRoleOptions {
		return &AlterDatabaseRoleOptions{
			name: id,
		}
	}
	// TODO: remove me
	_ = defaultOpts()

	// TODO: fill me

	// TODO: validate valid identifier for [opts.name]
	// TODO: validate exactly one field from [opts.Rename opts.Set opts.Unset] is present

	// TODO: validate valid identifier for [opts.Rename.Name]

	// TODO: validate at least one of fields [opts.Set.NestedThirdLevel.Field] set

	// TODO: validate at least one of fields [opts.Unset.Comment] set

}
