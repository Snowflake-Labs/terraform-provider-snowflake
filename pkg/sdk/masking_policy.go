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
	Create(ctx context.Context, id SchemaObjectIdentifier, signature []TableColumnSignature, returns DataType, expression string, opts *CreateMaskingPolicyOptions) error
	// Alter modifies an existing masking policy.
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterMaskingPolicyOptions) error
	// Drop removes a masking policy.
	Drop(ctx context.Context, id SchemaObjectIdentifier) error
	// Show returns a list of masking policies.
	Show(ctx context.Context, opts *ShowMaskingPolicyOptions) ([]*MaskingPolicy, error)
	// ShowByID returns a masking policy by ID
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicy, error)
	// Describe returns the details of a masking policy.
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicyDetails, error)
}

// maskingPolicies implements MaskingPolicies.
type maskingPolicies struct {
	client *Client
}

type CreateMaskingPolicyOptions struct {
	create        bool                   `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace     *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists   *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`

	// required
	signature []TableColumnSignature `ddl:"keyword,parentheses" sql:"AS"`
	returns   DataType               `ddl:"parameter,no_equals" sql:"RETURNS"`
	body      string                 `ddl:"parameter,no_equals" sql:"->"`

	// optional
	Comment             *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExemptOtherPolicies *bool   `ddl:"parameter" sql:"EXEMPT_OTHER_POLICIES"`
}

func (opts *CreateMaskingPolicyOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return errors.New("invalid object identifier")
	}

	return nil
}

func (v *maskingPolicies) Create(ctx context.Context, id SchemaObjectIdentifier, signature []TableColumnSignature, returns DataType, body string, opts *CreateMaskingPolicyOptions) error {
	if opts == nil {
		opts = &CreateMaskingPolicyOptions{}
	}
	opts.name = id
	opts.signature = signature
	opts.returns = returns
	opts.body = body
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

type AlterMaskingPolicyOptions struct {
	alter         bool                   `ddl:"static" sql:"ALTER"`          //lint:ignore U1000 This is used in the ddl tag
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfExists      *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
	NewName       SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	Set           *MaskingPolicySet      `ddl:"keyword" sql:"SET"`
	Unset         *MaskingPolicyUnset    `ddl:"keyword" sql:"UNSET"`
}

func (opts *AlterMaskingPolicyOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return errors.New("invalid object identifier")
	}

	if everyValueNil(opts.Set, opts.Unset) {
		if !validObjectidentifier(opts.NewName) {
			return ErrInvalidObjectIdentifier
		}
	}

	if !valueSet(opts.NewName) && !exactlyOneValueSet(opts.Set, opts.Unset) {
		return errors.New("cannot use both set and unset")
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

type MaskingPolicySet struct {
	Body    *string          `ddl:"parameter,no_equals" sql:"BODY ->"`
	Tag     []TagAssociation `ddl:"keyword" sql:"TAG"`
	Comment *string          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (v *MaskingPolicySet) validate() error {
	if !exactlyOneValueSet(v.Body, v.Tag, v.Comment) {
		return errors.New("only one parameter can be set at a time")
	}
	return nil
}

type MaskingPolicyUnset struct {
	Tag     []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
	Comment *bool              `ddl:"keyword" sql:"COMMENT"`
}

func (v *MaskingPolicyUnset) validate() error {
	if !exactlyOneValueSet(v.Tag, v.Comment) {
		return errors.New("only one parameter can be unset at a time")
	}
	return nil
}

func (v *maskingPolicies) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterMaskingPolicyOptions) error {
	if opts == nil {
		opts = &AlterMaskingPolicyOptions{}
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

type DropMaskingPolicyOptions struct {
	drop          bool                   `ddl:"static" sql:"DROP"`           //lint:ignore U1000 This is used in the ddl tag
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DropMaskingPolicyOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *maskingPolicies) Drop(ctx context.Context, id SchemaObjectIdentifier) error {
	// masking policy drop does not support [IF EXISTS] so there are no drop options.
	opts := &DropMaskingPolicyOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate drop options: %w", err)
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	if err != nil {
		return err
	}
	return err
}

// ShowMaskingPolicyOptions represents the options for listing masking policies.
type ShowMaskingPolicyOptions struct {
	show            bool  `ddl:"static" sql:"SHOW"`             //lint:ignore U1000 This is used in the ddl tag
	maskingPolicies bool  `ddl:"static" sql:"MASKING POLICIES"` //lint:ignore U1000 This is used in the ddl tag
	Like            *Like `ddl:"keyword" sql:"LIKE"`
	In              *In   `ddl:"keyword" sql:"IN"`
	Limit           *int  `ddl:"parameter,no_equals" sql:"LIMIT"`
}

func (input *ShowMaskingPolicyOptions) validate() error {
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

func (v *MaskingPolicy) ObjectType() ObjectType {
	return ObjectTypeMaskingPolicy
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
func (v *maskingPolicies) Show(ctx context.Context, opts *ShowMaskingPolicyOptions) ([]*MaskingPolicy, error) {
	if opts == nil {
		opts = &ShowMaskingPolicyOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []maskingPolicyDBRow{}

	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]*MaskingPolicy, len(dest))
	for i, row := range dest {
		resultList[i] = row.toMaskingPolicy()
	}

	return resultList, nil
}

func (v *maskingPolicies) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicy, error) {
	maskingPolicies, err := v.Show(ctx, &ShowMaskingPolicyOptions{
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

	for _, maskingPolicy := range maskingPolicies {
		if maskingPolicy.ID().name == id.Name() {
			return maskingPolicy, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type describeMaskingPolicyOptions struct {
	describe      bool                   `ddl:"static" sql:"DESCRIBE"`       //lint:ignore U1000 This is used in the ddl tag
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

func (v *describeMaskingPolicyOptions) validate() error {
	if !validObjectidentifier(v.name) {
		return ErrInvalidObjectIdentifier
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
	dataType, err := ToDataType(row.ReturnType)
	if err != nil {
		return nil
	}
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
		dType, err := ToDataType(p[1])
		if err != nil {
			continue
		}
		v.Signature = append(v.Signature, TableColumnSignature{
			Name: p[0],
			Type: dType,
		})
	}

	return v
}

func (v *maskingPolicies) Describe(ctx context.Context, id SchemaObjectIdentifier) (*MaskingPolicyDetails, error) {
	opts := &describeMaskingPolicyOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := maskingPolicyDetailsRow{}
	err = v.client.queryOne(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}

	return dest.toMaskingPolicyDetails(), nil
}
