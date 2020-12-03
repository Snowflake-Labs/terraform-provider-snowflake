package resources

import (
	"database/sql"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var networkPolicyAttachmentSchema = map[string]*schema.Schema{
	"network_policy_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the network policy; must be unique for the account in which the network policy is created.",
		ForceNew:    true,
	},
	"set_for_account": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether the network policy should be applied globally to your Snowflake account<br><br>**Note:** The Snowflake user running `terraform apply` must be on an IP address allowed by the network policy to set that policy globally on the Snowflake account.<br><br>Additionally, a Snowflake account can only have one network policy set globally at any given time. This resource does not enforce one-policy-per-account, it is the user's responsibility to enforce this. If multiple network policy resources have `set_for_account: true`, the final policy set on the account will be non-deterministic.",
		Default:     false,
	},
	"users": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies which users the network policy should be attached to",
	},
}

// NetworkPolicyAttachment returns a pointer to the resource representing a network policy attachment
func NetworkPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: CreateNetworkPolicyAttachment,
		Read:   ReadNetworkPolicyAttachment,
		Update: UpdateNetworkPolicyAttachment,
		Delete: DeleteNetworkPolicyAttachment,

		Schema: networkPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateNetworkPolicyAttachment implements schema.CreateFunc
func CreateNetworkPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	policyName := d.Get("network_policy_name").(string)
	d.SetId(policyName + "_attachment")

	if d.Get("set_for_account").(bool) {
		err := setOnAccount(d, meta)
		if err != nil {
			return errors.Wrapf(err, "error creating attachment for network policy %v", policyName)
		}
	}

	if u, ok := d.GetOk("users"); ok {
		users := expandStringList(u.(*schema.Set).List())

		err := ensureUserAlterPrivileges(users, meta)
		if err != nil {
			return err
		}

		err = setOnUsers(users, d, meta)
		if err != nil {
			return errors.Wrapf(err, "error creating attachment for network policy %v", policyName)
		}
	}

	return nil
}

// ReadNetworkPolicyAttachment implements schema.ReadFunc
func ReadNetworkPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	// HACK: InternalValidate requires Read to be implemented
	// There is no way of using SHOW/DESC on Network Policies/Users to pull attachment information, so we can't actually Read
	return nil
}

// UpdateNetworkPolicyAttachment implements schema.UpdateFunc
func UpdateNetworkPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("set_for_account") {
		oldAcctFlag, newAcctFlag := d.GetChange("set_for_account")
		if newAcctFlag.(bool) {
			if err := setOnAccount(d, meta); err != nil {
				return err
			}
		} else if !newAcctFlag.(bool) && oldAcctFlag == true {
			if err := unsetOnAccount(d, meta); err != nil {
				return err
			}
		}
	}

	if d.HasChange("users") {
		old, new := d.GetChange("users")
		oldUsersSet := old.(*schema.Set)
		newUsersSet := new.(*schema.Set)

		removedUsers := expandStringList(oldUsersSet.Difference(newUsersSet).List())
		addedUsers := expandStringList(newUsersSet.Difference(oldUsersSet).List())

		err := ensureUserAlterPrivileges(removedUsers, meta)
		if err != nil {
			return err
		}

		err = ensureUserAlterPrivileges(addedUsers, meta)
		if err != nil {
			return err
		}

		for _, user := range removedUsers {
			err := unsetOnUser(user, d, meta)
			if err != nil {
				return err
			}
		}

		for _, user := range addedUsers {
			err := setOnUser(user, d, meta)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteNetworkPolicyAttachment implements schema.DeleteFunc
func DeleteNetworkPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	policyName := d.Get("network_policy_name").(string)
	d.SetId(policyName + "_attachment")

	err := unsetOnAccount(d, meta)
	if err != nil {
		return errors.Wrapf(err, "error deleting attachment for network policy %v", policyName)
	}

	if u, ok := d.GetOk("users"); ok {
		users := expandStringList(u.(*schema.Set).List())

		err := ensureUserAlterPrivileges(users, meta)
		if err != nil {
			return err
		}

		err = unsetOnUsers(users, d, meta)
		if err != nil {
			return errors.Wrapf(err, "error deleting attachment for network policy %v", policyName)
		}
	}

	return nil
}

// setOnAccount sets the network policy globally for the Snowflake account
// Note: the ip address of the session executing this SQL must be allowed by the network policy being set
func setOnAccount(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	policyName := d.Get("network_policy_name").(string)

	acctSql := snowflake.NetworkPolicy(policyName).SetOnAccount()

	err := snowflake.Exec(db, acctSql)
	if err != nil {
		return errors.Wrapf(err, "error setting network policy %v on account", policyName)
	}

	return nil
}

// setOnAccount unsets the network policy globally for the Snowflake account
func unsetOnAccount(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	policyName := d.Get("network_policy_name").(string)

	acctSql := snowflake.NetworkPolicy(policyName).UnsetOnAccount()

	err := snowflake.Exec(db, acctSql)
	if err != nil {
		return errors.Wrapf(err, "error unsetting network policy %v on account", policyName)
	}

	return nil
}

// setOnUsers sets the network policy for list of users
func setOnUsers(users []string, data *schema.ResourceData, meta interface{}) error {
	policyName := data.Get("network_policy_name").(string)
	for _, user := range users {
		err := setOnUser(user, data, meta)
		if err != nil {
			return errors.Wrapf(err, "error setting network policy %v on user %v", policyName, user)
		}
	}

	return nil
}

// setOnUser sets the network policy for a given user
func setOnUser(user string, data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	policyName := data.Get("network_policy_name").(string)
	userSql := snowflake.NetworkPolicy(policyName).SetOnUser(user)
	err := snowflake.Exec(db, userSql)
	if err != nil {
		return errors.Wrapf(err, "error setting network policy %v on user %v", policyName, user)
	}

	return nil
}

// unsetOnUsers unsets the network policy for list of users
func unsetOnUsers(users []string, data *schema.ResourceData, meta interface{}) error {
	policyName := data.Get("network_policy_name").(string)
	for _, user := range users {
		err := unsetOnUser(user, data, meta)
		if err != nil {
			return errors.Wrapf(err, "error unsetting network policy %v on user %v", policyName, user)
		}
	}

	return nil
}

// unsetOnUser sets the network policy for a given user
func unsetOnUser(user string, data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	policyName := data.Get("network_policy_name").(string)
	userSql := snowflake.NetworkPolicy(policyName).UnsetOnUser(user)
	err := snowflake.Exec(db, userSql)
	if err != nil {
		return errors.Wrapf(err, "error unsetting network policy %v on user %v", policyName, user)
	}

	return nil
}

// ensureUserAlterPrivileges ensures the executing Snowflake user can alter each user in the set of users
func ensureUserAlterPrivileges(users []string, meta interface{}) error {
	db := meta.(*sql.DB)
	for _, user := range users {
		userDescSql := snowflake.User(user).Describe()
		err := snowflake.Exec(db, userDescSql)
		if err != nil {
			return errors.Wrapf(err, "error altering network policy of user %v", user)
		}
	}

	return nil
}
