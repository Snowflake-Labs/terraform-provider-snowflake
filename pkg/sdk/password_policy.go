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
	Create(ctx context.Context, name string, opts *PasswordPolicyCreateOptions) error
	// Alter modifies an existing password policy.
	Alter(ctx context.Context, name string, opts *PasswordPolicyAlterOptions) error
	// Drop removes a password policy.
	Drop(ctx context.Context, name string, opts *PasswordPolicyDropOptions) error
	// Show returns a list of password policies.
	Show(ctx context.Context, opts *PasswordPolicyShowOptions) ([]*PasswordPolicyShowResult, error)
	// Describe returns the details of a password policy.
	Describe(ctx context.Context, name string) (*PasswordPolicyDescribeResult, error)
}

// passwordPolicies implements PasswordPolicies.
type passwordPolicies struct {
	client *Client
}

type PasswordPolicyCreateOptions struct {
	OrReplace   *bool `ddl:"keyword" db:"OR REPLACE"`
	IfNotExists *bool

	PasswordMinLength         *int
	PasswordMaxLength         *int
	PasswordMinUpperCaseChars *int
	PasswordMinLowerCaseChars *int
	PasswordMinNumericChars   *int
	PasswordMinSpecialChars   *int
	PasswordMaxAgeDays        *int
	PasswordMaxRetries        *int
	PasswordLockoutTimeMins   *int

	Comment *string
}

