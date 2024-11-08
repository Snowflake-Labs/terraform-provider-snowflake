package resources_test

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	configvariable "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO(SNOW-1348116 - next pr): More tests for complicated DAGs
// TODO(SNOW-1348116 - next pr): Test for stored procedures passed to sql_statement (decide on name)
// TODO(SNOW-1348116 - next pr): More test with external changes

func TestAcc_Task_Basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	currentRole := acc.TestClient().Context.CurrentRole(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	configModel := model.TaskWithId("test", id, false, statement)

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
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfterIdsInOrder().
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, configModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment("").
						HasWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasNoSchedule().
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")).
						HasLastCommittedOn("").
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			{
				ResourceName: configModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedTaskResource(t, helpers.EncodeResourceIdentifier(id)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
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

	currentRole := acc.TestClient().Context.CurrentRole(t)

	errorNotificationIntegration, errorNotificationIntegrationCleanup := acc.TestClient().NotificationIntegration.Create(t)
	t.Cleanup(errorNotificationIntegrationCleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	taskConfig := `$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`
	// We have to do three $ at the beginning because Terraform will remove one $.
	// It's because `${` is a special pattern, and it's escaped by `$${`.
	expectedTaskConfig := strings.ReplaceAll(taskConfig, "$", "")
	comment := random.Comment()
	condition := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
	configModel := model.TaskWithId("test", id, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(10).
		WithConfigValue(configvariable.StringVariable(taskConfig)).
		WithAllowOverlappingExecution(r.BooleanTrue).
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasWarehouseString(acc.TestClient().Ids.WarehouseId().Name()).
						HasScheduleMinutes(10).
						HasConfigString(expectedTaskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasCommentString(comment).
						HasFinalizeString("").
						HasNoAfter().
						HasWhenString(condition).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, configModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment(comment).
						HasWarehouse(acc.TestClient().Ids.WarehouseId()).
						HasScheduleMinutes(10).
						HasPredecessors().
						HasState(sdk.TaskStateStarted).
						HasDefinition(statement).
						HasCondition(condition).
						HasAllowOverlappingExecution(true).
						HasErrorIntegration(errorNotificationIntegration.ID()).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig(expectedTaskConfig).
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			{
				ResourceName:    "snowflake_task.test",
				ImportState:     true,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedTaskResource(t, helpers.EncodeResourceIdentifier(id)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasWarehouseString(acc.TestClient().Ids.WarehouseId().Name()).
						HasScheduleMinutes(10).
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

	currentRole := acc.TestClient().Context.CurrentRole(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	basicConfigModel := model.TaskWithId("test", id, false, statement)

	// TODO(SNOW-1736173): New warehouse created, because the common one has lower-case letters that won't work
	warehouse, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	errorNotificationIntegration, errorNotificationIntegrationCleanup := acc.TestClient().NotificationIntegration.Create(t)
	t.Cleanup(errorNotificationIntegrationCleanup)

	taskConfig := `$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`
	// We have to do three $ at the beginning because Terraform will remove one $.
	// It's because `${` is a special pattern, and it's escaped by `$${`.
	expectedTaskConfig := strings.ReplaceAll(taskConfig, "$", "")
	comment := random.Comment()
	condition := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
	completeConfigModel := model.TaskWithId("test", id, true, statement).
		WithWarehouse(warehouse.ID().Name()).
		WithScheduleMinutes(5).
		WithConfigValue(configvariable.StringVariable(taskConfig)).
		WithAllowOverlappingExecution(r.BooleanTrue).
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
					resourceassert.TaskResource(t, basicConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfterIdsInOrder().
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, basicConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment("").
						HasWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasNoSchedule().
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")).
						HasLastCommittedOn("").
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
			// Set
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, completeConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, completeConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasWarehouseString(warehouse.ID().Name()).
						HasScheduleMinutes(5).
						HasConfigString(expectedTaskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasCommentString(comment).
						HasFinalizeString("").
						HasAfterIdsInOrder().
						HasWhenString(condition).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, completeConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasWarehouse(warehouse.ID()).
						HasComment(comment).
						HasScheduleMinutes(5).
						HasPredecessors().
						HasState(sdk.TaskStateStarted).
						HasDefinition(statement).
						HasCondition(condition).
						HasAllowOverlappingExecution(true).
						HasErrorIntegration(errorNotificationIntegration.ID()).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig(expectedTaskConfig).
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
			// Unset
			{
				Config: config.FromModel(t, basicConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, basicConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfterIdsInOrder().
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, basicConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment("").
						HasWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasNoSchedule().
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOnNotEmpty().
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
		},
	})
}

func TestAcc_Task_CronAndMinutes(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	minutes := 5
	cron := "*/5 * * * * UTC"
	configModelWithoutSchedule := model.TaskWithId("test", id, false, "SELECT 1")
	configModelWithMinutes := model.TaskWithId("test", id, true, "SELECT 1").WithScheduleMinutes(minutes)
	configModelWithCron := model.TaskWithId("test", id, true, "SELECT 1").WithScheduleCron(cron)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			// create with minutes
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithMinutes),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelWithMinutes.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(minutes).
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithMinutes.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleMinutes(minutes),
				),
			},
			// Unset schedule (from minutes)
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithoutSchedule),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelWithoutSchedule.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasNoScheduleSet().
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithoutSchedule.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasNoSchedule(),
				),
			},
			// Create with cron
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithCron),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelWithCron.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasScheduleCron(cron).
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithCron.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleCron(cron),
				),
			},
			// Change to minutes
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithMinutes),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelWithMinutes.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(minutes).
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithMinutes.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleMinutes(minutes),
				),
			},
			// Change back to cron
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithCron),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelWithCron.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasScheduleCron(cron).
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithCron.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleCron(cron),
				),
			},
			// Unset schedule (from cron)
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithoutSchedule),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelWithoutSchedule.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasNoScheduleSet().
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithoutSchedule.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasNoSchedule(),
				),
			},
		},
	})
}

