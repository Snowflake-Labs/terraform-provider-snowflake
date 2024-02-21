package resources

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
)

const (
	tagAttachmentPolicyIDDelimiter = "|"
)

var mpAttachmentPolicySchema = map[string]*schema.Schema{
	"tag_id": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the tag. Note: format must follow: \"databaseName\".\"schemaName\".\"tagName\" or \"databaseName.schemaName.tagName\" or \"databaseName|schemaName.tagName\" (snowflake_tag.tag.id)",
		ForceNew:    true,
	},
	"masking_policy_id": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The resource id of the masking policy",
	},
}

type attachmentID struct {
	TagDatabaseName           string
	TagSchemaName             string
	TagName                   string
	MaskingPolicyDatabaseName string
	MaskingPolicySchemaName   string
	MaskingPolicyName         string
}

func (v *attachmentID) String() string {
	return strings.Join([]string{
		v.TagDatabaseName,
		v.TagSchemaName,
		v.TagName,
		v.MaskingPolicyDatabaseName,
		v.MaskingPolicySchemaName,
		v.MaskingPolicyName,
	}, tagAttachmentPolicyIDDelimiter)
}

func parseAttachmentID(id string) (*attachmentID, error) {
	parts := strings.Split(id, tagAttachmentPolicyIDDelimiter)
	if len(parts) != 6 {
		return nil, fmt.Errorf("6 fields allowed")
	}
	return &attachmentID{
		TagDatabaseName:           parts[0],
		TagSchemaName:             parts[1],
		TagName:                   parts[2],
		MaskingPolicyDatabaseName: parts[3],
		MaskingPolicySchemaName:   parts[4],
		MaskingPolicyName:         parts[5],
	}, nil
}

// Schema returns a pointer to the resource representing a schema.
func TagMaskingPolicyAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextTagMaskingPolicyAssociation,
		ReadContext:   ReadContextTagMaskingPolicyAssociation,
		DeleteContext: DeleteContextTagMaskingPolicyAssociation,

		Schema: mpAttachmentPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Attach a masking policy to a tag. Requires a current warehouse to be set. Either with SNOWFLAKE_WAREHOUSE env variable or in current session. If no warehouse is provided, a temporary warehouse will be created.",
	}
}

func CreateContextTagMaskingPolicyAssociation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	value := d.Get("tag_id").(string)
	tid := helpers.DecodeSnowflakeID(value).(sdk.SchemaObjectIdentifier)
	value = d.Get("masking_policy_id").(string)
	mid := helpers.DecodeSnowflakeID(value).(sdk.SchemaObjectIdentifier)

	set := sdk.NewTagSetRequest().WithMaskingPolicies([]sdk.SchemaObjectIdentifier{mid})
	if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(tid).WithSet(set)); err != nil {
		return diag.FromErr(err)
	}
	aid := attachmentID{
		TagDatabaseName:           tid.DatabaseName(),
		TagSchemaName:             tid.SchemaName(),
		TagName:                   tid.Name(),
		MaskingPolicyDatabaseName: mid.DatabaseName(),
		MaskingPolicySchemaName:   mid.SchemaName(),
		MaskingPolicyName:         mid.Name(),
	}
	fmt.Printf("attachment id: %s\n", aid.String())
	d.SetId(aid.String())
	return ReadContextTagMaskingPolicyAssociation(ctx, d, meta)
}

func ReadContextTagMaskingPolicyAssociation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	aid, err := parseAttachmentID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// create temp warehouse to query the tag, and make sure to clean it up
	warehouse, err := client.ContextFunctions.CurrentWarehouse(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	if warehouse == "" {
		log.Printf("[DEBUG] no current warehouse set, creating a temporary warehouse")
		randomWarehouseName := fmt.Sprintf("terraform-provider-snowflake-%v", helpers.RandomString())
		wid := sdk.NewAccountObjectIdentifier(randomWarehouseName)
		if err := client.Warehouses.Create(ctx, wid, nil); err != nil {
			return diag.FromErr(err)
		}
		defer func() {
			if err := client.Warehouses.Drop(ctx, wid, nil); err != nil {
				log.Printf("[WARN] error cleaning up temp warehouse %v", err)
			}
			if err := client.Sessions.UseWarehouse(ctx, sdk.NewAccountObjectIdentifier(warehouse)); err != nil {
				log.Printf("[WARN] error resetting warehouse %v", err)
			}
		}()
		if err := client.Sessions.UseWarehouse(ctx, wid); err != nil {
			return diag.FromErr(err)
		}
	}
	// show attached masking policy
	tid := sdk.NewSchemaObjectIdentifier(aid.TagDatabaseName, aid.TagSchemaName, aid.TagName)
	mid := sdk.NewSchemaObjectIdentifier(aid.MaskingPolicyDatabaseName, aid.MaskingPolicySchemaName, aid.MaskingPolicyName)
	builder := snowflake.NewTagBuilder(tid).WithMaskingPolicy(mid)
	row := snowflake.QueryRow(db, builder.ShowAttachedPolicy())
	t, err := snowflake.ScanTagPolicy(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] attached policy (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}
	id := helpers.EncodeSnowflakeID(t.PolicyDB.String, t.PolicySchema.String, t.PolicyName.String)
	if err := d.Set("masking_policy_id", id); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func DeleteContextTagMaskingPolicyAssociation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	aid, err := parseAttachmentID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	tid := sdk.NewSchemaObjectIdentifier(aid.TagDatabaseName, aid.TagSchemaName, aid.TagName)
	mid := sdk.NewSchemaObjectIdentifier(aid.MaskingPolicyDatabaseName, aid.MaskingPolicySchemaName, aid.MaskingPolicyName)
	unset := sdk.NewTagUnsetRequest().WithMaskingPolicies([]sdk.SchemaObjectIdentifier{mid})
	if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(tid).WithUnset(unset)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
