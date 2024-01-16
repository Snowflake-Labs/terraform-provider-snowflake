package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var managedAccountDbRow = g.DbStruct("managedAccountDBRow").
	Text("name").
	Text("cloud").
	Text("region").
	Text("locator").
	Text("created_on").
	Text("url").
	Text("account_locator_url").
	Bool("is_reader").
	OptionalText("comment")

var managedAccount = g.PlainStruct("ManagedAccount").
	Text("Name").
	Text("Cloud").
	Text("Region").
	Text("Locator").
	Text("CreatedOn").
	Text("URL").
	Text("AccountLocatorURL").
	Bool("IsReader").
	OptionalText("Comment")

var ManagedAccountsDef = g.NewInterface(
	"ManagedAccounts",
	"ManagedAccount",
	g.KindOfT[AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-managed-account",
		g.NewQueryStruct("CreateManagedAccount").
			Create().
			SQL("MANAGED ACCOUNT").
			Name().
			QueryStructField(
				"CreateManagedAccountParams",
				g.NewQueryStruct("CreateManagedAccountParams").
					TextAssignment("ADMIN_NAME", g.ParameterOptions().NoQuotes().Required()).
					TextAssignment("ADMIN_PASSWORD", g.ParameterOptions().SingleQuotes().Required()).
					PredefinedQueryStructField("typeProvider", "string", g.StaticOptions().SQL("TYPE = READER")).
					OptionalComment().
					WithValidation(g.ValidateValueSet, "AdminName").
					WithValidation(g.ValidateValueSet, "AdminPassword"),
				g.ListOptions().NoParentheses().Required(),
			).
			WithValidation(g.ValidIdentifier, "name"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-managed-account",
		g.NewQueryStruct("DropManagedAccount").
			Drop().
			SQL("MANAGED ACCOUNT").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-managed-accounts",
		managedAccountDbRow,
		managedAccount,
		g.NewQueryStruct("ShowManagedAccounts").
			Show().
			SQL("MANAGED ACCOUNTS").
			OptionalLike(),
	).
	ShowByIdOperation()