func TestAcc_Task_AllParameters(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	configModel := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(5)
	configModelWithAllParametersSet := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(5).
		WithSuspendTaskAfterNumFailures(15).
		WithTaskAutoRetryAttempts(15).
		WithUserTaskManagedInitialWarehouseSizeEnum(sdk.WarehouseSizeXSmall).
		WithUserTaskMinimumTriggerIntervalInSeconds(30).
		WithUserTaskTimeoutMs(1000).
		WithAbortDetachedQuery(true).
		WithAutocommit(false).
		WithBinaryInputFormatEnum(sdk.BinaryInputFormatUTF8).
		WithBinaryOutputFormatEnum(sdk.BinaryOutputFormatBase64).
		WithClientMemoryLimit(1024).
		WithClientMetadataRequestUseConnectionCtx(true).
		WithClientPrefetchThreads(2).
		WithClientResultChunkSize(48).
		WithClientResultColumnCaseInsensitive(true).
		WithClientSessionKeepAlive(true).
		WithClientSessionKeepAliveHeartbeatFrequency(2400).
		WithClientTimestampTypeMappingEnum(sdk.ClientTimestampTypeMappingNtz).
		WithDateInputFormat("YYYY-MM-DD").
		WithDateOutputFormat("YY-MM-DD").
		WithEnableUnloadPhysicalTypeOptimization(false).
		WithErrorOnNondeterministicMerge(false).
		WithErrorOnNondeterministicUpdate(true).
		WithGeographyOutputFormatEnum(sdk.GeographyOutputFormatWKB).
		WithGeometryOutputFormatEnum(sdk.GeometryOutputFormatWKB).
		WithJdbcUseSessionTimezone(false).
		WithJsonIndent(4).
		WithLockTimeout(21222).
		WithLogLevelEnum(sdk.LogLevelError).
		WithMultiStatementCount(0).
		WithNoorderSequenceAsDefault(false).
		WithOdbcTreatDecimalAsInt(true).
		WithQueryTag("some_tag").
		WithQuotedIdentifiersIgnoreCase(true).
		WithRowsPerResultset(2).
		WithS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
		WithSearchPath("$public, $current").
		WithStatementQueuedTimeoutInSeconds(10).
		WithStatementTimeoutInSeconds(10).
		WithStrictJsonOutput(true).
		WithTimestampDayIsAlways24h(true).
		WithTimestampInputFormat("YYYY-MM-DD").
		WithTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampTypeMappingEnum(sdk.TimestampTypeMappingLtz).
		WithTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimezone("Europe/Warsaw").
		WithTimeInputFormat("HH24:MI").
		WithTimeOutputFormat("HH24:MI").
		WithTraceLevelEnum(sdk.TraceLevelOnEvent).
		WithTransactionAbortOnError(true).
		WithTransactionDefaultIsolationLevelEnum(sdk.TransactionDefaultIsolationLevelReadCommitted).
		WithTwoDigitCenturyStart(1980).
		WithUnsupportedDdlActionEnum(sdk.UnsupportedDDLActionFail).
		WithUseCachedResult(false).
		WithWeekOfYearPolicy(1).
		WithWeekStart(1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// create with default values for all the parameters
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				Check: assert.AssertThat(t,
					objectparametersassert.TaskParameters(t, id).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			// import when no parameter set
			{
				ResourceName:    configModel.ResourceReference(),
				ImportState:     true,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				ImportStateCheck: assert.AssertThatImport(t,
					resourceparametersassert.ImportedTaskResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasAllDefaults(),
				),
			},
			// set all parameters
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithAllParametersSet),
				Check: assert.AssertThat(t,
					objectparametersassert.TaskParameters(t, id).
						HasSuspendTaskAfterNumFailures(15).
						HasTaskAutoRetryAttempts(15).
						HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall).
						HasUserTaskMinimumTriggerIntervalInSeconds(30).
						HasUserTaskTimeoutMs(1000).
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelOnEvent).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1),
					resourceparametersassert.TaskResourceParameters(t, configModelWithAllParametersSet.ResourceReference()).
						HasSuspendTaskAfterNumFailures(15).
						HasTaskAutoRetryAttempts(15).
						HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall).
						HasUserTaskMinimumTriggerIntervalInSeconds(30).
						HasUserTaskTimeoutMs(1000).
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelOnEvent).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1),
				),
			},
			// import when all parameters set
			{
				ResourceName:    configModelWithAllParametersSet.ResourceReference(),
				ImportState:     true,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithAllParametersSet),
				ImportStateCheck: assert.AssertThatImport(t,
					resourceparametersassert.ImportedTaskResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasSuspendTaskAfterNumFailures(15).
						HasTaskAutoRetryAttempts(15).
						HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall).
						HasUserTaskMinimumTriggerIntervalInSeconds(30).
						HasUserTaskTimeoutMs(1000).
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelOnEvent).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1),
				),
			},
			// unset all the parameters
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				Check: assert.AssertThat(t,
					objectparametersassert.TaskParameters(t, id).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
		},
	})
}

