package sdk

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Shares interface {
	// Create creates a share.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateShareOptions) error
	// Alter modifies an existing share
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterShareOptions) error
	// Drop removes a share.
	Drop(ctx context.Context, id AccountObjectIdentifier) error
	// Show returns a list of shares.
	Show(ctx context.Context, opts *ShowShareOptions) ([]*Share, error)
	// ShowByID returns a share by ID.
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Share, error)
	// Describe returns the details of an outbound share.
	DescribeProvider(ctx context.Context, id AccountObjectIdentifier) (*ShareDetails, error)
	// Describe returns the details of an inbound share.
	DescribeConsumer(ctx context.Context, id ExternalObjectIdentifier) (*ShareDetails, error)
}

var _ Shares = (*shares)(nil)

type shares struct {
	client *Client
}

type ShareKind string

const (
	ShareKindInbound  ShareKind = "INBOUND"
	ShareKindOutbound ShareKind = "OUTBOUND"
)

type Share struct {
	CreatedOn    time.Time
	Kind         ShareKind
	Name         ExternalObjectIdentifier
	DatabaseName AccountObjectIdentifier
	To           []AccountIdentifier
	Owner        string
	Comment      string
}

func (v *Share) ID() AccountObjectIdentifier {
	return v.Name.objectIdentifier.(AccountObjectIdentifier)
}

func (v *Share) ExternalID() ExternalObjectIdentifier {
	return v.Name
}

func (v *Share) ObjectType() ObjectType {
	return ObjectTypeShare
}

type shareRow struct {
	CreatedOn    time.Time `db:"created_on"`
	Kind         string    `db:"kind"`
	Name         string    `db:"name"`
	DatabaseName string    `db:"database_name"`
	To           string    `db:"to"`
	Owner        string    `db:"owner"`
	Comment      string    `db:"comment"`
}

func (r *shareRow) toShare() *Share {
	toAccounts := strings.Split(r.To, ",")
	var to []AccountIdentifier
	if len(toAccounts) != 0 {
		for _, a := range toAccounts {
			if a == "" {
				continue
			}
			parts := strings.Split(a, ".")
			if len(parts) == 1 {
				accountLocator := parts[0]
				to = append(to, NewAccountIdentifierFromAccountLocator(accountLocator))
				continue
			}
			orgName := parts[0]
			accountName := strings.Join(parts[1:], ".")
			to = append(to, NewAccountIdentifier(orgName, accountName))
		}
	}
	return &Share{
		CreatedOn:    r.CreatedOn,
		Kind:         ShareKind(r.Kind),
		Name:         NewExternalObjectIdentifierFromFullyQualifiedName(r.Name),
		DatabaseName: NewAccountObjectIdentifier(r.DatabaseName),
		To:           to,
		Owner:        r.Owner,
		Comment:      r.Comment,
	}
}

type CreateShareOptions struct {
	create    bool                    `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	share     bool                    `ddl:"static" sql:"SHARE"` //lint:ignore U1000 This is used in the ddl tag
	name      AccountObjectIdentifier `ddl:"identifier"`
	Comment   *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreateShareOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return fmt.Errorf("not a valid object identifier: %s", opts.name)
	}
	return nil
}

func (v *shares) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateShareOptions) error {
	if opts == nil {
		opts = &CreateShareOptions{}
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

type shareDropOptions struct {
	drop  bool                    `ddl:"static" sql:"DROP"` //lint:ignore U1000 This is used in the ddl tag
	share bool                    `ddl:"static" sql:"SHARE"`
	name  AccountObjectIdentifier `ddl:"identifier"`
}

func (v *shares) Drop(ctx context.Context, id AccountObjectIdentifier) error {
	opts := &shareDropOptions{
		name: id,
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

type AlterShareOptions struct {
	alter    bool                    `ddl:"static" sql:"ALTER"` //lint:ignore U1000 This is used in the ddl tag
	share    bool                    `ddl:"static" sql:"SHARE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`
	Add      *ShareAdd               `ddl:"keyword" sql:"ADD"`
	Remove   *ShareRemove            `ddl:"keyword" sql:"REMOVE"`
	Set      *ShareSet               `ddl:"keyword" sql:"SET"`
	Unset    *ShareUnset             `ddl:"keyword" sql:"UNSET"`
}

