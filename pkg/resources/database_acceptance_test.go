package resources_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/snowflakechecks"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Database_Basic(t *testing.T) {
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
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "0"),

					resource.TestCheckResourceAttrPtr("snowflake_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(newId, newComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", newComment),

					resource.TestCheckResourceAttrPtr("snowflake_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables:   configVariables(newId, newComment),
				ResourceName:      "snowflake_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_Database_ComputedValues(t *testing.T) {
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
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),

					resource.TestCheckResourceAttrPtr("snowflake_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/complete_optionals_set"),
				ConfigVariables: completeConfigVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "max_data_extension_time_in_days", "30"),
					resource.TestCheckResourceAttr("snowflake_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_database.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_database.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "enable_console_output", "true"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),

					resource.TestCheckResourceAttrPtr("snowflake_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
		},
	})
}

func TestAcc_Database_Complete(t *testing.T) {
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
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/complete_optionals_set"),
				ConfigVariables: completeConfigVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "max_data_extension_time_in_days", "30"),
					resource.TestCheckResourceAttr("snowflake_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_database.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_database.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "enable_console_output", "true"),

					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.with_failover", "true"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_Database/complete_optionals_set"),
				ConfigVariables:         completeConfigVariables,
				ResourceName:            "snowflake_database.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication.0.ignore_edition_check"},
			},
		},
	})
}

func TestAcc_Database_Update(t *testing.T) {
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
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: basicConfigVariables(id, comment),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/complete_optionals_set"),
				ConfigVariables: completeConfigVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", newComment),

					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "max_data_extension_time_in_days", "30"),
					resource.TestCheckResourceAttr("snowflake_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_database.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_database.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_database.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "enable_console_output", "true"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: basicConfigVariables(id, comment),
			},
		},
	})
}

func TestAcc_Database_HierarchicalValues(t *testing.T) {
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
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					*paramDefault = acc.TestClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterMaxDataExtensionTimeInDays).Default
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "max_data_extension_time_in_days", paramDefault),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterToDefault = acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterMaxDataExtensionTimeInDays, "50")
					t.Cleanup(revertAccountParameterToDefault)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "max_data_extension_time_in_days", "50"),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterToDefault()
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "max_data_extension_time_in_days", paramDefault),
				),
			},
		},
	})
}

func TestAcc_Database_Replication(t *testing.T) {
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
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "0"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/replication"),
				ConfigVariables: configVariables(id, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.with_failover", "true"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/replication"),
				ConfigVariables: configVariables(id, true, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.with_failover", "false"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/basic"),
				ConfigVariables: configVariables(id, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "0"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/replication"),
				ConfigVariables: configVariables(id, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.with_failover", "true"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_Database/replication"),
				ConfigVariables:         configVariables(id, true, true, true),
				ResourceName:            "snowflake_database.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication.0.ignore_edition_check"},
			},
		},
	})
}

