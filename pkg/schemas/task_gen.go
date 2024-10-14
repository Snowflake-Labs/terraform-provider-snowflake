package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowTaskSchema represents output of SHOW query for the single Task.
var ShowTaskSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"database_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schema_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"warehouse": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schedule": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"predecessors": {
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
	"state": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"definition": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"condition": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"allow_overlapping_execution": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"error_integration": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"last_committed_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"last_suspended_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"config": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"budget": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"task_relations": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"predecessors": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"finalizer": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"finalized_root_task": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	},
	"last_suspended_reason": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowTaskSchema

func TaskToSchema(task *sdk.Task) map[string]any {
	taskSchema := make(map[string]any)
	taskSchema["created_on"] = task.CreatedOn
	taskSchema["name"] = task.Name
	taskSchema["id"] = task.Id
	taskSchema["database_name"] = task.DatabaseName
	taskSchema["schema_name"] = task.SchemaName
	taskSchema["owner"] = task.Owner
	taskSchema["comment"] = task.Comment
	if task.Warehouse != nil {
		taskSchema["warehouse"] = task.Warehouse.Name()
	}
	taskSchema["schedule"] = task.Schedule
	taskSchema["predecessors"] = collections.Map(task.Predecessors, sdk.SchemaObjectIdentifier.FullyQualifiedName)
	taskSchema["state"] = string(task.State)
	taskSchema["definition"] = task.Definition
	taskSchema["condition"] = task.Condition
	taskSchema["allow_overlapping_execution"] = task.AllowOverlappingExecution
	if task.ErrorIntegration != nil {
		taskSchema["error_integration"] = task.ErrorIntegration.Name()
	}
	taskSchema["last_committed_on"] = task.LastCommittedOn
	taskSchema["last_suspended_on"] = task.LastSuspendedOn
	taskSchema["owner_role_type"] = task.OwnerRoleType
	taskSchema["config"] = task.Config
	taskSchema["budget"] = task.Budget
	taskSchema["last_suspended_reason"] = task.LastSuspendedReason
	// This is manually edited, please don't re-generate this file
	finalizer := ""
	if task.TaskRelations.FinalizerTask != nil {
		finalizer = task.TaskRelations.FinalizerTask.FullyQualifiedName()
	}
	finalizedRootTask := ""
	if task.TaskRelations.FinalizedRootTask != nil {
		finalizedRootTask = task.TaskRelations.FinalizedRootTask.FullyQualifiedName()
	}
	taskSchema["task_relations"] = []any{
		map[string]any{
			"predecessors":        collections.Map(task.TaskRelations.Predecessors, sdk.SchemaObjectIdentifier.FullyQualifiedName),
			"finalizer":           finalizer,
			"finalized_root_task": finalizedRootTask,
		},
	}
	return taskSchema
}

var _ = TaskToSchema
