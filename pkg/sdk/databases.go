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
var _ Databases = (*databases)(nil)

// MaskingPolicies describes all the database related methods that the
// Snowflake API supports.
type Databases interface {
	// Create creates a new standard database.
	Create(ctx context.Context, id AccountLevelIdentifier, opts *DatabaseCreateOptions) error
	// CreateShared creates a new shared database.
	CreateShared(ctx context.Context, id AccountLevelIdentifier, shareID InboundShareIdentifier) error
	// CreateSecondary creates a new secondary database.
	CreateSecondary(ctx context.Context, id AccountLevelIdentifier, primaryID ExternalObjectIdentifier, opts *DatabaseCreateSecondaryOptions) error
	// Alter modifies an existing database.
	Alter(ctx context.Context, id AccountLevelIdentifier, opts *DatabaseAlterOptions) error
	// AlterReplication changes the replication settings of a database.
	AlterReplication(ctx context.Context, id AccountLevelIdentifier, opts *DatabaseAlterReplicationOptions) error
	// AlterFailover changes the failover settings of a database.
	AlterFailover(ctx context.Context, id AccountLevelIdentifier, opts *DatabaseAlterFailoverOptions) error
	// Describe returns the details of a database.
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*DatabaseDetails, error)
	// Drop removes a database.
	Drop(ctx context.Context, id SchemaObjectIdentifier) error
	// Undrop restores a dropped database.
	Undrop(ctx context.Context, id SchemaObjectIdentifier) error
	// Use sets the current database context.
	Use(ctx context.Context, id SchemaObjectIdentifier) error
	// Show returns a list of databases.
	Show(ctx context.Context, opts *DatabaseShowOptions) ([]*Database, error)
}

// databases implements MaskingPolicies.
type databases struct {
	client  *Client
	builder *sqlBuilder
}

type DatabaseCreateOptions struct {
	create      bool                   `ddl:"static" db:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace   *bool                  `ddl:"keyword" db:"OR REPLACE"`
	transient   *bool                  `ddl:"keyword" db:"TRANSIENT"`
	database    bool                   `ddl:"static" db:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists *bool                  `ddl:"keyword" db:"IF NOT EXISTS"`
	name        AccountLevelIdentifier `ddl:"identifier"`

	// optional
	Clone                      *Clone           `ddl:"keyword" db:"CLONE"`
	DataRetentionTimeInDays    *int             `ddl:"parameter" db:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int             `ddl:"parameter" db:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDDLCollation        *string          `ddl:"parameter,single_quotes" db:"DEFAULT_DDL_COLLATION"`
	Tag                        []TagAssociation `ddl:"keyword" db:"TAG"`
	Comment                    *string          `ddl:"parameter,single_quotes" db:"COMMENT"`
}

