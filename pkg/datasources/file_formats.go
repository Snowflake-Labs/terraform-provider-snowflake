package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var fileFormatsSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the file formats from.",
	},
	"file_formats": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The file formats in the schema",
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
				"format_type": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func FileFormats() *schema.Resource {
	return &schema.Resource{
		Read:   ReadFileFormats,
		Schema: fileFormatsSchema,
	}
}

func ReadFileFormats(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	currentFileFormats, err := snowflake.ListFileFormats(databaseName, schemaName, db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] file formats in schema (%s) not found", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse file formats in schema (%s)", d.Id())
		d.SetId("")
		return nil
	}

	fileFormats := []map[string]interface{}{}

	for _, fileFormat := range currentFileFormats {
		fileFormatMap := map[string]interface{}{}

		fileFormatMap["name"] = fileFormat.FileFormatName.String
		fileFormatMap["database"] = fileFormat.DatabaseName.String
		fileFormatMap["schema"] = fileFormat.SchemaName.String
		fileFormatMap["comment"] = fileFormat.Comment.String
		fileFormatMap["format_type"] = fileFormat.FormatType.String

		fileFormats = append(fileFormats, fileFormatMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("file_formats", fileFormats)
}
