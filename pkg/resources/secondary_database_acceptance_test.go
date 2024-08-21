package resources_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_CreateSecondaryDatabase_Basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	primaryDatabase, externalPrimaryId, _ := acc.SecondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		acc.TestClient().Account.GetAccountIdentifier(t),
	})
	t.Cleanup(func() {
		// TODO(SNOW-1562172): Create a better solution for this type of situations
		require.Eventually(t, func() bool { return acc.SecondaryTestClient().Database.DropDatabase(t, primaryDatabase.ID()) == nil }, time.Second*5, time.Second)
	})

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

	configVariables := func(id sdk.AccountObjectIdentifier, primaryDatabaseName sdk.ExternalObjectIdentifier, comment string) config.Variables {
		return config.Variables{
			"name":          config.StringVariable(id.Name()),
			"as_replica_of": config.StringVariable(primaryDatabaseName.FullyQualifiedName()),
			"comment":       config.StringVariable(comment),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.SharedDatabase),
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
				ConfigVariables: configVariables(id, externalPrimaryId, comment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", comment),

					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			// Rename + comment update
			{
				ConfigVariables: configVariables(newId, externalPrimaryId, newComment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", newComment),

					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			// Import all values
			{
				ConfigVariables:   configVariables(newId, externalPrimaryId, newComment),
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/basic"),
				ResourceName:      "snowflake_secondary_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CreateSecondaryDatabase_complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	primaryDatabase, externalPrimaryId, _ := acc.SecondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		sdk.NewAccountIdentifierFromAccountLocator(acc.Client(t).GetAccountLocator()),
	})
	t.Cleanup(func() {
		// TODO(SNOW-1562172): Create a better solution for this type of situations
		require.Eventually(t, func() bool { return acc.SecondaryTestClient().Database.DropDatabase(t, primaryDatabase.ID()) == nil }, time.Second*5, time.Second)
	})

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	newId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newComment := random.Comment()

	newExternalVolumeId, newExternalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(newExternalVolumeCleanup)

	newCatalogId, newCatalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(newCatalogCleanup)

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

	unsetConfigVariables := config.Variables{
		"name":          config.StringVariable(id.Name()),
		"as_replica_of": config.StringVariable(externalPrimaryId.FullyQualifiedName()),
	}

	setConfigVariables := config.Variables{
		"name":          config.StringVariable(id.Name()),
		"as_replica_of": config.StringVariable(externalPrimaryId.FullyQualifiedName()),
		"comment":       config.StringVariable(comment),

		"data_retention_time_in_days":                   config.IntegerVariable(20),
		"max_data_extension_time_in_days":               config.IntegerVariable(25),
		"external_volume":                               config.StringVariable(externalVolumeId.Name()),
		"catalog":                                       config.StringVariable(catalogId.Name()),
		"replace_invalid_characters":                    config.BoolVariable(true),
		"default_ddl_collation":                         config.StringVariable("en_US"),
		"storage_serialization_policy":                  config.StringVariable(string(sdk.StorageSerializationPolicyCompatible)),
		"log_level":                                     config.StringVariable(string(sdk.LogLevelDebug)),
		"trace_level":                                   config.StringVariable(string(sdk.TraceLevelAlways)),
		"suspend_task_after_num_failures":               config.IntegerVariable(20),
		"task_auto_retry_attempts":                      config.IntegerVariable(20),
		"user_task_managed_initial_warehouse_size":      config.StringVariable(string(sdk.WarehouseSizeLarge)),
		"user_task_timeout_ms":                          config.IntegerVariable(1200000),
		"user_task_minimum_trigger_interval_in_seconds": config.IntegerVariable(60),
		"quoted_identifiers_ignore_case":                config.BoolVariable(true),
		"enable_console_output":                         config.BoolVariable(true),
	}

	updatedConfigVariables := config.Variables{
		"name":          config.StringVariable(newId.Name()),
		"as_replica_of": config.StringVariable(externalPrimaryId.FullyQualifiedName()),
		"comment":       config.StringVariable(newComment),

		"data_retention_time_in_days":                   config.IntegerVariable(40),
		"max_data_extension_time_in_days":               config.IntegerVariable(45),
		"external_volume":                               config.StringVariable(newExternalVolumeId.Name()),
		"catalog":                                       config.StringVariable(newCatalogId.Name()),
		"replace_invalid_characters":                    config.BoolVariable(false),
		"default_ddl_collation":                         config.StringVariable("en_GB"),
		"storage_serialization_policy":                  config.StringVariable(string(sdk.StorageSerializationPolicyOptimized)),
		"log_level":                                     config.StringVariable(string(sdk.LogLevelInfo)),
		"trace_level":                                   config.StringVariable(string(sdk.TraceLevelOnEvent)),
		"suspend_task_after_num_failures":               config.IntegerVariable(40),
		"task_auto_retry_attempts":                      config.IntegerVariable(40),
		"user_task_managed_initial_warehouse_size":      config.StringVariable(string(sdk.WarehouseSizeXLarge)),
		"user_task_timeout_ms":                          config.IntegerVariable(2400000),
		"user_task_minimum_trigger_interval_in_seconds": config.IntegerVariable(120),
		"quoted_identifiers_ignore_case":                config.BoolVariable(false),
		"enable_console_output":                         config.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.SecondaryDatabase),
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
				ConfigVariables: setConfigVariables,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "max_data_extension_time_in_days", "25"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "log_level", string(sdk.LogLevelDebug)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "trace_level", string(sdk.TraceLevelAlways)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_managed_initial_warehouse_size", "LARGE"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_minimum_trigger_interval_in_seconds", "60"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "enable_console_output", "true"),
				),
			},
			{
				ConfigVariables: updatedConfigVariables,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", newComment),

					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", "40"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "max_data_extension_time_in_days", "45"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "external_volume", newExternalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "catalog", newCatalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "replace_invalid_characters", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "default_ddl_collation", "en_GB"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyOptimized)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "suspend_task_after_num_failures", "40"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "task_auto_retry_attempts", "40"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_managed_initial_warehouse_size", "XLARGE"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_timeout_ms", "2400000"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "quoted_identifiers_ignore_case", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "enable_console_output", "false"),
				),
			},
			{
				ConfigVariables: unsetConfigVariables,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", ""),

					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_secondary_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			{
				ConfigVariables: setConfigVariables,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", "20"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "max_data_extension_time_in_days", "25"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "log_level", string(sdk.LogLevelDebug)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "trace_level", string(sdk.TraceLevelAlways)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_managed_initial_warehouse_size", "LARGE"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_minimum_trigger_interval_in_seconds", "60"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "enable_console_output", "true"),
				),
			},
			// Import all values
			{
				ConfigVariables:   setConfigVariables,
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				ResourceName:      "snowflake_secondary_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CreateSecondaryDatabase_DataRetentionTimeInDays(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	primaryDatabase, externalPrimaryId, _ := acc.SecondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		sdk.NewAccountIdentifierFromAccountLocator(acc.Client(t).GetAccountLocator()),
	})
	t.Cleanup(func() {
		// TODO(SNOW-1562172): Create a better solution for this type of situations
		require.Eventually(t, func() bool { return acc.SecondaryTestClient().Database.DropDatabase(t, primaryDatabase.ID()) == nil }, time.Second*5, time.Second)
	})

	accountDataRetentionTimeInDays, err := acc.Client(t).Parameters.ShowAccountParameter(context.Background(), sdk.AccountParameterDataRetentionTimeInDays)
	require.NoError(t, err)

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	configVariables := func(
		id sdk.AccountObjectIdentifier,
		primaryDatabaseName sdk.ExternalObjectIdentifier,
		dataRetentionTimeInDays *int,
	) config.Variables {
		variables := config.Variables{
			"name":          config.StringVariable(id.Name()),
			"as_replica_of": config.StringVariable(primaryDatabaseName.FullyQualifiedName()),
			"transient":     config.BoolVariable(false),
			"comment":       config.StringVariable(""),

			"max_data_extension_time_in_days":               config.IntegerVariable(10),
			"external_volume":                               config.StringVariable(externalVolumeId.Name()),
			"catalog":                                       config.StringVariable(catalogId.Name()),
			"replace_invalid_characters":                    config.BoolVariable(true),
			"default_ddl_collation":                         config.StringVariable("en_US"),
			"storage_serialization_policy":                  config.StringVariable("OPTIMIZED"),
			"log_level":                                     config.StringVariable("OFF"),
			"trace_level":                                   config.StringVariable("OFF"),
			"suspend_task_after_num_failures":               config.IntegerVariable(10),
			"task_auto_retry_attempts":                      config.IntegerVariable(10),
			"user_task_managed_initial_warehouse_size":      config.StringVariable(string(sdk.WarehouseSizeSmall)),
			"user_task_timeout_ms":                          config.IntegerVariable(120000),
			"user_task_minimum_trigger_interval_in_seconds": config.IntegerVariable(120),
			"quoted_identifiers_ignore_case":                config.BoolVariable(true),
			"enable_console_output":                         config.BoolVariable(true),
		}
		if dataRetentionTimeInDays != nil {
			variables["data_retention_time_in_days"] = config.IntegerVariable(*dataRetentionTimeInDays)
		}
		return variables
	}

	var revertAccountParameterChange func()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.SecondaryDatabase),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(id, externalPrimaryId, sdk.Int(2)),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", "2"),
				),
			},
			{
				ConfigVariables: configVariables(id, externalPrimaryId, sdk.Int(1)),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", "1"),
				),
			},
			{
				ConfigVariables: configVariables(id, externalPrimaryId, nil),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays.Value),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterChange = acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterDataRetentionTimeInDays, "3")
					t.Cleanup(revertAccountParameterChange)
				},
				ConfigVariables: configVariables(id, externalPrimaryId, nil),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", "3"),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterChange()
				},
				ConfigVariables: configVariables(id, externalPrimaryId, nil),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays.Value),
				),
			},
			{
				ConfigVariables: configVariables(id, externalPrimaryId, sdk.Int(3)),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", "3"),
				),
			},
			{
				ConfigVariables: configVariables(id, externalPrimaryId, nil),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays.Value),
				),
			},
		},
	})
}

