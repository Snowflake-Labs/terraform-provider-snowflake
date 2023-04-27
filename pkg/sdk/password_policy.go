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
	Create(ctx context.Context, id SchemaObjectIdentifier, opts *PasswordPolicyCreateOptions) error
	// Alter modifies an existing password policy.
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *PasswordPolicyAlterOptions) error
	// Drop removes a password policy.
	Drop(ctx context.Context, id SchemaObjectIdentifier, opts *PasswordPolicyDropOptions) error
	// Show returns a list of password policies.
	Show(ctx context.Context, opts *PasswordPolicyShowOptions) ([]*PasswordPolicy, error)
	// Describe returns the details of a password policy.
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*PasswordPolicyDetails, error)
}

// passwordPolicies implements PasswordPolicies.
type passwordPolicies struct {
	client *Client
}

type PasswordPolicyCreateOptions struct {
	create         bool                   `ddl:"static" db:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace      *bool                  `ddl:"keyword" db:"OR REPLACE"`
	passwordPolicy bool                   `ddl:"static" db:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name           SchemaObjectIdentifier `ddl:"identifier"`
	IfNotExists    *bool                  `ddl:"keyword" db:"IF NOT EXISTS"`

	PasswordMinLength         *int `ddl:"parameter" db:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *int `ddl:"parameter" db:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *int `ddl:"parameter" db:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *int `ddl:"parameter" db:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *int `ddl:"parameter" db:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *int `ddl:"parameter" db:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMaxAgeDays        *int `ddl:"parameter" db:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *int `ddl:"parameter" db:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *int `ddl:"parameter" db:"PASSWORD_LOCKOUT_TIME_MINS"`

	Comment *string `ddl:"parameter,single_quotes" db:"COMMENT"`
}

