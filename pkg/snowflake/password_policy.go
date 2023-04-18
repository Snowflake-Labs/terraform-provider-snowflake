package snowflake

import (
	"context"
	"database/sql"
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
	Create(ctx context.Context, options PasswordPolicyCreateOptions) (*PasswordPolicy, error)
	// Update attributes of an existing role.
	// Alter(ctx context.Context, role string, options PasswordPolicyAlterOptions) (*Role, error)
	// Drop a role by its name.
	// Drop(ctx context.Context, role string) error
	// Show lists all the roles by pattern.
	// Show(ctx context.Context, options PasswordPolicyShowOptions) ([]*PasswordPolicy, error)
	// Describe an password policy by its name.
	// Describe(ctx context.Context, role string) (*PasswordPolicyDetails, error)
}

// passwordPolicies implements PasswordPolicies
type passwordPolicies struct {
	client *Client
}

// PasswordPolicy represents a Snowflake object.
type PasswordPolicy struct {
	Name      string
	CreatedOn time.Time
	Owner     string
	Comment   string
}

type passwordPolicyDB struct {
	Name      sql.NullString `db:"name"`
	CreatedOn sql.NullTime   `db:"created_on"`
	Owner     sql.NullString `db:"owner"`
	Comment   sql.NullString `db:"comment"`
}

func (v *passwordPolicyDB) toPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		Name:      v.Name.String,
		CreatedOn: v.CreatedOn.Time,
		Owner:     v.Owner.String,
	}
}

type StandardCreateOptions struct {
	// Required: Name of the object to create
	Name        string
	OrReplace   bool
	IfNotExists bool
	Comment     string
}

type PasswordPolicyCreateOptions struct {
	StandardCreateOptions
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

func (opts *PasswordPolicyCreateOptions) validate() error {
	if opts.Name == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func (v *passwordPolicies) Create(ctx context.Context, opts PasswordPolicyCreateOptions) (*PasswordPolicy, error) {
	if err := opts.validate(); err != nil {
		return nil, fmt.Errorf("validate create options: %w", err)
	}
	props := []ddlProperty{}

	v.client.sql(sqlOperationCreate, ObjectTypePasswordPolicy, opts.Name, opts)
	sql := fmt.Sprintf("CREATE %s %s", ObjectTypePasswordPolicy, opts.Name)
	if opts.OrReplace {
		sql += " OR REPLACE"
	}
	if opts.IfNotExists {
		sql += " IF NOT EXISTS"
	}
	if opts.Comment != "" {
		sql += fmt.Sprintf(" COMMENT = '%s'", opts.Comment)
	}
	if opts.PasswordMinLength != 0 {
		sql += fmt.Sprintf(" PASSWORD_MIN_LENGTH = %d", opts.PasswordMinLength)
	}
	if opts.PasswordMaxLength != 0 {
		sql += fmt.Sprintf(" PASSWORD_MAX_LENGTH = %d", opts.PasswordMaxLength)
	}
	if opts.PasswordMinUpperCaseChars != 0 {
		sql += fmt.Sprintf(" PASSWORD_MIN_UPPERCASE_CHARS = %d", opts.PasswordMinUpperCaseChars)
	}
	if opts.PasswordMinLowerCaseChars != 0 {
		sql += fmt.Sprintf(" PASSWORD_MIN_LOWERCASE_CHARS = %d", opts.PasswordMinLowerCaseChars)
	}
	if opts.PasswordMinNumericChars != 0 {
		sql += fmt.Sprintf(" PASSWORD_MIN_NUMERIC_CHARS = %d", opts.PasswordMinNumericChars)
	}
	if opts.PasswordMinSpecialChars != 0 {
		sql += fmt.Sprintf(" PASSWORD_MIN_SPECIAL_CHARS = %d", opts.PasswordMinSpecialChars)
	}
	if opts.PasswordMaxAgeDays != 0 {
		sql += fmt.Sprintf(" PASSWORD_MAX_AGE_DAYS = %d", opts.PasswordMaxAgeDays)
	}
	if opts.PasswordMaxRetries != 0 {
		sql += fmt.Sprintf(" PASSWORD_MAX_RETRIES = %d", opts.PasswordMaxRetries)
	}
	if opts.PasswordLockoutTimeMins != 0 {
		sql += fmt.Sprintf(" PASSWORD_LOCKOUT_TIME_MINS = %d", opts.PasswordLockoutTimeMins)
	}

	_, err := v.client.exec(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("do exec: %w", err)
	}

	return &PasswordPolicy{
		Name: opts.Name,
	}, nil
}

/*
// PasswordPolicyShowOptions represents the options for listing password policies.
type PasswordPolicyShowOptions struct {
	// Required: Filters the command output by object name
	Pattern string

	// Optional: Returns records for the entire account.
	InAccount bool

	// Optional: Returns records for the specified database
	InDatabase string

	// Optional: Returns records for the specified schema
	InSchema string

	// Optional: Limits the maximum number of rows returned
	Limit *int
}

func (opts *PasswordPolicyShowOptions) validate() error {
	if opts.Pattern == "" {
		return errors.New("pattern must not be empty")
	}
	return nil
}

// List all the password policies by pattern.
func (v *passwordPolicies) Show(ctx context.Context, opts PasswordPolicyShowOptions) ([]*PasswordPolicy, error) {
	if err := opts.validate(); err != nil {
		return nil, fmt.Errorf("validate list options: %w", err)
	}

	sql := fmt.Sprintf("SHOW %s LIKE '%s'", ResourceRoles, options.Pattern)
	rows, err := r.client.query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("do query: %w", err)
	}
	defer rows.Close()

	entities := []*Role{}
	for rows.Next() {
		var entity roleEntity
		if err := rows.StructScan(&entity); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		entities = append(entities, entity.toRole())
	}
	return entities, nil
}

// PasswordPolicyShowOptions represents the options for listing password policies.
type PasswordPolicyDropOptions struct {
	Name     string
	IfExists bool
}

func (opts *PasswordPolicyDropOptions) validate() error {
	if opts.Name == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func (v *passwordPolicies) Drop(ctx context.Context, opts PasswordPolicyDropOptions) error {
	sql := v.client.templater.drop(ObjectTypePasswordPolicy, opts.Name, opts.IfExists)
	_, err := v.client.exec(ctx, sql)
	return err
	return nil
}
*/
// PasswordPolicyDetails
type PasswordPolicyDetails struct {
	Name                      string
	Owner                     string
	PasswordMinLength         int
	PasswordMaxLength         int
	PasswordMinUpperCaseChars int
	PasswordMinLowerCaseChars int
	PasswordMinNumericChars   int
	PasswordMinSpecialChars   int
	PasswordMaxAgeDays        int
	PasswordMaxRetries        int
	PasswordLockoutTimeMins   int
	Comment                   string
}
