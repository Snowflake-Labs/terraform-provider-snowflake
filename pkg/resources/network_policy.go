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
	"set_for_account": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether this network policy should be applied globally to your Snowflake account<br><br>**Note:** The Snowflake user running `terraform apply` must be on an IP address allowed by the network policy to set that policy globally on the Snowflake account.<br><br>Additionally, a Snowflake account can only have one network policy set globally at any given time. This resource does not enforce one-policy-per-account, it is the user's responsibility to enforce this. If multiple network policy resources have `set_for_account: true`, the final policy set on the account will be non-deterministic.",
		Default:     false,
		ForceNew:    true,
	},
	"users": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies which users this network policy should be applied to.<br><br>**Note**: The network policy resource creates exclusive assignments of policies to users. Any users that have a network policy assigned outside of that policy's Terraform resource definition will have the policy unset.<br><br>For safety, you should *only* ever assign/unset network policies to users via Terraform.<br><br>Additionally, a Snowflake user can only have one network policy applied to it. This resource does not enforce one-policy-per-user, it is the user's responsibility to enforce this. If a user appears in multiple network policy resource definitions, the final policy set on the user will be non-deterministic.",
		ForceNew:    true,
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

	err := ensureUserAlterPrivileges(data, meta)
	if err != nil {
		return err
	}

	stmt := builder.Create()
	err = snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating network policy %v", name)
	}
	data.SetId(name)

	if data.Get("set_for_account").(bool) {
		err := setOnAccount(data, meta)
		if err != nil {
			return errors.Wrapf(err, "error creating network policy %v", name)
		}
	}

	err = setOnUsers(data, meta)
	if err != nil {
		return err
	}

	return ReadNetworkPolicy(data, meta)
}

// ReadNetworkPolicy implements schema.ReadFunc
func ReadNetworkPolicy(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Id()

	builder := snowflake.NetworkPolicy(name)
	showSql, err := builder.Show(meta)
	if err != nil {
		return err
	}

	row := snowflake.QueryRow(db, showSql)

	s, err := snowflake.ScanNetworkPolicy(row)
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

	err = data.Set("allowed_ip_list", strings.Split(s.AllowedIpList.String, ","))
	if err != nil {
		return err
	}

	err = data.Set("blocked_ip_list", strings.Split(s.BlockedIpList.String, ","))
	if err != nil {
		return err
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

	err := ensureUserAlterPrivileges(data, meta)
	if err != nil {
		return err
	}

	dropSql := snowflake.NetworkPolicy(name).Drop()
	err = snowflake.Exec(db, dropSql)
	if err != nil {
		return errors.Wrapf(err, "error deleting network policy %v", name)
	}

	data.SetId("")

	if data.Get("set_for_account").(bool) {
		err := unsetOnAccount(data, meta)
		if err != nil {
			return errors.Wrapf(err, "error deleting network policy %v", name)
		}
	}

	err = unsetOnUsers(data, meta)
	if err != nil {
		return err
	}

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

// setOnAccount sets the network policy globally for the Snowflake account
// Note: the ip address of the session executing this SQL must be allowed by the network policy being set
func setOnAccount(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Id()

	acctSql := snowflake.NetworkPolicy(name).SetOnAccount()

	err := snowflake.Exec(db, acctSql)
	if err != nil {
		return errors.Wrapf(err, "error setting network policy %v on account", name)
	}

	return nil
}

// setOnAccount unsets the network policy globally for the Snowflake account
func unsetOnAccount(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Id()

	acctSql := snowflake.NetworkPolicy(name).UnsetOnAccount()

	err := snowflake.Exec(db, acctSql)
	if err != nil {
		return errors.Wrapf(err, "error unsetting network policy %v on account", name)
	}

	return nil
}

// setOnUsers sets the network policy for all users currently specified in the network policy resource
func setOnUsers(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	if u, ok := data.GetOk("users"); ok {
		users := expandStringList(u.(*schema.Set).List())
		for _, user := range users {
			userSql := snowflake.NetworkPolicy(name).SetOnUser(user)
			err := snowflake.Exec(db, userSql)
			if err != nil {
				return errors.Wrapf(err, "error setting network policy %v on user %v", name, user)
			}
		}
	}

	return nil
}

// unsetOnUsers unsets the network policy for all users currently specified in the network policy resource
func unsetOnUsers(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	if u, ok := data.GetOk("users"); ok {
		users := expandStringList(u.(*schema.Set).List())
		for _, user := range users {
			userSql := snowflake.NetworkPolicy(name).UnsetOnUser(user)
			err := snowflake.Exec(db, userSql)
			if err != nil {
				return errors.Wrapf(err, "error unsetting network policy %v on user %v", name, user)
			}
		}
	}

	return nil
}

// ensureUserAlterPrivileges ensures the executing Snowflake user can alter all users in the network policy resource definition
func ensureUserAlterPrivileges(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	if u, ok := data.GetOk("users"); ok {
		users := expandStringList(u.(*schema.Set).List())
		for _, user := range users {
			userDescSql := snowflake.User(user).Describe()
			err := snowflake.Exec(db, userDescSql)
			if err != nil {
				return errors.Wrapf(err, "error altering network policy of user %v", user)
			}
		}
	}

	return nil
}
