package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

// UserTaskManagedInitialWarehouseSizeOptions is based on https://docs.snowflake.com/en/sql-reference/sql/USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE.

var TasksDef = g.NewInterface(
	"Tasks",
	"Task",
	g.KindOfT[SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-task",
		g.QueryStruct("CreateTask").
			Create().
			OrReplace().
			SQL("TASK").
			IfNotExists().
			Name().
			OptionalQueryStructField(
				"Warehouse",
				g.QueryStruct("CreateTaskWarehouse").
					OptionalIdentifier("Warehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("WAREHOUSE")).
					OptionalTextAssignment("USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE", g.ParameterOptions().SingleQuotes()).
					WithValidation(g.ValidIdentifier, "Warehouse").
					WithValidation(g.ExactlyOneValueSet, "Warehouse", "UserTaskManagedInitialWarehouseSize"),
				g.KeywordOptions(),
			).
			OptionalTextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("CONFIG", g.ParameterOptions().NoQuotes()).
			OptionalBooleanAssignment("ALLOW_OVERLAPPING_EXECUTION", nil).
			OptionalSessionParameters().
			OptionalNumberAssignment("USER_TASK_TIMEOUT_MS", nil).
			OptionalNumberAssignment("SUSPEND_TASK_AFTER_NUM_FAILURES", nil).
			OptionalTextAssignment("ERROR_INTEGRATION", g.ParameterOptions().NoQuotes()).
			OptionalSQL("COPY GRANTS").
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			ListAssignment("AFTER", "SchemaObjectIdentifier", nil).
			WithTags().
			OptionalTextAssignment("WHEN", g.ParameterOptions().NoQuotes().NoEquals()).
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	)
