package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var (
	onStreamDef = g.QueryStruct("OnStream").
		// TODO: AT / BEFORE enum
		OptionalSQL("AT").
		OptionalSQL("BEFORE").
		QueryStructField(
			// TODO: Rename
			"Statement",
			g.QueryStruct("OnStreamStatement").
				OptionalTextAssignment("TIMESTAMP", g.ParameterOptions().ArrowEquals().DoubleQuotes()).
				OptionalTextAssignment("OFFSET", g.ParameterOptions().ArrowEquals().DoubleQuotes()).
				OptionalTextAssignment("STATEMENT", g.ParameterOptions().ArrowEquals().DoubleQuotes()).
				OptionalTextAssignment("STREAM", g.ParameterOptions().ArrowEquals().SingleQuotes()).
				WithValidation(g.ExactlyOneValueSet, "Timestamp", "Offset", "Statement", "Stream"),
			g.ListOptions().Parentheses(),
		).
		WithValidation(g.ExactlyOneValueSet, "At", "Before")

	showStreamDbRowDef = g.DbStruct("showStreamsDbRow").
				Field("created_on", "time.Time").
				Field("name", "string").
				Field("database_name", "string").
				Field("schema_name", "string").
				Field("owner", "string").
				Field("comment", "string").
				Field("table_name", "string").
				Field("source_type", "string").
				Field("base_tables", "string").
				Field("type", "string").
				Field("stale", "string").
				Field("mode", "string").
				Field("stale_after", "sql.NullTime").
				Field("invalid_reason", "string").
				Field("owner_role_type", "string")

	streamPlainStructDef = g.PlainStruct("Stream").
				Field("CreatedOn", "time.Time").
				Field("Name", "string").
				Field("DatabaseName", "string").
				Field("SchemaName", "string").
				Field("Owner", "string").
				Field("Comment", "string").
				Field("TableName", "string").
				Field("SourceType", "string").
				Field("BaseTables", "string").
				Field("Type", "string").
				Field("Stale", "string").
				Field("Mode", "string").
				Field("StaleAfter", "*time.Time").
				Field("InvalidReason", "string").
				Field("OwnerRoleType", "string")

	StreamsDef = g.NewInterface(
		"Streams",
		"Stream",
		g.KindOfT[AccountObjectIdentifier](),
	).
		CustomOperation(
			"CreateOnTable",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stream",
			g.QueryStruct("CreateStreamOnTable").
				Create().
				OrReplace().
				SQL("STREAM").
				IfNotExists().
				Name().
				OptionalCopyGrants().
				SQL("ON TABLE").
				Identifier("TableId", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
				OptionalQueryStructField("On", onStreamDef, g.KeywordOptions()).
				OptionalBooleanAssignment("APPEND_ONLY", nil).
				OptionalBooleanAssignment("SHOW_INITIAL_ROWS", nil).
				OptionalComment().
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ValidIdentifier, "TableId").
				WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		).
		CustomOperation(
			"CreateOnExternalTable",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stream",
			g.QueryStruct("CreateStreamOnExternalTable").
				Create().
				OrReplace().
				SQL("STREAM").
				IfNotExists().
				Name().
				OptionalCopyGrants().
				SQL("ON EXTERNAL TABLE").
				Identifier("ExternalTableId", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
				OptionalQueryStructField("On", onStreamDef, g.KeywordOptions()).
				OptionalBooleanAssignment("INSERT_ONLY", nil).
				OptionalComment().
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ValidIdentifier, "ExternalTableId").
				WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		).
		CustomOperation(
			"CreateOnStage",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stream",
			g.QueryStruct("CreateStreamOnStage").
				Create().
				OrReplace().
				SQL("STREAM").
				IfNotExists().
				Name().
				OptionalCopyGrants().
				SQL("ON STAGE").
				Identifier("StageId", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
				OptionalComment().
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ValidIdentifier, "StageId").
				WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		).
		CustomOperation(
			"CreateOnView",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stream",
			g.QueryStruct("CreateStreamOnView").
				Create().
				OrReplace().
				SQL("STREAM").
				IfNotExists().
				Name().
				OptionalCopyGrants().
				SQL("ON VIEW").
				Identifier("ViewId", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
				OptionalQueryStructField("On", onStreamDef, g.KeywordOptions()).
				OptionalBooleanAssignment("APPEND_ONLY", nil).
				OptionalBooleanAssignment("SHOW_INITIAL_ROWS", nil).
				OptionalComment().
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ValidIdentifier, "ViewId").
				WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		).
		CustomOperation(
			"Clone",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stream#variant-syntax",
			g.QueryStruct("CloneStream").
				Create().
				OrReplace().
				SQL("STREAM").
				Name().
				Identifier("sourceStream", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("CLONE").Required()).
				OptionalCopyGrants().
				WithValidation(g.ValidIdentifier, "name"),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-stream",
			g.QueryStruct("AlterStream").
				Alter().
				SQL("STREAM").
				IfExists().
				Name().
				SetComment().
				UnsetComment().
				SetTags().
				UnsetTags().
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ConflictingFields, "IfExists", "UnsetTags").
				WithValidation(g.ExactlyOneValueSet, "SetComment", "UnsetComment", "SetTags", "UnsetTags"),
		).
		DropOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/drop-stream",
			g.QueryStruct("DropStream").
				Drop().
				SQL("STREAM").
				IfExists().
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		).
		ShowOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/show-streams",
			showStreamDbRowDef,
			streamPlainStructDef,
			g.QueryStruct("ShowStreams").
				Show().
				Terse().
				SQL("STREAMS").
				OptionalLike().
				OptionalIn().
				OptionalStartsWith().
				OptionalLimit(),
		).
		ShowByIdOperation().
		DescribeOperation(
			g.DescriptionMappingKindSingleValue,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-stream",
			showStreamDbRowDef,
			streamPlainStructDef,
			g.QueryStruct("DescribeStream").
				Describe().
				SQL("STREAM").
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		)
)
