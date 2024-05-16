package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var networkPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the network policy; must be unique for the account in which the network policy is created.",
		ForceNew:    true,
	},
	"allowed_network_rule_list": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		},
		Optional:    true,
		Description: "Specifies a list of fully qualified network rules that contain the network identifiers that are allowed access to Snowflake.",
	},
	"blocked_network_rule_list": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		},
		Optional:    true,
		Description: "Specifies a list of fully qualified network rules that contain the network identifiers that are denied access to Snowflake.",
	},
	"allowed_ip_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies one or more IPv4 addresses (CIDR notation) that are allowed access to your Snowflake account.",
	},
	// TODO: Add a ValidationFunc to ensure 0.0.0.0/0 is not in blocked_ip_list
	// See: https://docs.snowflake.com/en/user-guide/network-policies.html#create-an-account-level-network-policy
	"blocked_ip_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies one or more IPv4 addresses (CIDR notation) that are denied access to your Snowflake account<br><br>**Do not** add `0.0.0.0/0` to `blocked_ip_list`.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the network policy.",
	},
}

// NetworkPolicy returns a pointer to the resource representing a network policy.
func NetworkPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextNetworkPolicy,
		ReadContext:   ReadContextNetworkPolicy,
		UpdateContext: UpdateContextNetworkPolicy,
		DeleteContext: DeleteContextNetworkPolicy,

		Schema: networkPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextNetworkPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	req := sdk.NewCreateNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))

	if v, ok := d.GetOk("comment"); ok {
		req = req.WithComment(sdk.String(v.(string)))
	}

	if v, ok := d.GetOk("allowed_network_rule_list"); ok {
		networkRuleIdentifiers := parseNetworkRulesList(v)
		req = req.WithAllowedNetworkRuleList(networkRuleIdentifiers)
	}

	if v, ok := d.GetOk("blocked_network_rule_list"); ok {
		networkRuleIdentifiers := parseNetworkRulesList(v)
		req = req.WithBlockedNetworkRuleList(networkRuleIdentifiers)
	}

	if v, ok := d.GetOk("allowed_ip_list"); ok {
		ipRequests := parseIPList(v)
		req = req.WithAllowedIpList(ipRequests)
	}

	if v, ok := d.GetOk("blocked_ip_list"); ok {
		ipRequests := parseIPList(v)
		req = req.WithAllowedIpList(ipRequests)
	}

	client := meta.(*provider.Context).Client
	err := client.NetworkPolicies.Create(ctx, req)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error creating network policy",
				Detail:   fmt.Sprintf("error creating network policy %v err = %v", name, err),
			},
		}
	}
	d.SetId(name)

	return ReadContextNetworkPolicy(ctx, d, meta)
}

// NetworkRulesSnowflakeDTO is needed to unpack the applied network rules from the JSON response from Snowflake
type NetworkRulesSnowflakeDTO struct {
	FullyQualifiedRuleName string
}

func ReadContextNetworkPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	policyName := d.Id()
	client := meta.(*provider.Context).Client

	networkPolicy, err := client.NetworkPolicies.ShowByID(ctx, sdk.NewAccountObjectIdentifier(policyName))
	if networkPolicy == nil || err != nil {
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve network policy. Target object not found. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Id: %s", d.Id()),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve network policy",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}

	if err = d.Set("name", networkPolicy.Name); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("comment", networkPolicy.Comment); err != nil {
		return diag.FromErr(err)
	}

	policyDescriptions, err := client.NetworkPolicies.Describe(ctx, sdk.NewAccountObjectIdentifier(policyName))
	if err != nil {
		return diag.FromErr(err)
	}
	for _, desc := range policyDescriptions {
		switch desc.Name {
		case "ALLOWED_IP_LIST":
			if err = d.Set("allowed_ip_list", strings.Split(desc.Value, ",")); err != nil {
				return diag.FromErr(err)
			}
		case "BLOCKED_IP_LIST":
			if err = d.Set("blocked_ip_list", strings.Split(desc.Value, ",")); err != nil {
				return diag.FromErr(err)
			}
		case "ALLOWED_NETWORK_RULE_LIST":
			var networkRules []NetworkRulesSnowflakeDTO
			err := json.Unmarshal([]byte(desc.Value), &networkRules)
			if err != nil {
				return diag.FromErr(err)
			}
			networkRulesFullyQualified := make([]string, len(networkRules))
			for i, ele := range networkRules {
				networkRulesFullyQualified[i] = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(ele.FullyQualifiedRuleName).FullyQualifiedName()
			}

			if err = d.Set("allowed_network_rule_list", networkRulesFullyQualified); err != nil {
				return diag.FromErr(err)
			}
		case "BLOCKED_NETWORK_RULE_LIST":
			var networkRules []NetworkRulesSnowflakeDTO
			err := json.Unmarshal([]byte(desc.Value), &networkRules)
			if err != nil {
				return diag.FromErr(err)
			}
			networkRulesFullyQualified := make([]string, len(networkRules))
			for i, ele := range networkRules {
				networkRulesFullyQualified[i] = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(ele.FullyQualifiedRuleName).FullyQualifiedName()
			}

			if err = d.Set("blocked_network_rule_list", networkRulesFullyQualified); err != nil {
				return diag.FromErr(err)

			}
		}
	}

	return diags
}

func UpdateContextNetworkPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Id()
	client := meta.(*provider.Context).Client

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		baseReq := sdk.NewAlterNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))

		if comment == "" {
			unsetReq := sdk.NewNetworkPolicyUnsetRequest().WithComment(sdk.Bool(true))
			err := client.NetworkPolicies.Alter(ctx, baseReq.WithUnset(unsetReq))
			if err != nil {
				return getUpdateContextDiag("unsetting comment", name, err)
			}
		} else {
			setReq := sdk.NewNetworkPolicySetRequest().WithComment(sdk.String(comment))
			err := client.NetworkPolicies.Alter(ctx, baseReq.WithSet(setReq))
			if err != nil {
				return getUpdateContextDiag("updating comment", name, err)
			}
		}
	}

	if d.HasChange("allowed_network_rule_list") {
		baseReq := sdk.NewAlterNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))
		networkRuleIdentifiers := parseNetworkRulesList(d.Get("allowed_network_rule_list"))
		setReq := sdk.NewNetworkPolicySetRequest().WithAllowedNetworkRuleList(sdk.NewAllowedNetworkRuleListRequest().WithAllowedNetworkRuleList(networkRuleIdentifiers))
		err := client.NetworkPolicies.Alter(ctx, baseReq.WithSet(setReq))

		if err != nil {
			return getUpdateContextDiag("updating ALLOWED_NETWORK_RULE_LIST", name, err)
		}
	}

	if d.HasChange("blocked_network_rule_list") {
		baseReq := sdk.NewAlterNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))
		networkRuleIdentifiers := parseNetworkRulesList(d.Get("blocked_network_rule_list"))
		setReq := sdk.NewNetworkPolicySetRequest().WithBlockedNetworkRuleList(sdk.NewBlockedNetworkRuleListRequest().WithBlockedNetworkRuleList(networkRuleIdentifiers))
		err := client.NetworkPolicies.Alter(ctx, baseReq.WithSet(setReq))

		if err != nil {
			return getUpdateContextDiag("updating BLOCKED_NETWORK_RULE_LIST", name, err)
		}
	}

	if d.HasChange("allowed_ip_list") {
		baseReq := sdk.NewAlterNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))
		ipRequests := parseIPList(d.Get("allowed_ip_list"))
		setReq := sdk.NewNetworkPolicySetRequest().WithAllowedIpList(sdk.NewAllowedIPListRequest().WithAllowedIPList(ipRequests))
		err := client.NetworkPolicies.Alter(ctx, baseReq.WithSet(setReq))

		if err != nil {
			return getUpdateContextDiag("updating ALLOWED_IP_LIST", name, err)
		}
	}

	if d.HasChange("blocked_ip_list") {
		baseReq := sdk.NewAlterNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))
		ipRequests := parseIPList(d.Get("blocked_ip_list"))
		setReq := sdk.NewNetworkPolicySetRequest().WithBlockedIpList(sdk.NewBlockedIPListRequest().WithBlockedIPList(ipRequests))
		err := client.NetworkPolicies.Alter(ctx, baseReq.WithSet(setReq))

		if err != nil {
			return getUpdateContextDiag("updating BLOCKED_IP_LIST", name, err)
		}
	}

	return ReadContextNetworkPolicy(ctx, d, meta)
}

func getUpdateContextDiag(action string, name string, err error) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error updating network policy",
			Detail:   fmt.Sprintf("error %v for network policy %v err = %v", action, name, err),
		},
	}
}

func DeleteContextNetworkPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Id()
	client := meta.(*provider.Context).Client

	err := client.NetworkPolicies.Drop(ctx, sdk.NewDropNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name)).WithIfExists(sdk.Bool(true)))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting network policy",
				Detail:   fmt.Sprintf("error deleting network policy %v err = %v", name, err),
			},
		}
	}

	d.SetId("")
	return nil
}

// parseIPList is a helper function to parse a given ip list from ResourceData.
func parseIPList(v interface{}) []sdk.IPRequest {
	ipList := expandStringList(v.(*schema.Set).List())
	ipRequests := make([]sdk.IPRequest, len(ipList))
	for i, v := range ipList {
		ipRequests[i] = *sdk.NewIPRequest(v)
	}
	return ipRequests
}

// parseNetworkRulesList is a helper function to parse a given network rule list from ResourceData.
func parseNetworkRulesList(v interface{}) []sdk.SchemaObjectIdentifier {
	networkRules := expandStringList(v.(*schema.Set).List())
	networkRuleIdentifiers := make([]sdk.SchemaObjectIdentifier, len(networkRules))
	for i, v := range networkRules {
		networkRuleIdentifiers[i] = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(v)
	}
	return networkRuleIdentifiers
}