func (opts *AlterShareOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return fmt.Errorf("not a valid object identifier: %s", opts.name)
	}
	if ok := exactlyOneValueSet(opts.Add, opts.Remove, opts.Set, opts.Unset); !ok {
		return fmt.Errorf("exactly one of add, remove, set, unset must be set")
	}
	if valueSet(opts.Add) {
		if err := opts.Add.validate(); err != nil {
			return err
		}
	}
	if valueSet(opts.Remove) {
		if err := opts.Remove.validate(); err != nil {
			return err
		}
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

type ShareAdd struct {
	Accounts          []AccountIdentifier `ddl:"parameter" sql:"ACCOUNTS"`
	ShareRestrictions *bool               `ddl:"parameter" sql:"SHARE_RESTRICTIONS"`
}

func (v *ShareAdd) validate() error {
	if len(v.Accounts) == 0 {
		return fmt.Errorf("at least one account must be specified")
	}
	return nil
}

type ShareRemove struct {
	Accounts []AccountIdentifier `ddl:"parameter" sql:"ACCOUNTS"`
}

func (v *ShareRemove) validate() error {
	if len(v.Accounts) == 0 {
		return fmt.Errorf("at least one account must be specified")
	}
	return nil
}

type ShareSet struct {
	Accounts []AccountIdentifier `ddl:"parameter" sql:"ACCOUNTS"`
	Comment  *string             `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag      []TagAssociation    `ddl:"keyword" sql:"TAG"`
}

func (v *ShareSet) validate() error {
	if valueSet(v.Tag) && anyValueSet(v.Accounts, v.Comment) {
		return fmt.Errorf("accounts and comment cannot be set when tag is set")
	}
	return nil
}

type ShareUnset struct {
	Tag     []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
	Comment *bool              `ddl:"keyword" sql:"COMMENT"`
}

func (v *ShareUnset) validate() error {
	if ok := exactlyOneValueSet(v.Comment, v.Tag); !ok {
		return fmt.Errorf("exactly one of comment, tag must be set")
	}
	return nil
}

func (v *shares) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterShareOptions) error {
	if opts == nil {
		opts = &AlterShareOptions{}
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

type ShowShareOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`   //lint:ignore U1000 This is used in the ddl tag
	shares     bool       `ddl:"static" sql:"SHARES"` //lint:ignore U1000 This is used in the ddl tag
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	StartsWith *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

func (opts *ShowShareOptions) validate() error {
	return nil
}

func (s *shares) Show(ctx context.Context, opts *ShowShareOptions) ([]*Share, error) {
	if opts == nil {
		opts = &ShowShareOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []*shareRow
	err = s.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	shares := make([]*Share, 0, len(rows))
	for _, row := range rows {
		shares = append(shares, row.toShare())
	}
	return shares, nil
}

func (s *shares) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Share, error) {
	shares, err := s.Show(ctx, &ShowShareOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	for _, share := range shares {
		if share.Name.Name() == id.Name() {
			return share, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type ShareDetails struct {
	SharedObjects []ShareInfo
}

type ShareInfo struct {
	Kind     ObjectType
	Name     ObjectIdentifier
	SharedOn time.Time
}

type shareDetailsRow struct {
	Kind     string    `db:"kind"`
	Name     string    `db:"name"`
	SharedOn time.Time `db:"shared_on"`
}

func (row *shareDetailsRow) toShareInfo() *ShareInfo {
	objectType := ObjectType(row.Kind)
	trimmedS := strings.Trim(row.Name, "\"")
	id := objectType.GetObjectIdentifier(trimmedS)
	return &ShareInfo{
		Kind:     objectType,
		Name:     id,
		SharedOn: row.SharedOn,
	}
}

func shareDetailsFromRows(rows []shareDetailsRow) *ShareDetails {
	v := &ShareDetails{}
	for _, row := range rows {
		v.SharedObjects = append(v.SharedObjects, *row.toShareInfo())
	}
	return v
}

type shareDescribeOptions struct {
	describe bool             `ddl:"static" sql:"DESCRIBE"` //lint:ignore U1000 This is used in the ddl tag
	share    bool             `ddl:"static" sql:"SHARE"`    //lint:ignore U1000 This is used in the ddl tag
	name     ObjectIdentifier `ddl:"identifier"`
}

func (opts *shareDescribeOptions) validate() error {
	if ok := validObjectidentifier(opts.name); !ok {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (c *shares) DescribeProvider(ctx context.Context, id AccountObjectIdentifier) (*ShareDetails, error) {
	opts := &shareDescribeOptions{
		name: id,
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []shareDetailsRow
	err = c.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	return shareDetailsFromRows(rows), nil
}

func (c *shares) DescribeConsumer(ctx context.Context, id ExternalObjectIdentifier) (*ShareDetails, error) {
	opts := &shareDescribeOptions{
		name: id,
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []shareDetailsRow
	err = c.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	return shareDetailsFromRows(rows), nil
}
