//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Account_Minimal(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := acc.TestClient().Context.CurrentAccountId(t).OrganizationName()
	id := acc.TestClient().Ids.RandomSensitiveAccountObjectIdentifier()
	accountId := sdk.NewAccountIdentifier(organizationName, id.Name())
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)
	region := acc.TestClient().Context.CurrentRegion(t)

	configModel := model.Account("test", name, string(sdk.EditionStandard), email, 3, id.Name()).
		WithAdminRsaPublicKey(key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminNameString(name).
						HasAdminRsaPublicKeyString(key).
						HasNoAdminUserType().
						HasEmailString(email).
						HasNoFirstName().
						HasNoLastName().
						HasMustChangePasswordString(r.BooleanDefault).
						HasNoRegionGroup().
						HasNoRegion().
						HasNoComment().
						HasIsOrgAdminString(r.BooleanDefault).
						HasGracePeriodInDaysString("3"),
					resourceshowoutputassert.AccountShowOutput(t, configModel.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasSnowflakeRegion(region).
						HasRegionGroup("").
						HasEdition(sdk.EditionStandard).
						HasAccountUrlNotEmpty().
						HasCreatedOnNotEmpty().
						HasComment("SNOWFLAKE").
						HasAccountLocatorNotEmpty().
						HasAccountLocatorUrlNotEmpty().
						HasManagedAccounts(0).
						HasConsumptionBillingEntityNameNotEmpty().
						HasMarketplaceConsumerBillingEntityName("").
						HasMarketplaceProviderBillingEntityNameNotEmpty().
						HasOldAccountURL("").
						HasIsOrgAdmin(false).
						HasAccountOldUrlSavedOnEmpty().
						HasAccountOldUrlLastUsedEmpty().
						HasOrganizationOldUrl("").
						HasOrganizationOldUrlSavedOnEmpty().
						HasOrganizationOldUrlLastUsedEmpty().
						HasIsEventsAccount(false).
						HasIsOrganizationAccount(false).
						HasDroppedOnEmpty().
						HasScheduledDeletionTimeEmpty().
						HasRestoredOnEmpty().
						HasMovedToOrganization("").
						HasMovedOn("").
						HasOrganizationUrlExpirationOnEmpty(),
				),
			},
			{
				ResourceName: configModel.ResourceReference(),
				Config:       config.FromModels(t, configModel),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedAccountResource(t, helpers.EncodeResourceIdentifier(accountId)).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasNoAdminName().
						HasNoAdminRsaPublicKey().
						HasNoAdminUserType().
						HasNoEmail().
						HasNoFirstName().
						HasNoLastName().
						HasNoMustChangePassword().
						HasEditionString(string(sdk.EditionStandard)).
						HasNoRegionGroup().
						HasRegionString(region).
						HasCommentString("SNOWFLAKE").
						HasIsOrgAdminString(r.BooleanFalse).
						HasNoGracePeriodInDays(),
				),
			},
		},
	})
}

