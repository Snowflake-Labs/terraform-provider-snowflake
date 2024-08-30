package resources_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_User_BasicFlows(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	id2 := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()
	newComment := random.Comment()

	key1, _ := random.GenerateRSAPublicKey(t)
	key2, _ := random.GenerateRSAPublicKey(t)

	pass := random.Password()
	newPass := random.Password()

	userModelNoAttributes := model.User("w", id.Name())
	userModelNoAttributesRenamed := model.User("w", id2.Name()).
		WithComment(newComment)

	userModelAllAttributes := model.User("w", id.Name()).
		WithPassword(pass).
		WithLoginName(id.Name() + "_login").
		WithDisplayName("Display Name").
		WithFirstName("Jan").
		WithMiddleName("Jakub").
		WithLastName("Testowski").
		WithEmail("fake@email.com").
		WithMustChangePassword("true").
		WithDisabled("false").
		WithDaysToExpiry(8).
		WithMinsToUnlock(9).
		WithDefaultWarehouse("some_warehouse").
		WithDefaultNamespace("some.namespace").
		WithDefaultRole("some_role").
		WithDefaultSecondaryRolesStringList("ALL").
		WithMinsToBypassMfa(10).
		WithRsaPublicKey(key1).
		WithRsaPublicKey2(key2).
		WithComment(comment).
		WithDisableMfa("true")

	userModelAllAttributesChanged := model.User("w", id.Name()).
		WithPassword(newPass).
		WithLoginName(id.Name() + "_other_login").
		WithDisplayName("New Display Name").
		WithFirstName("Janek").
		WithMiddleName("Kuba").
		WithLastName("Terraformowski").
		WithEmail("fake@email.net").
		WithMustChangePassword("false").
		WithDisabled("true").
		WithDaysToExpiry(12).
		WithMinsToUnlock(13).
		WithDefaultWarehouse("other_warehouse").
		WithDefaultNamespace("one_part_namespace").
		WithDefaultRole("other_role").
		WithDefaultSecondaryRolesStringList("ALL").
		WithMinsToBypassMfa(14).
		WithRsaPublicKey(key2).
		WithRsaPublicKey2(key1).
		WithComment(newComment).
		WithDisableMfa("false")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// CREATE WITHOUT ATTRIBUTES
			{
				Config: config.FromModel(t, userModelNoAttributes),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasNoPassword().
						HasNoLoginName().
						HasNoDisplayName().
						HasNoFirstName().
						HasNoMiddleName().
						HasNoLastName().
						HasNoEmail().
						HasMustChangePasswordString(r.BooleanDefault).
						HasDisabledString(r.BooleanDefault).
						HasNoDaysToExpiry().
						HasMinsToUnlockString(r.IntDefaultString).
						HasNoDefaultWarehouse().
						HasNoDefaultNamespace().
						HasNoDefaultRole().
						HasNoDefaultSecondaryRoles().
						HasMinsToBypassMfaString(r.IntDefaultString).
						HasNoRsaPublicKey().
						HasNoRsaPublicKey2().
						HasNoComment().
						HasDisableMfaString(r.BooleanDefault).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
						HasLoginName(fmt.Sprintf(id.Name())).
						HasDisplayName(id.Name()),
				),
			},
			// RENAME AND CHANGE ONE PROP
			{
				Config: config.FromModel(t, userModelNoAttributesRenamed),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelNoAttributes.ResourceReference()).
						HasNameString(id2.Name()).
						HasCommentString(newComment),
					// default names stay the same
					resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
						HasLoginName(fmt.Sprintf(id.Name())).
						HasDisplayName(id.Name()),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelNoAttributesRenamed.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "disable_mfa", "days_to_expiry", "mins_to_unlock", "mins_to_bypass_mfa", "login_name", "display_name", "disabled", "must_change_password"},
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedUserResource(t, id2.Name()).
						HasLoginNameString(fmt.Sprintf(id.Name())).
						HasDisplayNameString(fmt.Sprintf(id.Name())).
						HasDisabled(false).
						HasMustChangePassword(false),
				),
			},
			// DESTROY
			{
				Config:  config.FromModel(t, userModelNoAttributes),
				Destroy: true,
			},
			// CREATE WITH ALL ATTRIBUTES
			{
				Config: config.FromModel(t, userModelAllAttributes),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelAllAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasPasswordString(pass).
						HasLoginNameString(fmt.Sprintf("%s_login", id.Name())).
						HasDisplayNameString("Display Name").
						HasFirstNameString("Jan").
						HasMiddleNameString("Jakub").
						HasLastNameString("Testowski").
						HasEmailString("fake@email.com").
						HasMustChangePassword(true).
						HasDisabled(false).
						HasDaysToExpiryString("8").
						HasMinsToUnlockString("9").
						HasDefaultWarehouseString("some_warehouse").
						HasDefaultNamespaceString("some.namespace").
						HasDefaultRoleString("some_role").
						HasDefaultSecondaryRoles("ALL").
						HasMinsToBypassMfaString("10").
						HasRsaPublicKeyString(key1).
						HasRsaPublicKey2String(key2).
						HasCommentString(comment).
						HasDisableMfaString(r.BooleanTrue).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: config.FromModel(t, userModelAllAttributesChanged),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelAllAttributesChanged.ResourceReference()).
						HasNameString(id.Name()).
						HasPasswordString(newPass).
						HasLoginNameString(fmt.Sprintf("%s_other_login", id.Name())).
						HasDisplayNameString("New Display Name").
						HasFirstNameString("Janek").
						HasMiddleNameString("Kuba").
						HasLastNameString("Terraformowski").
						HasEmailString("fake@email.net").
						HasMustChangePassword(false).
						HasDisabled(true).
						HasDaysToExpiryString("12").
						HasMinsToUnlockString("13").
						HasDefaultWarehouseString("other_warehouse").
						HasDefaultNamespaceString("one_part_namespace").
						HasDefaultRoleString("other_role").
						HasDefaultSecondaryRoles("ALL").
						HasMinsToBypassMfaString("14").
						HasRsaPublicKeyString(key2).
						HasRsaPublicKey2String(key1).
						HasCommentString(newComment).
						HasDisableMfaString(r.BooleanFalse).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelAllAttributesChanged.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "disable_mfa", "days_to_expiry", "mins_to_unlock", "mins_to_bypass_mfa", "default_namespace", "login_name", "show_output.0.days_to_expiry"},
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedUserResource(t, id.Name()).
						HasDefaultNamespaceString("ONE_PART_NAMESPACE").
						HasLoginNameString(fmt.Sprintf("%s_OTHER_LOGIN", id.Name())),
				),
			},
			// UNSET ALL
			{
				Config: config.FromModel(t, userModelNoAttributes),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasPasswordString("").
						HasLoginNameString("").
						HasDisplayNameString("").
						HasFirstNameString("").
						HasMiddleNameString("").
						HasLastNameString("").
						HasEmailString("").
						HasMustChangePasswordString(r.BooleanDefault).
						HasDisabledString(r.BooleanDefault).
						HasDaysToExpiryString("0").
						HasMinsToUnlockString(r.IntDefaultString).
						HasDefaultWarehouseString("").
						HasDefaultNamespaceString("").
						HasDefaultRoleString("").
						HasDefaultSecondaryRolesEmpty().
						HasMinsToBypassMfaString(r.IntDefaultString).
						HasRsaPublicKeyString("").
						HasRsaPublicKey2String("").
						HasCommentString("").
						HasDisableMfaString(r.BooleanDefault).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
						HasLoginName(fmt.Sprintf(id.Name())).
						HasDisplayName(""),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2481 has been fixed
func TestAcc_User_RemovedOutsideOfTerraform(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("u", userId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				PreConfig: acc.TestClient().User.DropUserFunc(t, userId),
				Config:    config.FromModel(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					func(state *terraform.State) error {
						if len(state.RootModule().Resources) != 1 {
							return errors.New("user should be created again and present in the state")
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAcc_User_issue2058 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2058 issue.
// The problem was with a dot in user identifier.
// Before the fix it results in panic: interface conversion: sdk.ObjectIdentifier is sdk.DatabaseObjectIdentifier, not sdk.AccountObjectIdentifier error.
func TestAcc_User_issue2058(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifierContaining(".")

	userModel1 := model.User("w", userId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, userModel1),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel1.ResourceReference()).HasNameString(userId.Name()),
				),
			},
		},
	})
}

