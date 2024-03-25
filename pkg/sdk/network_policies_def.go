package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var (
	ip = g.NewQueryStruct("IP").
		Text("IP", g.KeywordOptions().SingleQuotes().Required())

	networkPoliciesAddNetworkRule = g.NewQueryStruct("AddNetworkRule").
					ListAssignment("ALLOWED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
					ListAssignment("BLOCKED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
					WithValidation(g.ExactlyOneValueSet, "AllowedNetworkRuleList", "BlockedNetworkRuleList")

	networkPoliciesRemoveNetworkRule = g.NewQueryStruct("RemoveNetworkRule").
						ListAssignment("ALLOWED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
						ListAssignment("BLOCKED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
						WithValidation(g.ExactlyOneValueSet, "AllowedNetworkRuleList", "BlockedNetworkRuleList")

	NetworkPoliciesDef = g.NewInterface(
		"NetworkPolicies",
		"NetworkPolicy",
		g.KindOfT[AccountObjectIdentifier](),
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-network-policy",
			g.NewQueryStruct("CreateNetworkPolicies").
				Create().
				OrReplace().
				SQL("NETWORK POLICY").
				Name().
				ListAssignment("ALLOWED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
				ListAssignment("BLOCKED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
				ListQueryStructField("AllowedIpList", ip, g.ParameterOptions().SQL("ALLOWED_IP_LIST").Parentheses()).
				ListQueryStructField("BlockedIpList", ip, g.ParameterOptions().SQL("BLOCKED_IP_LIST").Parentheses()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ValidIdentifier, "name"),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-network-policy",
			g.NewQueryStruct("AlterNetworkPolicy").
				Alter().
				SQL("NETWORK POLICY").
				IfExists().
				Name().
				OptionalQueryStructField(
					"Set",
					g.NewQueryStruct("NetworkPolicySet").
						ListAssignment("ALLOWED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
						ListAssignment("BLOCKED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
						ListQueryStructField("AllowedIpList", ip, g.ParameterOptions().SQL("ALLOWED_IP_LIST").Parentheses()).
						ListQueryStructField("BlockedIpList", ip, g.ParameterOptions().SQL("BLOCKED_IP_LIST").Parentheses()).
						OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
						WithValidation(g.AtLeastOneValueSet, "AllowedIpList", "BlockedIpList", "Comment", "AllowedNetworkRuleList", "BlockedNetworkRuleList"),
					g.KeywordOptions().SQL("SET"),
				).
				OptionalQueryStructField(
					"Add",
					networkPoliciesAddNetworkRule,
					g.KeywordOptions().SQL("ADD"),
				).
				OptionalQueryStructField(
					"Remove",
					networkPoliciesRemoveNetworkRule,
					g.KeywordOptions().SQL("REMOVE"),
				).
				OptionalSQL("UNSET COMMENT").
				Identifier("RenameTo", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ExactlyOneValueSet, "Set", "UnsetComment", "RenameTo", "Add", "Remove").
				WithValidation(g.ValidIdentifierIfSet, "RenameTo"),
		).
		DropOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/drop-network-policy",
			g.NewQueryStruct("DropNetworkPolicy").
				Drop().
				SQL("NETWORK POLICY").
				IfExists().
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		).
		ShowOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/show-network-policies",
			g.DbStruct("showNetworkPolicyDBRow").
				Field("created_on", "string").
				Field("name", "string").
				Field("comment", "string").
				Field("entries_in_allowed_ip_list", "int").
				Field("entries_in_blocked_ip_list", "int").
				Field("entries_in_allowed_network_rules", "int").
				Field("entries_in_blocked_network_rules", "int"),
			g.PlainStruct("NetworkPolicy").
				Field("CreatedOn", "string").
				Field("Name", "string").
				Field("Comment", "string").
				Field("EntriesInAllowedIpList", "int").
				Field("EntriesInBlockedIpList", "int").
				Field("EntriesInAllowedNetworkRules", "int").
				Field("EntriesInBlockedNetworkRules", "int"),
			g.NewQueryStruct("ShowNetworkPolicies").
				Show().
				SQL("NETWORK POLICIES"),
		).
		DescribeOperation(
			g.DescriptionMappingKindSlice,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-network-policy",
			g.DbStruct("describeNetworkPolicyDBRow").
				Field("name", "string").
				Field("value", "string"),
			g.PlainStruct("NetworkPolicyDescription").
				Field("Name", "string").
				Field("Value", "string"),
			g.NewQueryStruct("DescribeNetworkPolicy").
				Describe().
				SQL("NETWORK POLICY").
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		)
)
