package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	_ validatable = new(CreateResourceMonitorOptions)
	_ validatable = new(AlterResourceMonitorOptions)
	_ validatable = new(dropResourceMonitorOptions)
	_ validatable = new(ShowResourceMonitorOptions)
)

type ResourceMonitors interface {
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateResourceMonitorOptions) error
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterResourceMonitorOptions) error
	Drop(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, opts *ShowResourceMonitorOptions) ([]ResourceMonitor, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ResourceMonitor, error)
}

var _ ResourceMonitors = (*resourceMonitors)(nil)

type resourceMonitors struct {
	client *Client
}

type ResourceMonitor struct {
	Name               string
	CreditQuota        float64
	UsedCredits        float64
	RemainingCredits   float64
	Frequency          Frequency
	StartTime          string
	EndTime            string
	SuspendAt          *int
	SuspendImmediateAt *int
	NotifyTriggers     []int
	Level              ResourceMonitorLevel
	Comment            string
	NotifyUsers        []string
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

func (row *resourceMonitorRow) convert() (*ResourceMonitor, error) {
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
		const YYMMDDhhmm = "2006-01-02 15:04"
		t, err := time.Parse(time.RFC3339, row.StartTime.String)
		if err != nil {
			return nil, err
		}
		localTime := t.Local()
		resourceMonitor.StartTime = localTime.Format(YYMMDDhhmm)
	}
	if row.EndTime.Valid {
		resourceMonitor.EndTime = row.EndTime.String
	}
	suspendTriggers, err := extractTriggerInts(row.SuspendAt)
	if err != nil {
		return nil, err
	}
	if len(suspendTriggers) > 0 {
		resourceMonitor.SuspendAt = &suspendTriggers[0]
	}
	suspendImmediateTriggers, err := extractTriggerInts(row.SuspendImmediateAt)
	if err != nil {
		return nil, err
	}
	if len(suspendImmediateTriggers) > 0 {
		resourceMonitor.SuspendImmediateAt = &suspendImmediateTriggers[0]
	}
	notifyTriggers, err := extractTriggerInts(row.NotifyAt)
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
// separated string with trailing %s) into a slice of ints.
func extractTriggerInts(s sql.NullString) ([]int, error) {
	// Check if this is NULL
	if !s.Valid {
		return []int{}, nil
	}
	ints := strings.Split(s.String, ",")
	out := make([]int, 0, len(ints))
	for _, i := range ints {
		myInt, err := strconv.Atoi(i[:len(i)-1])
		if err != nil {
			return out, fmt.Errorf("failed to convert %v to integer err = %w", i, err)
		}
		out = append(out, myInt)
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

// CreateResourceMonitorOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-resource-monitor.
type CreateResourceMonitorOptions struct {
	create          bool                    `ddl:"static" sql:"CREATE"`
	OrReplace       *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	With            *ResourceMonitorWith    `ddl:"keyword" sql:"WITH"`
}

type ResourceMonitorWith struct {
	CreditQuota    *int                `ddl:"parameter,equals" sql:"CREDIT_QUOTA"`
	Frequency      *Frequency          `ddl:"parameter,equals" sql:"FREQUENCY"`
	StartTimestamp *string             `ddl:"parameter,equals,single_quotes" sql:"START_TIMESTAMP"`
	EndTimestamp   *string             `ddl:"parameter,equals,single_quotes" sql:"END_TIMESTAMP"`
	NotifyUsers    *NotifyUsers        `ddl:"parameter,equals" sql:"NOTIFY_USERS"`
	Triggers       []TriggerDefinition `ddl:"keyword,no_comma" sql:"TRIGGERS"`
}

func (opts *CreateResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
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
	f := Frequency(s)
	switch f {
	case FrequencyDaily, FrequencyWeekly, FrequencyMonthly, FrequencyYearly, FrequencyNever:
		return &f, nil
	default:
		return nil, fmt.Errorf("invalid frequency type: %s", s)
	}
}

const (
	FrequencyMonthly Frequency = "MONTHLY"
	FrequencyDaily   Frequency = "DAILY"
	FrequencyWeekly  Frequency = "WEEKLY"
	FrequencyYearly  Frequency = "YEARLY"
	FrequencyNever   Frequency = "NEVER"
)

// AlterResourceMonitorOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-resource-monitor.
type AlterResourceMonitorOptions struct {
	alter           bool                    `ddl:"static" sql:"ALTER"`
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"`
	IfExists        *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	Set             *ResourceMonitorSet     `ddl:"keyword" sql:"SET"`
	NotifyUsers     *NotifyUsers            `ddl:"parameter,equals" sql:"NOTIFY_USERS"`
	Triggers        []TriggerDefinition     `ddl:"keyword,no_comma" sql:"TRIGGERS"`
}

func (opts *AlterResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Set) {
		if (opts.Set.Frequency != nil && opts.Set.StartTimestamp == nil) || (opts.Set.Frequency == nil && opts.Set.StartTimestamp != nil) {
			errs = append(errs, errors.New("must specify frequency and start time together"))
		}
	}
	if !exactlyOneValueSet(opts.Set, opts.NotifyUsers) && opts.Triggers == nil {
		errs = append(errs, errExactlyOneOf("AlterResourceMonitorOptions", "Set", "NotifyUsers", "Triggers"))
	}
	return errors.Join(errs...)
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
	StartTimestamp *string    `ddl:"parameter,equals,single_quotes" sql:"START_TIMESTAMP"`
	EndTimestamp   *string    `ddl:"parameter,equals,single_quotes" sql:"END_TIMESTAMP"`
}

// dropResourceMonitorOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-resource-monitor.
type dropResourceMonitorOptions struct {
	drop            bool                    `ddl:"static" sql:"DROP"`
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"`
	name            AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *dropResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
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

// ShowResourceMonitorOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-resource-monitors.
type ShowResourceMonitorOptions struct {
	show             bool  `ddl:"static" sql:"SHOW"`
	resourceMonitors bool  `ddl:"static" sql:"RESOURCE MONITORS"`
	Like             *Like `ddl:"keyword" sql:"LIKE"`
}

func (opts *ShowResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (v *resourceMonitors) Show(ctx context.Context, opts *ShowResourceMonitorOptions) ([]ResourceMonitor, error) {
	opts = createIfNil(opts)
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
	resourceMonitors := make([]ResourceMonitor, 0, len(rows))
	for _, row := range rows {
		resourceMonitor, err := row.convert()
		if err != nil {
			return nil, err
		}
		resourceMonitors = append(resourceMonitors, *resourceMonitor)
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
			return &resourceMonitor, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}
