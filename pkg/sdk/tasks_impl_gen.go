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
