package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
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
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentTasks, err := snowflake.ListTasks(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] tasks in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse tasks in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	tasks := []map[string]interface{}{}

	for _, task := range currentTasks {
		taskMap := map[string]interface{}{}

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