func TestAcc_Task_Enabled(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	configModelEnabled := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(5)
	configModelDisabled := model.TaskWithId("test", id, false, statement).
		WithScheduleMinutes(5)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelDisabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelDisabled.ResourceReference()).
						HasStartedString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, configModelDisabled.ResourceReference()).
						HasState(sdk.TaskStateSuspended),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelEnabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelEnabled.ResourceReference()).
						HasStartedString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, configModelEnabled.ResourceReference()).
						HasState(sdk.TaskStateStarted),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelDisabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelDisabled.ResourceReference()).
						HasStartedString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, configModelDisabled.ResourceReference()).
						HasState(sdk.TaskStateSuspended),
				),
			},
		},
	})
}

func TestAcc_Task_ConvertStandaloneTaskToSubtask(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	id2 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"

	firstTaskStandaloneModel := model.TaskWithId("main_task", id, true, statement).
		WithScheduleMinutes(5).
		WithSuspendTaskAfterNumFailures(1)
	secondTaskStandaloneModel := model.TaskWithId("second_task", id2, true, statement).
		WithScheduleMinutes(5)

	rootTaskModel := model.TaskWithId("main_task", id, true, statement).
		WithScheduleMinutes(5).
		WithSuspendTaskAfterNumFailures(2)
	childTaskModel := model.TaskWithId("second_task", id2, true, statement).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(id.FullyQualifiedName())))
	childTaskModel.SetDependsOn(rootTaskModel.ResourceReference())

	firstTaskStandaloneModelDisabled := model.TaskWithId("root", id, false, statement).
		WithScheduleMinutes(5)
	secondTaskStandaloneModelDisabled := model.TaskWithId("child", id2, false, statement).
		WithScheduleMinutes(5)
	secondTaskStandaloneModelDisabled.SetDependsOn(firstTaskStandaloneModelDisabled.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", firstTaskStandaloneModel, secondTaskStandaloneModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasScheduleMinutes(5).
						HasStartedString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("1"),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.root").
						HasScheduleMinutes(5).
						HasState(sdk.TaskStateStarted),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasScheduleMinutes(5).
						HasStartedString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasScheduleMinutes(5).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change the second task to run after the first one (creating a DAG)
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskModel, childTaskModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasScheduleMinutes(5).
						HasStartedString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("2"),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.root").
						HasScheduleMinutes(5).
						HasState(sdk.TaskStateStarted),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasAfterIdsInOrder(id).
						HasStartedString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasPredecessors(id).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change tasks in DAG to standalone tasks (disabled to check if resuming/suspending works correctly)
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", firstTaskStandaloneModelDisabled, secondTaskStandaloneModelDisabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasScheduleMinutes(5).
						HasStartedString(r.BooleanFalse).
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.root").
						HasScheduleMinutes(5).
						HasState(sdk.TaskStateSuspended),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasScheduleMinutes(5).
						HasStartedString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasScheduleMinutes(5).
						HasState(sdk.TaskStateSuspended),
				),
			},
		},
	})
}

