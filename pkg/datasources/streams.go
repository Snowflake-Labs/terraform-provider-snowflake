package datasources

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var streamsSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the streams from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the streams from.",
	},
	"streams": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The streams in the schema",
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
				"table": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func Streams() *schema.Resource {
	return &schema.Resource{
		Read:   ReadStreams,
		Schema: streamsSchema,
	}
}

func ReadStreams(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentStreams, err := client.Streams.Show(ctx, sdk.NewShowStreamRequest().
		WithIn(&sdk.In{
			Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName),
		}))
	if err != nil {
		log.Printf("[DEBUG] streams in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	streams := make([]map[string]any, len(currentStreams))
	for i, stream := range currentStreams {
		streams[i] = map[string]any{
			"name":     stream.Name,
			"database": stream.DatabaseName,
			"schema":   stream.SchemaName,
			"comment":  stream.Comment,
			"table":    stream.TableName,
		}
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("streams", streams)
}
