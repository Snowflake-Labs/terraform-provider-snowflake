package snowflake

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
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
	 Drop(ctx context.Context, opts PasswordPolicyDropOptions) error
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



func ddlClausesForObject(objectType ObjectType, s interface{}) ([]ddlClause, error) {
	clauses := []ddlClause{}
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", v.Kind())
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if !value.CanInterface() {
			continue
		}
		if field.Type.Kind() == reflect.Struct {
			innerClauses, err := ddlClausesForObject(objectType, value.Interface())
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, innerClauses...)
			continue
		}
		if field.Type.Kind() != reflect.String {
			return nil, fmt.Errorf("expected string, got %s", field.Type.Kind())
		}
		if field.Tag.Get("ddl") == "" {
			continue
		}
		tagParts := strings.Split(field.Tag.Get("ddl"), ",")
		ddlType := tagParts[0]
		switch ddlType {
	

		case "identifier":
			clauses = append(clauses, ddlClauseIdentifier{
				objectType: objectType,
				name : 	value.Interface().(string),
			})
			
		case "keyword":
			useKeyword := value.Interface().(bool)
			if !useKeyword {
				continue
			}
			if len(tagParts) != 2 {
				return nil, fmt.Errorf("expected 2 parts, got %d", len(tagParts))
			}
			clauses = append(clauses, ddlClauseKeyword(tagParts[1]))
		case "parameter":
			if len(tagParts) != 2 {
				return nil, fmt.Errorf("expected 2 parts, got %d", len(tagParts))
			}
			clause := ddlClauseParameter{
				key:   tagParts[1],
				value: value.Interface(),
			}
			clauses = append(clauses, clause)
		default:
			return nil, fmt.Errorf("unknown ddl type %s", ddlType)
		}
	}
	return clauses, nil
}

type PasswordPolicyCreateOptions struct {
	OrReplace   bool   `ddl:"keyword,OR REPLACE"`
	Name        string `ddl:"identifier"`
	IfNotExists bool   `ddl:"keyword,IF NOT EXISTS"`

	PasswordMinLength         int `ddl:"param,PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         int `ddl:"param,PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars int `ddl:"param,PASSWORD_MIN_UPPERCASE_CHARS"`
	PasswordMinLowerCaseChars int `ddl:"param,PASSWORD_MIN_LOWERCASE_CHARS"`
	PasswordMinNumericChars   int `ddl:"param,PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   int `ddl:"param,PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMaxAgeDays        int `ddl:"param,PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        int `ddl:"param,PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   int `ddl:"param,PASSWORD_LOCKOUT_TIME_MINS"`

	Comment string `ddl:"parameter,COMMENT"`
}

type PasswordPolicyAlterOptions struct {
	PasswordMinLength         *int `ddl:"param,PASSWORD_MIN_LENGTH"`
	PasswordMaxLength         int  `ddl:"param,PASSWORD_MAX_LENGTH"`
	PasswordMinUpperCaseChars int  `ddl:"param,PASSWORD_MIN_UPPERCASE_CHARS"`
	PasswordMinLowerCaseChars int  `ddl:"param,PASSWORD_MIN_LOWERCASE_CHARS"`
	PasswordMinNumericChars   int  `ddl:"param,PASSWORD_MIN_NUMERIC_CHARS"`
	PasswordMinSpecialChars   int  `ddl:"param,PASSWORD_MIN_SPECIAL_CHARS"`
	PasswordMaxAgeDays        int  `ddl:"param,PASSWORD_MAX_AGE_DAYS"`
	PasswordMaxRetries        int  `ddl:"param,PASSWORD_MAX_RETRIES"`
	PasswordLockoutTimeMins   int  `ddl:"param,PASSWORD_LOCKOUT_TIME_MINS"`
}

func (opts *PasswordPolicyCreateOptions) validate() error {
	if opts.Name == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func (v *passwordPolicies) Create(ctx context.Context, opts *PasswordPolicyCreateOptions) error {
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate create options: %w", err)
	}
	ddlClauses, err := ddlClausesForObject(ObjectTypePasswordPolicy, opts)
	if err != nil {
		return  err
	}
	stmt := v.client.sql(sqlOperationCreate, ddlClauses...)
	_, err = v.client.exec(ctx, stmt)
	return err
}

type PasswordPolicyDropOptions struct {
	Name     string `ddl:"identifier"`
	IfExists bool   `ddl:"keyword,IF EXISTS"`
}

func (opts *PasswordPolicyDropOptions) validate() error {
	if opts.Name == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func (v *passwordPolicies) Drop(ctx context.Context, opts *PasswordPolicyDropOptions) error {
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate drop options: %w", err)
	}
	ddlClauses, err := ddlClausesForObject(ObjectTypePasswordPolicy, opts)
	if err != nil {
		return  err
	}
	stmt := v.client.sql(sqlOperationDrop, ddlClauses...)
	_, err = v.client.exec(ctx, stmt)
	return err
}

// PasswordPolicyShowOptions represents the options for listing password policies.
type PasswordPolicyShowOptions struct {
	// Required: Filters the command output by object name
	Pattern *string `ddl:"param,LIKE"`

	In *struct  {
		// Optional: Returns records for the specified database
		Account *bool `ddl:"keyword,ACCOUNT"`
		Database *string `ddl:"command,DATABASE"`
		Schema *string `ddl:"command,SCHEMA"`
	}

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

	ddlClauses, err := ddlClausesForObject(ObjectTypePasswordPolicies, opts)
	if err != nil {
		return nil, err
	}
	stmt := v.client.sql(sqlOperationShow, ddlClauses...)
	rows, err := v.client.query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwordPolicies []*PasswordPolicy
	for rows.Next() {
		var passwordPolicy PasswordPolicy
		if err := rows.Scan(
			&passwordPolicy.Name,
			&passwordPolicy.Owner,
			&passwordPolicy.PasswordMinLength,
			&passwordPolicy.PasswordMaxLength,
			&passwordPolicy.PasswordMinUpperCaseChars,
			&passwordPolicy.PasswordMinLowerCaseChars,
			&passwordPolicy.PasswordMinNumericChars,
			&passwordPolicy.PasswordMinSpecialChars,
			&passwordPolicy.PasswordMaxAgeDays,
			&passwordPolicy.PasswordMaxRetries,
			&passwordPolicy.PasswordLockoutTimeMins,
		); err != nil {
			return nil, err
		}
		passwordPolicies = append(passwordPolicies, &passwordPolicy)
	}
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
