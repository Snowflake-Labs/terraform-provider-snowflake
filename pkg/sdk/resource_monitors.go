package sdk

import (
	"context"
)

type ResourceMonitors interface {
	// Create creates a resource monitor.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateResourceMonitorOptions) error
	// Alter modifies an existing resource monitor
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterResourceMonitorOptions) error
	// Drop removes a resource monitor.
	Drop(ctx context.Context, id AccountObjectIdentifier) error
	// Show returns a list of resource monitor.
	Show(ctx context.Context, opts *ShowResourceMonitorOptions) ([]*ResourceMonitor, error)
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

// CreateResourceMonitorOptions contains options for creating a resource monitor.
type CreateResourceMonitorOptions struct {
	create          bool                    `ddl:"static" sql:"CREATE"`           //lint:ignore U1000 This is used in the ddl tag
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"` //lint:ignore U1000 This is used in the ddl tag
	name            AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *CreateResourceMonitorOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *resourceMonitors) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateResourceMonitorOptions) error {
	if opts == nil {
		opts = &CreateResourceMonitorOptions{}
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

// AlterResourceMonitorOptions contains options for altering a resource monitor.
type AlterResourceMonitorOptions struct{}

func (v *resourceMonitors) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterResourceMonitorOptions) error {
	return nil
}

// resourceMonitorDropOptions contains options for dropping a resource monitor.
type resourceMonitorDropOptions struct {
	drop            bool                    `ddl:"static" sql:"DROP"`             //lint:ignore U1000 This is used in the ddl tag
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"` //lint:ignore U1000 This is used in the ddl tag
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

// ShowResourceMonitorOptions contains options for listing resource monitors.
type ShowResourceMonitorOptions struct {
	show             bool  `ddl:"static" sql:"SHOW"`              //lint:ignore U1000 This is used in the ddl tag
	resourceMonitors bool  `ddl:"static" sql:"RESOURCE MONITORS"` //lint:ignore U1000 This is used in the ddl tag
	Like             *Like `ddl:"keyword" sql:"LIKE"`
}

func (opts *ShowResourceMonitorOptions) validate() error {
	return nil
}

func (v *resourceMonitors) Show(ctx context.Context, opts *ShowResourceMonitorOptions) ([]*ResourceMonitor, error) {
	if opts == nil {
		opts = &ShowResourceMonitorOptions{}
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
	resourceMonitors, err := v.Show(ctx, &ShowResourceMonitorOptions{
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
