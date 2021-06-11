package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var currentAccountSchema = map[string]*schema.Schema{
	"account": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake Account ID; as returned by CURRENT_ACCOUNT().",
	},

	"region": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake Region; as returned by CURRENT_REGION()",
	},

	"url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake URL.",
	},
}

// CurrentAccount the Snowflake current account resource
func CurrentAccount() *schema.Resource {
	return &schema.Resource{
		Read:   ReadCurrentAccount,
		Schema: currentAccountSchema,
	}
}

// ReadCurrentAccount read the current snowflake account information
func ReadCurrentAccount(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	acc, err := snowflake.ReadCurrentAccount(db)

	if err != nil {
		log.Printf("[DEBUG] current_account failed to decode")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%s.%s", acc.Account, acc.Region))
	d.Set("account", acc.Account)
	d.Set("region", acc.Region)
	url, err := acc.AccountURL()

	if err != nil {
		log.Printf("[DEBUG] generating snowflake url failed")
		d.SetId("")
		return nil
	}

	d.Set("url", url)
	return nil
}
