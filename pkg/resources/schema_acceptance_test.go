package resources_test

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	acchelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Schema_basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	databaseId := acc.TestClient().Ids.DatabaseId()

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	basicConfigVariables := config.Variables{
		"name":     config.StringVariable(id.Name()),
		"comment":  config.StringVariable("foo"),
		"database": config.StringVariable(databaseId.Name()),
	}

	basicConfigVariablesWithTransient := func(isTransient bool) config.Variables {
		return config.Variables{
			"name":         config.StringVariable(id.Name()),
			"comment":      config.StringVariable("foo"),
			"database":     config.StringVariable(databaseId.Name()),
			"is_transient": config.BoolVariable(isTransient),
		}
	}

	completeConfigVariables := config.Variables{
		"name":                config.StringVariable(id.Name()),
		"comment":             config.StringVariable("foo"),
		"database":            config.StringVariable(databaseId.Name()),
		"with_managed_access": config.BoolVariable(true),
		"is_transient":        config.BoolVariable(false),

		"data_retention_time_in_days":                   config.IntegerVariable(1),
		"max_data_extension_time_in_days":               config.IntegerVariable(1),
		"external_volume":                               config.StringVariable(externalVolumeId.Name()),
		"catalog":                                       config.StringVariable(catalogId.Name()),
		"replace_invalid_characters":                    config.BoolVariable(true),
		"default_ddl_collation":                         config.StringVariable("en_US"),
		"storage_serialization_policy":                  config.StringVariable(string(sdk.StorageSerializationPolicyCompatible)),
		"log_level":                                     config.StringVariable(string(sdk.LogLevelInfo)),
		"trace_level":                                   config.StringVariable(string(sdk.TraceLevelOnEvent)),
		"suspend_task_after_num_failures":               config.IntegerVariable(20),
		"task_auto_retry_attempts":                      config.IntegerVariable(20),
		"user_task_managed_initial_warehouse_size":      config.StringVariable(string(sdk.WarehouseSizeXLarge)),
		"user_task_timeout_ms":                          config.IntegerVariable(1200000),
		"user_task_minimum_trigger_interval_in_seconds": config.IntegerVariable(120),
		"quoted_identifiers_ignore_case":                config.BoolVariable(true),
		"enable_console_output":                         config.BoolVariable(true),
		"pipe_execution_paused":                         config.BoolVariable(true),
	}

	var (
		accountDataRetentionTimeInDays                 = new(string)
		accountMaxDataExtensionTimeInDays              = new(string)
		accountExternalVolume                          = new(string)
		accountCatalog                                 = new(string)
		accountReplaceInvalidCharacters                = new(string)
		accountDefaultDdlCollation                     = new(string)
		accountStorageSerializationPolicy              = new(string)
		accountLogLevel                                = new(string)
		accountTraceLevel                              = new(string)
		accountSuspendTaskAfterNumFailures             = new(string)
		accountTaskAutoRetryAttempts                   = new(string)
		accountUserTaskMangedInitialWarehouseSize      = new(string)
		accountUserTaskTimeoutMs                       = new(string)
		accountUserTaskMinimumTriggerIntervalInSeconds = new(string)
		accountQuotedIdentifiersIgnoreCase             = new(string)
		accountEnableConsoleOutput                     = new(string)
		accountPipeExecutionPaused                     = new(string)
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					params := acc.TestClient().Parameter.ShowAccountParameters(t)
					*accountDataRetentionTimeInDays = acchelpers.FindParameter(t, params, sdk.AccountParameterDataRetentionTimeInDays).Value
					*accountMaxDataExtensionTimeInDays = acchelpers.FindParameter(t, params, sdk.AccountParameterMaxDataExtensionTimeInDays).Value
					*accountExternalVolume = acchelpers.FindParameter(t, params, sdk.AccountParameterExternalVolume).Value
					*accountCatalog = acchelpers.FindParameter(t, params, sdk.AccountParameterCatalog).Value
					*accountReplaceInvalidCharacters = acchelpers.FindParameter(t, params, sdk.AccountParameterReplaceInvalidCharacters).Value
					*accountDefaultDdlCollation = acchelpers.FindParameter(t, params, sdk.AccountParameterDefaultDDLCollation).Value
					*accountStorageSerializationPolicy = acchelpers.FindParameter(t, params, sdk.AccountParameterStorageSerializationPolicy).Value
					*accountLogLevel = acchelpers.FindParameter(t, params, sdk.AccountParameterLogLevel).Value
					*accountTraceLevel = acchelpers.FindParameter(t, params, sdk.AccountParameterTraceLevel).Value
					*accountSuspendTaskAfterNumFailures = acchelpers.FindParameter(t, params, sdk.AccountParameterSuspendTaskAfterNumFailures).Value
					*accountTaskAutoRetryAttempts = acchelpers.FindParameter(t, params, sdk.AccountParameterTaskAutoRetryAttempts).Value
					*accountUserTaskMangedInitialWarehouseSize = acchelpers.FindParameter(t, params, sdk.AccountParameterUserTaskManagedInitialWarehouseSize).Value
					*accountUserTaskTimeoutMs = acchelpers.FindParameter(t, params, sdk.AccountParameterUserTaskTimeoutMs).Value
					*accountUserTaskMinimumTriggerIntervalInSeconds = acchelpers.FindParameter(t, params, sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds).Value
					*accountQuotedIdentifiersIgnoreCase = acchelpers.FindParameter(t, params, sdk.AccountParameterQuotedIdentifiersIgnoreCase).Value
					*accountEnableConsoleOutput = acchelpers.FindParameter(t, params, sdk.AccountParameterEnableConsoleOutput).Value
					*accountPipeExecutionPaused = acchelpers.FindParameter(t, params, sdk.AccountParameterPipeExecutionPaused).Value
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema/basic"),
				ConfigVariables: basicConfigVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "with_managed_access", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_schema.test", "is_transient", r.BooleanDefault),

					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "enable_console_output", accountEnableConsoleOutput),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "pipe_execution_paused", accountPipeExecutionPaused),

					resource.TestCheckResourceAttrSet("snowflake_schema.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("snowflake_schema.test", "show_output.0.is_current"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttrSet("snowflake_schema.test", "show_output.0.owner"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_schema.test", "show_output.0.options", ""),
				),
			},
			// import - without optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema/basic"),
				ConfigVariables: basicConfigVariables,
				ResourceName:    "snowflake_schema.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", databaseId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "with_managed_access", "false"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "is_transient", "false"),
				),
			},
			// set other fields
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema/complete"),
				ConfigVariables: completeConfigVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "with_managed_access", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", "foo"),

					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "max_data_extension_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "enable_console_output", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "pipe_execution_paused", "true"),

					resource.TestCheckResourceAttrSet("snowflake_schema.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("snowflake_schema.test", "show_output.0.is_current"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttrSet("snowflake_schema.test", "show_output.0.owner"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "show_output.0.options", "MANAGED ACCESS"),

					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.data_retention_time_in_days.0.value", "1"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.max_data_extension_time_in_days.0.value", "1"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.external_volume.0.value", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.catalog.0.value", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.replace_invalid_characters.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.default_ddl_collation.0.value", "en_US"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.storage_serialization_policy.0.value", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.log_level.0.value", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.trace_level.0.value", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.suspend_task_after_num_failures.0.value", "20"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.task_auto_retry_attempts.0.value", "20"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.user_task_managed_initial_warehouse_size.0.value", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.user_task_timeout_ms.0.value", "1200000"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.user_task_minimum_trigger_interval_in_seconds.0.value", "120"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.quoted_identifiers_ignore_case.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.enable_console_output.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "parameters.0.pipe_execution_paused.0.value", "true"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_Schema/complete"),
				ConfigVariables:         completeConfigVariables,
				ResourceName:            "snowflake_schema.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"show_output.0.is_current"},
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema/basic_with_transient"),
				ConfigVariables: basicConfigVariablesWithTransient(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "with_managed_access", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_schema.test", "is_transient", "false"),

					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "enable_console_output", accountEnableConsoleOutput),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "pipe_execution_paused", accountPipeExecutionPaused),
				),
			},
			// set is_transient - recreate
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema/basic_with_transient"),
				ConfigVariables: basicConfigVariablesWithTransient(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "with_managed_access", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_schema.test", "is_transient", "true"),

					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "enable_console_output", accountEnableConsoleOutput),
					resource.TestCheckResourceAttrPtr("snowflake_schema.test", "pipe_execution_paused", accountPipeExecutionPaused),
				),
			},
		},
	})
}

