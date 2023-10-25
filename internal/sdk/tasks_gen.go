// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import (
	"context"
	"database/sql"
)

type Tasks interface {
	Create(ctx context.Context, request *CreateTaskRequest) error
	Clone(ctx context.Context, request *CloneTaskRequest) error
	Alter(ctx context.Context, request *AlterTaskRequest) error
	Drop(ctx context.Context, request *DropTaskRequest) error
	Show(ctx context.Context, request *ShowTaskRequest) ([]Task, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Task, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*Task, error)
	Execute(ctx context.Context, request *ExecuteTaskRequest) error
}

// CreateTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-task.
type CreateTaskOptions struct {
	create                      bool                     `ddl:"static" sql:"CREATE"`
	OrReplace                   *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	task                        bool                     `ddl:"static" sql:"TASK"`
	IfNotExists                 *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                        SchemaObjectIdentifier   `ddl:"identifier"`
	Warehouse                   *CreateTaskWarehouse     `ddl:"keyword"`
	Schedule                    *string                  `ddl:"parameter,single_quotes" sql:"SCHEDULE"`
	Config                      *string                  `ddl:"parameter,no_quotes" sql:"CONFIG"`
	AllowOverlappingExecution   *bool                    `ddl:"parameter" sql:"ALLOW_OVERLAPPING_EXECUTION"`
	SessionParameters           *SessionParameters       `ddl:"list,no_parentheses"`
	UserTaskTimeoutMs           *int                     `ddl:"parameter" sql:"USER_TASK_TIMEOUT_MS"`
	SuspendTaskAfterNumFailures *int                     `ddl:"parameter" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	ErrorIntegration            *string                  `ddl:"parameter,no_quotes" sql:"ERROR_INTEGRATION"`
	CopyGrants                  *bool                    `ddl:"keyword" sql:"COPY GRANTS"`
	Comment                     *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
	After                       []SchemaObjectIdentifier `ddl:"parameter,no_equals" sql:"AFTER"`
	Tag                         []TagAssociation         `ddl:"keyword,parentheses" sql:"TAG"`
	When                        *string                  `ddl:"parameter,no_quotes,no_equals" sql:"WHEN"`
	as                          bool                     `ddl:"static" sql:"AS"`
	sql                         string                   `ddl:"keyword,no_quotes"`
}

type CreateTaskWarehouse struct {
	Warehouse                           *AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
	UserTaskManagedInitialWarehouseSize *WarehouseSize           `ddl:"parameter,single_quotes" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
}

// CloneTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-task#variant-syntax.
type CloneTaskOptions struct {
	create     bool                   `ddl:"static" sql:"CREATE"`
	OrReplace  *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	task       bool                   `ddl:"static" sql:"TASK"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
	clone      bool                   `ddl:"static" sql:"CLONE"`
	sourceTask SchemaObjectIdentifier `ddl:"identifier"`
	CopyGrants *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
}

// AlterTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-task.
type AlterTaskOptions struct {
	alter       bool                     `ddl:"static" sql:"ALTER"`
	task        bool                     `ddl:"static" sql:"TASK"`
	IfExists    *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name        SchemaObjectIdentifier   `ddl:"identifier"`
	Resume      *bool                    `ddl:"keyword" sql:"RESUME"`
	Suspend     *bool                    `ddl:"keyword" sql:"SUSPEND"`
	RemoveAfter []SchemaObjectIdentifier `ddl:"parameter,no_equals" sql:"REMOVE AFTER"`
	AddAfter    []SchemaObjectIdentifier `ddl:"parameter,no_equals" sql:"ADD AFTER"`
	Set         *TaskSet                 `ddl:"keyword" sql:"SET"`
	Unset       *TaskUnset               `ddl:"keyword" sql:"UNSET"`
	SetTags     []TagAssociation         `ddl:"keyword" sql:"SET TAG"`
	UnsetTags   []ObjectIdentifier       `ddl:"keyword" sql:"UNSET TAG"`
	ModifyAs    *string                  `ddl:"parameter,no_quotes,no_equals" sql:"MODIFY AS"`
	ModifyWhen  *string                  `ddl:"parameter,no_quotes,no_equals" sql:"MODIFY WHEN"`
}

type TaskSet struct {
	Warehouse                           *AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
	UserTaskManagedInitialWarehouseSize *WarehouseSize           `ddl:"parameter,single_quotes" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
	Schedule                            *string                  `ddl:"parameter,single_quotes" sql:"SCHEDULE"`
	Config                              *string                  `ddl:"parameter,no_quotes" sql:"CONFIG"`
	AllowOverlappingExecution           *bool                    `ddl:"parameter" sql:"ALLOW_OVERLAPPING_EXECUTION"`
	UserTaskTimeoutMs                   *int                     `ddl:"parameter" sql:"USER_TASK_TIMEOUT_MS"`
	SuspendTaskAfterNumFailures         *int                     `ddl:"parameter" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	ErrorIntegration                    *string                  `ddl:"parameter,no_quotes" sql:"ERROR_INTEGRATION"`
	Comment                             *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
	SessionParameters                   *SessionParameters       `ddl:"list,no_parentheses"`
}

type TaskUnset struct {
	Warehouse                   *bool                   `ddl:"keyword" sql:"WAREHOUSE"`
	Schedule                    *bool                   `ddl:"keyword" sql:"SCHEDULE"`
	Config                      *bool                   `ddl:"keyword" sql:"CONFIG"`
	AllowOverlappingExecution   *bool                   `ddl:"keyword" sql:"ALLOW_OVERLAPPING_EXECUTION"`
	UserTaskTimeoutMs           *bool                   `ddl:"keyword" sql:"USER_TASK_TIMEOUT_MS"`
	SuspendTaskAfterNumFailures *bool                   `ddl:"keyword" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	ErrorIntegration            *bool                   `ddl:"keyword" sql:"ERROR_INTEGRATION"`
	Comment                     *bool                   `ddl:"keyword" sql:"COMMENT"`
	SessionParametersUnset      *SessionParametersUnset `ddl:"list,no_parentheses"`
}

// DropTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-task.
type DropTaskOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	task     bool                   `ddl:"static" sql:"TASK"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-tasks.
type ShowTaskOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	Terse      *bool      `ddl:"keyword" sql:"TERSE"`
	tasks      bool       `ddl:"static" sql:"TASKS"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	In         *In        `ddl:"keyword" sql:"IN"`
	StartsWith *string    `ddl:"parameter,no_equals,single_quotes" sql:"STARTS WITH"`
	RootOnly   *bool      `ddl:"keyword" sql:"ROOT ONLY"`
	Limit      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type taskDBRow struct {
	CreatedOn                 string         `db:"created_on"`
	Name                      string         `db:"name"`
	Id                        sql.NullString `db:"id"`
	DatabaseName              string         `db:"database_name"`
	SchemaName                string         `db:"schema_name"`
	Owner                     sql.NullString `db:"owner"`
	Comment                   sql.NullString `db:"comment"`
	Warehouse                 sql.NullString `db:"warehouse"`
	Schedule                  sql.NullString `db:"schedule"`
	Predecessors              sql.NullString `db:"predecessors"`
	State                     sql.NullString `db:"state"`
	Definition                sql.NullString `db:"definition"`
	Condition                 sql.NullString `db:"condition"`
	AllowOverlappingExecution sql.NullString `db:"allow_overlapping_execution"`
	ErrorIntegration          sql.NullString `db:"error_integration"`
	LastCommittedOn           sql.NullString `db:"last_committed_on"`
	LastSuspendedOn           sql.NullString `db:"last_suspended_on"`
	OwnerRoleType             sql.NullString `db:"owner_role_type"`
	Config                    sql.NullString `db:"config"`
	Budget                    sql.NullString `db:"budget"`
}

type Task struct {
	CreatedOn                 string
	Name                      string
	Id                        string
	DatabaseName              string
	SchemaName                string
	Owner                     string
	Comment                   string
	Warehouse                 string
	Schedule                  string
	Predecessors              []SchemaObjectIdentifier
	State                     TaskState
	Definition                string
	Condition                 string
	AllowOverlappingExecution bool
	ErrorIntegration          string
	LastCommittedOn           string
	LastSuspendedOn           string
	OwnerRoleType             string
	Config                    string
	Budget                    string
}

// DescribeTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-task.
type DescribeTaskOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	task     bool                   `ddl:"static" sql:"TASK"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// ExecuteTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/execute-task.
type ExecuteTaskOptions struct {
	execute   bool                   `ddl:"static" sql:"EXECUTE"`
	task      bool                   `ddl:"static" sql:"TASK"`
	name      SchemaObjectIdentifier `ddl:"identifier"`
	RetryLast *bool                  `ddl:"keyword" sql:"RETRY LAST"`
}

func (v *Task) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

type TaskState string

const (
	TaskStateStarted   TaskState = "started"
	TaskStateSuspended TaskState = "suspended"
)

func (v *Task) IsStarted() bool {
	return v.State == TaskStateStarted
}
