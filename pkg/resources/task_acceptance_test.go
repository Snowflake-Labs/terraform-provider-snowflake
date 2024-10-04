package resources_test

import (
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO(SNOW-1348116 - next pr): More tests for complicated DAGs
// TODO(SNOW-1348116 - next pr): Test for stored procedures passed to sql_statement (decide on name)
// TODO(SNOW-1348116 - next pr): Test with cron schedule

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
						HasEnabledString(r.BooleanFalse).
						HasWarehouseString("").
						HasScheduleString("").
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfterIds().
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
						HasWarehouse("").
						HasSchedule("").
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration("").
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

	currentRole := acc.TestClient().Context.CurrentRole(t)

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
	configModel := model.TaskWithId("test", id, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithSchedule("10 MINUTES").
		WithConfigValue(configvariable.StringVariable(taskConfigVariableValue)).
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
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
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
					resourceshowoutputassert.TaskShowOutput(t, configModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment(comment).
						HasWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
						HasSchedule("10 MINUTES").
						HasPredecessors().
						HasState(sdk.TaskStateStarted).
						HasDefinition(statement).
						HasCondition(condition).
						HasAllowOverlappingExecution(true).
						HasErrorIntegration(errorNotificationIntegration.ID().Name()).
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

	currentRole := acc.TestClient().Context.CurrentRole(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	basicConfigModel := model.TaskWithId("test", id, false, statement)

	// New warehouse created, because the common one has lower-case letters that won't work
	warehouse, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	errorNotificationIntegration, errorNotificationIntegrationCleanup := acc.TestClient().NotificationIntegration.Create(t)
	t.Cleanup(errorNotificationIntegrationCleanup)

	taskConfig := `$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`
	// We have to do three $ at the beginning because Terraform will remove one $.
	// It's because `${` is a special pattern, and it's escaped by `$${`.
	expectedTaskConfig := strings.ReplaceAll(taskConfig, "$", "")
	taskConfigVariableValue := "$" + taskConfig
	comment := random.Comment()
	condition := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
	completeConfigModel := model.TaskWithId("test", id, true, statement).
		WithWarehouse(warehouse.ID().Name()).
		WithSchedule("5 MINUTES").
		WithConfigValue(configvariable.StringVariable(taskConfigVariableValue)).
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
						HasEnabledString(r.BooleanFalse).
						HasWarehouseString("").
						HasScheduleString("").
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfterIds().
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
						HasWarehouse("").
						HasSchedule("").
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration("").
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
				Config: config.FromModel(t, completeConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, completeConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasEnabledString(r.BooleanTrue).
						HasWarehouseString(warehouse.ID().Name()).
						HasScheduleString("5 MINUTES").
						HasConfigString(expectedTaskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasCommentString(comment).
						HasFinalizeString("").
						HasAfterIds().
						HasWhenString(condition).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, completeConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasWarehouse(warehouse.ID().Name()).
						HasComment(comment).
						HasSchedule("5 MINUTES").
						HasPredecessors().
						HasState(sdk.TaskStateStarted).
						HasDefinition(statement).
						HasCondition(condition).
						HasAllowOverlappingExecution(true).
						HasErrorIntegration(errorNotificationIntegration.ID().Name()).
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
						HasEnabledString(r.BooleanFalse).
						HasWarehouseString("").
						HasScheduleString("").
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfterIds().
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
						HasWarehouse("").
						HasSchedule("").
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration("").
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

func TestAcc_Task_AllParameters(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	configModel := model.TaskWithId("test", id, true, statement).
		WithSchedule("5 MINUTES")
	configModelWithAllParametersSet := model.TaskWithId("test", id, true, statement).
		WithSchedule("5 MINUTES").
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
				Config: config.FromModel(t, configModel),
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
				ResourceName: configModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceparametersassert.ImportedTaskResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasAllDefaults(),
				),
			},
			// set all parameters
			{
				Config: config.FromModel(t, configModelWithAllParametersSet),
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
				ResourceName: configModelWithAllParametersSet.ResourceReference(),
				ImportState:  true,
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
				Config: config.FromModel(t, configModel),
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
		WithSchedule("5 MINUTES")
	configModelDisabled := model.TaskWithId("test", id, false, statement).
		WithSchedule("5 MINUTES")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModelDisabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelDisabled.ResourceReference()).
						HasEnabledString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, configModelDisabled.ResourceReference()).
						HasState(sdk.TaskStateSuspended),
				),
			},
			{
				Config: config.FromModel(t, configModelEnabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelEnabled.ResourceReference()).
						HasEnabledString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, configModelEnabled.ResourceReference()).
						HasState(sdk.TaskStateStarted),
				),
			},
			{
				Config: config.FromModel(t, configModelDisabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, configModelDisabled.ResourceReference()).
						HasEnabledString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, configModelDisabled.ResourceReference()).
						HasState(sdk.TaskStateSuspended),
				),
			},
		},
	})
}