func TestAcc_Account_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := acc.TestClient().Context.CurrentAccountId(t).OrganizationName()
	id := acc.TestClient().Ids.RandomSensitiveAccountObjectIdentifier()
	accountId := sdk.NewAccountIdentifier(organizationName, id.Name())
	firstName := random.AlphaN(30)
	lastName := random.AlphaN(30)
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)
	region := acc.TestClient().Context.CurrentRegion(t)
	comment := random.Comment()

	configModel := model.Account("test", name, string(sdk.EditionStandard), email, 3, id.Name()).
		WithAdminUserTypeEnum(sdk.UserTypePerson).
		WithAdminRsaPublicKey(key).
		WithFirstName(firstName).
		WithLastName(lastName).
		WithMustChangePassword(r.BooleanTrue).
		WithRegionGroup("PUBLIC").
		WithRegion(region).
		WithComment(comment).
		WithIsOrgAdmin(r.BooleanFalse)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(sdk.NewAccountIdentifier(organizationName, id.Name()).FullyQualifiedName()).
						HasAdminNameString(name).
						HasAdminRsaPublicKeyString(key).
						HasAdminUserType(sdk.UserTypePerson).
						HasEmailString(email).
						HasFirstNameString(firstName).
						HasLastNameString(lastName).
						HasMustChangePasswordString(r.BooleanTrue).
						HasRegionGroupString("PUBLIC").
						HasRegionString(region).
						HasCommentString(comment).
						HasIsOrgAdminString(r.BooleanFalse).
						HasGracePeriodInDaysString("3"),
					resourceshowoutputassert.AccountShowOutput(t, configModel.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasSnowflakeRegion(region).
						HasRegionGroup("").
						HasEdition(sdk.EditionStandard).
						HasAccountUrlNotEmpty().
						HasCreatedOnNotEmpty().
						HasComment(comment).
						HasAccountLocatorNotEmpty().
						HasAccountLocatorUrlNotEmpty().
						HasManagedAccounts(0).
						HasConsumptionBillingEntityNameNotEmpty().
						HasMarketplaceConsumerBillingEntityName("").
						HasMarketplaceProviderBillingEntityNameNotEmpty().
						HasOldAccountURL("").
						HasIsOrgAdmin(false).
						HasAccountOldUrlSavedOnEmpty().
						HasAccountOldUrlLastUsedEmpty().
						HasOrganizationOldUrl("").
						HasOrganizationOldUrlSavedOnEmpty().
						HasOrganizationOldUrlLastUsedEmpty().
						HasIsEventsAccount(false).
						HasIsOrganizationAccount(false).
						HasDroppedOnEmpty().
						HasScheduledDeletionTimeEmpty().
						HasRestoredOnEmpty().
						HasMovedToOrganization("").
						HasMovedOn("").
						HasOrganizationUrlExpirationOnEmpty(),
				),
			},
			{
				ResourceName: configModel.ResourceReference(),
				Config:       config.FromModels(t, configModel),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedAccountResource(t, helpers.EncodeResourceIdentifier(accountId)).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(sdk.NewAccountIdentifier(organizationName, id.Name()).FullyQualifiedName()).
						HasNoAdminName().
						HasNoAdminRsaPublicKey().
						HasNoEmail().
						HasNoFirstName().
						HasNoLastName().
						HasNoAdminUserType().
						HasNoMustChangePassword().
						HasEditionString(string(sdk.EditionStandard)).
						HasNoRegionGroup().
						HasRegionString(region).
						HasCommentString(comment).
						HasIsOrgAdminString(r.BooleanFalse).
						HasNoGracePeriodInDays(),
				),
			},
		},
	})
}

func TestAcc_Account_Rename(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := acc.TestClient().Context.CurrentAccountId(t).OrganizationName()
	id := acc.TestClient().Ids.RandomSensitiveAccountObjectIdentifier()
	accountId := sdk.NewAccountIdentifier(organizationName, id.Name())

	newId := acc.TestClient().Ids.RandomSensitiveAccountObjectIdentifier()
	newAccountId := sdk.NewAccountIdentifier(organizationName, newId.Name())

	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)

	configModel := model.Account("test", name, string(sdk.EditionStandard), email, 3, id.Name()).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)
	newConfigModel := model.Account("test", name, string(sdk.EditionStandard), email, 3, newId.Name()).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService),
					resourceshowoutputassert.AccountShowOutput(t, configModel.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(newConfigModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, newConfigModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, newConfigModel.ResourceReference()).
						HasNameString(newId.Name()).
						HasFullyQualifiedNameString(newAccountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService),
					resourceshowoutputassert.AccountShowOutput(t, newConfigModel.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(newId.Name()),
				),
			},
		},
	})
}

