package resources_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

	basicConfigVariablesWithTransient := config.Variables{
		"name":         config.StringVariable(id.Name()),
		"comment":      config.StringVariable("foo"),
		"database":     config.StringVariable(databaseId.Name()),
		"is_transient": config.BoolVariable(true),
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
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					params := acc.TestClient().Parameter.ShowAccountParameters(t)
					*accountDataRetentionTimeInDays = helpers.FindParameter(t, params, sdk.AccountParameterDataRetentionTimeInDays).Value
					*accountMaxDataExtensionTimeInDays = helpers.FindParameter(t, params, sdk.AccountParameterMaxDataExtensionTimeInDays).Value
					*accountExternalVolume = helpers.FindParameter(t, params, sdk.AccountParameterExternalVolume).Value
					*accountCatalog = helpers.FindParameter(t, params, sdk.AccountParameterCatalog).Value
					*accountReplaceInvalidCharacters = helpers.FindParameter(t, params, sdk.AccountParameterReplaceInvalidCharacters).Value
					*accountDefaultDdlCollation = helpers.FindParameter(t, params, sdk.AccountParameterDefaultDDLCollation).Value
					*accountStorageSerializationPolicy = helpers.FindParameter(t, params, sdk.AccountParameterStorageSerializationPolicy).Value
					*accountLogLevel = helpers.FindParameter(t, params, sdk.AccountParameterLogLevel).Value
					*accountTraceLevel = helpers.FindParameter(t, params, sdk.AccountParameterTraceLevel).Value
					*accountSuspendTaskAfterNumFailures = helpers.FindParameter(t, params, sdk.AccountParameterSuspendTaskAfterNumFailures).Value
					*accountTaskAutoRetryAttempts = helpers.FindParameter(t, params, sdk.AccountParameterTaskAutoRetryAttempts).Value
					*accountUserTaskMangedInitialWarehouseSize = helpers.FindParameter(t, params, sdk.AccountParameterUserTaskManagedInitialWarehouseSize).Value
					*accountUserTaskTimeoutMs = helpers.FindParameter(t, params, sdk.AccountParameterUserTaskTimeoutMs).Value
					*accountUserTaskMinimumTriggerIntervalInSeconds = helpers.FindParameter(t, params, sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds).Value
					*accountQuotedIdentifiersIgnoreCase = helpers.FindParameter(t, params, sdk.AccountParameterQuotedIdentifiersIgnoreCase).Value
					*accountEnableConsoleOutput = helpers.FindParameter(t, params, sdk.AccountParameterEnableConsoleOutput).Value
					*accountPipeExecutionPaused = helpers.FindParameter(t, params, sdk.AccountParameterPipeExecutionPaused).Value
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema/basic"),
				ConfigVariables: basicConfigVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "with_managed_access", "false"),
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
			// import - without optionals
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_Schema/basic"),
				ConfigVariables:   basicConfigVariables,
				ResourceName:      "snowflake_schema.test",
				ImportState:       true,
				ImportStateVerify: true,
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
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_Schema/complete"),
				ConfigVariables:   completeConfigVariables,
				ResourceName:      "snowflake_schema.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema/basic"),
				ConfigVariables: basicConfigVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "with_managed_access", "false"),
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
				ConfigVariables: basicConfigVariablesWithTransient,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_schema.test", "with_managed_access", "false"),
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
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
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
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_Schema/complete"),
				ConfigVariables:   completeConfigVariables,
				ResourceName:      "snowflake_schema.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_Schema_Rename(t *testing.T) {
	oldSchemaName := acc.TestClient().Ids.Alpha()
	newSchemaName := acc.TestClient().Ids.Alpha()
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
					"name":     config.StringVariable(oldSchemaName),
					"database": config.StringVariable(acc.TestDatabaseName),
					"comment":  config.StringVariable(comment),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", oldSchemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", comment),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(newSchemaName),
					"database": config.StringVariable(acc.TestDatabaseName),
					"comment":  config.StringVariable(comment),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", newSchemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", comment),
				),
			},
		},
	})
}