func TestAcc_User_AllParameters(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)

	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("u", userId.Name())
	userModelWithAllParametersSet := model.User("u", userId.Name()).
		WithAbortDetachedQuery(true).
		WithAutocommit(false).
		WithBinaryInputFormatEnum(sdk.BinaryInputFormatUTF8).
		WithBinaryOutputFormatEnum(sdk.BinaryOutputFormatBase64).
		WithClientMemoryLimit(1024).
		WithClientMetadataRequestUseConnectionCtx(true).
		WithClientPrefetchThreads(2).
		WithClientResultChunkSize(48).
		WithClientResultColumnCaseInsensitive(true).
		WithClientSessionKeepAlive(true).
		WithClientSessionKeepAliveHeartbeatFrequency(2400).
		WithClientTimestampTypeMappingEnum(sdk.ClientTimestampTypeMappingNtz).
		WithDateInputFormat("YYYY-MM-DD").
		WithDateOutputFormat("YY-MM-DD").
		WithEnableUnloadPhysicalTypeOptimization(false).
		WithErrorOnNondeterministicMerge(false).
		WithErrorOnNondeterministicUpdate(true).
		WithGeographyOutputFormatEnum(sdk.GeographyOutputFormatWKB).
		WithGeometryOutputFormatEnum(sdk.GeometryOutputFormatWKB).
		WithJdbcTreatDecimalAsInt(false).
		WithJdbcTreatTimestampNtzAsUtc(true).
		WithJdbcUseSessionTimezone(false).
		WithJsonIndent(4).
		WithLockTimeout(21222).
		WithLogLevelEnum(sdk.LogLevelError).
		WithMultiStatementCount(0).
		WithNoorderSequenceAsDefault(false).
		WithOdbcTreatDecimalAsInt(true).
		WithQueryTag("some_tag").
		WithQuotedIdentifiersIgnoreCase(true).
		WithRowsPerResultset(2).
		WithS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
		WithSearchPath("$public, $current").
		WithSimulatedDataSharingConsumer("some_consumer").
		WithStatementQueuedTimeoutInSeconds(10).
		WithStatementTimeoutInSeconds(10).
		WithStrictJsonOutput(true).
		WithTimestampDayIsAlways24h(true).
		WithTimestampInputFormat("YYYY-MM-DD").
		WithTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampTypeMappingEnum(sdk.TimestampTypeMappingLtz).
		WithTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimezone("Europe/Warsaw").
		WithTimeInputFormat("HH24:MI").
		WithTimeOutputFormat("HH24:MI").
		WithTraceLevelEnum(sdk.TraceLevelOnEvent).
		WithTransactionAbortOnError(true).
		WithTransactionDefaultIsolationLevelEnum(sdk.TransactionDefaultIsolationLevelReadCommitted).
		WithTwoDigitCenturyStart(1980).
		WithUnsupportedDdlActionEnum(sdk.UnsupportedDDLActionFail).
		WithUseCachedResult(false).
		WithWeekOfYearPolicy(1).
		WithWeekStart(1).
		WithEnableUnredactedQuerySyntaxError(true).
		WithNetworkPolicyId(networkPolicy.ID()).
		WithPreventUnloadToInternalStages(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// create with default values for all the parameters
			{
				Config: config.FromModel(t, userModel),
				Check: assert.AssertThat(t,
					objectparametersassert.UserParameters(t, userId).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.UserResourceParameters(t, userModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			// import when no parameter set
			{
				ResourceName: userModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceparametersassert.ImportedUserResourceParameters(t, userId.Name()).
						HasAllDefaults(),
				),
			},
			// set all parameters
			{
				Config: config.FromModel(t, userModelWithAllParametersSet),
				Check: assert.AssertThat(t,
					objectparametersassert.UserParameters(t, userId).
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcTreatDecimalAsInt(false).
						HasJdbcTreatTimestampNtzAsUtc(true).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasSimulatedDataSharingConsumer("some_consumer").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelOnEvent).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1).
						HasEnableUnredactedQuerySyntaxError(true).
						HasNetworkPolicy(networkPolicy.ID().Name()).
						HasPreventUnloadToInternalStages(true),
					resourceparametersassert.UserResourceParameters(t, "snowflake_user.u").
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcTreatDecimalAsInt(false).
						HasJdbcTreatTimestampNtzAsUtc(true).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasSimulatedDataSharingConsumer("some_consumer").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelOnEvent).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1).
						HasEnableUnredactedQuerySyntaxError(true).
						HasNetworkPolicy(networkPolicy.ID().Name()).
						HasPreventUnloadToInternalStages(true),
				),
			},
			// import when all parameters set
			{
				ResourceName: userModelWithAllParametersSet.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceparametersassert.ImportedUserResourceParameters(t, userId.Name()).
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcTreatDecimalAsInt(false).
						HasJdbcTreatTimestampNtzAsUtc(true).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasSimulatedDataSharingConsumer("some_consumer").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelOnEvent).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1).
						HasEnableUnredactedQuerySyntaxError(true).
						HasNetworkPolicy(networkPolicy.ID().Name()).
						HasPreventUnloadToInternalStages(true),
				),
			},
			// unset all the parameters
			{
				Config: config.FromModel(t, userModel),
				Check: assert.AssertThat(t,
					objectparametersassert.UserParameters(t, userId).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.UserResourceParameters(t, userModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
		},
	})
}

