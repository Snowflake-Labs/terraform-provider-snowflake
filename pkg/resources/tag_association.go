package resources

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var tagAssociationSchema = map[string]*schema.Schema{
	"object_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the object identifier for the tag association.",
		ForceNew:    true,
		Deprecated:  "Use `object_identifier` instead",
	},
	"object_identifier": {
		Type:        schema.TypeSet,
		MinItems:    1,
		Required:    true,
		Description: "Specifies the object identifiers for the tag association.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("object_identifier"),
	},
	"object_type": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  fmt.Sprintf("Specifies the type of object to add a tag. Allowed object types: %v.", sdk.TagAssociationAllowedObjectTypesString),
		ValidateFunc: validation.StringInSlice(sdk.TagAssociationAllowedObjectTypesString, true),
		ForceNew:     true,
	},
	"tag_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the identifier for the tag.",
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"tag_value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the value of the tag, (e.g. 'finance' or 'engineering')",
		ForceNew:    true,
	},
	"skip_validation": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If true, skips validation of the tag association.",
		Default:     true,
	},
}

// TagAssociation returns a pointer to the resource representing a schema.
func TagAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextTagAssociation,
		ReadContext:   ReadContextTagAssociation,
		UpdateContext: UpdateContextTagAssociation,
		DeleteContext: DeleteContextTagAssociation,
		Description:   "Resource used to manage tag associations. For more information, check [object tagging documentation](https://docs.snowflake.com/en/user-guide/object-tagging).",

		Schema: tagAssociationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Minute),
		},
	}
}

func TagIdentifierAndObjectIdentifier(d *schema.ResourceData) (sdk.SchemaObjectIdentifier, []sdk.ObjectIdentifier, sdk.ObjectType, error) {
	tag := d.Get("tag_id").(string)
	tagId, err := sdk.ParseSchemaObjectIdentifier(tag)
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, nil, "", fmt.Errorf("invalid tag id: %w", err)
	}

	objectType := sdk.ObjectType(d.Get("object_type").(string))

	idsRaw := expandStringList(d.Get("object_identifier").(*schema.Set).List())
	ids := make([]sdk.ObjectIdentifier, len(idsRaw))
	for i, idRaw := range idsRaw {
		id, err := sdk.ParseObjectIdentifierString(idRaw)
		if err != nil {
			return sdk.SchemaObjectIdentifier{}, nil, "", fmt.Errorf("invalid object id: %w", err)
		}
		ids[i] = id
	}
	return tagId, ids, objectType, nil
}

func CreateContextTagAssociation(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	tagValue := d.Get("tag_value").(string)

	tagId, ids, objectType, err := TagIdentifierAndObjectIdentifier(d)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, oid := range ids {
		request := sdk.NewSetTagRequest(objectType, oid).WithSetTags([]sdk.TagAssociation{
			{
				Name:  tagId,
				Value: tagValue,
			},
		})
		if err := client.Tags.Set(ctx, request); err != nil {
			return diag.FromErr(err)
		}
		skipValidate := d.Get("skip_validation").(bool)
		if !skipValidate {
			log.Println("[DEBUG] validating tag creation")
			if err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate)-time.Minute, func() *retry.RetryError {
				tag, err := client.SystemFunctions.GetTag(ctx, tagId, oid, objectType)
				if err != nil {
					return retry.NonRetryableError(fmt.Errorf("error getting tag: %w", err))
				}
				// if length of response is zero, tag association was not found. retry
				if tag == nil {
					return retry.RetryableError(fmt.Errorf("expected tag association to be created but not yet created"))
				}
				return nil
			}); err != nil {
				return diag.FromErr(fmt.Errorf("error validating tag creation: %w", err))
			}
		}
	}
	d.SetId(helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue, string(objectType)))
	return ReadContextTagAssociation(ctx, d, meta)
}

func ReadContextTagAssociation(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	tagValue := d.Get("tag_value").(string)

	tagId, ids, objectType, err := TagIdentifierAndObjectIdentifier(d)
	if err != nil {
		return diag.FromErr(err)
	}
	var correctObjectIds []string
	for _, oid := range ids {
		objectTagValue, err := client.SystemFunctions.GetTag(ctx, tagId, oid, objectType)
		if err != nil {
			return diag.FromErr(err)
		}
		if objectTagValue != nil && *objectTagValue == tagValue {
			correctObjectIds = append(correctObjectIds, oid.FullyQualifiedName())
		}
	}
	if err := d.Set("object_identifier", correctObjectIds); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateContextTagAssociation(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	tagId, _, objectType, err := TagIdentifierAndObjectIdentifier(d)
	if err != nil {
		return diag.FromErr(err)
	}
	tagValue := d.Get("tag_value").(string)
	if d.HasChange("object_identifier") {
		o, n := d.GetChange("object_identifier")

		oldAllowedValues, err := expandStringListWithMapping(o.(*schema.Set).List(), sdk.ParseObjectIdentifierString)
		if err != nil {
			return diag.FromErr(err)
		}
		newAllowedValues, err := expandStringListWithMapping(n.(*schema.Set).List(), sdk.ParseObjectIdentifierString)
		if err != nil {
			return diag.FromErr(err)
		}

		addedids, removedids := ListDiff(oldAllowedValues, newAllowedValues)

		for _, id := range addedids {
			request := sdk.NewSetTagRequest(objectType, id).WithSetTags([]sdk.TagAssociation{
				{
					Name:  tagId,
					Value: tagValue,
				},
			})
			if err := client.Tags.Set(ctx, request); err != nil {
				return diag.FromErr(err)
			}
		}

		for _, id := range removedids {
			request := sdk.NewUnsetTagRequest(objectType, id).WithUnsetTags([]sdk.ObjectIdentifier{tagId}).WithIfExists(true)
			if err := client.Tags.Unset(ctx, request); err != nil {
				return diag.FromErr(err)
			}
		}

	}

	return ReadContextTagAssociation(ctx, d, meta)
}

func DeleteContextTagAssociation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	tid, ids, ot, err := TagIdentifierAndObjectIdentifier(d)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, oid := range ids {
		request := sdk.NewUnsetTagRequest(ot, oid).WithUnsetTags([]sdk.ObjectIdentifier{tid}).WithIfExists(true)
		if err := client.Tags.Unset(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return nil
}