func TestAcc_Schema_ManagePublic(t *testing.T) {
	name := "PUBLIC"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema/basic_with_pipe_execution_paused"),
				ConfigVariables: map[string]config.Variable{
					"name":                  config.StringVariable(name),
					"database":              config.StringVariable(acc.TestDatabaseName),
					"pipe_execution_paused": config.BoolVariable(true),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "pipe_execution_paused", "true"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema/basic_with_pipe_execution_paused"),
				ConfigVariables: map[string]config.Variable{
					"name":                  config.StringVariable(name),
					"database":              config.StringVariable(acc.TestDatabaseName),
					"pipe_execution_paused": config.BoolVariable(false),
				},
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
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	id := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)

	configVariablesWithoutSchemaDataRetentionTime := func(databaseDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                     config.StringVariable(databaseName),
			"schema":                       config.StringVariable(schemaName),
			"database_data_retention_time": config.IntegerVariable(databaseDataRetentionTime),
		}
	}

	configVariablesWithSchemaDataRetentionTime := func(databaseDataRetentionTime int, schemaDataRetentionTime int) config.Variables {
		vars := configVariablesWithoutSchemaDataRetentionTime(databaseDataRetentionTime)
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
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "5"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 15),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "15"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 15),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "0"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 0),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "3"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 3),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Schema_DefaultDataRetentionTime_SetOutsideOfTerraform(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	id := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)

	configVariablesWithoutSchemaDataRetentionTime := func(databaseDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                     config.StringVariable(databaseName),
			"schema":                       config.StringVariable(schemaName),
			"database_data_retention_time": config.IntegerVariable(databaseDataRetentionTime),
		}
	}

	configVariablesWithSchemaDataRetentionTime := func(databaseDataRetentionTime int, schemaDataRetentionTime int) config.Variables {
		vars := configVariablesWithoutSchemaDataRetentionTime(databaseDataRetentionTime)
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
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 5, 5),
				),
			},
			{
				PreConfig:       acc.TestClient().Schema.UpdateDataRetentionTime(t, id, 20),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithoutDataRetentionSet"),
				ConfigVariables: configVariablesWithoutSchemaDataRetentionTime(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schema_DefaultDataRetentionTime/WithDataRetentionSet"),
				ConfigVariables: configVariablesWithSchemaDataRetentionTime(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "data_retention_time_in_days", "3"),
					checkDatabaseAndSchemaDataRetentionTime(t, id, 10, 3),
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
	schemaName := schemaId.Name()
	configVariables := map[string]config.Variable{
		"schema_name":   config.StringVariable(schemaName),
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
				ExpectError: regexp.MustCompile("error creating schema"),
			},
		},
	})
}

func checkDatabaseAndSchemaDataRetentionTime(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, expectedDatabaseRetentionsDays int, expectedSchemaRetentionDays int) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		schema, err := acc.TestClient().Schema.Show(t, schemaId)
		if err != nil {
			return err
		}

		database, err := acc.TestClient().Database.Show(t, schemaId.DatabaseId())
		if err != nil {
			return err
		}

		// "retention_time" may sometimes be an empty string instead of an integer
		var schemaRetentionTime int64
		{
			rt := schema.RetentionTime
			if rt == "" {
				rt = "0"
			}

			schemaRetentionTime, err = strconv.ParseInt(rt, 10, 64)
			if err != nil {
				return err
			}
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

func TestAcc_Schema_migrateFromVersion093(t *testing.T) {
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
				Config: schemav093(id.Name(), databaseId.Name(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "true"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   schemav094(id.Name(), databaseId.Name(), true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckNoResourceAttr(resourceName, "is_managed"),
					resource.TestCheckResourceAttr(resourceName, "with_managed_access", "true"),
				),
			},
		},
	})
}

func schemav093(name, database string, isManaged bool) string {
	s := `
resource "snowflake_schema" "test" {
	name             = "%s"
	database		 = "%s"
	is_managed		 = %t
}
`
	return fmt.Sprintf(s, name, database, isManaged)
}

func schemav094(name, database string, isManaged bool) string {
	s := `
resource "snowflake_schema" "test" {
	name             = "%s"
	database		 = "%s"
	with_managed_access		 = %t
}
`
	return fmt.Sprintf(s, name, database, isManaged)
}