type passwordPolicyCreateInput struct {
	create         bool   `ddl:"static" db:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace      *bool  `ddl:"keyword" db:"OR REPLACE"`
	passwordPolicy bool   `ddl:"static" db:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name           string `ddl:"name"`
	IfNotExists    *bool  `ddl:"keyword" db:"IF NOT EXISTS"`

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

func (v *passwordPolicyCreateInput) validate() error {
	if v.name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (v *passwordPolicies) Create(ctx context.Context, name string, opts *PasswordPolicyCreateOptions) error {
	input := &passwordPolicyCreateInput{
		name: name,
	}
	copyFields(opts, input)
	if err := input.validate(); err != nil {
		return err
	}
	clauses, err := v.client.parseStruct(input)
	if err != nil {
		return err
	}
	stmt := v.client.sql(clauses...)
	_, err = v.client.exec(ctx, stmt)
	return err
}

type PasswordPolicyAlterOptions struct {
	IfExists *bool
	Set      *PasswordPolicyAlterSet
	Unset    *PasswordPolicyAlterSet
}

type passwordPolicyAlterInput struct {
	alter          bool                    `ddl:"static" db:"ALTER"`           //lint:ignore U1000 This is used in the ddl tag
	passwordPolicy bool                    `ddl:"static" db:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfExists       *bool                   `ddl:"keyword" db:"IF EXISTS"`
	name           string                  `ddl:"name"`
	Set            *PasswordPolicyAlterSet `ddl:"keyword" db:"SET"`
	Unset          *PasswordPolicyAlterSet `ddl:"keyword" db:"UNSET"`
}

type PasswordPolicyAlterSet struct {
	Name                      *string `ddl:"parameter"`
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

func (v *passwordPolicyAlterInput) validate() error {
	if v.name == "" {
		return errors.New("name must not be empty")
	}

	return nil
}

func (v *passwordPolicies) Alter(ctx context.Context, name string, opts *PasswordPolicyAlterOptions) error {
	input := &passwordPolicyAlterInput{
		name: name,
	}
	copyFields(opts, input)
	if err := input.validate(); err != nil {
		return err
	}
	clauses, err := v.client.parseStruct(input)
	if err != nil {
		return err
	}
	stmt := v.client.sql(clauses...)
	_, err = v.client.exec(ctx, stmt)
	return err
}

type PasswordPolicyDropOptions struct {
	IfExists *bool
}

type PasswordPolicyDropInput struct {
	drop           bool   `ddl:"static" db:"DROP"`            //lint:ignore U1000 This is used in the ddl tag
	passwordPolicy bool   `ddl:"static" db:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfExists       *bool  `ddl:"keyword" db:"IF EXISTS"`
	name           string `ddl:"name"`
}

func (v *PasswordPolicyDropInput) validate() error {
	if v.name == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func (v *passwordPolicies) Drop(ctx context.Context, name string, opts *PasswordPolicyDropOptions) error {
	input := &PasswordPolicyDropInput{
		name: name,
	}
	copyFields(opts, input)
	if err := input.validate(); err != nil {
		return fmt.Errorf("validate drop input: %w", err)
	}
	clauses, err := v.client.parseStruct(input)
	if err != nil {
		return err
	}
	stmt := v.client.sql(clauses...)
	_, err = v.client.exec(ctx, stmt)
	return err
}

// PasswordPolicyShowOptions represents the options for listing password policies.
type PasswordPolicyShowOptions struct {
	Like  *Like
	In    *In
	Limit *int
}

type passwordPolicyShowInput struct {
	show             bool  `ddl:"static" db:"SHOW"`              //lint:ignore U1000 This is used in the ddl tag
	passwordPolicies bool  `ddl:"static" db:"PASSWORD POLICIES"` //lint:ignore U1000 This is used in the ddl tag
	Like             *Like `ddl:"keyword" db:"LIKE"`
	In               *In   `ddl:"keyword" db:"IN"`
	Limit            *int  `ddl:"command,single_quotes" db:"LIMIT"`
}

func (input *passwordPolicyShowInput) validate() error {
	return nil
}

// PasswordPolicyCreateOptions is a user friendly result for a CREATE PASSWORD POLICY query.
type PasswordPolicyShowResult struct {
	CreatedOn    time.Time
	Name         string
	DatabaseName string
	SchemaName   string
	Kind         string
	Owner        string
	Comment      string
}

// passwordPolicyCreateDBRow is used to decode the result of a CREATE PASSWORD POLICY query.
type passwordPolicyShowDBRow struct {
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

func passwordPolicyShowResultFromRow(row *passwordPolicyShowDBRow) *PasswordPolicyShowResult {
	return &PasswordPolicyShowResult{
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
func (v *passwordPolicies) Show(ctx context.Context, opts *PasswordPolicyShowOptions) ([]*PasswordPolicyShowResult, error) {
	input := &passwordPolicyShowInput{}
	copyFields(opts, input)
	if err := input.validate(); err != nil {
		return nil, err
	}
	clauses, err := v.client.parseStruct(input)
	if err != nil {
		return nil, err
	}
	stmt := v.client.sql(clauses...)
	dest := []passwordPolicyShowDBRow{}

	err = v.client.query(ctx, &dest, stmt)
	if err != nil {
		return nil, err
	}
	resultList := make([]*PasswordPolicyShowResult, len(dest))
	for i, row := range dest {
		resultList[i] = passwordPolicyShowResultFromRow(&row)
	}

	return resultList, nil
}

type passwordPolicyDetailsInput struct {
	show           bool   `ddl:"static" db:"DESCRIBE"`        //lint:ignore U1000 This is used in the ddl tag
	passwordPolicy bool   `ddl:"static" db:"PASSWORD POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name           string `ddl:"name"`
}

func (v *passwordPolicyDetailsInput) validate() error {
	if v.name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type PasswordPolicyDescribeResult struct {
	Name                      *DescribeStringProperty
	Owner                     *DescribeStringProperty
	Comment                   *DescribeStringProperty
	PasswordMinLength         *DescribeIntProperty
	PasswordMaxLength         *DescribeIntProperty
	PasswordMinUpperCaseChars *DescribeIntProperty
	PasswordMinLowerCaseChars *DescribeIntProperty
	PasswordMinNumericChars   *DescribeIntProperty
	PasswordMinSpecialChars   *DescribeIntProperty
	PasswordMaxAgeDays        *DescribeIntProperty
	PasswordMaxRetries        *DescribeIntProperty
	PasswordLockoutTimeMins   *DescribeIntProperty
}

func passwordPolicyDescribeResultFromRows(rows []describePropertyRow) *PasswordPolicyDescribeResult {
	v := &PasswordPolicyDescribeResult{}
	for _, row := range rows {
		switch row.Property {
		case "NAME":
			v.Name = row.toDescribeStringProperty()
		case "OWNER":
			v.Owner = row.toDescribeStringProperty()
		case "COMMENT":
			v.Comment = row.toDescribeStringProperty()
		case "PASSWORD_MIN_LENGTH":
			v.PasswordMinLength = row.toDescribeIntProperty()
		case "PASSWORD_MAX_LENGTH":
			v.PasswordMaxLength = row.toDescribeIntProperty()
		case "PASSWORD_MIN_UPPER_CASE_CHARS":
			v.PasswordMinUpperCaseChars = row.toDescribeIntProperty()
		case "PASSWORD_MIN_LOWER_CASE_CHARS":
			v.PasswordMinLowerCaseChars = row.toDescribeIntProperty()
		case "PASSWORD_MIN_NUMERIC_CHARS":
			v.PasswordMinNumericChars = row.toDescribeIntProperty()
		case "PASSWORD_MIN_SPECIAL_CHARS":
			v.PasswordMinSpecialChars = row.toDescribeIntProperty()
		case "PASSWORD_MAX_AGE_DAYS":
			v.PasswordMaxAgeDays = row.toDescribeIntProperty()
		case "PASSWORD_MAX_RETRIES":
			v.PasswordMaxRetries = row.toDescribeIntProperty()
		case "PASSWORD_LOCKOUT_TIME_MINS":
			v.PasswordLockoutTimeMins = row.toDescribeIntProperty()
		}
	}
	return v
}

func (v *passwordPolicies) Describe(ctx context.Context, name string) (*PasswordPolicyDescribeResult, error) {
	input := &passwordPolicyDetailsInput{
		name: name,
	}
	if err := input.validate(); err != nil {
		return nil, err
	}

	clauses, err := v.client.parseStruct(input)
	if err != nil {
		return nil, err
	}
	stmt := v.client.sql(clauses...)
	dest := []describePropertyRow{}
	err = v.client.query(ctx, &dest, stmt)
	if err != nil {
		return nil, err
	}

	return passwordPolicyDescribeResultFromRows(dest), nil
}
