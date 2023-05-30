package sdk

import (
	"context"
	"fmt"
)

type Databases interface {
	// Create creates a database.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateDatabaseOptions) error
	// Alter modifies an existing database
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseOptions) error
	// Drop removes a database.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropDatabaseOptions) error
	// Show returns a list of databases.
	Show(ctx context.Context, opts *ShowDatabaseOptions) ([]*Database, error)
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
	Name string
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
type ShowDatabaseOptions struct{}

func (c *databases) Show(ctx context.Context, _ *ShowDatabaseOptions) ([]*Database, error) {
	sql := `SHOW DATABASES`
	var databases []*Database
	err := c.client.query(ctx, &databases, sql)
	return databases, err
}

func (c *databases) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Database, error) {
	sql := fmt.Sprintf(`SHOW DATABASES LIKE '%s'`, id.Name())
	var database Database
	err := c.client.queryOne(ctx, &database, sql)
	return &database, err
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
