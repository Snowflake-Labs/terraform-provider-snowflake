package resources_test

import (
	"fmt"
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	configvariable "github.com/hashicorp/terraform-plugin-testing/config"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO: More tests for complicated DAGs

func TestAcc_Task_Basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	configModel := model.TaskWithId("test", id, statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.test").
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanDefault).
						HasWarehouseString("").
						HasScheduleString("").
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfterLen(0).
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner("ACCOUNTADMIN"). // TODO: Current role
						HasComment("").
						HasWarehouse("").
						HasSchedule("").
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")). // TODO: *sdk.AOI
						HasLastCommittedOn("").
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget(""),
					//HasTaskRelations(sdk.TaskRelations{}). // TODO:
					resourceparametersassert.TaskResourceParameters(t, "snowflake_task.test").
						HasAllDefaults(),
				),
			},
			{
				ResourceName: "snowflake_task.test",
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedTaskResource(t, helpers.EncodeResourceIdentifier(id)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanFalse).
						HasWarehouseString("").
						HasScheduleString("").
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanFalse).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasNoAfter().
						HasWhenString("").
						HasSqlStatementString(statement),
				),
			},
		},
	})
}

func TestAcc_Task_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	errorNotificationIntegration, errorNotificationIntegrationCleanup := acc.TestClient().NotificationIntegration.Create(t)
	t.Cleanup(errorNotificationIntegrationCleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	taskConfig := `$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`
	// We have to do three $ at the beginning because Terraform will remove one $.
	// It's because `${` is a special pattern, and it's escaped by `$${`.
	expectedTaskConfig := strings.ReplaceAll(taskConfig, "$", "")
	taskConfigVariableValue := "$" + taskConfig
	comment := random.Comment()
	condition := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
	configModel := model.TaskWithId("test", id, statement).
		WithEnabled(r.BooleanTrue).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithSchedule("10 MINUTES").
		WithConfigValue(configvariable.StringVariable(taskConfigVariableValue)).
		WithAllowOverlappingExecution(true).
		WithErrorIntegration(errorNotificationIntegration.ID().Name()).
		WithComment(comment).
		WithWhen(condition)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.test").
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanTrue).
						HasWarehouseString(acc.TestClient().Ids.WarehouseId().Name()).
						HasScheduleString("10 MINUTES").
						HasConfigString(expectedTaskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasCommentString(comment).
						HasFinalizeString("").
						HasNoAfter().
						HasWhenString(condition).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						//HasId(id.FullyQualifiedName()). // TODO: not empty
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner("ACCOUNTADMIN"). // TODO: Current role
						HasComment(comment).
						HasWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
						HasSchedule("10 MINUTES").
						//HasPredecessors(nil). // TODO:
						HasState(sdk.TaskStateStarted).
						HasDefinition(statement).
						HasCondition(condition).
						HasAllowOverlappingExecution(true).
						HasErrorIntegration(errorNotificationIntegration.ID()).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig(expectedTaskConfig).
						HasBudget(""),
					//HasTaskRelations(sdk.TaskRelations{}). // TODO:
					resourceparametersassert.TaskResourceParameters(t, "snowflake_task.test").
						HasAllDefaults(),
				),
			},
			{
				ResourceName: "snowflake_task.test",
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedTaskResource(t, helpers.EncodeResourceIdentifier(id)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanTrue).
						HasWarehouseString(acc.TestClient().Ids.WarehouseId().Name()).
						HasScheduleString("10 MINUTES").
						HasConfigString(expectedTaskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasCommentString(comment).
						HasFinalizeString("").
						HasNoAfter().
						HasWhenString(condition).
						HasSqlStatementString(statement),
				),
			},
		},
	})
}

