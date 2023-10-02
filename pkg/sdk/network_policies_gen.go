package sdk

import "context"

type NetworkPolicies interface {
	Create(ctx context.Context, request *CreateNetworkPolicyRequest) error
	Alter(ctx context.Context, request *AlterNetworkPolicyRequest) error
	Drop(ctx context.Context, request *DropNetworkPolicyRequest) error
	Show(ctx context.Context, request *ShowNetworkPolicyRequest) ([]NetworkPolicy, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*NetworkPolicy, error)
	Describe(ctx context.Context, id AccountObjectIdentifier) ([]NetworkPolicyDescription, error)
}

// CreateNetworkPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-network-policy.
type CreateNetworkPolicyOptions struct {
	create        bool                    `ddl:"static" sql:"CREATE"`
	OrReplace     *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	networkPolicy bool                    `ddl:"static" sql:"NETWORK POLICY"`
	name          AccountObjectIdentifier `ddl:"identifier"`
	AllowedIpList []IP                    `ddl:"parameter,parentheses" sql:"ALLOWED_IP_LIST"`
	BlockedIpList []IP                    `ddl:"parameter,parentheses" sql:"BLOCKED_IP_LIST"`
	Comment       *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type IP struct {
	IP string `ddl:"keyword,single_quotes"`
}

// AlterNetworkPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-network-policy.
type AlterNetworkPolicyOptions struct {
	alter         bool                     `ddl:"static" sql:"ALTER"`
	networkPolicy bool                     `ddl:"static" sql:"NETWORK POLICY"`
	IfExists      *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name          AccountObjectIdentifier  `ddl:"identifier"`
	Set           *NetworkPolicySet        `ddl:"keyword" sql:"SET"`
	UnsetComment  *bool                    `ddl:"keyword" sql:"UNSET COMMENT"`
	RenameTo      *AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
}

type NetworkPolicySet struct {
	AllowedIpList []IP    `ddl:"parameter,parentheses" sql:"ALLOWED_IP_LIST"`
	BlockedIpList []IP    `ddl:"parameter,parentheses" sql:"BLOCKED_IP_LIST"`
	Comment       *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// DropNetworkPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-network-policy.
type DropNetworkPolicyOptions struct {
	drop          bool                    `ddl:"static" sql:"DROP"`
	networkPolicy bool                    `ddl:"static" sql:"NETWORK POLICY"`
	IfExists      *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name          AccountObjectIdentifier `ddl:"identifier"`
}

// ShowNetworkPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-network-policies.
type ShowNetworkPolicyOptions struct {
	show            bool `ddl:"static" sql:"SHOW"`
	networkPolicies bool `ddl:"static" sql:"NETWORK POLICIES"`
}

type showNetworkPolicyDBRow struct {
	CreatedOn              string `db:"created_on"`
	Name                   string `db:"name"`
	Comment                string `db:"comment"`
	EntriesInAllowedIpList int    `db:"entries_in_allowed_ip_list"`
	EntriesInBlockedIpList int    `db:"entries_in_blocked_ip_list"`
}

type NetworkPolicy struct {
	CreatedOn              string
	Name                   string
	Comment                string
	EntriesInAllowedIpList int
	EntriesInBlockedIpList int
}

// DescribeNetworkPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-network-policy.
type DescribeNetworkPolicyOptions struct {
	describe      bool                    `ddl:"static" sql:"DESCRIBE"`
	networkPolicy bool                    `ddl:"static" sql:"NETWORK POLICY"`
	name          AccountObjectIdentifier `ddl:"identifier"`
}

type describeNetworkPolicyDBRow struct {
	Name  string `db:"name"`
	Value string `db:"value"`
}

type NetworkPolicyDescription struct {
	Name  string
	Value string
}
