package sdk

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Schemas interface {
	// Create creates a schema.
	Create(ctx context.Context, id SchemaIdentifier, opts *CreateSchemaOptions) error
	// Alter modifies an existing schema.
	Alter(ctx context.Context, id SchemaIdentifier, opts *AlterSchemaOptions) error
	// Drop removes a schema.
	Drop(ctx context.Context, id SchemaIdentifier, opts *DropSchemaOptions) error
	// Undrop restores the most recent version of a dropped schema.
	Undrop(ctx context.Context, id SchemaIdentifier) error
	// Describe lists objects in the schema.
	Describe(ctx context.Context, id SchemaIdentifier) ([]SchemaDetails, error)
	// Show returns a list of schemas.
	Show(ctx context.Context, opts *ShowSchemaOptions) ([]Schema, error)
	// ShowByID returns a schema by ID.
	ShowByID(ctx context.Context, id SchemaIdentifier) (*Schema, error)
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

func (v *Schema) ID() SchemaIdentifier {
	return NewSchemaIdentifier(v.DatabaseName, v.Name)
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
		CreatedOn:    row.CreatedOn,
		Name:         row.Name,
		IsDefault:    row.IsDefault == "Y",
		IsCurrent:    row.IsCurrent == "Y",
		DatabaseName: row.DatabaseName,
		Owner:        row.Owner,
		Comment:      comment,
		Options:      options,
	}
}

type CreateSchemaOptions struct {
	create                     bool             `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace                  *bool            `ddl:"keyword" sql:"OR REPLACE"`
	Transient                  *bool            `ddl:"keyword" sql:"TRANSIENT"`
	schema                     bool             `ddl:"static" sql:"SCHEMA"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists                *bool            `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaIdentifier `ddl:"identifier"` //lint:ignore U1000 This is used in the ddl tag
	Clone                      *Clone           `ddl:"-"`
	WithManagedAccess          *bool            `ddl:"keyword" sql:"WITH MANAGED ACCESS"`
	DataRetentionTimeInDays    *int             `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int             `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDDLCollation        *string          `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	Tag                        []TagAssociation `ddl:"keyword,parentheses" sql:"TAG"`
	Comment                    *string          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreateSchemaOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
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

func (v *schemas) Create(ctx context.Context, id SchemaIdentifier, opts *CreateSchemaOptions) error {
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

type AlterSchemaOptions struct {
	alter                bool             `ddl:"static" sql:"ALTER"`  //lint:ignore U1000 This is used in the ddl tag
	schema               bool             `ddl:"static" sql:"SCHEMA"` //lint:ignore U1000 This is used in the ddl tag
	IfExists             *bool            `ddl:"keyword" sql:"IF EXISTS"`
	name                 SchemaIdentifier `ddl:"identifier"`
	NewName              SchemaIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SwapWith             SchemaIdentifier `ddl:"identifier" sql:"SWAP WITH"`
	Set                  *SchemaSet       `ddl:"list,no_parentheses" sql:"SET"`
	Unset                *SchemaUnset     `ddl:"list,no_parentheses" sql:"UNSET"`
	EnableManagedAccess  *bool            `ddl:"keyword" sql:"ENABLE MANAGED ACCESS"`
	DisableMangaedAccess *bool            `ddl:"keyword" sql:"DISABLE MANAGED ACCESS"`
}

func (opts *AlterSchemaOptions) validate() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.NewName, opts.SwapWith, opts.Set, opts.Unset, opts.EnableManagedAccess, opts.DisableMangaedAccess) {
		errs = append(errs, errors.New("Only one of the fields [ NewName | SwapWith | Set | Unset | EnableManagedAcccess | DisableManagedAcccess ] can be set at once"))
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
	DataRetentionTimeInDays    *int             `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int             `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDDLCollation        *string          `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *string          `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag                        []TagAssociation `ddl:"keyword" sql:"TAG"`
}

func (v *SchemaSet) validate() error {
	if valueSet(v.Tag) && anyValueSet(v.DataRetentionTimeInDays, v.MaxDataExtensionTimeInDays, v.DefaultDDLCollation, v.Comment) {
		return errors.New("Tag cannot be set with other options")
	}
	return nil
}

type SchemaUnset struct {
	DataRetentionTimeInDays    *bool              `ddl:"keyword" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *bool              `ddl:"keyword" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	DefaultDDLCollation        *bool              `ddl:"keyword" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *bool              `ddl:"keyword" sql:"COMMENT"`
	Tag                        []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
}

func (v *SchemaUnset) validate() error {
	if valueSet(v.Tag) && anyValueSet(v.DataRetentionTimeInDays, v.MaxDataExtensionTimeInDays, v.DefaultDDLCollation, v.Comment) {
		return errors.New("Tag cannot be set with other options")
	}
	return nil
}

func (v *schemas) Alter(ctx context.Context, id SchemaIdentifier, opts *AlterSchemaOptions) error {
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

type DropSchemaOptions struct {
	drop     bool             `ddl:"static" sql:"DROP"`
	schema   bool             `ddl:"static" sql:"SCHEMA"`
	IfExists *bool            `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaIdentifier `ddl:"identifier"`
	Cascade  *bool            `ddl:"static" sql:"CASCADE"`
	Restrict *bool            `ddl:"static" sql:"RESTRICT"`
}

func (opts *DropSchemaOptions) validate() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.Cascade, opts.Restrict) {
		errs = append(errs, errors.New("Only one of the fields [ Cascade | Restrict ] can be set at once"))
	}
	return errors.Join(errs...)
}