func TestAcc_Task_ConvertStandaloneTaskToFinalizer(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	rootTaskId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	finalizerTaskId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	firstTaskStandaloneModel := model.TaskWithId("root", rootTaskId, true, statement).
		WithScheduleMinutes(schedule).
		WithSuspendTaskAfterNumFailures(1)
	secondTaskStandaloneModel := model.TaskWithId("child", finalizerTaskId, true, statement).
		WithScheduleMinutes(schedule)

	rootTaskModel := model.TaskWithId("root", rootTaskId, true, statement).
		WithScheduleMinutes(schedule).
		WithSuspendTaskAfterNumFailures(2)
	childTaskModel := model.TaskWithId("child", finalizerTaskId, true, statement).
		WithFinalize(rootTaskId.FullyQualifiedName())
	childTaskModel.SetDependsOn(rootTaskModel.ResourceReference())

	firstTaskStandaloneModelDisabled := model.TaskWithId("root", rootTaskId, false, statement).
		WithScheduleMinutes(schedule)
	secondTaskStandaloneModelDisabled := model.TaskWithId("child", finalizerTaskId, false, statement).
		WithScheduleMinutes(schedule)
	secondTaskStandaloneModelDisabled.SetDependsOn(firstTaskStandaloneModelDisabled.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", firstTaskStandaloneModel, secondTaskStandaloneModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasScheduleMinutes(schedule).
						HasStartedString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("1"),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.root").
						HasScheduleMinutes(schedule).
						HasTaskRelations(sdk.TaskRelations{}).
						HasState(sdk.TaskStateStarted),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasScheduleMinutes(schedule).
						HasStartedString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasScheduleMinutes(schedule).
						HasTaskRelations(sdk.TaskRelations{}).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change the second task to run after the first one (creating a DAG)
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskModel, childTaskModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasScheduleMinutes(schedule).
						HasStartedString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("2"),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.root").
						HasScheduleMinutes(schedule).
						// TODO(SNOW-1348116 - next pr): Create ticket and report; this field in task relations seems to have mixed chances of appearing (needs deeper digging, doesn't affect the resource; could be removed for now)
						// HasTaskRelations(sdk.TaskRelations{FinalizerTask: &finalizerTaskId}).
						HasState(sdk.TaskStateStarted),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasTaskRelations(sdk.TaskRelations{FinalizedRootTask: &rootTaskId}).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change tasks in DAG to standalone tasks (disabled to check if resuming/suspending works correctly)
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", firstTaskStandaloneModelDisabled, secondTaskStandaloneModelDisabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasScheduleMinutes(schedule).
						HasStartedString(r.BooleanFalse).
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.root").
						HasScheduleMinutes(schedule).
						HasTaskRelations(sdk.TaskRelations{}).
						HasState(sdk.TaskStateSuspended),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasScheduleMinutes(schedule).
						HasStartedString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasScheduleMinutes(schedule).
						HasTaskRelations(sdk.TaskRelations{}).
						HasState(sdk.TaskStateSuspended),
				),
			},
		},
	})
}