func (opts *DatabaseCreateOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (v *databases) Create(ctx context.Context, id AccountLevelIdentifier, opts *DatabaseCreateOptions) error {
	if opts == nil {
		opts = &DatabaseCreateOptions{}
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

type databaseCreateSharedOptions struct {
	create   bool                   `ddl:"static" db:"CREATE"`   //lint:ignore U1000 This is used in the ddl tag
	database bool                   `ddl:"static" db:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	name     AccountLevelIdentifier `ddl:"identifier"`
	share    InboundShareIdentifier `ddl:"identifier" db:"FROM SHARE"`
}

func (opts *databaseCreateSharedOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (v *databases) CreateShared(ctx context.Context, id AccountLevelIdentifier, shareID InboundShareIdentifier) error {
	opts := &databaseCreateSharedOptions{
		name:  id,
		share: shareID,
	}
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

type DatabaseCreateSecondaryOptions struct {
	create                  bool                     `ddl:"static" db:"CREATE"`   //lint:ignore U1000 This is used in the ddl tag
	database                bool                     `ddl:"static" db:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	name                    AccountLevelIdentifier   `ddl:"identifier"`
	replicationIdentifier   ExternalObjectIdentifier `ddl:"identifier" db:"AS REPLICA OF"`
	DataRetentionTimeInDays *int                     `ddl:"parameter" db:"DATA_RETENTION_TIME_IN_DAYS"`
}

func (opts *DatabaseCreateSecondaryOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (v *databases) CreateSecondary(ctx context.Context, id AccountLevelIdentifier, primaryID ExternalObjectIdentifier, opts *DatabaseCreateSecondaryOptions) error {
	if opts == nil {
		opts = &DatabaseCreateSecondaryOptions{}
	}
	opts.name = id
	opts.replicationIdentifier = primaryID
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

type DatabaseAlterOptions struct {
	alter    bool                   `ddl:"static" db:"ALTER"`    //lint:ignore U1000 This is used in the ddl tag
	database bool                   `ddl:"static" db:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                  `ddl:"keyword" db:"IF EXISTS"`
	name     AccountLevelIdentifier `ddl:"identifier"`
	NewName  AccountLevelIdentifier `ddl:"identifier" db:"RENAME TO"`
	Swap     AccountLevelIdentifier `ddl:"identifier" db:"SWAP WITH"`
	Set      *DatabaseSet           `ddl:"keyword" db:"SET"`
	Unset    *DatabaseUnset         `ddl:"keyword" db:"UNSET"`
}

func (opts *DatabaseAlterOptions) validate() error {
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

	if opts.Swap.FullyQualifiedName() != "" && opts.NewName.FullyQualifiedName() != "" {
		return errors.New("cannot swap and rename in the same ALTER statement")
	}

	return nil
}

type DatabaseSet struct {
	DataRetentionTimeInDays    *int             `ddl:"parameter" db:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int             `ddl:"parameter" db:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDDLCollation        *string          `ddl:"parameter" db:"DEFAULT_DDL_COLLATION"`
	Comment                    *string          `ddl:"parameter,single_quotes" db:"COMMENT"`
	Tag                        []TagAssociation `ddl:"list,no_parentheses" db:"TAG"`
}

type DatabaseUnset struct {
	DataRetentionTimeInDays    *bool        `ddl:"keyword" db:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *bool        `ddl:"keyword" db:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDDLCollation        *bool        `ddl:"keyword" db:"DEFAULT_DDL_COLLATION"`
	Comment                    *bool        `ddl:"keyword" db:"COMMENT"`
	Tag                        []Identifier `ddl:"list,no_parentheses" db:"TAG"`
}

func (v *databases) Alter(ctx context.Context, id AccountLevelIdentifier, opts *DatabaseAlterOptions) error {
	if opts == nil {
		opts = &DatabaseAlterOptions{}
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

type DatabaseAlterReplicationOptions struct {
	alter    bool                   `ddl:"static" db:"ALTER"`    //lint:ignore U1000 This is used in the ddl tag
	database bool                   `ddl:"static" db:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	name    AccountLevelIdentifier `ddl:"identifier"`
	Enable []AccountLevelIdentifier `ddl:"list,no_parentheses" db:"ENABLE REPLICATION TO ACCOUNTS"`
	Disable []AccountLevelIdentifier `ddl:"list,no_parentheses" db:"DISABLE REPLICATION TO ACCOUNTS"`
	Refresh *bool					`ddl:"keyword" db:"REFRESH"`
}

func (opts *DatabaseAlterReplicationOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return errors.New("name must not be empty")
	}

	if len(opts.Enable) > 0 && len(opts.Disable) > 0 {
		return errors.New("cannot enable and disable replication in the same ALTER statement")
	}
	return nil
}

func (v *databases) AlterReplication(ctx context.Context, id AccountLevelIdentifier, opts *DatabaseAlterReplicationOptions) error {
	if opts == nil {
		opts = &DatabaseAlterReplicationOptions{}
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

type DatabaseAlterFailoverOptions struct {
	alter    bool                   `ddl:"static" db:"ALTER"`    //lint:ignore U1000 This is used in the ddl tag
	database bool                   `ddl:"static" db:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	name    AccountLevelIdentifier `ddl:"identifier"`
	Enable []AccountLevelIdentifier `ddl:"list,no_parentheses" db:"ENABLE FAILOVER TO ACCOUNTS"`
	Disable []AccountLevelIdentifier `ddl:"list,no_parentheses" db:"DISABLE FAILOVER TO ACCOUNTS"`
	Primary *bool					`ddl:"keyword" db:"PRIMARY"`
}

func (opts *DatabaseAlterFailoverOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return errors.New("name must not be empty")
	}

	if len(opts.Enable) > 0 && len(opts.Disable) > 0 {
		return errors.New("cannot enable and disable failover in the same ALTER statement")
	}
	return nil
}

func (v *databases) AlterFailover(ctx context.Context, id AccountLevelIdentifier, opts *DatabaseAlterFailoverOptions) error {
	if opts == nil {
		opts = &DatabaseAlterFailoverOptions{}
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

type DatabaseDropOptions struct {
	drop     bool                   `ddl:"static" db:"DROP"`           //lint:ignore U1000 This is used in the ddl tag
	database bool                   `ddl:"static" db:"MASKING POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DatabaseDropOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func (v *databases) Drop(ctx context.Context, id SchemaObjectIdentifier) error {
	// database drop does not support [IF EXISTS] so there are no drop options.
	opts := &DatabaseDropOptions{
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

// DatabaseShowOptions represents the options for listing databases.
type DatabaseShowOptions struct {
	show      bool  `ddl:"static" db:"SHOW"`             //lint:ignore U1000 This is used in the ddl tag
	databases bool  `ddl:"static" db:"MASKING POLICIES"` //lint:ignore U1000 This is used in the ddl tag
	Like      *Like `ddl:"keyword" db:"LIKE"`
	In        *In   `ddl:"keyword" db:"IN"`
	Limit     *int  `ddl:"command,no_quotes" db:"LIMIT"`
}

func (input *DatabaseShowOptions) validate() error {
	return nil
}

// Databases is a user friendly result for a CREATE MASKING POLICY query.
type Database struct {
	CreatedOn           time.Time
	Name                string
	DatabaseName        string
	SchemaName          string
	Kind                string
	Owner               string
	Comment             string
	ExemptOtherPolicies bool
}

func (v *Database) ID() AccountLevelIdentifier {
	return NewAccountLevelIdentifier(v.Name, ObjectTypeDatabase)
}

// databaseDBRow is used to decode the result of a CREATE MASKING POLICY query.
type databaseDBRow struct {
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

func (row databaseDBRow) toDatabase() *Database {
	exemptOtherPolicies, err := jsonparser.GetBoolean([]byte(row.Options), "EXEMPT_OTHER_POLICIES")
	if err != nil {
		exemptOtherPolicies = false
	}
	return &Database{
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

// List all the databases by pattern.
func (v *databases) Show(ctx context.Context, opts *DatabaseShowOptions) ([]*Database, error) {
	if opts == nil {
		opts = &DatabaseShowOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	clauses, err := v.builder.parseStruct(opts)
	if err != nil {
		return nil, err
	}
	stmt := v.builder.sql(clauses...)
	dest := []databaseDBRow{}

	err = v.client.query(ctx, &dest, stmt)
	if err != nil {
		return nil, decodeDriverError(err)
	}
	resultList := make([]*Database, len(dest))
	for i, row := range dest {
		resultList[i] = row.toDatabase()
	}

	return resultList, nil
}

type databaseDescribeOptions struct {
	describe bool                   `ddl:"static" db:"DESCRIBE"`       //lint:ignore U1000 This is used in the ddl tag
	database bool                   `ddl:"static" db:"MASKING POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

func (v *databaseDescribeOptions) validate() error {
	if v.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type DatabaseDetails struct {
	Name       string
	Signature  []TableColumnSignature
	ReturnType DataType
	Body       string
}

type databaseDetailsRow struct {
	Name       string `db:"name"`
	Signature  string `db:"signature"`
	ReturnType string `db:"return_type"`
	Body       string `db:"body"`
}

func (row databaseDetailsRow) toDatabaseDetails() *DatabaseDetails {
	dataType := DataTypeFromString(row.ReturnType)
	v := &DatabaseDetails{
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

func (v *databases) Describe(ctx context.Context, id SchemaObjectIdentifier) (*DatabaseDetails, error) {
	opts := &databaseDescribeOptions{
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
	dest := databaseDetailsRow{}
	err = v.client.queryOne(ctx, &dest, stmt)
	if err != nil {
		return nil, decodeDriverError(err)
	}

	return dest.toDatabaseDetails(), nil
}
