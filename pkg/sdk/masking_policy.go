package sdk

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

// Compile-time proof of interface implementation.
var _ MaskingPolicies = (*maskingPolicies)(nil)

// MaskingPolicies describes all the masking policy related methods that the
// Snowflake API supports.
type MaskingPolicies interface {
	// Create creates a new masking policy.
	Create(ctx context.Context, id SchemaObjectIdentifier, signature []TableColumnSignature, returns DataType, expression string, opts *MaskingPolicyCreateOptions) error
	// Alter modifies an existing masking policy.
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *MaskingPolicyAlterOptions) error
	// Drop removes a masking policy.
	Drop(ctx context.Context, id SchemaObjectIdentifier) error
	// Show returns a list of masking policies.
	Show(ctx context.Context, opts *MaskingPolicyShowOptions) ([]*MaskingPolicy, error)
	// Describe returns the details of a masking policy.
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicyDetails, error)
}

// maskingPolicies implements MaskingPolicies.
type maskingPolicies struct {
	client  *Client
	builder *sqlBuilder
}

type MaskingPolicyCreateOptions struct {
	create        bool                   `ddl:"static" db:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace     *bool                  `ddl:"keyword" db:"OR REPLACE"`
	maskingPolicy bool                   `ddl:"static" db:"MASKING POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists   *bool                  `ddl:"keyword" db:"IF NOT EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`

	// required
	signature []TableColumnSignature `ddl:"keyword,parentheses" db:"AS"`
	returns   DataType               `ddl:"parameter,no_equals" db:"RETURNS"`
	body      string                 `ddl:"parameter,no_equals" db:"->"`

	// optional
	Comment             *string `ddl:"parameter,single_quotes" db:"COMMENT"`
	ExemptOtherPolicies *bool   `ddl:"parameter" db:"EXEMPT_OTHER_POLICIES"`
}

func (opts *MaskingPolicyCreateOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (v *maskingPolicies) Create(ctx context.Context, id SchemaObjectIdentifier, signature []TableColumnSignature, returns DataType, body string, opts *MaskingPolicyCreateOptions) error {
	if opts == nil {
		opts = &MaskingPolicyCreateOptions{}
	}
	opts.name = id
	opts.signature = signature
	opts.returns = returns
	opts.body = body
	if err := opts.validate(); err != nil {
		return err
	}
	clauses, err := v.builder.parseStruct(opts)
	if err != nil {
		return err
	}
	stmt := v.builder.sql(clauses...)
	_, err = v.client.exec(ctx, stmt)
	return err
}

type MaskingPolicyAlterOptions struct {
	alter         bool                   `ddl:"static" db:"ALTER"`          //lint:ignore U1000 This is used in the ddl tag
	maskingPolicy bool                   `ddl:"static" db:"MASKING POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfExists      *bool                  `ddl:"keyword" db:"IF EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
	NewName       SchemaObjectIdentifier `ddl:"identifier" db:"RENAME TO"`
	Set           *MaskingPolicySet      `ddl:"keyword" db:"SET"`
	Unset         *MaskingPolicyUnset    `ddl:"keyword" db:"UNSET"`
}

func (opts *MaskingPolicyAlterOptions) validate() error {
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
		if opts.Set.Body != nil {
			count++
		}
		if opts.Set.Tag != nil {
			count++
		}
		if opts.Set.Comment != nil {
			count++
		}
		if count != 1 {
			return errors.New("only one parameter must be set")
		}
	}

	if opts.Unset != nil {
		count := 0
		if opts.Unset.Tag != nil {
			count++
		}
		if opts.Unset.Comment != nil {
			count++
		}
		if count != 1 {
			return errors.New("only one parameter can be unset at a time")
		}
	}

	return nil
}

type MaskingPolicySet struct {
	Body    *string          `ddl:"parameter,no_equals" db:"BODY ->"`
	Tag     []TagAssociation `ddl:"keyword" db:"TAG"`
	Comment *string          `ddl:"parameter,single_quotes" db:"COMMENT"`
}

type MaskingPolicyUnset struct {
	Tag     []ObjectIdentifier `ddl:"keyword" db:"TAG"`
	Comment *bool              `ddl:"keyword" db:"COMMENT"`
}