func TestAcc_Task_SwitchScheduledWithAfter(t *testing.T) {
	rootId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	childId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5
	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithScheduleMinutes(schedule).
		WithSuspendTaskAfterNumFailures(1)
	childTaskConfigModel := model.TaskWithId("child", childId, true, statement).
		WithScheduleMinutes(schedule)

	rootTaskConfigModelAfterSuspendFailuresUpdate := model.TaskWithId("root", rootId, true, statement).
		WithScheduleMinutes(schedule).
		WithSuspendTaskAfterNumFailures(2)
	childTaskConfigModelWithAfter := model.TaskWithId("child", childId, true, statement).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName())))
	childTaskConfigModelWithAfter.SetDependsOn(rootTaskConfigModelAfterSuspendFailuresUpdate.ResourceReference())

	rootTaskConfigModelDisabled := model.TaskWithId("root", rootId, false, statement).
		WithScheduleMinutes(schedule)
	childTaskConfigModelDisabled := model.TaskWithId("child", childId, false, statement).
		WithScheduleMinutes(schedule)
	childTaskConfigModelDisabled.SetDependsOn(rootTaskConfigModelDisabled.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule).
						HasSuspendTaskAfterNumFailuresString("1"),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule).
						HasAfterIdsInOrder().
						HasSuspendTaskAfterNumFailuresString("10"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModelAfterSuspendFailuresUpdate, childTaskConfigModelWithAfter),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule).
						HasSuspendTaskAfterNumFailuresString("2"),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasNoScheduleSet().
						HasAfterIdsInOrder(rootId).
						HasSuspendTaskAfterNumFailuresString("10"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule).
						HasSuspendTaskAfterNumFailuresString("1"),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule).
						HasAfterIdsInOrder().
						HasSuspendTaskAfterNumFailuresString("10"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModelDisabled, childTaskConfigModelDisabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasStartedString(r.BooleanFalse).
						HasScheduleMinutes(schedule).
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanFalse).
						HasScheduleMinutes(schedule).
						HasAfterIdsInOrder().
						HasSuspendTaskAfterNumFailuresString("10"),
				),
			},
		},
	})
}

func TestAcc_Task_WithAfter(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	rootId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	childId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	childTaskConfigModelWithAfter := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName()))).
		WithSqlStatement(statement)

	childTaskConfigModelWithoutAfter := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithAfter),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasAfterIdsInOrder(rootId),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutAfter),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasAfterIdsInOrder(),
				),
			},
		},
	})
}

func TestAcc_Task_WithFinalizer(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	rootId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	childId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	childTaskConfigModelWithFinalizer := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithFinalize(rootId.FullyQualifiedName()).
		WithSqlStatement(statement)

	childTaskConfigModelWithoutFinalizer := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithFinalizer),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(rootId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutFinalizer),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(""),
				),
			},
		},
	})
}

func TestAcc_Task_UpdateFinalizerExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	rootId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	childId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	childTaskConfigModelWithoutFinalizer := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithComment("abc").
		WithSqlStatement(statement)

	childTaskConfigModelWithFinalizer := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithFinalize(rootId.FullyQualifiedName()).
		WithComment("abc").
		WithSqlStatement(statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutFinalizer),
			},
			// Set finalizer externally
			{
				PreConfig: func() {
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithSuspend(true))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSuspend(true))

					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithUnset(*sdk.NewTaskUnsetRequest().WithSchedule(true)))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSetFinalize(rootId))

					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithResume(true))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithResume(true))
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutFinalizer),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(""),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
			// Set finalizer in config
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithFinalizer),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(rootId.FullyQualifiedName()),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{FinalizedRootTask: &rootId}),
				),
			},
			// Unset finalizer externally
			{
				PreConfig: func() {
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithSuspend(true))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSuspend(true))

					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithUnsetFinalize(true))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSet(*sdk.NewTaskSetRequest().WithSchedule(fmt.Sprintf("%d minutes", schedule))))

					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithResume(true))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithResume(true))
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithFinalizer),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(rootId.FullyQualifiedName()),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{FinalizedRootTask: &rootId}),
				),
			},
			// Unset finalizer in config
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutFinalizer),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(""),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
		},
	})
}

