package datasources

import (
	"context"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	extractedPipes, err := client.Pipes.Show(ctx, &sdk.ShowPipeOptions{
		In: &sdk.In{
			Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName),
		},
	})
	if err != nil {
		log.Printf("[DEBUG] unable to parse pipes in schema (%s)", d.Id())
		d.SetId("")
		return err
	}

	pipes := make([]map[string]any, 0, len(extractedPipes))
	for _, pipe := range extractedPipes {
		pipeMap := map[string]any{}

		pipeMap["name"] = pipe.Name
		pipeMap["database"] = pipe.DatabaseName
		pipeMap["schema"] = pipe.SchemaName
		pipeMap["comment"] = pipe.Comment
		pipeMap["integration"] = pipe.Integration

		pipes = append(pipes, pipeMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("pipes", pipes)
}
