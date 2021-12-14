package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var systemGetSnowflakePlatformInfoSchema = map[string]*schema.Schema{
	"azure_vnet_subnet_ids": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Computed:    true,
		Description: "Snowflake Azure Virtual Network Subnet IDs",
	},
	"aws_vpc_ids": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Computed:    true,
		Description: "Snowflake AWS Virtual Private Cloud IDs",
	},
}

func SystemGetSnowflakePlatformInfo() *schema.Resource {
	return &schema.Resource{
		Read:   ReadSystemGetSnowflakePlatformInfo,
		Schema: systemGetSnowflakePlatformInfoSchema,
	}
}

// ReadSystemGetSnowflakePlatformInfo implements schema.ReadFunc
func ReadSystemGetSnowflakePlatformInfo(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	sel := snowflake.SystemGetSnowflakePlatformInfoQuery()
	row := snowflake.QueryRow(db, sel)

	acc, err := snowflake.ReadCurrentAccount(db)
	if err != nil {
		// If not found, mark resource to be removed from statefile during apply or refresh
		d.SetId("")
		log.Printf("[DEBUG] current_account failed to decode")
		return errors.Wrap(err, "error current_account")
	}

	d.SetId(fmt.Sprintf("%s.%s", acc.Account, acc.Region))

	rawInfo, err := snowflake.ScanSnowflakePlatformInfo(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Print("[DEBUG] system_get_snowflake_platform_info not found")
		return errors.Wrap(err, "error system_get_snowflake_platform_info")
	}

	info, err := rawInfo.GetStructuredConfig()
	if err != nil {
		log.Printf("[DEBUG] system_get_snowflake_platform_info failed to decode")
		d.SetId("")
		return errors.Wrap(err, "error system_get_snowflake_platform_info")
	}

	if err = d.Set("azure_vnet_subnet_ids", info.AzureVnetSubnetIds); err != nil {
		return errors.Wrap(err, "error system_get_snowflake_platform_info")
	}

	if err = d.Set("aws_vpc_ids", info.AwsVpcIds); err != nil {
		return errors.Wrap(err, "error system_get_snowflake_platform_info")
	}

	return nil
}
