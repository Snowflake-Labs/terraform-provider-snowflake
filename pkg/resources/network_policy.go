package resources

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
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

// NetworkPolicy returns a pointer to the resource representing a network policy
func NetworkPolicy() *schema.Resource {
	return &schema.Resource{
		Create: CreateNetworkPolicy,
		Read:   ReadNetworkPolicy,
		Update: UpdateNetworkPolicy,
		Delete: DeleteNetworkPolicy,

		Schema: networkPolicySchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateNetworkPolicy implements schema.CreateFunc
func CreateNetworkPolicy(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	builder := snowflake.NetworkPolicy(name)

	// Set optionals
	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := data.GetOk("allowed_ip_list"); ok {
		builder.WithAllowedIpList(expandStringList(v.(*schema.Set).List()))
	}

	if v, ok := data.GetOk("blocked_ip_list"); ok {
		builder.WithBlockedIpList(expandStringList(v.(*schema.Set).List()))
	}

	stmt := builder.Create()
	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating network policy %v", name)
	}
	data.SetId(name)

	return ReadNetworkPolicy(data, meta)
}

// ReadNetworkPolicy implements schema.ReadFunc
func ReadNetworkPolicy(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	policyName := data.Id()

	builder := snowflake.NetworkPolicy(policyName)

	// There is no way to SHOW a single Network Policy, so we have to read *all* network policies and filter in memory
	showSql := builder.ShowAllNetworkPolicies()

	rows, err := snowflake.Query(db, showSql)
	if err != nil {
		return err
	}

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

	descSql := builder.Describe()
	rows, err = snowflake.Query(db, descSql)
	if err != nil {
		return err
	}

	err = data.Set("name", s.Name.String)
	if err != nil {
		return err
	}

	err = data.Set("comment", s.Comment.String)
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
			err = data.Set("allowed_ip_list", strings.Split(value, ","))
			if err != nil {
				return err
			}
		} else if name == "BLOCKED_IP_LIST" {
			err = data.Set("blocked_ip_list", strings.Split(value, ","))
			if err != nil {
				return err
			}
		}
	}

	return err
}

// UpdateNetworkPolicy implements schema.UpdateFunc
func UpdateNetworkPolicy(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Id()
	builder := snowflake.NetworkPolicy(name)

	if data.HasChange("comment") {
		_, comment := data.GetChange("comment")

		if c := comment.(string); c == "" {
			q := builder.RemoveComment()
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error unsetting comment for network policy %v", name)
			}
		} else {
			q := builder.ChangeComment(c)
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error updating comment for network policy %v", name)
			}
		}
	}

	if data.HasChange("allowed_ip_list") {
		newIps := ipChangeParser(data, "allowed_ip_list")
		q := builder.ChangeIpList("ALLOWED", newIps)
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating ALLOWED_IP_LIST for network policy %v", name)
		}
	}

	if data.HasChange("blocked_ip_list") {
		newIps := ipChangeParser(data, "blocked_ip_list")
		q := builder.ChangeIpList("BLOCKED", newIps)
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating BLOCKED_IP_LIST for network policy %v", name)
		}
	}

	return ReadNetworkPolicy(data, meta)
}

// DeleteNetworkPolicy implements schema.DeleteFunc
func DeleteNetworkPolicy(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Id()

	dropSql := snowflake.NetworkPolicy(name).Drop()
	err := snowflake.Exec(db, dropSql)
	if err != nil {
		return errors.Wrapf(err, "error deleting network policy %v", name)
	}

	data.SetId("")
	return nil
}

// ipChangeParser is a helper function to parse a given ip list change from ResourceData
func ipChangeParser(data *schema.ResourceData, key string) []string {
	_, ipChangeSet := data.GetChange(key)
	ipList := ipChangeSet.(*schema.Set).List()
	newIps := make([]string, len(ipList))
	for idx, value := range ipList {
		newIps[idx] = fmt.Sprintf("%v", value)
	}
	return newIps
}