func TestAcc_User_issue2836(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	defaultRole := "SOME ROLE WITH SPACE case sensitive"
	defaultRoleQuoted := fmt.Sprintf(`"%s"`, defaultRole)

	userModel := model.User("u", userId.Name()).
		WithDefaultRole(defaultRoleQuoted)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, userModel),
				Check: assert.AssertThat(t,
					objectassert.User(t, userId).
						HasDefaultRole(defaultRole),
				),
			},
		},
	})
}

func TestAcc_User_issue2970(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	pass := random.Password()
	key, _ := random.GenerateRSAPublicKey(t)
	resourceName := "u"

	newPass := random.Password()
	newKey, _ := random.GenerateRSAPublicKey(t)
	incorrectlyFormattedNewKey := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s-----END PUBLIC KEY-----\n", newKey)

	userModel := model.User(resourceName, userId.Name()).
		WithPassword(pass).
		WithRsaPublicKey(key)

	newUserModelIncorrectNewKey := model.User(resourceName, userId.Name()).
		WithPassword(newPass).
		WithRsaPublicKey(incorrectlyFormattedNewKey)

	newUserModel := model.User(resourceName, userId.Name()).
		WithPassword(newPass).
		WithRsaPublicKey(newKey)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, userModel),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasPasswordString(pass).
						HasRsaPublicKeyString(key),
				),
			},
			{
				Config: config.FromModel(t, newUserModelIncorrectNewKey),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(newUserModelIncorrectNewKey.ResourceReference(), "password", tfjson.ActionUpdate, sdk.String(pass), sdk.String(newPass)),
						planchecks.ExpectChange(newUserModelIncorrectNewKey.ResourceReference(), "rsa_public_key", tfjson.ActionUpdate, sdk.String(key), sdk.String(incorrectlyFormattedNewKey)),
					},
				},
				ExpectError: regexp.MustCompile("New public key rejected by current policy"),
			},
			{
				Config: config.FromModel(t, newUserModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(newUserModel.ResourceReference(), "password", tfjson.ActionUpdate, sdk.String(pass), sdk.String(newPass)),
						planchecks.ExpectChange(newUserModel.ResourceReference(), "rsa_public_key", tfjson.ActionUpdate, sdk.String(key), sdk.String(newKey)),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, newUserModel.ResourceReference()).
						HasPasswordString(newPass).
						HasRsaPublicKeyString(newKey),
				),
			},
		},
	})
}

