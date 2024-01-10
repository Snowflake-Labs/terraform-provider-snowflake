package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var StorageIntegrationDef = g.NewInterface(
	"StorageIntegrations",
	"StorageIntegration",
	g.KindOfT[AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-storage-integration",
		g.NewQueryStruct("CreateStorageIntegration").
			Create().
			OrReplace().
			SQL("STORAGE INTEGRATION").
			IfNotExists().
			Name().
			PredefinedQueryStructField("externalStageType", "string", g.StaticOptions().SQL("TYPE = EXTERNAL STAGE")).
			OptionalQueryStructField(
				"S3StorageProviderParams",
				g.NewQueryStruct("S3StorageParams").
					PredefinedQueryStructField("storageProvider", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'S3'")).
					TextAssignment("STORAGE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("STORAGE_AWS_OBJECT_ACL", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GCSStorageProviderParams",
				g.NewQueryStruct("GCSStorageParams").
					PredefinedQueryStructField("storageProvider", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'GCS'")),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"AzureStorageProviderParams",
				g.NewQueryStruct("AzureStorageParams").
					PredefinedQueryStructField("storageProvider", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'AZURE'")).
					OptionalTextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()),
				g.KeywordOptions(),
			).
			BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
			ListAssignment("STORAGE_ALLOWED_LOCATIONS", "string", g.ParameterOptions().Required()).
			ListAssignment("STORAGE_BLOCKED_LOCATIONS", "string", g.ParameterOptions()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithValidation(g.ExactlyOneValueSet, "S3StorageProviderParams", "GCSStorageProviderParams", "AzureStorageProviderParams"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-storage-integration",
		g.NewQueryStruct("AlterStorageIntegration").
			Alter().
			SQL("STORAGE INTEGRATION").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("StorageIntegrationSet").
					OptionalQueryStructField(
						"SetS3Params",
						g.NewQueryStruct("SetS3StorageParams").
							TextAssignment("STORAGE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
							OptionalTextAssignment("STORAGE_AWS_OBJECT_ACL", g.ParameterOptions().SingleQuotes()),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"SetAzureParams",
						g.NewQueryStruct("SetAzureStorageParams").
							OptionalTextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					BooleanAssignment("ENABLED", g.ParameterOptions()).
					ListAssignment("STORAGE_ALLOWED_LOCATIONS", "string", g.ParameterOptions()).
					ListAssignment("STORAGE_BLOCKED_LOCATIONS", "string", g.ParameterOptions()).
					OptionalComment(),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("StorageIntegrationUnset").
					OptionalSQL("ENABLED").
					OptionalSQL("STORAGE_BLOCKED_LOCATIONS").
					OptionalSQL("COMMENT"),
				g.KeywordOptions().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfExists", "UnsetTags").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
	)
