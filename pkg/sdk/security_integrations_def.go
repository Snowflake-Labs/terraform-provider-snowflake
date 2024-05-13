package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

func createSecurityIntegrationOperation(structName string, apply func(qs *g.QueryStruct) *g.QueryStruct) *g.QueryStruct {
	qs := g.NewQueryStruct(structName).
		Create().
		OrReplace().
		SQL("SECURITY INTEGRATION").
		IfNotExists().
		Name()
	qs = apply(qs)
	return qs.
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists")
}

func alterSecurityIntegrationOperation(structName string, apply func(qs *g.QueryStruct) *g.QueryStruct) *g.QueryStruct {
	qs := g.NewQueryStruct(structName).
		Alter().
		SQL("SECURITY INTEGRATION").
		IfExists().
		Name()
	qs = apply(qs)
	return qs.
		NamedList("SET TAG", g.KindOfT[TagAssociation]()).
		NamedList("UNSET TAG", g.KindOfT[ObjectIdentifier]()).
		WithValidation(g.ValidIdentifier, "name")
}

var scimIntegrationSetDef = g.NewQueryStruct("SCIMIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
	OptionalBooleanAssignment("SYNC_PASSWORD", g.ParameterOptions()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "NetworkPolicy", "SyncPassword", "Comment")

var scimIntegrationUnsetDef = g.NewQueryStruct("SCIMIntegrationUnset").
	OptionalSQL("NETWORK_POLICY").
	OptionalSQL("SYNC_PASSWORD").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "NetworkPolicy", "SyncPassword", "Comment")

var SecurityIntegrationsDef = g.NewInterface(
	"SecurityIntegrations",
	"SecurityIntegration",
	g.KindOfT[AccountObjectIdentifier](),
).
	CustomOperation(
		"CreateSCIM",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-scim",
		createSecurityIntegrationOperation("CreateSCIMIntegration", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = SCIM")).
				BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
				TextAssignment("SCIM_CLIENT", g.ParameterOptions().Required().SingleQuotes()).
				TextAssignment("RUN_AS_ROLE", g.ParameterOptions().Required().SingleQuotes()).
				OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
				OptionalBooleanAssignment("SYNC_PASSWORD", g.ParameterOptions())
		}),
	).
	CustomOperation(
		"AlterSCIMIntegration",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-scim",
		alterSecurityIntegrationOperation("AlterSCIMIntegration", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				scimIntegrationSetDef,
				g.KeywordOptions().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				scimIntegrationUnsetDef,
				g.KeywordOptions().SQL("UNSET"),
			)
		}),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropSecurityIntegration").
			Drop().
			SQL("SECURITY INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		g.DbStruct("securityIntegrationDescRow").
			Field("property", "string").
			Field("property_type", "string").
			Field("property_value", "string").
			Field("property_default", "string"),
		g.PlainStruct("SecurityIntegrationProperty").
			Field("Name", "string").
			Field("Type", "string").
			Field("Value", "string").
			Field("Default", "string"),
		g.NewQueryStruct("DescSecurityIntegration").
			Describe().
			SQL("SECURITY INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.DbStruct("securityIntegrationShowRow").
			Text("name").
			Text("type").
			Text("category").
			Bool("enabled").
			OptionalText("comment").
			Time("created_on"),
		g.PlainStruct("SecurityIntegration").
			Text("Name").
			Text("IntegrationType").
			Text("Category").
			Bool("Enabled").
			Text("Comment").
			Time("CreatedOn"),
		g.NewQueryStruct("ShowSecurityIntegrations").
			Show().
			SQL("SECURITY INTEGRATIONS").
			OptionalLike(),
	).
	ShowByIdOperation()
