package snowflake

import (
	"database/sql"
	"fmt"
	"strings"

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
	npb.allowedIpList = IpListToString(allowedIps)
	return npb
}

// WithBlockedIpList adds a blockedIpList to the NetworkPolicyBuilder
func (npb *NetworkPolicyBuilder) WithBlockedIpList(blockedIps []string) *NetworkPolicyBuilder {
	npb.blockedIpList = IpListToString(blockedIps)
	return npb
}

// NetworkPolicy returns a pointer to a Builder that abstracts the DDL operations for a network policy.
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/sql/alter-network-policy.html)
func NetworkPolicy(name string) *NetworkPolicyBuilder {
	return &NetworkPolicyBuilder{
		name: name,
	}
}

// Create returns the SQL query that will create a network policy.
func (npb *NetworkPolicyBuilder) Create() string {
	return fmt.Sprintf(`CREATE NETWORK POLICY "%v" ALLOWED_IP_LIST=%v BLOCKED_IP_LIST=%v COMMENT="%v"`, npb.name, npb.allowedIpList, npb.blockedIpList, npb.comment)
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
	return fmt.Sprintf(`ALTER NETWORK POLICY "%v" SET %v_IP_LIST = %v`, npb.name, listType, IpListToString(ips))
}

// Drop returns the SQL query that will drop a network policy.
func (npb *NetworkPolicyBuilder) Drop() string {
	return fmt.Sprintf(`DROP NETWORK POLICY "%v"`, npb.name)
}

// Show returns the SQL query that will show the network policy
// Note: this requires a non-standard approach because Snowflake's implementation of Network Policies makes it
//   difficult to read a given policy's details
func (npb *NetworkPolicyBuilder) Show(meta interface{}) (string, error) {
	db := meta.(*sql.DB)

	// Run DESC + fetch Query ID because we'll later post-process the output separately via RESULT_SCAN
	descSql := fmt.Sprintf(`DESC NETWORK POLICY "%v"`, npb.name)
	descQueryId, err := ExecAndGetId(db, descSql)
	if err != nil {
		return "", err
	}

	// Run SHOW + fetch Query ID because we'll later post-process the output separately via RESULT_SCAN
	// Snowflake SHOW only supports showing *all* network policies, so we have to filter with SQL later
	showAllQueryID, err := ExecAndGetId(db, "SHOW NETWORK POLICIES")
	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf(`
	WITH
	desc_output AS (SELECT * FROM TABLE(RESULT_SCAN('%v'))),
    show_output AS (SELECT * FROM TABLE(RESULT_SCAN('%v')) WHERE "name" = '%v'),
	allowed_ips AS (SELECT "value" AS "allowed_ip_list" FROM desc_output WHERE "name" = 'ALLOWED_IP_LIST'),
	blocked_ips AS (SELECT "value" AS "blocked_ip_list" FROM desc_output WHERE "name" = 'BLOCKED_IP_LIST')
	SELECT *
      FROM show_output
      LEFT JOIN allowed_ips ON TRUE
      LEFT JOIN blocked_ips ON TRUE
	`, descQueryId, showAllQueryID, npb.name)

	return sql, nil
}

// SetOnAccount returns the SQL query that will set the network policy globally on your Snowflake account
func (npb *NetworkPolicyBuilder) SetOnAccount() string {
	return fmt.Sprintf(`ALTER ACCOUNT SET NETWORK_POLICY = "%v"`, npb.name)
}

// UnsetOnAccount returns the SQL query that will unset the network policy globally on your Snowflake account
func (npb *NetworkPolicyBuilder) UnsetOnAccount() string {
	return fmt.Sprintf(`ALTER ACCOUNT UNSET NETWORK_POLICY`)
}

// SetOnUser returns the SQL query that will set the network policy on a given user
func (npb *NetworkPolicyBuilder) SetOnUser(u string) string {
	return fmt.Sprintf(`ALTER USER "%v" SET NETWORK_POLICY = "%v"`, u, npb.name)
}

// UnsetOnUser returns the SQL query that will unset the network policy of a given user
func (npb *NetworkPolicyBuilder) UnsetOnUser(u string) string {
	return fmt.Sprintf(`ALTER USER "%v" UNSET NETWORK_POLICY`, u)
}

// IpListToString formats a list of IPs into a Snowflake-DDL friendly string, e.g. ('192.168.1.0', '192.168.1.100')
func IpListToString(ips []string) string {
	for index, element := range ips {
		ips[index] = fmt.Sprintf(`'%v'`, element)
	}

	return fmt.Sprintf("(%v)", strings.Join(ips, ", "))
}

type networkPolicy struct {
	CreatedOn              sql.NullString `db:"created_on"`
	Name                   sql.NullString `db:"name"`
	Comment                sql.NullString `db:"comment"`
	EntriesInAllowedIpList sql.NullString `db:"entries_in_allowed_ip_list"`
	EntriesInBlockedIpList sql.NullString `db:"entries_in_blocked_ip_list"`
	AllowedIpList          sql.NullString `db:"allowed_ip_list"`
	BlockedIpList          sql.NullString `db:"blocked_ip_list"`
}

func ScanNetworkPolicy(row *sqlx.Row) (*networkPolicy, error) {
	r := &networkPolicy{}
	err := row.StructScan(r)
	return r, err
}