func TestAcc_Schema_complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	databaseId := acc.TestClient().Ids.DatabaseId()

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	completeConfigVariables := config.Variables{
		"name":                config.StringVariable(id.Name()),
		"comment":             config.StringVariable("foo"),
		"database":            config.StringVariable(databaseId.Name()),
		"with_managed_access": config.BoolVariable(true),
		"is_transient":        config.BoolVariable(true),

		"data_retention_time_in_days":                   config.IntegerVariable(1),
		"max_data_extension_time_in_days":               config.IntegerVariable(1),
		"external_volume":                               config.StringVariable(externalVolumeId.Name()),
		"catalog":                                       config.StringVariable(catalogId.Name()),
		"replace_invalid_characters":                    config.BoolVariable(true),
		"default_ddl_collation":                         config.StringVariable("en_US"),
		"storage_serialization_policy":                  config.StringVariable(string(sdk.StorageSerializationPolicyCompatible)),
		"log_level":                                     config.StringVariable(string(sdk.LogLevelInfo)),
		"trace_level":                                   config.StringVariable(string(sdk.TraceLevelOnEvent)),
		"suspend_task_after_num_failures":               config.IntegerVariable(20),
		"task_auto_retry_attempts":                      config.IntegerVariable(20),
		"user_task_managed_initial_warehouse_size":      config.StringVariable(string(sdk.WarehouseSizeXLarge)),
		"user_task_timeout_ms":                          config.IntegerVariable(1200000),
		"user_task_minimum_trigger_interval_in_seconds": config.IntegerVariable(120),
		"quoted_identifiers_ignore_case":                config.BoolVariable(true),
		"enable_console_output":                         config.BoolVariable(true),
		"pipe_execution_paused":                         config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema/complete"),
				ConfigVariables: completeConfigVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "is_transient", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "with_managed_access", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", "foo"),

					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "max_data_extension_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_schema.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "enable_console_output", "true"),
					resource.TestCheckResourceAttr("snowflake_schema.test", "pipe_execution_paused", "true"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_Schema/complete"),
				ConfigVariables:         completeConfigVariables,
				ResourceName:            "snowflake_schema.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"show_output.0.is_current"},
			},
		},
	})
}

