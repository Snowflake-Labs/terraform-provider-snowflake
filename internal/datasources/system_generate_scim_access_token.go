// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package datasources

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var systemGenerateSCIMAccesstokenSchema = map[string]*schema.Schema{
	"integration_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "SCIM Integration Name",
	},
	"access_token": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "SCIM Access Token",
	},
}

func SystemGenerateSCIMAccessToken() *schema.Resource {
	return &schema.Resource{
		Read:   ReadSystemGenerateSCIMAccessToken,
		Schema: systemGenerateSCIMAccesstokenSchema,
	}
}

// ReadSystemGetAWSSNSIAMPolicy implements schema.ReadFunc.
func ReadSystemGenerateSCIMAccessToken(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	integrationName := d.Get("integration_name").(string)

	sel := snowflake.NewSystemGenerateSCIMAccessTokenBuilder(integrationName).Select()
	row := snowflake.QueryRow(db, sel)
	accessToken, err := snowflake.ScanSCIMAccessToken(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] system_generate_scim_access_token (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		log.Printf("[DEBUG] system_generate_scim_access_token (%s) failed to generate (%q)", d.Id(), err.Error())
		d.SetId("")
		return nil
	}

	d.SetId(integrationName)
	return d.Set("access_token", accessToken.Token)
}