func TestAcc_Database_IntParameter(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	databaseBasicConfig := config.Variables{
		"name": config.StringVariable(id.Name()),
	}

	databaseWithIntParameterConfig := func(dataRetentionTimeInDays int) config.Variables {
		return config.Variables{
			"name":                        config.StringVariable(id.Name()),
			"data_retention_time_in_days": config.IntegerVariable(dataRetentionTimeInDays),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			// create with setting one param
			{
				ConfigVariables: databaseWithIntParameterConfig(50),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionCreate, nil, sdk.String("50")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "50"),
				),
			},
			// do not make any change (to check if there is no drift)
			{
				ConfigVariables: databaseWithIntParameterConfig(50),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// import when param in config
			{
				ResourceName:    "snowflake_database.test",
				ImportState:     true,
				ConfigVariables: databaseWithIntParameterConfig(50),
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "data_retention_time_in_days", "50"),
				),
			},
			// change the param value in config
			{
				ConfigVariables: databaseWithIntParameterConfig(25),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("50"), sdk.String("25")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "25"),
				),
			},
			// change param value on account - expect no changes
			{
				PreConfig: func() {
					param := acc.TestClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
					require.Equal(t, "", string(param.Level))
					revert := acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterDataRetentionTimeInDays, "50")
					t.Cleanup(revert)
				},
				ConfigVariables: databaseWithIntParameterConfig(25),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionNoop, sdk.String("25"), sdk.String("25")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "25"),
				),
			},
			// change the param value externally
			{
				PreConfig: func() {
					// clean after the previous step
					acc.TestClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
					// update externally
					acc.TestClient().Database.UpdateDataRetentionTime(t, id, 50)
				},
				ConfigVariables: databaseWithIntParameterConfig(25),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectDrift("snowflake_database.test", "data_retention_time_in_days", sdk.String("25"), sdk.String("50")),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("50"), sdk.String("25")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "25"),
					snowflakechecks.CheckDatabaseDataRetentionTimeInDays(t, id, sdk.ParameterTypeDatabase, "25"),
				),
			},
			// remove the param from config
			{
				PreConfig: func() {
					param := acc.TestClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
					require.Equal(t, "", string(param.Level))
				},
				ConfigVariables: databaseBasicConfig,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/int_parameter/unset"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("25"), nil),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "1"),
					snowflakechecks.CheckDatabaseDataRetentionTimeInDays(t, id, "", "1"),
				),
			},
			// import when param not in config (snowflake default)
			{
				ResourceName:    "snowflake_database.test",
				ImportState:     true,
				ConfigVariables: databaseBasicConfig,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "data_retention_time_in_days", "1"),
				),
			},
			// change the param value in config to snowflake default
			{
				ConfigVariables: databaseWithIntParameterConfig(1),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/int_parameter/set"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("1"), nil),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "1"),
					snowflakechecks.CheckDatabaseDataRetentionTimeInDays(t, id, sdk.ParameterTypeDatabase, "1"),
				),
			},
			// remove the param from config
			{
				PreConfig: func() {
					param := acc.TestClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
					require.Equal(t, "", string(param.Level))
				},
				ConfigVariables: databaseBasicConfig,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/int_parameter/unset"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("1"), nil),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "1"), // Database default
					snowflakechecks.CheckDatabaseDataRetentionTimeInDays(t, id, "", "1"),
				),
			},
			// change param value on account - change expected to be noop
			{
				PreConfig: func() {
					param := acc.TestClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
					require.Equal(t, "", string(param.Level))
					revert := acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterDataRetentionTimeInDays, "50")
					t.Cleanup(revert)
				},
				ConfigVariables: databaseBasicConfig,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/int_parameter/unset"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectDrift("snowflake_database.test", "data_retention_time_in_days", sdk.String("1"), sdk.String("50")),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionNoop, sdk.String("50"), sdk.String("50")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "50"),
					snowflakechecks.CheckDatabaseDataRetentionTimeInDays(t, id, sdk.ParameterTypeAccount, "50"),
				),
			},
			// import when param not in config (set on account)
			{
				ResourceName:    "snowflake_database.test",
				ImportState:     true,
				ConfigVariables: databaseBasicConfig,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "data_retention_time_in_days", "50"),
				),
				Check: resource.ComposeTestCheckFunc(
					snowflakechecks.CheckDatabaseDataRetentionTimeInDays(t, id, sdk.ParameterTypeAccount, "50"),
				),
			},
			// change param value on database
			{
				PreConfig: func() {
					acc.TestClient().Database.UpdateDataRetentionTime(t, id, 50)
				},
				ConfigVariables: databaseBasicConfig,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/int_parameter/unset"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionUpdate, sdk.String("50"), nil),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "50"),
					snowflakechecks.CheckDatabaseDataRetentionTimeInDays(t, id, sdk.ParameterTypeAccount, "50"),
				),
			},
			// unset param on account
			{
				PreConfig: func() {
					acc.TestClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
				},
				ConfigVariables: databaseBasicConfig,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/int_parameter/unset"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "data_retention_time_in_days"),
						planchecks.ExpectDrift("snowflake_database.test", "data_retention_time_in_days", sdk.String("50"), sdk.String("1")),
						planchecks.ExpectChange("snowflake_database.test", "data_retention_time_in_days", tfjson.ActionNoop, sdk.String("1"), sdk.String("1")),
						planchecks.ExpectComputed("snowflake_database.test", "data_retention_time_in_days", false),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "1"),
					snowflakechecks.CheckDatabaseDataRetentionTimeInDays(t, id, "", "1"),
				),
			},
		},
	})
}