func TestAcc_SecondaryDatabase_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	primaryDatabase, externalPrimaryId, _ := acc.SecondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		sdk.NewAccountIdentifierFromAccountLocator(acc.Client(t).GetAccountLocator()),
	})
	t.Cleanup(func() {
		// TODO(SNOW-1562172): Create a better solution for this type of situations
		require.Eventually(t, func() bool { return acc.SecondaryTestClient().Database.DropDatabase(t, primaryDatabase.ID()) == nil }, time.Second*5, time.Second)
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecondaryDatabase),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: secondaryDatabaseConfigBasic(id.Name(), externalPrimaryId.FullyQualifiedName()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "id", id.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   secondaryDatabaseConfigBasic(id.Name(), externalPrimaryId.FullyQualifiedName()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

func secondaryDatabaseConfigBasic(name, externalDatabaseId string) string {
	return fmt.Sprintf(`resource "snowflake_secondary_database" "test" {
		name = "%v"
		as_replica_of = %v
	}`, name, strconv.Quote(externalDatabaseId))
}

func TestAcc_SecondaryDatabase_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())

	primaryDatabase, externalPrimaryId, _ := acc.SecondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		sdk.NewAccountIdentifierFromAccountLocator(acc.Client(t).GetAccountLocator()),
	})
	unquotedExternalPrimaryId := fmt.Sprintf("%s.%s.%s", externalPrimaryId.AccountIdentifier().OrganizationName(), externalPrimaryId.AccountIdentifier().AccountName(), externalPrimaryId.Name())
	t.Cleanup(func() {
		// TODO(SNOW-1562172): Create a better solution for this type of situations
		require.Eventually(t, func() bool { return acc.SecondaryTestClient().Database.DropDatabase(t, primaryDatabase.ID()) == nil }, time.Second*5, time.Second)
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecondaryDatabase),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				ExpectNonEmptyPlan: true,
				Config:             secondaryDatabaseConfigBasic(quotedId, unquotedExternalPrimaryId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "id", id.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   secondaryDatabaseConfigBasic(quotedId, unquotedExternalPrimaryId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_secondary_database.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_secondary_database.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}
