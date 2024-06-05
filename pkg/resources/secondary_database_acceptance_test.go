package resources_test

import (
	"context"
	"slices"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_CreateSecondaryDatabase_minimal(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	_, externalPrimaryId, primaryDatabaseCleanup := acc.SecondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		acc.TestClient().Account.GetAccountIdentifier(t),
	})
	t.Cleanup(primaryDatabaseCleanup)

	newId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newComment := random.Comment()

	accountDataRetentionTimeInDays, err := acc.Client(t).Parameters.ShowAccountParameter(context.Background(), sdk.AccountParameterDataRetentionTimeInDays)
	require.NoError(t, err)

	accountMaxDataExtensionTimeInDays, err := acc.Client(t).Parameters.ShowAccountParameter(context.Background(), sdk.AccountParameterMaxDataExtensionTimeInDays)
	require.NoError(t, err)

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
				ConfigVariables: configVariables(id, externalPrimaryId, comment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays.Value),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays.Value),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "external_volume", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "catalog", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "replace_invalid_characters", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "default_ddl_collation", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "storage_serialization_policy", "OPTIMIZED"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "log_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "trace_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "suspend_task_after_num_failures", "10"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "task_auto_retry_attempts", "0"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_managed_initial_warehouse_size", "Medium"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_timeout_ms", "3600000"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_minimum_trigger_interval_in_seconds", "30"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "quoted_identifiers_ignore_case", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "enable_console_output", "false"),
				),
			},
			// Rename + comment update
			{
				ConfigVariables: configVariables(newId, externalPrimaryId, newComment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", newComment),

					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days", accountDataRetentionTimeInDays.Value),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "max_data_extension_time_in_days", accountMaxDataExtensionTimeInDays.Value),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "external_volume", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "catalog", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "replace_invalid_characters", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "default_ddl_collation", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "storage_serialization_policy", "OPTIMIZED"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "log_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "trace_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "suspend_task_after_num_failures", "10"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "task_auto_retry_attempts", "0"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_managed_initial_warehouse_size", "Medium"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_timeout_ms", "3600000"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_minimum_trigger_interval_in_seconds", "30"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "quoted_identifiers_ignore_case", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "enable_console_output", "false"),
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

	_, externalPrimaryId, primaryDatabaseCleanup := acc.SecondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		sdk.NewAccountIdentifierFromAccountLocator(acc.Client(t).GetAccountLocator()),
	})
	t.Cleanup(primaryDatabaseCleanup)

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

	params, err := acc.Client(t).Parameters.ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Account: sdk.Bool(true),
		},
	})
	require.NoError(t, err)

	findParamValue := func(searchedParameter sdk.AccountParameter) string {
		idx := slices.IndexFunc(params, func(parameter *sdk.Parameter) bool {
			return parameter.Key == string(searchedParameter)
		})
		require.NotEqual(t, -1, idx, string(searchedParameter))
		return params[idx].Value
	}

	accountDataRetentionTimeInDays := findParamValue(sdk.AccountParameterDataRetentionTimeInDays)
	accountMaxDataExtensionTimeInDays := findParamValue(sdk.AccountParameterMaxDataExtensionTimeInDays)
	accountExternalVolume := findParamValue(sdk.AccountParameterExternalVolume)
	accountCatalog := findParamValue(sdk.AccountParameterCatalog)
	accountReplaceInvalidCharacters := findParamValue(sdk.AccountParameterReplaceInvalidCharacters)
	accountDefaultDdlCollation := findParamValue(sdk.AccountParameterDefaultDDLCollation)
	accountStorageSerializationPolicy := findParamValue(sdk.AccountParameterStorageSerializationPolicy)
	accountLogLevel := findParamValue(sdk.AccountParameterLogLevel)
	accountTraceLevel := findParamValue(sdk.AccountParameterTraceLevel)
	accountSuspendTaskAfterNumFailures := findParamValue(sdk.AccountParameterSuspendTaskAfterNumFailures)
	accountTaskAutoRetryAttempts := findParamValue(sdk.AccountParameterTaskAutoRetryAttempts)
	accountUserTaskMangedInitialWarehouseSize := findParamValue(sdk.AccountParameterUserTaskManagedInitialWarehouseSize)
	accountUserTaskTimeoutMs := findParamValue(sdk.AccountParameterUserTaskTimeoutMs)
	accountUserTaskMinimumTriggerIntervalInSeconds := findParamValue(sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds)
	accountQuotedIdentifiersIgnoreCase := findParamValue(sdk.AccountParameterQuotedIdentifiersIgnoreCase)
	accountEnableConsoleOutput := findParamValue(sdk.AccountParameterEnableConsoleOutput)

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
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", newComment),

					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", accountDataRetentionTimeInDays),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "max_data_extension_time_in_days.0.value", accountMaxDataExtensionTimeInDays),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "enable_console_output", accountEnableConsoleOutput),
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

	_, externalPrimaryId, primaryDatabaseCleanup := acc.SecondaryTestClient().Database.CreatePrimaryDatabase(t, []sdk.AccountIdentifier{
		sdk.NewAccountIdentifierFromAccountLocator(acc.Client(t).GetAccountLocator()),
	})
	t.Cleanup(primaryDatabaseCleanup)

	accountDataRetentionTimeInDays, err := acc.Client(t).Parameters.ShowAccountParameter(context.Background(), sdk.AccountParameterDataRetentionTimeInDays)
	require.NoError(t, err)

	configVariables := func(
		id sdk.AccountObjectIdentifier,
		primaryDatabaseName sdk.ExternalObjectIdentifier,
		dataRetentionTimeInDays *int,
	) config.Variables {
		variables := config.Variables{
			"name":                         config.StringVariable(id.Name()),
			"as_replica_of":                config.StringVariable(primaryDatabaseName.FullyQualifiedName()),
			"transient":                    config.BoolVariable(false),
			"external_volume":              config.StringVariable(""),
			"catalog":                      config.StringVariable(""),
			"replace_invalid_characters":   config.StringVariable("false"),
			"default_ddl_collation":        config.StringVariable(""),
			"storage_serialization_policy": config.StringVariable("OPTIMIZED"),
			"log_level":                    config.StringVariable("OFF"),
			"trace_level":                  config.StringVariable("OFF"),
			"comment":                      config.StringVariable(""),
		}
		if dataRetentionTimeInDays != nil {
			variables["data_retention_time_in_days"] = config.IntegerVariable(*dataRetentionTimeInDays)
			variables["max_data_extension_time_in_days"] = config.IntegerVariable(10)
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
