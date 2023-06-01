package sdk

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
)

type Databases interface {
	// Create creates a database.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateDatabaseOptions) error
	// CreateShared creates a database from a shared database.
	CreateShared(ctx context.Context, id AccountObjectIdentifier, shareID ExternalObjectIdentifier, opts *CreateSharedDatabaseOptions) error
	// CreateSecondary creates a secondary database.
	CreateSecondary(ctx context.Context, id AccountObjectIdentifier, primaryID ExternalObjectIdentifier, opts *CreateSecondaryDatabaseOptions) error
	// Alter modifies an existing database
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseOptions) error
	// AlterReplication modifies an existing database replica
	AlterReplication(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseReplicationOptions) error
	// AlterFailover modifies an existing database failover group
	AlterFailover(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseFailoverOptions) error
	// Drop removes a database.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropDatabaseOptions) error
	// Undrop restores the most recent version of a dropped database
	Undrop(ctx context.Context, id AccountObjectIdentifier) error
	// Show returns a list of databases.
	Show(ctx context.Context, opts *ShowDatabasesOptions) ([]*Database, error)
	// ShowByID returns a database by ID
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Database, error)
	// Describe returns the details of a database.
	Describe(ctx context.Context, id AccountObjectIdentifier) (*DatabaseDetails, error)
	// Use sets the active database for the current session.
	Use(ctx context.Context, id AccountObjectIdentifier) error
}

var _ Databases = (*databases)(nil)

type databases struct {
	client *Client
}

type Database struct {
	CreatedOn     time.Time
	Name          string
	IsDefault     bool
	IsCurrent     bool
	Origin        string
	Owner         string
	Comment       string
	Options       string
	RetentionTime int
	ResourceGroup string
	DroppedOn     time.Time
	Transient     bool
}

func (v *Database) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *Database) ObjectType() ObjectType {
	return ObjectTypeDatabase
}

type databaseRow struct {
	CreatedOn     time.Time      `db:"created_on"`
	Name          string         `db:"name"`
	IsDefault     sql.NullString `db:"is_default"`
	IsCurrent     sql.NullString `db:"is_current"`
	Origin        sql.NullString `db:"origin"`
	Owner         sql.NullString `db:"owner"`
	Comment       sql.NullString `db:"comment"`
	Options       sql.NullString `db:"options"`
	RetentionTime sql.NullString `db:"retention_time"`
	ResourceGroup sql.NullString `db:"resource_group"`
	DroppedOn     sql.NullTime   `db:"dropped_on"`
}

func (row *databaseRow) toDatabase() *Database {
	database := Database{
		CreatedOn: row.CreatedOn,
		Name:      row.Name,
	}
	if row.IsDefault.Valid {
		database.IsDefault = row.IsDefault.String == "Y"
	}
	if row.IsCurrent.Valid {
		database.IsCurrent = row.IsCurrent.String == "Y"
	}
	if row.Origin.Valid {
		database.Origin = row.Origin.String
	}
	if row.Owner.Valid {
		database.Owner = row.Owner.String
	}
	if row.Comment.Valid {
		database.Comment = row.Comment.String
	}
	if row.Options.Valid {
		database.Options = row.Options.String
	}
	if row.RetentionTime.Valid {
		retentionTimeInt, err := strconv.Atoi(row.RetentionTime.String)
		if err != nil {
			database.RetentionTime = 0
		}
		database.RetentionTime = retentionTimeInt
	}
	if row.ResourceGroup.Valid {
		database.ResourceGroup = row.ResourceGroup.String
	}
	if row.DroppedOn.Valid {
		database.DroppedOn = row.DroppedOn.Time
	}
	if row.Options.Valid {
		parts := strings.Split(row.Options.String, ", ")
		for _, part := range parts {
			if part == "TRANSIENT" {
				database.Transient = true
			}
		}
	}
	return &database
}

type CreateDatabaseOptions struct {
	create                     bool                    `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace                  *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	Transient                  *bool                   `ddl:"keyword" sql:"TRANSIENT"`
	database                   bool                    `ddl:"static" sql:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	name                       AccountObjectIdentifier `ddl:"identifier"`            //lint:ignore U1000 This is used in the ddl tag
	IfNotExists                *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	Clone                      *Clone                  `ddl:"-"`
	DataRetentionTimeInDays    *int                    `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int                    `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	Comment                    *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag                        []TagAssociation        `ddl:"keyword,parentheses" sql:"TAG"`
}

func (opts *CreateDatabaseOptions) validate() error {
	if valueSet(opts.Clone) {
		if err := opts.Clone.validate(); err != nil {
			return err
		}
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		return errors.New("IF NOT EXISTS and OR REPLACE are incompatible.")
	}
	return nil
}

