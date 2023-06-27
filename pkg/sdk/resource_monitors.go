package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	Name                     string
	CreditQuota              float64
	UsedCredits              float64
	RemainingCredits         float64
	Frequency                Frequency
	StartTime                string
	EndTime                  string
	SuspendTriggers          []TriggerDefinition
	SuspendImmediateTriggers []TriggerDefinition
	NotifyTriggers           []TriggerDefinition
	Level                    ResourceMonitorLevel
	Comment                  string
	NotifyUsers              []string
}

type resourceMonitorRow struct {
	Name               string         `db:"name"`
	CreditQuota        sql.NullString `db:"credit_quota"`
	UsedCredits        sql.NullString `db:"used_credits"`
	RemainingCredits   sql.NullString `db:"remaining_credits"`
	Level              sql.NullString `db:"level"`
	Frequency          sql.NullString `db:"frequency"`
	StartTime          sql.NullString `db:"start_time"`
	EndTime            sql.NullString `db:"end_time"`
	NotifyAt           sql.NullString `db:"notify_at"`
	SuspendAt          sql.NullString `db:"suspend_at"`
	SuspendImmediateAt sql.NullString `db:"suspend_immediately_at"`
	Owner              sql.NullString `db:"owner"`
	Comment            sql.NullString `db:"comment"`
	NotifyUsers        sql.NullString `db:"notify_users"`
}

func (row *resourceMonitorRow) toResourceMonitor() (*ResourceMonitor, error) {
	resourceMonitor := &ResourceMonitor{
		Name: row.Name,
	}
	if row.CreditQuota.Valid {
		creditQuota, err := strconv.ParseFloat(row.CreditQuota.String, 64)
		if err != nil {
			return nil, err
		}
		resourceMonitor.CreditQuota = creditQuota
	}
	if row.UsedCredits.Valid {
		usedCredits, err := strconv.ParseFloat(row.UsedCredits.String, 64)
		if err != nil {
			return nil, err
		}
		resourceMonitor.UsedCredits = usedCredits
	}
	if row.RemainingCredits.Valid {
		remainingCredits, err := strconv.ParseFloat(row.RemainingCredits.String, 64)
		if err != nil {
			return nil, err
		}
		resourceMonitor.RemainingCredits = remainingCredits
	}

	if row.Frequency.Valid {
		frequency, err := FrequencyFromString(row.Frequency.String)
		if err != nil {
			return nil, err
		}
		resourceMonitor.Frequency = *frequency
	}
	if row.StartTime.Valid {
		resourceMonitor.StartTime = row.StartTime.String
	}
	if row.EndTime.Valid {
		resourceMonitor.EndTime = row.EndTime.String
	}
	suspendTriggers, err := extractTriggers(row.SuspendAt, TriggerActionSuspend)
	if err != nil {
		return nil, err
	}
	resourceMonitor.SuspendTriggers = suspendTriggers
	suspendImmediateTriggers, err := extractTriggers(row.SuspendImmediateAt, TriggerActionSuspendImmediate)
	if err != nil {
		return nil, err
	}
	resourceMonitor.SuspendImmediateTriggers = suspendImmediateTriggers
	notifyTriggers, err := extractTriggers(row.NotifyAt, TriggerActionNotify)
	if err != nil {
		return nil, err
	}
	resourceMonitor.NotifyTriggers = notifyTriggers
	if row.Comment.Valid {
		resourceMonitor.Comment = row.Comment.String
	}
	resourceMonitor.NotifyUsers = extractUsers(row.NotifyUsers)

	if row.Level.Valid {
		switch row.Level.String {
		case "ACCOUNT":
			resourceMonitor.Level = ResourceMonitorLevelAccount
		case "WAREHOUSE":
			resourceMonitor.Level = ResourceMonitorLevelWarehouse
		default:
			resourceMonitor.Level = ResourceMonitorLevelNull
		}
	} else {
		resourceMonitor.Level = ResourceMonitorLevelNull
	}

	return resourceMonitor, nil
}

// extractTriggerInts converts the triggers in the DB (stored as a comma
// separated string with trailing %s) into a slice of TriggerDefinitions.
func extractTriggers(s sql.NullString, trigger TriggerAction) ([]TriggerDefinition, error) {
	// Check if this is NULL
	if !s.Valid {
		return []TriggerDefinition{}, nil
	}
	ints := strings.Split(s.String, ",")
	out := make([]TriggerDefinition, 0, len(ints))
	for _, i := range ints {
		threshold, err := strconv.Atoi(i[:len(i)-1])
		if err != nil {
			return out, fmt.Errorf("failed to convert %v to integer err = %w", i, err)
		}
		out = append(out, TriggerDefinition{Threshold: threshold, TriggerAction: trigger})
	}
	return out, nil
}

