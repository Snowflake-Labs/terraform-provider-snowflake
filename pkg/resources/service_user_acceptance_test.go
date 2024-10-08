package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ServiceUser_BasicFlows(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	id2 := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()
	newComment := random.Comment()

	key1, _ := random.GenerateRSAPublicKey(t)
	key2, _ := random.GenerateRSAPublicKey(t)

	userModelNoAttributes := model.ServiceUser("w", id.Name())
	userModelNoAttributesRenamed := model.ServiceUser("w", id2.Name()).
		WithComment(newComment)

	userModelAllAttributes := model.ServiceUser("w", id.Name()).
		WithLoginName(id.Name() + "_login").
		WithDisplayName("Display Name").
		WithEmail("fake@email.com").
		WithDisabled("false").
		WithDaysToExpiry(8).
		WithMinsToUnlock(9).
		WithDefaultWarehouse("some_warehouse").
		WithDefaultNamespace("some.namespace").
		WithDefaultRole("some_role").
		WithDefaultSecondaryRolesOptionEnum(sdk.SecondaryRolesOptionAll).
		WithRsaPublicKey(key1).
		WithRsaPublicKey2(key2).
		WithComment(comment)

	userModelAllAttributesChanged := func(loginName string) *model.ServiceUserModel {
		return model.ServiceUser("w", id.Name()).
			WithLoginName(loginName).
			WithDisplayName("New Display Name").
			WithEmail("fake@email.net").
			WithDisabled("true").
			WithDaysToExpiry(12).
			WithMinsToUnlock(13).
			WithDefaultWarehouse("other_warehouse").
			WithDefaultNamespace("one_part_namespace").
			WithDefaultRole("other_role").
			WithDefaultSecondaryRolesOptionEnum(sdk.SecondaryRolesOptionAll).
			WithRsaPublicKey(key2).
			WithRsaPublicKey2(key1).
			WithComment(newComment)
	}

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
					resourceassert.ServiceUserResource(t, userModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasNoLoginName().
						HasNoDisplayName().
						HasNoEmail().
						HasDisabledString(r.BooleanDefault).
						HasNoDaysToExpiry().
						HasMinsToUnlockString(r.IntDefaultString).
						HasNoDefaultWarehouse().
						HasNoDefaultNamespace().
						HasNoDefaultRole().
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault).
						HasNoRsaPublicKey().
						HasNoRsaPublicKey2().
						HasNoComment().
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
						HasLoginName(strings.ToUpper(id.Name())).
						HasDisplayName(id.Name()),
				),
			},
			// RENAME AND CHANGE ONE PROP
			{
				Config: config.FromModel(t, userModelNoAttributesRenamed),
				Check: assert.AssertThat(t,
					resourceassert.ServiceUserResource(t, userModelNoAttributes.ResourceReference()).
						HasNameString(id2.Name()).
						HasCommentString(newComment),
					// default names stay the same
					resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
						HasLoginName(strings.ToUpper(id.Name())).
						HasDisplayName(id.Name()),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelNoAttributesRenamed.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_unlock", "mins_to_bypass_mfa", "login_name", "display_name", "disabled"},
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedServiceUserResource(t, id2.Name()).
						HasLoginNameString(strings.ToUpper(id.Name())).
						HasDisplayNameString(id.Name()).
						HasDisabled(false),
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
					resourceassert.ServiceUserResource(t, userModelAllAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasLoginNameString(fmt.Sprintf("%s_login", id.Name())).
						HasDisplayNameString("Display Name").
						HasEmailString("fake@email.com").
						HasDisabled(false).
						HasDaysToExpiryString("8").
						HasMinsToUnlockString("9").
						HasDefaultWarehouseString("some_warehouse").
						HasDefaultNamespaceString("some.namespace").
						HasDefaultRoleString("some_role").
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionAll).
						HasRsaPublicKeyString(key1).
						HasRsaPublicKey2String(key2).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: config.FromModel(t, userModelAllAttributesChanged(id.Name()+"_other_login")),
				Check: assert.AssertThat(t,
					resourceassert.ServiceUserResource(t, userModelAllAttributesChanged(id.Name()+"_other_login").ResourceReference()).
						HasNameString(id.Name()).
						HasLoginNameString(fmt.Sprintf("%s_other_login", id.Name())).
						HasDisplayNameString("New Display Name").
						HasEmailString("fake@email.net").
						HasDisabled(true).
						HasDaysToExpiryString("12").
						HasMinsToUnlockString("13").
						HasDefaultWarehouseString("other_warehouse").
						HasDefaultNamespaceString("one_part_namespace").
						HasDefaultRoleString("other_role").
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionAll).
						HasRsaPublicKeyString(key2).
						HasRsaPublicKey2String(key1).
						HasCommentString(newComment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelAllAttributesChanged(id.Name() + "_other_login").ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_unlock", "default_namespace", "login_name", "show_output.0.days_to_expiry"},
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedServiceUserResource(t, id.Name()).
						HasDefaultNamespaceString("ONE_PART_NAMESPACE").
						HasLoginNameString(fmt.Sprintf("%s_OTHER_LOGIN", id.Name())),
				),
			},
			// CHANGE PROP TO THE CURRENT SNOWFLAKE VALUE
			{
				PreConfig: func() {
					acc.TestClient().User.SetLoginName(t, id, id.Name()+"_different_login")
				},
				Config: config.FromModel(t, userModelAllAttributesChanged(id.Name()+"_different_login")),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// UNSET ALL
			{
				Config: config.FromModel(t, userModelNoAttributes),
				Check: assert.AssertThat(t,
					resourceassert.ServiceUserResource(t, userModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasLoginNameString("").
						HasDisplayNameString("").
						HasEmailString("").
						HasDisabledString(r.BooleanDefault).
						HasDaysToExpiryString("0").
						HasMinsToUnlockString(r.IntDefaultString).
						HasDefaultWarehouseString("").
						HasDefaultNamespaceString("").
						HasDefaultRoleString("").
						HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault).
						HasRsaPublicKeyString("").
						HasRsaPublicKey2String("").
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
						HasLoginName(strings.ToUpper(id.Name())).
						HasDisplayName(""),
				),
			},
		},
	})
}

func TestAcc_ServiceUser_AllParameters(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)

	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.ServiceUser("u", userId.Name())
	userModelWithAllParametersSet := model.ServiceUser("u", userId.Name()).
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
					resourceparametersassert.UserResourceParameters(t, userModelWithAllParametersSet.ResourceReference()).
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

func TestAcc_ServiceUser_handleExternalTypeChange(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.ServiceUserWithDefaultMeta(userId.Name())

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
					resourceassert.ServiceUserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString("SERVICE"),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType("SERVICE"),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().User.SetType(t, userId, sdk.UserTypePerson)
					objectassert.User(t, userId).HasType(string(sdk.UserTypePerson))
				},
				Config: config.FromModel(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString("SERVICE"),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType("SERVICE"),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().User.UnsetType(t, userId)
					objectassert.User(t, userId).HasType("")
				},
				Config: config.FromModel(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).HasNameString(userId.Name()).HasUserTypeString("SERVICE"),
					resourceshowoutputassert.UserShowOutput(t, userModel.ResourceReference()).HasType("SERVICE"),
				),
			},
		},
	})
}
