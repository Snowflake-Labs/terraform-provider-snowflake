package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sequencesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the sequences from.",
	},
	"sequences": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The sequences in the schema",
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
			},
		},
	},
}

func Sequences() *schema.Resource {
	return &schema.Resource{
		Read:   ReadSequences,
		Schema: sequencesSchema,
	}
}

func ReadSequences(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentSequences, err := snowflake.ListSequences(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] sequences in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse sequences in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	sequences := []map[string]interface{}{}

	for _, sequence := range currentSequences {
		sequenceMap := map[string]interface{}{}

		sequenceMap["name"] = sequence.Name.String
		sequenceMap["database"] = sequence.DBName.String
		sequenceMap["schema"] = sequence.SchemaName.String
		sequenceMap["comment"] = sequence.Comment.String

		sequences = append(sequences, sequenceMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("sequences", sequences)
}
