package sdk

import (
	"context"
)

type ResourceMonitors interface {
	// Create creates a resource monitor.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *ResourceMonitorCreateOptions) error
	// Alter modifies an existing resource monitor
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *ResourceMonitorAlterOptions) error
	// Drop removes a resource monitor.
	Drop(ctx context.Context, id AccountObjectIdentifier) error
	// Show returns a list of resource monitor.
	Show(ctx context.Context, opts *ResourceMonitorShowOptions) ([]*ResourceMonitor, error)
	// ShowByID returns a resource monitor by ID
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ResourceMonitor, error)
}

var _ ResourceMonitors = (*resourceMonitors)(nil)

type resourceMonitors struct {
	client *Client
}

type ResourceMonitor struct {
	Name string
}

type resourceMonitorRow struct {
	Name string `db:"name"`
}

func (row *resourceMonitorRow) toResourceMonitor() *ResourceMonitor {
	return &ResourceMonitor{
		Name: row.Name,
	}
}

func (v *ResourceMonitor) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *ResourceMonitor) ObjectType() ObjectType {
	return ObjectTypeResourceMonitor
}

// ResourceMonitorCreateOptions contains options for creating a resource monitor.
type ResourceMonitorCreateOptions struct {
	create          bool                    `ddl:"static" db:"CREATE"`           //lint:ignore U1000 This is used in the ddl tag
	resourceMonitor bool                    `ddl:"static" db:"RESOURCE MONITOR"` //lint:ignore U1000 This is used in the ddl tag
	name            AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *ResourceMonitorCreateOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *resourceMonitors) Create(ctx context.Context, id AccountObjectIdentifier, opts *ResourceMonitorCreateOptions) error {
	if opts == nil {
		opts = &ResourceMonitorCreateOptions{}
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

// ResourceMonitorAlterOptions contains options for altering a resource monitor.
type ResourceMonitorAlterOptions struct{}

func (v *resourceMonitors) Alter(ctx context.Context, id AccountObjectIdentifier, opts *ResourceMonitorAlterOptions) error {
	return nil
}

// resourceMonitorDropOptions contains options for dropping a resource monitor.
type resourceMonitorDropOptions struct {
	drop            bool                    `ddl:"static" db:"DROP"`             //lint:ignore U1000 This is used in the ddl tag
	resourceMonitor bool                    `ddl:"static" db:"RESOURCE MONITOR"` //lint:ignore U1000 This is used in the ddl tag
	name            AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *resourceMonitorDropOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *resourceMonitors) Drop(ctx context.Context, id AccountObjectIdentifier) error {
	opts := &resourceMonitorDropOptions{
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

// ResourceMonitorShowOptions contains options for listing resource monitors.
type ResourceMonitorShowOptions struct {
	show             bool  `ddl:"static" db:"SHOW"`              //lint:ignore U1000 This is used in the ddl tag
	resourceMonitors bool  `ddl:"static" db:"RESOURCE MONITORS"` //lint:ignore U1000 This is used in the ddl tag
	Like             *Like `ddl:"keyword" db:"LIKE"`
}

func (opts *ResourceMonitorShowOptions) validate() error {
	return nil
}

func (v *resourceMonitors) Show(ctx context.Context, opts *ResourceMonitorShowOptions) ([]*ResourceMonitor, error) {
	if opts == nil {
		opts = &ResourceMonitorShowOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []*resourceMonitorRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	resourceMonitors := make([]*ResourceMonitor, 0, len(rows))
	for _, row := range rows {
		resourceMonitors = append(resourceMonitors, row.toResourceMonitor())
	}
	return resourceMonitors, nil
}

func (v *resourceMonitors) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ResourceMonitor, error) {
	resourceMonitors, err := v.Show(ctx, &ResourceMonitorShowOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	for _, resourceMonitor := range resourceMonitors {
		if resourceMonitor.Name == id.Name() {
			return resourceMonitor, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}
