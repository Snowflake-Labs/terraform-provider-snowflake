package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Databases interface {
	// Create creates a database.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateDatabaseOptions) error
	// Alter modifies an existing database
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseOptions) error
	// Drop removes a database.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropDatabaseOptions) error
	// Show returns a list of databases.
	Show(ctx context.Context, opts *ShowDatabasesOptions) ([]*Database, error)
	// ShowByID returns a database by ID
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Database, error)
	// Describe returns the details of a database.
	Describe(ctx context.Context, id AccountObjectIdentifier) (*DatabaseDetails, error)
}

var _ Databases = (*databases)(nil)

type databases struct {
	client *Client
}

type Database struct {
	CreatedOn     time.Time
	Name          string
	IsDefault     string
	IsCurrent     string
	Origin        string
	Owner         string
	Comment       string
	Options       string
	RetentionTime string
	ResourceGroup string
	DroppedOn     time.Time
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
		database.IsDefault = row.IsDefault.String
	}
	if row.IsCurrent.Valid {
		database.IsCurrent = row.IsCurrent.String
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
		database.RetentionTime = row.RetentionTime.String
	}
	if row.ResourceGroup.Valid {
		database.ResourceGroup = row.ResourceGroup.String
	}
	if row.DroppedOn.Valid {
		database.DroppedOn = row.DroppedOn.Time
	}
	return &database
}

// placeholder for the real implementation.
type CreateDatabaseOptions struct{}

func (c *databases) Create(ctx context.Context, id AccountObjectIdentifier, _ *CreateDatabaseOptions) error {
	sql := fmt.Sprintf(`CREATE DATABASE %s`, id.FullyQualifiedName())
	_, err := c.client.exec(ctx, sql)
	return err
}

// placeholder for the real implementation.
type AlterDatabaseOptions struct{}

func (c *databases) Alter(ctx context.Context, id AccountObjectIdentifier, _ *AlterDatabaseOptions) error {
	sql := fmt.Sprintf(`ALTER DATABASE %s`, id.FullyQualifiedName())
	_, err := c.client.exec(ctx, sql)
	return err
}

// placeholder for the real implementation.
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

func (c *databases) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropDatabaseOptions) error {
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
	_, err = c.client.exec(ctx, sql)
	return err
}

// placeholder for the real implementation.
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

func (c *databases) Show(ctx context.Context, opts *ShowDatabasesOptions) ([]*Database, error) {
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
	err = c.client.query(ctx, &rows, sql)
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
	Name string
}

func (c *databases) Describe(ctx context.Context, id AccountObjectIdentifier) (*DatabaseDetails, error) {
	sql := fmt.Sprintf(`DESCRIBE DATABASE %s`, id.FullyQualifiedName())
	var details DatabaseDetails
	err := c.client.queryOne(ctx, &details, sql)
	return &details, err
}

func (v *Database) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *Database) ObjectType() ObjectType {
	return ObjectTypeDatabase
}