func TestAcc_Task_Updates(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	basicConfigModel := model.TaskWithId("test", id, statement)

	// TODO: Assert the rest of fields (e.g. parameters)

	errorNotificationIntegration, errorNotificationIntegrationCleanup := acc.TestClient().NotificationIntegration.Create(t)
	t.Cleanup(errorNotificationIntegrationCleanup)

	taskConfig := `$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`
	// We have to do three $ at the beginning because Terraform will remove one $.
	// It's because `${` is a special pattern, and it's escaped by `$${`.
	expectedTaskConfig := strings.ReplaceAll(taskConfig, "$", "")
	taskConfigVariableValue := "$" + taskConfig
	comment := random.Comment()
	condition := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
	completeConfigModel := model.TaskWithId("test", id, statement).
		WithEnabled(r.BooleanTrue).
		// TODO: Warehouse cannot be set (error)
		//WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithSchedule("10 MINUTES").
		WithConfigValue(configvariable.StringVariable(taskConfigVariableValue)).
		WithAllowOverlappingExecution(true).
		WithErrorIntegration(errorNotificationIntegration.ID().Name()).
		WithComment(comment).
		WithWhen(condition)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, basicConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.test").
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanDefault).
						HasWarehouseString("").
						HasScheduleString("").
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfterLen(0).
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.test").
						//HasCreatedOnNotEmpty(),
						HasName(id.Name()).
						//HasId(id.FullyQualifiedName()). // TODO: not empty
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner("ACCOUNTADMIN"). // TODO: Current role
						HasComment("").
						HasWarehouse("").
						HasSchedule("").
						//HasPredecessors(nil). // TODO:
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")). // TODO: *sdk.AOI
						HasLastCommittedOn("").
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget(""),
					//HasTaskRelations(sdk.TaskRelations{}). // TODO:
				),
			},
			// Set
			{
				Config: config.FromModel(t, completeConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.test").
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanTrue).
						//HasWarehouseString(acc.TestClient().Ids.WarehouseId().Name()).
						HasScheduleString("10 MINUTES").
						HasConfigString(expectedTaskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasCommentString(comment).
						HasFinalizeString("").
						HasAfterLen(0).
						HasWhenString(condition).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						//HasId(id.FullyQualifiedName()). // TODO: not empty
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner("ACCOUNTADMIN"). // TODO: Current role
						HasComment(comment).
						//HasWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
						HasSchedule("10 MINUTES").
						//HasPredecessors(nil). // TODO:
						HasState(sdk.TaskStateStarted).
						HasDefinition(statement).
						HasCondition(condition).
						HasAllowOverlappingExecution(true).
						HasErrorIntegration(errorNotificationIntegration.ID()).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig(expectedTaskConfig).
						HasBudget(""),
					//HasTaskRelations(sdk.TaskRelations{}). // TODO:
				),
			},
			// Unset
			{
				Config: config.FromModel(t, basicConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.test").
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanDefault).
						HasWarehouseString("").
						HasScheduleString("").
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfterLen(0).
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.test").
						//HasCreatedOnNotEmpty(),
						HasName(id.Name()).
						//HasId(id.FullyQualifiedName()). // TODO: not empty
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner("ACCOUNTADMIN"). // TODO: Current role
						HasComment("").
						HasWarehouse("").
						HasSchedule("").
						//HasPredecessors(nil). // TODO:
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")). // TODO: *sdk.AOI
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOnNotEmpty().
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget(""),
					//HasTaskRelations(sdk.TaskRelations{}). // TODO:
				),
			},
		},
	})
}

func TestAcc_Task_AllParameters(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
}

// TODO: Test other paths (alter finalize, after, itd)

