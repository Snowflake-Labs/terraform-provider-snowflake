package sdk

import (
	"errors"
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
	_ validatable = new(GrantOwnershipOptions)
	_ validatable = new(ShowGrantOptions)
)

func (opts *GrantPrivilegesToAccountRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !valueSet(opts.privileges) {
		errs = append(errs, errNotSet("GrantPrivilegesToAccountRoleOptions", "privileges"))
	} else {
		if err := opts.privileges.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !valueSet(opts.on) {
		errs = append(errs, errNotSet("GrantPrivilegesToAccountRoleOptions", "on"))
	} else {
		if err := opts.on.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *AccountRoleGrantPrivileges) validate() error {
	if !exactlyOneValueSet(v.AllPrivileges, v.GlobalPrivileges, v.AccountObjectPrivileges, v.SchemaPrivileges, v.SchemaObjectPrivileges) {
		return errExactlyOneOf("AccountRoleGrantPrivileges", "AllPrivileges", "GlobalPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges")
	}
	return nil
}

func (v *AccountRoleGrantOn) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.Account, v.AccountObject, v.Schema, v.SchemaObject) {
		errs = append(errs, errExactlyOneOf("AccountRoleGrantOn", "Account", "AccountObject", "Schema", "SchemaObject"))
	}
	if valueSet(v.AccountObject) {
		if err := v.AccountObject.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.Schema) {
		if err := v.Schema.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.SchemaObject) {
		if err := v.SchemaObject.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *GrantOnAccountObject) validate() error {
	if !exactlyOneValueSet(v.User, v.ResourceMonitor, v.Warehouse, v.Database, v.Integration, v.FailoverGroup, v.ReplicationGroup) {
		return errExactlyOneOf("GrantOnAccountObject", "User", "ResourceMonitor", "Warehouse", "Database", "Integration", "FailoverGroup", "ReplicationGroup")
	}
	return nil
}

func (v *GrantOnSchema) validate() error {
	if !exactlyOneValueSet(v.Schema, v.AllSchemasInDatabase, v.FutureSchemasInDatabase) {
		return errExactlyOneOf("GrantOnSchema", "Schema", "AllSchemasInDatabase", "FutureSchemasInDatabase")
	}
	return nil
}

func (v *GrantOnSchemaObject) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.SchemaObject, v.All, v.Future) {
		errs = append(errs, errExactlyOneOf("GrantOnSchemaObject", "SchemaObject", "All", "Future"))
	}
	if valueSet(v.All) {
		if err := v.All.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.Future) {
		if err := v.Future.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *GrantOnSchemaObjectIn) validate() error {
	if !exactlyOneValueSet(v.InDatabase, v.InSchema) {
		return errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema")
	}
	return nil
}

func (opts *RevokePrivilegesFromAccountRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !valueSet(opts.privileges) {
		errs = append(errs, errNotSet("RevokePrivilegesFromAccountRoleOptions", "privileges"))
	} else {
		if err := opts.privileges.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !valueSet(opts.on) {
		errs = append(errs, errNotSet("RevokePrivilegesFromAccountRoleOptions", "on"))
	} else {
		if err := opts.on.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !ValidObjectIdentifier(opts.accountRole) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.Restrict, opts.Cascade) {
		errs = append(errs, errOneOf("RevokePrivilegesFromAccountRoleOptions", "Restrict", "Cascade"))
	}
	return errors.Join(errs...)
}

func (opts *GrantPrivilegesToDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !valueSet(opts.privileges) {
		errs = append(errs, errNotSet("GrantPrivilegesToDatabaseRoleOptions", "privileges"))
	} else {
		if err := opts.privileges.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !valueSet(opts.on) {
		errs = append(errs, errNotSet("GrantPrivilegesToDatabaseRoleOptions", "on"))
	} else {
		if err := opts.on.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *DatabaseRoleGrantPrivileges) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.DatabasePrivileges, v.SchemaPrivileges, v.SchemaObjectPrivileges) {
		errs = append(errs, errExactlyOneOf("DatabaseRoleGrantPrivileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
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
				errs = append(errs, fmt.Errorf("privilege %s is not allowed", p.String()))
			}
		}
	}
	return errors.Join(errs...)
}

func (v *DatabaseRoleGrantOn) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.Database, v.Schema, v.SchemaObject) {
		errs = append(errs, errExactlyOneOf("DatabaseRoleGrantOn", "Database", "Schema", "SchemaObject"))
	}
	if valueSet(v.Schema) {
		if err := v.Schema.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.SchemaObject) {
		if err := v.SchemaObject.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (opts *RevokePrivilegesFromDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !valueSet(opts.privileges) {
		errs = append(errs, errNotSet("RevokePrivilegesFromDatabaseRoleOptions", "privileges"))
	} else {
		if err := opts.privileges.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !valueSet(opts.on) {
		errs = append(errs, errNotSet("RevokePrivilegesFromDatabaseRoleOptions", "on"))
	} else {
		if err := opts.on.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !ValidObjectIdentifier(opts.databaseRole) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.Restrict, opts.Cascade) {
		errs = append(errs, errOneOf("RevokePrivilegesFromDatabaseRoleOptions", "Restrict", "Cascade"))
	}
	return errors.Join(errs...)
}

func (opts *grantPrivilegeToShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.to) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.On) || opts.privilege == "" {
		errs = append(errs, fmt.Errorf("on and privilege are required"))
	}
	if valueSet(opts.On) {
		if err := opts.On.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *GrantPrivilegeToShareOn) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.Database, v.Schema, v.Function, v.Table, v.View) {
		errs = append(errs, errExactlyOneOf("GrantPrivilegeToShareOn", "Database", "Schema", "Function", "Table", "View"))
	}
	if valueSet(v.Table) {
		if err := v.Table.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *OnTable) validate() error {
	if !exactlyOneValueSet(v.Name, v.AllInSchema) {
		return errExactlyOneOf("OnTable", "Name", "AllInSchema")
	}
	return nil
}

func (opts *revokePrivilegeFromShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.from) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.On) || opts.privilege == "" {
		errs = append(errs, errNotSet("revokePrivilegeFromShareOptions", "On", "privilege"))
	}
	if !exactlyOneValueSet(opts.On.Database, opts.On.Schema, opts.On.Table, opts.On.View) {
		errs = append(errs, errExactlyOneOf("revokePrivilegeFromShareOptions", "On.Database", "On.Schema", "On.Table", "On.View"))
	}
	if valueSet(opts.On) {
		if err := opts.On.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *RevokePrivilegeFromShareOn) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.Database, v.Schema, v.Table, v.View) {
		errs = append(errs, errExactlyOneOf("RevokePrivilegeFromShareOn", "Database", "Schema", "Table", "View"))
	}
	if valueSet(v.Table) {
		if err := v.Table.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.View) {
		if err := v.View.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *OnView) validate() error {
	if !exactlyOneValueSet(v.Name, v.AllInSchema) {
		return errExactlyOneOf("OnView", "Name", "AllInSchema")
	}
	return nil
}

func (opts *GrantOwnershipOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if valueSet(opts.On) {
		if err := opts.On.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.To) {
		if err := opts.To.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *OwnershipGrantOn) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.Object, v.All, v.Future) {
		errs = append(errs, errExactlyOneOf("OwnershipGrantOn", "Object", "AllIn", "Future"))
	}
	if valueSet(v.All) {
		if err := v.All.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.Future) {
		if err := v.Future.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *OwnershipGrantTo) validate() error {
	if !exactlyOneValueSet(v.DatabaseRoleName, v.AccountRoleName) {
		return errExactlyOneOf("OwnershipGrantTo", "databaseRoleName", "accountRoleName")
	}
	return nil
}

// TODO: add validations for ShowGrantsOn, ShowGrantsTo, ShowGrantsOf and ShowGrantsIn
func (opts *ShowGrantOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if everyValueNil(opts.On, opts.To, opts.Of, opts.In) {
		errs = append(errs, errExactlyOneOf("ShowGrantOptions", "On", "To", "Of", "In"))
	}
	if moreThanOneValueSet(opts.On, opts.To, opts.Of, opts.In) {
		errs = append(errs, errOneOf("ShowGrantOptions", "On", "To", "Of", "In"))
	}
	return errors.Join(errs...)
}
