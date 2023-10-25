package datasources

import (
	"context"
	"database/sql"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountsSchema = map[string]*schema.Schema{
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies an account name pattern. If a pattern is specified, only accounts matching the pattern are returned.",
	},
	"accounts": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of all the accounts available in the organization.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"organization_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the organization.",
				},
				"account_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "User-defined name that identifies an account within the organization.",
				},
				"region_group": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Region group where the account is located. Note: this column is only visible to organizations that span multiple Region Groups.",
				},
				"snowflake_region": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Snowflake Region where the account is located. A Snowflake Region is a distinct location within a cloud platform region that is isolated from other Snowflake Regions. A Snowflake Region can be either multi-tenant or single-tenant (for a Virtual Private Snowflake account).",
				},
				"edition": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Snowflake Edition of the account.",
				},
				"account_url": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Preferred Snowflake access URL that includes the values of organization_name and account_name.",
				},
				"created_on": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Date and time when the account was created.",
				},
				"comment": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Comment for the account.",
				},
				"account_locator": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "System-assigned identifier of the acccount.",
				},
				"account_locator_url": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Legacy Snowflake access URL syntax that includes the region_name and account_locator.",
				},
				"managed_accounts": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Indicates how many managed accounts have been created by the account.",
				},
				"consumption_billing_entity_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the consumption billing entity.",
				},
				"marketplace_consumer_billing_entity_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the marketplace consumer billing entity.",
				},
				"marketplace_provider_billing_entity_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the marketplace provider billing entity.",
				},
				"old_account_url": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The previous account URL for a given account.",
				},
				"is_org_admin": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates whether the ORGADMIN role is enabled in an account. If TRUE, the role is enabled.",
				},
			},
		},
	},
}

// Accounts Snowflake Accounts resource.
func Accounts() *schema.Resource {
	return &schema.Resource{
		Read:   ReadAccounts,
		Schema: accountsSchema,
	}
}

// ReadAccounts lists accounts.
func ReadAccounts(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	ok, err := client.ContextFunctions.IsRoleInSession(ctx, sdk.NewAccountObjectIdentifier("ORGADMIN"))
	if err != nil {
		return err
	}
	if !ok {
		log.Printf("[DEBUG] ORGADMIN role is not in current session, cannot read accounts")
		return nil
	}
	opts := &sdk.ShowAccountOptions{}
	if pattern, ok := d.GetOk("pattern"); ok {
		opts.Like = &sdk.Like{
			Pattern: sdk.String(pattern.(string)),
		}
	}
	accounts, err := client.Accounts.Show(ctx, opts)
	if err != nil {
		return err
	}
	d.SetId("accounts")
	accountsFlatten := []map[string]interface{}{}
	for _, account := range accounts {
		m := map[string]interface{}{}
		m["organization_name"] = account.OrganizationName
		m["account_name"] = account.AccountName
		m["region_group"] = account.RegionGroup
		m["snowflake_region"] = account.SnowflakeRegion
		m["edition"] = string(account.Edition)
		m["account_url"] = account.AccountURL
		m["created_on"] = account.CreatedOn.String()
		m["comment"] = account.Comment
		m["account_locator"] = account.AccountLocator
		m["account_locator_url"] = account.AccountLocatorURL
		m["managed_accounts"] = account.ManagedAccounts
		m["consumption_billing_entity_name"] = account.ConsumptionBillingEntityName
		m["marketplace_consumer_billing_entity_name"] = account.MarketplaceConsumerBillingEntityName
		m["marketplace_provider_billing_entity_name"] = account.MarketplaceProviderBillingEntityName
		m["old_account_url"] = account.OldAccountURL
		m["is_org_admin"] = account.IsOrgAdmin
		accountsFlatten = append(accountsFlatten, m)
	}
	if err := d.Set("accounts", accountsFlatten); err != nil {
		return err
	}
	return nil
}
