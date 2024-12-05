package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Tasks = (*tasks)(nil)

type tasks struct {
	client *Client
}

func (v *tasks) Create(ctx context.Context, request *CreateTaskRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tasks) CreateOrAlter(ctx context.Context, request *CreateOrAlterTaskRequest) error {
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
	tasks, err := v.Show(ctx, NewShowTaskRequest().WithIn(ExtendedIn{
		In: In{
			Schema: id.SchemaId(),
		},
	}).WithLike(Like{
		Pattern: String(id.Name()),
	}))
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(tasks, func(r Task) bool { return r.Name == id.Name() })
}

func (v *tasks) ShowParameters(ctx context.Context, id SchemaObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Task: id,
		},
	})
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

// TODO(SNOW-1277135): See if depId is necessary or could be removed
func (v *tasks) SuspendRootTasks(ctx context.Context, taskId SchemaObjectIdentifier, id SchemaObjectIdentifier) ([]SchemaObjectIdentifier, error) {
	rootTasks, err := GetRootTasks(v.client.Tasks, ctx, taskId)
	if err != nil {
		return nil, err
	}

	tasksToResume := make([]SchemaObjectIdentifier, 0)
	suspendErrs := make([]error, 0)

	for _, rootTask := range rootTasks {
		// If a root task is started, then it needs to be suspended before the child tasks can be created
		if rootTask.IsStarted() {
			err := v.client.Tasks.Alter(ctx, NewAlterTaskRequest(rootTask.ID()).WithSuspend(true))
			if err != nil {
				log.Printf("[WARN] failed to suspend task %s", rootTask.ID().FullyQualifiedName())
				suspendErrs = append(suspendErrs, err)
			}

			// Resume the task after modifications are complete as long as it is not a standalone task
			// TODO(SNOW-1277135): Document the purpose of this check and why we need different value for GetRootTasks (depId).
			if rootTask.Name != id.Name() {
				tasksToResume = append(tasksToResume, rootTask.ID())
			}
		}
	}

	return tasksToResume, errors.Join(suspendErrs...)
}

func (v *tasks) ResumeTasks(ctx context.Context, ids []SchemaObjectIdentifier) error {
	resumeErrs := make([]error, 0)
	for _, id := range ids {
		err := v.client.Tasks.Alter(ctx, NewAlterTaskRequest(id).WithResume(true))
		if err != nil {
			log.Printf("[WARN] failed to resume task %s", id.FullyQualifiedName())
			resumeErrs = append(resumeErrs, err)
		}
	}
	return errors.Join(resumeErrs...)
}

// GetRootTasks is a way to get all root tasks for the given tasks.
// Snowflake does not have (yet) a method to do it without traversing the task graph manually.
// Task DAG should have a single root but this is checked when the root task is being resumed; that's why we return here multiple roots.
// Cycles should not be possible in a task DAG, but it is checked when the root task is being resumed; that's why this method has to be cycle-proof.
func GetRootTasks(v Tasks, ctx context.Context, id SchemaObjectIdentifier) ([]Task, error) {
	tasksToExamine := collections.NewQueue[SchemaObjectIdentifier]()
	alreadyExaminedTasksNames := make([]string, 0)
	rootTasks := make([]Task, 0)

	tasksToExamine.Push(id)

	for tasksToExamine.Head() != nil {
		current := tasksToExamine.Pop()

		if slices.Contains(alreadyExaminedTasksNames, current.Name()) {
			continue
		}

		task, err := v.ShowByID(ctx, *current)
		if err != nil {
			return nil, err
		}

		if task.TaskRelations.FinalizedRootTask != nil {
			tasksToExamine.Push(*task.TaskRelations.FinalizedRootTask)
			alreadyExaminedTasksNames = append(alreadyExaminedTasksNames, current.Name())
			continue
		}

		predecessors := task.Predecessors
		if len(predecessors) == 0 {
			rootTasks = append(rootTasks, *task)
		} else {
			for _, p := range predecessors {
				tasksToExamine.Push(p)
			}
		}
		alreadyExaminedTasksNames = append(alreadyExaminedTasksNames, current.Name())
	}

	return rootTasks, nil
}