//type (
//	AccTaskTestSettings struct {
//		DatabaseName  string
//		WarehouseName string
//		RootTask      *TaskSettings
//		ChildTask     *TaskSettings
//		SoloTask      *TaskSettings
//	}
//
//	TaskSettings struct {
//		Name              string
//		Enabled           bool
//		Schema            string
//		SQL               string
//		Schedule          string
//		Comment           string
//		When              string
//		SessionParams     map[string]string
//		UserTaskTimeoutMs int64
//	}
//)
//
//var (
//	rootname  = acc.TestClient().Ids.AlphaContaining("_root_task")
//	rootId    = sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, rootname)
//	childname = acc.TestClient().Ids.AlphaContaining("_child_task")
//	childId   = sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, childname)
//	soloname  = acc.TestClient().Ids.AlphaContaining("_standalone_task")
//
//	initialState = &AccTaskTestSettings{ //nolint
//		WarehouseName: acc.TestWarehouseName,
//		DatabaseName:  acc.TestDatabaseName,
//		RootTask: &TaskSettings{
//			Name:              rootname,
//			Schema:            acc.TestSchemaName,
//			SQL:               "SHOW FUNCTIONS",
//			Enabled:           true,
//			Schedule:          "5 MINUTE",
//			UserTaskTimeoutMs: 1800000,
//			SessionParams: map[string]string{
//				string(sdk.SessionParameterLockTimeout):      "1000",
//				string(sdk.SessionParameterStrictJSONOutput): "true",
//			},
//		},
//
//		ChildTask: &TaskSettings{
//			Name:    childname,
//			SQL:     "SELECT 1",
//			Enabled: false,
//			Comment: "initial state",
//		},
//
//		SoloTask: &TaskSettings{
//			Name:    soloname,
//			Schema:  acc.TestSchemaName,
//			SQL:     "SELECT 1",
//			When:    "TRUE",
//			Enabled: false,
//		},
//	}
//
//	// Enables the Child and changes the SQL.
//	stepOne = &AccTaskTestSettings{ //nolint
//		WarehouseName: acc.TestWarehouseName,
//		DatabaseName:  acc.TestDatabaseName,
//		RootTask: &TaskSettings{
//			Name:              rootname,
//			Schema:            acc.TestSchemaName,
//			SQL:               "SHOW FUNCTIONS",
//			Enabled:           true,
//			Schedule:          "5 MINUTE",
//			UserTaskTimeoutMs: 1800000,
//			SessionParams: map[string]string{
//				string(sdk.SessionParameterLockTimeout):      "1000",
//				string(sdk.SessionParameterStrictJSONOutput): "true",
//			},
//		},
//
//		ChildTask: &TaskSettings{
//			Name:    childname,
//			SQL:     "SELECT *",
//			Enabled: true,
//			Comment: "secondary state",
//		},
//
//		SoloTask: &TaskSettings{
//			Name:    soloname,
//			Schema:  acc.TestSchemaName,
//			SQL:     "SELECT *",
//			When:    "TRUE",
//			Enabled: true,
//			SessionParams: map[string]string{
//				string(sdk.SessionParameterTimestampInputFormat): "YYYY-MM-DD HH24",
//			},
//			Schedule:          "5 MINUTE",
//			UserTaskTimeoutMs: 1800000,
//		},
//	}
//
//	// Changes Root Schedule and SQL.
//	stepTwo = &AccTaskTestSettings{ //nolint
//		WarehouseName: acc.TestWarehouseName,
//		DatabaseName:  acc.TestDatabaseName,
//		RootTask: &TaskSettings{
//			Name:              rootname,
//			Schema:            acc.TestSchemaName,
//			SQL:               "SHOW TABLES",
//			Enabled:           true,
//			Schedule:          "15 MINUTE",
//			UserTaskTimeoutMs: 1800000,
//			SessionParams: map[string]string{
//				string(sdk.SessionParameterLockTimeout):      "1000",
//				string(sdk.SessionParameterStrictJSONOutput): "true",
//			},
//		},
//
//		ChildTask: &TaskSettings{
//			Name:    childname,
//			SQL:     "SELECT 1",
//			Enabled: true,
//			Comment: "third state",
//		},
//
//		SoloTask: &TaskSettings{
//			Name:              soloname,
//			Schema:            acc.TestSchemaName,
//			SQL:               "SELECT *",
//			When:              "FALSE",
//			Enabled:           true,
//			Schedule:          "15 MINUTE",
//			UserTaskTimeoutMs: 900000,
//		},
//	}
//
//	stepThree = &AccTaskTestSettings{ //nolint
//		WarehouseName: acc.TestWarehouseName,
//		DatabaseName:  acc.TestDatabaseName,
//
//		RootTask: &TaskSettings{
//			Name:              rootname,
//			Schema:            acc.TestSchemaName,
//			SQL:               "SHOW FUNCTIONS",
//			Enabled:           false,
//			Schedule:          "5 MINUTE",
//			UserTaskTimeoutMs: 1800000,
//			// Changes session params: one is updated, one is removed, one is added
//			SessionParams: map[string]string{
//				string(sdk.SessionParameterLockTimeout):         "2000",
//				string(sdk.SessionParameterMultiStatementCount): "5",
//			},
//		},
//
//		ChildTask: &TaskSettings{
//			Name:    childname,
//			SQL:     "SELECT 1",
//			Enabled: false,
//			Comment: "reset",
//		},
//
//		SoloTask: &TaskSettings{
//			Name:    soloname,
//			Schema:  acc.TestSchemaName,
//			SQL:     "SELECT 1",
//			When:    "TRUE",
//			Enabled: true,
//			SessionParams: map[string]string{
//				string(sdk.SessionParameterTimestampInputFormat): "YYYY-MM-DD HH24",
//			},
//			Schedule:          "5 MINUTE",
//			UserTaskTimeoutMs: 0,
//		},
//	}
//)

