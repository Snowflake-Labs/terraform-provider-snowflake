package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
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
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentStreams, err := snowflake.ListStreams(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] streams in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse streams in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	streams := []map[string]interface{}{}

	for _, stream := range currentStreams {
		streamMap := map[string]interface{}{}

		streamMap["name"] = stream.StreamName.String
		streamMap["database"] = stream.DatabaseName.String
		streamMap["schema"] = stream.SchemaName.String
		streamMap["comment"] = stream.Comment.String
		streamMap["table"] = stream.TableName.String

		streams = append(streams, streamMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("streams", streams)
}
