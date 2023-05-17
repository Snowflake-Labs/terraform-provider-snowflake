package sdk

import (
	"context"
	"fmt"
)

type Databases interface {
	// Create creates a database.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *DatabaseCreateOptions) error
	// Alter modifies an existing database
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *DatabaseAlterOptions) error
	// Drop removes a database.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DatabaseDropOptions) error
	// Show returns a list of databases.
	Show(ctx context.Context, opts *DatabaseShowOptions) ([]*Database, error)
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
type DatabaseCreateOptions struct{}

func (c *databases) Create(ctx context.Context, id AccountObjectIdentifier, _ *DatabaseCreateOptions) error {
	sql := fmt.Sprintf(`CREATE DATABASE %s`, id.FullyQualifiedName())
	_, err := c.client.exec(ctx, sql)
	return err
}

// placeholder for the real implementation.
type DatabaseAlterOptions struct{}

func (c *databases) Alter(ctx context.Context, id AccountObjectIdentifier, _ *DatabaseAlterOptions) error {
	sql := fmt.Sprintf(`ALTER DATABASE %s`, id.FullyQualifiedName())
	_, err := c.client.exec(ctx, sql)
	return err
}

// placeholder for the real implementation.
type DatabaseDropOptions struct {
	drop     bool                    `ddl:"static" db:"DROP"`     //lint:ignore U1000 This is used in the ddl tag
	database bool                    `ddl:"static" db:"DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                   `ddl:"keyword" db:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"` //lint:ignore U1000 This is used in the ddl tag
}

func (opts *DatabaseDropOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (c *databases) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DatabaseDropOptions) error {
	if opts == nil {
		opts = &DatabaseDropOptions{}
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
type DatabaseShowOptions struct{}

func (c *databases) Show(ctx context.Context, _ *DatabaseShowOptions) ([]*Database, error) {
	sql := `SHOW DATABASES`
	var databases []*Database
	err := c.client.query(ctx, &databases, sql)
	return databases, err
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