//func TestAcc_Task(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
//		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
//			tfversion.RequireAbove(tfversion.Version1_5_0),
//		},
//		PreCheck:     func() { acc.TestAccPreCheck(t) },
//		CheckDestroy: acc.CheckDestroy(t, resources.Task),
//		Steps: []resource.TestStep{
//			{
//				Config: taskConfig(initialState),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "enabled", "true"),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "enabled", "false"),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "fully_qualified_name", rootId.FullyQualifiedName()),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "fully_qualified_name", childId.FullyQualifiedName()),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", acc.TestDatabaseName),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", acc.TestDatabaseName),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", acc.TestSchemaName),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", acc.TestSchemaName),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql_statement", initialState.RootTask.SQL),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql_statement", initialState.ChildTask.SQL),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "after.0", rootname),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", initialState.ChildTask.Comment),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "schedule", initialState.RootTask.Schedule),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "schedule", initialState.ChildTask.Schedule),
//					checkInt64("snowflake_task.root_task", "user_task_timeout_ms", initialState.RootTask.UserTaskTimeoutMs),
//					resource.TestCheckNoResourceAttr("snowflake_task.solo_task", "user_task_timeout_ms"),
//					checkInt64("snowflake_task.root_task", "session_parameters.LOCK_TIMEOUT", 1000),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "session_parameters.STRICT_JSON_OUTPUT", "true"),
//					resource.TestCheckNoResourceAttr("snowflake_task.root_task", "session_parameters.MULTI_STATEMENT_COUNT"),
//				),
//			},
//			{
//				Config: taskConfig(stepOne),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "enabled", "true"),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "enabled", "true"),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "fully_qualified_name", rootId.FullyQualifiedName()),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "fully_qualified_name", childId.FullyQualifiedName()),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", acc.TestDatabaseName),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", acc.TestDatabaseName),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", acc.TestSchemaName),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", acc.TestSchemaName),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql_statement", stepOne.RootTask.SQL),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql_statement", stepOne.ChildTask.SQL),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", stepOne.ChildTask.Comment),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "schedule", stepOne.RootTask.Schedule),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "schedule", stepOne.ChildTask.Schedule),
//					checkInt64("snowflake_task.root_task", "user_task_timeout_ms", stepOne.RootTask.UserTaskTimeoutMs),
//					checkInt64("snowflake_task.solo_task", "user_task_timeout_ms", stepOne.SoloTask.UserTaskTimeoutMs),
//					checkInt64("snowflake_task.root_task", "session_parameters.LOCK_TIMEOUT", 1000),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "session_parameters.STRICT_JSON_OUTPUT", "true"),
//					resource.TestCheckNoResourceAttr("snowflake_task.root_task", "session_parameters.MULTI_STATEMENT_COUNT"),
//				),
//			},
//			{
//				Config: taskConfig(stepTwo),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "enabled", "true"),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "enabled", "true"),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "fully_qualified_name", rootId.FullyQualifiedName()),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "fully_qualified_name", childId.FullyQualifiedName()),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", acc.TestDatabaseName),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", acc.TestDatabaseName),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", acc.TestSchemaName),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", acc.TestSchemaName),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql_statement", stepTwo.RootTask.SQL),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql_statement", stepTwo.ChildTask.SQL),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", stepTwo.ChildTask.Comment),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "schedule", stepTwo.RootTask.Schedule),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "schedule", stepTwo.ChildTask.Schedule),
//					checkInt64("snowflake_task.root_task", "user_task_timeout_ms", stepTwo.RootTask.UserTaskTimeoutMs),
//					checkInt64("snowflake_task.solo_task", "user_task_timeout_ms", stepTwo.SoloTask.UserTaskTimeoutMs),
//					checkInt64("snowflake_task.root_task", "session_parameters.LOCK_TIMEOUT", 1000),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "session_parameters.STRICT_JSON_OUTPUT", "true"),
//					resource.TestCheckNoResourceAttr("snowflake_task.root_task", "session_parameters.MULTI_STATEMENT_COUNT"),
//				),
//			},
//			{
//				Config: taskConfig(stepThree),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "enabled", "false"),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "enabled", "false"),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "fully_qualified_name", rootId.FullyQualifiedName()),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "fully_qualified_name", childId.FullyQualifiedName()),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", acc.TestDatabaseName),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", acc.TestDatabaseName),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", acc.TestSchemaName),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", acc.TestSchemaName),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql_statement", stepThree.RootTask.SQL),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql_statement", stepThree.ChildTask.SQL),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", stepThree.ChildTask.Comment),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "schedule", stepThree.RootTask.Schedule),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "schedule", stepThree.ChildTask.Schedule),
//					checkInt64("snowflake_task.root_task", "user_task_timeout_ms", stepThree.RootTask.UserTaskTimeoutMs),
//					checkInt64("snowflake_task.solo_task", "user_task_timeout_ms", stepThree.SoloTask.UserTaskTimeoutMs),
//					checkInt64("snowflake_task.root_task", "session_parameters.LOCK_TIMEOUT", 2000),
//					resource.TestCheckNoResourceAttr("snowflake_task.root_task", "session_parameters.STRICT_JSON_OUTPUT"),
//					checkInt64("snowflake_task.root_task", "session_parameters.MULTI_STATEMENT_COUNT", 5),
//				),
//			},
//			{
//				Config: taskConfig(initialState),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "enabled", "true"),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "enabled", "false"),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "fully_qualified_name", rootId.FullyQualifiedName()),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "fully_qualified_name", childId.FullyQualifiedName()),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", acc.TestDatabaseName),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", acc.TestDatabaseName),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", acc.TestSchemaName),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", acc.TestSchemaName),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql_statement", initialState.RootTask.SQL),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql_statement", initialState.ChildTask.SQL),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", initialState.ChildTask.Comment),
//					checkInt64("snowflake_task.root_task", "user_task_timeout_ms", stepOne.RootTask.UserTaskTimeoutMs),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "schedule", initialState.RootTask.Schedule),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "schedule", initialState.ChildTask.Schedule),
//					// Terraform SDK is not able to differentiate if the
//					// attribute has deleted or set to zero value.
//					// ResourceData.GetChange returns the zero value of defined
//					// type in schema as new the value. Provider handles 0 for
//					// `user_task_timeout_ms` by unsetting the
//					// USER_TASK_TIMEOUT_MS session variable.
//					checkInt64("snowflake_task.solo_task", "user_task_timeout_ms", initialState.ChildTask.UserTaskTimeoutMs),
//					checkInt64("snowflake_task.root_task", "session_parameters.LOCK_TIMEOUT", 1000),
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "session_parameters.STRICT_JSON_OUTPUT", "true"),
//					resource.TestCheckNoResourceAttr("snowflake_task.root_task", "session_parameters.MULTI_STATEMENT_COUNT"),
//				),
//			},
//		},
//	})
//}

