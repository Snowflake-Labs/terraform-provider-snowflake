// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// ResourceMonitorBuilder extends the generic builder to provide support for triggers.
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
func NewResourceMonitorBuilder(name string) *ResourceMonitorBuilder {
	return &ResourceMonitorBuilder{
		Builder{
			entityType: ResourceMonitorType,
			name:       name,
		},
	}
}

// @TODO support for a ResourceMonitorAlterBuilder so that we can alter triggers

// ResourceMonitorCreateBuilder extends the generic create builder to provide support for triggers.
type ResourceMonitorCreateBuilder struct {
	CreateBuilder

	// triggers consist of the type (DO SUSPEND | SUSPEND_IMMEDIATE | NOTIFY) and
	// the threshold (a percentage value)
	triggers []trigger
}

type ResourceMonitorAlterBuilder struct {
	AlterPropertiesBuilder

	// triggers consist of the type (DO SUSPEND | SUSPEND_IMMEDIATE | NOTIFY) and
	// the threshold (a percentage value)
	triggers []trigger
}

type action string

type trigger struct {
	action    action
	threshold int
}

const (
	// SuspendTrigger suspends all assigned warehouses while allowing currently running queries to complete.
	SuspendTrigger action = "SUSPEND"
	// SuspendImmediatelyTrigger suspends all assigned warehouses immediately and cancel any currently running queries or statements using the warehouses.
	SuspendImmediatelyTrigger action = "SUSPEND_IMMEDIATE"
	// NotifyTrigger sends an alert (to all users who have enabled notifications for themselves), but do not take any other action.
	NotifyTrigger action = "NOTIFY"
)

// Create returns a pointer to a ResourceMonitorCreateBuilder.
func (rb *ResourceMonitorBuilder) Create() *ResourceMonitorCreateBuilder {
	return &ResourceMonitorCreateBuilder{
		CreateBuilder{
			name:                 rb.name,
			entityType:           rb.entityType,
			stringProperties:     make(map[string]string),
			boolProperties:       make(map[string]bool),
			intProperties:        make(map[string]int),
			floatProperties:      make(map[string]float64),
			stringListProperties: make(map[string][]string),
		},
		make([]trigger, 0),
	}
}

// NotifyAt adds a notify trigger at the specified percentage threshold.
func (rcb *ResourceMonitorCreateBuilder) NotifyAt(pct int) *ResourceMonitorCreateBuilder {
	rcb.triggers = append(rcb.triggers, trigger{NotifyTrigger, pct})
	return rcb
}

// SuspendAt adds a suspend trigger at the specified percentage threshold.
func (rcb *ResourceMonitorCreateBuilder) SuspendAt(pct int) *ResourceMonitorCreateBuilder {
	rcb.triggers = append(rcb.triggers, trigger{SuspendTrigger, pct})
	return rcb
}

// SuspendImmediatelyAt adds a suspend immediately trigger at the specified percentage threshold.
func (rcb *ResourceMonitorCreateBuilder) SuspendImmediatelyAt(pct int) *ResourceMonitorCreateBuilder {
	rcb.triggers = append(rcb.triggers, trigger{SuspendImmediatelyTrigger, pct})
	return rcb
}

// Statement returns the SQL statement needed to actually create the resource.
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

	for k, v := range rcb.stringListProperties {
		sb.WriteString(fmt.Sprintf(" %s=%s", strings.ToUpper(k), formatStringList(v)))
	}

	if len(rcb.triggers) > 0 {
		sb.WriteString(" TRIGGERS")
	}

	for _, trig := range rcb.triggers {
		sb.WriteString(fmt.Sprintf(` ON %d PERCENT DO %v`, trig.threshold, trig.action))
	}

	return sb.String()
}

// SetOnAccount returns the SQL query that will set the resource monitor globally on your Snowflake account.
func (rcb *ResourceMonitorCreateBuilder) SetOnAccount() string {
	return fmt.Sprintf(`ALTER ACCOUNT SET RESOURCE_MONITOR = "%v"`, rcb.name)
}

// SetOnWarehouse returns the SQL query that will set the resource monitor on the specified warehouse.
func (rcb *ResourceMonitorCreateBuilder) SetOnWarehouse(warehouse string) string {
	return fmt.Sprintf(`ALTER WAREHOUSE "%v" SET RESOURCE_MONITOR = "%v"`, warehouse, rcb.name)
}

