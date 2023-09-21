package sdk

import (
	"fmt"

	// TODO: change to slices with go 1.21
	"golang.org/x/exp/slices"
)

var (
	_ validatable = new(GrantPrivilegesToAccountRoleOptions)
	_ validatable = new(RevokePrivilegesFromAccountRoleOptions)
	_ validatable = new(GrantPrivilegesToDatabaseRoleOptions)
	_ validatable = new(RevokePrivilegesFromDatabaseRoleOptions)
	_ validatable = new(grantPrivilegeToShareOptions)
	_ validatable = new(revokePrivilegeFromShareOptions)
	_ validatable = new(ShowGrantOptions)
)

func (opts *GrantPrivilegesToAccountRoleOptions) validate() error {
	if !valueSet(opts.privileges) {
		return fmt.Errorf("privileges must be set")
	}
	if err := opts.privileges.validate(); err != nil {
		return err
	}
	if !valueSet(opts.on) {
		return fmt.Errorf("on must be set")
	}
	if err := opts.on.validate(); err != nil {
		return err
	}
	return nil
}

func (v *AccountRoleGrantPrivileges) validate() error {
	if !exactlyOneValueSet(v.AllPrivileges, v.GlobalPrivileges, v.AccountObjectPrivileges, v.SchemaPrivileges, v.SchemaObjectPrivileges) {
		return fmt.Errorf("exactly one of AllPrivileges, GlobalPrivileges, AccountObjectPrivileges, SchemaPrivileges, or SchemaObjectPrivileges must be set")
	}
	return nil
}