//func taskConfig(settings *AccTaskTestSettings) string { //nolint
//	config, err := template.New("task_acceptance_test_config").Parse(`
//resource "snowflake_warehouse" "wh" {
//	name = "{{ .WarehouseName }}-{{ .RootTask.Name }}"
//}
//resource "snowflake_task" "root_task" {
//	name     	  = "{{ .RootTask.Name }}"
//	database  	  = "{{ .DatabaseName }}"
//	schema   	  = "{{ .RootTask.Schema }}"
//	warehouse 	  = "${snowflake_warehouse.wh.name}"
//	sql_statement = "{{ .RootTask.SQL }}"
//	enabled  	  = {{ .RootTask.Enabled }}
//	schedule 	  = "{{ .RootTask.Schedule }}"
//	{{ if .RootTask.UserTaskTimeoutMs }}
//	user_task_timeout_ms = {{ .RootTask.UserTaskTimeoutMs }}
//	{{- end }}
//
//	{{ if .RootTask.SessionParams }}
//	session_parameters = {
//	{{ range $key, $value := .RootTask.SessionParams}}
//        {{ $key }} = "{{ $value }}",
//	{{- end }}
//	}
//	{{- end }}
//}
//resource "snowflake_task" "child_task" {
//	name     	  = "{{ .ChildTask.Name }}"
//	database   	  = snowflake_task.root_task.database
//	schema    	  = snowflake_task.root_task.schema
//	warehouse 	  = snowflake_task.root_task.warehouse
//	sql_statement = "{{ .ChildTask.SQL }}"
//	enabled  	  = {{ .ChildTask.Enabled }}
//	after    	  = [snowflake_task.root_task.name]
//	comment 	  = "{{ .ChildTask.Comment }}"
//	{{ if .ChildTask.UserTaskTimeoutMs }}
//	user_task_timeout_ms = {{ .ChildTask.UserTaskTimeoutMs }}
//	{{- end }}
//
//	{{ if .ChildTask.SessionParams }}
//	session_parameters = {
//	{{ range $key, $value := .ChildTask.SessionParams}}
//        {{ $key }} = "{{ $value }}",
//	{{- end }}
//	}
//	{{- end }}
//}
//resource "snowflake_task" "solo_task" {
//	name     	  = "{{ .SoloTask.Name }}"
//	database  	  = "{{ .DatabaseName }}"
//	schema    	  = "{{ .SoloTask.Schema }}"
//	warehouse 	  = "{{ .WarehouseName }}"
//	sql_statement = "{{ .SoloTask.SQL }}"
//	enabled  	  = {{ .SoloTask.Enabled }}
//	when     	  = "{{ .SoloTask.When }}"
//	{{ if .SoloTask.Schedule }}
//	schedule    = "{{ .SoloTask.Schedule }}"
//	{{- end }}
//
//	{{ if .SoloTask.UserTaskTimeoutMs }}
//	user_task_timeout_ms = {{ .SoloTask.UserTaskTimeoutMs }}
//	{{- end }}
//
//	{{ if .SoloTask.SessionParams }}
//	session_parameters = {
//	{{ range $key, $value :=  .SoloTask.SessionParams}}
//        {{ $key }} = "{{ $value }}",
//	{{- end }}
//	}
//	{{- end }}
//}
//	`)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	var result bytes.Buffer
//	config.Execute(&result, settings) //nolint
//
//	return result.String()
//}

