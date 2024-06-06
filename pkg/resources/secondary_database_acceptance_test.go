package resources_test

import (
	"context"
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
	t.Skip("To be unskipped in the next database pr")

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
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", accountDataRetentionTimeInDays.Value),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "max_data_extension_time_in_days.0.value", accountMaxDataExtensionTimeInDays.Value),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "external_volume", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "catalog", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "replace_invalid_characters", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "default_ddl_collation", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "storage_serialization_policy", "OPTIMIZED"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "log_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "trace_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", comment),
				),
			},
			// Rename + comment update
			{
				ConfigVariables: configVariables(newId, externalPrimaryId, newComment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", accountDataRetentionTimeInDays.Value),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "max_data_extension_time_in_days.0.value", accountMaxDataExtensionTimeInDays.Value),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "external_volume", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "catalog", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "replace_invalid_characters", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "default_ddl_collation", ""),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "storage_serialization_policy", "OPTIMIZED"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "log_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "trace_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", newComment),
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
	t.Skip("To be unskipped in the next database pr")

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

	accountDataRetentionTimeInDays, err := acc.Client(t).Parameters.ShowAccountParameter(context.Background(), sdk.AccountParameterDataRetentionTimeInDays)
	require.NoError(t, err)

	accountMaxDataExtensionTimeInDays, err := acc.Client(t).Parameters.ShowAccountParameter(context.Background(), sdk.AccountParameterMaxDataExtensionTimeInDays)
	require.NoError(t, err)

	configVariables := func(
		id sdk.AccountObjectIdentifier,
		primaryDatabaseName sdk.ExternalObjectIdentifier,
		transient bool,
		dataRetentionTimeInDays *int,
		maxDataExtensionTimeInDays *int,
		externalVolume string,
		catalog string,
		replaceInvalidCharacters bool,
		defaultDdlCollation string,
		storageSerializationPolicy sdk.StorageSerializationPolicy,
		logLevel sdk.LogLevel,
		traceLevel sdk.TraceLevel,
		comment string,
	) config.Variables {
		variables := config.Variables{
			"name":                         config.StringVariable(id.Name()),
			"as_replica_of":                config.StringVariable(primaryDatabaseName.FullyQualifiedName()),
			"transient":                    config.BoolVariable(transient),
			"external_volume":              config.StringVariable(externalVolume),
			"catalog":                      config.StringVariable(catalog),
			"replace_invalid_characters":   config.BoolVariable(replaceInvalidCharacters),
			"default_ddl_collation":        config.StringVariable(defaultDdlCollation),
			"storage_serialization_policy": config.StringVariable(string(storageSerializationPolicy)),
			"log_level":                    config.StringVariable(string(logLevel)),
			"trace_level":                  config.StringVariable(string(traceLevel)),
			"comment":                      config.StringVariable(comment),
		}
		if dataRetentionTimeInDays != nil {
			variables["data_retention_time_in_days"] = config.IntegerVariable(*dataRetentionTimeInDays)
		}
		if maxDataExtensionTimeInDays != nil {
			variables["max_data_extension_time_in_days"] = config.IntegerVariable(*maxDataExtensionTimeInDays)
		}
		return variables
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
				ConfigVariables: configVariables(
					id,
					externalPrimaryId,
					false,
					sdk.Int(2),
					sdk.Int(5),
					externalVolumeId.Name(),
					catalogId.Name(),
					true,
					"en_US",
					sdk.StorageSerializationPolicyOptimized,
					sdk.LogLevelInfo,
					sdk.TraceLevelOnEvent,
					comment,
				),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", "2"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "max_data_extension_time_in_days.0.value", "5"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyOptimized)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", comment),
				),
			},
			{
				ConfigVariables: configVariables(
					newId,
					externalPrimaryId,
					false,
					nil,
					nil,
					newExternalVolumeId.Name(),
					newCatalogId.Name(),
					false,
					"en_GB",
					sdk.StorageSerializationPolicyOptimized,
					sdk.LogLevelDebug,
					sdk.TraceLevelAlways,
					newComment,
				),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", accountDataRetentionTimeInDays.Value),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "max_data_extension_time_in_days.0.value", accountMaxDataExtensionTimeInDays.Value),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "external_volume", newExternalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "catalog", newCatalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "replace_invalid_characters", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "default_ddl_collation", "en_GB"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyOptimized)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "log_level", string(sdk.LogLevelDebug)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "trace_level", string(sdk.TraceLevelAlways)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", newComment),
				),
			},
			{
				ConfigVariables: configVariables(
					id,
					externalPrimaryId,
					false,
					sdk.Int(2),
					sdk.Int(5),
					externalVolumeId.Name(),
					catalogId.Name(),
					true,
					"en_US",
					sdk.StorageSerializationPolicyCompatible,
					sdk.LogLevelInfo,
					sdk.TraceLevelOnEvent,
					comment,
				),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "as_replica_of", externalPrimaryId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", "2"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "max_data_extension_time_in_days.0.value", "5"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "comment", comment),
				),
			},
			// Import all values
			{
				ConfigVariables: configVariables(
					id,
					externalPrimaryId,
					false,
					sdk.Int(2),
					sdk.Int(5),
					externalVolumeId.Name(),
					catalogId.Name(),
					true,
					"en_US",
					sdk.StorageSerializationPolicyCompatible,
					sdk.LogLevelInfo,
					sdk.TraceLevelOnEvent,
					comment,
				),
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				ResourceName:      "snowflake_secondary_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CreateSecondaryDatabase_DataRetentionTimeInDays(t *testing.T) {
	t.Skip("To be unskipped in the next database pr")

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
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", "2"),
				),
			},
			{
				ConfigVariables: configVariables(id, externalPrimaryId, sdk.Int(1)),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", "1"),
				),
			},
			{
				ConfigVariables: configVariables(id, externalPrimaryId, nil),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", accountDataRetentionTimeInDays.Value),
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
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", "3"),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterChange()
				},
				ConfigVariables: configVariables(id, externalPrimaryId, nil),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", accountDataRetentionTimeInDays.Value),
				),
			},
			{
				ConfigVariables: configVariables(id, externalPrimaryId, sdk.Int(3)),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-set"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", "3"),
				),
			},
			{
				ConfigVariables: configVariables(id, externalPrimaryId, nil),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecondaryDatabase/complete-optionals-unset"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_secondary_database.test", "data_retention_time_in_days.0.value", accountDataRetentionTimeInDays.Value),
				),
			},
		},
	})
}