func TestAcc_Account_IsOrgAdmin(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := acc.TestClient().Context.CurrentAccountId(t).OrganizationName()
	id := acc.TestClient().Ids.RandomSensitiveAccountObjectIdentifier()
	accountId := sdk.NewAccountIdentifier(organizationName, id.Name())

	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)

	configModelWithOrgAdminTrue := model.Account("test", name, string(sdk.EditionStandard), email, 3, id.Name()).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key).
		WithIsOrgAdmin(r.BooleanTrue)

	configModelWithOrgAdminFalse := model.Account("test", name, string(sdk.EditionStandard), email, 3, id.Name()).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key).
		WithIsOrgAdmin(r.BooleanFalse)

	configModelWithoutOrgAdmin := model.Account("test", name, string(sdk.EditionStandard), email, 3, id.Name()).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			// Create with ORGADMIN enabled
			{
				Config: config.FromModels(t, configModelWithOrgAdminTrue),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModelWithOrgAdminTrue.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService).
						HasIsOrgAdminString(r.BooleanTrue),
					resourceshowoutputassert.AccountShowOutput(t, configModelWithOrgAdminTrue.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasIsOrgAdmin(true),
				),
			},
			// Disable ORGADMIN
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configModelWithOrgAdminFalse.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, configModelWithOrgAdminFalse),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModelWithOrgAdminFalse.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService).
						HasIsOrgAdminString(r.BooleanFalse),
					resourceshowoutputassert.AccountShowOutput(t, configModelWithOrgAdminFalse.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasIsOrgAdmin(false),
				),
			},
			// Remove is_org_admin from the config and go back to default (disabled)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configModelWithoutOrgAdmin.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, configModelWithoutOrgAdmin),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModelWithoutOrgAdmin.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService).
						HasIsOrgAdminString(r.BooleanDefault),
					resourceshowoutputassert.AccountShowOutput(t, configModelWithoutOrgAdmin.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasIsOrgAdmin(false),
				),
			},
			// External change (enable ORGADMIN)
			{
				PreConfig: func() {
					acc.TestClient().Account.Alter(t, &sdk.AlterAccountOptions{
						SetIsOrgAdmin: &sdk.AccountSetIsOrgAdmin{
							Name:     accountId.AsAccountObjectIdentifier(),
							OrgAdmin: true,
						},
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configModelWithoutOrgAdmin.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, configModelWithoutOrgAdmin),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModelWithoutOrgAdmin.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminUserType(sdk.UserTypeService).
						HasIsOrgAdminString(r.BooleanDefault),
					resourceshowoutputassert.AccountShowOutput(t, configModelWithoutOrgAdmin.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id.Name()).
						HasIsOrgAdmin(false),
				),
			},
		},
	})
}

func TestAcc_Account_IgnoreUpdateAfterCreationOnCertainFields(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := acc.TestClient().Context.CurrentAccountId(t).OrganizationName()
	id := acc.TestClient().Ids.RandomSensitiveAccountObjectIdentifier()
	accountId := sdk.NewAccountIdentifier(organizationName, id.Name())

	firstName := random.AlphaN(30)
	lastName := random.AlphaN(30)
	email := random.Email()
	name := random.AdminName()
	pass := random.Password()

	newFirstName := random.AlphaN(30)
	newLastName := random.AlphaN(30)
	newEmail := random.Email()
	newName := random.AdminName()
	newPass := random.Password()

	configModel := model.Account("test", name, string(sdk.EditionStandard), email, 3, id.Name()).
		WithAdminUserTypeEnum(sdk.UserTypePerson).
		WithFirstName(firstName).
		WithLastName(lastName).
		WithMustChangePassword(r.BooleanTrue).
		WithAdminPassword(pass)

	newConfigModel := model.Account("test", newName, string(sdk.EditionStandard), newEmail, 3, id.Name()).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminPassword(newPass).
		WithFirstName(newFirstName).
		WithLastName(newLastName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminNameString(name).
						HasAdminPasswordString(pass).
						HasAdminUserType(sdk.UserTypePerson).
						HasEmailString(email).
						HasFirstNameString(firstName).
						HasLastNameString(lastName).
						HasMustChangePasswordString(r.BooleanTrue),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(newConfigModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, newConfigModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, newConfigModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(accountId.FullyQualifiedName()).
						HasAdminNameString(name).
						HasAdminPasswordString(pass).
						HasAdminUserType(sdk.UserTypePerson).
						HasEmailString(email).
						HasFirstNameString(firstName).
						HasLastNameString(lastName).
						HasMustChangePasswordString(r.BooleanTrue),
				),
			},
		},
	})
}

