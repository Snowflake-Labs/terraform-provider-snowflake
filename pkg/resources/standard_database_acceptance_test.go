package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_StandardDatabase_Basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	newId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newComment := random.Comment()

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
	)

	configVariables := func(id sdk.AccountObjectIdentifier, comment string) config.Variables {
		return config.Variables{
			"name":    config.StringVariable(id.Name()),
			"comment": config.StringVariable(comment),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StandardDatabase),
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
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.#", "0"),

					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(newId, newComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", newComment),

					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables:   configVariables(newId, newComment),
				ResourceName:      "snowflake_standard_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_StandardDatabase_ComputedValues(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	configVariables := func(id sdk.AccountObjectIdentifier, comment string) config.Variables {
		return config.Variables{
			"name":    config.StringVariable(id.Name()),
			"comment": config.StringVariable(comment),
		}
	}

	secondaryAccountIdentifier := acc.SecondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

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
	)

	completeConfigVariables := config.Variables{
		"name":                                     config.StringVariable(id.Name()),
		"comment":                                  config.StringVariable(comment),
		"transient":                                config.BoolVariable(false),
		"account_identifier":                       config.StringVariable(secondaryAccountIdentifier),
		"with_failover":                            config.BoolVariable(true),
		"ignore_edition_check":                     config.BoolVariable(true),
		"data_retention_time_in_days":              config.IntegerVariable(20),
		"max_data_extension_time_in_days":          config.IntegerVariable(30),
		"external_volume":                          config.StringVariable(externalVolumeId.Name()),
		"catalog":                                  config.StringVariable(catalogId.Name()),
		"replace_invalid_characters":               config.BoolVariable(true),
		"default_ddl_collation":                    config.StringVariable("en_US"),
		"storage_serialization_policy":             config.StringVariable(string(sdk.StorageSerializationPolicyCompatible)),
		"log_level":                                config.StringVariable(string(sdk.LogLevelInfo)),
		"trace_level":                              config.StringVariable(string(sdk.TraceLevelOnEvent)),
		"suspend_task_after_num_failures":          config.IntegerVariable(20),
		"task_auto_retry_attempts":                 config.IntegerVariable(20),
		"user_task_managed_initial_warehouse_size": config.StringVariable(string(sdk.WarehouseSizeXLarge)),
		"user_task_timeout_ms":                     config.IntegerVariable(1200000),
		"user_task_minimum_trigger_interval_in_seconds": config.IntegerVariable(120),
		"quoted_identifiers_ignore_case":                config.BoolVariable(true),
		"enable_console_output":                         config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StandardDatabase),
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
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", comment),

					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/complete-optionals-set"),
				ConfigVariables: completeConfigVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_standard_database.test", "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days", "30"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "enable_console_output", "true"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", comment),

					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
		},
	})
}

func TestAcc_StandardDatabase_Complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	secondaryAccountIdentifier := acc.SecondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()
	comment := random.Comment()

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	completeConfigVariables := config.Variables{
		"name":                 config.StringVariable(id.Name()),
		"comment":              config.StringVariable(comment),
		"transient":            config.BoolVariable(false),
		"account_identifier":   config.StringVariable(secondaryAccountIdentifier),
		"with_failover":        config.BoolVariable(true),
		"ignore_edition_check": config.BoolVariable(true),

		"data_retention_time_in_days":                   config.IntegerVariable(20),
		"max_data_extension_time_in_days":               config.IntegerVariable(30),
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
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StandardDatabase),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/complete-optionals-set"),
				ConfigVariables: completeConfigVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_standard_database.test", "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days", "30"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "enable_console_output", "true"),

					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.0.with_failover", "true"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_StandardDatabase/complete-optionals-set"),
				ConfigVariables:         completeConfigVariables,
				ResourceName:            "snowflake_standard_database.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication.0.ignore_edition_check"},
			},
		},
	})
}

