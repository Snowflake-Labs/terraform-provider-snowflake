package sdk

import "context"

type NetworkPolicies interface {
	Create(ctx context.Context, request *CreateNetworkPolicyRequest) error
	Show(ctx context.Context, request *ShowNetworkPolicyRequest) error
}

// CreateNetworkPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-network-policy.
type CreateNetworkPolicyOptions struct {
	create        bool                    `ddl:"static" sql:"CREATE"`
	OrReplace     *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	networkPolicy bool                    `ddl:"static" sql:"NETWORK POLICY"`
	name          AccountObjectIdentifier `ddl:"identifier"`
	AllowedIpList []string                `ddl:"parameter,parentheses" sql:"ALLOWED_IP_LIST"`
	Comment       *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// ShowNetworkPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-network-policies.
type ShowNetworkPolicyOptions struct {
	show            bool `ddl:"static" sql:"SHOW"`
	networkPolicies bool `ddl:"static" sql:"NETWORK POLICIES"`
}

// databaseNetworkPolicyDBRow is used to decode the result of a Show NetworkPolicies query.
type databaseNetworkPolicyDBRow struct {
	CreatedOn              string `db:"created_on"`
	Name                   string `db:"name"`
	Comment                string `db:"comment"`
	EntriesInAllowedIpList int    `db:"entries_in_allowed_ip_list"`
	EntriesInBlockedIpList int    `db:"entries_in_blocked_ip_list"`
}

// NetworkPolicy is used to decode the result of a Show NetworkPolicies query.
type NetworkPolicy struct {
	CreatedOn              string
	Name                   string
	Comment                string
	EntriesInAllowedIpList int
	EntriesInBlockedIpList int
}
