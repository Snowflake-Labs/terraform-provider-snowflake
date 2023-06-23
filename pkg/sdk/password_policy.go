package sdk

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Compile-time proof of interface implementation.
var _ PasswordPolicies = (*passwordPolicies)(nil)

// PasswordPolicies describes all the password policy related methods that the
// Snowflake API supports.
type PasswordPolicies interface {
	// Create creates a new password policy.
	Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreatePasswordPolicyOptions) error
	// Alter modifies an existing password policy.
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterPasswordPolicyOptions) error
	// Drop removes a password policy.
	Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropPasswordPolicyOptions) error
	// Show returns a list of password policies.
	Show(ctx context.Context, opts *PasswordPolicyShowOptions) ([]*PasswordPolicy, error)
	// ShowByID returns a password policy by ID.
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicy, error)
	// Describe returns the details of a password policy.
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicyDetails, error)
}

// passwordPolicies implements PasswordPolicies.
type passwordPolicies struct {
	client *Client
}

type CreatePasswordPolicyOptions struct {
	create         bool                   `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace      *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	passwordPolicy bool                   `ddl:"static" sql:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists    *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name           SchemaObjectIdentifier `ddl:"identifier"`

	PasswordMinLength         *int `ddl:"parameter" sql:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *int `ddl:"parameter" sql:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *int `ddl:"parameter" sql:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *int `ddl:"parameter" sql:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *int `ddl:"parameter" sql:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *int `ddl:"parameter" sql:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMaxAgeDays        *int `ddl:"parameter" sql:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *int `ddl:"parameter" sql:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *int `ddl:"parameter" sql:"PASSWORD_LOCKOUT_TIME_MINS"`

	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreatePasswordPolicyOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}

	return nil
}

