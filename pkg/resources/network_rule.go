package resources

import (
	"context"
	"errors"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var networkRuleSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the network rule; must be unique for the database and schema in which the network rule is created.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the network rule.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the network rule.",
	},
	"type": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Specifies the type of network identifiers being allowed or blocked. A network rule can have only one type. Allowed values are IPV4, AWSVPCEID, AZURELINKID and HOST_PORT; see https://docs.snowflake.com/en/sql-reference/sql/create-network-rule#required-parameters for details.",
		ValidateFunc: validation.StringInSlice([]string{"IPV4", "AWSVPCEID", "AZURELINKID", "HOST_PORT"}, false),
	},
	"value_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Specifies the network identifiers that will be allowed or blocked. Valid values in the list are determined by the type of network rule, see https://docs.snowflake.com/en/sql-reference/sql/create-network-rule#required-parameters for details.",
	},
	"mode": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Specifies what is restricted by the network rule. Valid values are INGRESS, INTERNAL_STAGE and EGRESS; see https://docs.snowflake.com/en/sql-reference/sql/create-network-rule#required-parameters for details.",
		ValidateFunc: validation.StringInSlice([]string{"INGRESS", "INTERNAL_STAGE", "EGRESS"}, false),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the network rule.",
	},
	"qualified_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Qualified name of the network rule.",
	},
}

// NetworkRule returns a pointer to the resource representing a network rule.
func NetworkRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextNetworkRule,
		ReadContext:   ReadContextNetworkRule,
		UpdateContext: UpdateContextNetworkRule,
		DeleteContext: DeleteContextNetworkRule,

		Schema: networkRuleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.Sequence(
			func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
				// Plan time validation for az_mode
				// InvalidParameterCombination: Must specify at least two cache nodes in order to specify AZ Mode of 'cross-az'.

				ruleTypeRaw, ok := diff.GetOk("type")
				if !ok {
					return nil
				}
				ruleType := sdk.NetworkRuleType(ruleTypeRaw.(string))
				ruleModeRaw, ok := diff.GetOk("mode")
				if !ok {
					return nil
				}
				ruleMode := sdk.NetworkRuleMode(ruleModeRaw.(string))
				//valueListRaw, ok := diff.GetOk("value_list")
				//if !ok {
				//	return nil
				//}
				//
				//valueList := expandStringList(valueListRaw.(*schema.Set).List())

				if ruleType == sdk.NetworkRuleTypeIpv4 && ruleMode == sdk.NetworkRuleModeEgress {
					return errors.New("the network rule mode EGRESS is not supported by the network rule type IPv4. The network rule mode must be one of [INGRESS].")
				}

				return nil
			},
		),
	}
}

func CreateContextNetworkRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	ruleType := sdk.NetworkRuleType(d.Get("type").(string))
	ruleMode := sdk.NetworkRuleMode(d.Get("mode").(string))

	valueList := expandStringList(d.Get("value_list").(*schema.Set).List())
	networkRuleValues := make([]sdk.NetworkRuleValue, len(valueList))
	for i, v := range valueList {
		networkRuleValues[i] = sdk.NetworkRuleValue{Value: v}
	}

	req := sdk.NewCreateNetworkRuleRequest(
		id,
		ruleType,
		networkRuleValues,
		ruleMode,
	)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		req = req.WithComment(sdk.String(v.(string)))
	}

	client := meta.(*provider.Context).Client
	if err := client.NetworkRules.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadContextNetworkRule(ctx, d, meta)
}

func ReadContextNetworkRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	networkRule, err := client.NetworkRules.ShowByID(ctx, id)
	if networkRule == nil || err != nil {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] network rule (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	networkRuleDescriptions, err := client.NetworkRules.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", networkRule.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("database", networkRule.DatabaseName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("schema", networkRule.SchemaName); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("type", networkRule.Type); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("value_list", networkRuleDescriptions.ValueList); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("mode", networkRule.Mode); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("comment", networkRule.Comment); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("qualified_name", id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func UpdateContextNetworkRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	baseReq := sdk.NewAlterNetworkRuleRequest(id)

	if d.HasChange("comment") || d.HasChange("value_list") {
		valueList := expandStringList(d.Get("value_list").(*schema.Set).List())
		networkRuleValues := make([]sdk.NetworkRuleValue, len(valueList))
		for i, v := range valueList {
			networkRuleValues[i] = sdk.NetworkRuleValue{Value: v}
		}

		setReq := sdk.NewNetworkRuleSetRequest(networkRuleValues)

		if d.HasChange("comment") {
			comment := d.Get("comment").(string)
			setReq.WithComment(sdk.String(comment))
		}
		if err := client.NetworkRules.Alter(ctx, baseReq.WithSet(setReq)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextNetworkRule(ctx, d, meta)
}

func DeleteContextNetworkRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	if err := client.NetworkRules.Drop(ctx, sdk.NewDropNetworkRuleRequest(id)); err != nil {
		diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

//// ipChangeParser is a helper function to parse a given ip list change from ResourceData.
//func ipChangeParser(data *schema.ResourceData, key string) []string {
//	ipChangeSet := data.Get(key)
//	ipList := ipChangeSet.(*schema.Set).List()
//	newIps := make([]string, len(ipList))
//	for idx, value := range ipList {
//		newIps[idx] = fmt.Sprintf("%v", value)
//	}
//	return newIps
//}
