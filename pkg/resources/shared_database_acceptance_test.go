//go:build account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CreateSharedDatabase_Basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	shareExternalId := createShareableDatabase(t)

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

	sharedDatabaseModel := model.SharedDatabase("test", shareExternalId.FullyQualifiedName(), id.Name()).
		WithComment(comment)
	sharedDatabaseModelRenamed := model.SharedDatabase("test", shareExternalId.FullyQualifiedName(), newId.Name()).
		WithComment(newComment)

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
				Config: accconfig.FromModels(t, sharedDatabaseModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(sharedDatabaseModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(sharedDatabaseModel.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(sharedDatabaseModel.ResourceReference(), "from_share", shareExternalId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(sharedDatabaseModel.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModel.ResourceReference(), "enable_console_output", accountEnableConsoleOutput),
				),
			},
			{
				Config: accconfig.FromModels(t, sharedDatabaseModelRenamed),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(sharedDatabaseModelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(sharedDatabaseModelRenamed.ResourceReference(), "name", newId.Name()),
					resource.TestCheckResourceAttr(sharedDatabaseModelRenamed.ResourceReference(), "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(sharedDatabaseModelRenamed.ResourceReference(), "from_share", shareExternalId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(sharedDatabaseModelRenamed.ResourceReference(), "comment", newComment),

					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "external_volume", accountExternalVolume),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "catalog", accountCatalog),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "replace_invalid_characters", accountReplaceInvalidCharacters),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "default_ddl_collation", accountDefaultDdlCollation),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "storage_serialization_policy", accountStorageSerializationPolicy),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "log_level", accountLogLevel),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "trace_level", accountTraceLevel),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "suspend_task_after_num_failures", accountSuspendTaskAfterNumFailures),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "task_auto_retry_attempts", accountTaskAutoRetryAttempts),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "user_task_managed_initial_warehouse_size", accountUserTaskMangedInitialWarehouseSize),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "user_task_timeout_ms", accountUserTaskTimeoutMs),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", accountUserTaskMinimumTriggerIntervalInSeconds),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "quoted_identifiers_ignore_case", accountQuotedIdentifiersIgnoreCase),
					resource.TestCheckResourceAttrPtr(sharedDatabaseModelRenamed.ResourceReference(), "enable_console_output", accountEnableConsoleOutput),
				),
			},
			// Import all values
			{
				Config:            accconfig.FromModels(t, sharedDatabaseModelRenamed),
				ResourceName:      sharedDatabaseModelRenamed.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CreateSharedDatabase_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	externalShareId := createShareableDatabase(t)

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	sharedDatabaseModelComplete := model.SharedDatabase("test", externalShareId.FullyQualifiedName(), id.Name()).
		WithComment(comment).
		WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithReplaceInvalidCharacters(true).
		WithDefaultDdlCollation("en_US").
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyOptimized)).
		WithLogLevel(string(sdk.LogLevelInfo)).
		WithTraceLevel(string(sdk.TraceLevelOnEvent)).
		WithSuspendTaskAfterNumFailures(20).
		WithTaskAutoRetryAttempts(20).
		WithUserTaskManagedInitialWarehouseSize(string(sdk.WarehouseSizeXLarge)).
		WithUserTaskTimeoutMs(1200000).
		WithUserTaskMinimumTriggerIntervalInSeconds(120).
		WithQuotedIdentifiersIgnoreCase(true).
		WithEnableConsoleOutput(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.SharedDatabase),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, sharedDatabaseModelComplete),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "from_share", externalShareId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "storage_serialization_policy", string(sdk.StorageSerializationPolicyOptimized)),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "suspend_task_after_num_failures", "20"),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "task_auto_retry_attempts", "20"),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "user_task_managed_initial_warehouse_size", string(sdk.WarehouseSizeXLarge)),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "user_task_timeout_ms", "1200000"),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "user_task_minimum_trigger_interval_in_seconds", "120"),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "quoted_identifiers_ignore_case", "true"),
					resource.TestCheckResourceAttr(sharedDatabaseModelComplete.ResourceReference(), "enable_console_output", "true"),
				),
			},
			// Import all values
			{
				Config:            accconfig.FromModels(t, sharedDatabaseModelComplete),
				ResourceName:      sharedDatabaseModelComplete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CreateSharedDatabase_InvalidValues(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	comment := random.Comment()

	sharedDatabaseModelInvalid := model.SharedDatabase("test", "org.acc.name", "org.acc.name").
		WithComment(comment).
		WithStorageSerializationPolicy("invalid_value").
		WithLogLevel("invalid_value").
		WithTraceLevel("invalid_value").
		WithUserTaskManagedInitialWarehouseSize("invalid_value")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.SharedDatabase),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, sharedDatabaseModelInvalid),
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

	sharedDatabase, sharedDatabaseCleanup := acc.SecondaryTestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(sharedDatabaseCleanup)

	revoke := acc.SecondaryTestClient().Grant.GrantPrivilegeOnDatabaseToShare(t, sharedDatabase.ID(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage})
	t.Cleanup(revoke)

	acc.SecondaryTestClient().Share.SetAccountOnShare(t, acc.TestClient().Account.GetAccountIdentifier(t), share.ID())

	externalShareId := sdk.NewExternalObjectIdentifier(acc.SecondaryTestClient().Account.GetAccountIdentifier(t), share.ID())

	acc.TestClient().Database.CreateDatabaseFromShareTemporarily(t, externalShareId)

	return externalShareId
}

func TestAcc_SharedDatabase_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	externalShareId := createShareableDatabase(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	sharedDatabaseModel := model.SharedDatabase("test", externalShareId.FullyQualifiedName(), id.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SharedDatabase),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, sharedDatabaseModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(sharedDatabaseModel.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, sharedDatabaseModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(sharedDatabaseModel.ResourceReference(), "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_SharedDatabase_IdentifierQuotingDiffSuppression(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	externalShareId := createShareableDatabase(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`"%s"`, id.Name())
	unquotedExternalShareId := fmt.Sprintf("%s.%s.%s", externalShareId.AccountIdentifier().OrganizationName(), externalShareId.AccountIdentifier().AccountName(), externalShareId.Name())

	sharedDatabaseModel := model.SharedDatabase("test", unquotedExternalShareId, quotedId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SharedDatabase),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  acc.ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             accconfig.FromModels(t, sharedDatabaseModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(sharedDatabaseModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(sharedDatabaseModel.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, sharedDatabaseModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(sharedDatabaseModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(sharedDatabaseModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(sharedDatabaseModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(sharedDatabaseModel.ResourceReference(), "id", id.Name()),
				),
			},
		},
	})
}