func (v *databases) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateDatabaseOptions) error {
	if opts == nil {
		opts = &CreateDatabaseOptions{}
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

type CreateSharedDatabaseOptions struct {
	create    bool                     `ddl:"static" sql:"CREATE"`   //lint:ignore U1000 This is used in the ddl tag
	database  bool                     `ddl:"static" sql:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	name      AccountObjectIdentifier  `ddl:"identifier"`            //lint:ignore U1000 This is used in the ddl tag
	fromShare ExternalObjectIdentifier `ddl:"identifier" sql:"FROM SHARE"`
	Comment   *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreateSharedDatabaseOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if !validObjectidentifier(opts.fromShare) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *databases) CreateShared(ctx context.Context, id AccountObjectIdentifier, shareID ExternalObjectIdentifier, opts *CreateSharedDatabaseOptions) error {
	if opts == nil {
		opts = &CreateSharedDatabaseOptions{}
	}

	opts.name = id
	opts.fromShare = shareID

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

type CreateSecondaryDatabaseOptions struct {
	create                  bool                     `ddl:"static" sql:"CREATE"`   //lint:ignore U1000 This is used in the ddl tag
	database                bool                     `ddl:"static" sql:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	name                    AccountObjectIdentifier  `ddl:"identifier"`            //lint:ignore U1000 This is used in the ddl tag
	primaryDatabase         ExternalObjectIdentifier `ddl:"identifier" sql:"AS REPLICA OF"`
	DataRetentionTimeInDays *int                     `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
}

func (opts *CreateSecondaryDatabaseOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if !validObjectidentifier(opts.primaryDatabase) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *databases) CreateSecondary(ctx context.Context, id AccountObjectIdentifier, primaryID ExternalObjectIdentifier, opts *CreateSecondaryDatabaseOptions) error {
	if opts == nil {
		opts = &CreateSecondaryDatabaseOptions{}
	}
	opts.name = id
	opts.primaryDatabase = primaryID
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

type AlterDatabaseOptions struct {
	alter    bool                    `ddl:"static" sql:"ALTER"`    //lint:ignore U1000 This is used in the ddl tag
	database bool                    `ddl:"static" sql:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`
	NewName  AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SwapWith AccountObjectIdentifier `ddl:"identifier" sql:"SWAP WITH"`
	Set      *DatabaseSet            `ddl:"list,no_parentheses" sql:"SET"`
	Unset    *DatabaseUnset          `ddl:"list,no_parentheses" sql:"UNSET"`
}

func (opts *AlterDatabaseOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if validObjectidentifier(opts.NewName) && anyValueSet(opts.Set, opts.Unset, opts.SwapWith) {
		return errors.New("RENAME TO cannot be set with other options")
	}

	if validObjectidentifier(opts.SwapWith) && anyValueSet(opts.Set, opts.Unset, opts.NewName) {
		return errors.New("SWAP WITH cannot be set with other options")
	}

	if valueSet(opts.Set) && valueSet(opts.Unset) {
		return errors.New("only one of SET or UNSET can be set")
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

type DatabaseSet struct {
	DataRetentionTimeInDays    *int    `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int    `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDDLCollation        *string `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (v *DatabaseSet) validate() error {
	return nil
}

type DatabaseUnset struct {
	DataRetentionTimeInDays    *bool              `ddl:"keyword" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *bool              `ddl:"keyword" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDDLCollation        *bool              `ddl:"keyword" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *bool              `ddl:"keyword" sql:"COMMENT"`
	Tag                        []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
}

func (v *DatabaseUnset) validate() error {
	if valueSet(v.Tag) {
		if anyValueSet(v.DataRetentionTimeInDays, v.MaxDataExtensionTimeInDays, v.DefaultDDLCollation, v.Comment) {
			return errors.New("Tag cannot be set with other options")
		}
	}
	return nil
}

func (v *databases) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseOptions) error {
	if opts == nil {
		opts = &AlterDatabaseOptions{}
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

type AlterDatabaseReplicationOptions struct {
	alter              bool                    `ddl:"static" sql:"ALTER"`    //lint:ignore U1000 This is used in the ddl tag
	database           bool                    `ddl:"static" sql:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	name               AccountObjectIdentifier `ddl:"identifier"`
	EnableReplication  *EnableReplication      `ddl:"keyword" sql:"ENABLE REPLICATION"`
	DisableReplication *DisableReplication     `ddl:"keyword" sql:"DISABLE REPLICATION"`
	Refresh            *bool                   `ddl:"keyword" sql:"REFRESH"`
}

func (opts *AlterDatabaseReplicationOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if everyValueNil(opts.EnableReplication, opts.DisableReplication, opts.Refresh) {
		return errors.New("one of ENABLE REPLICATION, DISABLE REPLICATION or REFRESH must be set")
	}
	if anyValueSet(opts.EnableReplication, opts.DisableReplication) {
		return errors.New("only one of ENABLE REPLICATION or DISABLE REPLICATION can be set")
	}
	if valueSet(opts.EnableReplication) {
		if err := opts.EnableReplication.validate(); err != nil {
			return err
		}
	}
	if valueSet(opts.DisableReplication) {
		if err := opts.DisableReplication.validate(); err != nil {
			return err
		}
	}
	return nil
}

type EnableReplication struct {
	ToAccounts         []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
	IgnoreEditionCheck *bool               `ddl:"keyword" sql:"IGNORE EDITION CHECK"`
}

func (v *EnableReplication) validate() error {
	return nil
}

type DisableReplication struct {
	ToAccounts []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
}

func (v *DisableReplication) validate() error {
	return nil
}

func (v *databases) AlterReplication(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseReplicationOptions) error {
	if opts == nil {
		opts = &AlterDatabaseReplicationOptions{}
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

type AlterDatabaseFailoverOptions struct {
	alter           bool                    `ddl:"static" sql:"ALTER"`    //lint:ignore U1000 This is used in the ddl tag
	database        bool                    `ddl:"static" sql:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	name            AccountObjectIdentifier `ddl:"identifier"`
	EnableFailover  *EnableFailover         `ddl:"keyword" sql:"ENABLE FAILOVER"`
	DisableFailover *DisableFailover        `ddl:"keyword" sql:"DISABLE FAILOVER"`
	Primary         *bool                   `ddl:"keyword" sql:"PRIMARY"`
}

func (opts *AlterDatabaseFailoverOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if everyValueNil(opts.EnableFailover, opts.DisableFailover, opts.Primary) {
		return errors.New("one of ENABLE FAILOVER, DISABLE FAILOVER or PRIMARY must be set")
	}
	if !exactlyOneValueSet(opts.EnableFailover, opts.DisableFailover, opts.Primary) {
		return errors.New("only one of ENABLE FAILOVER, DISABLE FAILOVER or PRIMARY can be set")
	}
	if valueSet(opts.EnableFailover) {
		if err := opts.EnableFailover.validate(); err != nil {
			return err
		}
	}
	if valueSet(opts.DisableFailover) {
		if err := opts.DisableFailover.validate(); err != nil {
			return err
		}
	}
	return nil
}

type EnableFailover struct {
	ToAccounts []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
}

func (v *EnableFailover) validate() error {
	return nil
}

type DisableFailover struct {
	ToAccounts []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
}

func (v *DisableFailover) validate() error {
	return nil
}

func (v *databases) AlterFailover(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseFailoverOptions) error {
	if opts == nil {
		opts = &AlterDatabaseFailoverOptions{}
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

type DropDatabaseOptions struct {
	drop     bool                    `ddl:"static" sql:"DROP"`     //lint:ignore U1000 This is used in the ddl tag
	database bool                    `ddl:"static" sql:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"` //lint:ignore U1000 This is used in the ddl tag
}

func (opts *DropDatabaseOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *databases) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropDatabaseOptions) error {
	if opts == nil {
		opts = &DropDatabaseOptions{}
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

type undropDatabaseOptions struct {
	undrop   bool                    `ddl:"static" sql:"UNDROP"`   //lint:ignore U1000 This is used in the ddl tag
	database bool                    `ddl:"static" sql:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	name     AccountObjectIdentifier `ddl:"identifier"`            //lint:ignore U1000 This is used in the ddl tag
}

func (opts *undropDatabaseOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *databases) Undrop(ctx context.Context, id AccountObjectIdentifier) error {
	opts := &undropDatabaseOptions{
		name: id,
	}
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

type ShowDatabasesOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"` //lint:ignore U1000 This is used in the ddl tag
	Terse      *bool      `ddl:"keyword" sql:"TERSE"`
	databases  bool       `ddl:"static" sql:"DATABASES"` //lint:ignore U1000 This is used in the ddl tag
	History    *bool      `ddl:"keyword" sql:"HISTORY"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	StartsWith *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	LimitFrom  *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

func (opts *ShowDatabasesOptions) validate() error {
	return nil
}

func (v *databases) Show(ctx context.Context, opts *ShowDatabasesOptions) ([]*Database, error) {
	if opts == nil {
		opts = &ShowDatabasesOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []databaseRow
	err = v.client.query(ctx, &rows, sql)
	databases := make([]*Database, len(rows))
	for i, row := range rows {
		databases[i] = row.toDatabase()
	}
	return databases, err
}

func (v *databases) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Database, error) {
	databases, err := v.client.Databases.Show(ctx, &ShowDatabasesOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	for _, database := range databases {
		if database.ID() == id {
			return database, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type DatabaseDetails struct {
	Rows []DatabaseDetailsRow
}

type DatabaseDetailsRow struct {
	CreatedOn time.Time
	Name      string
	Kind      string
}

type describeDatabaseOptions struct {
	describe bool                    `ddl:"static" sql:"DESCRIBE"` //lint:ignore U1000 This is used in the ddl tag
	database bool                    `ddl:"static" sql:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	name     AccountObjectIdentifier `ddl:"identifier"`            //lint:ignore U1000 This is used in the ddl tag
}

func (opts *describeDatabaseOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *databases) Describe(ctx context.Context, id AccountObjectIdentifier) (*DatabaseDetails, error) {
	opts := &describeDatabaseOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []DatabaseDetailsRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	details := DatabaseDetails{
		Rows: rows,
	}
	return &details, err
}

func (v *databases) Use(ctx context.Context, id AccountObjectIdentifier) error {
	// proxy to sessions
	return v.client.Sessions.UseDatabase(ctx, id)
}