func TestAcc_User_issue1572(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.UserWithDefaultMeta(userId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, userModel),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasDisabledString(r.BooleanDefault),
					objectassert.User(t, userId).HasDisabled(false),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().User.Disable(t, userId)
					objectassert.User(t, userId).HasDisabled(true)
				},
				Config: config.FromModel(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectDrift(userModel.ResourceReference(), "disabled", sdk.String(r.BooleanDefault), sdk.String(r.BooleanTrue)),
						planchecks.ExpectChange(userModel.ResourceReference(), "disabled", tfjson.ActionUpdate, sdk.String(r.BooleanTrue), sdk.String(r.BooleanDefault)),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasDisabledString(r.BooleanDefault),
					objectassert.User(t, userId).HasDisabled(false),
				),
			},
		},
	})
}

func TestAcc_User_issue1535_withNullPassword(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	pass := random.Password()

	userModel := model.UserWithDefaultMeta(userId.Name()).
		WithPassword(pass)

	userWithNullPasswordModel := model.UserWithDefaultMeta(userId.Name()).
		WithNullPassword()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, userModel),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasPasswordString(pass),
				),
			},
			{
				Config: config.FromModel(t, userWithNullPasswordModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userWithNullPasswordModel.ResourceReference(), "password", tfjson.ActionUpdate, sdk.String(pass), nil),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userWithNullPasswordModel.ResourceReference()).
						HasEmptyPassword(),
				),
			},
		},
	})
}