func TestAcc_Account_TryToCreateWithoutOrgadmin(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	id := acc.TestClient().Ids.RandomSensitiveAccountObjectIdentifier()
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)

	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	t.Setenv(snowflakeenvs.Role, snowflakeroles.Accountadmin.Name())

	configModel := model.Account("test", name, string(sdk.EditionStandard), email, 3, id.Name()).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, configModel),
				ExpectError: regexp.MustCompile("Error: current user doesn't have the orgadmin role in session"),
			},
		},
	})
}

func TestAcc_Account_InvalidValues(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	id := acc.TestClient().Ids.RandomSensitiveAccountObjectIdentifier()
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)

	configModelInvalidUserType := model.Account("test", name, string(sdk.EditionStandard), email, 3, id.Name()).
		WithAdminUserType("invalid_user_type").
		WithAdminRsaPublicKey(key)

	configModelInvalidAccountEdition := model.Account("test", name, "invalid_account_edition", email, 3, id.Name()).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)

	configModelInvalidGracePeriodInDays := model.Account("test", name, string(sdk.EditionStandard), email, 2, id.Name()).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminRsaPublicKey(key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, configModelInvalidUserType),
				ExpectError: regexp.MustCompile("invalid user type: invalid_user_type"),
			},
			{
				Config:      config.FromModels(t, configModelInvalidAccountEdition),
				ExpectError: regexp.MustCompile("unknown account edition: invalid_account_edition"),
			},
			{
				Config:      config.FromModels(t, configModelInvalidGracePeriodInDays),
				ExpectError: regexp.MustCompile(`Error: expected grace_period_in_days to be at least \(3\), got 2`),
			},
		},
	})
}

func TestAcc_Account_UpgradeFrom_v0_99_0(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	id := acc.TestClient().Ids.RandomSensitiveAccountObjectIdentifier()
	email := random.Email()
	adminName := random.AdminName()
	adminPassword := random.Password()
	firstName := random.AlphaN(30)
	lastName := random.AlphaN(30)
	region := acc.TestClient().Context.CurrentRegion(t)
	comment := random.Comment()

	configModel := model.Account("test", adminName, string(sdk.EditionStandard), email, 3, id.Name()).
		WithAdminUserTypeEnum(sdk.UserTypeService).
		WithAdminPassword(adminPassword).
		WithFirstName(firstName).
		WithLastName(lastName).
		WithMustChangePasswordValue(tfconfig.BoolVariable(true)).
		WithRegion(region).
		WithIsOrgAdmin(r.BooleanFalse).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.99.0"),
				Config:            accountConfig_v0_99_0(id.Name(), adminName, adminPassword, email, sdk.EditionStandard, firstName, lastName, true, region, 3, comment),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasAdminNameString(adminName).
						HasAdminPasswordString(adminPassword).
						HasEmailString(email).
						HasFirstNameString(firstName).
						HasLastNameString(lastName).
						HasMustChangePasswordString(r.BooleanTrue).
						HasRegionGroupString("").
						HasRegionString(region).
						HasCommentString(comment).
						HasIsOrgAdminString(r.BooleanFalse).
						HasGracePeriodInDaysString("3"),
				),
			},
		},
	})
}

func accountConfig_v0_99_0(
	name string,
	adminName string,
	adminPassword string,
	email string,
	edition sdk.AccountEdition,
	firstName string,
	lastName string,
	mustChangePassword bool,
	region string,
	gracePeriodInDays int,
	comment string,
) string {
	return fmt.Sprintf(`
resource "snowflake_account" "test" {
	name = "%[1]s"
	admin_name = "%[2]s"
	admin_password = "%[3]s"
	email = "%[4]s"
	edition = "%[5]s"
	first_name = "%[6]s"
	last_name = "%[7]s"
	must_change_password = %[8]t
	region = "%[9]s"
	grace_period_in_days = %[10]d
	comment = "%[11]s"
}
`,
		name,
		adminName,
		adminPassword,
		email,
		edition,
		firstName,
		lastName,
		mustChangePassword,
		region,
		gracePeriodInDays,
		comment,
	)
}

// TODO(SNOW-1875369): add a state upgrader test for an imported account with optional parameters