/*
todo: this test is failing due to error message below. Need to figure out why this is happening
=== RUN   TestAcc_Task_Managed

	task_acceptance_test.go:371: Step 2/4 error: Error running apply: exit status 1

	    Error: error updating warehouse on task "terraform_test_database"."terraform_test_schema"."tst-terraform-DBMPMESYJB" err = 091083 (42601): Nonexistent warehouse terraform_test_warehouse-tst-terraform-DBMPMESYJB was specified.

	      with snowflake_task.managed_task,
	      on terraform_plugin_test.tf line 7, in resource "snowflake_task" "managed_task":
	       7: resource "snowflake_task" "managed_task" {


	func TestAcc_Task_Managed(t *testing.T) {
		accName := acc.TestClient().Ids.Alpha()
		resource.Test(t, resource.TestCase{
					ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
			PreCheck:     func() { acc.TestAccPreCheck(t) },
			CheckDestroy: acc.CheckDestroy(t, resources.Task),
			Steps: []resource.TestStep{
				{
					Config: taskConfigManaged1(accName, acc.TestDatabaseName, acc.TestSchemaName),
					Check: resource.ComposeTestCheckFunc(
						checkBool("snowflake_task.managed_task", "enabled", true),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "database", acc.TestDatabaseName),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "schema", acc.TestSchemaName),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "sql_statement", "SELECT 1"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "schedule", "5 MINUTE"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "user_task_managed_initial_warehouse_size", "XSMALL"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task_no_init", "user_task_managed_initial_warehouse_size", ""),
						resource.TestCheckResourceAttr("snowflake_task.managed_task_no_init", "session_parameters.TIMESTAMP_INPUT_FORMAT", "YYYY-MM-DD HH24"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "warehouse", ""),
					),
				},
				{
					Config: taskConfigManaged2(accName, acc.TestDatabaseName, acc.TestSchemaName, acc.TestWarehouseName),
					Check: resource.ComposeTestCheckFunc(
						checkBool("snowflake_task.managed_task", "enabled", true),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "database", acc.TestDatabaseName),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "schema", acc.TestSchemaName),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "sql_statement", "SELECT 1"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "schedule", "5 MINUTE"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "user_task_managed_initial_warehouse_size", ""),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "warehouse", fmt.Sprintf("%s-%s", acc.TestWarehouseName, accName)),
					),
				},
				{
					Config: taskConfigManaged1(accName, acc.TestDatabaseName, acc.TestSchemaName),
					Check: resource.ComposeTestCheckFunc(
						checkBool("snowflake_task.managed_task", "enabled", true),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "database", acc.TestDatabaseName),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "schema", acc.TestSchemaName),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "sql_statement", "SELECT 1"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "schedule", "5 MINUTE"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task_no_init", "session_parameters.TIMESTAMP_INPUT_FORMAT", "YYYY-MM-DD HH24"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task_no_init", "user_task_managed_initial_warehouse_size", ""),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "warehouse", ""),
					),
				},
				{
					Config: taskConfigManaged3(accName, acc.TestDatabaseName, acc.TestSchemaName),
					Check: resource.ComposeTestCheckFunc(
						checkBool("snowflake_task.managed_task", "enabled", true),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "database", acc.TestDatabaseName),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "schema", acc.TestSchemaName),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "sql_statement", "SELECT 1"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "schedule", "5 MINUTE"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "user_task_managed_initial_warehouse_size", "SMALL"),
						resource.TestCheckResourceAttr("snowflake_task.managed_task", "warehouse", ""),
					),
				},
			},
		})
	}
*/
func taskConfigManaged1(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_task" "managed_task" {
	name     	                             = "%s"
	database  	                             = "%s"
	schema    	                             = "%s"
	sql_statement                            = "SELECT 1"
	enabled  	                             = true
	schedule                                 = "5 MINUTE"
    user_task_managed_initial_warehouse_size = "XSMALL"
}
resource "snowflake_task" "managed_task_no_init" {
	name     	  = "%s_no_init"
	database  	  = "%s"
	schema    	  = "%s"
	sql_statement = "SELECT 1"
	enabled  	  = true
	schedule      = "5 MINUTE"
	session_parameters = {
		TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24",
	}
}

