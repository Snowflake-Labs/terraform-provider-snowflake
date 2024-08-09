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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
)

const (
	tagAttachmentPolicyIDDelimiter = "|"
)

var mpAttachmentPolicySchema = map[string]*schema.Schema{
	"tag_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the identifier for the tag. Note: format must follow: \"databaseName\".\"schemaName\".\"tagName\" or \"databaseName.schemaName.tagName\" or \"databaseName|schemaName.tagName\" (snowflake_tag.tag.id)",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
	"masking_policy_id": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "The resource id of the masking policy",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
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
	client := meta.(*provider.Context).Client

	value := d.Get("tag_id").(string)
	tagObjectIdentifier, err := helpers.DecodeSnowflakeParameterID(value)
	if err != nil {
		return diag.FromErr(err)
	}
	tagId := tagObjectIdentifier.(sdk.SchemaObjectIdentifier)

	value = d.Get("masking_policy_id").(string)
	maskingPolicyObjectIdentifier, err := helpers.DecodeSnowflakeParameterID(value)
	if err != nil {
		return diag.FromErr(err)
	}
	maskingPolicyId := maskingPolicyObjectIdentifier.(sdk.SchemaObjectIdentifier)

	set := sdk.NewTagSetRequest().WithMaskingPolicies([]sdk.SchemaObjectIdentifier{maskingPolicyId})
	if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(tagId).WithSet(set)); err != nil {
		return diag.FromErr(err)
	}
	aid := attachmentID{
		TagDatabaseName:           tagId.DatabaseName(),
		TagSchemaName:             tagId.SchemaName(),
		TagName:                   tagId.Name(),
		MaskingPolicyDatabaseName: maskingPolicyId.DatabaseName(),
		MaskingPolicySchemaName:   maskingPolicyId.SchemaName(),
		MaskingPolicyName:         maskingPolicyId.Name(),
	}
	fmt.Printf("attachment id: %s\n", aid.String())
	d.SetId(aid.String())
	return ReadContextTagMaskingPolicyAssociation(ctx, d, meta)
}

func ReadContextTagMaskingPolicyAssociation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	aid, err := parseAttachmentID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// create temp warehouse to query the tag, and make sure to clean it up
	cleanupWarehouse, err := ensureWarehouse(ctx, client)
	if err != nil {
		return diag.FromErr(err)
	}
	defer cleanupWarehouse()
	// show attached masking policy
	tid := sdk.NewSchemaObjectIdentifier(aid.TagDatabaseName, aid.TagSchemaName, aid.TagName)
	mid := sdk.NewSchemaObjectIdentifier(aid.MaskingPolicyDatabaseName, aid.MaskingPolicySchemaName, aid.MaskingPolicyName)
	builder := snowflake.NewTagBuilder(tid).WithMaskingPolicy(mid)
	row := snowflake.QueryRow(db, builder.ShowAttachedPolicy())
	_, err = snowflake.ScanTagPolicy(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] attached policy (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func DeleteContextTagMaskingPolicyAssociation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
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
