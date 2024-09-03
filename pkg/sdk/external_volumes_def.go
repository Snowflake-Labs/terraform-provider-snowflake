package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

type (
	S3EncryptionType  string
	S3StorageProvider string
	GCSEncryptionType string
)

var (
	S3EncryptionTypeSseS3   S3EncryptionType  = "AWS_SSE_S3"
	S3EncryptionTypeSseKms  S3EncryptionType  = "AWS_SSE_KMS"
	S3EncryptionNone        S3EncryptionType  = "NONE"
	GCSEncryptionTypeSseKms GCSEncryptionType = "GCS_SSE_KMS"
	GCSEncryptionTypeNone   GCSEncryptionType = "NONE"
	S3StorageProviderS3     S3StorageProvider = "S3"
	S3StorageProviderS3GOV  S3StorageProvider = "S3GOV"
)

var externalS3StorageLocationDef = g.NewQueryStruct("S3StorageLocationParams").
	TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
	Assignment("STORAGE_PROVIDER", g.KindOfT[S3StorageProvider](), g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	OptionalTextAssignment("STORAGE_AWS_EXTERNAL_ID", g.ParameterOptions().SingleQuotes()).
	OptionalQueryStructField(
		"Encryption",
		g.NewQueryStruct("ExternalVolumeS3Encryption").
			Assignment("TYPE", g.KindOfT[S3EncryptionType](), g.ParameterOptions().SingleQuotes().Required()).
			OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
		g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
	)

var externalGCSStorageLocationDef = g.NewQueryStruct("GCSStorageLocationParams").
	TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
	PredefinedQueryStructField("StorageProviderGcs", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'GCS'")).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	OptionalQueryStructField(
		"Encryption",
		g.NewQueryStruct("ExternalVolumeGCSEncryption").
			Assignment("TYPE", g.KindOfT[GCSEncryptionType](), g.ParameterOptions().SingleQuotes().Required()).
			OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
		g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
	)

var externalAzureStorageLocationDef = g.NewQueryStruct("AzureStorageLocationParams").
	TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
	PredefinedQueryStructField("StorageProviderAzure", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'AZURE'")).
	TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required())

// Can't name StorageLocation due to naming clash with type in storage integration
var storageLocationDef = g.NewQueryStruct("ExternalVolumeStorageLocation").
	OptionalQueryStructField(
		"S3StorageLocationParams",
		externalS3StorageLocationDef,
		g.ListOptions().Parentheses().NoComma(),
	).
	OptionalQueryStructField(
		"GCSStorageLocationParams",
		externalGCSStorageLocationDef,
		g.ListOptions().Parentheses().NoComma(),
	).
	OptionalQueryStructField(
		"AzureStorageLocationParams",
		externalAzureStorageLocationDef,
		g.ListOptions().Parentheses().NoComma(),
	).
	WithValidation(g.ExactlyOneValueSet, "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams")

var ExternalVolumesDef = g.NewInterface(
	"ExternalVolumes",
	"ExternalVolume",
	g.KindOfT[AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-external-volume",
		g.NewQueryStruct("CreateExternalVolume").
			Create().
			OrReplace().
			SQL("EXTERNAL VOLUME").
			IfNotExists().
			Name().
			ListAssignment("STORAGE_LOCATIONS", "ExternalVolumeStorageLocation", g.ParameterOptions().Parentheses().Required()).
			OptionalBooleanAssignment("ALLOW_WRITES", nil).
			OptionalComment().
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
			WithValidation(g.ValidIdentifier, "name"),
		storageLocationDef,
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-external-volume",
		g.NewQueryStruct("AlterExternalVolume").
			Alter().
			SQL("EXTERNAL VOLUME").
			IfExists().
			Name().
			OptionalTextAssignment("REMOVE STORAGE_LOCATION", g.ParameterOptions().SingleQuotes().NoEquals()).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("AlterExternalVolumeSet").
					OptionalBooleanAssignment("ALLOW_WRITES", g.ParameterOptions()).
					OptionalComment(),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"AddStorageLocation",
				storageLocationDef,
				g.ParameterOptions().SQL("ADD STORAGE_LOCATION"),
			).
			WithValidation(g.ExactlyOneValueSet, "RemoveStorageLocation", "Set", "AddStorageLocation").
			WithValidation(g.ValidIdentifier, "name"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-external-volume",
		g.NewQueryStruct("DropExternalVolume").
			Drop().
			SQL("EXTERNAL VOLUME").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-external-volume",
		g.DbStruct("externalVolumeDescRow").
			Field("parent_property", "string").
			Field("property", "string").
			Field("property_type", "string").
			Field("property_value", "string").
			Field("property_default", "string"),
		g.PlainStruct("ExternalVolumeProperty").
			Field("Parent", "string").
			Field("Name", "string").
			Field("Type", "string").
			Field("Value", "string").
			Field("Default", "string"),
		g.NewQueryStruct("DescExternalVolume").
			Describe().
			SQL("EXTERNAL VOLUME").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-external-volumes",
		g.DbStruct("externalVolumeShowRow").
			Field("name", "string").
			Field("allow_writes", "string").
			Field("comment", "string"),
		g.PlainStruct("ExternalVolume").
			Field("Name", "string").
			Field("AllowWrites", "string").
			Field("Comment", "string"),
		g.NewQueryStruct("ShowExternalVolumes").
			Show().
			SQL("EXTERNAL VOLUMES").
			OptionalLike(),
	).
	ShowByIdOperation()
