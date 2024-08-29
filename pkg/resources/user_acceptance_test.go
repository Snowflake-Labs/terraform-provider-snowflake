package resources_test

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func checkBool(path, attr string, value bool) func(*terraform.State) error {
	return func(state *terraform.State) error {
		is := state.RootModule().Resources[path].Primary
		d := is.Attributes[attr]
		b, err := strconv.ParseBool(d)
		if err != nil {
			return err
		}
		if b != value {
			return fmt.Errorf("at %s expected %t but got %t", path, value, b)
		}
		return nil
	}
}

// TODO [SNOW-1348101]: handle 1-part default namespace
func TestAcc_User(t *testing.T) {
	r := require.New(t)
	prefix := acc.TestClient().Ids.Alpha()
	prefix2 := acc.TestClient().Ids.Alpha()
	id := sdk.NewAccountObjectIdentifier(prefix)
	id2 := sdk.NewAccountObjectIdentifier(prefix2)

	comment := random.Comment()
	newComment := random.Comment()

	sshkey1, err := testhelpers.Fixture("userkey1")
	r.NoError(err)

	sshkey2, err := testhelpers.Fixture("userkey2")
	r.NoError(err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: uConfig(prefix, sshkey1, sshkey2, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix),
					resource.TestCheckResourceAttr("snowflake_user.w", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_user.w", "login_name", strings.ToUpper(fmt.Sprintf("%s_login", prefix))),
					resource.TestCheckResourceAttr("snowflake_user.w", "display_name", "Display Name"),
					resource.TestCheckResourceAttr("snowflake_user.w", "first_name", "Marcin"),
					resource.TestCheckResourceAttr("snowflake_user.w", "last_name", "Zukowski"),
					resource.TestCheckResourceAttr("snowflake_user.w", "email", "fake@email.com"),
					checkBool("snowflake_user.w", "disabled", false),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_warehouse", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_role", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_secondary_roles.0", "ALL"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_namespace", "foo.bar"),
					resource.TestCheckResourceAttr("snowflake_user.w", "fully_qualified_name", id.FullyQualifiedName()),
					checkBool("snowflake_user.w", "must_change_password", true),
				),
			},
			// RENAME
			{
				Config: uConfig(prefix2, sshkey1, sshkey2, newComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_user.w", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_user.w", "comment", newComment),
					resource.TestCheckResourceAttr("snowflake_user.w", "login_name", strings.ToUpper(fmt.Sprintf("%s_login", prefix2))),
					resource.TestCheckResourceAttr("snowflake_user.w", "display_name", "Display Name"),
					resource.TestCheckResourceAttr("snowflake_user.w", "first_name", "Marcin"),
					resource.TestCheckResourceAttr("snowflake_user.w", "last_name", "Zukowski"),
					resource.TestCheckResourceAttr("snowflake_user.w", "email", "fake@email.com"),
					checkBool("snowflake_user.w", "disabled", false),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_warehouse", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_role", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_secondary_roles.0", "ALL"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_namespace", "foo.bar"),
					resource.TestCheckResourceAttr("snowflake_user.w", "fully_qualified_name", id2.FullyQualifiedName()),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: uConfig2(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_user.w", "comment", "test comment 2"),
					resource.TestCheckResourceAttr("snowflake_user.w", "password", "best password"),
					resource.TestCheckResourceAttr("snowflake_user.w", "login_name", strings.ToUpper(fmt.Sprintf("%s_login", prefix2))),
					resource.TestCheckResourceAttr("snowflake_user.w", "display_name", "New Name"),
					resource.TestCheckResourceAttr("snowflake_user.w", "first_name", "Benoit"),
					resource.TestCheckResourceAttr("snowflake_user.w", "last_name", "Dageville"),
					resource.TestCheckResourceAttr("snowflake_user.w", "email", "fake@email.net"),
					checkBool("snowflake_user.w", "disabled", true),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_warehouse", "bar"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_role", "bar"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_secondary_roles.#", "0"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_namespace", "bar.baz"),
					resource.TestCheckResourceAttr("snowflake_user.w", "fully_qualified_name", id2.FullyQualifiedName()),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_user.w",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rsa_public_key", "rsa_public_key_2", "must_change_password"},
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2481 has been fixed
func TestAcc_User_RemovedOutsideOfTerraform(t *testing.T) {
	userName := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	config := fmt.Sprintf(`
resource "snowflake_user" "test" {
	name = "%s"
}
`, userName.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				PreConfig: acc.TestClient().User.DropUserFunc(t, userName),
				Config:    config,
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

func uConfig(prefix, key1, key2, comment string) string {
	s := `
resource "snowflake_user" "w" {
	name = "%s"
	comment = "%s"
	login_name = "%s_login"
	display_name = "Display Name"
	first_name = "Marcin"
	last_name = "Zukowski"
	email = "fake@email.com"
	disabled = false
	default_warehouse="foo"
	default_role="foo"
	default_secondary_roles=["ALL"]
	default_namespace="foo.bar"
	rsa_public_key = <<KEY
%s
KEY
	rsa_public_key_2 = <<KEY
%s
KEY
	must_change_password = true
}
`
	s = fmt.Sprintf(s, prefix, comment, prefix, key1, key2)
	log.Printf("[DEBUG] s %s", s)
	return s
}

func uConfig2(prefix string) string {
	s := `
resource "snowflake_user" "w" {
	name = "%s"
	comment = "test comment 2"
	password = "best password"
	login_name = "%s_login"
	display_name = "New Name"
	first_name = "Benoit"
	last_name = "Dageville"
	email = "fake@email.net"
	disabled = true
	default_warehouse="bar"
	default_role="bar"
	default_secondary_roles=[]
	default_namespace="bar.baz"
}
`
	s = fmt.Sprintf(s, prefix, prefix)
	log.Printf("[DEBUG] s2 %s", s)
	return s
}

// TestAcc_User_issue2058 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2058 issue.
// The problem was with a dot in user identifier.
// Before the fix it results in panic: interface conversion: sdk.ObjectIdentifier is sdk.DatabaseObjectIdentifier, not sdk.AccountObjectIdentifier error.
func TestAcc_User_issue2058(t *testing.T) {
	r := require.New(t)
	prefix := acc.TestClient().Ids.AlphaContaining(".")
	sshkey1, err := testhelpers.Fixture("userkey1")
	r.NoError(err)
	sshkey2, err := testhelpers.Fixture("userkey2")
	r.NoError(err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: uConfig(prefix, sshkey1, sshkey2, "test_comment"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix),
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
