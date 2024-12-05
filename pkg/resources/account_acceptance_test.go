package resources_test

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Account_minimal(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := acc.TestClient().Context.CurrentAccountId(t).OrganizationName()
	id := random.AdminName()
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)
	region := acc.TestClient().Context.CurrentRegion(t)

	configModel := model.Account("test", name, string(sdk.UserTypeService), string(sdk.EditionStandard), email, 3, id).
		WithAdminRsaPublicKey(key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id).
						HasFullyQualifiedNameString(sdk.NewAccountIdentifier(organizationName, id).FullyQualifiedName()).
						HasAdminNameString(name).
						HasAdminRsaPublicKeyString(key).
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
						HasAccountName(id).
						HasSnowflakeRegion(region).
						HasRegionGroup("").
						HasEdition(sdk.EditionStandard).
						//HasAccountURL().
						//HasCreatedOn().
						HasComment("SNOWFLAKE").
						//HasAccountLocator().
						//HasAccountLocatorURL().
						HasManagedAccounts(0).
						//HasConsumptionBillingEntityName().
						//HasMarketplaceConsumerBillingEntityName().
						//HasMarketplaceProviderBillingEntityName().
						//HasOldAccountURL().
						HasIsOrgAdmin(false).
						//HasAccountOldUrlSavedOn().
						//HasAccountOldUrlLastUsed().
						//HasOrganizationOldUrl().
						//HasOrganizationOldUrlSavedOn().
						//HasOrganizationOldUrlLastUsed().
						HasIsEventsAccount(false).
						HasIsOrganizationAccount(false),
					//HasDroppedOn().
					//HasScheduledDeletionTime().
					//HasRestoredOn().
					//HasMovedToOrganization().
					//HasMovedOn().
					//HasOrganizationUrlExpirationOn(),
				),
			},
		},
	})
}

func TestAcc_Account_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	organizationName := acc.TestClient().Context.CurrentAccountId(t).OrganizationName()
	id := random.AdminName()
	firstName := acc.TestClient().Ids.Alpha()
	lastName := acc.TestClient().Ids.Alpha()
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)
	region := acc.TestClient().Context.CurrentRegion(t)
	comment := random.Comment()

	configModel := model.Account("test", name, string(sdk.UserTypePerson), string(sdk.EditionStandard), email, 3, id).
		WithAdminRsaPublicKey(key).
		WithFirstName(firstName).
		WithLastName(lastName).
		WithMustChangePassword(r.BooleanTrue).
		//WithRegionGroup("PUBLIC").
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
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id).
						HasFullyQualifiedNameString(sdk.NewAccountIdentifier(organizationName, id).FullyQualifiedName()).
						HasAdminNameString(name).
						HasAdminRsaPublicKeyString(key).
						HasEmailString(email).
						HasFirstNameString(firstName).
						HasLastNameString(lastName).
						HasMustChangePasswordString(r.BooleanTrue).
						HasNoRegionGroup(). // TODO
						HasRegionString(region).
						HasCommentString(comment).
						HasIsOrgAdminString(r.BooleanFalse).
						HasGracePeriodInDaysString("3"),
					resourceshowoutputassert.AccountShowOutput(t, configModel.ResourceReference()).
						HasOrganizationName(organizationName).
						HasAccountName(id).
						HasSnowflakeRegion(region).
						HasRegionGroup("").
						HasEdition(sdk.EditionStandard).
						//HasAccountURL().
						//HasCreatedOn().
						HasComment(comment).
						//HasAccountLocator().
						//HasAccountLocatorURL().
						HasManagedAccounts(0).
						//HasConsumptionBillingEntityName().
						//HasMarketplaceConsumerBillingEntityName().
						//HasMarketplaceProviderBillingEntityName().
						//HasOldAccountURL().
						HasIsOrgAdmin(false).
						//HasAccountOldUrlSavedOn().
						//HasAccountOldUrlLastUsed().
						//HasOrganizationOldUrl().
						//HasOrganizationOldUrlSavedOn().
						//HasOrganizationOldUrlLastUsed().
						HasIsEventsAccount(false).
						HasIsOrganizationAccount(false),
					//HasDroppedOn().
					//HasScheduledDeletionTime().
					//HasRestoredOn().
					//HasMovedToOrganization().
					//HasMovedOn().
					//HasOrganizationUrlExpirationOn(),
				),
			},
		},
	})
}

// TODO: All show outputs in minimal and complete
// TODO: Imports
// TODO: Alters
// TODO: Not orgadmin role
// TODO: Invalid values
// TODO: State upgrader
