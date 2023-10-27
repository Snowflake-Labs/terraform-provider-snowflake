package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

type InternalStageEncryptionOption string

const (
	InternalStageEncryptionFull InternalStageEncryptionOption = "SNOWFLAKE_FULL"
	InternalStageEncryptionSSE  InternalStageEncryptionOption = "SNOWFLAKE_SSE"
)

type StageCopyColumnMapOption string

const (
	StageCopyColumnMapCaseSensitive   StageCopyColumnMapOption = "CASE_SENSITIVE"
	StageCopyColumnMapCaseInsensitive StageCopyColumnMapOption = "CASE_INSENSITIVE"
	StageCopyColumnMapCaseNone        StageCopyColumnMapOption = "NONE"
)

// TODO PUT, GET, LS, etc. ???

var (
	StagesDef = g.NewInterface(
		"Stages",
		"Stage",
		g.KindOfT[SchemaObjectIdentifier](),
	).
		CustomOperation(
			"CreateInternal",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stage",
			g.QueryStruct("CreateInternalStage").
				Create().
				OrReplace().
				OptionalSQL("TEMPORARY").
				SQL("STAGE").
				IfNotExists().
				Name().
				OptionalQueryStructField(
					"Encryption",
					g.QueryStruct("InternalStageEncryption").
						OptionalAssignment(
							"TYPE",
							g.KindOfT[InternalStageEncryptionOption](),
							g.ParameterOptions().SingleQuotes(),
						),
					g.ParameterOptions().Parentheses().SQL("ENCRYPTION"),
				).
				OptionalQueryStructField(
					"DirectoryTableOptions",
					g.QueryStruct("InternalDirectoryTableOptions").
						OptionalBooleanAssignment("ENABLE", nil).
						OptionalBooleanAssignment("REFRESH_ON_CREATE", nil),
					g.ParameterOptions().Parentheses(),
				).
				OptionalQueryStructField(
					"FileFormat",
					g.QueryStruct("StageFileFormat").
						OptionalTextAssignment("FORMAT_NAME", g.ParameterOptions().SingleQuotes()).
						OptionalAssignment("TYPE", g.KindOfTPointer[FileFormatType](), g.ParameterOptions()).
						List("TYPE", g.KindOfT[FileFormatType](), g.ListOptions()),
					g.ParameterOptions().Parentheses().SQL("FILE_FORMAT"),
				).
				OptionalQueryStructField(
					"CopyOptions",
					g.QueryStruct("StageCopyOptions").
						OptionalQueryStructField(
							"OnError",
							g.QueryStruct("StageCopyOnErrorOptions").
								OptionalSQL("CONTINUE").
								OptionalSQL("SKIP_FILE").
								OptionalSQL("SKIP_FILE"). // TODO templated value
								OptionalSQL("SKIP_FILE"). // TODO templated value with %
								OptionalSQL("ABORT_STATEMENT"),
							g.ParameterOptions().SQL("ON_ERROR"),
						).
						OptionalNumberAssignment("SIZE_LIMIT", nil).
						OptionalBooleanAssignment("PURGE", nil).
						OptionalBooleanAssignment("RETURN_FAILED_ONLY", nil).
						OptionalAssignment("MATCH_BY_COLUMN_NAME", g.KindOfTPointer[StageCopyColumnMapOption](), nil).
						OptionalBooleanAssignment("ENFORCE_LENGTH", nil).
						OptionalBooleanAssignment("TRUNCATECOLUMNS", nil).
						OptionalBooleanAssignment("FORCE", nil).
						OptionalAssignment("TYPE", g.KindOfTPointer[FileFormatType](), g.ParameterOptions()).
						List("TYPE", g.KindOfT[FileFormatType](), g.ListOptions()),
					g.ParameterOptions().Parentheses().SQL("FILE_FORMAT"),
				).
				OptionalComment().
				OptionalTags(),
		).
		DropOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/drop-stage",
			g.QueryStruct("DropStage").
				Drop().
				SQL("STAGE").
				IfExists().
				Name(),
		).
		DescribeOperation(
			g.DescriptionMappingKindSlice,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-stage",
			g.DbStruct("stageDescRow").
				Field("parent_property", "string").
				Field("property", "string").
				Field("property_type", "string").
				Field("property_value", "sql.NullString").
				Field("property_default", "sql.NullString"),
			g.PlainStruct("StageProperty").
				Field("Parent", "string").
				Field("Name", "string").
				Field("Type", "string").
				Field("Value", "*string").
				Field("Default", "*string"),
			g.QueryStruct("DescStage").
				Describe().
				SQL("STAGE").
				Name(),
		).
		ShowOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/show-stages",
			g.DbStruct("stageShowRow").
				Field("created_on", "time.Time").
				Field("name", "string").
				Field("database_name", "string").
				Field("schema_name", "string").
				Field("url", "string").
				Field("has_credentials", "string").
				Field("has_encryption_key", "string").
				Field("owner", "string").
				Field("comment", "string").
				Field("region", "string").
				Field("type", "string").
				Field("cloud", "sql.NullString").
				Field("storage_integration", "sql.NullString").
				Field("endpoint", "sql.NullString").
				Field("owner_role_type", "sql.NullString").
				Field("directory_enabled", "string"),
			g.PlainStruct("Stage").
				Field("CreatedOn", "time.Time").
				Field("Name", "string").
				Field("DatabaseName", "string").
				Field("SchemaName", "string").
				Field("Url", "string").
				Field("HasCredentials", "bool").
				Field("HasEncryptionKey", "bool").
				Field("Owner", "string").
				Field("Comment", "string").
				Field("Region", "string").
				Field("Type", "string").
				Field("Cloud", "*string").
				Field("StorageIntegration", "*string").
				Field("Endpoint", "*string").
				Field("OwnerRoleType", "*string").
				Field("DirectoryEnabled", "bool"),
			g.QueryStruct("ShowStages").
				Show().
				SQL("STAGES").
				OptionalLike().
				OptionalIn(),
		)
)
