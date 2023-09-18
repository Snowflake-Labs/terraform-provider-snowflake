package sdk

// TODO expose only needed types (field, interface, operation could be not exposed - only building functions)
import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var (
	networkPolicyRepresentation = g.PlainStruct("NetworkPolicy").
					Field("CreatedOn", "string").
					Field("Name", "string").
					Field("Comment", "string").
					Field("EntriesInAllowedIpList", "int").
					Field("EntriesInBlockedIpList", "int")

	ip = g.QueryStruct("IP").
		Text("IP", nil)

	NetworkPoliciesDef = g.NewInterface(
		"NetworkPolicies",
		"NetworkPolicy",
		g.KindOfT[AccountObjectIdentifier](),
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-network-policy",
			// top level QueryStruct name doesn't matter because it's created from op name + interface field
			g.QueryStruct("CreateNetworkPolicies").
				Create().
				OrReplace().
				SQL("NETWORK POLICY").
				// by convention field is named "name" and type is derived from interface field
				SelfIdentifier().
				ListQueryStructField("AllowedIpList", ip, g.ParameterOptions().SQL("ALLOWED_IP_LIST").Parentheses()).
				ListQueryStructField("BlockedIpList", ip, g.ParameterOptions().SQL("BLOCKED_IP_LIST").Parentheses()).
				// for those cases better to pass name and sql prefix in Options ?
				// ListAssignment is Optional in its nature, because it's a slice which can be null,
				// thus should it be ListAssignment or OptionalListAssignment ?
				//ListAssignment("ALLOWED_IP_LIST", "string", g.ParameterOptions().Parentheses().SingleQuotes().Required()).
				//ListAssignment("BLOCKED_IP_LIST", "string", g.ParameterOptions().Parentheses().SingleQuotes()).
				// for those cases better to pass name and sql prefix in Options ? prefix could be empty
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ValidIdentifier, "name"),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-network-policy",
			g.QueryStruct("AlterNetworkPolicy").
				Alter().
				SQL("NETWORK POLICY").
				IfExists().
				SelfIdentifier().
				OptionalQueryStructField(
					// We can omit name and derive it from type, in this case field could be NetworkPolicySet
					// Or we can have a convention of <resource name><type> and remove prefix
					"Set",
					g.QueryStruct("NetworkPolicySet").
						// should we pass plain kinds or instead there should be interface with [Kind() string] func in it
						// then we would force users to use g.KindOf... functions family, and it would look more consistent
						// with places where we would use g.KindOfT[type]()
						ListAssignment("ALLOWED_IP_LIST", "string", g.ParameterOptions().Parentheses().SingleQuotes()).
						ListAssignment("BLOCKED_IP_LIST", "string", g.ParameterOptions().Parentheses().SingleQuotes()).
						OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
						WithValidation(g.AtLeastOneValueSet, "AllowedIpList", "BlockedIpList", "Comment"),
					g.KeywordOptions().SQL("SET"),
				).
				OptionalSQL("UNSET COMMENT").
				Identifier("RenameTo", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
				// generator.ValidIdentifier validation can be implicit (we can add it when calling SelfIdentifier)
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
			networkPolicyRepresentation,
			g.QueryStruct("ShowNetworkPolicies").
				Show().
				SQL("NETWORK POLICIES"),
		).
		// Should describe be always Describe(context, id) or Describe(context, request) ? or we support both ?
		//	e.g. external tables are expecting more inputs than id
		DescribeOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/desc-network-policy",
			g.DbStruct("describeNetworkPolicyDBRow").
				Field("name", "string").
				Field("value", "string"),
			networkPolicyRepresentation,
			g.QueryStruct("DescribeNetworkPolicy").
				Describe().
				SQL("NETWORK POLICY").
				SelfIdentifier().
				WithValidation(g.ValidIdentifier, "name"),
		)
)