func TestAcc_Schema_Rename(t *testing.T) {
	oldId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	newId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	comment := "Terraform acceptance test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(oldId.Name()),
					"database": config.StringVariable(acc.TestDatabaseName),
					"comment":  config.StringVariable(comment),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", oldId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "fully_qualified_name", oldId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", comment),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(newId.Name()),
					"database": config.StringVariable(acc.TestDatabaseName),
					"comment":  config.StringVariable(comment),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", comment),
				),
			},
		},
	})
}

func TestAcc_Schema_ManagePublicVersion_0_94_0(t *testing.T) {
	name := "PUBLIC"
	schemaId := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// PUBLIC can not be created in v0.93
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.93.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config:      schemav093(name, acc.TestDatabaseName),
				ExpectError: regexp.MustCompile("Error: error creating schema PUBLIC"),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: schemav094WithPipeExecutionPaused(name, acc.TestDatabaseName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "pipe_execution_paused", "true"),
				),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				PreConfig: func() {
					// In v0.94 `CREATE OR REPLACE` was called, so we should see a DROP event.
					schemas := acc.TestClient().Schema.ShowWithOptions(t, &sdk.ShowSchemaOptions{
						History: sdk.Pointer(true),
						Like: &sdk.Like{
							Pattern: sdk.String(schemaId.Name()),
						},
					})
					require.Len(t, schemas, 2)
					slices.SortFunc(schemas, func(x, y sdk.Schema) int {
						return cmp.Compare(x.DroppedOn.Unix(), y.DroppedOn.Unix())
					})
					require.Zero(t, schemas[0].DroppedOn)
					require.NotZero(t, schemas[1].DroppedOn)
				},
				Config: schemav094WithPipeExecutionPaused(name, acc.TestDatabaseName, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "pipe_execution_paused", "false"),
				),
			},
		},
	})
}

