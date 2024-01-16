package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var rowAccessPolicyDbRow = g.DbStruct("rowAccessPolicyDBRow").
	Text("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Text("kind").
	Text("owner").
	OptionalText("comment").
	Text("options").
	Bool("owner_role_type")

var rowAccessPolicy = g.PlainStruct("RowAccessPolicy").
	Text("CreatedOn").
	Text("Name").
	Text("DatabaseName").
	Text("SchemaName").
	Text("Kind").
	Text("Owner").
	OptionalText("Comment").
	Text("Options").
	Bool("OwnerRoleType")

var RowAccessPoliciesDef = g.NewInterface(
	"RowAccessPolicies",
	"RowAccessPolicy",
	g.KindOfT[SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-row-access-policy",
		g.NewQueryStruct("CreateRowAccessPolicy").
			Create().
			SQL("ROW ACCESS POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-row-access-policy",
		g.NewQueryStruct("AlterRowAccessPolicy").
			Drop().
			SQL("ROW ACCESS POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-row-access-policy",
		g.NewQueryStruct("DropRowAccessPolicy").
			Drop().
			SQL("ROW ACCESS POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-row-access-policies",
		rowAccessPolicyDbRow,
		rowAccessPolicy,
		g.NewQueryStruct("ShowRowAccessPolicies").
			Show().
			SQL("ROW ACCESS POLICIES").
			OptionalLike().
			OptionalIn(),
	).
	ShowByIdOperation().
	DescribeOperation(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-row-access-policy",
		g.DbStruct("describeRowAccessPolicyDBRow").
			Field("name", "string").
			Field("signature", "string").
			Field("return_type", "string").
			Field("body", "string"),
		g.PlainStruct("RowAccessPolicyDescription").
			Field("Name", "string").
			Field("Signature", "string").
			Field("ReturnType", "string").
			Field("Body", "string"),
		g.NewQueryStruct("DescribeRowAccessPolicy").
			Describe().
			SQL("ROW ACCESS POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
