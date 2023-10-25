// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import (
	"context"
	"time"
)

var _ validatable = new(ShowRegionsOptions)

type ReplicationFunctions interface {
	ShowReplicationAccounts(ctx context.Context) ([]*ReplicationAccount, error)
	// todo: ShowReplicationDatabases(ctx context.Context, opts *ShowReplicationDatabasesOptions) ([]*ReplicationDatabase, error)
	ShowRegions(ctx context.Context, opts *ShowRegionsOptions) ([]*Region, error)
}

var _ ReplicationFunctions = (*replicationFunctions)(nil)

type replicationFunctions struct {
	client *Client
}

type ReplicationAccount struct {
	SnowflakeRegion  string    `db:"snowflake_region"`
	CreatedOn        time.Time `db:"created_on"`
	AccountName      string    `db:"account_name"`
	AccountLocator   string    `db:"account_locator"`
	Comment          string    `db:"comment"`
	OrganizationName string    `db:"organization_name"`
	IsOrgAdmin       bool      `db:"is_org_admin"`
}

func (v *ReplicationAccount) ID() AccountIdentifier {
	return AccountIdentifier{
		organizationName: v.OrganizationName,
		accountName:      v.AccountName,
		accountLocator:   v.AccountLocator,
	}
}

// ShowReplicationAccounts is based on https://docs.snowflake.com/en/sql-reference/sql/show-replication-accounts.
func (c *replicationFunctions) ShowReplicationAccounts(ctx context.Context) ([]*ReplicationAccount, error) {
	rows := []ReplicationAccount{}
	sql := "SHOW REPLICATION ACCOUNTS"
	err := c.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	replicationAccounts := make([]*ReplicationAccount, len(rows))
	for i, row := range rows {
		replicationAccount := row
		replicationAccounts[i] = &replicationAccount
	}
	return replicationAccounts, nil
}

type CloudType string

const (
	CloudTypeAWS   CloudType = "aws"
	CloudTypeAzure CloudType = "azure"
	CloudTypeGCP   CloudType = "gcp"
)

type Region struct {
	RegionGroup     string
	SnowflakeRegion string
	CloudType       CloudType
	Region          string
	DisplayName     string
}

type regionRow struct {
	RegionGroup     string `db:"region_group"`
	SnowflakeRegion string `db:"snowflake_region"`
	Cloud           string `db:"cloud"`
	Region          string `db:"region"`
	DisplayName     string `db:"display_name"`
}

func (row *regionRow) toRegion() *Region {
	return &Region{
		RegionGroup:     row.RegionGroup,
		SnowflakeRegion: row.SnowflakeRegion,
		CloudType:       CloudType(row.Cloud),
		Region:          row.Region,
		DisplayName:     row.DisplayName,
	}
}

type ShowRegionsOptions struct {
	show    bool  `ddl:"static" sql:"SHOW"`
	regions bool  `ddl:"static" sql:"REGIONS"`
	Like    *Like `ddl:"keyword" sql:"LIKE"`
}

func (opts *ShowRegionsOptions) validate() error {
	return nil
}

// ShowRegions is based on https://docs.snowflake.com/en/sql-reference/sql/show-regions.
func (c *replicationFunctions) ShowRegions(ctx context.Context, opts *ShowRegionsOptions) ([]*Region, error) {
	if opts == nil {
		opts = &ShowRegionsOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	rows := []regionRow{}
	err = c.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	regions := make([]*Region, len(rows))
	for i, row := range rows {
		regions[i] = row.toRegion()
	}
	return regions, nil
}