func (v *maskingPolicies) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *MaskingPolicyAlterOptions) error {
	if opts == nil {
		opts = &MaskingPolicyAlterOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	clauses, err := v.builder.parseStruct(opts)
	if err != nil {
		return err
	}
	stmt := v.builder.sql(clauses...)
	_, err = v.client.exec(ctx, stmt)
	return err
}

type MaskingPolicyDropOptions struct {
	drop          bool                   `ddl:"static" db:"DROP"`           //lint:ignore U1000 This is used in the ddl tag
	maskingPolicy bool                   `ddl:"static" db:"MASKING POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *MaskingPolicyDropOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func (v *maskingPolicies) Drop(ctx context.Context, id SchemaObjectIdentifier) error {
	// masking policy drop does not support [IF EXISTS] so there are no drop options.
	opts := &MaskingPolicyDropOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate drop options: %w", err)
	}
	clauses, err := v.builder.parseStruct(opts)
	if err != nil {
		return err
	}
	stmt := v.builder.sql(clauses...)
	_, err = v.client.exec(ctx, stmt)
	if err != nil {
		return decodeDriverError(err)
	}
	return err
}

// MaskingPolicyShowOptions represents the options for listing masking policies.
type MaskingPolicyShowOptions struct {
	show            bool  `ddl:"static" db:"SHOW"`             //lint:ignore U1000 This is used in the ddl tag
	maskingPolicies bool  `ddl:"static" db:"MASKING POLICIES"` //lint:ignore U1000 This is used in the ddl tag
	Like            *Like `ddl:"keyword" db:"LIKE"`
	In              *In   `ddl:"keyword" db:"IN"`
	Limit           *int  `ddl:"parameter,no_equals" db:"LIMIT"`
}

func (input *MaskingPolicyShowOptions) validate() error {
	return nil
}

// MaskingPolicys is a user friendly result for a CREATE MASKING POLICY query.
type MaskingPolicy struct {
	CreatedOn           time.Time
	Name                string
	DatabaseName        string
	SchemaName          string
	Kind                string
	Owner               string
	Comment             string
	ExemptOtherPolicies bool
}

func (v *MaskingPolicy) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

// maskingPolicyDBRow is used to decode the result of a CREATE MASKING POLICY query.
type maskingPolicyDBRow struct {
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

func (row maskingPolicyDBRow) toMaskingPolicy() *MaskingPolicy {
	exemptOtherPolicies, err := jsonparser.GetBoolean([]byte(row.Options), "EXEMPT_OTHER_POLICIES")
	if err != nil {
		exemptOtherPolicies = false
	}
	return &MaskingPolicy{
		CreatedOn:           row.CreatedOn,
		Name:                row.Name,
		DatabaseName:        row.DatabaseName,
		SchemaName:          row.SchemaName,
		Kind:                row.Kind,
		Owner:               row.Owner,
		Comment:             row.Comment,
		ExemptOtherPolicies: exemptOtherPolicies,
	}
}

// List all the masking policies by pattern.
func (v *maskingPolicies) Show(ctx context.Context, opts *MaskingPolicyShowOptions) ([]*MaskingPolicy, error) {
	if opts == nil {
		opts = &MaskingPolicyShowOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	clauses, err := v.builder.parseStruct(opts)
	if err != nil {
		return nil, err
	}
	stmt := v.builder.sql(clauses...)
	dest := []maskingPolicyDBRow{}

	err = v.client.query(ctx, &dest, stmt)
	if err != nil {
		return nil, decodeDriverError(err)
	}
	resultList := make([]*MaskingPolicy, len(dest))
	for i, row := range dest {
		resultList[i] = row.toMaskingPolicy()
	}

	return resultList, nil
}

type maskingPolicyDescribeOptions struct {
	describe      bool                   `ddl:"static" db:"DESCRIBE"`       //lint:ignore U1000 This is used in the ddl tag
	maskingPolicy bool                   `ddl:"static" db:"MASKING POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

func (v *maskingPolicyDescribeOptions) validate() error {
	if v.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type MaskingPolicyDetails struct {
	Name       string
	Signature  []TableColumnSignature
	ReturnType DataType
	Body       string
}

type maskingPolicyDetailsRow struct {
	Name       string `db:"name"`
	Signature  string `db:"signature"`
	ReturnType string `db:"return_type"`
	Body       string `db:"body"`
}

func (row maskingPolicyDetailsRow) toMaskingPolicyDetails() *MaskingPolicyDetails {
	dataType := DataTypeFromString(row.ReturnType)
	v := &MaskingPolicyDetails{
		Name:       row.Name,
		Signature:  []TableColumnSignature{},
		ReturnType: dataType,
		Body:       row.Body,
	}
	s := strings.Trim(row.Signature, "()")
	parts := strings.Split(s, ",")
	for _, part := range parts {
		p := strings.Split(strings.TrimSpace(part), " ")
		if len(p) != 2 {
			continue
		}
		dType := DataTypeFromString(p[1])
		v.Signature = append(v.Signature, TableColumnSignature{
			Name: p[0],
			Type: dType,
		})
	}

	return v
}

func (v *maskingPolicies) Describe(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicyDetails, error) {
	opts := &maskingPolicyDescribeOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	clauses, err := v.builder.parseStruct(opts)
	if err != nil {
		return nil, err
	}
	stmt := v.builder.sql(clauses...)
	dest := maskingPolicyDetailsRow{}
	err = v.client.queryOne(ctx, &dest, stmt)
	if err != nil {
		return nil, decodeDriverError(err)
	}

	return dest.toMaskingPolicyDetails(), nil
}