func (opts *PasswordPolicyCreateOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (v *passwordPolicies) Create(ctx context.Context, id SchemaObjectIdentifier, opts *PasswordPolicyCreateOptions) error {
	if opts == nil {
		opts = &PasswordPolicyCreateOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	clauses, err := v.client.parseStruct(opts)
	if err != nil {
		return err
	}
	stmt := v.client.sql(clauses...)
	_, err = v.client.exec(ctx, stmt)
	return err
}

type PasswordPolicyAlterOptions struct {
	alter          bool                      `ddl:"static" db:"ALTER"`           //lint:ignore U1000 This is used in the ddl tag
	passwordPolicy bool                      `ddl:"static" db:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfExists       *bool                     `ddl:"keyword" db:"IF EXISTS"`
	name           SchemaObjectIdentifier    `ddl:"identifier"`
	NewName        SchemaObjectIdentifier    `ddl:"identifier" db:"RENAME TO"`
	Set            *PasswordPolicyAlterSet   `ddl:"keyword" db:"SET"`
	Unset          *PasswordPolicyAlterUnset `ddl:"keyword" db:"UNSET"`
}

func (opts *PasswordPolicyAlterOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return errors.New("name must not be empty")
	}

	if opts.Set == nil && opts.Unset == nil {
		if opts.NewName.FullyQualifiedName() == "" {
			return errors.New("new name must not be empty")
		}
	}

	if opts.Set != nil && opts.Unset != nil {
		return errors.New("cannot set and unset parameters in the same ALTER statement")
	}

	if opts.Set != nil {
		count := 0
		if opts.Set.PasswordMinLength != nil {
			count++
		}
		if opts.Set.PasswordMaxLength != nil {
			count++
		}
		if opts.Set.PasswordMinUpperCaseChars != nil {
			count++
		}
		if opts.Set.PasswordMinLowerCaseChars != nil {
			count++
		}
		if opts.Set.PasswordMinNumericChars != nil {
			count++
		}
		if opts.Set.PasswordMinSpecialChars != nil {
			count++
		}
		if opts.Set.PasswordMaxAgeDays != nil {
			count++
		}
		if opts.Set.PasswordMaxRetries != nil {
			count++
		}
		if opts.Set.PasswordLockoutTimeMins != nil {
			count++
		}
		if count == 0 {
			return errors.New("at least one parameter must be set")
		}
		if count > 1 {
			return errors.New("only one parameter can be set at a time")
		}
	}

	if opts.Unset != nil {
		count := 0
		if opts.Unset.PasswordMinLength != nil {
			count++
		}
		if opts.Unset.PasswordMaxLength != nil {
			count++
		}
		if opts.Unset.PasswordMinUpperCaseChars != nil {
			count++
		}
		if opts.Unset.PasswordMinLowerCaseChars != nil {
			count++
		}
		if opts.Unset.PasswordMinNumericChars != nil {
			count++
		}
		if opts.Unset.PasswordMinSpecialChars != nil {
			count++
		}
		if opts.Unset.PasswordMaxAgeDays != nil {
			count++
		}
		if opts.Unset.PasswordMaxRetries != nil {
			count++
		}
		if opts.Unset.PasswordLockoutTimeMins != nil {
			count++
		}
		if count == 0 {
			return errors.New("at least one parameter must be unset")
		}
		if count > 1 {
			return errors.New("only one parameter can be unset at a time")
		}
	}

	return nil
}

type PasswordPolicyAlterSet struct {
	PasswordMinLength         *int    `ddl:"parameter" db:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *int    `ddl:"parameter" db:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *int    `ddl:"parameter" db:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *int    `ddl:"parameter" db:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *int    `ddl:"parameter" db:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *int    `ddl:"parameter" db:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMaxAgeDays        *int    `ddl:"parameter" db:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *int    `ddl:"parameter" db:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *int    `ddl:"parameter" db:"PASSWORD_LOCKOUT_TIME_MINS"`
	Comment                   *string `ddl:"parameter,single_quotes" db:"COMMENT"`
}

type PasswordPolicyAlterUnset struct {
	PasswordMinLength         *bool `ddl:"keyword" db:"PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         *bool `ddl:"keyword" db:"PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars *bool `ddl:"keyword" db:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	PasswordMinLowerCaseChars *bool `ddl:"keyword" db:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	PasswordMinNumericChars   *bool `ddl:"keyword" db:"PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   *bool `ddl:"keyword" db:"PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMaxAgeDays        *bool `ddl:"keyword" db:"PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        *bool `ddl:"keyword" db:"PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   *bool `ddl:"keyword" db:"PASSWORD_LOCKOUT_TIME_MINS"`
	Comment                   *bool `ddl:"keyword" db:"COMMENT"`
}

func (v *passwordPolicies) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *PasswordPolicyAlterOptions) error {
	if opts == nil {
		opts = &PasswordPolicyAlterOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	clauses, err := v.client.parseStruct(opts)
	if err != nil {
		return err
	}
	stmt := v.client.sql(clauses...)
	_, err = v.client.exec(ctx, stmt)
	return err
}

type PasswordPolicyDropOptions struct {
	drop           bool                   `ddl:"static" db:"DROP"`            //lint:ignore U1000 This is used in the ddl tag
	passwordPolicy bool                   `ddl:"static" db:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfExists       *bool                  `ddl:"keyword" db:"IF EXISTS"`
	name           SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *PasswordPolicyDropOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func (v *passwordPolicies) Drop(ctx context.Context, id SchemaObjectIdentifier, opts *PasswordPolicyDropOptions) error {
	if opts == nil {
		opts = &PasswordPolicyDropOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate drop options: %w", err)
	}
	clauses, err := v.client.parseStruct(opts)
	if err != nil {
		return err
	}
	stmt := v.client.sql(clauses...)
	_, err = v.client.exec(ctx, stmt)
	if err != nil {
		return decodeError(err)
	}
	return err
}

// PasswordPolicyShowOptions represents the options for listing password policies.
type PasswordPolicyShowOptions struct {
	show             bool  `ddl:"static" db:"SHOW"`              //lint:ignore U1000 This is used in the ddl tag
	passwordPolicies bool  `ddl:"static" db:"PASSWORD POLICIES"` //lint:ignore U1000 This is used in the ddl tag
	Like             *Like `ddl:"keyword" db:"LIKE"`
	In               *In   `ddl:"keyword" db:"IN"`
	Limit            *int  `ddl:"command,no_quotes" db:"LIMIT"`
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

func (v *PasswordPolicy) Identifier() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
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

func passwordPolicyFromRow(row passwordPolicyDBRow) *PasswordPolicy {
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
	clauses, err := v.client.parseStruct(opts)
	if err != nil {
		return nil, err
	}
	stmt := v.client.sql(clauses...)
	dest := []passwordPolicyDBRow{}

	err = v.client.query(ctx, &dest, stmt)
	if err != nil {
		return nil, decodeError(err)
	}
	resultList := make([]*PasswordPolicy, len(dest))
	for i, row := range dest {
		resultList[i] = passwordPolicyFromRow(row)
	}

	return resultList, nil
}

type passwordPolicyDetailsOptions struct {
	show           bool                   `ddl:"static" db:"DESCRIBE"`        //lint:ignore U1000 This is used in the ddl tag
	passwordPolicy bool                   `ddl:"static" db:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name           SchemaObjectIdentifier `ddl:"identifier"`
}

func (v *passwordPolicyDetailsOptions) validate() error {
	if v.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
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
	opts := &passwordPolicyDetailsOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	clauses, err := v.client.parseStruct(opts)
	if err != nil {
		return nil, err
	}
	stmt := v.client.sql(clauses...)
	dest := []propertyRow{}
	err = v.client.query(ctx, &dest, stmt)
	if err != nil {
		return nil, decodeError(err)
	}

	return passwordPolicyDetailsFromRows(dest), nil
}