func (v *passwordPolicies) Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreatePasswordPolicyOptions) error {
	if opts == nil {
		opts = &CreatePasswordPolicyOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

type AlterPasswordPolicyOptions struct {
	alter          bool                   `ddl:"static" sql:"ALTER"`           //lint:ignore U1000 This is used in the ddl tag
	passwordPolicy bool                   `ddl:"static" sql:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfExists       *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name           SchemaObjectIdentifier `ddl:"identifier"`
	NewName        SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	Set            *PasswordPolicySet     `ddl:"keyword" sql:"SET"`
	Unset          *PasswordPolicyUnset   `ddl:"keyword" sql:"UNSET"`
}

func (opts *AlterPasswordPolicyOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}

	if everyValueNil(opts.Set, opts.Unset) {
		if !validObjectidentifier(opts.NewName) {
			return ErrInvalidObjectIdentifier
		}
	}

	if !valueSet(opts.NewName) && !exactlyOneValueSet(opts.Set, opts.Unset) {
		return errors.New("cannot set and unset parameters in the same ALTER statement")
	}
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			return err
		}
	}

	if valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			return err
		}
	}

	return nil
}

type PasswordPolicySet struct {
	PasswordMinLength         *int    `ddl:"parameter" sql:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *int    `ddl:"parameter" sql:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *int    `ddl:"parameter" sql:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *int    `ddl:"parameter" sql:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *int    `ddl:"parameter" sql:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *int    `ddl:"parameter" sql:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMaxAgeDays        *int    `ddl:"parameter" sql:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *int    `ddl:"parameter" sql:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *int    `ddl:"parameter" sql:"PASSWORD_LOCKOUT_TIME_MINS"`
	Comment                   *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (v *PasswordPolicySet) validate() error {
	if everyValueNil(
		v.PasswordMinLength,
		v.PasswordMaxLength,
		v.PasswordMinUpperCaseChars,
		v.PasswordMinLowerCaseChars,
		v.PasswordMinNumericChars,
		v.PasswordMinSpecialChars,
		v.PasswordMaxAgeDays,
		v.PasswordMaxRetries,
		v.PasswordLockoutTimeMins,
		v.Comment) {
		return errors.New("must set at least one parameter")
	}
	return nil
}

type PasswordPolicyUnset struct {
	PasswordMinLength         *bool `ddl:"keyword" sql:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *bool `ddl:"keyword" sql:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *bool `ddl:"keyword" sql:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *bool `ddl:"keyword" sql:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *bool `ddl:"keyword" sql:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *bool `ddl:"keyword" sql:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMaxAgeDays        *bool `ddl:"keyword" sql:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *bool `ddl:"keyword" sql:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *bool `ddl:"keyword" sql:"PASSWORD_LOCKOUT_TIME_MINS"`
	Comment                   *bool `ddl:"keyword" sql:"COMMENT"`
}

func (v *PasswordPolicyUnset) validate() error {
	if everyValueNil(
		v.PasswordMinLength,
		v.PasswordMaxLength,
		v.PasswordMinUpperCaseChars,
		v.PasswordMinLowerCaseChars,
		v.PasswordMinNumericChars,
		v.PasswordMinSpecialChars,
		v.PasswordMaxAgeDays,
		v.PasswordMaxRetries,
		v.PasswordLockoutTimeMins,
		v.Comment) {
		return errors.New("must unset at least one parameter")
	}
	if !exactlyOneValueSet(
		v.PasswordMinLength,
		v.PasswordMaxLength,
		v.PasswordMinUpperCaseChars,
		v.PasswordMinLowerCaseChars,
		v.PasswordMinNumericChars,
		v.PasswordMinSpecialChars,
		v.PasswordMaxAgeDays,
		v.PasswordMaxRetries,
		v.PasswordLockoutTimeMins,
		v.Comment) {
		return errors.New("cannot unset more than one parameter in the same ALTER statement")
	}
	return nil
}

func (v *passwordPolicies) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterPasswordPolicyOptions) error {
	if opts == nil {
		opts = &AlterPasswordPolicyOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

type DropPasswordPolicyOptions struct {
	drop           bool                   `ddl:"static" sql:"DROP"`            //lint:ignore U1000 This is used in the ddl tag
	passwordPolicy bool                   `ddl:"static" sql:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfExists       *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name           SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DropPasswordPolicyOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *passwordPolicies) Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropPasswordPolicyOptions) error {
	if opts == nil {
		opts = &DropPasswordPolicyOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate drop options: %w", err)
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

// PasswordPolicyShowOptions represents the options for listing password policies.
type PasswordPolicyShowOptions struct {
	show             bool  `ddl:"static" sql:"SHOW"`              //lint:ignore U1000 This is used in the ddl tag
	passwordPolicies bool  `ddl:"static" sql:"PASSWORD POLICIES"` //lint:ignore U1000 This is used in the ddl tag
	Like             *Like `ddl:"keyword" sql:"LIKE"`
	In               *In   `ddl:"keyword" sql:"IN"`
	Limit            *int  `ddl:"parameter,no_equals" sql:"LIMIT"`
}

func (input *PasswordPolicyShowOptions) validate() error {
	return nil
}

// PasswordPolicys is a user friendly result for a CREATE PASSWORD POLICY query.
type PasswordPolicy struct {
	CreatedOn    time.Time
	Name         string
	DatabaseName string
	SchemaName   string
	Kind         string
	Owner        string
	Comment      string
}

func (v *PasswordPolicy) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *PasswordPolicy) ObjectType() ObjectType {
	return ObjectTypePasswordPolicy
}

// passwordPolicyDBRow is used to decode the result of a CREATE PASSWORD POLICY query.
type passwordPolicyDBRow struct {
	CreatedOn     time.Time `db:"created_on"`
	Name          string    `db:"name"`
	DatabaseName  string    `db:"database_name"`
	SchemaName    string    `db:"schema_name"`
	Kind          string    `db:"kind"`
	Owner         string    `db:"owner"`
	Comment       string    `db:"comment"`
	OwnerRoleType string    `db:"owner_role_type"`
	Options       string    `db:"options"`
}

func (row passwordPolicyDBRow) toPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		CreatedOn:    row.CreatedOn,
		Name:         row.Name,
		DatabaseName: row.DatabaseName,
		SchemaName:   row.SchemaName,
		Kind:         row.Kind,
		Owner:        row.Owner,
		Comment:      row.Comment,
	}
}