func TestAcc_Database_StringValueSetOnDifferentParameterLevelWithSameValue(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	configVariables := config.Variables{
		"name":    config.StringVariable(id.Name()),
		"catalog": config.StringVariable(catalogId.Name()),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/catalog"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "catalog", catalogId.Name()),
				),
			},
			{
				PreConfig: func() {
					require.Empty(t, acc.TestClient().Parameter.ShowAccountParameter(t, sdk.AccountParameterCatalog).Level)
					acc.TestClient().Database.UnsetCatalog(t, id)
					t.Cleanup(acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterCatalog, catalogId.Name()))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "catalog"),
						planchecks.ExpectChange("snowflake_database.test", "catalog", tfjson.ActionUpdate, sdk.String(catalogId.Name()), nil),
						planchecks.ExpectComputed("snowflake_database.test", "catalog", true),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Database/catalog"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "catalog", catalogId.Name()),
				),
			},
		},
	})
}

func TestAcc_Database_UpgradeWithTheSameFieldsAsInTheOldOne(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	dataRetentionTimeInDays := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: databaseStateUpgraderBasic(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "-1"),
				),
			},
			{
				PreConfig: func() {
					*dataRetentionTimeInDays = helpers.FindParameter(t, acc.TestClient().Parameter.ShowDatabaseParameters(t, id), sdk.AccountParameterDataRetentionTimeInDays).Value
				},
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   databaseStateUpgraderBasic(id, comment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),
					resource.TestCheckResourceAttrPtr("snowflake_database.test", "data_retention_time_in_days", dataRetentionTimeInDays),
				),
			},
		},
	})
}

func databaseStateUpgraderBasic(id sdk.AccountObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	is_transient = true
	comment = "%s"
}
`, id.Name(), comment)
}

func TestAcc_Database_UpgradeWithDataRetentionSet(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: databaseStateUpgraderDataRetentionSet(id, comment, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "10"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   databaseStateUpgraderDataRetentionSet(id, comment, 10),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_database.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_database.test", "data_retention_time_in_days", "10"),
				),
			},
		},
	})
}

func databaseStateUpgraderDataRetentionSet(id sdk.AccountObjectIdentifier, comment string, dataRetention int) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	comment = "%s"
	data_retention_time_in_days = %d
}
`, id.Name(), comment, dataRetention)
}

func TestAcc_Database_WithReplication(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	secondaryAccountLocator := acc.SecondaryClient(t).GetAccountLocator()
	secondaryAccountIdentifier := acc.SecondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: databaseStateUpgraderWithReplicationOld(id, secondaryAccountLocator),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.0.accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.0.accounts.0", secondaryAccountLocator),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   databaseStateUpgraderWithReplicationNew(id, secondaryAccountIdentifier),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails("snowflake_database.test", "replication"),
						// Updates in place (no ALTER DATABASE is called)
						planchecks.ExpectChange("snowflake_database.test", "replication.0.ignore_edition_check", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectChange("snowflake_database.test", "replication.0.enable_to_account", tfjson.ActionUpdate, sdk.String(fmt.Sprintf("[map[account_identifier:%s with_failover:false]]", secondaryAccountIdentifier)), sdk.String(fmt.Sprintf("[map[account_identifier:%s with_failover:false]]", secondaryAccountIdentifier))),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.with_failover", "false"),
					resource.TestCheckNoResourceAttr("snowflake_database.test", "replication_configuration"),
				),
			},
		},
	})
}

func databaseStateUpgraderWithReplicationOld(id sdk.AccountObjectIdentifier, enableToAccount string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	replication_configuration {
		accounts = ["%s"] 
		ignore_edition_check = true
	}
}
`, id.Name(), enableToAccount)
}

func databaseStateUpgraderWithReplicationNew(id sdk.AccountObjectIdentifier, enableToAccount string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	replication {
		enable_to_account {
			account_identifier = %s
			with_failover = false
		}
		ignore_edition_check = true
	}
}
`, id.Name(), strconv.Quote(enableToAccount))
}

func TestAcc_Database_UpgradeFromShare(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	secondaryClientLocator := acc.SecondaryClient(t).GetAccountLocator()

	shareExternalId := createShareableDatabase(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: databaseStateUpgraderFromShareOld(id, secondaryClientLocator, shareExternalId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "from_share.provider", secondaryClientLocator),
					resource.TestCheckResourceAttr("snowflake_database.test", "from_share.share", shareExternalId.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   databaseStateUpgraderFromShareNewAfterUpgrade(id),
				ExpectError:              regexp.MustCompile("failed to upgrade the state with database created from share, please use snowflake_shared_database or deprecated snowflake_database_old instead"),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   databaseStateUpgraderFromShareNew(id, shareExternalId),
				ResourceName:             "snowflake_shared_database.test",
				ImportStateId:            id.FullyQualifiedName(),
				ImportState:              true,
			},
		},
	})
}

