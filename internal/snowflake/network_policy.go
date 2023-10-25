// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// NetworkPolicyBuilder abstracts the creation of SQL queries for a Snowflake Network Policy.
type NetworkPolicyBuilder struct {
	name          string
	comment       string
	allowedIPList string
	blockedIPList string
}

// NetworkPolicy returns a pointer to a Builder that abstracts the DDL operations for a network policy.
//
// Supported DDL operations are:
//   - CREATE NETWORK POLICY
//   - DROP NETWORK POLICY
//   - SHOW NETWORK POLICIES
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/network-policies.html)
func NetworkPolicy(name string) *NetworkPolicyBuilder {
	return &NetworkPolicyBuilder{
		name: name,
	}
}

// SetOnAccount returns the SQL query that will set the network policy globally on your Snowflake account.
func (npb *NetworkPolicyBuilder) SetOnAccount() string {
	return fmt.Sprintf(`ALTER ACCOUNT SET NETWORK_POLICY = "%v"`, npb.name)
}

// UnsetOnAccount returns the SQL query that will unset the network policy globally on your Snowflake account.
func (npb *NetworkPolicyBuilder) UnsetOnAccount() string {
	return `ALTER ACCOUNT UNSET NETWORK_POLICY`
}

// SetOnUser returns the SQL query that will set the network policy on a given user.
func (npb *NetworkPolicyBuilder) SetOnUser(u string) string {
	return fmt.Sprintf(`ALTER USER "%v" SET NETWORK_POLICY = "%v"`, u, npb.name)
}

// UnsetOnUser returns the SQL query that will unset the network policy of a given user.
func (npb *NetworkPolicyBuilder) UnsetOnUser(u string) string {
	return fmt.Sprintf(`ALTER USER "%v" UNSET NETWORK_POLICY`, u)
}

// ShowOnUser returns the SQL query that will SHOW network policy set on a specific User.
func (npb *NetworkPolicyBuilder) ShowOnUser(u string) string {
	return fmt.Sprintf(`SHOW PARAMETERS LIKE 'network_policy' IN USER "%v"`, u)
}

// ShowOnAccount returns the SQL query that will SHOW network policy set on Account.
func (npb *NetworkPolicyBuilder) ShowOnAccount() string {
	return `SHOW PARAMETERS LIKE 'network_policy' IN ACCOUNT`
}

type NetworkPolicyAttachmentStruct struct {
	Key   sql.NullString `db:"key"`
	Value sql.NullString `db:"value"`
	Level sql.NullString `db:"level"`
}

func ScanNetworkPolicyAttachment(row *sqlx.Row) (*NetworkPolicyAttachmentStruct, error) {
	r := &NetworkPolicyAttachmentStruct{}
	err := row.StructScan(r)
	return r, err
}
