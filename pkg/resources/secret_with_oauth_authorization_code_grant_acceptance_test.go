package resources_test

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecretWithAuthorizationCodeGrant_BasicFlow(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := random.Comment()
	newComment := random.Comment()
	refreshTokenExpiryDateTime := time.Now().Add(24 * time.Hour).Format(time.DateTime)
	newRefreshTokenExpiryDateOnly := time.Now().Add(4 * 24 * time.Hour).Format(time.DateOnly)
	externallyChangedRefreshTokenExpiryTime := time.Now().Add(10 * 24 * time.Hour).Format(time.DateOnly)
	refreshToken := "test_token"
	newRefreshToken := "new_test_token"

	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := acc.TestClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "foo", "foo"),
	)
	t.Cleanup(apiIntegrationCleanup)

	secretModel := model.SecretWithAuthorizationCodeGrant("s", integrationId.Name(), id.DatabaseName(), id.SchemaName(), name, refreshToken, refreshTokenExpiryDateTime).WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithAuthorizationCodeGrant),
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
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryDateTime).
							HasCommentString(comment),

						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasName(name).
							HasDatabaseName(id.DatabaseName()).
							HasSecretType("OAUTH2").
							HasSchemaName(id.SchemaName()).
							HasComment(comment),
					),

					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttrSet(secretModel.ResourceReference(), "describe_output.0.created_on"),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.name", name),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.secret_type", "OAUTH2"),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.username", ""),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.oauth_access_token_expiry_time", ""),
					resource.TestCheckResourceAttrSet(secretModel.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time"),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.integration_name", integrationId.Name()),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.oauth_scopes.#", "0"),
				),
			},
			// set all
			{
				Config: config.FromModel(t, secretModel.
					WithOauthRefreshTokenExpiryTime(newRefreshTokenExpiryDateOnly).
					WithOauthRefreshToken(newRefreshToken).
					WithComment(newComment),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(secretModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(comment), sdk.String(newComment)),
						planchecks.ExpectChange(secretModel.ResourceReference(), "oauth_refresh_token", tfjson.ActionUpdate, sdk.String(refreshToken), sdk.String(newRefreshToken)),
						planchecks.ExpectChange(secretModel.ResourceReference(), "oauth_refresh_token_expiry_time", tfjson.ActionUpdate, sdk.String(refreshTokenExpiryDateTime), sdk.String(newRefreshTokenExpiryDateOnly)),
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
							HasOauthRefreshTokenExpiryTimeString(newRefreshTokenExpiryDateOnly).
							HasCommentString(newComment),
						assert.Check(resource.TestCheckResourceAttrSet(secretModel.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
			// set comment and refresh_token_expiry_time externally
			{
				PreConfig: func() {
					acc.TestClient().Secret.Alter(t, sdk.NewAlterSecretRequest(id).WithSet(*sdk.NewSecretSetRequest().
						WithComment("secret resource - changed comment").
						WithSetForFlow(*sdk.NewSetForFlowRequest().WithSetForOAuthAuthorization(*sdk.NewSetForOAuthAuthorizationRequest().WithOauthRefreshTokenExpiryTime(time.Now().Add(24 * time.Hour).Format(time.DateOnly)))),
					))
				},
				Config: config.FromModel(t, secretModel.WithOauthRefreshTokenExpiryTime(externallyChangedRefreshTokenExpiryTime).WithComment("comment")),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(secretModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String("secret resource - changed comment"), sdk.String("comment")),
						planchecks.ExpectComputed(secretModel.ResourceReference(), r.DescribeOutputAttributeName, true),
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
							HasOauthRefreshTokenExpiryTimeString(externallyChangedRefreshTokenExpiryTime).
							HasCommentString("comment"),
						assert.Check(resource.TestCheckResourceAttrSet(secretModel.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
			// import
			{
				ResourceName:            secretModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_refresh_token"},
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedSecretWithAuthorizationCodeResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasApiAuthenticationString(integrationId.Name()).
						HasCommentString("comment").
						HasOauthRefreshTokenExpiryTimeNotEmpty(),
				),
			},
		},
	})
}

func TestAcc_SecretWithAuthorizationCodeGrant_DifferentTimeFormats(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := acc.TestClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "foo", "foo"),
	)
	t.Cleanup(apiIntegrationCleanup)

	refreshTokenExpiryDateOnly := time.Now().Add(4 * 24 * time.Hour).Format(time.DateOnly)
	refreshTokenExpiryWithoutSeconds := time.Now().Add(4 * 24 * time.Hour).Format("2006-01-02 15:04")
	refreshTokenExpiryDateTime := time.Now().Add(4 * 24 * time.Hour).Format(time.DateTime)
	refreshTokenExpiryWithPDT := fmt.Sprintf("%s %s", time.Now().Add(4*24*time.Hour).Format("2006-01-02 15:04"), "PDT")

	secretModelDateOnly := model.SecretWithAuthorizationCodeGrant("s", integrationId.Name(), id.DatabaseName(), id.SchemaName(), name, "test_token", refreshTokenExpiryDateOnly)
	secretModelWithoutSeconds := model.SecretWithAuthorizationCodeGrant("s", integrationId.Name(), id.DatabaseName(), id.SchemaName(), name, "test_token", refreshTokenExpiryWithoutSeconds)
	secretModelDateTime := model.SecretWithAuthorizationCodeGrant("s", integrationId.Name(), id.DatabaseName(), id.SchemaName(), name, "test_token", refreshTokenExpiryDateTime)
	secretModelWithPDT := model.SecretWithAuthorizationCodeGrant("s", integrationId.Name(), id.DatabaseName(), id.SchemaName(), name, "test_token", refreshTokenExpiryWithPDT)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			// create with DateOnly
			{
				Config: config.FromModel(t, secretModelDateOnly),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithAuthorizationCodeResource(t, secretModelDateOnly.ResourceReference()).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryDateOnly),
						assert.Check(resource.TestCheckResourceAttrSet(secretModelDateOnly.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
			// update with DateTime without seconds
			{
				Config: config.FromModel(t, secretModelWithoutSeconds),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithAuthorizationCodeResource(t, secretModelWithoutSeconds.ResourceReference()).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryWithoutSeconds),
						assert.Check(resource.TestCheckResourceAttrSet(secretModelWithoutSeconds.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
			// update with DateTime
			{
				Config: config.FromModel(t, secretModelDateTime),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithAuthorizationCodeResource(t, secretModelDateTime.ResourceReference()).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryDateTime),
						assert.Check(resource.TestCheckResourceAttrSet(secretModelDateTime.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
			// update with DateTime with PDT timezone
			{
				Config: config.FromModel(t, secretModelWithPDT),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithAuthorizationCodeResource(t, secretModelWithPDT.ResourceReference()).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryWithPDT),
						assert.Check(resource.TestCheckResourceAttrSet(secretModelWithPDT.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
		},
	})
}

func TestAcc_SecretWithAuthorizationCodeGrant_ExternalChange(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := random.Comment()
	//newComment := random.Comment()
	refreshTokenExpiryDateTime := time.Now().Add(24 * time.Hour).Format(time.DateTime)
	newRefreshTokenExpiryDateOnly := time.Now().Add(4 * 24 * time.Hour).Format(time.DateOnly)
	refreshToken := "test_token"

	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := acc.TestClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "foo", "foo"),
	)
	t.Cleanup(apiIntegrationCleanup)

	secretModel := model.SecretWithAuthorizationCodeGrant("s", integrationId.Name(), id.DatabaseName(), id.SchemaName(), name, refreshToken, refreshTokenExpiryDateTime).WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithAuthorizationCodeGrant),
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
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryDateTime).
							HasCommentString(comment),

						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasName(name).
							HasDatabaseName(id.DatabaseName()).
							HasSecretType("OAUTH2").
							HasSchemaName(id.SchemaName()).
							HasComment(comment),
					),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().Secret.Alter(t, sdk.NewAlterSecretRequest(id).WithSet(*sdk.NewSecretSetRequest().
						WithSetForFlow(*sdk.NewSetForFlowRequest().WithSetForOAuthAuthorization(*sdk.NewSetForOAuthAuthorizationRequest().WithOauthRefreshTokenExpiryTime(time.Now().Add(10 * 24 * time.Hour).Format(time.DateOnly)))),
					))
				},
				Config: config.FromModel(t, secretModel.
					WithOauthRefreshTokenExpiryTime(newRefreshTokenExpiryDateOnly),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectComputed(secretModel.ResourceReference(), r.DescribeOutputAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithAuthorizationCodeResource(t, secretModel.ResourceReference()).
							HasOauthRefreshTokenExpiryTimeString(newRefreshTokenExpiryDateOnly),
						assert.Check(resource.TestCheckResourceAttrSet(secretModel.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
		},
	})
}
