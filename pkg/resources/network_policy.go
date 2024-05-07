package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
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
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies a list of network rules that contain the network identifiers that are allowed access to Snowflake.",
		// TODO: Add a ValidationFunc to ensure that each entry in the list is a fully qualified name
	},
	"blocked_network_rule_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies a list of network rules that contain the network identifiers that are denied access to Snowflake.",
		// TODO: Add a ValidationFunc to ensure that each entry in the list is a fully qualified name
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
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] network policy (%s) not found", d.Id())
		d.SetId("")
		return nil
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
		oldList, newList := d.GetChange("allowed_network_rule_list")
		addedNetworkRuleIdentifiers, removedNetworkRuleIdentifiers := getAddedAndRemovedIdentifiers(oldList, newList)

		if len(addedNetworkRuleIdentifiers) > 0 {
			baseReq := sdk.NewAlterNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))
			addReq := sdk.NewAddNetworkRuleRequest().WithAllowedNetworkRuleList(addedNetworkRuleIdentifiers)
			err := client.NetworkPolicies.Alter(ctx, baseReq.WithAdd(addReq))

			if err != nil {
				return getUpdateContextDiag("adding to ALLOWED_NETWORK_RULE_LIST", name, err)
			}
		}
		if len(removedNetworkRuleIdentifiers) > 0 {
			baseReq := sdk.NewAlterNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))
			removeReq := sdk.NewRemoveNetworkRuleRequest().WithAllowedNetworkRuleList(removedNetworkRuleIdentifiers)
			err := client.NetworkPolicies.Alter(ctx, baseReq.WithRemove(removeReq))
			if err != nil {
				return getUpdateContextDiag("removing from ALLOWED_NETWORK_RULE_LIST", name, err)
			}
		}
	}

	if d.HasChange("blocked_network_rule_list") {
		oldList, newList := d.GetChange("blocked_network_rule_list")
		addedNetworkRuleIdentifiers, removedNetworkRuleIdentifiers := getAddedAndRemovedIdentifiers(oldList, newList)

		if len(addedNetworkRuleIdentifiers) > 0 {
			baseReq := sdk.NewAlterNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))
			addReq := sdk.NewAddNetworkRuleRequest().WithBlockedNetworkRuleList(addedNetworkRuleIdentifiers)
			err := client.NetworkPolicies.Alter(ctx, baseReq.WithAdd(addReq))

			if err != nil {
				return getUpdateContextDiag("adding to BLOCKED_NETWORK_RULE_LIST", name, err)
			}
		}
		if len(removedNetworkRuleIdentifiers) > 0 {
			baseReq := sdk.NewAlterNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))
			removeReq := sdk.NewRemoveNetworkRuleRequest().WithBlockedNetworkRuleList(removedNetworkRuleIdentifiers)
			err := client.NetworkPolicies.Alter(ctx, baseReq.WithRemove(removeReq))
			if err != nil {
				return getUpdateContextDiag("removing from BLOCKED_NETWORK_RULE_LIST", name, err)
			}
		}
	}

	if d.HasChange("allowed_ip_list") {
		baseReq := sdk.NewAlterNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))
		ipRequests := parseIPList(d.Get("allowed_ip_list"))

		var err error
		if len(ipRequests) == 0 {
			unsetReq := sdk.NewNetworkPolicyUnsetRequest().WithAllowedIpList(sdk.Bool(true))
			err = client.NetworkPolicies.Alter(ctx, baseReq.WithUnset(unsetReq))
		} else {
			setReq := sdk.NewNetworkPolicySetRequest().WithAllowedIpList(sdk.NewAllowedIPListRequest().WithAllowedIPList(ipRequests))
			err = client.NetworkPolicies.Alter(ctx, baseReq.WithSet(setReq))
		}

		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Error updating network policy",
					Detail:   fmt.Sprintf("error updating ALLOWED_IP_LIST for network policy %v err = %v", name, err),
				},
			}
		}
	}

	if d.HasChange("blocked_ip_list") {
		baseReq := sdk.NewAlterNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))
		ipRequests := parseIPList(d.Get("blocked_ip_list"))

		var err error
		if len(ipRequests) == 0 {
			unsetReq := sdk.NewNetworkPolicyUnsetRequest().WithBlockedIpList(sdk.Bool(true))
			err = client.NetworkPolicies.Alter(ctx, baseReq.WithUnset(unsetReq))
		} else {
			setReq := sdk.NewNetworkPolicySetRequest().WithBlockedIpList(sdk.NewBlockedIPListRequest().WithBlockedIPList(ipRequests))
			err = client.NetworkPolicies.Alter(ctx, baseReq.WithSet(setReq))
		}

		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Error updating network policy",
					Detail:   fmt.Sprintf("error updating BLOCKED_IP_LIST for network policy %v err = %v", name, err),
				},
			}
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

	err := client.NetworkPolicies.Drop(ctx, sdk.NewDropNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name)))
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

// getAddedAndRemovedIdentifiers returns the identifiers that were added and removed from the old and new network rule lists.
func getAddedAndRemovedIdentifiers(oldList interface{}, newList interface{}) ([]sdk.SchemaObjectIdentifier, []sdk.SchemaObjectIdentifier) {
	oldNetworkRuleIdentifiers := parseNetworkRulesList(oldList)
	newNetworkRuleIdentifiers := parseNetworkRulesList(newList)

	var addedNetworkRuleIdentifiers []sdk.SchemaObjectIdentifier
	var removedNetworkRuleIdentifiers []sdk.SchemaObjectIdentifier

	for _, identifier := range oldNetworkRuleIdentifiers {
		if !contains(newNetworkRuleIdentifiers, identifier) {
			removedNetworkRuleIdentifiers = append(removedNetworkRuleIdentifiers, identifier)
		}
	}
	log.Printf("removedNetworkRuleIdentifiers: %v", removedNetworkRuleIdentifiers)
	for _, identifier := range newNetworkRuleIdentifiers {
		if !contains(oldNetworkRuleIdentifiers, identifier) {
			addedNetworkRuleIdentifiers = append(addedNetworkRuleIdentifiers, identifier)
		}
	}
	return addedNetworkRuleIdentifiers, removedNetworkRuleIdentifiers
}

// contains checks if a given identifier is in a list of identifiers.
func contains(identifierList []sdk.SchemaObjectIdentifier, identifier sdk.SchemaObjectIdentifier) bool {
	for _, objectIdentifier := range identifierList {
		if objectIdentifier.FullyQualifiedName() == identifier.FullyQualifiedName() {
			return true
		}
	}
	return false
}
