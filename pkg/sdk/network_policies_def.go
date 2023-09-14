package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var NetworkPoliciesDef = g.NewInterface(
	"NetworkPolicies",
	"NetworkPolicy",
	g.KindOfT[AccountObjectIdentifier](), // TODO Do we need this ?
	g.CreateOperation(
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
	),
	g.ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-network-policies",
		g.DbStruct("databaseNetworkPolicyDBRow").
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
		// TODO Mapping between DbStruct and PlainStruct
		g.QueryStruct("ShowNetworkPolicies").
			Show().
			SQL("NETWORK POLICIES"),
	),
)
