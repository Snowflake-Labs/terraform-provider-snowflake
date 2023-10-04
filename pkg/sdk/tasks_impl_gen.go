package sdk

import "context"

var _ Tasks = (*tasks)(nil)

type tasks struct {
	client *Client
}

func (v *tasks) Create(ctx context.Context, request *CreateTaskRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tasks) Alter(ctx context.Context, request *AlterTaskRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
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
			Warehouse:                   r.Set.Warehouse,
			Schedule:                    r.Set.Schedule,
			Config:                      r.Set.Config,
			AllowOverlappingExecution:   r.Set.AllowOverlappingExecution,
			UserTaskTimeoutMs:           r.Set.UserTaskTimeoutMs,
			SuspendTaskAfterNumFailures: r.Set.SuspendTaskAfterNumFailures,
			Comment:                     r.Set.Comment,
			SessionParameters:           r.Set.SessionParameters,
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
			Comment:                     r.Unset.Comment,
			SessionParametersUnset:      r.Unset.SessionParametersUnset,
		}
	}
	return opts
}