// TODO(SNOW-1348116 - analyze in next pr): This test may also be not deterministic and sometimes it fail when resuming a task while other task is modifying DAG (removing after)
func TestAcc_Task_ConvertStandaloneTaskToSubtask(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	id2 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := "5 MINUTES"

	firstTaskStandaloneModel := model.TaskWithId("main_task", id, true, statement).
		WithSchedule(schedule).
		WithSuspendTaskAfterNumFailures(1)
	secondTaskStandaloneModel := model.TaskWithId("second_task", id2, true, statement).
		WithSchedule(schedule)

	rootTaskModel := model.TaskWithId("main_task", id, true, statement).
		WithSchedule(schedule).
		WithSuspendTaskAfterNumFailures(2)
	childTaskModel := model.TaskWithId("second_task", id2, true, statement).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(id.FullyQualifiedName())))
	childTaskModel.SetDependsOn([]string{rootTaskModel.ResourceReference()})

	firstTaskStandaloneModelDisabled := model.TaskWithId("main_task", id, false, statement).
		WithSchedule(schedule)
	secondTaskStandaloneModelDisabled := model.TaskWithId("second_task", id2, false, statement).
		WithSchedule(schedule)
	secondTaskStandaloneModelDisabled.SetDependsOn([]string{firstTaskStandaloneModelDisabled.ResourceReference()})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, firstTaskStandaloneModel) + config.FromModel(t, secondTaskStandaloneModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, firstTaskStandaloneModel.ResourceReference()).
						HasScheduleString(schedule).
						HasEnabledString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("1"),
					resourceshowoutputassert.TaskShowOutput(t, firstTaskStandaloneModel.ResourceReference()).
						HasSchedule(schedule).
						HasState(sdk.TaskStateStarted),
					resourceassert.TaskResource(t, secondTaskStandaloneModel.ResourceReference()).
						HasScheduleString(schedule).
						HasEnabledString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, secondTaskStandaloneModel.ResourceReference()).
						HasSchedule(schedule).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change the second task to run after the first one (creating a DAG)
			{
				Config: config.FromModel(t, rootTaskModel) + config.FromModel(t, childTaskModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, rootTaskModel.ResourceReference()).
						HasScheduleString(schedule).
						HasEnabledString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("2"),
					resourceshowoutputassert.TaskShowOutput(t, rootTaskModel.ResourceReference()).
						HasSchedule(schedule).
						HasState(sdk.TaskStateStarted),
					resourceassert.TaskResource(t, childTaskModel.ResourceReference()).
						HasAfterIds(id).
						HasEnabledString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, childTaskModel.ResourceReference()).
						HasPredecessors(id).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change tasks in DAG to standalone tasks (disabled to check if resuming/suspending works correctly)
			{
				Config: config.FromModel(t, firstTaskStandaloneModelDisabled) + config.FromModel(t, secondTaskStandaloneModelDisabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, firstTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleString(schedule).
						HasEnabledString(r.BooleanFalse).
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceshowoutputassert.TaskShowOutput(t, firstTaskStandaloneModelDisabled.ResourceReference()).
						HasSchedule(schedule).
						HasState(sdk.TaskStateSuspended),
					resourceassert.TaskResource(t, secondTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleString(schedule).
						HasEnabledString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, secondTaskStandaloneModelDisabled.ResourceReference()).
						HasSchedule(schedule).
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
	schedule := "5 MINUTES"

	firstTaskStandaloneModel := model.TaskWithId("main_task", rootTaskId, true, statement).
		WithSchedule(schedule).
		WithSuspendTaskAfterNumFailures(1)
	secondTaskStandaloneModel := model.TaskWithId("second_task", finalizerTaskId, true, statement).
		WithSchedule(schedule)

	rootTaskModel := model.TaskWithId("main_task", rootTaskId, true, statement).
		WithSchedule(schedule).
		WithSuspendTaskAfterNumFailures(2)
	childTaskModel := model.TaskWithId("second_task", finalizerTaskId, true, statement).
		WithFinalize(rootTaskId.FullyQualifiedName())
	childTaskModel.SetDependsOn([]string{rootTaskModel.ResourceReference()})

	firstTaskStandaloneModelDisabled := model.TaskWithId("main_task", rootTaskId, false, statement).
		WithSchedule(schedule)
	secondTaskStandaloneModelDisabled := model.TaskWithId("second_task", finalizerTaskId, false, statement).
		WithSchedule(schedule)
	secondTaskStandaloneModelDisabled.SetDependsOn([]string{firstTaskStandaloneModelDisabled.ResourceReference()})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, firstTaskStandaloneModel) + config.FromModel(t, secondTaskStandaloneModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, firstTaskStandaloneModel.ResourceReference()).
						HasScheduleString(schedule).
						HasEnabledString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("1"),
					resourceshowoutputassert.TaskShowOutput(t, firstTaskStandaloneModel.ResourceReference()).
						HasSchedule(schedule).
						HasState(sdk.TaskStateStarted),
					resourceassert.TaskResource(t, secondTaskStandaloneModel.ResourceReference()).
						HasScheduleString(schedule).
						HasEnabledString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, secondTaskStandaloneModel.ResourceReference()).
						HasSchedule(schedule).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change the second task to run after the first one (creating a DAG)
			{
				Config: config.FromModel(t, rootTaskModel) + config.FromModel(t, childTaskModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, rootTaskModel.ResourceReference()).
						HasScheduleString(schedule).
						HasEnabledString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("2"),
					resourceshowoutputassert.TaskShowOutput(t, rootTaskModel.ResourceReference()).
						HasSchedule(schedule).
						// HasTaskRelations(sdk.TaskRelations{FinalizerTask: &finalizerTaskId}).
						HasState(sdk.TaskStateStarted),
					resourceassert.TaskResource(t, childTaskModel.ResourceReference()).
						HasEnabledString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, childTaskModel.ResourceReference()).
						// HasTaskRelations(sdk.TaskRelations{FinalizedRootTask: &rootTaskId}).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change tasks in DAG to standalone tasks (disabled to check if resuming/suspending works correctly)
			{
				Config: config.FromModel(t, firstTaskStandaloneModelDisabled) + config.FromModel(t, secondTaskStandaloneModelDisabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, firstTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleString(schedule).
						HasEnabledString(r.BooleanFalse).
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceshowoutputassert.TaskShowOutput(t, firstTaskStandaloneModelDisabled.ResourceReference()).
						HasSchedule(schedule).
						HasState(sdk.TaskStateSuspended),
					resourceassert.TaskResource(t, secondTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleString(schedule).
						HasEnabledString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, secondTaskStandaloneModelDisabled.ResourceReference()).
						HasSchedule(schedule).
						HasState(sdk.TaskStateSuspended),
				),
			},
		},
	})
}