func TestAcc_Schema_ManagePublicVersion_0_94_1(t *testing.T) {
	name := "PUBLIC"

	// use a separate db because this test relies on schema history
	db, cleanupDb := acc.TestClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDb)
	schemaId := sdk.NewDatabaseObjectIdentifier(db.ID().Name(), name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// PUBLIC can not be created in v0.93
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.93.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config:      schemav093(name, db.ID().Name()),
				ExpectError: regexp.MustCompile("Error: error creating schema PUBLIC"),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   schemav094WithPipeExecutionPaused(name, db.ID().Name(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", db.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "pipe_execution_paused", "true"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					// In newer versions, ALTER was called, so we should not see a DROP event.
					schemas := acc.TestClient().Schema.ShowWithOptions(t, &sdk.ShowSchemaOptions{
						History: sdk.Pointer(true),
						Like: &sdk.Like{
							Pattern: sdk.String(schemaId.Name()),
						},
					})
					require.Len(t, schemas, 1)
					require.Zero(t, schemas[0].DroppedOn)
				},
				Config: schemav094WithPipeExecutionPaused(name, db.ID().Name(), true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", db.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "pipe_execution_paused", "true"),
				),
			},
		},
	})
}

// TestAcc_Schema_TwoSchemasWithTheSameNameOnDifferentDatabases proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2209 issue.
func TestAcc_Schema_TwoSchemasWithTheSameNameOnDifferentDatabases(t *testing.T) {
	name := "test_schema"
	// It seems like Snowflake orders the output of SHOW command based on names, so they do matter
	newDatabaseName := "SELDQBXEKC"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(name),
					"database": config.StringVariable(acc.TestDatabaseName),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":         config.StringVariable(name),
					"database":     config.StringVariable(acc.TestDatabaseName),
					"new_database": config.StringVariable(newDatabaseName),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test_2", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test_2", "database", newDatabaseName),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Schema_DefaultDataRetentionTime(t *testing.T) {
	db, dbCleanup := acc.TestClient().Database.CreateDatabase(t)
	t.Cleanup(dbCleanup)

	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(db.ID())

	configVariablesWithoutSchemaDataRetentionTime := func() config.Variables {
		return config.Variables{
			"database": config.StringVariable(db.ID().Name()),
			"schema":   config.StringVariable(id.Name()),
		}
	}

	configVariablesWithSchemaDataRetentionTime := func(schemaDataRetentionTime int) config.Variables {
		vars := configVariablesWithoutSchemaDataRetentionTime()
		vars["schema_data_retention_time"] = config.IntegerVariable(schemaDataRetentionTime)
		return vars
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "1"),
				),
			},
			// change param value in database
			{
				PreConfig: func() {
					acc.TestClient().Database.UpdateDataRetentionTime(t, db.ID(), 50)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_schema.test", "data_retention_time_in_days"),
						planchecks.ExpectDrift("snowflake_schema.test", "data_retention_time_in_days", sdk.String("1"), sdk.String("50")),
						planchecks.ExpectChange("snowflake_schema.test", "data_retention_time_in_days", tfjson.ActionNoop, sdk.String("50"), sdk.String("50")),
						planchecks.ExpectComputed("snowflake_schema.test", "data_retention_time_in_days", false),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "50"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "5"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 50, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(15),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "15"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 50, 15),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "50"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 50, 50),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "0"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 50, 0),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(3),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "3"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 50, 3),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Schema_DefaultDataRetentionTime_SetOutsideOfTerraform(t *testing.T) {
	databaseId := acc.TestClient().Ids.DatabaseId()
	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariablesWithoutSchemaDataRetentionTime := func() config.Variables {
		return config.Variables{
			"database": config.StringVariable(databaseId.Name()),
			"schema":   config.StringVariable(id.Name()),
		}
	}

	configVariablesWithSchemaDataRetentionTime := func(schemaDataRetentionTime int) config.Variables {
		vars := configVariablesWithoutSchemaDataRetentionTime()
		vars["schema_data_retention_time"] = config.IntegerVariable(schemaDataRetentionTime)
		return vars
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "1"),
				),
			},
			{
				PreConfig:       acc.TestClient().Schema.UpdateDataRetentionTime(t, id, 20),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "1"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "3"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAcc_Schema_RemoveDatabaseOutsideOfTerraform(t *testing.T) {
	schemaId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	configVariables := map[string]config.Variable{
		"schema_name":   config.StringVariable(schemaId.Name()),
		"database_name": config.StringVariable(acc.TestDatabaseName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_RemoveOutsideOfTerraform"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					acc.TestClient().Schema.DropSchemaFunc(t, schemaId)()
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				RefreshPlanChecks: resource.RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

func TestAcc_Schema_RemoveSchemaOutsideOfTerraform(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	schemaName := acc.TestClient().Ids.Alpha()
	configVariables := map[string]config.Variable{
		"schema_name":   config.StringVariable(schemaName),
		"database_name": config.StringVariable(databaseName),
	}

	var cleanupDatabase func()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanupDatabase = acc.TestClient().Database.CreateDatabaseWithIdentifier(t, databaseId)
					t.Cleanup(cleanupDatabase)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_RemoveOutsideOfTerraform"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					cleanupDatabase()
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_RemoveOutsideOfTerraform"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("Failed to create schema"),
			},
		},
	})
}