func extractUsers(s sql.NullString) []string {
	if s.Valid && s.String != "" {
		return strings.Split(s.String, ", ")
	} else {
		return []string{}
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
	With            *bool                   `ddl:"keyword" sql:"WITH"`

	CreditQuota    *int                 `ddl:"parameter,equals" sql:"CREDIT_QUOTA"`
	Frequency      *Frequency           `ddl:"parameter,equals" sql:"FREQUENCY"`
	StartTimeStamp *string              `ddl:"parameter,equals,single_quotes" sql:"START_TIMESTAMP"`
	EndTimeStamp   *string              `ddl:"parameter,equals,single_quotes" sql:"END_TIMESTAMP"`
	NotifyUsers    *NotifyUsers         `ddl:"parameter,equals" sql:"NOTIFY_USERS"`
	Triggers       *[]TriggerDefinition `ddl:"keyword,no_comma" sql:"TRIGGERS"`
}

func (opts *CreateResourceMonitorOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *resourceMonitors) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateResourceMonitorOptions) error {
	if opts == nil || everyValueNil(
		opts.CreditQuota,
		opts.Frequency,
		opts.StartTimeStamp,
		opts.EndTimeStamp,
		opts.NotifyUsers,
		opts.Triggers,
	) {
		opts = &CreateResourceMonitorOptions{}
		opts.With = Bool(false)
	} else {
		opts.With = Bool(true)
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

type ResourceMonitorLevel int

const (
	ResourceMonitorLevelAccount = iota
	ResourceMonitorLevelWarehouse
	ResourceMonitorLevelNull
)

type TriggerDefinition struct {
	Threshold     int           `ddl:"parameter,no_equals" sql:"ON"`
	TriggerAction TriggerAction `ddl:"parameter,no_equals" sql:"PERCENT DO"`
}

type TriggerAction string

const (
	TriggerActionSuspend          TriggerAction = "SUSPEND"
	TriggerActionSuspendImmediate TriggerAction = "SUSPEND_IMMEDIATE"
	TriggerActionNotify           TriggerAction = "NOTIFY"
)

type NotifyUsers struct {
	Users []NotifiedUser `ddl:"list,parentheses,comma"`
}

type NotifiedUser struct {
	Name string `ddl:"keyword,double_quotes"`
}

type Frequency string

func FrequencyFromString(s string) (*Frequency, error) {
	s = strings.ToUpper(s)
	if monthly := string(FrequencyMonthly); monthly == s {
		return (*Frequency)(&monthly), nil
	}
	if daily := string(FrequencyDaily); daily == s {
		return (*Frequency)(&daily), nil
	}
	if weekly := string(FrequencyWeekly); weekly == s {
		return (*Frequency)(&weekly), nil
	}
	if yearly := string(FrequencyYearly); yearly == s {
		return (*Frequency)(&yearly), nil
	}
	if never := string(FrequencyNever); never == s {
		return (*Frequency)(&never), nil
	}
	return nil, fmt.Errorf("Invalid frequency type: %s", s)
}

const (
	FrequencyMonthly Frequency = "MONTHLY"
	FrequencyDaily   Frequency = "DAILY"
	FrequencyWeekly  Frequency = "WEEKLY"
	FrequencyYearly  Frequency = "YEARLY"
	FrequencyNever   Frequency = "NEVER"
)

// AlterResourceMonitorOptions contains options for altering a resource monitor.
type AlterResourceMonitorOptions struct {
	alter           bool                    `ddl:"static" sql:"ALTER"`            //lint:ignore U1000 This is used in the ddl tag
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"` //lint:ignore U1000 This is used in the ddl tag
	IfExists        *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	Set             *ResourceMonitorSet     `ddl:"keyword" sql:"SET"`
	NotifyUsers     *NotifyUsers            `ddl:"parameter,equals" sql:"NOTIFY_USERS"`
	Triggers        *[]TriggerDefinition    `ddl:"keyword,no_comma" sql:"TRIGGERS"`
}

func (opts *AlterResourceMonitorOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if opts.Set == nil {
		return nil
	}
	if (opts.Set.Frequency != nil && opts.Set.StartTimeStamp == nil) || (opts.Set.Frequency == nil && opts.Set.StartTimeStamp != nil) {
		return errors.New("must specify frequency and start time together")
	}

	return nil
}

func (v *resourceMonitors) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterResourceMonitorOptions) error {
	if opts == nil {
		opts = &AlterResourceMonitorOptions{}
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

type ResourceMonitorSet struct {
	// at least one
	CreditQuota    *int       `ddl:"parameter,equals" sql:"CREDIT_QUOTA"`
	Frequency      *Frequency `ddl:"parameter,equals" sql:"FREQUENCY"`
	StartTimeStamp *string    `ddl:"parameter,equals,single_quotes" sql:"START_TIMESTAMP"`
	EndTimeStamp   *string    `ddl:"parameter,equals,single_quotes" sql:"END_TIMESTAMP"`
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
		resourceMonitor, err := row.toResourceMonitor()
		if err != nil {
			return nil, err
		}
		resourceMonitors = append(resourceMonitors, resourceMonitor)
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