// Alter returns a pointer to a ResourceMonitorAlterBuilder.
func (rb *ResourceMonitorBuilder) Alter() *ResourceMonitorAlterBuilder {
	return &ResourceMonitorAlterBuilder{
		AlterPropertiesBuilder{
			name:                 rb.name,
			entityType:           rb.entityType,
			stringProperties:     make(map[string]string),
			boolProperties:       make(map[string]bool),
			intProperties:        make(map[string]int),
			floatProperties:      make(map[string]float64),
			stringListProperties: make(map[string][]string),
		},
		make([]trigger, 0),
	}
}

// NotifyAt adds a notify trigger at the specified percentage threshold.
func (rcb *ResourceMonitorAlterBuilder) NotifyAt(pct int) *ResourceMonitorAlterBuilder {
	rcb.triggers = append(rcb.triggers, trigger{NotifyTrigger, pct})
	return rcb
}

// SuspendAt adds a suspend trigger at the specified percentage threshold.
func (rcb *ResourceMonitorAlterBuilder) SuspendAt(pct int) *ResourceMonitorAlterBuilder {
	rcb.triggers = append(rcb.triggers, trigger{SuspendTrigger, pct})
	return rcb
}

// SuspendImmediatelyAt adds a suspend immediately trigger at the specified percentage threshold.
func (rcb *ResourceMonitorAlterBuilder) SuspendImmediatelyAt(pct int) *ResourceMonitorAlterBuilder {
	rcb.triggers = append(rcb.triggers, trigger{SuspendImmediatelyTrigger, pct})
	return rcb
}

// Statement returns the SQL statement needed to actually alter the resource.
func (rcb *ResourceMonitorAlterBuilder) Statement() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`ALTER %v "%v" SET`, rcb.entityType, rcb.name))

	for k, v := range rcb.stringProperties {
		sb.WriteString(fmt.Sprintf(` %v='%v'`, strings.ToUpper(k), EscapeString(v)))
	}

	for k, v := range rcb.intProperties {
		sb.WriteString(fmt.Sprintf(` %v=%d`, strings.ToUpper(k), v))
	}

	for k, v := range rcb.floatProperties {
		sb.WriteString(fmt.Sprintf(` %v=%.2f`, strings.ToUpper(k), v))
	}

	for k, v := range rcb.stringListProperties {
		sb.WriteString(fmt.Sprintf(" %s=%s", strings.ToUpper(k), formatStringList(v)))
	}

	// If the only change is the trigges, then we do not need to add the SET keyword
	if strings.HasSuffix(sb.String(), "SET") {
		sb.Reset()
		sb.WriteString(fmt.Sprintf(`ALTER %v "%v"`, rcb.entityType, rcb.name))
	}
	if len(rcb.triggers) > 0 {
		sb.WriteString(" TRIGGERS")
	}

	for _, trig := range rcb.triggers {
		sb.WriteString(fmt.Sprintf(` ON %d PERCENT DO %v`, trig.threshold, trig.action))
	}

	return sb.String()
}

// SetOnAccount returns the SQL query that will set the resource monitor globally on your Snowflake account.
func (rcb *ResourceMonitorAlterBuilder) SetOnAccount() string {
	return fmt.Sprintf(`ALTER ACCOUNT SET RESOURCE_MONITOR = "%v"`, rcb.name)
}

func (rcb *ResourceMonitorAlterBuilder) UnsetOnAccount() string {
	return `ALTER ACCOUNT SET RESOURCE_MONITOR = NULL`
}

// SetOnWarehouse returns the SQL query that will set the resource monitor on the specified warehouse.
func (rcb *ResourceMonitorAlterBuilder) SetOnWarehouse(warehouse string) string {
	return fmt.Sprintf(`ALTER WAREHOUSE "%v" SET RESOURCE_MONITOR = "%v"`, warehouse, rcb.name)
}

// UnsetOnWarehouse returns the SQL query that will unset the resource monitor on the specified warehouse.
func (rcb *ResourceMonitorAlterBuilder) UnsetOnWarehouse(warehouse string) string {
	return fmt.Sprintf(`ALTER WAREHOUSE "%v" SET RESOURCE_MONITOR = NULL`, warehouse)
}

type ResourceMonitor struct {
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
	NotifyUsers          sql.NullString `db:"notify_users"`
}

func ScanResourceMonitor(row *sqlx.Row) (*ResourceMonitor, error) {
	rm := &ResourceMonitor{}
	err := row.StructScan(rm)
	return rm, err
}

func ListResourceMonitors(db *sql.DB) ([]ResourceMonitor, error) {
	stmt := "SHOW RESOURCE MONITORS"
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []ResourceMonitor{}
	if err := sqlx.StructScan(rows, &dbs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no resource monitors found")
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return dbs, nil
}