func TestAcc_User_issue1535_withRemovedPassword(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	pass := random.Password()

	userModel := model.UserWithDefaultMeta(userId.Name()).
		WithPassword(pass)

	userWithoutPasswordModel := model.UserWithDefaultMeta(userId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, userModel),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasPasswordString(pass),
				),
			},
			{
				Config: config.FromModel(t, userWithoutPasswordModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userWithoutPasswordModel.ResourceReference(), "password", tfjson.ActionUpdate, sdk.String(pass), nil),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userWithoutPasswordModel.ResourceReference()).
						HasEmptyPassword(),
				),
			},
		},
	})
}

func TestAcc_User_issue1155_handleChangesToDaysToExpiry(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	userModelWithoutDaysToExpiry := model.UserWithDefaultMeta(userId.Name())
	userModelDaysToExpiry10 := model.UserWithDefaultMeta(userId.Name()).WithDaysToExpiry(10)
	userModelDaysToExpiry5 := model.UserWithDefaultMeta(userId.Name()).WithDaysToExpiry(5)
	userModelDaysToExpiry0 := model.UserWithDefaultMeta(userId.Name()).WithDaysToExpiry(0)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// 1. create without days_to_expiry
			{
				Config: config.FromModel(t, userModelWithoutDaysToExpiry),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelWithoutDaysToExpiry.ResourceReference()).HasNoDaysToExpiry(),
					objectassert.User(t, userId).HasDaysToExpiryEmpty(),
				),
			},
			// 2. change to 10 (no plan after)
			{
				Config: config.FromModel(t, userModelDaysToExpiry10),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelDaysToExpiry10.ResourceReference()).HasDaysToExpiryString("10"),
					objectassert.User(t, userId).HasDaysToExpiryNotEmpty(),
				),
			},
			// 3. change externally to 2 (no changes)
			{
				PreConfig: func() {
					acc.TestClient().User.SetDaysToExpiry(t, userId, 2)
				},
				Config: config.FromModel(t, userModelDaysToExpiry10),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// 4. change externally to 0 (no changes)
			{
				PreConfig: func() {
					acc.TestClient().User.SetDaysToExpiry(t, userId, 0)
				},
				Config: config.FromModel(t, userModelDaysToExpiry10),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// 5. change in config to 5 (change)
			{
				Config: config.FromModel(t, userModelDaysToExpiry5),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelDaysToExpiry5.ResourceReference(), "days_to_expiry", tfjson.ActionUpdate, sdk.String("10"), sdk.String("5")),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelDaysToExpiry10.ResourceReference()).HasDaysToExpiryString("5"),
					objectassert.User(t, userId).HasDaysToExpiryNotEmpty(),
				),
			},
			// 6. change in config to 0 (change)
			{
				Config: config.FromModel(t, userModelDaysToExpiry0),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelDaysToExpiry0.ResourceReference(), "days_to_expiry", tfjson.ActionUpdate, sdk.String("5"), sdk.String("0")),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelDaysToExpiry10.ResourceReference()).HasDaysToExpiryString("0"),
					objectassert.User(t, userId).HasDaysToExpiryEmpty(),
				),
			},
			// 7. remove from config (no change)
			{
				Config: config.FromModel(t, userModelWithoutDaysToExpiry),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelWithoutDaysToExpiry.ResourceReference()).HasDaysToExpiryString("0"),
					objectassert.User(t, userId).HasDaysToExpiryEmpty(),
				),
			},
		},
	})
}