`
	return fmt.Sprintf(s, name, databaseName, schemaName, name, databaseName, schemaName)
}

func taskConfigManaged2(name, databaseName, schemaName, warehouseName string) string {
	s := `
resource "snowflake_warehouse" "wh" {
	name = "%s-%s"
}

resource "snowflake_task" "managed_task" {
	name     	  = "%s"
	database  	  = "%s"
	schema    	  = "%s"
	sql_statement = "SELECT 1"
	enabled  	  = true
	schedule      = "5 MINUTE"
	warehouse     = snowflake_warehouse.wh.name
}
`
	return fmt.Sprintf(s, warehouseName, name, name, databaseName, schemaName)
}

func taskConfigManaged3(name, databaseName, schemaName string) string {
	s := `
resource "snowflake_task" "managed_task" {
	name     	                             = "%s"
	database  	                             = "%s"
	schema    	                             = "%s"
	sql_statement                            = "SELECT 1"
	enabled  	                             = true
	schedule                                 = "5 MINUTE"
    user_task_managed_initial_warehouse_size = "SMALL"
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName)
}

func TestAcc_Task_SwitchScheduled(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()
	taskRootName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: taskConfigManagedScheduled(accName, taskRootName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task.test_task", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schedule", "5 MINUTE"),
					resource.TestCheckResourceAttr("snowflake_task.test_task_root", "suspend_task_after_num_failures", "1"),
				),
			},
			{
				Config: taskConfigManagedScheduled2(accName, taskRootName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task.test_task", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schedule", ""),
					resource.TestCheckResourceAttr("snowflake_task.test_task_root", "suspend_task_after_num_failures", "2"),
				),
			},
			{
				Config: taskConfigManagedScheduled(accName, taskRootName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task.test_task", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schedule", "5 MINUTE"),
					resource.TestCheckResourceAttr("snowflake_task.test_task_root", "suspend_task_after_num_failures", "1"),
				),
			},
			{
				Config: taskConfigManagedScheduled3(accName, taskRootName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task.test_task", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schedule", ""),
					resource.TestCheckResourceAttr("snowflake_task.test_task_root", "suspend_task_after_num_failures", "0"),
				),
			},
		},
	})
}

func taskConfigManagedScheduled(name string, taskRootName string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_task" "test_task_root" {
	name     	                    = "%[1]s"
	database  	                    = "%[2]s"
	schema    	                    = "%[3]s"
	sql_statement                   = "SELECT 1"
	enabled  	                    = true
	schedule                        = "5 MINUTE"
    suspend_task_after_num_failures = 1
}

resource "snowflake_task" "test_task" {
	depends_on = [snowflake_task.test_task_root]
	name     	  = "%[4]s"
	database  	  = "%[2]s"
	schema    	  = "%[3]s"
	sql_statement = "SELECT 1"
	enabled  	  = true
	schedule      = "5 MINUTE"
}
`, taskRootName, databaseName, schemaName, name)
}

func taskConfigManagedScheduled2(name string, taskRootName string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_task" "test_task_root" {
	name     	                    = "%[1]s"
	database  	                    = "%[2]s"
	schema    	                    = "%[3]s"
	sql_statement                   = "SELECT 1"
	enabled  	                    = true
	schedule                        = "5 MINUTE"
    suspend_task_after_num_failures = 2
}

resource "snowflake_task" "test_task" {
	name     	  = "%[4]s"
	database  	  = "%[2]s"
	schema    	  = "%[3]s"
	sql_statement = "SELECT 1"
	enabled  	  = true
	after         = [snowflake_task.test_task_root.name]
}
`, taskRootName, databaseName, schemaName, name)
}