func (v *AccountRoleGrantOn) validate() error {
	if !exactlyOneValueSet(v.Account, v.AccountObject, v.Schema, v.SchemaObject) {
		return fmt.Errorf("exactly one of Account, AccountObject, Schema, or SchemaObject must be set")
	}
	if valueSet(v.AccountObject) {
		if err := v.AccountObject.validate(); err != nil {
			return err
		}
	}
	if valueSet(v.Schema) {
		if err := v.Schema.validate(); err != nil {
			return err
		}
	}
	if valueSet(v.SchemaObject) {
		if err := v.SchemaObject.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (v *GrantOnAccountObject) validate() error {
	if !exactlyOneValueSet(v.User, v.ResourceMonitor, v.Warehouse, v.Database, v.Integration, v.FailoverGroup, v.ReplicationGroup) {
		return fmt.Errorf("exactly one of User, ResourceMonitor, Warehouse, Database, Integration, FailoverGroup, or ReplicationGroup must be set")
	}
	return nil
}

func (v *GrantOnSchema) validate() error {
	if !exactlyOneValueSet(v.Schema, v.AllSchemasInDatabase, v.FutureSchemasInDatabase) {
		return fmt.Errorf("exactly one of Schema, AllSchemasInDatabase, or FutureSchemasInDatabase must be set")
	}
	return nil
}

func (v *GrantOnSchemaObject) validate() error {
	if !exactlyOneValueSet(v.SchemaObject, v.All, v.Future) {
		return fmt.Errorf("exactly one of Object, AllIn or Future must be set")
	}
	if valueSet(v.All) {
		if err := v.All.validate(); err != nil {
			return err
		}
	}
	if valueSet(v.Future) {
		if err := v.Future.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (v *GrantOnSchemaObjectIn) validate() error {
	if !exactlyOneValueSet(v.InDatabase, v.InSchema) {
		return fmt.Errorf("exactly one of InDatabase, or InSchema must be set")
	}
	return nil
}

func (opts *RevokePrivilegesFromAccountRoleOptions) validate() error {
	if !valueSet(opts.privileges) {
		return fmt.Errorf("privileges must be set")
	}
	if err := opts.privileges.validate(); err != nil {
		return err
	}
	if !valueSet(opts.on) {
		return fmt.Errorf("on must be set")
	}
	if err := opts.on.validate(); err != nil {
		return err
	}
	if !validObjectidentifier(opts.accountRole) {
		return ErrInvalidObjectIdentifier
	}
	if everyValueSet(opts.Restrict, opts.Cascade) {
		return fmt.Errorf("either Restrict or Cascade can be set, or neither but not both")
	}
	return nil
}

func (opts *GrantPrivilegesToDatabaseRoleOptions) validate() error {
	if !valueSet(opts.privileges) {
		return fmt.Errorf("privileges must be set")
	}
	if err := opts.privileges.validate(); err != nil {
		return err
	}
	if !valueSet(opts.on) {
		return fmt.Errorf("on must be set")
	}
	if err := opts.on.validate(); err != nil {
		return err
	}
	return nil
}

func (v *DatabaseRoleGrantPrivileges) validate() error {
	if !exactlyOneValueSet(v.DatabasePrivileges, v.SchemaPrivileges, v.SchemaObjectPrivileges) {
		return fmt.Errorf("exactly one of DatabasePrivileges, SchemaPrivileges, or SchemaObjectPrivileges must be set")
	}
	if valueSet(v.DatabasePrivileges) {
		allowedPrivileges := []AccountObjectPrivilege{
			AccountObjectPrivilegeCreateSchema,
			AccountObjectPrivilegeModify,
			AccountObjectPrivilegeMonitor,
			AccountObjectPrivilegeUsage,
		}
		for _, p := range v.DatabasePrivileges {
			if !slices.Contains(allowedPrivileges, p) {
				return fmt.Errorf("privilege %s is not allowed", p.String())
			}
		}
	}
	return nil
}

func (v *DatabaseRoleGrantOn) validate() error {
	if !exactlyOneValueSet(v.Database, v.Schema, v.SchemaObject) {
		return fmt.Errorf("exactly one of Database, Schema, or SchemaObject must be set")
	}
	if valueSet(v.Schema) {
		if err := v.Schema.validate(); err != nil {
			return err
		}
	}
	if valueSet(v.SchemaObject) {
		if err := v.SchemaObject.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (opts *RevokePrivilegesFromDatabaseRoleOptions) validate() error {
	if !valueSet(opts.privileges) {
		return fmt.Errorf("privileges must be set")
	}
	if err := opts.privileges.validate(); err != nil {
		return err
	}
	if !valueSet(opts.on) {
		return fmt.Errorf("on must be set")
	}
	if err := opts.on.validate(); err != nil {
		return err
	}
	if !validObjectidentifier(opts.databaseRole) {
		return ErrInvalidObjectIdentifier
	}
	if everyValueSet(opts.Restrict, opts.Cascade) {
		return fmt.Errorf("either Restrict or Cascade can be set, or neither but not both")
	}
	return nil
}

func (opts *grantPrivilegeToShareOptions) validate() error {
	if !validObjectidentifier(opts.to) {
		return ErrInvalidObjectIdentifier
	}
	if !valueSet(opts.On) || opts.privilege == "" {
		return fmt.Errorf("on and privilege are required")
	}
	if err := opts.On.validate(); err != nil {
		return err
	}
	return nil
}

func (v *GrantPrivilegeToShareOn) validate() error {
	if !exactlyOneValueSet(v.Database, v.Schema, v.Function, v.Table, v.View) {
		return fmt.Errorf("only one of database, schema, function, table, or view can be set")
	}
	if valueSet(v.Table) {
		if err := v.Table.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (v *OnTable) validate() error {
	if !exactlyOneValueSet(v.Name, v.AllInSchema) {
		return fmt.Errorf("only one of name or allInSchema can be set")
	}
	return nil
}

func (opts *revokePrivilegeFromShareOptions) validate() error {
	if !validObjectidentifier(opts.from) {
		return ErrInvalidObjectIdentifier
	}
	if !valueSet(opts.On) || opts.privilege == "" {
		return fmt.Errorf("on and privilege are required")
	}
	if !exactlyOneValueSet(opts.On.Database, opts.On.Schema, opts.On.Table, opts.On.View) {
		return fmt.Errorf("only one of database, schema, function, table, or view can be set")
	}

	if err := opts.On.validate(); err != nil {
		return err
	}

	return nil
}

func (v *RevokePrivilegeFromShareOn) validate() error {
	if !exactlyOneValueSet(v.Database, v.Schema, v.Table, v.View) {
		return fmt.Errorf("only one of database, schema, table, or view can be set")
	}
	if valueSet(v.Table) {
		return v.Table.validate()
	}
	if valueSet(v.View) {
		return v.View.validate()
	}
	return nil
}

func (v *OnView) validate() error {
	if !exactlyOneValueSet(v.Name, v.AllInSchema) {
		return fmt.Errorf("only one of name or allInSchema can be set")
	}
	return nil
}

// TODO: add validations for ShowGrantsOn, ShowGrantsTo, ShowGrantsOf and ShowGrantsIn
func (opts *ShowGrantOptions) validate() error {
	if everyValueNil(opts.On, opts.To, opts.Of, opts.In) {
		return nil
	}
	if !exactlyOneValueSet(opts.On, opts.To, opts.Of, opts.In) {
		return fmt.Errorf("only one of [on, to, of, in] can be set")
	}
	return nil
}
