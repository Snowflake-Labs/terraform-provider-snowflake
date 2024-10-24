package sdk

//go:generate go run ./poc/main.go

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

var ConnectionDef = g.NewInterface(
	"Connections",
	"Connection",
	g.KindOfT[AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-connection",
	g.NewQueryStruct("CreateConnection").
		Create().
		SQL("CONNECTION").
		IfNotExists().
		Name().
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateReplicated",
	"https://docs.snowflake.com/en/sql-reference/sql/create-connection",
	g.NewQueryStruct("CreateReplicated").
		Create().
		SQL("CONNECTION").
		IfNotExists().
		Name().
		//SQL("AS REPLICA OF").
		// external reference to connection: <orgnization_name>.<account_name>.<name>
		Identifier("ReplicaOf", g.KindOfT[ExternalObjectIdentifier](), g.IdentifierOptions().Required().SQL("AS REPLICA OF")).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "ReplicaOf"),
).CustomOperation(
	"AlterFailover",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-connection",
	g.NewQueryStruct("AlterFailover").
		Alter().
		SQL("CONNECTION").
		Name().
		OptionalQueryStructField(
			"EnableConnectionFailover",
			g.NewQueryStruct("EnableConnectionFailover").
				List("ToAccounts", "AccountIdentifier", g.ListOptions().NoParentheses()).
				OptionalSQL("IGNORE EDITION CHECK"),
			g.KeywordOptions().SQL("ENABLE FAILOVER TO ACCOUNTS"),
		).
		OptionalQueryStructField(
			"DisableConnectionFailover",
			g.NewQueryStruct("DisableConnectionFailover").
				OptionalSQL("TO ACCOUNTS").
				List("Accounts", "AccountIdentifier", g.ListOptions().NoParentheses()),
			g.KeywordOptions().SQL("DISABLE FAILOVER"),
		).
		OptionalSQL("PRIMARY").
		WithValidation(g.ExactlyOneValueSet, "EnableConnectionFailover", "DisableConnectionFailover", "Primary"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-connection",
	g.NewQueryStruct("Alter").
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
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-connection",
	g.NewQueryStruct("DropConnection").
		Drop().
		SQL("CONNECTION").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-connections",
	g.DbStruct("connectionRow").
		Field("snowflake_region", "string").
		Field("created_on", "time.Time").
		Field("account_name", "string").
		Field("name", "string").
		Field("comment", "sql.NullString").
		Field("is_primary", "string").
		Field("primary", "string").
		Field("failover_allowed_to_accounts", "string").
		Field("connection_url", "string").
		Field("organization_name", "string").
		Field("account_locator", "string"),
	g.PlainStruct("Connection").
		Field("SnowflakeRegion", "string").
		Field("CreatedOn", "time.Time").
		Field("AccountName", "string").
		Field("Name", "string").
		Field("Comment", "*string").
		Field("IsPrimary", "bool").
		Field("Primary", "string").
		Field("FailoverAllowedToAccounts", "[]string").
		Field("ConnectionUrl", "string").
		Field("OrganizationName", "string").
		Field("AccountLocator", "string"),
	g.NewQueryStruct("ShowConnections").
		Show().
		SQL("CONNECTIONS").
		OptionalLike(),
).ShowByIdOperation()