func checkDatabaseAndSchemaDataRetentionTime(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, expectedDatabaseRetentionsDays int, expectedSchemaRetentionDays int) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		schema, err := acc.TestClient().Schema.Show(t, schemaId)
		require.NoError(t, err)

		database, err := acc.TestClient().Database.Show(t, schemaId.DatabaseId())
		require.NoError(t, err)

		// "retention_time" may sometimes be an empty string instead of an integer
		var schemaRetentionTime int64
		{
			rt := schema.RetentionTime
			if rt == "" {
				rt = "0"
			}

			schemaRetentionTime, err = strconv.ParseInt(rt, 10, 64)
			require.NoError(t, err)
		}

		if database.RetentionTime != expectedDatabaseRetentionsDays {
			return fmt.Errorf("invalid database retention time, expected: %d, got: %d", expectedDatabaseRetentionsDays, database.RetentionTime)
		}

		if schemaRetentionTime != int64(expectedSchemaRetentionDays) {
			return fmt.Errorf("invalid schema retention time, expected: %d, got: %d", expectedSchemaRetentionDays, schemaRetentionTime)
		}

		return nil
	}
}

func TestAcc_Schema_migrateFromVersion093WithoutManagedAccess(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseId := acc.TestClient().Ids.DatabaseId()
	resourceName := "snowflake_schema.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.93.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: schemav093(id.Name(), databaseId.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "false"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   schemav094(id.Name(), databaseId.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "with_managed_access", r.BooleanDefault),
				),
			},
		},
	})
}

