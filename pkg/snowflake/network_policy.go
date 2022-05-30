package snowflake

import (
	"database/sql"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/jmoiron/sqlx"
)

// NetworkPolicyBuilder abstracts the creation of SQL queries for a Snowflake Network Policy
type NetworkPolicyBuilder struct {
	name          string
	comment       string
	allowedIpList string
	blockedIpList string
}

// WithComment adds a comment to the NetworkPolicyBuilder
func (npb *NetworkPolicyBuilder) WithComment(c string) *NetworkPolicyBuilder {
	npb.comment = EscapeString(c)
	return npb
}

// WithAllowedIpList adds an allowedIpList to the NetworkPolicyBuilder
func (npb *NetworkPolicyBuilder) WithAllowedIpList(allowedIps []string) *NetworkPolicyBuilder {
	npb.allowedIpList = helpers.IpListToSnowflakeString(allowedIps)
	return npb
}

// WithBlockedIpList adds a blockedIpList to the NetworkPolicyBuilder
func (npb *NetworkPolicyBuilder) WithBlockedIpList(blockedIps []string) *NetworkPolicyBuilder {
	npb.blockedIpList = helpers.IpListToSnowflakeString(blockedIps)
	return npb
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

// Create returns the SQL query that will create a network policy.
func (npb *NetworkPolicyBuilder) Create() string {
	createSql := fmt.Sprintf(`CREATE NETWORK POLICY "%v" ALLOWED_IP_LIST=%v`, npb.name, npb.allowedIpList)
	if npb.blockedIpList != "" {
		createSql = createSql + fmt.Sprintf(" BLOCKED_IP_LIST=%v", npb.blockedIpList)
	}
	if npb.comment != "" {
		createSql = createSql + fmt.Sprintf(` COMMENT="%v"`, npb.comment)
	}

	return createSql
}

// Describe returns the SQL query that will describe a network policy
func (npb *NetworkPolicyBuilder) Describe() string {
	return fmt.Sprintf(`DESC NETWORK POLICY "%v"`, npb.name)
}

// ChangeComment returns the SQL query that will update the comment on the network policy.
func (npb *NetworkPolicyBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER NETWORK POLICY "%v" SET COMMENT = '%v'`, npb.name, EscapeString(c))
}

// RemoveComment returns the SQL query that will remove the comment on the network policy.
func (npb *NetworkPolicyBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER NETWORK POLICY "%v" UNSET COMMENT`, npb.name)
}

// ChangeIpList returns the SQL query that will update the ip list (of the specified listType) on the network policy.
func (npb *NetworkPolicyBuilder) ChangeIpList(listType string, ips []string) string {
	return fmt.Sprintf(`ALTER NETWORK POLICY "%v" SET %v_IP_LIST = %v`, npb.name, listType, helpers.IpListToSnowflakeString(ips))
}

// Drop returns the SQL query that will drop a network policy.
func (npb *NetworkPolicyBuilder) Drop() string {
	return fmt.Sprintf(`DROP NETWORK POLICY "%v"`, npb.name)
}

// SetOnAccount returns the SQL query that will set the network policy globally on your Snowflake account
func (npb *NetworkPolicyBuilder) SetOnAccount() string {
	return fmt.Sprintf(`ALTER ACCOUNT SET NETWORK_POLICY = "%v"`, npb.name)
}

// UnsetOnAccount returns the SQL query that will unset the network policy globally on your Snowflake account
func (npb *NetworkPolicyBuilder) UnsetOnAccount() string {
	return `ALTER ACCOUNT UNSET NETWORK_POLICY`
}

// SetOnUser returns the SQL query that will set the network policy on a given user
func (npb *NetworkPolicyBuilder) SetOnUser(u string) string {
	return fmt.Sprintf(`ALTER USER "%v" SET NETWORK_POLICY = "%v"`, u, npb.name)
}

// UnsetOnUser returns the SQL query that will unset the network policy of a given user
func (npb *NetworkPolicyBuilder) UnsetOnUser(u string) string {
	return fmt.Sprintf(`ALTER USER "%v" UNSET NETWORK_POLICY`, u)
}

// ShowAllNetworkPolicies returns the SQL query that will SHOW *all* network policies in the Snowflake account
// Snowflake's implementation of SHOW for network policies does *not* support limiting results with LIKE
func (npb *NetworkPolicyBuilder) ShowAllNetworkPolicies() string {
	return `SHOW NETWORK POLICIES`
}

//ShowOnUser returns the SQL query that will SHOW network policy set on a specific User
func (npb *NetworkPolicyBuilder) ShowOnUser(u string) string {
	return fmt.Sprintf(`SHOW PARAMETERS LIKE 'network_policy' IN USER "%v"`, u)
}

//ShowOnAccount returns the SQL query that will SHOW network policy set on Account
func (npb *NetworkPolicyBuilder) ShowOnAccount() string {
	return `SHOW PARAMETERS LIKE 'network_policy' IN ACCOUNT`
}

type NetworkPolicyStruct struct {
	CreatedOn              sql.NullString `db:"created_on"`
	Name                   sql.NullString `db:"name"`
	Comment                sql.NullString `db:"comment"`
	EntriesInAllowedIpList sql.NullString `db:"entries_in_allowed_ip_list"`
	EntriesInBlockedIpList sql.NullString `db:"entries_in_blocked_ip_list"`
}

type NetworkPolicyAttachmentStruct struct {
	Key   sql.NullString `db:"key"`
	Value sql.NullString `db:"value"`
	Level sql.NullString `db:"level"`
}

// ScanNetworkPolicies takes database rows and converts them to a list of NetworkPolicyStruct pointers
func ScanNetworkPolicies(rows *sqlx.Rows) ([]*NetworkPolicyStruct, error) {
	var n []*NetworkPolicyStruct

	for rows.Next() {
		r := &NetworkPolicyStruct{}
		err := rows.StructScan(r)
		if err != nil {
			return nil, err
		}
		n = append(n, r)
	}
	return n, nil
}

func ScanNetworkPolicyAttachment(row *sqlx.Row) (*NetworkPolicyAttachmentStruct, error) {
	r := &NetworkPolicyAttachmentStruct{}
	err := row.StructScan(r)
	return r, err
}
