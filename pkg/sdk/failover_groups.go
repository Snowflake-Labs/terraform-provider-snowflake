package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// Compile-time proof of interface implementation.
var _ FailoverGroups = (*failoverGroups)(nil)

// FailoverGroups describes all the failover group related methods that the
// Snowflake API supports.
type FailoverGroups interface {
	// Create creates a new failover group.
	Create(ctx context.Context, id AccountObjectIdentifier, objectTypes []ObjectType, allowedAccounts []AccountIdentifier, opts *FailoverGroupCreateOptions) error
	// CreateSecondaryReplicationGroup creates a new secondary replication group.
	CreateSecondaryReplicationGroup(ctx context.Context, id AccountObjectIdentifier, primaryFailoverGroupID ExternalObjectIdentifier, opts *FailoverGroupCreateSecondaryReplicationGroupOptions) error
	// Alter modifies an existing failover group in a source acount.
	AlterSource(ctx context.Context, id AccountObjectIdentifier, opts *FailoverGroupAlterSourceOptions) error
	// AlterTarget modifies an existing failover group in a target acount.
	AlterTarget(ctx context.Context, id AccountObjectIdentifier, opts *FailoverGroupAlterTargetOptions) error
	// Drop removes a failover group.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *FailoverGroupDropOptions) error
	// Show returns a list offailover groups.
	Show(ctx context.Context, opts *FailoverGroupShowOptions) ([]*FailoverGroup, error)
	// ShowDatabases returns a list of databases in a failover group.
	ShowDatabases(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error)
	// ShowShares returns a list of shares in a failover group.
	ShowShares(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error)
}

// FailoverGroups implements FailoverGroups.
type failoverGroups struct {
	client  *Client
	builder *sqlBuilder
}

// IntegrationType is the type of integration.
type IntegrationType string

const (
	IntegrationTypeSecurityIntegrations     IntegrationType = "SECURITY INTEGRATIONS"
	IntegrationTypeAPIIntegrations          IntegrationType = "API INTEGRATIONS"
	IntegrationTypeNotificationIntegrations IntegrationType = "NOTIFICATION INTEGRATIONS"
)

type FailoverGroupCreateOptions struct {
	create        bool                    `ddl:"static" db:"CREATE"`         //lint:ignore U1000 This is used in the ddl tag
	failoverGroup bool                    `ddl:"static" db:"FAILOVER GROUP"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists   *bool                   `ddl:"keyword" db:"IF NOT EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`

	objectTypes             []string                  `ddl:"list,no_parentheses" db:"OBJECT_TYPES ="`
	AllowedDatabases        []AccountObjectIdentifier `ddl:"list,no_parentheses" db:"ALLOWED_DATABASES ="`
	AllowedShares           []AccountObjectIdentifier `ddl:"list,no_parentheses" db:"ALLOWED_SHARES ="`
	AllowedIntegrationTypes []IntegrationType         `ddl:"list,no_parentheses" db:"ALLOWED_INTEGRATION_TYPES ="`
	allowedAccounts         []AccountIdentifier       `ddl:"list,no_parentheses" db:"ALLOWED_ACCOUNTS ="`
	IgnoreEditionCheck      *bool                     `ddl:"keyword" db:"IGNORE EDITION CHECK"`
	ReplicationSchedule     *string                   `ddl:"parameter,single_quotes" db:"REPLICATION_SCHEDULE"`
}