func TestAcc_Task_UpdateAfterExternally(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	rootId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	childId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	childTaskConfigModelWithoutAfter := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithComment("abc").
		WithSqlStatement(statement)

	childTaskConfigModelWithAfter := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName()))).
		WithComment("abc").
		WithSqlStatement(statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutAfter),
			},
			// Set after externally
			{
				PreConfig: func() {
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithSuspend(true))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSuspend(true))

					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithUnset(*sdk.NewTaskUnsetRequest().WithSchedule(true)))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithAddAfter([]sdk.SchemaObjectIdentifier{rootId}))

					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithResume(true))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithResume(true))
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutAfter),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasAfterIdsInOrder(),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
			// Set after in config
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithAfter),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasAfterIdsInOrder(rootId),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{Predecessors: []sdk.SchemaObjectIdentifier{rootId}}),
				),
			},
			// Unset after externally
			{
				PreConfig: func() {
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithSuspend(true))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSuspend(true))

					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithRemoveAfter([]sdk.SchemaObjectIdentifier{rootId}))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSet(*sdk.NewTaskSetRequest().WithSchedule(fmt.Sprintf("%d MINUTES", schedule))))

					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithResume(true))
					acc.TestClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithResume(true))
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithAfter),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasAfterIdsInOrder(rootId),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{Predecessors: []sdk.SchemaObjectIdentifier{rootId}}),
				),
			},
			// Unset after in config
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutAfter),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasAfterIdsInOrder(),
					resourceshowoutputassert.TaskShowOutput(t, "snowflake_task.child").
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
		},
	})
}

func TestAcc_Task_issue2207(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	rootId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	childId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	childTaskConfigModel := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName()))).
		WithComment("abc").
		WithSqlStatement(statement)

	childTaskConfigModelWithDifferentComment := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName()))).
		WithComment("def").
		WithSqlStatement(statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasAfterIdsInOrder(rootId).
						HasCommentString("abc"),
				),
			},
			// change comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(childTaskConfigModelWithDifferentComment.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithDifferentComment),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasStartedString(r.BooleanTrue).
						HasAfterIdsInOrder(rootId).
						HasCommentString("def"),
				),
			},
		},
	})
}

func TestAcc_Task_issue2036(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5
	when := "TRUE"

	taskConfigModelWithoutWhen := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	taskConfigModelWithWhen := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement).
		WithWhen(when)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			// create without when
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, taskConfigModelWithoutWhen),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, taskConfigModelWithoutWhen.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasWhenString(""),
				),
			},
			// add when
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, taskConfigModelWithWhen),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, taskConfigModelWithWhen.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasWhenString("TRUE"),
				),
			},
			// remove when
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, taskConfigModelWithoutWhen),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, taskConfigModelWithoutWhen.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasWhenString(""),
				),
			},
		},
	})
}

func TestAcc_Task_issue3113(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	errorNotificationIntegration, errorNotificationIntegrationCleanup := acc.TestClient().NotificationIntegration.Create(t)
	t.Cleanup(errorNotificationIntegrationCleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5
	configModel := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement).
		WithErrorIntegration(errorNotificationIntegration.ID().Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.97.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config:      taskConfigWithErrorIntegration(id, errorNotificationIntegration.ID()),
				ExpectError: regexp.MustCompile("error_integration: '' expected type 'string', got unconvertible type 'sdk.AccountObjectIdentifier'"),
			},
			{
				PreConfig: func() {
					acc.TestClient().Task.DropFunc(t, id)()
				},
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          acc.ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables:          config.ConfigVariablesFromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()),
				),
			},
		},
	})
}

func taskConfigWithErrorIntegration(id sdk.SchemaObjectIdentifier, errorIntegrationId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_task" "test" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	schedule = "5 MINUTES"
	sql_statement = "SELECT 1"
	enabled = true
	error_integration = "%[4]s"
}
`, id.DatabaseName(), id.SchemaName(), id.Name(), errorIntegrationId.Name())
}
