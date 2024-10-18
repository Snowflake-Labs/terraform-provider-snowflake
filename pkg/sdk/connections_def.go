package sdk

//go:generate go run ./poc/main.go
import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

var enabledFailoverAccounts = g.NewQueryStruct("EnabledFailoverAccounts").
	Text("Account", g.KeywordOptions().NoQuotes())

var ConnectionDef = g.NewInterface(
	"Conntections",
	"Connection",
	g.KindOfT[AccountObjectIdentifier](),
).CustomOperation(
	"CreateConnection",
	"https://docs.snowflake.com/en/sql-reference/sql/create-connection",
	g.NewQueryStruct("CreateConnection").
		Create().
		SQL("CONNECTION").
		IfNotExists().
		Name().
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateReplicatedConnection",
	"https://docs.snowflake.com/en/sql-reference/sql/create-connection",
	g.NewQueryStruct("CreateReplicatedConnection").
		Create().
		SQL("CONNECTION").
		IfNotExists().
		Name().
		SQL("AS REPLICA OF").
		// external reference to connection: <orgnization_name>.<account_name>.<name>
		Identifier("ReplicaOf", g.KindOfT[ExternalObjectIdentifier](), g.IdentifierOptions().Required()).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "ReplicaOf"),
).CustomOperation(
	"AlterConnectionFailover",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-connection",
	g.NewQueryStruct("AlterConnectionFailover").
		Alter().
		SQL("CONNECTION").
		Name().
		OptionalQueryStructField(
			"EnableConnectionFailover",
			g.NewQueryStruct("EnableConnectionFailover").
				// ListQueryStructField("Accounts", enabledFailoverAccounts, g.ListOptions().NoParentheses()).
				List("Accounts", "ExternalObjectIdentifier", g.ListOptions().NoParentheses()).
				OptionalSQL("IGNORE EDITION CHECK"),
			g.KeywordOptions().SQL("ENABLE FAILOVER TO ACCOUNTS"),
		).
		OptionalQueryStructField(
			"DisableConnectionFailover",
			g.NewQueryStruct("DisableConnectionFailover").
				OptionalSQL("TO ACCOUNTS").
				//		ListQueryStructField("Accounts", enabledFailoverAccounts, g.ListOptions().NoParentheses()),
				List("Accounts", "ExternalObjectIdentifier", g.ListOptions().NoParentheses()),
			g.KeywordOptions().SQL("DISABLE FAILOVER"),
		).
		OptionalQueryStructField(
			"Primary",
			g.NewQueryStruct("Primary").
				SQL("PRIMARY"),
			g.KeywordOptions(),
		).
		WithValidation(g.ExactlyOneValueSet, "EnableConnectionFailover", "DisableConnectionFailover", "Primary"),
).CustomOperation(
	"AlterConnection",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-connection",
	g.NewQueryStruct("AlterConnection").
		Alter().
		SQL("CONNECTION").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("Set").
				OptionalComment().
				WithValidation(g.AtLeastOneValueSet, "Comment"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("Unset").
				OptionalSQL("COMMENT").
				WithValidation(g.AtLeastOneValueSet, "Comment"),
			g.KeywordOptions().SQL("UNSET"),
		).
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset"),
)