func (v *schemas) Drop(ctx context.Context, id SchemaIdentifier, opts *DropSchemaOptions) error {
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

type undropSchemaOptions struct {
	undrop bool             `ddl:"static" sql:"UNDROP"` //lint:ignore U1000 This is used in the ddl tag
	schema bool             `ddl:"static" sql:"SCHEMA"` //lint:ignore U1000 This is used in the ddl tag
	name   SchemaIdentifier `ddl:"identifier"`          //lint:ignore U1000 This is used in the ddl tag
}

func (opts *undropSchemaOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *schemas) Undrop(ctx context.Context, id SchemaIdentifier) error {
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

type describeSchemaOptions struct {
	describe bool             `ddl:"static" sql:"DESCRIBE"` //lint:ignore U1000 This is used in the ddl tag
	database bool             `ddl:"static" sql:"SCHEMA"`   //lint:ignore U1000 This is used in the ddl tag
	name     SchemaIdentifier `ddl:"identifier"`            //lint:ignore U1000 This is used in the ddl tag
}

func (opts *describeSchemaOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

type SchemaDetails struct {
	CreatedOn time.Time `db:"created_on"`
	Name      string    `db:"name"`
	Kind      string    `db:"kind"`
}

func (v *schemas) Describe(ctx context.Context, id SchemaIdentifier) ([]SchemaDetails, error) {
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

type InSchema struct {
	Account  *bool                   `ddl:"keyword" sql:"ACCOUNT"`
	Database *bool                   `ddl:"keyword" sql:"DATABASE"`
	Name     AccountObjectIdentifier `ddl:"identifier"`
}

type ShowSchemaOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"` //lint:ignore U1000 This is used in the ddl tag
	Terse      *bool      `ddl:"keyword" sql:"TERSE"`
	schemas    bool       `ddl:"static" sql:"SCHEMAS"` //lint:ignore U1000 This is used in the ddl tag
	History    *bool      `ddl:"keyword" sql:"HISTORY"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	In         *InSchema  `ddl:"keyword" sql:"IN"`
	StartsWith *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	LimitFrom  *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

func (opts *ShowSchemaOptions) validate() error {
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

func (v *schemas) ShowByID(ctx context.Context, id SchemaIdentifier) (*Schema, error) {
	schemas, err := v.client.Schemas.Show(ctx, &ShowSchemaOptions{
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
