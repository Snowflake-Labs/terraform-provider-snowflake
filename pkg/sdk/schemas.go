package sdk

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	_ validatable = new(CreateSchemaOptions)
	_ validatable = new(AlterSchemaOptions)
	_ validatable = new(DropSchemaOptions)
	_ validatable = new(undropSchemaOptions)
	_ validatable = new(describeSchemaOptions)
	_ validatable = new(ShowSchemaOptions)
)

type Schemas interface {
	Create(ctx context.Context, id DatabaseObjectIdentifier, opts *CreateSchemaOptions) error
	Alter(ctx context.Context, id DatabaseObjectIdentifier, opts *AlterSchemaOptions) error
	Drop(ctx context.Context, id DatabaseObjectIdentifier, opts *DropSchemaOptions) error
	Undrop(ctx context.Context, id DatabaseObjectIdentifier) error
	Describe(ctx context.Context, id DatabaseObjectIdentifier) ([]SchemaDetails, error)
	Show(ctx context.Context, opts *ShowSchemaOptions) ([]Schema, error)
	ShowByID(ctx context.Context, id DatabaseObjectIdentifier) (*Schema, error)
	Use(ctx context.Context, id DatabaseObjectIdentifier) error
}

var _ Schemas = (*schemas)(nil)

type schemas struct {
	client *Client
}

type Schema struct {
	CreatedOn     time.Time
	Name          string
	IsDefault     bool
	IsCurrent     bool
	DatabaseName  string
	Owner         string
	Comment       *string
	Options       *string
	RetentionTime string
	OwnerRoleType string
}

func (v *Schema) ID() DatabaseObjectIdentifier {
	return NewDatabaseObjectIdentifier(v.DatabaseName, v.Name)
}

func (v *Schema) ObjectType() ObjectType {
	return ObjectTypeSchema
}

type schemaDBRow struct {
	CreatedOn     time.Time      `db:"created_on"`
	Name          string         `db:"name"`
	IsDefault     string         `db:"is_default"`
	IsCurrent     string         `db:"is_current"`
	DatabaseName  string         `db:"database_name"`
	Owner         string         `db:"owner"`
	Comment       sql.NullString `db:"comment"`
	Options       sql.NullString `db:"options"`
	RetentionTime string         `db:"retention_time"`
	OwnerRoleType string         `db:"owner_role_type"`
}

func (row schemaDBRow) toSchema() Schema {
	var comment *string
	var options *string
	if row.Comment.Valid {
		comment = &row.Comment.String
	}
	if row.Options.Valid {
		options = &row.Options.String
	}
	return Schema{
		CreatedOn:     row.CreatedOn,
		Name:          row.Name,
		IsDefault:     row.IsDefault == "Y",
		IsCurrent:     row.IsCurrent == "Y",
		DatabaseName:  row.DatabaseName,
		Owner:         row.Owner,
		Comment:       comment,
		Options:       options,
		RetentionTime: row.RetentionTime,
		OwnerRoleType: row.OwnerRoleType,
	}
}

// CreateSchemaOptions based on https://docs.snowflake.com/en/sql-reference/sql/create-schema
type CreateSchemaOptions struct {
	create                     bool                     `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	Transient                  *bool                    `ddl:"keyword" sql:"TRANSIENT"`
	schema                     bool                     `ddl:"static" sql:"SCHEMA"`
	IfNotExists                *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       DatabaseObjectIdentifier `ddl:"identifier"`
	Clone                      *Clone                   `ddl:"-"`
	WithManagedAccess          *bool                    `ddl:"keyword" sql:"WITH MANAGED ACCESS"`
	DataRetentionTimeInDays    *int                     `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int                     `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDDLCollation        *string                  `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	Tag                        []TagAssociation         `ddl:"keyword,parentheses" sql:"TAG"`
	Comment                    *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreateSchemaOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Clone) {
		if err := opts.Clone.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateSchemaOptions", "IfNotExists", "OrReplace"))
	}
	return errors.Join(errs...)
}