// List all the password policies by pattern.
func (v *passwordPolicies) Show(ctx context.Context, opts *PasswordPolicyShowOptions) ([]*PasswordPolicy, error) {
	if opts == nil {
		opts = &PasswordPolicyShowOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []passwordPolicyDBRow{}

	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]*PasswordPolicy, len(dest))
	for i, row := range dest {
		resultList[i] = row.toPasswordPolicy()
	}

	return resultList, nil
}

func (v *passwordPolicies) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicy, error) {
	passwordPolicies, err := v.Show(ctx, &PasswordPolicyShowOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
		In: &In{
			Schema: NewSchemaIdentifier(id.DatabaseName(), id.SchemaName()),
		},
	})
	if err != nil {
		return nil, err
	}

	for _, passwordPolicy := range passwordPolicies {
		if passwordPolicy.ID().name == id.Name() {
			return passwordPolicy, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type describePasswordPolicyOptions struct {
	describe       bool                   `ddl:"static" sql:"DESCRIBE"`        //lint:ignore U1000 This is used in the ddl tag
	passwordPolicy bool                   `ddl:"static" sql:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name           SchemaObjectIdentifier `ddl:"identifier"`
}

func (v *describePasswordPolicyOptions) validate() error {
	if !validObjectidentifier(v.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

type PasswordPolicyDetails struct {
	Name                      *StringProperty
	Owner                     *StringProperty
	Comment                   *StringProperty
	PasswordMinLength         *IntProperty
	PasswordMaxLength         *IntProperty
	PasswordMinUpperCaseChars *IntProperty
	PasswordMinLowerCaseChars *IntProperty
	PasswordMinNumericChars   *IntProperty
	PasswordMinSpecialChars   *IntProperty
	PasswordMaxAgeDays        *IntProperty
	PasswordMaxRetries        *IntProperty
	PasswordLockoutTimeMins   *IntProperty
}

func passwordPolicyDetailsFromRows(rows []propertyRow) *PasswordPolicyDetails {
	v := &PasswordPolicyDetails{}
	for _, row := range rows {
		switch row.Property {
		case "NAME":
			v.Name = row.toStringProperty()
		case "OWNER":
			v.Owner = row.toStringProperty()
		case "COMMENT":
			v.Comment = row.toStringProperty()
		case "PASSWORD_MIN_LENGTH":
			v.PasswordMinLength = row.toIntProperty()
		case "PASSWORD_MAX_LENGTH":
			v.PasswordMaxLength = row.toIntProperty()
		case "PASSWORD_MIN_UPPER_CASE_CHARS":
			v.PasswordMinUpperCaseChars = row.toIntProperty()
		case "PASSWORD_MIN_LOWER_CASE_CHARS":
			v.PasswordMinLowerCaseChars = row.toIntProperty()
		case "PASSWORD_MIN_NUMERIC_CHARS":
			v.PasswordMinNumericChars = row.toIntProperty()
		case "PASSWORD_MIN_SPECIAL_CHARS":
			v.PasswordMinSpecialChars = row.toIntProperty()
		case "PASSWORD_MAX_AGE_DAYS":
			v.PasswordMaxAgeDays = row.toIntProperty()
		case "PASSWORD_MAX_RETRIES":
			v.PasswordMaxRetries = row.toIntProperty()
		case "PASSWORD_LOCKOUT_TIME_MINS":
			v.PasswordLockoutTimeMins = row.toIntProperty()
		}
	}
	return v
}

func (v *passwordPolicies) Describe(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicyDetails, error) {
	opts := &describePasswordPolicyOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []propertyRow{}
	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}

	return passwordPolicyDetailsFromRows(dest), nil
}
