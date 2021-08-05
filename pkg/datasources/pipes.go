package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var pipesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the pipes from.",
	},
	"pipes": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The pipes in the schema",
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
				"integration": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func Pipes() *schema.Resource {
	return &schema.Resource{
		Read:   ReadPipes,
		Schema: pipesSchema,
	}
}

func ReadPipes(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentPipes, err := snowflake.ListPipes(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] pipes in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse pipes in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	pipes := []map[string]interface{}{}

	for _, pipe := range currentPipes {
		pipeMap := map[string]interface{}{}

		pipeMap["name"] = pipe.Name
		pipeMap["database"] = pipe.DatabaseName
		pipeMap["schema"] = pipe.SchemaName
		pipeMap["comment"] = pipe.Comment
		pipeMap["integration"] = pipe.Integration.String

		pipes = append(pipes, pipeMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("pipes", pipes)
}
