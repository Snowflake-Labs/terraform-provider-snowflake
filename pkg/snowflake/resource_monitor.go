package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ResourceMonitorBuilder extends the generic builder to provide support for triggers
type ResourceMonitorBuilder struct {
	Builder
}

// ResourceMonitor returns a pointer to a ResourceMonitorBuilder that abstracts the DDL operations for a resource monitor.
//
// Supported DDL operations are:
//   - CREATE RESOURCE MONITOR
//   - ALTER RESOURCE MONITOR
//   - DROP RESOURCE MONITOR
//   - SHOW RESOURCE MONITOR
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/user-guide/resource-monitors.html#ddl-for-resource-monitors)
func ResourceMonitor(name string) *ResourceMonitorBuilder {
	return &ResourceMonitorBuilder{
		Builder{
			entityType: ResourceMonitorType,
			name:       name,
		},
	}
}

// @TODO support for a ResourceMonitorAlterBuilder so that we can alter triggers

// ResourceMonitorCreateBuilder extends the generic create builder to provide support for triggers
type ResourceMonitorCreateBuilder struct {
	CreateBuilder

	// triggers consist of the type (DO SUSPEND | SUSPEND_IMMEDIATE | NOTIFY) and
	// the threshold (a percentage value)
	triggers []trigger
}

type trigger struct {
	action    string
	threshold int
}

const (
	// SuspendTrigger suspends all assigned warehouses while allowing currently running queries to complete.
	SuspendTrigger = "SUSPEND"
	// SuspendImmediatelyTrigger suspends all assigned warehouses immediately and cancel any currently running queries or statements using the warehouses.
	SuspendImmediatelyTrigger = "SUSPEND_IMMEDIATE"
	// NotifyTrigger sends an alert (to all users who have enabled notifications for themselves), but do not take any other action.
	NotifyTrigger = "NOTIFY"
)

// Create returns a pointer to a ResourceMonitorCreateBuilder
func (rb *ResourceMonitorBuilder) Create() *ResourceMonitorCreateBuilder {
	return &ResourceMonitorCreateBuilder{
		CreateBuilder{
			name:             rb.name,
			entityType:       rb.entityType,
			stringProperties: make(map[string]string),
			boolProperties:   make(map[string]bool),
			intProperties:    make(map[string]int),
			floatProperties:  make(map[string]float64),
		},
		make([]trigger, 0),
	}
}

// NotifyAt adds a notify trigger at the specified percentage threshold
func (rcb *ResourceMonitorCreateBuilder) NotifyAt(pct int) *ResourceMonitorCreateBuilder {
	rcb.triggers = append(rcb.triggers, trigger{NotifyTrigger, pct})
	return rcb
}

// SuspendAt adds a suspend trigger at the specified percentage threshold
func (rcb *ResourceMonitorCreateBuilder) SuspendAt(pct int) *ResourceMonitorCreateBuilder {
	rcb.triggers = append(rcb.triggers, trigger{SuspendTrigger, pct})
	return rcb
}

// SuspendImmediatelyAt adds a suspend immediately trigger at the specified percentage threshold
func (rcb *ResourceMonitorCreateBuilder) SuspendImmediatelyAt(pct int) *ResourceMonitorCreateBuilder {
	rcb.triggers = append(rcb.triggers, trigger{SuspendImmediatelyTrigger, pct})
	return rcb
}

// Statement returns the SQL statement needed to actually create the resource
func (rcb *ResourceMonitorCreateBuilder) Statement() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`CREATE %v "%v"`, rcb.entityType, rcb.name))

	for k, v := range rcb.stringProperties {
		sb.WriteString(fmt.Sprintf(` %v='%v'`, strings.ToUpper(k), EscapeString(v)))
	}

	for k, v := range rcb.intProperties {
		sb.WriteString(fmt.Sprintf(` %v=%d`, strings.ToUpper(k), v))
	}

	for k, v := range rcb.floatProperties {
		sb.WriteString(fmt.Sprintf(` %v=%.2f`, strings.ToUpper(k), v))
	}

	if len(rcb.triggers) > 0 {
		sb.WriteString(" TRIGGERS")
	}

	for _, trig := range rcb.triggers {
		sb.WriteString(fmt.Sprintf(` ON %d PERCENT DO %v`, trig.threshold, trig.action))
	}

	return sb.String()
}

// SetOnAccount returns the SQL query that will set the resource monitor globally on your Snowflake account
func (rcb *ResourceMonitorCreateBuilder) SetOnAccount() string {
	return fmt.Sprintf(`ALTER ACCOUNT SET RESOURCE_MONITOR = "%v"`, rcb.name)
}

// SetOnWarehouse returns the SQL query that will set the resource monitor on the specified warehouse
func (rcb *ResourceMonitorCreateBuilder) SetOnWarehouse(warehouse string) string {
	return fmt.Sprintf(`ALTER WAREHOUSE "%v" SET RESOURCE_MONITOR = "%v"`, warehouse, rcb.name)
}

type resourceMonitor struct {
	Name                 sql.NullString `db:"name"`
	CreditQuota          sql.NullString `db:"credit_quota"`
	UsedCredits          sql.NullString `db:"used_credits"`
	RemainingCredits     sql.NullString `db:"remaining_credits"`
	Level                sql.NullString `db:"level"`
	Frequency            sql.NullString `db:"frequency"`
	StartTime            sql.NullString `db:"start_time"`
	EndTime              sql.NullString `db:"end_time"`
	NotifyAt             sql.NullString `db:"notify_at"`
	SuspendAt            sql.NullString `db:"suspend_at"`
	SuspendImmediatelyAt sql.NullString `db:"suspend_immediately_at"`
	CreatedOn            sql.NullString `db:"created_on"`
	Owner                sql.NullString `db:"owner"`
	Comment              sql.NullString `db:"comment"`
}

func ScanResourceMonitor(row *sqlx.Row) (*resourceMonitor, error) {
	rm := &resourceMonitor{}
	err := row.StructScan(rm)
	return rm, err
}

func ListResourceMonitors(db *sql.DB) ([]resourceMonitor, error) {
	stmt := "SHOW RESOURCE MONITORS"
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []resourceMonitor{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no resouce monitors found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