func (opts *FailoverGroupCreateOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (v *failoverGroups) Create(ctx context.Context, id AccountObjectIdentifier, objectTypes []ObjectType, allowedAccounts []AccountIdentifier, opts *FailoverGroupCreateOptions) error {
	if opts == nil {
		opts = &FailoverGroupCreateOptions{}
	}
	opts.name = id
	opts.allowedAccounts = allowedAccounts

	// convert objectTypes to plural.
	var objectTypesStrList []string
	for _, objectType := range objectTypes {
		objectTypesStrList = append(objectTypesStrList, objectType.Plural())
	}
	opts.objectTypes = objectTypesStrList
	if err := opts.validate(); err != nil {
		return err
	}
	log.Printf("[DEBUG] creating failover group: %v", opts.name)
	clauses, err := v.builder.parseStruct(opts)
	if err != nil {
		return err
	}
	stmt := v.builder.sql(clauses...)
	_, err = v.client.exec(ctx, stmt)
	return err
}

type FailoverGroupCreateSecondaryReplicationGroupOptions struct {
	create               bool                     `ddl:"static" db:"CREATE"`         //lint:ignore U1000 This is used in the ddl tag
	failoverGroup        bool                     `ddl:"static" db:"FAILOVER GROUP"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists          *bool                    `ddl:"keyword" db:"IF NOT EXISTS"`
	name                 AccountObjectIdentifier  `ddl:"identifier"`
	primaryFailoverGroup ExternalObjectIdentifier `ddl:"identifier" db:"AS REPLICA OF"`
}

func (opts *FailoverGroupCreateSecondaryReplicationGroupOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
	}
	if opts.primaryFailoverGroup.FullyQualifiedName() == "" {
		return fmt.Errorf("primaryFailoverGroup is required")
	}
	return nil
}

func (v *failoverGroups) CreateSecondaryReplicationGroup(ctx context.Context, id AccountObjectIdentifier, primaryFailoverGroupID ExternalObjectIdentifier, opts *FailoverGroupCreateSecondaryReplicationGroupOptions) error {
	if opts == nil {
		opts = &FailoverGroupCreateSecondaryReplicationGroupOptions{}
	}
	opts.name = id
	opts.primaryFailoverGroup = primaryFailoverGroupID
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

type FailoverGroupAlterSourceOptions struct {
	alter         bool                    `ddl:"static" db:"ALTER"`          //lint:ignore U1000 This is used in the ddl tag
	failoverGroup bool                    `ddl:"static" db:"FAILOVER GROUP"` //lint:ignore U1000 This is used in the ddl tag
	IfExists      *bool                   `ddl:"keyword" db:"IF EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`
	NewName       AccountObjectIdentifier `ddl:"identifier" db:"RENAME TO"`
	Set           *FailoverGroupSet       `ddl:"keyword" db:"SET"`
	Add           *FailoverGroupAdd       `ddl:"keyword" db:"ADD"`
	Move          *FailoverGroupMove      `ddl:"keyword" db:"MOVE"`
	Remove        *FailoverGroupRemove    `ddl:"keyword" db:"REMOVE"`
}

func (opts *FailoverGroupAlterSourceOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

type FailoverGroupSet struct {
	ObjectTypes             []ObjectType      `ddl:"list" db:"OBJECT_TYPES"`
	ReplicationSchedule     *string           `ddl:"parameter,single_quotes" db:"REPLICATION_SCHEDULE"`
	AllowedIntegrationTypes []IntegrationType `ddl:"list" db:"ALLOWED_INTEGRATION_TYPES"`
}

type FailoverGroupAdd struct {
	AllowedDatabases   []AccountObjectIdentifier `ddl:"command,reverse" db:"TO ALLOWED DATABASES"`
	AllowedShares      []AccountObjectIdentifier `ddl:"command,reverse" db:"TO ALLOWED SHARES"`
	AllowedAccounts    []AccountIdentifier       `ddl:"command,reverse" db:"TO ALLOWED ACCOUNTS"`
	IgnoreEditionCheck *bool                     `ddl:"keyword" db:"IGNORE_EDITION_CHECK"`
}

type FailoverGroupMove struct {
	Databases []AccountObjectIdentifier `ddl:"command" db:"DATABASES"`
	Shares    []AccountObjectIdentifier `ddl:"command" db:"SHARES"`
	To        ExternalObjectIdentifier  `ddl:"identifier" db:"TO"`
}

type FailoverGroupRemove struct {
	AllowedDatabases []AccountObjectIdentifier `ddl:"command,reverse" db:"FROM ALLOWED DATABASES"`
	AllowedShares    []AccountObjectIdentifier `ddl:"command,reverse" db:"FROM ALLOWED SHARES"`
	AllowedAccounts  []AccountIdentifier       `ddl:"command,reverse" db:"FROM ALLOWED ACCOUNTS"`
}

func (v *failoverGroups) AlterSource(ctx context.Context, id AccountObjectIdentifier, opts *FailoverGroupAlterSourceOptions) error {
	if opts == nil {
		opts = &FailoverGroupAlterSourceOptions{}
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

type FailoverGroupAlterTargetOptions struct {
	alter         bool                    `ddl:"static" db:"ALTER"`          //lint:ignore U1000 This is used in the ddl tag
	failoverGroup bool                    `ddl:"static" db:"FAILOVER GROUP"` //lint:ignore U1000 This is used in the ddl tag
	IfExists      *bool                   `ddl:"keyword" db:"IF EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`
	Refresh       *bool                   `ddl:"keyword" db:"REFRESH"`
	Primary       *bool                   `ddl:"keyword" db:"PRIMARY"`
	Suspend       *bool                   `ddl:"keyword" db:"SUSPEND"`
	Resume        *bool                   `ddl:"keyword" db:"RESUME"`
}

func (opts *FailoverGroupAlterTargetOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return errors.New("name must not be empty")
	}
	// can only choose one of [Refresh, Primary, Suspend, Resume]
	setKeys := []bool{*opts.Refresh, *opts.Primary, *opts.Suspend, *opts.Resume}
	count := 0
	for _, v := range setKeys {
		if v {
			count++
		}
	}
	if count > 1 {
		return errors.New("can only choose one of [Refresh, Primary, Suspend, Resume]")
	}

	return nil
}

func (v *failoverGroups) AlterTarget(ctx context.Context, id AccountObjectIdentifier, opts *FailoverGroupAlterTargetOptions) error {
	if opts == nil {
		opts = &FailoverGroupAlterTargetOptions{}
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

type FailoverGroupDropOptions struct {
	drop          bool                    `ddl:"static" db:"DROP"`           //lint:ignore U1000 This is used in the ddl tag
	failoverGroup bool                    `ddl:"static" db:"FAILOVER GROUP"` //lint:ignore U1000 This is used in the ddl tag
	IfExists      *bool                   `ddl:"keyword" db:"IF EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *FailoverGroupDropOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func (v *failoverGroups) Drop(ctx context.Context, id AccountObjectIdentifier, opts *FailoverGroupDropOptions) error {
	if opts == nil {
		opts = &FailoverGroupDropOptions{}
	}
	opts.name = id
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

// FailoverGroupShowOptions represents the options for listing failover groups.
type FailoverGroupShowOptions struct {
	show           bool              `ddl:"static" db:"SHOW"`            //lint:ignore U1000 This is used in the ddl tag
	failoverGroups bool              `ddl:"static" db:"FAILOVER GROUPS"` //lint:ignore U1000 This is used in the ddl tag
	InAcccount     AccountIdentifier `ddl:"identifier" db:"IN ACCOUNT"`
}

func (opts *FailoverGroupShowOptions) validate() error {
	return nil
}

// FailoverGroups is a user friendly result for a CREATE FAILOVER GROUP query.
type FailoverGroup struct {
	SnowflakeRegion         string `db:"snowflake_region"`
	CreatedOn               time.Time
	AccountName             string
	Name                    string
	Comment                 string
	IsPrimary               bool
	Primary                 ExternalObjectIdentifier
	ObjectTypes             []ObjectType
	AllowedIntegrationTypes []IntegrationType
	AllowedAccounts         []AccountIdentifier
	OrganizationName        string
	ReplicationSchedule     string `db:"replication_schedule"`
}

func (v *FailoverGroup) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

// failoverGroupDBRow is used to decode the result of a CREATE FAILOVER GROUP query.
type failoverGroupDBRow struct {
	SnowflakeRegion         string         `db:"snowflake_region"`
	CreatedOn               time.Time      `db:"created_on"`
	AccountName             string         `db:"account_name"`
	Name                    string         `db:"name"`
	Type                    string         `db:"type"`
	Comment                 sql.NullString `db:"comment"`
	IsPrimary               bool           `db:"is_primary"`
	Primary                 string         `db:"primary"`
	ObjectTypes             string         `db:"object_types"`
	AllowedIntegrationTypes string         `db:"allowed_integration_types"`
	AllowedAccounts         string         `db:"allowed_accounts"`
	OrganizationName        string         `db:"organization_name"`
	AccountLocator          string         `db:"account_locator"`
	ReplicationSchedule     string         `db:"replication_schedule"`
	SecondaryState          sql.NullString `db:"secondary_state"`
	NextScheduledRefresh    sql.NullString `db:"next_scheduled_refresh"`
	Owner                   string         `db:"owner"`
}

func (row failoverGroupDBRow) toFailoverGroup() *FailoverGroup {
	var objectTypes []ObjectType
	for _, v := range strings.Split(row.ObjectTypes, ",") {
		objectTypes = append(objectTypes, ObjectTypeFromPluralString(strings.TrimSpace(v)))
	}
	var allowedIntegrationTypes []IntegrationType
	for _, v := range strings.Split(row.AllowedIntegrationTypes, ",") {
		if v == "" {
			continue
		}
		allowedIntegrationTypes = append(allowedIntegrationTypes, IntegrationType(strings.TrimSpace(v)))
	}
	log.Printf("allowedIntegrationTypes: %+v", allowedIntegrationTypes)
	var allowedAccounts []AccountIdentifier
	for _, v := range strings.Split(row.AllowedAccounts, ",") {
		s := strings.TrimSpace(v)
		p := strings.Split(s, ".")
		orgName := p[0]
		accountName := p[1]
		allowedAccounts = append(allowedAccounts, NewAccountIdentifier(orgName, accountName))
	}
	var comment string
	if row.Comment.Valid {
		comment = row.Comment.String
	}
	return &FailoverGroup{
		SnowflakeRegion:         row.SnowflakeRegion,
		CreatedOn:               row.CreatedOn,
		AccountName:             row.AccountName,
		Name:                    row.Name,
		Comment:                 comment,
		IsPrimary:               row.IsPrimary,
		Primary:                 NewExternalObjectIdentifierFromFullyQualifiedName(row.Primary),
		ObjectTypes:             objectTypes,
		AllowedIntegrationTypes: allowedIntegrationTypes,
		AllowedAccounts:         allowedAccounts,
		OrganizationName:        row.OrganizationName,
		ReplicationSchedule:     row.ReplicationSchedule,
	}
}

// List all the failover groups by pattern.
func (v *failoverGroups) Show(ctx context.Context, opts *FailoverGroupShowOptions) ([]*FailoverGroup, error) {
	if opts == nil {
		opts = &FailoverGroupShowOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	clauses, err := v.builder.parseStruct(opts)
	if err != nil {
		return nil, err
	}
	stmt := v.builder.sql(clauses...)
	dest := []failoverGroupDBRow{}

	err = v.client.query(ctx, &dest, stmt)
	if err != nil {
		return nil, decodeDriverError(err)
	}
	resultList := make([]*FailoverGroup, len(dest))
	for i, row := range dest {
		resultList[i] = row.toFailoverGroup()
	}

	return resultList, nil
}

func (v *failoverGroups) ShowDatabases(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error) {
	stmt := fmt.Sprintf("SHOW DATABASES IN FAILOVER GROUP %s", id.FullyQualifiedName())
	dest := []struct {
		Name string `db:"name"`
	}{}
	err := v.client.query(ctx, &dest, stmt)
	if err != nil {
		return nil, decodeDriverError(err)
	}
	resultList := make([]AccountObjectIdentifier, len(dest))
	for i, row := range dest {
		resultList[i] = NewAccountObjectIdentifier(row.Name)
	}
	return resultList, nil
}

func (v *failoverGroups) ShowShares(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error) {
	stmt := fmt.Sprintf("SHOW SHARES IN FAILOVER GROUP %s", id.FullyQualifiedName())
	dest := []struct {
		Name string `db:"name"`
	}{}
	err := v.client.query(ctx, &dest, stmt)
	if err != nil {
		return nil, decodeDriverError(err)
	}
	resultList := make([]AccountObjectIdentifier, len(dest))
	for i, row := range dest {
		resultList[i] = NewExternalObjectIdentifierFromFullyQualifiedName(row.Name).objectIdentifier.(AccountObjectIdentifier)
	}
	return resultList, nil
}
