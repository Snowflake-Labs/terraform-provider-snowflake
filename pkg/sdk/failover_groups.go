package sdk

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

// Compile-time proof of interface implementation.
var _ FailoverGroups = (*failoverGroups)(nil)

// FailoverGroups describes all the failover group related methods that the
// Snowflake API supports.
type FailoverGroups interface {
	// Create creates a new failover group.
	Create(ctx context.Context, id AccountObjectIdentifier, objectTypes []PluralObjectType, allowedAccounts []AccountIdentifier, opts *CreateFailoverGroupOptions) error
	// CreateSecondaryReplicationGroup creates a new secondary replication group.
	CreateSecondaryReplicationGroup(ctx context.Context, id AccountObjectIdentifier, primaryFailoverGroupID ExternalObjectIdentifier, opts *CreateSecondaryReplicationGroupOptions) error
	// Alter modifies an existing failover group in a source acount.
	AlterSource(ctx context.Context, id AccountObjectIdentifier, opts *AlterSourceFailoverGroupOptions) error
	// AlterTarget modifies an existing failover group in a target acount.
	AlterTarget(ctx context.Context, id AccountObjectIdentifier, opts *AlterTargetFailoverGroupOptions) error
	// Drop removes a failover group.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropFailoverGroupOptions) error
	// Show returns a list of failover groups.
	Show(ctx context.Context, opts *ShowFailoverGroupOptions) ([]*FailoverGroup, error)
	// ShowByID returns a failover group by ID
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*FailoverGroup, error)
	// ShowDatabases returns a list of databases in a failover group.
	ShowDatabases(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error)
	// ShowShares returns a list of shares in a failover group.
	ShowShares(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error)
}

// FailoverGroups implements FailoverGroups.
type failoverGroups struct {
	client *Client
}

// IntegrationType is the type of integration.
type IntegrationType string

const (
	IntegrationTypeSecurityIntegrations     IntegrationType = "SECURITY INTEGRATIONS"
	IntegrationTypeAPIIntegrations          IntegrationType = "API INTEGRATIONS"
	IntegrationTypeNotificationIntegrations IntegrationType = "NOTIFICATION INTEGRATIONS"
)

type CreateFailoverGroupOptions struct {
	create        bool                    `ddl:"static" sql:"CREATE"`         //lint:ignore U1000 This is used in the ddl tag
	failoverGroup bool                    `ddl:"static" sql:"FAILOVER GROUP"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists   *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`

	objectTypes             []PluralObjectType        `ddl:"parameter" sql:"OBJECT_TYPES"`
	AllowedDatabases        []AccountObjectIdentifier `ddl:"parameter" sql:"ALLOWED_DATABASES"`
	AllowedShares           []AccountObjectIdentifier `ddl:"parameter" sql:"ALLOWED_SHARES"`
	AllowedIntegrationTypes []IntegrationType         `ddl:"parameter" sql:"ALLOWED_INTEGRATION_TYPES"`
	allowedAccounts         []AccountIdentifier       `ddl:"parameter" sql:"ALLOWED_ACCOUNTS"`
	IgnoreEditionCheck      *bool                     `ddl:"keyword" sql:"IGNORE EDITION CHECK"`
	ReplicationSchedule     *string                   `ddl:"parameter,single_quotes" sql:"REPLICATION_SCHEDULE"`
}

