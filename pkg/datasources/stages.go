package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var stagesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the stages from.",
	},
	"stages": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The stages in the schema",
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
				"storage_integration": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func Stages() *schema.Resource {
	return &schema.Resource{
		Read:   ReadStages,
		Schema: stagesSchema,
	}
}

func ReadStages(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentStages, err := snowflake.ListStages(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] stages in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse stages in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	stages := []map[string]interface{}{}

	for _, stage := range currentStages {
		stageMap := map[string]interface{}{}

		stageMap["name"] = stage.Name
		stageMap["database"] = stage.DatabaseName
		stageMap["schema"] = stage.SchemaName
		stageMap["comment"] = stage.Comment
		stageMap["storage_integration"] = stage.StorageIntegration

		stages = append(stages, stageMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("stages", stages)
}