func taskConfigManagedScheduled3(name string, taskRootName string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_task" "test_task_root" {
	name     	  = "%s"
	database  	  = "%s"
	schema    	  = "%s"
	sql_statement = "SELECT 1"
	enabled  	  = false
	schedule      = "5 MINUTE"
}

resource "snowflake_task" "test_task" {
	name     	  = "%s"
	database  	  = "%s"
	schema    	  = "%s"
	sql_statement = "SELECT 1"
	enabled  	  = false
	after         = [snowflake_task.test_task_root.name]
}
`
	return fmt.Sprintf(s, taskRootName, databaseName, schemaName, name, databaseName, schemaName)
}

func checkInt64(name, key string, value int64) func(*terraform.State) error {
	return func(state *terraform.State) error {
		return resource.TestCheckResourceAttr(name, key, fmt.Sprintf("%v", value))(state)
	}
}

//func TestAcc_Task_issue2207(t *testing.T) {
//	prefix := acc.TestClient().Ids.Alpha()
//	rootName := prefix + "_root_task"
//	childName := prefix + "_child_task"
//
//	m := func() map[string]config.Variable {
//		return map[string]config.Variable{
//			"root_name":  config.StringVariable(rootName),
//			"database":   config.StringVariable(acc.TestDatabaseName),
//			"schema":     config.StringVariable(acc.TestSchemaName),
//			"warehouse":  config.StringVariable(acc.TestWarehouseName),
//			"child_name": config.StringVariable(childName),
//			"comment":    config.StringVariable("abc"),
//		}
//	}
//	m2 := m()
//	m2["comment"] = config.StringVariable("def")
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
//		PreCheck:                 func() { acc.TestAccPreCheck(t) },
//		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
//			tfversion.RequireAbove(tfversion.Version1_5_0),
//		},
//		CheckDestroy: acc.CheckDestroy(t, resources.Task),
//		Steps: []resource.TestStep{
//			{
//				ConfigDirectory: config.TestStepDirectory(),
//				ConfigVariables: m(),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "enabled", "true"),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "enabled", "true"),
//				),
//				ConfigPlanChecks: resource.ConfigPlanChecks{
//					PostApplyPostRefresh: []plancheck.PlanCheck{
//						plancheck.ExpectEmptyPlan(),
//					},
//				},
//			},
//			// change comment
//			{
//				ConfigDirectory: acc.ConfigurationSameAsStepN(1),
//				ConfigVariables: m2,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_task.root_task", "enabled", "true"),
//					resource.TestCheckResourceAttr("snowflake_task.child_task", "enabled", "true"),
//				),
//			},
//		},
//	})
//}
//
//func TestAcc_Task_issue2036(t *testing.T) {
//	name := acc.TestClient().Ids.Alpha()
//
//	m := func() map[string]config.Variable {
//		return map[string]config.Variable{
//			"name":      config.StringVariable(name),
//			"database":  config.StringVariable(acc.TestDatabaseName),
//			"schema":    config.StringVariable(acc.TestSchemaName),
//			"warehouse": config.StringVariable(acc.TestWarehouseName),
//		}
//	}
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
//		PreCheck:                 func() { acc.TestAccPreCheck(t) },
//		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
//			tfversion.RequireAbove(tfversion.Version1_5_0),
//		},
//		CheckDestroy: acc.CheckDestroy(t, resources.Task),
//		Steps: []resource.TestStep{
//			// create without when
//			{
//				ConfigDirectory: config.TestStepDirectory(),
//				ConfigVariables: m(),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_task.test_task", "enabled", "true"),
//					resource.TestCheckResourceAttr("snowflake_task.test_task", "when", ""),
//				),
//			},
//			// add when
//			{
//				ConfigDirectory: config.TestStepDirectory(),
//				ConfigVariables: m(),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_task.test_task", "enabled", "true"),
//					resource.TestCheckResourceAttr("snowflake_task.test_task", "when", "TRUE"),
//				),
//			},
//			// remove when
//			{
//				ConfigDirectory: acc.ConfigurationSameAsStepN(1),
//				ConfigVariables: m(),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("snowflake_task.test_task", "enabled", "true"),
//					resource.TestCheckResourceAttr("snowflake_task.test_task", "when", ""),
//				),
//			},
//		},
//	})
//}