func (opts *CreateFailoverGroupOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *failoverGroups) Create(ctx context.Context, id AccountObjectIdentifier, objectTypes []PluralObjectType, allowedAccounts []AccountIdentifier, opts *CreateFailoverGroupOptions) error {
	if opts == nil {
		opts = &CreateFailoverGroupOptions{}
	}
	opts.name = id
	opts.allowedAccounts = allowedAccounts
	opts.objectTypes = objectTypes
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

type CreateSecondaryReplicationGroupOptions struct {
	create               bool                     `ddl:"static" sql:"CREATE"`         //lint:ignore U1000 This is used in the ddl tag
	failoverGroup        bool                     `ddl:"static" sql:"FAILOVER GROUP"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists          *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                 AccountObjectIdentifier  `ddl:"identifier"`
	primaryFailoverGroup ExternalObjectIdentifier `ddl:"identifier" sql:"AS REPLICA OF"`
}

func (opts *CreateSecondaryReplicationGroupOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if !validObjectidentifier(opts.primaryFailoverGroup) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *failoverGroups) CreateSecondaryReplicationGroup(ctx context.Context, id AccountObjectIdentifier, primaryFailoverGroupID ExternalObjectIdentifier, opts *CreateSecondaryReplicationGroupOptions) error {
	if opts == nil {
		opts = &CreateSecondaryReplicationGroupOptions{}
	}
	opts.name = id
	opts.primaryFailoverGroup = primaryFailoverGroupID
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

type AlterSourceFailoverGroupOptions struct {
	alter         bool                    `ddl:"static" sql:"ALTER"`          //lint:ignore U1000 This is used in the ddl tag
	failoverGroup bool                    `ddl:"static" sql:"FAILOVER GROUP"` //lint:ignore U1000 This is used in the ddl tag
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`
	NewName       AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	Set           *FailoverGroupSet       `ddl:"keyword" sql:"SET"`
	Add           *FailoverGroupAdd       `ddl:"keyword" sql:"ADD"`
	Move          *FailoverGroupMove      `ddl:"keyword" sql:"MOVE"`
	Remove        *FailoverGroupRemove    `ddl:"keyword" sql:"REMOVE"`
}

func (opts *AlterSourceFailoverGroupOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if !exactlyOneValueSet(opts.Set, opts.Add, opts.Move, opts.Remove, opts.NewName) {
		return errors.New("exactly one of SET, ADD, MOVE, REMOVE, or NewName must be specified")
	}
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			return err
		}
	}
	if valueSet(opts.Add) {
		if err := opts.Add.validate(); err != nil {
			return err
		}
	}
	if valueSet(opts.Move) {
		if err := opts.Move.validate(); err != nil {
			return err
		}
	}
	if valueSet(opts.Remove) {
		if err := opts.Remove.validate(); err != nil {
			return err
		}
	}
	return nil
}

type FailoverGroupSet struct {
	ObjectTypes             []PluralObjectType `ddl:"parameter" sql:"OBJECT_TYPES"`
	ReplicationSchedule     *string            `ddl:"parameter,single_quotes" sql:"REPLICATION_SCHEDULE"`
	AllowedIntegrationTypes []IntegrationType  `ddl:"parameter" sql:"ALLOWED_INTEGRATION_TYPES"`
}

func (v *FailoverGroupSet) validate() error {
	if len(v.AllowedIntegrationTypes) > 0 {
		// INTEGRATIONS must be set in object types
		if !slices.Contains(v.ObjectTypes, PluralObjectTypeIntegrations) {
			return errors.New("INTEGRATIONS must be set in OBJECT_TYPES when setting allowed integration types")
		}
	}
	return nil
}

type FailoverGroupAdd struct {
	AllowedDatabases   []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"TO ALLOWED_DATABASES"`
	AllowedShares      []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"TO ALLOWED_SHARES"`
	AllowedAccounts    []AccountIdentifier       `ddl:"parameter,reverse" sql:"TO ALLOWED_ACCOUNTS"`
	IgnoreEditionCheck *bool                     `ddl:"keyword" sql:"IGNORE_EDITION_CHECK"`
}

func (v *FailoverGroupAdd) validate() error {
	return nil
}

type FailoverGroupMove struct {
	Databases []AccountObjectIdentifier `ddl:"parameter,no_equals" sql:"DATABASES"`
	Shares    []AccountObjectIdentifier `ddl:"parameter,no_equals" sql:"SHARES"`
	To        AccountObjectIdentifier   `ddl:"identifier" sql:"TO FAILOVER GROUP"`
}

func (v *FailoverGroupMove) validate() error {
	return nil
}

type FailoverGroupRemove struct {
	AllowedDatabases []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"FROM ALLOWED_DATABASES"`
	AllowedShares    []AccountObjectIdentifier `ddl:"parameter,reverse" sql:"FROM ALLOWED_SHARES"`
	AllowedAccounts  []AccountIdentifier       `ddl:"parameter,reverse" sql:"FROM ALLOWED_ACCOUNTS"`
}

func (v *FailoverGroupRemove) validate() error {
	return nil
}

func (v *failoverGroups) AlterSource(ctx context.Context, id AccountObjectIdentifier, opts *AlterSourceFailoverGroupOptions) error {
	if opts == nil {
		opts = &AlterSourceFailoverGroupOptions{}
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

type AlterTargetFailoverGroupOptions struct {
	alter         bool                    `ddl:"static" sql:"ALTER"`          //lint:ignore U1000 This is used in the ddl tag
	failoverGroup bool                    `ddl:"static" sql:"FAILOVER GROUP"` //lint:ignore U1000 This is used in the ddl tag
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`
	Refresh       *bool                   `ddl:"keyword" sql:"REFRESH"`
	Primary       *bool                   `ddl:"keyword" sql:"PRIMARY"`
	Suspend       *bool                   `ddl:"keyword" sql:"SUSPEND"`
	Resume        *bool                   `ddl:"keyword" sql:"RESUME"`
}

func (opts *AlterTargetFailoverGroupOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if !exactlyOneValueSet(opts.Refresh, opts.Primary, opts.Suspend, opts.Resume) {
		return errors.New("must set one of [Refresh, Primary, Suspend, Resume]")
	}
	return nil
}

func (v *failoverGroups) AlterTarget(ctx context.Context, id AccountObjectIdentifier, opts *AlterTargetFailoverGroupOptions) error {
	if opts == nil {
		opts = &AlterTargetFailoverGroupOptions{}
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

type DropFailoverGroupOptions struct {
	drop          bool                    `ddl:"static" sql:"DROP"`           //lint:ignore U1000 This is used in the ddl tag
	failoverGroup bool                    `ddl:"static" sql:"FAILOVER GROUP"` //lint:ignore U1000 This is used in the ddl tag
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *DropFailoverGroupOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *failoverGroups) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropFailoverGroupOptions) error {
	if opts == nil {
		opts = &DropFailoverGroupOptions{}
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
	if err != nil {
		return err
	}
	return err
}

// ShowFailoverGroupOptions represents the options for listing failover groups.
type ShowFailoverGroupOptions struct {
	show           bool              `ddl:"static" sql:"SHOW"`            //lint:ignore U1000 This is used in the ddl tag
	failoverGroups bool              `ddl:"static" sql:"FAILOVER GROUPS"` //lint:ignore U1000 This is used in the ddl tag
	InAccount      AccountIdentifier `ddl:"identifier" sql:"IN ACCOUNT"`
}

func (opts *ShowFailoverGroupOptions) validate() error {
	return nil
}

type FailoverGroupSecondaryState string

const (
	FailoverGroupSecondaryStateSuspended FailoverGroupSecondaryState = "SUSPENDED"
	FailoverGroupSecondaryStateStarted   FailoverGroupSecondaryState = "STARTED"
	FailoverGroupSecondaryStateNull      FailoverGroupSecondaryState = "NULL"
)

// FailoverGroups is a user friendly result for a CREATE FAILOVER GROUP query.
type FailoverGroup struct {
	RegionGroup             string
	SnowflakeRegion         string
	CreatedOn               time.Time
	AccountName             string
	Name                    string
	Type                    string
	Comment                 string
	IsPrimary               bool
	Primary                 ExternalObjectIdentifier
	ObjectTypes             []PluralObjectType
	AllowedIntegrationTypes []IntegrationType
	AllowedAccounts         []AccountIdentifier
	OrganizationName        string
	AccountLocator          string
	ReplicationSchedule     string
	SecondaryState          FailoverGroupSecondaryState
	NextScheduledRefresh    string
	Owner                   string
}

func (v *FailoverGroup) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *FailoverGroup) ExternalID() ExternalObjectIdentifier {
	return NewExternalObjectIdentifier(AccountIdentifier{
		organizationName: v.OrganizationName,
		accountName:      v.AccountName,
		accountLocator:   v.AccountLocator,
	}, v.ID())
}

func (v *FailoverGroup) ObjectType() ObjectType {
	return ObjectTypeFailoverGroup
}

// failoverGroupDBRow is used to decode the result of a CREATE FAILOVER GROUP query.
type failoverGroupDBRow struct {
	RegionGroup             string         `db:"region_group"`
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
	ReplicationSchedule     sql.NullString `db:"replication_schedule"`
	SecondaryState          sql.NullString `db:"secondary_state"`
	NextScheduledRefresh    sql.NullString `db:"next_scheduled_refresh"`
	Owner                   sql.NullString `db:"owner"`
}

func (row failoverGroupDBRow) toFailoverGroup() *FailoverGroup {
	ots := strings.Split(row.ObjectTypes, ",")
	pluralObjectTypes := make([]PluralObjectType, 0, len(ots))
	for _, ot := range ots {
		pluralObjectTypes = append(pluralObjectTypes, PluralObjectType(strings.TrimSpace(ot)))
	}
	its := strings.Split(row.AllowedIntegrationTypes, ",")
	allowedIntegrationTypes := make([]IntegrationType, 0, len(its))
	for _, it := range its {
		if it == "" {
			continue
		}
		allowedIntegrationTypes = append(allowedIntegrationTypes, IntegrationType(strings.TrimSpace(it)+" INTEGRATIONS"))
	}
	aas := strings.Split(row.AllowedAccounts, ",")
	allowedAccounts := make([]AccountIdentifier, 0, len(aas))
	for _, aa := range aas {
		s := strings.TrimSpace(aa)
		p := strings.Split(s, ".")
		if len(p) != 2 {
			continue
		}
		orgName := p[0]
		accountName := p[1]
		allowedAccounts = append(allowedAccounts, NewAccountIdentifier(orgName, accountName))
	}
	var comment string
	if row.Comment.Valid {
		comment = row.Comment.String
	}
	var replicationSchedule string
	if row.ReplicationSchedule.Valid {
		replicationSchedule = row.ReplicationSchedule.String
	}

	secondaryState := FailoverGroupSecondaryStateNull
	if row.SecondaryState.Valid {
		secondaryState = FailoverGroupSecondaryState(row.SecondaryState.String)
	}
	nextScheduledRefresh := ""
	if row.NextScheduledRefresh.Valid {
		nextScheduledRefresh = row.NextScheduledRefresh.String
	}
	return &FailoverGroup{
		RegionGroup:             row.RegionGroup,
		SnowflakeRegion:         row.SnowflakeRegion,
		CreatedOn:               row.CreatedOn,
		AccountName:             row.AccountName,
		OrganizationName:        row.OrganizationName,
		AccountLocator:          row.AccountLocator,
		Name:                    row.Name,
		Comment:                 comment,
		IsPrimary:               row.IsPrimary,
		Primary:                 NewExternalObjectIdentifierFromFullyQualifiedName(row.Primary),
		ObjectTypes:             pluralObjectTypes,
		AllowedIntegrationTypes: allowedIntegrationTypes,
		AllowedAccounts:         allowedAccounts,
		ReplicationSchedule:     replicationSchedule,
		SecondaryState:          secondaryState,
		NextScheduledRefresh:    nextScheduledRefresh,
		Owner:                   row.Owner.String,
		Type:                    row.Type,
	}
}

// List all the failover groups by pattern.
func (v *failoverGroups) Show(ctx context.Context, opts *ShowFailoverGroupOptions) ([]*FailoverGroup, error) {
	if opts == nil {
		opts = &ShowFailoverGroupOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []failoverGroupDBRow{}

	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]*FailoverGroup, len(dest))
	for i, row := range dest {
		resultList[i] = row.toFailoverGroup()
	}

	return resultList, nil
}

