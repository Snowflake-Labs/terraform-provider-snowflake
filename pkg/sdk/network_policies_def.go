package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var (
	// TODO: prefix / postfix for top level definitions - readme
	ip = g.QueryStruct("IP").
		Text("IP", g.KeywordOptions().SingleQuotes().Required())

	NetworkPoliciesDef = g.NewInterface(
		"NetworkPolicies",
		"NetworkPolicy",
		g.KindOfT[AccountObjectIdentifier](),
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-network-policy",
			// Change
			// TODO top level QueryStruct name doesn't matter because it's created from op name + interface field
			g.QueryStruct("CreateNetworkPolicies").
				Create().
				OrReplace().
				SQL("NETWORK POLICY").
				SelfIdentifier().
				ListQueryStructField("AllowedIpList", ip, g.ParameterOptions().SQL("ALLOWED_IP_LIST").Parentheses()).
				ListQueryStructField("BlockedIpList", ip, g.ParameterOptions().SQL("BLOCKED_IP_LIST").Parentheses()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ValidIdentifier, "name"),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-network-policy",
			g.QueryStruct("AlterNetworkPolicy").
				Alter().
				SQL("NETWORK POLICY").
				IfExists().
				// TODO is it ok - add to readme
				SelfIdentifier().
				OptionalQueryStructField(
					// TODO We can omit name and derive it from type, in this case field could be NetworkPolicySet - yes, already in readme - add something
					// Or we can have a convention of <resource name><type> and remove prefix
					"Set",
					g.QueryStruct("NetworkPolicySet").
						ListQueryStructField("AllowedIpList", ip, g.ParameterOptions().SQL("ALLOWED_IP_LIST").Parentheses()).
						ListQueryStructField("BlockedIpList", ip, g.ParameterOptions().SQL("BLOCKED_IP_LIST").Parentheses()).
						OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
						WithValidation(g.AtLeastOneValueSet, "AllowedIpList", "BlockedIpList", "Comment"),
					g.KeywordOptions().SQL("SET"),
				).
				OptionalSQL("UNSET COMMENT").
				Identifier("RenameTo", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.AtLeastOneValueSet, "Set", "UnsetComment", "RenameTo").
				WithValidation(g.ValidIdentifierIfSet, "RenameTo"),
		).
		DropOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/drop-network-policy",
			g.QueryStruct("DropNetworkPolicy").
				Drop().
				SQL("NETWORK POLICY").
				IfExists().
				SelfIdentifier().
				WithValidation(g.ValidIdentifier, "name"),
		).
		ShowOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/show-network-policies",
			g.DbStruct("showNetworkPolicyDBRow").
				Field("created_on", "string").
				Field("name", "string").
				Field("comment", "string").
				Field("entries_in_allowed_ip_list", "int").
				Field("entries_in_blocked_ip_list", "int"),
			g.PlainStruct("NetworkPolicy").
				Field("CreatedOn", "string").
				Field("Name", "string").
				Field("Comment", "string").
				Field("EntriesInAllowedIpList", "int").
				Field("EntriesInBlockedIpList", "int"),
			g.QueryStruct("ShowNetworkPolicies").
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
			g.QueryStruct("DescribeNetworkPolicy").
				Describe().
				SQL("NETWORK POLICY").
				SelfIdentifier().
				WithValidation(g.ValidIdentifier, "name"),
		)
)
