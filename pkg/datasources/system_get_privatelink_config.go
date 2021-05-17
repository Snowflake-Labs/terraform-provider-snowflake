package datasources

import (
	"database/sql"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var systemGetPrivateLinkConfigSchema = map[string]*schema.Schema{
	"account_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The name of your Snowflake account.",
	},

	"account_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The URL used to connect to Snowflake through AWS PrivateLink or Azure Private Link.",
	},

	"oscp_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The OCSP URL corresponding to your Snowflake account that uses AWS PrivateLink or Azure Private Link.",
	},

	"aws_vpce_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The AWS VPCE ID for your account.",
	},

	"azure_pls_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Azure Private Link Service ID for your account.",
	},
}

func SystemGetPrivateLinkConfig() *schema.Resource {
	return &schema.Resource{
		Read:   ReadSystemGetPrivateLinkConfig,
		Schema: systemGetPrivateLinkConfigSchema,
	}
}

// ReadSystemGetPrivateLinkConfig implements schema.ReadFunc
func ReadSystemGetPrivateLinkConfig(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	sel := snowflake.SystemGetPrivateLinkConfigQuery()
	row := snowflake.QueryRow(db, sel)
	rawConfig, err := snowflake.ScanPrivateLinkConfig(row)

	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Print("[DEBUG] system_get_privatelink_config not found")
		d.SetId("")
		return nil
	}

	config, err := rawConfig.GetStructuredConfig()
	if err != nil {
		log.Printf("[DEBUG] system_get_privatelink_config failed to decode")
		d.SetId("")
		return nil
	}

	d.SetId(config.AccountName)
	d.Set("account_name", config.AccountName)
	d.Set("account_url", config.AccountURL)
	d.Set("oscp_url", config.OSCPURL)

	if config.AwsVpceID != "" {
		d.Set("aws_vpce_id", config.AwsVpceID)
	}

	if config.AzurePrivateLinkServiceID != "" {
		d.Set("azure_pls_id", config.AzurePrivateLinkServiceID)
	}

	return nil
}
