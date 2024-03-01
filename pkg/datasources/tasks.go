package datasources

import (
	"context"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var tasksSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the tasks from.",
	},
	"tasks": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The tasks in the schema",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"database": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"schema": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"comment": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"warehouse": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func Tasks() *schema.Resource {
	return &schema.Resource{
		Read:   ReadTasks,
		Schema: tasksSchema,
	}
}

func ReadTasks(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	extractedTasks, err := client.Tasks.Show(ctx, sdk.NewShowTaskRequest().WithIn(&sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)}))
	if err != nil {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] tasks in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	tasks := make([]map[string]any, 0, len(extractedTasks))
	for _, task := range extractedTasks {
		taskMap := map[string]any{}

		taskMap["name"] = task.Name
		taskMap["database"] = task.DatabaseName
		taskMap["schema"] = task.SchemaName
		taskMap["comment"] = task.Comment
		taskMap["warehouse"] = task.Warehouse

		tasks = append(tasks, taskMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("tasks", tasks)
}