func (r *CreateTaskRequest) toOpts() *CreateTaskOptions {
	opts := &CreateTaskOptions{
		OrReplace:                               r.OrReplace,
		IfNotExists:                             r.IfNotExists,
		name:                                    r.name,
		Schedule:                                r.Schedule,
		AllowOverlappingExecution:               r.AllowOverlappingExecution,
		SessionParameters:                       r.SessionParameters,
		UserTaskTimeoutMs:                       r.UserTaskTimeoutMs,
		SuspendTaskAfterNumFailures:             r.SuspendTaskAfterNumFailures,
		ErrorIntegration:                        r.ErrorIntegration,
		Comment:                                 r.Comment,
		Finalize:                                r.Finalize,
		TaskAutoRetryAttempts:                   r.TaskAutoRetryAttempts,
		Tag:                                     r.Tag,
		UserTaskMinimumTriggerIntervalInSeconds: r.UserTaskMinimumTriggerIntervalInSeconds,
		After:                                   r.After,
		When:                                    r.When,
		sql:                                     r.sql,
	}
	if r.Warehouse != nil {
		opts.Warehouse = &CreateTaskWarehouse{
			Warehouse:                           r.Warehouse.Warehouse,
			UserTaskManagedInitialWarehouseSize: r.Warehouse.UserTaskManagedInitialWarehouseSize,
		}
	}
	if r.Config != nil {
		opts.Config = String(fmt.Sprintf("$$%s$$", *r.Config))
	}
	return opts
}

