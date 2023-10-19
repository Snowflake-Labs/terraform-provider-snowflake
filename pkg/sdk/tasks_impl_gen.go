package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

var _ Tasks = (*tasks)(nil)

type tasks struct {
	client *Client
}

func (v *tasks) Create(ctx context.Context, request *CreateTaskRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tasks) Clone(ctx context.Context, request *CloneTaskRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tasks) Alter(ctx context.Context, request *AlterTaskRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tasks) Drop(ctx context.Context, request *DropTaskRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tasks) Show(ctx context.Context, request *ShowTaskRequest) ([]Task, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[taskDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[taskDBRow, Task](dbRows)
	return resultList, nil
}

func (v *tasks) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Task, error) {
	return v.Describe(ctx, id)
}

func (v *tasks) Describe(ctx context.Context, id SchemaObjectIdentifier) (*Task, error) {
	opts := &DescribeTaskOptions{
		name: id,
	}
	result, err := validateAndQueryOne[taskDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (v *tasks) Execute(ctx context.Context, request *ExecuteTaskRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

// TODO [SNOW-884987]: handle cycles
func (v *tasks) GetRootTasks(ctx context.Context, id SchemaObjectIdentifier) ([]Task, error) {
	task, err := v.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	predecessors := task.Predecessors
	// no predecessors mean this is a root task
	if len(predecessors) == 0 {
		return []Task{*task}, nil
	}

	rootTasks := make([]Task, 0, len(predecessors))
	for _, predecessor := range predecessors {
		predecessorTasks, err := v.GetRootTasks(ctx, predecessor)
		if err != nil {
			return nil, fmt.Errorf("unable to get predecessors for task %s err = %w", predecessor.FullyQualifiedName(), err)
		}
		rootTasks = append(rootTasks, predecessorTasks...)
	}

	// TODO [SNOW-884987]: extract unique function in our collection helper (if cycle-proof algorithm still needs it)
	keys := make(map[string]bool)
	uniqueRootTasks := make([]Task, 0, len(rootTasks))
	for _, rootTask := range rootTasks {
		if _, value := keys[rootTask.ID().FullyQualifiedName()]; !value {
			keys[rootTask.ID().FullyQualifiedName()] = true
			uniqueRootTasks = append(uniqueRootTasks, rootTask)
		}
	}

	return uniqueRootTasks, nil
}

func (r *CreateTaskRequest) toOpts() *CreateTaskOptions {
	opts := &CreateTaskOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		Schedule:                    r.Schedule,
		Config:                      r.Config,
		AllowOverlappingExecution:   r.AllowOverlappingExecution,
		SessionParameters:           r.SessionParameters,
		UserTaskTimeoutMs:           r.UserTaskTimeoutMs,
		SuspendTaskAfterNumFailures: r.SuspendTaskAfterNumFailures,
		ErrorIntegration:            r.ErrorIntegration,
		CopyGrants:                  r.CopyGrants,
		Comment:                     r.Comment,
		After:                       r.After,
		Tag:                         r.Tag,
		When:                        r.When,
		sql:                         r.sql,
	}
	if r.Warehouse != nil {
		opts.Warehouse = &CreateTaskWarehouse{
			Warehouse:                           r.Warehouse.Warehouse,
			UserTaskManagedInitialWarehouseSize: r.Warehouse.UserTaskManagedInitialWarehouseSize,
		}
	}
	return opts
}

func (r *CloneTaskRequest) toOpts() *CloneTaskOptions {
	opts := &CloneTaskOptions{
		OrReplace:  r.OrReplace,
		name:       r.name,
		sourceTask: r.sourceTask,
		CopyGrants: r.CopyGrants,
	}
	return opts
}

func (r *AlterTaskRequest) toOpts() *AlterTaskOptions {
	opts := &AlterTaskOptions{
		IfExists:    r.IfExists,
		name:        r.name,
		Resume:      r.Resume,
		Suspend:     r.Suspend,
		RemoveAfter: r.RemoveAfter,
		AddAfter:    r.AddAfter,

		SetTags:    r.SetTags,
		UnsetTags:  r.UnsetTags,
		ModifyAs:   r.ModifyAs,
		ModifyWhen: r.ModifyWhen,
	}
	if r.Set != nil {
		opts.Set = &TaskSet{
			Warehouse:                           r.Set.Warehouse,
			UserTaskManagedInitialWarehouseSize: r.Set.UserTaskManagedInitialWarehouseSize,
			Schedule:                            r.Set.Schedule,
			Config:                              r.Set.Config,
			AllowOverlappingExecution:           r.Set.AllowOverlappingExecution,
			UserTaskTimeoutMs:                   r.Set.UserTaskTimeoutMs,
			SuspendTaskAfterNumFailures:         r.Set.SuspendTaskAfterNumFailures,
			ErrorIntegration:                    r.Set.ErrorIntegration,
			Comment:                             r.Set.Comment,
			SessionParameters:                   r.Set.SessionParameters,
		}
	}
	if r.Unset != nil {
		opts.Unset = &TaskUnset{
			Warehouse:                   r.Unset.Warehouse,
			Schedule:                    r.Unset.Schedule,
			Config:                      r.Unset.Config,
			AllowOverlappingExecution:   r.Unset.AllowOverlappingExecution,
			UserTaskTimeoutMs:           r.Unset.UserTaskTimeoutMs,
			SuspendTaskAfterNumFailures: r.Unset.SuspendTaskAfterNumFailures,
			ErrorIntegration:            r.Unset.ErrorIntegration,
			Comment:                     r.Unset.Comment,
			SessionParametersUnset:      r.Unset.SessionParametersUnset,
		}
	}
	return opts
}

func (r *DropTaskRequest) toOpts() *DropTaskOptions {
	opts := &DropTaskOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowTaskRequest) toOpts() *ShowTaskOptions {
	opts := &ShowTaskOptions{
		Terse:      r.Terse,
		Like:       r.Like,
		In:         r.In,
		StartsWith: r.StartsWith,
		RootOnly:   r.RootOnly,
		Limit:      r.Limit,
	}
	return opts
}

func (r taskDBRow) convert() *Task {
	task := Task{
		CreatedOn:    r.CreatedOn,
		Name:         r.Name,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
	}
	if r.Id.Valid {
		task.Id = r.Id.String
	}
	if r.Owner.Valid {
		task.Owner = r.Owner.String
	}
	if r.Comment.Valid {
		task.Comment = r.Comment.String
	}
	if r.Warehouse.Valid {
		task.Warehouse = r.Warehouse.String
	}
	if r.Schedule.Valid {
		task.Schedule = r.Schedule.String
	}
	if r.Predecessors.Valid {
		names, err := getPredecessors(r.Predecessors.String)
		ids := make([]SchemaObjectIdentifier, len(names))
		if err == nil {
			for i, name := range names {
				ids[i] = NewSchemaObjectIdentifier(r.DatabaseName, r.SchemaName, name)
			}
		}
		task.Predecessors = ids
	}
	if r.State.Valid {
		if strings.ToLower(r.State.String) == string(TaskStateStarted) {
			task.State = TaskStateStarted
		} else {
			task.State = TaskStateSuspended
		}
	}
	if r.Definition.Valid {
		task.Definition = r.Definition.String
	}
	if r.Condition.Valid {
		task.Condition = r.Condition.String
	}
	if r.AllowOverlappingExecution.Valid {
		task.AllowOverlappingExecution = r.AllowOverlappingExecution.String == "true"
	}
	if r.ErrorIntegration.Valid && r.ErrorIntegration.String != "null" {
		task.ErrorIntegration = r.ErrorIntegration.String
	}
	if r.LastCommittedOn.Valid {
		task.LastCommittedOn = r.LastCommittedOn.String
	}
	if r.LastSuspendedOn.Valid {
		task.LastSuspendedOn = r.LastSuspendedOn.String
	}
	if r.OwnerRoleType.Valid {
		task.OwnerRoleType = r.OwnerRoleType.String
	}
	if r.Config.Valid {
		task.Config = r.Config.String
	}
	if r.Budget.Valid {
		task.Budget = r.Budget.String
	}
	return &task
}

// TODO [SNOW-884987]: test this method with different task names (like "Ls.T7-(bt{.lWd@DRWkyA6<6hNdh")
func getPredecessors(predecessors string) ([]string, error) {
	// Since 2022_03, Snowflake returns this as a JSON array (even empty)
	// The list is formatted, e.g.:
	// e.g. `[\n  \"\\\"qgb)Z1KcNWJ(\\\".\\\"glN@JtR=7dzP$7\\\".\\\"_XEL(7N_F?@frgT5>dQS>V|vSy,J\\\"\"\n]`.
	predecessorNames := make([]string, 0)
	err := json.Unmarshal([]byte(predecessors), &predecessorNames)
	if err == nil {
		for i, predecessorName := range predecessorNames {
			formattedName := predecessorName[strings.LastIndex(predecessorName, ".\"")+1:]
			formattedName = strings.Trim(formattedName, "\\\"")
			predecessorNames[i] = formattedName
		}
	}
	return predecessorNames, err
}

func (r *DescribeTaskRequest) toOpts() *DescribeTaskOptions {
	opts := &DescribeTaskOptions{
		name: r.name,
	}
	return opts
}

func (r *ExecuteTaskRequest) toOpts() *ExecuteTaskOptions {
	opts := &ExecuteTaskOptions{
		name:      r.name,
		RetryLast: r.RetryLast,
	}
	return opts
}
