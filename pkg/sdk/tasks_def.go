package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var taskDbRow = g.DbStruct("taskDBRow").
	Field("created_on", "string").
	Field("name", "string").
	Field("id", "string").
	Field("database_name", "string").
	Field("schema_name", "string").
	Field("owner", "string").
	Field("comment", "string").
	Field("warehouse", "string").
	Field("schedule", "string").
	Field("predecessors", "string").
	Field("state", "string").
	Field("definition", "string").
	Field("condition", "string").
	Field("allow_overlapping_execution", "string").
	Field("error_integration", "string").
	Field("last_committed_on", "string").
	Field("last_suspended_on", "string").
	Field("owner_role_type", "string").
	Field("config", "string").
	Field("budget", "string")

var task = g.PlainStruct("Task").
	Field("CreatedOn", "string").
	Field("Name", "string").
	Field("Id", "string").
	Field("DatabaseName", "string").
	Field("SchemaName", "string").
	Field("Owner", "string").
	Field("Comment", "string").
	Field("Warehouse", "string").
	Field("Schedule", "string").
	Field("Predecessors", "string").
	Field("State", "string").
	Field("Definition", "string").
	Field("Condition", "string").
	Field("AllowOverlappingExecution", "string").
	Field("ErrorIntegration", "string").
	Field("LastCommittedOn", "string").
	Field("LastSuspendedOn", "string").
	Field("OwnerRoleType", "string").
	Field("Config", "string").
	Field("Budget", "string")

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
					OptionalIdentifier("Warehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WAREHOUSE")).
					OptionalAssignment("USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE", "WarehouseSize", g.ParameterOptions().SingleQuotes()).
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
			ListAssignment("AFTER", "SchemaObjectIdentifier", g.ParameterOptions().NoEquals()).
			WithTags().
			OptionalTextAssignment("WHEN", g.ParameterOptions().NoQuotes().NoEquals()).
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-task",
		g.QueryStruct("AlterTask").
			Alter().
			SQL("TASK").
			IfExists().
			Name().
			OptionalSQL("RESUME").
			OptionalSQL("SUSPEND").
			ListAssignment("REMOVE AFTER", "SchemaObjectIdentifier", g.ParameterOptions().NoEquals()).
			ListAssignment("ADD AFTER", "SchemaObjectIdentifier", g.ParameterOptions().NoEquals()).
			OptionalQueryStructField(
				"Set",
				g.QueryStruct("TaskSet").
					OptionalIdentifier("Warehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("WAREHOUSE")).
					OptionalTextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("CONFIG", g.ParameterOptions().NoQuotes()).
					OptionalBooleanAssignment("ALLOW_OVERLAPPING_EXECUTION", nil).
					OptionalNumberAssignment("USER_TASK_TIMEOUT_MS", nil).
					OptionalNumberAssignment("SUSPEND_TASK_AFTER_NUM_FAILURES", nil).
					OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
					OptionalSessionParameters().
					WithValidation(g.AtLeastOneValueSet, "Warehouse", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "Comment", "SessionParameters"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.QueryStruct("TaskUnset").
					OptionalSQL("WAREHOUSE").
					OptionalSQL("SCHEDULE").
					OptionalSQL("CONFIG").
					OptionalSQL("ALLOW_OVERLAPPING_EXECUTION").
					OptionalSQL("USER_TASK_TIMEOUT_MS").
					OptionalSQL("SUSPEND_TASK_AFTER_NUM_FAILURES").
					OptionalSQL("COMMENT").
					OptionalSessionParametersUnset().
					WithValidation(g.AtLeastOneValueSet, "Warehouse", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "Comment", "SessionParametersUnset"),
				g.KeywordOptions().SQL("UNSET"),
			).
			SetTags().
			UnsetTags().
			OptionalTextAssignment("MODIFY AS", g.ParameterOptions().NoQuotes().NoEquals()).
			OptionalTextAssignment("MODIFY WHEN", g.ParameterOptions().NoQuotes().NoEquals()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Resume", "Suspend", "RemoveAfter", "AddAfter", "Set", "Unset", "SetTags", "UnsetTags", "ModifyAs", "ModifyWhen"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-task",
		g.QueryStruct("DropTask").
			Drop().
			SQL("TASK").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-tasks",
		taskDbRow,
		task,
		g.QueryStruct("ShowTasks").
			Show().
			Terse().
			SQL("TASKS").
			OptionalLike().
			OptionalIn().
			OptionalStartsWith().
			OptionalSQL("ROOT ONLY").
			OptionalLimit(),
	).
	ShowByIdOperation().
	DescribeOperation(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-task",
		taskDbRow,
		task,
		g.QueryStruct("DescribeTask").
			Describe().
			SQL("TASK").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	CustomOperation(
		"Execute",
		"https://docs.snowflake.com/en/sql-reference/sql/execute-task",
		g.QueryStruct("ExecuteTask").
			SQL("EXECUTE").
			SQL("TASK").
			Name().
			OptionalSQL("RETRY LAST").
			WithValidation(g.ValidIdentifier, "name"),
	)