func TestAcc_StandardDatabase_Update(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	newId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newComment := random.Comment()

	secondaryAccountIdentifier := acc.SecondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	basicConfigVariables := func(id sdk.AccountObjectIdentifier, comment string) config.Variables {
		return config.Variables{
			"name":    config.StringVariable(id.Name()),
			"comment": config.StringVariable(comment),
		}
	}

	completeConfigVariables := config.Variables{
		"name":                                     config.StringVariable(newId.Name()),
		"comment":                                  config.StringVariable(newComment),
		"transient":                                config.BoolVariable(false),
		"account_identifier":                       config.StringVariable(secondaryAccountIdentifier),
		"with_failover":                            config.BoolVariable(true),
		"ignore_edition_check":                     config.BoolVariable(true),
		"data_retention_time_in_days":              config.IntegerVariable(20),
		"max_data_extension_time_in_days":          config.IntegerVariable(30),
		"external_volume":                          config.StringVariable(externalVolumeId.Name()),
		"catalog":                                  config.StringVariable(catalogId.Name()),
		"replace_invalid_characters":               config.BoolVariable(true),
		"default_ddl_collation":                    config.StringVariable("en_US"),
		"storage_serialization_policy":             config.StringVariable(string(sdk.StorageSerializationPolicyCompatible)),
		"log_level":                                config.StringVariable(string(sdk.LogLevelInfo)),
		"trace_level":                              config.StringVariable(string(sdk.TraceLevelOnEvent)),
		"suspend_task_after_num_failures":          config.IntegerVariable(20),
		"task_auto_retry_attempts":                 config.IntegerVariable(20),
		"user_task_managed_initial_warehouse_size": config.StringVariable(string(sdk.WarehouseSizeXLarge)),
		"user_task_timeout_ms":                     config.IntegerVariable(1200000),
		"user_task_minimum_trigger_interval_in_seconds": config.IntegerVariable(120),
		"quoted_identifiers_ignore_case":                config.BoolVariable(true),
		"enable_console_output":                         config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StandardDatabase),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: basicConfigVariables(id, comment),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/complete-optionals-set"),
				ConfigVariables: completeConfigVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", newComment),

					resource.TestCheckResourceAttr("snowflake_standard_database.test", "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days", "30"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "enable_console_output", "true"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: basicConfigVariables(id, comment),
			},
		},
	})
}

func TestAcc_StandardDatabase_HierarchicalValues(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	configVariables := func(id sdk.AccountObjectIdentifier, comment string) config.Variables {
		return config.Variables{
			"name":    config.StringVariable(id.Name()),
			"comment": config.StringVariable(comment),
		}
	}

	paramDefault := new(string)
	var revertAccountParameterToDefault func()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StandardDatabase),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					*paramDefault = acc.TestClient().Parameter.GetAccountParameter(t, sdk.AccountParameterMaxDataExtensionTimeInDays).Default
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "max_data_extension_time_in_days", paramDefault),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterToDefault = acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterMaxDataExtensionTimeInDays, "50")
					t.Cleanup(revertAccountParameterToDefault)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days", "50"),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterToDefault()
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "max_data_extension_time_in_days", paramDefault),
				),
			},
		},
	})
}

func TestAcc_StandardDatabase_Replication(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	secondaryAccountIdentifier := acc.SecondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()

	configVariables := func(id sdk.AccountObjectIdentifier, withReplication bool, withFailover bool, ignoreEditionCheck bool) config.Variables {
		if withReplication {
			return config.Variables{
				"name":                 config.StringVariable(id.Name()),
				"account_identifier":   config.StringVariable(secondaryAccountIdentifier),
				"with_failover":        config.BoolVariable(withFailover),
				"ignore_edition_check": config.BoolVariable(ignoreEditionCheck),
			}
		}
		return config.Variables{
			"name":    config.StringVariable(id.Name()),
			"comment": config.StringVariable(""),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StandardDatabase),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.#", "0"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/replication"),
				ConfigVariables: configVariables(id, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.0.with_failover", "true"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/replication"),
				ConfigVariables: configVariables(id, true, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.0.with_failover", "false"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.#", "0"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/replication"),
				ConfigVariables: configVariables(id, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_to_account.0.with_failover", "true"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_StandardDatabase/replication"),
				ConfigVariables:         configVariables(id, true, true, true),
				ResourceName:            "snowflake_standard_database.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication.0.ignore_edition_check"},
			},
		},
	})
}
