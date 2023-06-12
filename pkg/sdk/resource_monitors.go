package sdk

import (
	"context"
	"errors"
	"time"
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
	create          bool                    `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace       *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"` //lint:ignore U1000 This is used in the ddl tag
	name            AccountObjectIdentifier `ddl:"identifier"`
	with            bool                    `ddl:"static" sql:"WITH"` //lint:ignore U1000 This is used in the ddl tag

	//optional, at least one
	CreditQuota    *int                 `ddl:"parameter,equals" sql:"CREDIT_QUOTA"`
	Frequency      *Frequency           `ddl:"parameter,equals" sql:"FREQUENCY"`
	StartTimeStamp *string              `ddl:"parameter,equals" sql:"START_TIMESTAMP"`
	EndTimeStamp   *string              `ddl:"parameter,equals" sql:"END_TIMESTAMP"`
	NotifyUsers    *NotifyUsers         `ddl:"parameter,equals" sql:"NOTIFY_USERS"`
	Triggers       *[]TriggerDefinition `ddl:"keyword,no_comma" sql:"TRIGGERS"`
}

func (opts *CreateResourceMonitorOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}

	if opts == nil || everyValueNil(
		opts.CreditQuota,
		opts.Frequency,
		opts.StartTimeStamp,
		opts.EndTimeStamp,
		opts.NotifyUsers,
		opts.Triggers,
	) {
		return errors.New("No alter action specified")
	}
	if (opts.Frequency != nil && opts.StartTimeStamp == nil) || (opts.Frequency == nil && opts.StartTimeStamp != nil) {
		return errors.New("must specify frequency and start time together")
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

type TriggerDefinition struct {
	on            bool          `ddl:"static" sql:"ON"`      //lint:ignore U1000 This is used in the ddl tag
	Threshold     int           `ddl:"keyword" `             //lint:ignore U1000 This is used in the ddl tag
	percent       bool          `ddl:"static" sql:"PERCENT"` //lint:ignore U1000 This is used in the ddl tag
	do            bool          `ddl:"static" sql:"DO"`      //lint:ignore U1000 This is used in the ddl tag
	TriggerAction triggerAction `ddl:"keyword" `             //lint:ignore U1000 This is used in the ddl tag

}

type triggerAction string

const (
	Suspend          triggerAction = "SUSPEND"
	SuspendImmediate triggerAction = "SUSPEND_IMMEDIATE"
	Notify           triggerAction = "NOTIFY"
)

type NotifyUsers struct {
	//tutaj walidacja zwykla chyba, a nie w typie przekazane ze jest min 1
	Users []string `ddl:"list,parentheses,no_comma"`
}

type StartTimeStamp interface {
	startTimeStamp()
	String() string
}
type Immediately struct{}

func (Immediately) startTimeStamp() {}
func (Immediately) String() string {
	return "IMMEDIATELY"
}

// time na pewno do poprawki, zobacz jak te timestampy ograć
type TimeStamp time.Time

func (TimeStamp) startTimeStamp() {}

func (ts TimeStamp) String() string {
	return (time.Time)(ts).String()

}

type Frequency string

const (
	Monthly Frequency = "MONTHLY"
	Daily   Frequency = "DAILY"
	Weekly  Frequency = "WEEKLY"
	Yearly  Frequency = "YEARLY"
	Never   Frequency = "NEVER"
)

// AlterResourceMonitorOptions contains options for altering a resource monitor.
type AlterResourceMonitorOptions struct {
	alter           bool                    `ddl:"static" sql:"ALTER"`            //lint:ignore U1000 This is used in the ddl tag
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"` //lint:ignore U1000 This is used in the ddl tag
	IfExists        *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	Set             *ResourceMonitorSet     `ddl:"keyword" sql:"SET"`
}

func (opts *AlterResourceMonitorOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if opts == nil || opts.Set == nil || everyValueNil(
		opts.Set.CreditQuota,
		opts.Set.Frequency,
		opts.Set.StartTimeStamp,
		opts.Set.EndTimeStamp,
		opts.Set.NotifyUsers,
		opts.Set.Triggers,
	) {
		return errors.New("No alter action specified")
	}
	if (opts.Set.Frequency != nil && opts.Set.StartTimeStamp == nil) || (opts.Set.Frequency == nil && opts.Set.StartTimeStamp != nil) {
		return errors.New("must specify frequency and start time together")
	}

	return nil
}

func (v *resourceMonitors) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterResourceMonitorOptions) error {

	return nil
}

type ResourceMonitorSet struct {
	//at least one
	CreditQuota    *int                 `ddl:"parameter,equals" sql:"CREDIT_QUOTA"`
	Frequency      *Frequency           `ddl:"parameter,equals" sql:"FREQUENCY"`
	StartTimeStamp *string              `ddl:"parameter,equals" sql:"START_TIMESTAMP"`
	EndTimeStamp   *string              `ddl:"parameter,equals" sql:"END_TIMESTAMP"`
	NotifyUsers    *NotifyUsers         `ddl:"parameter,equals" sql:"NOTIFY_USERS"`
	Triggers       *[]TriggerDefinition `ddl:"keyword,no_comma" sql:"TRIGGERS"`
}

// resourceMonitorDropOptions contains options for dropping a resource monitor.
type dropResourceMonitorOptions struct {
	drop            bool                    `ddl:"static" sql:"DROP"`             //lint:ignore U1000 This is used in the ddl tag
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"` //lint:ignore U1000 This is used in the ddl tag
	name            AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *dropResourceMonitorOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *resourceMonitors) Drop(ctx context.Context, id AccountObjectIdentifier) error {
	opts := &dropResourceMonitorOptions{
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
