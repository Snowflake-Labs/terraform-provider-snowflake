package sdk

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	ObjectTypePasswordPolicy   ObjectType = "PASSWORD POLICY"
	ObjectTypePasswordPolicies ObjectType = "PASSWORD POLICIES"
)

// Compile-time proof of interface implementation.
var _ PasswordPolicies = (*passwordPolicies)(nil)

// PasswordPolicies describes all the roles related methods that the
// Snowflake API supports.
type PasswordPolicies interface {
	// Create a new role with the given options.
	Create(ctx context.Context, name string, opts *PasswordPolicyCreateOptions) error
	// Update attributes of an existing role.
	Alter(ctx context.Context, name string, opts *PasswordPolicyAlterOptions) error
	// Drop a role by its name.
	Drop(ctx context.Context, name string, opts *PasswordPolicyDropOptions) error
	// Show lists all the roles by pattern.
	Show(ctx context.Context, opts *PasswordPolicyShowOptions) ([]*PasswordPolicy, error)
	// Describe an password policy by its name.
	Describe(ctx context.Context, name string) (*PasswordPolicyDetails, error)
}

// passwordPolicies implements PasswordPolicies.
type passwordPolicies struct {
	client *Client
}

// PasswordPolicy represents a Snowflake object.
type PasswordPolicy struct {
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

type passwordPolicyDetailsRow struct {
	Property     string `db:"property"`
	Value        string `db:"value"`
	DefaultValue string `db:"default"`
	Description  string `db:"description"`
}

type PasswordPolicyDetails struct {
	Name                      string
	Owner                     string
	Comment                   string
	PasswordMinLength         int
	PasswordMaxLength         int
	PasswordMinUpperCaseChars int
	PasswordMinLowerCaseChars int
	PasswordMinNumericChars   int
	PasswordMinSpecialChars   int
	PasswordMaxAgeDays        int
	PasswordMaxRetries        int
	PasswordLockoutTimeMins   int
}

func NewPasswordPolicyDetails(rows []passwordPolicyDetailsRow) *PasswordPolicyDetails {
	v := &PasswordPolicyDetails{}
	for _, row := range rows {
		switch row.Property {
		case "NAME":
			v.Name = row.Value
		case "OWNER":
			v.Owner = row.Value
		case "COMMENT":
			v.Comment = row.Value
		case "PASSWORD_MIN_LENGTH":
			v.PasswordMinLength = toInt(row.Value)
		case "PASSWORD_MAX_LENGTH":
			v.PasswordMaxLength = toInt(row.Value)
		case "PASSWORD_MIN_UPPER_CASE_CHARS":
			v.PasswordMinUpperCaseChars = toInt(row.Value)
		case "PASSWORD_MIN_LOWER_CASE_CHARS":
			v.PasswordMinLowerCaseChars = toInt(row.Value)
		case "PASSWORD_MIN_NUMERIC_CHARS":
			v.PasswordMinNumericChars = toInt(row.Value)
		case "PASSWORD_MIN_SPECIAL_CHARS":
			v.PasswordMinSpecialChars = toInt(row.Value)
		case "PASSWORD_MAX_AGE_DAYS":
			v.PasswordMaxAgeDays = toInt(row.Value)
		case "PASSWORD_MAX_RETRIES":
			v.PasswordMaxRetries = toInt(row.Value)
		case "PASSWORD_LOCKOUT_TIME_MINS":
			v.PasswordLockoutTimeMins = toInt(row.Value)
		}
	}
	return v
}

type PasswordPolicyCreateOptions struct {
	OrReplace   *bool      `ddl:"keyword" db:"OR REPLACE"`
	objectType  ObjectType `ddl:"object_type"`
	name        string     `ddl:"name"`
	IfNotExists *bool      `ddl:"keyword" db:"IF NOT EXISTS"`

	PasswordMinLength         *int `ddl:"param" db:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *int `ddl:"param" db:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *int `ddl:"param" db:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *int `ddl:"param" db:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *int `ddl:"param" db:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *int `ddl:"param" db:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMaxAgeDays        *int `ddl:"param" db:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *int `ddl:"param" db:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *int `ddl:"param" db:"PASSWORD_LOCKOUT_TIME_MINS"`

	Comment *string `ddl:"param,single_quotes" db:"COMMENT"`
}

func (opts *PasswordPolicyCreateOptions) validate() error {
	if opts.name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (v *passwordPolicies) Create(ctx context.Context, name string, opts *PasswordPolicyCreateOptions) error {
	if opts == nil {
		opts = &PasswordPolicyCreateOptions{}
	}
	opts.name = name
	opts.objectType = ObjectTypePasswordPolicy
	if err := opts.validate(); err != nil {
		return err
	}
	clauses, err := ddlClausesForObject(opts)
	if err != nil {
		return err
	}
	stmt := v.client.sql(sqlOperationCreate, clauses...)
	_, err = v.client.execContext(ctx, stmt)
	return err
}

type PasswordPolicyAlterOptions struct {
	objectType ObjectType              `ddl:"object_type"`
	IfExists   *bool                   `ddl:"keyword" db:"IF EXISTS"`
	name       string                  `ddl:"name"`
	Set        *PasswordPolicyAlterSet `ddl:"keyword" db:"SET"`
	Unset      *PasswordPolicyAlterSet `ddl:"keyword" db:"UNSET"`
}

type PasswordPolicyAlterSet struct {
	Name                      *string `ddl:"param"`
	PasswordMinLength         *int    `ddl:"param" db:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *int    `ddl:"param" db:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *int    `ddl:"param" db:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *int    `ddl:"param" db:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *int    `ddl:"param" db:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *int    `ddl:"param" db:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMaxAgeDays        *int    `ddl:"param" db:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *int    `ddl:"param" db:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *int    `ddl:"param" db:"PASSWORD_LOCKOUT_TIME_MINS"`
	Comment                   *string `ddl:"param,single_quotes" db:"COMMENT"`
}

type PasswordPolicyAlterUnset struct {
	PasswordMinLength         *int    `ddl:"param" db:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *int    `ddl:"param" db:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *int    `ddl:"param" db:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *int    `ddl:"param" db:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *int    `ddl:"param" db:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *int    `ddl:"param" db:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMaxAgeDays        *int    `ddl:"param" db:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *int    `ddl:"param" db:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *int    `ddl:"param" db:"PASSWORD_LOCKOUT_TIME_MINS"`
	Comment                   *string `ddl:"param,single_quotes" db:"COMMENT"`
}

func (opts *PasswordPolicyAlterOptions) validate() error {
	if opts.name == "" {
		return errors.New("name must not be empty")
	}

	return nil
}

func (v *passwordPolicies) Alter(ctx context.Context, name string, opts *PasswordPolicyAlterOptions) error {
	if opts == nil {
		opts = &PasswordPolicyAlterOptions{}
	}
	opts.name = name
	opts.objectType = ObjectTypePasswordPolicy
	if err := opts.validate(); err != nil {
		return err
	}
	ddlClauses, err := ddlClausesForObject(opts)
	if err != nil {
		return err
	}
	stmt := v.client.sql(sqlOperationAlter, ddlClauses...)
	_, err = v.client.execContext(ctx, stmt)
	return err
}

type PasswordPolicyDropOptions struct {
	objectType ObjectType `ddl:"object_type"`
	IfExists   *bool      `ddl:"keyword" db:"IF EXISTS"`
	name       string     `ddl:"name"`
}

func (opts *PasswordPolicyDropOptions) validate() error {
	if opts.name == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func (v *passwordPolicies) Drop(ctx context.Context, name string, opts *PasswordPolicyDropOptions) error {
	if opts == nil {
		opts = &PasswordPolicyDropOptions{}
	}
	opts.name = name
	opts.objectType = ObjectTypePasswordPolicy
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate drop options: %w", err)
	}
	ddlClauses, err := ddlClausesForObject(opts)
	if err != nil {
		return err
	}
	stmt := v.client.sql(sqlOperationDrop, ddlClauses...)
	_, err = v.client.execContext(ctx, stmt)
	return err
}

// PasswordPolicyShowOptions represents the options for listing password policies.
type PasswordPolicyShowOptions struct {
	objectType ObjectType            `ddl:"object_type"`
	Pattern    *string               `ddl:"command_param,single_quotes" db:"LIKE"`
	In         *PasswordPolicyShowIn `ddl:"keyword" db:"IN"`

	// Optional: Limits the maximum number of rows returned.
	Limit *int `ddl:"command_param,single_quotes" db:"LIMIT"`
}

type PasswordPolicyShowIn struct {
	Account  *bool   `ddl:"keyword" db:"ACCOUNT"`
	Database *string `ddl:"command_param,no_quotes" db:"DATABASE"`
	Schema   *string `ddl:"command_param,no_quotes" db:"SCHEMA"`
}

// todo: implement this function to validate that the combination of options is valid.
func (opts *PasswordPolicyShowOptions) validate() error {
	return nil
}

// List all the password policies by pattern.
func (v *passwordPolicies) Show(ctx context.Context, opts *PasswordPolicyShowOptions) ([]*PasswordPolicy, error) {
	if opts == nil {
		opts = &PasswordPolicyShowOptions{}
	}
	opts.objectType = ObjectTypePasswordPolicies
	if err := opts.validate(); err != nil {
		return nil, err
	}
	clauses, err := ddlClausesForObject(opts)
	if err != nil {
		return nil, err
	}
	stmt := v.client.sql(sqlOperationShow, clauses...)
	dest := []PasswordPolicy{}

	err = v.client.selectContext(ctx, &dest, stmt)
	if err != nil {
		return nil, err
	}
	return pSlice(dest), err
}

// PasswordPolicyDetailsOptions represents the options for listing password policies.
type PasswordPolicyDetailsOptions struct {
	objectType ObjectType `ddl:"object_type"`
	name       string     `ddl:"name"`
}

func (opts *PasswordPolicyDetailsOptions) validate() error {
	if opts.name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func (v *passwordPolicies) Describe(ctx context.Context, name string) (*PasswordPolicyDetails, error) {
	opts := &PasswordPolicyDetailsOptions{
		name: name,
	}
	opts.objectType = ObjectTypePasswordPolicy
	if err := opts.validate(); err != nil {
		return nil, err
	}

	clauses, err := ddlClausesForObject(opts)
	if err != nil {
		return nil, err
	}
	stmt := v.client.sql(sqlOperationDescribe, clauses...)
	dest := []passwordPolicyDetailsRow{}
	err = v.client.selectContext(ctx, &dest, stmt)
	if err != nil {
		return nil, err
	}
	return NewPasswordPolicyDetails(dest), nil
}
