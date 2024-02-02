package datasources

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// ReadSystemGetSnowflakePlatformInfo implements schema.ReadFunc.
func ReadSystemGetSnowflakePlatformInfo(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	sel := snowflake.SystemGetSnowflakePlatformInfoQuery()
	row := snowflake.QueryRow(db, sel)

	acc, err := client.ContextFunctions.Current(context.Background())
	if err != nil {
		// If not found, mark resource to be removed from state file during apply or refresh
		d.SetId("")
		log.Println("[DEBUG] current_account failed to decode")
		return fmt.Errorf("error current_account err = %w", err)
	}

	d.SetId(fmt.Sprintf("%s.%s", acc.Account, acc.Region))

	rawInfo, err := snowflake.ScanSnowflakePlatformInfo(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Println("[DEBUG] system_get_snowflake_platform_info not found")
		return fmt.Errorf("error system_get_snowflake_platform_info err = %w", err)
	}

	info, err := rawInfo.GetStructuredConfig()
	if err != nil {
		log.Println("[DEBUG] system_get_snowflake_platform_info failed to decode")
		d.SetId("")
		return fmt.Errorf("error system_get_snowflake_platform_info err = %w", err)
	}

	if err := d.Set("azure_vnet_subnet_ids", info.AzureVnetSubnetIds); err != nil {
		return fmt.Errorf("error system_get_snowflake_platform_info err = %w", err)
	}

	if err := d.Set("aws_vpc_ids", info.AwsVpcIds); err != nil {
		return fmt.Errorf("error system_get_snowflake_platform_info err = %w", err)
	}

	return nil
}
