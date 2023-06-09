package sdk

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Compile-time proof of interface implementation.
var _ Alerts = (*alerts)(nil)

type Alerts interface {
	// Create creates a new alert.
	Create(ctx context.Context, id SchemaObjectIdentifier, warehouse AccountObjectIdentifier, schedule AlertSchedule, condition string, action string, opts *CreateAlertOptions) error
	// Alter modifies an existing alert.
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterAlertOptions) error
	// Drop removes an alert.
	Drop(ctx context.Context, id SchemaObjectIdentifier) error
	// Show returns a list of alerts
	Show(ctx context.Context, opts *ShowAlertOptions) ([]*Alert, error)
	// ShowByID returns an alert by ID
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Alert, error)
	// Describe returns the details of an alert.
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*AlertDetails, error)
}

// alerts implements Alerts
type alerts struct {
	client *Client
}

type CreateAlertOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	alert       bool                   `ddl:"static" sql:"ALERT"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`

	// required
	warehouse AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
	schedule  string                  `ddl:"parameter,single_quotes" sql:"SCHEDULE"`

	// optional
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`

	// required
	condition AlertCondition `ddl:"condition,parentheses"   sql:"IF"`
	action    string         `ddl:"parameter,no_equals" sql:"THEN"`
}

type AlertCondition struct {
	Condition string `ddl:"condition,parentheses" sql:"EXISTS"`
}

func (opts *CreateAlertOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return errors.New("invalid object identifier")
	}

	return nil
}

func (v *alerts) Create(ctx context.Context, id SchemaObjectIdentifier, warehouse AccountObjectIdentifier, schedule AlertSchedule, condition string, action string, opts *CreateAlertOptions) error {
	if opts == nil {
		opts = &CreateAlertOptions{}
	}
	opts.name = id
	opts.warehouse = warehouse
	opts.schedule = schedule.String()
	opts.name = id
	opts.condition = AlertCondition{condition}
	opts.action = action
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

type state string

var (
	//Resume makes a suspended alert active.
	Resume state = "RESUME"
	//Suspend puts the alert into a “Suspended” state.
	Suspend state = "SUSPEND"
)

type AlterAlertOptions struct {
	alter    bool                   `ddl:"static" sql:"ALTER"` //lint:ignore U1000 This is used in the ddl tag
	alert    bool                   `ddl:"static" sql:"ALERT"` //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`

	//one of
	State           *state      `ddl:"keyword"`
	Set             *AlertSet   `ddl:"keyword" sql:"SET"`
	Unset           *AlertUnset `ddl:"keyword" sql:"UNSET"`
	ModifyCondition *string     `ddl:"condition,parentheses" sql:"MODIFY CONDITION EXISTS"`
	ModifyAction    *string     `ddl:"parameter,no_equals" sql:"MODIFY ACTION"`
}

func (opts *AlterAlertOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return errors.New("invalid object identifier")
	}

	if everyValueNil(opts.State, opts.Set, opts.Unset, opts.ModifyCondition, opts.ModifyAction) {
		return errors.New("No alter action specified")
	}
	if !exactlyOneValueSet(opts.State, opts.Set, opts.Unset, opts.ModifyCondition, opts.ModifyAction) {
		return errors.New(`
		Only one of the following actions can be performed at a time:
		{
			RESUME | SUSPEND,
			SET,
			UNSET,
			MODIFY CONDITION EXISTS,
			MODIFY ACTION
		}
		`)
	}

	return nil
}

