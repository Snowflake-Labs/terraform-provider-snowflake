package datasources

import (
	"context"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

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

// CurrentAccount the Snowflake current account resource.
func CurrentAccount() *schema.Resource {
	return &schema.Resource{
		Read:   ReadCurrentAccount,
		Schema: currentAccountSchema,
	}
}

// ReadCurrentAccount read the current snowflake account information.
func ReadCurrentAccount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	current, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	if err != nil {
		log.Println("[DEBUG] current_account failed to decode")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%s.%s", current.Account, current.Region))
	accountErr := d.Set("account", current.Account)
	if accountErr != nil {
		return accountErr
	}
	regionErr := d.Set("region", current.Region)
	if regionErr != nil {
		return regionErr
	}
	url, err := current.AccountURL()
	if err != nil {
		log.Println("[DEBUG] generating snowflake url failed")
		return nil
	}

	urlErr := d.Set("url", url)
	if urlErr != nil {
		return urlErr
	}
	return nil
}
