// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/sdk/poc/generator"

//go:generate go run ./poc/main.go

var (
	onStreamDef = g.QueryStruct("OnStream").
			OptionalSQL("AT").
			OptionalSQL("BEFORE").
			QueryStructField(
			"Statement",
			g.QueryStruct("OnStreamStatement").
				OptionalTextAssignment("TIMESTAMP", g.ParameterOptions().ArrowEquals()).
				OptionalTextAssignment("OFFSET", g.ParameterOptions().ArrowEquals()).
				OptionalTextAssignment("STATEMENT", g.ParameterOptions().ArrowEquals()).
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
				Field("tableOn", "sql.NullString").
				Field("owner", "sql.NullString").
				Field("comment", "sql.NullString").
				Field("table_name", "sql.NullString").
				Field("source_type", "sql.NullString").
				Field("base_tables", "sql.NullString").
				Field("type", "sql.NullString").
				Field("stale", "sql.NullString").
				Field("mode", "sql.NullString").
				Field("stale_after", "sql.NullTime").
				Field("invalid_reason", "sql.NullString").
				Field("owner_role_type", "sql.NullString")

	streamPlainStructDef = g.PlainStruct("Stream").
				Field("CreatedOn", "time.Time").
				Field("Name", "string").
				Field("DatabaseName", "string").
				Field("SchemaName", "string").
				Field("TableOn", "*string").
				Field("Owner", "*string").
				Field("Comment", "*string").
				Field("TableName", "*string").
				Field("SourceType", "*string").
				Field("BaseTables", "*string").
				Field("Type", "*string").
				Field("Stale", "*string").
				Field("Mode", "*string").
				Field("StaleAfter", "*time.Time").
				Field("InvalidReason", "*string").
				Field("OwnerRoleType", "*string")

	StreamsDef = g.NewInterface(
		"Streams",
		"Stream",
		g.KindOfT[SchemaObjectIdentifier](),
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
			"CreateOnDirectoryTable",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stream",
			g.QueryStruct("CreateStreamOnDirectoryTable").
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
				Identifier("sourceStream", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("CLONE").Required()).
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
				OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
				OptionalSQL("UNSET COMMENT").
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
