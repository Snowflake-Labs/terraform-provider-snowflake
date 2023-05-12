package sdk

import (
	"context"
	"fmt"
)

type Warehouses interface {
	// Create creates a warehouse.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *WarehouseCreateOptions) error
	// Alter modifies an existing warehouse
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *WarehouseAlterOptions) error
	// Drop removes a warehouse.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *WarehouseDropOptions) error
	// Show returns a list of warehouses.
	Show(ctx context.Context, opts *WarehouseShowOptions) ([]*Warehouse, error)
	// Describe returns the details of a warehouse.
	Describe(ctx context.Context, id AccountObjectIdentifier) (*WarehouseDetails, error)
}

var _ Warehouses = (*warehouses)(nil)

type warehouses struct {
	client  *Client
	builder *sqlBuilder
}

type Warehouse struct {
	Name string
}

// placeholder for the real implementation.
type WarehouseCreateOptions struct{}

func (c *warehouses) Create(ctx context.Context, id AccountObjectIdentifier, _ *WarehouseCreateOptions) error {
	sql := fmt.Sprintf(`CREATE WAREHOUSE %s`, id.FullyQualifiedName())
	_, err := c.client.exec(ctx, sql)
	return err
}

// placeholder for the real implementation.
type WarehouseAlterOptions struct{}

func (c *warehouses) Alter(ctx context.Context, id AccountObjectIdentifier, _ *WarehouseAlterOptions) error {
	sql := fmt.Sprintf(`ALTER WAREHOUSE %s`, id.FullyQualifiedName())
	_, err := c.client.exec(ctx, sql)
	return err
}

// placeholder for the real implementation.
type WarehouseDropOptions struct {
	drop      bool                   `ddl:"static" db:"DROP"`      //lint:ignore U1000 This is used in the ddl tag
	warehouse bool                   `ddl:"static" db:"WAREHOUSE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists  *bool                   `ddl:"keyword" db:"IF EXISTS"`
	name      AccountObjectIdentifier `ddl:"identifier"` //lint:ignore U1000 This is used in the ddl tag
}

func (opts *WarehouseDropOptions) validate() error {
	if opts.name.FullyQualifiedName() == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func (c *warehouses) Drop(ctx context.Context, id AccountObjectIdentifier, opts *WarehouseDropOptions) error {
	if opts == nil {
		opts = &WarehouseDropOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	clauses, err := c.builder.parseStruct(opts)
	if err != nil {
		return err
	}
	sql := c.builder.sql(clauses...)
	_, err = c.client.exec(ctx, sql)
	return err
}

// placeholder for the real implementation.
type WarehouseShowOptions struct{}

func (c *warehouses) Show(ctx context.Context, _ *WarehouseShowOptions) ([]*Warehouse, error) {
	sql := `SHOW WAREHOUSES`
	var warehouses []*Warehouse
	err := c.client.query(ctx, &warehouses, sql)
	return warehouses, err
}

type WarehouseDetails struct {
	Name string
}

func (c *warehouses) Describe(ctx context.Context, id AccountObjectIdentifier) (*WarehouseDetails, error) {
	sql := fmt.Sprintf(`DESCRIBE WAREHOUSE %s`, id.FullyQualifiedName())
	var details WarehouseDetails
	err := c.client.queryOne(ctx, &details, sql)
	return &details, err
}

func (v *Warehouse) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}
