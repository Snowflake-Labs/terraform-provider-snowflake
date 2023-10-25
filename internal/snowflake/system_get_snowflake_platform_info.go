// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func SystemGetSnowflakePlatformInfoQuery() string {
	return `SELECT SYSTEM$GET_SNOWFLAKE_PLATFORM_INFO() AS "INFO"`
}

type RawPlatformInfo struct {
	Info string `db:"INFO"`
}

type platformInfoInternal struct {
	AzureVnetSubnetIds []string `json:"snowflake-vnet-subnet-id,omitempty"`
	AwsVpcIds          []string `json:"snowflake-vpc-id,omitempty"`
}

type PlatformInfo struct {
	AzureVnetSubnetIds []string
	AwsVpcIds          []string
}

func ScanSnowflakePlatformInfo(row *sqlx.Row) (*RawPlatformInfo, error) {
	info := &RawPlatformInfo{}
	err := row.StructScan(info)
	return info, err
}

func (r *RawPlatformInfo) GetStructuredConfig() (*PlatformInfo, error) {
	info := &platformInfoInternal{}
	err := json.Unmarshal([]byte(r.Info), info)
	if err != nil {
		return nil, err
	}

	return info.getSnowflakePlatformInfo()
}

func (i *platformInfoInternal) getSnowflakePlatformInfo() (*PlatformInfo, error) {
	config := &PlatformInfo{
		i.AzureVnetSubnetIds,
		i.AwsVpcIds,
	}

	return config, nil
}