func TestAcc_Schema_migrateFromVersion093(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseId := acc.TestClient().Ids.DatabaseId()
	resourceName := "snowflake_schema.test"

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.93.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: schemav093WithIsManagedAndDataRetentionDays(id.Name(), databaseId.Name(), tag.SchemaName, tag.Name, "foo", true, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "true"),
					resource.TestCheckResourceAttr(resourceName, "data_retention_days", "10"),
					resource.TestCheckResourceAttr(resourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tag.0.name", tag.Name),
					resource.TestCheckResourceAttr(resourceName, "tag.0.value", "foo"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   schemav094WithManagedAccessAndDataRetentionTimeInDays(id.Name(), databaseId.Name(), true, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckNoResourceAttr(resourceName, "is_managed"),
					resource.TestCheckResourceAttr(resourceName, "with_managed_access", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "data_retention_days"),
					resource.TestCheckResourceAttr(resourceName, "data_retention_time_in_days", "10"),
					resource.TestCheckNoResourceAttr(resourceName, "tag.#"),
				),
			},
		},
	})
}

func schemav093WithIsManagedAndDataRetentionDays(name, database, tagSchema, tagName, tagValue string, isManaged bool, dataRetentionDays int) string {
	s := `
resource "snowflake_schema" "test" {
	name					= "%[1]s"
	database				= "%[2]s"
	is_managed				= %[6]t
	data_retention_days		= %[7]d
	tag {
		name = "%[4]s"
		value = "%[5]s"
		schema = "%[3]s"
		database = "%[2]s"
	}
}
`
	return fmt.Sprintf(s, name, database, tagSchema, tagName, tagValue, isManaged, dataRetentionDays)
}

func schemav093(name, database string) string {
	s := `
resource "snowflake_schema" "test" {
	name					= "%s"
	database				= "%s"
}
`
	return fmt.Sprintf(s, name, database)
}

func schemav094WithManagedAccessAndDataRetentionTimeInDays(name, database string, isManaged bool, dataRetentionDays int) string {
	s := `
resource "snowflake_schema" "test" {
	name             				= "%s"
	database		 				= "%s"
	with_managed_access				= %t
	data_retention_time_in_days		= %d
}
`
	return fmt.Sprintf(s, name, database, isManaged, dataRetentionDays)
}

func schemav094(name, database string) string {
	s := `
resource "snowflake_schema" "test" {
	name             				= "%s"
	database		 				= "%s"
}
`
	return fmt.Sprintf(s, name, database)
}

func schemav094WithPipeExecutionPaused(name, database string, pipeExecutionPaused bool) string {
	s := `
resource "snowflake_schema" "test" {
	name             				= "%s"
	database		 				= "%s"
	pipe_execution_paused			= %t
}
`
	return fmt.Sprintf(s, name, database, pipeExecutionPaused)
}

func TestAcc_Schema_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: schemaBasicConfig(id.DatabaseName(), id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "id", helpers.EncodeSnowflakeID(id)),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   schemaBasicConfig(id.DatabaseName(), id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

func schemaBasicConfig(databaseName string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	database = "%s"
	name     = "%s"
}
`, databaseName, name)
}

func TestAcc_Schema_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	quotedDatabaseName := fmt.Sprintf(`\"%s\"`, id.DatabaseName())
	quotedName := fmt.Sprintf(`\"%s\"`, id.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				ExpectNonEmptyPlan: true,
				Config:             schemaBasicConfig(quotedDatabaseName, quotedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "id", fmt.Sprintf(`"%s"|"%s"`, id.DatabaseName(), id.Name())),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   schemaBasicConfig(quotedDatabaseName, quotedName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "id", id.Name()),
				),
			},
		},
	})
}
