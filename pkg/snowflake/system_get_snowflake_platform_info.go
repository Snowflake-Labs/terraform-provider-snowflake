package snowflake

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func SystemGetSnowflakePlatformInfoQuery() string {
	return `SELECT SYSTEM$GET_SNOWFLAKE_PLATFORM_INFO() AS "info"`
}

type RawSnowflakePlatformInfo struct {
	Info string `db:"info"`
}

type snowflakePlatformInfoInternal struct {
	AzureVnetSubnetIds []string `json:"snowflake-vnet-subnet-id,omitempty"`
	AwsVpcIds          []string `json:"snowflake-vpc-id,omitempty"`
}

type SnowflakePlatformInfo struct {
	AzureVnetSubnetIds []string
	AwsVpcIds          []string
}

func ScanSnowflakePlatformInfo(row *sqlx.Row) (*RawSnowflakePlatformInfo, error) {
	info := &RawSnowflakePlatformInfo{}
	err := row.StructScan(info)
	return info, err
}

func (r *RawSnowflakePlatformInfo) GetStructuredConfig() (*SnowflakePlatformInfo, error) {
	info := &snowflakePlatformInfoInternal{}
	err := json.Unmarshal([]byte(r.Info), info)
	if err != nil {
		return nil, err
	}

	return info.getSnowflakePlatformInfo()
}

func (i *snowflakePlatformInfoInternal) getSnowflakePlatformInfo() (*SnowflakePlatformInfo, error) {
	config := &SnowflakePlatformInfo{
		i.AzureVnetSubnetIds,
		i.AwsVpcIds,
	}

	return config, nil
}
