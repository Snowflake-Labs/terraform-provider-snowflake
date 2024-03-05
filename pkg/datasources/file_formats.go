package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	result, err := client.FileFormats.Show(ctx, &sdk.ShowFileFormatsOptions{
		In: &sdk.In{
			Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName),
		},
	})
	if err != nil {
		d.SetId("")
		return err
	}

	fileFormats := []map[string]interface{}{}

	for _, fileFormat := range result {
		fileFormatMap := map[string]interface{}{}

		fileFormatMap["name"] = fileFormat.Name.Name()
		fileFormatMap["database"] = fileFormat.Name.DatabaseName()
		fileFormatMap["schema"] = fileFormat.Name.SchemaName()
		fileFormatMap["comment"] = fileFormat.Comment
		fileFormatMap["format_type"] = fileFormat.Type

		fileFormats = append(fileFormats, fileFormatMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return d.Set("file_formats", fileFormats)
}
