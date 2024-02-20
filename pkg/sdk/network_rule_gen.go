package sdk

import (
	"context"
	"time"
)

type NetworkRules interface {
	Create(ctx context.Context, request *CreateNetworkRuleRequest) error
	Alter(ctx context.Context, request *AlterNetworkRuleRequest) error
	Drop(ctx context.Context, request *DropNetworkRuleRequest) error
	Show(ctx context.Context, request *ShowNetworkRuleRequest) ([]NetworkRule, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*NetworkRule, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*NetworkRuleDetails, error)
}

// CreateNetworkRuleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-network-rule.
type CreateNetworkRuleOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	networkRule bool                   `ddl:"static" sql:"NETWORK RULE"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	Type        NetworkRuleType        `ddl:"parameter,no_quotes" sql:"TYPE"`
	ValueList   []NetworkRuleValue     `ddl:"parameter,parentheses" sql:"VALUE_LIST"`
	Mode        NetworkRuleMode        `ddl:"parameter,no_quotes" sql:"MODE"`
	Comment     *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type NetworkRuleValue struct {
	Value string `ddl:"keyword,single_quotes"`
}

// AlterNetworkRuleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-network-rule.
type AlterNetworkRuleOptions struct {
	alter       bool                   `ddl:"static" sql:"ALTER"`
	networkRule bool                   `ddl:"static" sql:"NETWORK RULE"`
	IfExists    *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	Set         *NetworkRuleSet        `ddl:"list" sql:"SET"`
	Unset       *NetworkRuleUnset      `ddl:"list" sql:"UNSET"`
}

type NetworkRuleSet struct {
	ValueList []NetworkRuleValue `ddl:"parameter,parentheses" sql:"VALUE_LIST"`
	Comment   *string            `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type NetworkRuleUnset struct {
	ValueList *bool `ddl:"keyword" sql:"VALUE_LIST"`
	Comment   *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropNetworkRuleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-network-rule.
type DropNetworkRuleOptions struct {
	drop        bool                   `ddl:"static" sql:"DROP"`
	networkRule bool                   `ddl:"static" sql:"NETWORK RULE"`
	IfExists    *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowNetworkRuleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-network-rules.
type ShowNetworkRuleOptions struct {
	show         bool       `ddl:"static" sql:"SHOW"`
	networkRules bool       `ddl:"static" sql:"NETWORK RULES"`
	Like         *Like      `ddl:"keyword" sql:"LIKE"`
	In           *In        `ddl:"keyword" sql:"IN"`
	StartsWith   *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit        *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type ShowNetworkRulesRow struct {
	CreatedOn          time.Time `db:"created_on"`
	Name               string    `db:"name"`
	DatabaseName       string    `db:"database_name"`
	SchemaName         string    `db:"schema_name"`
	Owner              string    `db:"owner"`
	Comment            string    `db:"comment"`
	Type               string    `db:"type"`
	Mode               string    `db:"mode"`
	EntriesInValueList int       `db:"entries_in_valuelist"`
	OwnerRoleType      string    `db:"owner_role_type"`
}

type NetworkRule struct {
	CreatedOn          time.Time
	Name               string
	DatabaseName       string
	SchemaName         string
	Owner              string
	Comment            string
	Type               NetworkRuleType
	Mode               NetworkRuleMode
	EntriesInValueList int
	OwnerRoleType      string
}

// DescribeNetworkRuleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-network-rule.
type DescribeNetworkRuleOptions struct {
	describe    bool                   `ddl:"static" sql:"DESCRIBE"`
	networkRule bool                   `ddl:"static" sql:"NETWORK RULE"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
}

type DescNetworkRulesRow struct {
	CreatedOn    time.Time `db:"created_on"`
	Name         string    `db:"name"`
	DatabaseName string    `db:"database_name"`
	SchemaName   string    `db:"schema_name"`
	Owner        string    `db:"owner"`
	Comment      string    `db:"comment"`
	Type         string    `db:"type"`
	Mode         string    `db:"mode"`
	ValueList    string    `db:"value_list"`
}

type NetworkRuleDetails struct {
	CreatedOn    time.Time
	Name         string
	DatabaseName string
	SchemaName   string
	Owner        string
	Comment      string
	Type         NetworkRuleType
	Mode         NetworkRuleMode
	ValueList    []string
}
