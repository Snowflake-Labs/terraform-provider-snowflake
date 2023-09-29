package resources

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var networkPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the network policy; must be unique for the account in which the network policy is created.",
		ForceNew:    true,
	},
	"allowed_ip_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Specifies one or more IPv4 addresses (CIDR notation) that are allowed access to your Snowflake account",
	},
	// TODO: Add a ValidationFunc to ensure 0.0.0.0/0 is not in blocked_ip_list
	// See: https://docs.snowflake.com/en/user-guide/network-policies.html#create-an-account-level-network-policy
	"blocked_ip_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies one or more IPv4 addresses (CIDR notation) that are denied access to your Snowflake account<br><br>**Do not** add `0.0.0.0/0` to `blocked_ip_list`",
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
		Create: CreateNetworkPolicy,
		Read:   ReadNetworkPolicy,
		Update: UpdateNetworkPolicy,
		Delete: DeleteNetworkPolicy,

		Schema: networkPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateNetworkPolicy implements schema.CreateFunc.
func CreateNetworkPolicy(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	req := sdk.NewCreateNetworkPolicyRequest(sdk.NewAccountObjectIdentifier(name))

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		req = req.WithComment(sdk.String(v.(string)))
	}

	if v, ok := d.GetOk("allowed_ip_list"); ok {
		ipList := expandStringList(v.(*schema.Set).List())
		ipRequests := make([]sdk.IPRequest, len(ipList))
		for i, v := range ipList {
			ipRequests[i] = *sdk.NewIPRequest(v)
		}
		req = req.WithAllowedIpList(ipRequests)
	}

	if v, ok := d.GetOk("blocked_ip_list"); ok {
		ipList := expandStringList(v.(*schema.Set).List())
		ipRequests := make([]sdk.IPRequest, len(ipList))
		for i, v := range ipList {
			ipRequests[i] = *sdk.NewIPRequest(v)
		}
		req = req.WithAllowedIpList(ipRequests)
	}

	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	err := client.NetworkPolicies.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("error creating network policy %v err = %w", name, err)
	}
	d.SetId(name)

	return ReadNetworkPolicy(d, meta)
}

// ReadNetworkPolicy implements schema.ReadFunc.
func ReadNetworkPolicy(d *schema.ResourceData, meta interface{}) error {
	policyName := d.Id()
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	networkPolicies, err := client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
	if err != nil {
		return err
	}

	var networkPolicy *sdk.NetworkPolicy
	for _, np := range networkPolicies {
		if np.Name == policyName {
			networkPolicy = &np
			break
		}
	}

	if networkPolicy == nil {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] network policy (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	builder := snowflake.NetworkPolicy(policyName)

	// There is no way to SHOW a single Network Policy, so we have to read *all* network policies and filter in memory
	showSQL := builder.ShowAllNetworkPolicies()

	rows, err := snowflake.Query(db, showSQL)

	allPolicies, err := snowflake.ScanNetworkPolicies(rows)
	if err != nil {
		return err
	}

	var s *snowflake.NetworkPolicyStruct
	for _, value := range allPolicies {
		if value.Name.String == policyName {
			s = value
		}
	}

	if s == nil {
		// The network policy was not found, the Terraform state does not reflect the Snowflake state
		log.Printf("[DEBUG] network policy (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	descSQL := builder.Describe()
	rows, err = snowflake.Query(db, descSQL)
	if err != nil {
		return err
	}

	err = d.Set("name", s.Name.String)
	if err != nil {
		return err
	}

	err = d.Set("comment", s.Comment.String)
	if err != nil {
		return err
	}

	var (
		name  string
		value string
	)
	for rows.Next() {
		err := rows.Scan(&name, &value)
		if err != nil {
			return err
		}

		if name == "ALLOWED_IP_LIST" {
			err = d.Set("allowed_ip_list", strings.Split(value, ","))
			if err != nil {
				return err
			}
		} else if name == "BLOCKED_IP_LIST" {
			err = d.Set("blocked_ip_list", strings.Split(value, ","))
			if err != nil {
				return err
			}
		}
	}

	return err
}

// UpdateNetworkPolicy implements schema.UpdateFunc.
func UpdateNetworkPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Id()
	builder := snowflake.NetworkPolicy(name)

	if d.HasChange("comment") {
		comment := d.Get("comment")

		if c := comment.(string); c == "" {
			q := builder.RemoveComment()
			err := snowflake.Exec(db, q)
			if err != nil {
				return fmt.Errorf("error unsetting comment for network policy %v err = %w", name, err)
			}
		} else {
			q := builder.ChangeComment(c)
			err := snowflake.Exec(db, q)
			if err != nil {
				return fmt.Errorf("error updating comment for network policy %v err = %w", name, err)
			}
		}
	}

	if d.HasChange("allowed_ip_list") {
		newIps := ipChangeParser(d, "allowed_ip_list")
		q := builder.ChangeIPList("ALLOWED", newIps)
		err := snowflake.Exec(db, q)
		if err != nil {
			return fmt.Errorf("error updating ALLOWED_IP_LIST for network policy %v err = %w", name, err)
		}
	}

	if d.HasChange("blocked_ip_list") {
		newIps := ipChangeParser(d, "blocked_ip_list")
		q := builder.ChangeIPList("BLOCKED", newIps)
		err := snowflake.Exec(db, q)
		if err != nil {
			return fmt.Errorf("error updating BLOCKED_IP_LIST for network policy %v err = %w", name, err)
		}
	}

	return ReadNetworkPolicy(d, meta)
}

// DeleteNetworkPolicy implements schema.DeleteFunc.
func DeleteNetworkPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Id()

	dropSQL := snowflake.NetworkPolicy(name).Drop()
	err := snowflake.Exec(db, dropSQL)
	if err != nil {
		return fmt.Errorf("error deleting network policy %v err = %w", name, err)
	}

	d.SetId("")
	return nil
}

// ipChangeParser is a helper function to parse a given ip list change from ResourceData.
func ipChangeParser(data *schema.ResourceData, key string) []string {
	ipChangeSet := data.Get(key)
	ipList := ipChangeSet.(*schema.Set).List()
	newIps := make([]string, len(ipList))
	for idx, value := range ipList {
		newIps[idx] = fmt.Sprintf("%v", value)
	}
	return newIps
}
