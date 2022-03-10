package datasources

import (
	"database/sql"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var organizationAccountsSchema = map[string]*schema.Schema{
	"accounts": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The accounts in the organization",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"region_group": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The Region Group where the account is located; if it exists",
				},
				"region": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The Snowflake Region; as returned by CURRENT_REGION()",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The name of the account.",
				},
				"edition": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The Snowflake Edition of the account.",
				},
				"created_on": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The Date and Time the account was created.",
				},
				"account_url": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Snowflake access URL",
				},
				"comment": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Comment for the account",
				},
				"managed_accounts": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Snowflake access URL",
				},
				"account_locator_url": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Older Snowflake access URL syntax",
				},
			},
		},
	},
}

// OrganizationAccounts the Snowflake  account resource
func OrganizationAccounts() *schema.Resource {
	return &schema.Resource{
		Read:   ReadOrganizationAccounts,
		Schema: organizationAccountsSchema,
	}
}

// ReadCurrentAccount read the current snowflake account information
func ReadOrganizationAccounts(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	currentOrgAccounts, err := snowflake.ListOrganizationAccounts(db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Print("[DEBUG] no organization accounts found")
		d.SetId("")
		return nil
	} else if err != nil {
		log.Print("[DEBUG] unable to parse organization accounts")
		d.SetId("")
		return nil
	}

	orgAccounts := []map[string]interface{}{}

	for _, orgAccount := range currentOrgAccounts {
		orgAccountMap := map[string]interface{}{}

		orgAccountMap["region_group"] = orgAccount.RegionGroup.String
		orgAccountMap["region"] = orgAccount.SnowflakeRegion.String
		orgAccountMap["name"] = orgAccount.Name.String
		orgAccountMap["edition"] = orgAccount.Edition.String
		orgAccountMap["created_on"] = orgAccount.CreatedOn.String
		orgAccountMap["account_url"] = orgAccount.AccountUrl.String
		orgAccountMap["comment"] = orgAccount.Comment.String
		orgAccountMap["managed_accounts"] = orgAccount.ManagedAccounts.Int32
		orgAccountMap["account_locator_url"] = orgAccount.AccountLocatorUrl.String

		orgAccounts = append(orgAccounts, orgAccountMap)
	}

	return d.Set("accounts", orgAccounts)
}
