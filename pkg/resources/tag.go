package resources

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
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
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// TODO: remove after rework of external table, materialized views stage and table
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

func CreateContextTag(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schema, name)

	request := sdk.NewCreateTagRequest(id)
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(sdk.String(v.(string)))
	}
	if v, ok := d.GetOk("allowed_values"); ok {
		request.WithAllowedValues(expandStringListAllowEmpty(v.([]any)))
	}
	if err := client.Tags.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeSnowflakeID(id))
	return ReadContextTag(ctx, d, meta)
}

func ReadContextTag(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	tag, err := client.Tags.ShowByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", tag.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("database", tag.DatabaseName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", tag.SchemaName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("comment", tag.Comment); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allowed_values", tag.AllowedValues); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func UpdateContextTag(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
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
		o, n := d.GetChange("allowed_values")
		oldAllowedValues := expandStringListAllowEmpty(o.([]any))
		newAllowedValues := expandStringListAllowEmpty(n.([]any))
		var allowedValuesToAdd, allowedValuesToRemove []string

		for _, oldAllowedValue := range oldAllowedValues {
			if !slices.Contains(newAllowedValues, oldAllowedValue) {
				allowedValuesToRemove = append(allowedValuesToRemove, oldAllowedValue)
			}
		}

		for _, newAllowedValue := range newAllowedValues {
			if !slices.Contains(oldAllowedValues, newAllowedValue) {
				allowedValuesToAdd = append(allowedValuesToAdd, newAllowedValue)
			}
		}

		if len(allowedValuesToAdd) > 0 {
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithAdd(allowedValuesToAdd)); err != nil {
				return diag.FromErr(err)
			}
		}

		if len(allowedValuesToRemove) > 0 {
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithDrop(allowedValuesToRemove)); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return ReadContextTag(ctx, d, meta)
}

func DeleteContextTag(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	if err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

// Returns the slice of strings for inputed allowed values.
func expandAllowedValues(avChangeSet any) []string {
	avList := avChangeSet.([]any)
	newAvs := make([]string, len(avList))
	for idx, value := range avList {
		newAvs[idx] = fmt.Sprintf("%v", value)
	}

	return newAvs
}