func (v *schemas) Create(ctx context.Context, id DatabaseObjectIdentifier, opts *CreateSchemaOptions) error {
	if opts == nil {
		opts = &CreateSchemaOptions{}
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

// AlterSchemaOptions based on https://docs.snowflake.com/en/sql-reference/sql/alter-schema
type AlterSchemaOptions struct {
	alter    bool                     `ddl:"static" sql:"ALTER"`
	schema   bool                     `ddl:"static" sql:"SCHEMA"`
	IfExists *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name     DatabaseObjectIdentifier `ddl:"identifier"`
	NewName  DatabaseObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SwapWith DatabaseObjectIdentifier `ddl:"identifier" sql:"SWAP WITH"`
	Set      *SchemaSet               `ddl:"list,no_parentheses" sql:"SET"`
	Unset    *SchemaUnset             `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTag   []TagAssociation         `ddl:"keyword" sql:"SET TAG"`
	UnsetTag []ObjectIdentifier       `ddl:"keyword" sql:"UNSET TAG"`
	// One of
	EnableManagedAccess  *bool `ddl:"keyword" sql:"ENABLE MANAGED ACCESS"`
	DisableManagedAccess *bool `ddl:"keyword" sql:"DISABLE MANAGED ACCESS"`
}

func (opts *AlterSchemaOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.NewName, opts.SwapWith, opts.Set, opts.Unset, opts.SetTag, opts.UnsetTag, opts.EnableManagedAccess, opts.DisableManagedAccess) {
		errs = append(errs, errOneOf("NewName", "SwapWith", "Set", "Unset", "SetTag", "UnsetTag", "EnableManagedAccess", "DisableManagedAccess"))
	}
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type SchemaSet struct {
	DataRetentionTimeInDays    *int    `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int    `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDDLCollation        *string `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (v *SchemaSet) validate() error {
	if !anyValueSet(v.DataRetentionTimeInDays, v.MaxDataExtensionTimeInDays, v.DefaultDDLCollation, v.Comment) {
		return errAtLeastOneOf("SchemaSet", "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "DefaultDDLCollation", "Comment")
	}
	return nil
}

type SchemaUnset struct {
	DataRetentionTimeInDays    *bool `ddl:"keyword" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *bool `ddl:"keyword" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDDLCollation        *bool `ddl:"keyword" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *bool `ddl:"keyword" sql:"COMMENT"`
}

func (v *SchemaUnset) validate() error {
	if !anyValueSet(v.DataRetentionTimeInDays, v.MaxDataExtensionTimeInDays, v.DefaultDDLCollation, v.Comment) {
		return errAtLeastOneOf("SchemaUnset", "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "DefaultDDLCollation", "Comment")
	}
	return nil
}

func (v *schemas) Alter(ctx context.Context, id DatabaseObjectIdentifier, opts *AlterSchemaOptions) error {
	if opts == nil {
		opts = &AlterSchemaOptions{}
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

// DropSchemaOptions Based on https://docs.snowflake.com/en/sql-reference/sql/drop-schema
type DropSchemaOptions struct {
	drop     bool                     `ddl:"static" sql:"DROP"`
	schema   bool                     `ddl:"static" sql:"SCHEMA"`
	IfExists *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name     DatabaseObjectIdentifier `ddl:"identifier"`
	// one of
	Cascade  *bool `ddl:"static" sql:"CASCADE"`
	Restrict *bool `ddl:"static" sql:"RESTRICT"`
}

func (opts *DropSchemaOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.Cascade, opts.Restrict) {
		errs = append(errs, errOneOf("DropSchemaOptions", "Cascade", "Restrict"))
	}
	return errors.Join(errs...)
}

func (v *schemas) Drop(ctx context.Context, id DatabaseObjectIdentifier, opts *DropSchemaOptions) error {
	if opts == nil {
		opts = &DropSchemaOptions{}
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

// undropSchemaOptions is based on https://docs.snowflake.com/en/sql-reference/sql/undrop-schema.
type undropSchemaOptions struct {
	undrop bool                     `ddl:"static" sql:"UNDROP"`
	schema bool                     `ddl:"static" sql:"SCHEMA"`
	name   DatabaseObjectIdentifier `ddl:"identifier"`
}

func (opts *undropSchemaOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *schemas) Undrop(ctx context.Context, id DatabaseObjectIdentifier) error {
	opts := &undropSchemaOptions{
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

// describeSchemaOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-schema.
type describeSchemaOptions struct {
	describe bool                     `ddl:"static" sql:"DESCRIBE"`
	database bool                     `ddl:"static" sql:"SCHEMA"`
	name     DatabaseObjectIdentifier `ddl:"identifier"`
}

func (opts *describeSchemaOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

type SchemaDetails struct {
	CreatedOn time.Time `db:"created_on"`
	Name      string    `db:"name"`
	Kind      string    `db:"kind"`
}

func (v *schemas) Describe(ctx context.Context, id DatabaseObjectIdentifier) ([]SchemaDetails, error) {
	opts := &describeSchemaOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var details []SchemaDetails
	err = v.client.query(ctx, &details, sql)
	if err != nil {
		return nil, err
	}
	return details, err
}

type SchemaIn struct {
	Account  *bool                   `ddl:"keyword" sql:"ACCOUNT"`
	Database *bool                   `ddl:"keyword" sql:"DATABASE"`
	Name     AccountObjectIdentifier `ddl:"identifier"`
}

// ShowSchemaOptions based on https://docs.snowflake.com/en/sql-reference/sql/show-schemas
type ShowSchemaOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	Terse      *bool      `ddl:"keyword" sql:"TERSE"`
	schemas    bool       `ddl:"static" sql:"SCHEMAS"`
	History    *bool      `ddl:"keyword" sql:"HISTORY"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	In         *SchemaIn  `ddl:"keyword" sql:"IN"`
	StartsWith *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	LimitFrom  *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

func (opts *ShowSchemaOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (v *schemas) Show(ctx context.Context, opts *ShowSchemaOptions) ([]Schema, error) {
	if opts == nil {
		opts = &ShowSchemaOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []schemaDBRow
	err = v.client.query(ctx, &rows, sql)
	schemas := make([]Schema, len(rows))
	for i, row := range rows {
		schemas[i] = row.toSchema()
	}
	return schemas, err
}

func (v *schemas) ShowByID(ctx context.Context, id DatabaseObjectIdentifier) (*Schema, error) {
	schemas, err := v.client.Schemas.Show(ctx, &ShowSchemaOptions{
		In: &SchemaIn{
			Database: Bool(true),
			Name:     NewAccountObjectIdentifier(id.DatabaseName()),
		},
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	for _, s := range schemas {
		if s.ID() == id {
			return &s, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

func (v *schemas) Use(ctx context.Context, id DatabaseObjectIdentifier) error {
	return v.client.Sessions.UseSchema(ctx, id)
}
