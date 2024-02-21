package resources

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var tagSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the tag; must be unique for the database in which the tag is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the tag.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the tag.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the tag.",
	},
	"allowed_values": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "List of allowed values for the tag.",
	},
}

var tagReferenceSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	MinItems:    0,
	Description: "Definitions of a tag to associate with the resource.",
	Deprecated:  "Use the 'snowflake_tag_association' resource instead.",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag name, e.g. department.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag value, e.g. marketing_info.",
			},
			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the database that the tag was created in.",
			},
			"schema": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the schema that the tag was created in.",
			},
		},
	},
}

// Schema returns a pointer to the resource representing a schema.
func Tag() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextTag,
		ReadContext:   ReadContextTag,
		UpdateContext: UpdateContextTag,
		DeleteContext: DeleteContextTag,

		Schema: tagSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextTag(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schema, name)

	request := sdk.NewCreateTagRequest(id)
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(sdk.String(v.(string)))
	}
	if v, ok := d.GetOk("allowed_values"); ok {
		request.WithAllowedValues(expandStringList(v.([]interface{})))
	}
	if err := client.Tags.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeSnowflakeID(id))
	return ReadContextTag(ctx, d, meta)
}

func ReadContextTag(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	request := sdk.NewShowTagRequest().WithIn(&sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())}).WithLike(id.Name())
	tags, err := client.Tags.Show(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, item := range tags {
		if err := d.Set("name", item.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("database", item.DatabaseName); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("schema", item.SchemaName); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("comment", item.Comment); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("allowed_values", item.AllowedValues); err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func UpdateContextTag(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	if d.HasChange("comment") {
		comment, ok := d.GetOk("comment")
		if ok {
			set := sdk.NewTagSetRequest().WithComment(comment.(string))
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(set)); err != nil {
				return diag.FromErr(err)
			}
		} else {
			unset := sdk.NewTagUnsetRequest().WithComment(true)
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(unset)); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("allowed_values") {
		v, ok := d.GetOk("allowed_values")
		if ok {
			allowedValues := expandAllowedValues(v)
			// unset the allowed values
			unset := sdk.NewTagUnsetRequest().WithAllowedValues(true)
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(unset)); err != nil {
				return diag.FromErr(err)
			}
			// add the allowed values
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithAdd(allowedValues)); err != nil {
				return diag.FromErr(err)
			}
		} else {
			unset := sdk.NewTagUnsetRequest().WithAllowedValues(true)
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(unset)); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return ReadContextTag(ctx, d, meta)
}

func DeleteContextTag(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	if err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

// Returns the slice of strings for inputed allowed values.
func expandAllowedValues(avChangeSet interface{}) []string {
	avList := avChangeSet.([]interface{})
	newAvs := make([]string, len(avList))
	for idx, value := range avList {
		newAvs[idx] = fmt.Sprintf("%v", value)
	}

	return newAvs
}
