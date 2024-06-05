package resources_test

import (
	"context"
	"slices"
	"strconv"
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

func TestAcc_StandardDatabase_Minimal(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	newId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newComment := random.Comment()

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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "data_retention_time_in_days.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "external_volume.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "catalog.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replace_invalid_characters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "default_ddl_collation.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "storage_serialization_policy.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "log_level.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "trace_level.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.#", "0"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(newId, newComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", newComment),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "data_retention_time_in_days.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "external_volume.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "catalog.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replace_invalid_characters.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "default_ddl_collation.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "storage_serialization_policy.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "log_level.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "trace_level.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.#", "0"),
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

	completeConfigVariables := func(
		id sdk.AccountObjectIdentifier,
		comment string,
		dataRetention int,
		maxDataExtension int,
		replaceInvalidCharacters bool,
		defaultDdlCollation string,
		storageSerializationPolicy sdk.StorageSerializationPolicy,
		logLevel sdk.LogLevel,
		traceLevel sdk.TraceLevel,
	) config.Variables {
		return config.Variables{
			"name":                            config.StringVariable(id.Name()),
			"comment":                         config.StringVariable(comment),
			"transient":                       config.BoolVariable(false),
			"data_retention_time_in_days":     config.IntegerVariable(dataRetention),
			"max_data_extension_time_in_days": config.IntegerVariable(maxDataExtension),
			"external_volume":                 config.StringVariable(externalVolumeId.Name()),
			"catalog":                         config.StringVariable(catalogId.Name()),
			"replace_invalid_characters":      config.BoolVariable(replaceInvalidCharacters),
			"default_ddl_collation":           config.StringVariable(defaultDdlCollation),
			"storage_serialization_policy":    config.StringVariable(string(storageSerializationPolicy)),
			"log_level":                       config.StringVariable(string(logLevel)),
			"trace_level":                     config.StringVariable(string(traceLevel)),
			"account_identifier":              config.StringVariable(secondaryAccountIdentifier),
			"with_failover":                   config.BoolVariable(true),
			"ignore_edition_check":            config.BoolVariable(true),
		}
	}

	var (
		dataRetentionTimeInDays    = new(string)
		maxDataExtensionTimeInDays = new(string)
		externalVolume             = new(string)
		catalog                    = new(string)
		replaceInvalidCharacters   = new(string)
		defaultDdlCollation        = new(string)
		storageSerializationPolicy = new(string)
		logLevel                   = new(string)
		traceLevel                 = new(string)
	)

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
				ConfigVariables: configVariables(id, comment),
			},
			{
				PreConfig: func() {
					params, err := acc.Client(t).Parameters.ShowParameters(context.Background(), &sdk.ShowParametersOptions{
						In: &sdk.ParametersIn{
							Database: id,
						},
					})
					require.NoError(t, err)

					findParam := func(key string) string {
						idx := slices.IndexFunc(params, func(parameter *sdk.Parameter) bool { return parameter.Key == key })
						require.NotEqual(t, -1, idx)
						return params[idx].Value
					}

					*dataRetentionTimeInDays = findParam("DATA_RETENTION_TIME_IN_DAYS")
					*maxDataExtensionTimeInDays = findParam("MAX_DATA_EXTENSION_TIME_IN_DAYS")
					*externalVolume = findParam("EXTERNAL_VOLUME")
					*catalog = findParam("CATALOG")
					*replaceInvalidCharacters = findParam("REPLACE_INVALID_CHARACTERS")
					*defaultDdlCollation = findParam("DEFAULT_DDL_COLLATION")
					*storageSerializationPolicy = findParam("STORAGE_SERIALIZATION_POLICY")
					*logLevel = findParam("LOG_LEVEL")
					*traceLevel = findParam("TRACE_LEVEL")
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", comment),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "data_retention_time_in_days.0.value", dataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "max_data_extension_time_in_days.0.value", maxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "external_volume.0.value", externalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "catalog.0.value", catalog),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "replace_invalid_characters.0.value", replaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "default_ddl_collation.0.value", defaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "storage_serialization_policy.0.value", storageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "log_level.0.value", logLevel),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "trace_level.0.value", traceLevel),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/complete-optionals-set"),
				ConfigVariables: completeConfigVariables(id, comment, 20, 30, true, "en_US", sdk.StorageSerializationPolicyCompatible, sdk.LogLevelInfo, sdk.TraceLevelOnEvent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "data_retention_time_in_days.0.value", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days.0.value", "30"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "external_volume.0.value", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "catalog.0.value", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replace_invalid_characters.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "default_ddl_collation.0.value", "en_US"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "storage_serialization_policy.0.value", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "log_level.0.value", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "trace_level.0.value", string(sdk.TraceLevelOnEvent)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", comment),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "data_retention_time_in_days.0.value", dataRetentionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "max_data_extension_time_in_days.0.value", maxDataExtensionTimeInDays),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "external_volume.0.value", externalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "catalog.0.value", catalog),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "replace_invalid_characters.0.value", replaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "default_ddl_collation.0.value", defaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "storage_serialization_policy.0.value", storageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "log_level.0.value", logLevel),
					resource.TestCheckResourceAttrPtr("snowflake_standard_database.test", "trace_level.0.value", traceLevel),
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

	configVariables := func(
		id sdk.AccountObjectIdentifier,
		comment string,
		dataRetention int,
		maxDataExtension int,
		replaceInvalidCharacters bool,
		defaultDdlCollation string,
		storageSerializationPolicy sdk.StorageSerializationPolicy,
		logLevel sdk.LogLevel,
		traceLevel sdk.TraceLevel,
		withFailover bool,
		ignoreEditionCheck bool,
	) config.Variables {
		return config.Variables{
			"name":                            config.StringVariable(id.Name()),
			"comment":                         config.StringVariable(comment),
			"transient":                       config.BoolVariable(false),
			"data_retention_time_in_days":     config.IntegerVariable(dataRetention),
			"max_data_extension_time_in_days": config.IntegerVariable(maxDataExtension),
			"external_volume":                 config.StringVariable(externalVolumeId.Name()),
			"catalog":                         config.StringVariable(catalogId.Name()),
			"replace_invalid_characters":      config.BoolVariable(replaceInvalidCharacters),
			"default_ddl_collation":           config.StringVariable(defaultDdlCollation),
			"storage_serialization_policy":    config.StringVariable(string(storageSerializationPolicy)),
			"log_level":                       config.StringVariable(string(logLevel)),
			"trace_level":                     config.StringVariable(string(traceLevel)),
			"account_identifier":              config.StringVariable(secondaryAccountIdentifier),
			"with_failover":                   config.BoolVariable(withFailover),
			"ignore_edition_check":            config.BoolVariable(ignoreEditionCheck),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/complete-optionals-set"),
				ConfigVariables: configVariables(id, comment, 20, 30, true, "en_US", sdk.StorageSerializationPolicyCompatible, sdk.LogLevelInfo, sdk.TraceLevelOnEvent, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "data_retention_time_in_days.0.value", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days.0.value", "30"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "external_volume.0.value", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "catalog.0.value", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replace_invalid_characters.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "default_ddl_collation.0.value", "en_US"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "storage_serialization_policy.0.value", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "log_level.0.value", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "trace_level.0.value", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.0.with_failover", "true"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_StandardDatabase/complete-optionals-set"),
				ConfigVariables:         configVariables(id, comment, 20, 30, true, "en_US", sdk.StorageSerializationPolicyCompatible, sdk.LogLevelInfo, sdk.TraceLevelOnEvent, true, true),
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
	secondaryAccountIdentifier := acc.SecondaryTestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()
	comment := random.Comment()

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	configVariables := func(
		id sdk.AccountObjectIdentifier,
		comment string,
		dataRetention int,
		maxDataExtension int,
		replaceInvalidCharacters bool,
		defaultDdlCollation string,
		storageSerializationPolicy sdk.StorageSerializationPolicy,
		logLevel sdk.LogLevel,
		traceLevel sdk.TraceLevel,
		withFailover bool,
		ignoreEditionCheck bool,
	) config.Variables {
		return config.Variables{
			"name":                            config.StringVariable(id.Name()),
			"comment":                         config.StringVariable(comment),
			"transient":                       config.BoolVariable(false),
			"data_retention_time_in_days":     config.IntegerVariable(dataRetention),
			"max_data_extension_time_in_days": config.IntegerVariable(maxDataExtension),
			"external_volume":                 config.StringVariable(externalVolumeId.Name()),
			"catalog":                         config.StringVariable(catalogId.Name()),
			"replace_invalid_characters":      config.BoolVariable(replaceInvalidCharacters),
			"default_ddl_collation":           config.StringVariable(defaultDdlCollation),
			"storage_serialization_policy":    config.StringVariable(string(storageSerializationPolicy)),
			"log_level":                       config.StringVariable(string(logLevel)),
			"trace_level":                     config.StringVariable(string(traceLevel)),
			"account_identifier":              config.StringVariable(secondaryAccountIdentifier),
			"with_failover":                   config.BoolVariable(withFailover),
			"ignore_edition_check":            config.BoolVariable(ignoreEditionCheck),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/complete-optionals-set"),
				ConfigVariables: configVariables(id, comment, 20, 30, true, "en_US", sdk.StorageSerializationPolicyCompatible, sdk.LogLevelInfo, sdk.TraceLevelOnEvent, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "is_transient", "false"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "data_retention_time_in_days.0.value", "20"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days.0.value", "30"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "external_volume.0.value", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "catalog.0.value", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replace_invalid_characters.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "default_ddl_collation.0.value", "en_US"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "storage_serialization_policy.0.value", string(sdk.StorageSerializationPolicyCompatible)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "log_level.0.value", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "trace_level.0.value", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.0.with_failover", "true"),
				),
			},
			{
				ConfigDirectory:         acc.ConfigurationDirectory("TestAcc_StandardDatabase/complete-optionals-set"),
				ConfigVariables:         configVariables(id, comment, 20, 30, true, "en_US", sdk.StorageSerializationPolicyCompatible, sdk.LogLevelInfo, sdk.TraceLevelOnEvent, true, true),
				ResourceName:            "snowflake_standard_database.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication.0.ignore_edition_check"},
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

	param, err := acc.Client(t).Parameters.ShowAccountParameter(context.Background(), sdk.AccountParameterMaxDataExtensionTimeInDays)
	require.NoError(t, err)

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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days.0.value", param.Default),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterToDefault = acc.TestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterMaxDataExtensionTimeInDays, strconv.Itoa(50))
					t.Cleanup(revertAccountParameterToDefault)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days.0.value", "50"),
				),
			},
			{
				PreConfig: func() {
					revertAccountParameterToDefault()
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/basic"),
				ConfigVariables: configVariables(id, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "max_data_extension_time_in_days.0.value", param.Default),
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
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.0.with_failover", "true"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StandardDatabase/replication"),
				ConfigVariables: configVariables(id, true, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.ignore_edition_check", "true"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.0.with_failover", "false"),
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
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.#", "1"),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.0.account_identifier", secondaryAccountIdentifier),
					resource.TestCheckResourceAttr("snowflake_standard_database.test", "replication.0.enable_for_account.0.with_failover", "true"),
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