func TestAcc_User_handleExternalTypeChange(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.UserWithDefaultMeta(userId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, userModel),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString(""),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType(""),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().User.SetType(t, userId, "SERVICE")
					objectassert.User(t, userId).HasType("SERVICE")
				},
				Config: config.FromModel(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString(""),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType(""),
				),
			},
		},
	})
}

func TestAcc_User_handleChangesToDefaultSecondaryRoles(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	userModelEmpty := model.UserWithDefaultMeta(userId.Name())
	userModelWithDefaultSecondaryRole := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesStringList("ALL")
	userModelLowercaseValue := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesStringList("all")
	userModelIncorrectValue := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesStringList("OTHER")
	userModelNoValues := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesStringList()
	userModelMultipleValues := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesStringList("ALL", "OTHER")
	userModelRepeatedDifferentCasing := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesStringList("ALL", "all")
	userModelRepeatedValues := model.UserWithDefaultMeta(userId.Name()).WithDefaultSecondaryRolesStringList("ALL", "ALL")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// 1. create without default secondary roles
			{
				Config: config.FromModel(t, userModelEmpty),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelEmpty.ResourceReference()).HasNoDefaultSecondaryRoles(),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(""),
				),
			},
			// 2. add default secondary roles
			{
				Config: config.FromModel(t, userModelWithDefaultSecondaryRole),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelWithDefaultSecondaryRole.ResourceReference()).HasDefaultSecondaryRoles("ALL"),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(`["ALL"]`),
				),
			},
			// 3. change to lowercase (no changes)
			{
				Config: config.FromModel(t, userModelLowercaseValue),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// 4. unset externally
			{
				PreConfig: func() {
					acc.TestClient().User.UnsetDefaultSecondaryRoles(t, userId)
				},
				Config: config.FromModel(t, userModelWithDefaultSecondaryRole),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelWithDefaultSecondaryRole.ResourceReference(), "default_secondary_roles", tfjson.ActionUpdate, sdk.String("[]"), sdk.String("[ALL]")),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelWithDefaultSecondaryRole.ResourceReference()).HasDefaultSecondaryRoles("ALL"),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(`["ALL"]`),
				),
			},
			// 5. unset in config to 5 (change)
			{
				Config: config.FromModel(t, userModelEmpty),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(userModelEmpty.ResourceReference(), "default_secondary_roles", tfjson.ActionUpdate, sdk.String("[ALL]"), sdk.String("[]")),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelEmpty.ResourceReference()).HasDefaultSecondaryRolesEmpty(),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(""),
				),
			},
			// 6. incorrect value used
			{
				Config:      config.FromModel(t, userModelIncorrectValue),
				ExpectError: regexp.MustCompile("Unsupported secondary role 'OTHER'"),
			},
			// 7. empty set used
			{
				Config:      config.FromModel(t, userModelNoValues),
				ExpectError: regexp.MustCompile("Attribute default_secondary_roles requires 1 item minimum"),
			},
			// 8. multiple values (correct and incorrect)
			{
				Config:      config.FromModel(t, userModelMultipleValues),
				ExpectError: regexp.MustCompile("Attribute default_secondary_roles supports 1 item maximum"),
			},
			// 9. multiple values (different casing)
			{
				Config:      config.FromModel(t, userModelRepeatedDifferentCasing),
				ExpectError: regexp.MustCompile("Attribute default_secondary_roles supports 1 item maximum"),
			},
			// 10. multiple values (two same) - no error
			{
				Config: config.FromModel(t, userModelRepeatedValues),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelWithDefaultSecondaryRole.ResourceReference()).HasDefaultSecondaryRoles("ALL"),
					objectassert.User(t, userId).HasDefaultSecondaryRoles(`["ALL"]`),
				),
			},
		},
	})
}
