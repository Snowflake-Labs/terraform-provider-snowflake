// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
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
	db := meta.(*sql.DB)
	acc, err := snowflake.ReadCurrentAccount(db)
	if err != nil {
		log.Println("[DEBUG] current_account failed to decode")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%s.%s", acc.Account, acc.Region))
	accountErr := d.Set("account", acc.Account)
	if accountErr != nil {
		return accountErr
	}
	regionErr := d.Set("region", acc.Region)
	if regionErr != nil {
		return regionErr
	}
	url, err := acc.AccountURL()
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
