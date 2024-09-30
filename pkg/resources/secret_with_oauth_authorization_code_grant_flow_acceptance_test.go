package resources_test

import (
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"testing"
	"time"
)

func TestAcc_SecretWithAuthorizationCode_BasicFlow(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := random.Comment()
	newComment := random.Comment()
	refreshTokenExpiryTime := time.Now().Add(24 * time.Hour).Format(time.DateOnly)
	newRefreshTokenExpiryTime := time.Now().Add(4 * 24 * time.Hour).Format(time.DateOnly)
	refreshToken := "test_token"
	newRefreshToken := "new_test_token"

	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := acc.TestClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "foo", "foo"),
	)
	t.Cleanup(apiIntegrationCleanup)

	secretModel := model.SecretWithAuthorizationCode("s", integrationId.Name(), id.DatabaseName(), id.SchemaName(), name, refreshToken, refreshTokenExpiryTime).WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithAuthorizationCode),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModel(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithAuthorizationCodeResource(t, secretModel.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(integrationId.Name()).
							HasOauthRefreshTokenString(refreshToken).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryTime).
							HasCommentString(comment),
					),
				),
			},
			// set all
			{
				Config: config.FromModel(t, secretModel.
					WithOauthRefreshTokenExpiryTime(newRefreshTokenExpiryTime).
					WithOauthRefreshToken(newRefreshToken).
					WithComment(newComment),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(secretModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(comment), sdk.String(newComment)),
						planchecks.ExpectChange(secretModel.ResourceReference(), "oauth_refresh_token", tfjson.ActionUpdate, sdk.String(refreshToken), sdk.String(newRefreshToken)),
						planchecks.ExpectChange(secretModel.ResourceReference(), "oauth_refresh_token_expiry_time", tfjson.ActionUpdate, sdk.String(refreshTokenExpiryTime), sdk.String(newRefreshTokenExpiryTime)),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithAuthorizationCodeResource(t, secretModel.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(integrationId.Name()).
							HasOauthRefreshTokenString(newRefreshToken).
							HasOauthRefreshTokenExpiryTimeString(newRefreshTokenExpiryTime).
							HasCommentString(newComment),
					),
				),
			},
			// import
			{
				ResourceName:            secretModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_refresh_token"},
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "oauth_refresh_token_expiry_time", newRefreshTokenExpiryTime),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", newComment),
				),
			},
		},
	})
}