type AlertSet struct {
	Warehouse *AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
	Schedule  *string                  `ddl:"parameter,single_quotes" sql:"SCHEDULE"`
	Comment   *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type AlertUnset struct {
	warehouse *bool `ddl:"keyword" sql:"WAREHOUSE"`
	schedule  *bool `ddl:"keyword" sql:"SCHEDULE"`
	comment   *bool `ddl:"keyword" sql:"COMMENT"`
}

func (v *alerts) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterAlertOptions) error {
	if opts == nil {
		return errors.New("alter alert options cannot be empty")
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

type DropAlertOptions struct {
	drop  bool                   `ddl:"static" sql:"DROP"`  //lint:ignore U1000 This is used in the ddl tag
	alert bool                   `ddl:"static" sql:"ALERT"` //lint:ignore U1000 This is used in the ddl tag
	name  SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DropAlertOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *alerts) Drop(ctx context.Context, id SchemaObjectIdentifier) error {
	// alert drop does not support [IF EXISTS] so there are no drop options.
	opts := &DropAlertOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate alert options: %w", err)
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	if err != nil {
		return err
	}
	return err
}

type ShowAlertOptions struct {
	show   bool  `ddl:"static" sql:"SHOW"` //lint:ignore U1000 This is used in the ddl tag
	Terse  *bool `ddl:"keyword" sql:"TERSE"`
	alerts bool  `ddl:"static" sql:"ALERTS"` //lint:ignore U1000 This is used in the ddl tag

	//optional
	Like       *Like   `ddl:"keyword" sql:"LIKE"`
	In         *In     `ddl:"keyword" sql:"IN"`
	StartsWith *string `ddl:"parameter,no_equals,single_quotes" sql:"STARTS WITH"`
	Limit      *int    `ddl:"parameter,no_equals" sql:"LIMIT"`
}

func (v *Alert) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Alert) ObjectType() ObjectType {
	return ObjectTypeAlert
}

type Alert struct {
	CreatedOn    time.Time
	Name         string
	DatabaseName string
	SchemaName   string
	Owner        string
	Comment      *string
	Warehouse    string
	Schedule     AlertSchedule
	State        string
	Condition    string
	Action       string
}

type AlertSchedule interface {
	alertSchedule()
	String() string
}

func AlertScheduleFromString(s string) (AlertSchedule, bool) {
	parts := strings.Fields(s)
	if len(parts) == 2 && parts[len(parts)-1] == "MINUTE" {
		return AlertScheduleInterval(toInt(parts[0])), true
	} else if withoutCron := dropUsingCron(parts); len(withoutCron) == 6 {
		return AlertScheduleCronExpression{
			Expression: strings.Join(withoutCron[0:len(withoutCron)-1], " "),
			TimeZone:   withoutCron[len(withoutCron)-1],
		}, true
	} else {
		return nil, false
	}

}
func dropUsingCron(parts []string) []string {
	if parts[0] == "USING" && parts[1] == "CRON" {
		return parts[2:]
	} else {
		return parts
	}
}

type AlertScheduleInterval int

func (interval AlertScheduleInterval) String() string {
	return fmt.Sprintf("%d MINUTE", interval)
}

func (AlertScheduleInterval) alertSchedule() {}

type AlertScheduleCronExpression struct {
	Expression string
	TimeZone   string
}

func (expression AlertScheduleCronExpression) String() string {
	return fmt.Sprintf("USING CRON %v %v", expression.Expression, expression.TimeZone)
}

func (AlertScheduleCronExpression) alertSchedule() {}

func (v *Alert) IsEnabled() bool {
	return strings.ToLower(v.State) == "started"
}

func (v *Alert) IsSuspended() bool {
	return strings.ToLower(v.State) == "suspended"
}

type alertDBRow struct {
	CreatedOn    time.Time `db:"created_on"`
	Name         string    `db:"name"`
	DatabaseName string    `db:"database_name"`
	SchemaName   string    `db:"schema_name"`
	Owner        string    `db:"owner"`
	Comment      *string   `db:"comment"`
	Warehouse    string    `db:"warehouse"`
	Schedule     string    `db:"schedule"`
	State        string    `db:"state"` // suspended, started
	Condition    string    `db:"condition"`
	Action       string    `db:"action"`
}

func (row alertDBRow) toAlert() *Alert {
	schedule, _ := AlertScheduleFromString(row.Schedule)
	return &Alert{
		CreatedOn:    row.CreatedOn,
		Name:         row.Name,
		DatabaseName: row.DatabaseName,
		SchemaName:   row.SchemaName,
		Owner:        row.Owner,
		Comment:      row.Comment,
		Warehouse:    row.Warehouse,
		Schedule:     schedule,
		State:        row.State,
		Condition:    row.Condition,
		Action:       row.Action,
	}
}

func (v *alerts) Show(ctx context.Context, opts *ShowAlertOptions) ([]*Alert, error) {
	if opts == nil {
		opts = &ShowAlertOptions{}
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []alertDBRow{}

	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]*Alert, len(dest))
	for i, row := range dest {
		resultList[i] = row.toAlert()
	}

	return resultList, nil
}

func (v *alerts) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Alert, error) {
	alerts, err := v.Show(ctx, &ShowAlertOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
		In: &In{
			Schema: NewSchemaIdentifier(id.DatabaseName(), id.SchemaName()),
		},
	})
	if err != nil {
		return nil, err
	}

	for _, alert := range alerts {
		if alert.ID().name == id.Name() {
			return alert, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type describeAlertOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"` //lint:ignore U1000 This is used in the ddl tag
	alert    bool                   `ddl:"static" sql:"ALERT"`    //lint:ignore U1000 This is used in the ddl tag
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

func (v *describeAlertOptions) validate() error {
	if !validObjectidentifier(v.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

type AlertDetails struct {
	CreatedOn    time.Time
	Name         string
	DatabaseName string
	SchemaName   string
	Owner        string
	Comment      *string
	Warehouse    string
	Schedule     AlertSchedule
	State        string
	Condition    string
	Action       string
}

func (row alertDBRow) toAlertDetails() *AlertDetails {
	schedule, _ := AlertScheduleFromString(row.Schedule)

	return &AlertDetails{
		CreatedOn:    row.CreatedOn,
		Name:         row.Name,
		DatabaseName: row.DatabaseName,
		SchemaName:   row.SchemaName,
		Owner:        row.Owner,
		Comment:      row.Comment,
		Warehouse:    row.Warehouse,
		Schedule:     schedule,
		State:        row.State,
		Condition:    row.Condition,
		Action:       row.Action,
	}
}

func (v *alerts) Describe(ctx context.Context, id SchemaObjectIdentifier) (*AlertDetails, error) {
	opts := &describeAlertOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}

	// SHOW ALERTS and DESCRIBE ALERT SQL statements return the same output
	dest := alertDBRow{}
	err = v.client.queryOne(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}

	return dest.toAlertDetails(), nil
}