func databaseStateUpgraderFromShareOld(id sdk.AccountObjectIdentifier, secondaryClientLocator string, externalShare sdk.ExternalObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	data_retention_time_in_days = 0 # to avoid in-place update to -1
	from_share = {
		provider = "%s"
		share = "%s"
	}
}
`, id.Name(), secondaryClientLocator, externalShare.Name())
}

func databaseStateUpgraderFromShareNewAfterUpgrade(id sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	data_retention_time_in_days = 0 # to avoid in-place update to -1
}
`, id.Name())
}

func databaseStateUpgraderFromShareNew(id sdk.AccountObjectIdentifier, externalShare sdk.ExternalObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_shared_database" "test" {
	name = "%s"
	from_share = %s
}
`, id.Name(), strconv.Quote(externalShare.FullyQualifiedName()))
}

func TestAcc_Database_UpgradeFromReplica(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, primaryDatabaseId, databaseCleanup := acc.SecondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		acc.TestClient().Account.GetAccountIdentifier(t),
	})
	t.Cleanup(databaseCleanup)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: databaseStateUpgraderFromReplicaOld(id, primaryDatabaseId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "id", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database.test", "from_replica", primaryDatabaseId.FullyQualifiedName()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   databaseStateUpgraderFromReplicaNewAfterUpgrade(id),
				ExpectError:              regexp.MustCompile("failed to upgrade the state with database created from replica, please use snowflake_secondary_database or deprecated snowflake_database_old instead"),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   databaseStateUpgraderFromReplicaNew(id, primaryDatabaseId),
				ResourceName:             "snowflake_secondary_database.test",
				ImportStateId:            id.FullyQualifiedName(),
				ImportState:              true,
			},
		},
	})
}

func databaseStateUpgraderFromReplicaOld(id sdk.AccountObjectIdentifier, primaryDatabaseId sdk.ExternalObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	data_retention_time_in_days = 0 # to avoid in-place update to -1
	from_replica = %s
}
`, id.Name(), strconv.Quote(primaryDatabaseId.FullyQualifiedName()))
}

func databaseStateUpgraderFromReplicaNewAfterUpgrade(id sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	data_retention_time_in_days = 0
}
`, id.Name())
}

func databaseStateUpgraderFromReplicaNew(id sdk.AccountObjectIdentifier, primaryDatabaseId sdk.ExternalObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_secondary_database" "test" {
	name = "%s"
	as_replica_of = %s
}
`, id.Name(), strconv.Quote(id.FullyQualifiedName()))
}

func TestAcc_Database_UpgradeFromClonedDatabase(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	cloneId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: databaseStateUpgraderFromDatabaseOld(id, cloneId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.cloned", "id", cloneId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.cloned", "name", cloneId.Name()),
					resource.TestCheckResourceAttr("snowflake_database.cloned", "from_database", id.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   databaseStateUpgraderFromDatabaseNewAfterUpgrade(id, cloneId),
				ExpectError:              regexp.MustCompile("failed to upgrade the state with database created from database, please use snowflake_database or deprecated snowflake_database_old instead"),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   databaseStateUpgraderFromDatabaseNew(id, cloneId),
				ResourceName:             "snowflake_database.cloned",
				ImportStateId:            cloneId.FullyQualifiedName(),
				ImportState:              true,
			},
		},
	})
}

func databaseStateUpgraderFromDatabaseOld(id sdk.AccountObjectIdentifier, secondId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	data_retention_time_in_days = 0 # to avoid in-place update to -1
}

resource "snowflake_database" "cloned" {
	name = "%s"
	data_retention_time_in_days = 0 # to avoid in-place update to -1
	from_database = snowflake_database.test.name
}
`, id.Name(), secondId.Name())
}

func databaseStateUpgraderFromDatabaseNewAfterUpgrade(id sdk.AccountObjectIdentifier, secondId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
	data_retention_time_in_days = 0
}

resource "snowflake_database" "cloned" {
	name = "%s"
	data_retention_time_in_days = 0
}
`, id.Name(), secondId.Name())
}

func databaseStateUpgraderFromDatabaseNew(id sdk.AccountObjectIdentifier, secondId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%s"
}

resource "snowflake_database" "cloned" {
	name = "%s"
}
`, id.Name(), secondId.Name())
}