func (v *failoverGroups) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*FailoverGroup, error) {
	currentAccount, err := v.client.ContextFunctions.CurrentAccount(ctx)
	if err != nil {
		return nil, err
	}
	failoverGroups, err := v.Show(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, failoverGroup := range failoverGroups {
		if failoverGroup.ID() == id && failoverGroup.AccountLocator == currentAccount {
			return failoverGroup, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type showFailoverGroupDatabasesOptions struct {
	show      bool                    `ddl:"static" sql:"SHOW"`      //lint:ignore U1000 This is used in the ddl tag
	databases bool                    `ddl:"static" sql:"DATABASES"` //lint:ignore U1000 This is used in the ddl tag
	in        AccountObjectIdentifier `ddl:"identifier" sql:"IN FAILOVER GROUP"`
}

func (opts *showFailoverGroupDatabasesOptions) validate() error {
	if !validObjectidentifier(opts.in) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *failoverGroups) ShowDatabases(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error) {
	opts := &showFailoverGroupDatabasesOptions{
		in: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []struct {
		Name string `db:"name"`
	}{}
	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]AccountObjectIdentifier, len(dest))
	for i, row := range dest {
		resultList[i] = NewAccountObjectIdentifier(row.Name)
	}
	return resultList, nil
}

type showFailoverGroupSharesOptions struct {
	show      bool                    `ddl:"static" sql:"SHOW"`   //lint:ignore U1000 This is used in the ddl tag
	databases bool                    `ddl:"static" sql:"SHARES"` //lint:ignore U1000 This is used in the ddl tag
	in        AccountObjectIdentifier `ddl:"identifier" sql:"IN FAILOVER GROUP"`
}

func (opts *showFailoverGroupSharesOptions) validate() error {
	if !validObjectidentifier(opts.in) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *failoverGroups) ShowShares(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error) {
	opts := &showFailoverGroupSharesOptions{
		in: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []struct {
		Name string `db:"name"`
	}{}
	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]AccountObjectIdentifier, len(dest))
	for i, row := range dest {
		resultList[i] = NewExternalObjectIdentifierFromFullyQualifiedName(row.Name).objectIdentifier.(AccountObjectIdentifier)
	}
	return resultList, nil
}