// TODO(SNOW-1348116 - analyze in next pr): This test is not deterministic and sometimes it fails when resuming a task while other task is modifying DAG (removing after)
func TestAcc_Task_SwitchScheduledWithAfter(t *testing.T) {
	rootId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	childId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := "5 MINUTES"
	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithSchedule(schedule).
		WithSuspendTaskAfterNumFailures(1)
	childTaskConfigModel := model.TaskWithId("child", childId, true, statement).
		WithSchedule(schedule)

	rootTaskConfigModelAfterSuspendFailuresUpdate := model.TaskWithId("root", rootId, true, statement).
		WithSchedule(schedule).
		WithSuspendTaskAfterNumFailures(2)
	childTaskConfigModelWithAfter := model.TaskWithId("child", childId, true, statement).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName())))
	childTaskConfigModelWithAfter.SetDependsOn([]string{rootTaskConfigModelAfterSuspendFailuresUpdate.ResourceReference()})

	rootTaskConfigModelDisabled := model.TaskWithId("root", rootId, false, statement).
		WithSchedule(schedule)
	childTaskConfigModelDisabled := model.TaskWithId("child", childId, false, statement).
		WithSchedule(schedule)
	childTaskConfigModelDisabled.SetDependsOn([]string{rootTaskConfigModelDisabled.ResourceReference()})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, rootTaskConfigModel) + config.FromModel(t, childTaskConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasEnabledString(r.BooleanTrue).
						HasScheduleString(schedule).
						HasAfterIds().
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasEnabledString(r.BooleanTrue).
						HasScheduleString(schedule).
						HasSuspendTaskAfterNumFailuresString("1"),
				),
			},
			{
				Config: config.FromModel(t, rootTaskConfigModelAfterSuspendFailuresUpdate) + config.FromModel(t, childTaskConfigModelWithAfter),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasEnabledString(r.BooleanTrue).
						HasScheduleString("").
						HasAfterIds(rootId).
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasEnabledString(r.BooleanTrue).
						HasScheduleString(schedule).
						HasSuspendTaskAfterNumFailuresString("2"),
				),
			},
			{
				Config: config.FromModel(t, rootTaskConfigModel) + config.FromModel(t, childTaskConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasEnabledString(r.BooleanTrue).
						HasScheduleString(schedule).
						HasAfterIds().
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasEnabledString(r.BooleanTrue).
						HasScheduleString(schedule).
						HasSuspendTaskAfterNumFailuresString("1"),
				),
			},
			{
				Config: config.FromModel(t, rootTaskConfigModelDisabled) + config.FromModel(t, childTaskConfigModelDisabled),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, "snowflake_task.child").
						HasEnabledString(r.BooleanFalse).
						HasScheduleString(schedule).
						HasAfterIds().
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceassert.TaskResource(t, "snowflake_task.root").
						HasEnabledString(r.BooleanFalse).
						HasScheduleString(schedule).
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
	schedule := "5 MINUTES"

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithSchedule(schedule).
		WithSqlStatement(statement)

	childTaskConfigModelWithAfter := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName()))).
		WithSqlStatement(statement)
	childTaskConfigModelWithAfter.SetDependsOn([]string{rootTaskConfigModel.ResourceReference()})

	childTaskConfigModelWithoutAfter := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithSchedule(schedule).
		WithSqlStatement(statement)
	childTaskConfigModelWithoutAfter.SetDependsOn([]string{rootTaskConfigModel.ResourceReference()})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, rootTaskConfigModel) + config.FromModel(t, childTaskConfigModelWithAfter),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasScheduleString(schedule),
					resourceassert.TaskResource(t, childTaskConfigModelWithAfter.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasAfterIds(rootId),
				),
			},
			{
				Config: config.FromModel(t, rootTaskConfigModel) + config.FromModel(t, childTaskConfigModelWithoutAfter),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasScheduleString(schedule),
					resourceassert.TaskResource(t, childTaskConfigModelWithoutAfter.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasAfterIds(),
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
	schedule := "5 MINUTES"

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithSchedule(schedule).
		WithSqlStatement(statement)

	childTaskConfigModelWithFinalizer := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithFinalize(rootId.FullyQualifiedName()).
		WithSqlStatement(statement)
	childTaskConfigModelWithFinalizer.SetDependsOn([]string{rootTaskConfigModel.ResourceReference()})

	childTaskConfigModelWithoutFinalizer := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithSchedule(schedule).
		WithSqlStatement(statement)
	childTaskConfigModelWithoutFinalizer.SetDependsOn([]string{rootTaskConfigModel.ResourceReference()})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, rootTaskConfigModel) + config.FromModel(t, childTaskConfigModelWithFinalizer),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasScheduleString(schedule),
					resourceassert.TaskResource(t, childTaskConfigModelWithFinalizer.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasFinalizeString(rootId.FullyQualifiedName()),
				),
			},
			{
				Config: config.FromModel(t, rootTaskConfigModel) + config.FromModel(t, childTaskConfigModelWithoutFinalizer),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasScheduleString(schedule),
					resourceassert.TaskResource(t, childTaskConfigModelWithoutFinalizer.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasFinalizeString(""),
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
	schedule := "5 MINUTES"

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithSchedule(schedule).
		WithSqlStatement(statement)

	childTaskConfigModel := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName()))).
		WithComment("abc").
		WithSqlStatement(statement)
	childTaskConfigModel.SetDependsOn([]string{rootTaskConfigModel.ResourceReference()})

	childTaskConfigModelWithDifferentComment := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(acc.TestClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName()))).
		WithComment("def").
		WithSqlStatement(statement)
	childTaskConfigModelWithDifferentComment.SetDependsOn([]string{rootTaskConfigModel.ResourceReference()})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, rootTaskConfigModel) + config.FromModel(t, childTaskConfigModel),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasScheduleString(schedule),
					resourceassert.TaskResource(t, childTaskConfigModel.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasAfterIds(rootId).
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
				Config: config.FromModel(t, rootTaskConfigModel) + config.FromModel(t, childTaskConfigModelWithDifferentComment),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasScheduleString(schedule),
					resourceassert.TaskResource(t, childTaskConfigModelWithDifferentComment.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasAfterIds(rootId).
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
	schedule := "5 MINUTES"
	when := "TRUE"

	taskConfigModelWithoutWhen := model.TaskWithId("test", id, true, statement).
		WithSchedule(schedule).
		WithSqlStatement(statement)

	taskConfigModelWithWhen := model.TaskWithId("test", id, true, statement).
		WithSchedule(schedule).
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
				Config: config.FromModel(t, taskConfigModelWithoutWhen),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, taskConfigModelWithoutWhen.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasWhenString(""),
				),
			},
			// add when
			{
				Config: config.FromModel(t, taskConfigModelWithWhen),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, taskConfigModelWithWhen.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasWhenString("TRUE"),
				),
			},
			// remove when
			{
				Config: config.FromModel(t, taskConfigModelWithoutWhen),
				Check: assert.AssertThat(t,
					resourceassert.TaskResource(t, taskConfigModelWithoutWhen.ResourceReference()).
						HasEnabledString(r.BooleanTrue).
						HasWhenString(""),
				),
			},
		},
	})
}
