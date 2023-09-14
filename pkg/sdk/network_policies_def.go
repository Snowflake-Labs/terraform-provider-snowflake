package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var (
	networkPolicyRepresentation = g.PlainStruct("NetworkPolicy").
					Field("CreatedOn", "string").
					Field("Name", "string").
					Field("Comment", "string").
					Field("EntriesInAllowedIpList", "int").
					Field("EntriesInBlockedIpList", "int")
	NetworkPoliciesDef = g.NewInterface(
		"NetworkPolicies",
		"NetworkPolicy",
		g.KindOfT[AccountObjectIdentifier](), // TODO Do we need this ?
		// We can use identifier kind above if we'll create fluent api for interface .CreateOperation().AlterOperation()...
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-network-policy",
			// TODO Could be type of Query struct that is converted into Field under the hood
			g.QueryStruct("CreateNetworkPolicies").
				Create().
				OrReplace().
				SQL("NETWORK POLICY").
				Identifier("name", g.KindOfT[AccountObjectIdentifier]()).
				// TODO list assignment
				ListAssignment("ALLOWED_IP_LIST", "string", g.ParameterOptions().Parentheses()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ValidIdentifier, "name"),
		).
		DropOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/drop-network-policy",
			g.QueryStruct("DropNetworkPolicy").
				Drop().
				SQL("NETWORK RULE").
				IfExists().
				Identifier("name", g.KindOfT[AccountObjectIdentifier]()),
		).
		ShowOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/show-network-policies",
			g.DbStruct("showNetworkPolicyDBRow").
				Field("created_on", "string").
				Field("name", "string").
				Field("comment", "string").
				Field("entries_in_allowed_ip_list", "int").
				Field("entries_in_blocked_ip_list", "int"),
			networkPolicyRepresentation,
			// TODO Mapping between DbStruct and PlainStruct
			g.QueryStruct("ShowNetworkPolicies").
				Show().
				SQL("NETWORK POLICIES"),
		).
		DescribeOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/desc-network-policy",
			g.DbStruct("describeNetworkPolicyDBRow").
				Field("name", "string").
				Field("value", "string"),
			// TODO because of that network policy will be generated twice
			networkPolicyRepresentation,
			g.QueryStruct("DescribeNetworkPolicy").
				SQL("DESCRIBE NETWORK POLICY").
				Identifier("name", g.KindOfT[AccountObjectIdentifier]()),
		)
)