func (r *CreateOrAlterTaskRequest) toOpts() *CreateOrAlterTaskOptions {
	opts := &CreateOrAlterTaskOptions{
		name:                        r.name,
		Schedule:                    r.Schedule,
		AllowOverlappingExecution:   r.AllowOverlappingExecution,
		UserTaskTimeoutMs:           r.UserTaskTimeoutMs,
		SessionParameters:           r.SessionParameters,
		SuspendTaskAfterNumFailures: r.SuspendTaskAfterNumFailures,
		ErrorIntegration:            r.ErrorIntegration,
		Comment:                     r.Comment,
		Finalize:                    r.Finalize,
		TaskAutoRetryAttempts:       r.TaskAutoRetryAttempts,
		After:                       r.After,
		When:                        r.When,
		sql:                         r.sql,
	}
	if r.Warehouse != nil {
		opts.Warehouse = &CreateTaskWarehouse{
			Warehouse:                           r.Warehouse.Warehouse,
			UserTaskManagedInitialWarehouseSize: r.Warehouse.UserTaskManagedInitialWarehouseSize,
		}
	}
	if r.Config != nil {
		opts.Config = String(fmt.Sprintf("$$%s$$", *r.Config))
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

		SetTags:       r.SetTags,
		UnsetTags:     r.UnsetTags,
		SetFinalize:   r.SetFinalize,
		UnsetFinalize: r.UnsetFinalize,
		ModifyAs:      r.ModifyAs,
		ModifyWhen:    r.ModifyWhen,
		RemoveWhen:    r.RemoveWhen,
	}
	if r.Set != nil {
		opts.Set = &TaskSet{
			Warehouse:                               r.Set.Warehouse,
			UserTaskManagedInitialWarehouseSize:     r.Set.UserTaskManagedInitialWarehouseSize,
			Schedule:                                r.Set.Schedule,
			AllowOverlappingExecution:               r.Set.AllowOverlappingExecution,
			UserTaskTimeoutMs:                       r.Set.UserTaskTimeoutMs,
			SuspendTaskAfterNumFailures:             r.Set.SuspendTaskAfterNumFailures,
			ErrorIntegration:                        r.Set.ErrorIntegration,
			Comment:                                 r.Set.Comment,
			SessionParameters:                       r.Set.SessionParameters,
			TaskAutoRetryAttempts:                   r.Set.TaskAutoRetryAttempts,
			UserTaskMinimumTriggerIntervalInSeconds: r.Set.UserTaskMinimumTriggerIntervalInSeconds,
		}
		if r.Set.Config != nil {
			opts.Set.Config = String(fmt.Sprintf("$$%s$$", *r.Set.Config))
		}
	}
	if r.Unset != nil {
		opts.Unset = &TaskUnset{
			Warehouse:                               r.Unset.Warehouse,
			UserTaskManagedInitialWarehouseSize:     r.Unset.UserTaskManagedInitialWarehouseSize,
			Schedule:                                r.Unset.Schedule,
			Config:                                  r.Unset.Config,
			AllowOverlappingExecution:               r.Unset.AllowOverlappingExecution,
			UserTaskTimeoutMs:                       r.Unset.UserTaskTimeoutMs,
			SuspendTaskAfterNumFailures:             r.Unset.SuspendTaskAfterNumFailures,
			ErrorIntegration:                        r.Unset.ErrorIntegration,
			Comment:                                 r.Unset.Comment,
			TaskAutoRetryAttempts:                   r.Unset.TaskAutoRetryAttempts,
			UserTaskMinimumTriggerIntervalInSeconds: r.Unset.UserTaskMinimumTriggerIntervalInSeconds,
			SessionParametersUnset:                  r.Unset.SessionParametersUnset,
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
		CreatedOn:                 r.CreatedOn,
		Id:                        r.Id,
		Name:                      r.Name,
		DatabaseName:              r.DatabaseName,
		SchemaName:                r.SchemaName,
		Owner:                     r.Owner,
		Definition:                r.Definition,
		AllowOverlappingExecution: r.AllowOverlappingExecution == "true",
		OwnerRoleType:             r.OwnerRoleType,
	}
	taskRelations, err := ToTaskRelations(r.TaskRelations)
	if err != nil {
		log.Printf("[DEBUG] failed to convert task relations: %v", err)
	} else {
		task.TaskRelations = taskRelations
	}
	if r.Comment.Valid {
		task.Comment = r.Comment.String
	}
	if r.Warehouse.Valid && r.Warehouse.String != "null" {
		id, err := ParseAccountObjectIdentifier(r.Warehouse.String)
		if err != nil {
			log.Printf("[DEBUG] failed to parse warehouse: %v", err)
		} else {
			task.Warehouse = &id
		}
	}
	if r.Schedule.Valid {
		task.Schedule = r.Schedule.String
	}
	if len(r.Predecessors) > 0 {
		names, err := getPredecessors(r.Predecessors)
		ids := make([]SchemaObjectIdentifier, len(names))
		if err == nil {
			for i, name := range names {
				ids[i] = NewSchemaObjectIdentifier(r.DatabaseName, r.SchemaName, name)
			}
		}
		task.Predecessors = ids
	} else {
		task.Predecessors = make([]SchemaObjectIdentifier, 0)
	}
	if len(r.State) > 0 {
		taskState, err := ToTaskState(r.State)
		if err != nil {
			log.Printf("[DEBUG] failed to convert to task state: %v", err)
		} else {
			task.State = taskState
		}
	}
	if r.Condition.Valid {
		task.Condition = r.Condition.String
	}
	if r.ErrorIntegration.Valid && r.ErrorIntegration.String != "null" {
		id, err := ParseAccountObjectIdentifier(r.ErrorIntegration.String)
		if err != nil {
			log.Printf("[DEBUG] failed to parse error_integration: %v", err)
		} else {
			task.ErrorIntegration = &id
		}
	}
	if r.LastCommittedOn.Valid {
		task.LastCommittedOn = r.LastCommittedOn.String
	}
	if r.LastSuspendedOn.Valid {
		task.LastSuspendedOn = r.LastSuspendedOn.String
	}
	if r.Config.Valid {
		task.Config = r.Config.String
	}
	if r.Budget.Valid {
		task.Budget = r.Budget.String
	}
	if r.LastSuspendedReason.Valid {
		task.LastSuspendedReason = r.LastSuspendedReason.String
	}
	return &task
}

func getPredecessors(predecessors string) ([]string, error) {
	// Since 2022_03, Snowflake returns this as a JSON array (even empty)
	// The list is formatted, e.g.:
	// e.g. `[\n  \"\\\"qgb)Z1KcNWJ(\\\".\\\"glN@JtR=7dzP$7\\\".\\\"_XEL(7N_F?@frgT5>dQS>V|vSy,J\\\"\"\n]`.
	predecessorNames := make([]string, 0)
	err := json.Unmarshal([]byte(predecessors), &predecessorNames)
	if err == nil {
		for i, predecessorName := range predecessorNames {
			formattedName := strings.Trim(predecessorName, "\\\"")
			idx := strings.LastIndex(formattedName, "\"") + 1
			// -1 because of not found +1 is 0
			if idx == 0 {
				idx = strings.LastIndex(formattedName, ".") + 1
			} else if strings.LastIndex(formattedName, ".\"")+2 < idx {
				idx++
			}
			formattedName = formattedName[idx:]
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
