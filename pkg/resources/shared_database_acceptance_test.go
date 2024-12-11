package resources_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CreateSharedDatabase_Basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	newId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newComment := random.Comment()

	var (
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

	configVariables := func(id sdk.AccountObjectIdentifier, shareName sdk.ExternalObjectIdentifier, comment string) config.Variables {
		return config.Variables{
			"name":       config.StringVariable(id.Name()),
			"from_share": config.StringVariable(shareName.FullyQualifiedName()),
			"comment":    config.StringVariable(comment),
		}
	}

	shareExternalId := createShareableDatabase(t)

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
				ConfigVariables: configVariables(id, shareExternalId, comment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SharedDatabase/basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "from_share", shareExternalId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "comment", comment),

					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			{
				ConfigVariables: configVariables(newId, shareExternalId, newComment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SharedDatabase/basic"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_shared_database.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "from_share", shareExternalId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "comment", newComment),

					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr("snowflake_shared_database.test", "enable_console_output", accountEnableConsoleOutput),
				),
			},
			// Import all values
			{
				ConfigVariables:   configVariables(newId, shareExternalId, newComment),
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_SharedDatabase/basic"),
				ResourceName:      "snowflake_shared_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CreateSharedDatabase_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	externalShareId := createShareableDatabase(t)

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	configVariables := config.Variables{
		"name":       config.StringVariable(id.Name()),
		"from_share": config.StringVariable(externalShareId.FullyQualifiedName()),
		"comment":    config.StringVariable(comment),

		"external_volume":                               config.StringVariable(externalVolumeId.Name()),
		"catalog":                                       config.StringVariable(catalogId.Name()),
		"replace_invalid_characters":                    config.BoolVariable(true),
		"default_ddl_collation":                         config.StringVariable("en_US"),
		"storage_serialization_policy":                  config.StringVariable(string(sdk.StorageSerializationPolicyOptimized)),
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

	acc.TestClient().Database.CreateDatabaseFromShareTemporarily(t, externalShareId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.SharedDatabase),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SharedDatabase/complete"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "from_share", externalShareId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_shared_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyOptimized)),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "enable_console_output", "true"),
				),
			},
			// Import all values
			{
				ConfigVariables:   configVariables,
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_SharedDatabase/complete"),
				ResourceName:      "snowflake_shared_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CreateSharedDatabase_InvalidValues(t *testing.T) {
	comment := random.Comment()

	configVariables := config.Variables{
		"name":       config.StringVariable("name"),
		"from_share": config.StringVariable("org.acc.name"),
		"comment":    config.StringVariable(comment),

		"external_volume":                               config.StringVariable(""),
		"catalog":                                       config.StringVariable(""),
		"replace_invalid_characters":                    config.BoolVariable(false),
		"default_ddl_collation":                         config.StringVariable(""),
		"storage_serialization_policy":                  config.StringVariable("invalid_value"),
		"log_level":                                     config.StringVariable("invalid_value"),
		"trace_level":                                   config.StringVariable("invalid_value"),
		"suspend_task_after_num_failures":               config.IntegerVariable(0),
		"task_auto_retry_attempts":                      config.IntegerVariable(0),
		"user_task_managed_initial_warehouse_size":      config.StringVariable(""),
		"user_task_timeout_ms":                          config.IntegerVariable(0),
		"user_task_minimum_trigger_interval_in_seconds": config.IntegerVariable(0),
		"quoted_identifiers_ignore_case":                config.BoolVariable(false),
		"enable_console_output":                         config.BoolVariable(false),
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
				ConfigVariables: configVariables,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SharedDatabase/complete"),
				ExpectError: regexp.MustCompile(`(unknown log level: invalid_value)|` +
					`(unknown trace level: invalid_value)|` +
					`(unknown storage serialization policy: invalid_value)|` +
					`(invalid warehouse size:)`),
			},
		},
	})
}

// createShareableDatabase creates a database on the secondary account and enables database sharing on the primary account.
// TODO(SNOW-1431726): Later on, this function should be moved to more sophisticated helpers.
func createShareableDatabase(t *testing.T) sdk.ExternalObjectIdentifier {
	t.Helper()

	share, shareCleanup := acc.SecondaryTestClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	sharedDatabase, sharedDatabaseCleanup := acc.SecondaryTestClient().Database.CreateDatabase(t)
	t.Cleanup(sharedDatabaseCleanup)

	revoke := acc.SecondaryTestClient().Grant.GrantPrivilegeOnDatabaseToShare(t, sharedDatabase.ID(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage})
	t.Cleanup(revoke)

	acc.SecondaryTestClient().Share.SetAccountOnShare(t, acc.TestClient().Account.GetAccountIdentifier(t), share.ID())

	return sdk.NewExternalObjectIdentifier(acc.SecondaryTestClient().Account.GetAccountIdentifier(t), share.ID())
}

func TestAcc_SharedDatabase_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	externalShareId := createShareableDatabase(t)

	acc.TestClient().Database.CreateDatabaseFromShareTemporarily(t, externalShareId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SharedDatabase),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: sharedDatabaseConfigBasic(id.Name(), externalShareId.FullyQualifiedName()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   sharedDatabaseConfigBasic(id.Name(), externalShareId.FullyQualifiedName()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "id", id.Name()),
				),
			},
		},
	})
}

func sharedDatabaseConfigBasic(name, externalShareId string) string {
	return fmt.Sprintf(`resource "snowflake_shared_database" "test" {
		name = "%v"
		from_share = %v
	}`, name, strconv.Quote(externalShareId))
}

func TestAcc_SharedDatabase_IdentifierQuotingDiffSuppression(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())

	externalShareId := createShareableDatabase(t)
	unquotedExternalShareId := fmt.Sprintf("%s.%s.%s", externalShareId.AccountIdentifier().OrganizationName(), externalShareId.AccountIdentifier().AccountName(), externalShareId.Name())

	acc.TestClient().Database.CreateDatabaseFromShareTemporarily(t, externalShareId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SharedDatabase),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				ExpectNonEmptyPlan: true,
				Config:             sharedDatabaseConfigBasic(quotedId, unquotedExternalShareId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   sharedDatabaseConfigBasic(quotedId, unquotedExternalShareId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_shared_database.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_shared_database.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "id", id.Name()),
				),
			},
		},
	})
}
