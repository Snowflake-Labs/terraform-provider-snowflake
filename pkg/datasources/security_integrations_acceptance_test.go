package datasources_test

import (
	"fmt"
	"maps"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecurityIntegrations_MultipleTypes(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "1")
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "2")
	issuer := acc.TestClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "https://example.com"
	role := snowflakeroles.GenericScimProvisioner

	saml2Model := model.Saml2SecurityIntegration("test", idOne.Name(), issuer, string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom), validUrl, cert)
	scimModel := model.ScimSecurityIntegration("test", true, idTwo.Name(), role.Name(), string(sdk.ScimSecurityIntegrationScimClientGeneric))
	securityIntegrationsModel := datasourcemodel.SecurityIntegrations("test").
		WithLike(prefix+"%").
		WithDependsOn(saml2Model.ResourceReference(), scimModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, scimModel, saml2Model, securityIntegrationsModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.#", "2"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.name", idTwo.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.integration_type", "SCIM - GENERIC"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.run_as_role.0.value", "GENERIC_SCIM_PROVISIONER"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.sync_password.0.value", "true"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.show_output.0.name", idOne.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.show_output.0.enabled", "false"),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.show_output.0.created_on"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.describe_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.describe_output.0.saml2_issuer.0.value", issuer),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.describe_output.0.saml2_provider.0.value", "CUSTOM"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.describe_output.0.saml2_sso_url.0.value", validUrl),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.1.describe_output.0.saml2_x509_cert.0.value", cert),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_ApiAuthenticationWithAuthorizationCodeGrant(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	pass := random.Password()
	comment := random.Comment()

	resourceModel := model.ApiAuthenticationIntegrationWithAuthorizationCodeGrant("test", true, id.Name(), "foo", pass).
		WithComment(comment).
		WithOauthAccessTokenValidity(42).
		WithOauthAuthorizationEndpoint("https://example.com").
		WithOauthClientAuthMethod(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)).
		WithOauthRefreshTokenValidity(12345).
		WithOauthTokenEndpoint("https://example.com").
		WithOauthAllowedScopesValue(config.SetVariable(config.StringVariable("foo")))
	securityIntegrationsModel := datasourcemodel.SecurityIntegrations("test").
		WithLike(id.Name()).
		WithDependsOn(resourceModel.ResourceReference())
	securityIntegrationsModelWithoutDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(resourceModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, resourceModel, securityIntegrationsModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.#", "1"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.integration_type", "API_AUTHENTICATION"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_access_token_validity.0.value", "42"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "12345"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_id.0.value", "foo"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_auth_method.0.value", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_authorization_endpoint.0.value", "https://example.com"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value", "https://example.com"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_allowed_scopes.0.value", "[foo]"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.parent_integration.0.value", ""),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.auth_type.0.value", "OAUTH2"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.comment.0.value", comment)),
			},
			{
				Config: accconfig.FromModels(t, resourceModel, securityIntegrationsModelWithoutDescribe),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.#", "1"),

					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.integration_type", "API_AUTHENTICATION"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_ApiAuthenticationWithClientCredentials(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	pass1 := random.Password()
	pass2 := random.Password()
	comment := random.Comment()

	resourceModel := model.ApiAuthenticationIntegrationWithClientCredentials("test", true, id.Name(), pass1, pass2).
		WithComment(comment).
		WithOauthAccessTokenValidity(42).
		WithOauthClientAuthMethod(string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)).
		WithOauthRefreshTokenValidity(12345).
		WithOauthTokenEndpoint("https://example.com").
		WithOauthAllowedScopesValue(config.SetVariable(config.StringVariable("foo")))
	securityIntegrationsModel := datasourcemodel.SecurityIntegrations("test").
		WithLike(id.Name()).
		WithDependsOn(resourceModel.ResourceReference())
	securityIntegrationsModelWithoutDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(resourceModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ApiAuthenticationIntegrationWithClientCredentials),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, resourceModel, securityIntegrationsModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.#", "1"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.integration_type", "API_AUTHENTICATION"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_access_token_validity.0.value", "42"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "12345"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_id.0.value", pass1),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_auth_method.0.value", string(sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost)),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value", "https://example.com"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_allowed_scopes.0.value", "[foo]"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_grant.0.value", sdk.ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.parent_integration.0.value", ""),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.auth_type.0.value", "OAUTH2"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.comment.0.value", comment)),
			},
			{
				Config: accconfig.FromModels(t, resourceModel, securityIntegrationsModelWithoutDescribe),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.#", "1"),

					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.integration_type", "API_AUTHENTICATION"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_ExternalOauth(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	issuer := random.String()
	comment := random.Comment()
	claim := random.AlphaN(6)
	mappingAttribute := random.AlphaN(6)
	audience := random.AlphaN(6)

	resourceModel := model.ExternalOauthSecurityIntegration("test", true, issuer, string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress), claim, string(sdk.ExternalOauthSecurityIntegrationTypeCustom), id.Name()).
		WithComment(comment).
		WithExternalOauthAllowedRoles(role.ID()).
		WithExternalOauthAnyRoleMode(string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)).
		WithExternalOauthAudiences(audience).
		WithExternalOauthJwsKeysUrls("https://example.com").
		WithExternalOauthScopeDelimiter(".").
		WithExternalOauthScopeMappingAttribute(mappingAttribute)
	securityIntegrationsModel := datasourcemodel.SecurityIntegrations("test").
		WithLike(id.Name()).
		WithDependsOn(resourceModel.ResourceReference())
	securityIntegrationsModelWithoutDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(resourceModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ExternalOauthSecurityIntegration),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, resourceModel, securityIntegrationsModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.#", "1"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_issuer.0.value", issuer),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_jws_keys_url.0.value", "https://example.com"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_any_role_mode.0.value", string(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_allowed_roles_list.0.value", role.ID().Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_audience_list.0.value", audience),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_token_user_mapping_claim.0.value", fmt.Sprintf("['%s']", claim)),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_snowflake_user_mapping_attribute.0.value", string(sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress)),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.external_oauth_scope_delimiter.0.value", "."),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.comment.0.value", comment),
				),
			},
			{
				Config: accconfig.FromModels(t, resourceModel, securityIntegrationsModelWithoutDescribe),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.#", "1"),

					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.integration_type", "EXTERNAL_OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_OauthForCustomClients(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)

	preAuthorizedRole, preauthorizedRoleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(preauthorizedRoleCleanup)

	blockedRole, blockedRoleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(blockedRoleCleanup)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	validUrl := "https://example.com"
	key, _ := random.GenerateRSAPublicKey(t)
	comment := random.Comment()

	resourceModel := model.OauthIntegrationForCustomClients("test", id.Name(), string(sdk.OauthSecurityIntegrationClientTypeConfidential), validUrl).
		WithComment(comment).
		WithEnabled(datasources.BooleanTrue).
		WithBlockedRolesList("ACCOUNTADMIN", "SECURITYADMIN", blockedRole.ID().Name()).
		WithNetworkPolicy(networkPolicy.ID().Name()).
		WithOauthAllowNonTlsRedirectUri(datasources.BooleanTrue).
		WithOauthClientRsaPublicKey(key).
		WithOauthClientRsaPublicKey2(key).
		WithOauthEnforcePkce(datasources.BooleanTrue).
		WithOauthIssueRefreshTokens(datasources.BooleanTrue).
		WithOauthRefreshTokenValidity(86400).
		WithOauthUseSecondaryRoles(string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)).
		WithPreAuthorizedRoles(preAuthorizedRole.ID())
	securityIntegrationsModel := datasourcemodel.SecurityIntegrations("test").
		WithLike(id.Name()).
		WithDependsOn(resourceModel.ResourceReference())
	securityIntegrationsModelWithoutDescribe := datasourcemodel.SecurityIntegrations("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(resourceModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForCustomClients),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, resourceModel, securityIntegrationsModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.#", "1"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypeConfidential)),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_redirect_uri.0.value", validUrl),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_allow_non_tls_redirect_uri.0.value", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_enforce_pkce.0.value", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesNone)),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.pre_authorized_roles_list.0.value", preAuthorizedRole.ID().Name()),
					// Not asserted, because it also contains other default roles
					// resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.blocked_roles_list.0.value"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_issue_refresh_tokens.0.value", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.network_policy.0.value", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_rsa_public_key_fp.0.value"),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_rsa_public_key_2_fp.0.value"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.comment.0.value", comment),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.describe_output.0.oauth_allowed_token_endpoints.0.value"),

					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(securityIntegrationsModel.DatasourceReference(), "security_integrations.0.show_output.0.created_on"),
				),
			},
			{
				Config: accconfig.FromModels(t, resourceModel, securityIntegrationsModelWithoutDescribe),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.#", "1"),

					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.integration_type", "OAUTH - CUSTOM"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr(securityIntegrationsModelWithoutDescribe.DatasourceReference(), "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_OauthForPartnerApplications(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                         config.StringVariable(id.Name()),
			"oauth_client":                 config.StringVariable(string(sdk.OauthSecurityIntegrationClientTableauServer)),
			"blocked_roles_list":           config.SetVariable(config.StringVariable("ACCOUNTADMIN"), config.StringVariable("SECURITYADMIN")),
			"enabled":                      config.BoolVariable(true),
			"oauth_issue_refresh_tokens":   config.BoolVariable(false),
			"oauth_refresh_token_validity": config.IntegerVariable(86400),
			"oauth_use_secondary_roles":    config.StringVariable(string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
			"comment":                      config.StringVariable(comment),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.OauthIntegrationForPartnerApplications),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/oauth_for_partner_applications/optionals_set"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_type.0.value", string(sdk.OauthSecurityIntegrationClientTypePublic)),
					resource.TestCheckNoResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_redirect_uri.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_use_secondary_roles.0.value", string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.blocked_roles_list.0.value", "ACCOUNTADMIN,SECURITYADMIN"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_issue_refresh_tokens.0.value", "false"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_refresh_token_validity.0.value", "86400"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.comment.0.value", comment),
					resource.TestCheckNoResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_client_id.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_authorization_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_token_endpoint.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_allowed_authorization_endpoints.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.oauth_allowed_token_endpoints.0.value"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "OAUTH - TABLEAU_SERVER"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/oauth_for_partner_applications/optionals_unset"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "OAUTH - TABLEAU_SERVER"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_Saml2(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	issuer := acc.TestClient().Ids.Alpha()
	cert := random.GenerateX509(t)
	validUrl := "http://example.com"
	acsURL := acc.TestClient().Context.ACSURL(t)
	issuerURL := acc.TestClient().Context.IssuerURL(t)

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"allowed_email_patterns":              config.ListVariable(config.StringVariable("^(.+dev)@example.com$")),
			"allowed_user_domains":                config.ListVariable(config.StringVariable("example.com")),
			"comment":                             config.StringVariable("foo"),
			"enabled":                             config.BoolVariable(true),
			"name":                                config.StringVariable(id.Name()),
			"saml2_enable_sp_initiated":           config.BoolVariable(true),
			"saml2_force_authn":                   config.BoolVariable(true),
			"saml2_issuer":                        config.StringVariable(issuer),
			"saml2_post_logout_redirect_url":      config.StringVariable(validUrl),
			"saml2_provider":                      config.StringVariable(string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
			"saml2_requested_nameid_format":       config.StringVariable(string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
			"saml2_sign_request":                  config.BoolVariable(true),
			"saml2_snowflake_acs_url":             config.StringVariable(acsURL),
			"saml2_snowflake_issuer_url":          config.StringVariable(issuerURL),
			"saml2_sp_initiated_login_page_label": config.StringVariable("foo"),
			"saml2_sso_url":                       config.StringVariable(validUrl),
			"saml2_x509_cert":                     config.StringVariable(cert),
			// TODO(SNOW-1479617): set saml2_snowflake_x509_cert
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Saml2SecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/saml2/optionals_set"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_issuer.0.value", issuer),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_sso_url.0.value", validUrl),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_provider.0.value", string(sdk.Saml2SecurityIntegrationSaml2ProviderCustom)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_x509_cert.0.value", cert),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_sp_initiated_login_page_label.0.value", "foo"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_enable_sp_initiated.0.value", "true"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_snowflake_x509_cert.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_sign_request.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_requested_nameid_format.0.value", string(sdk.Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified)),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_post_logout_redirect_url.0.value", "http://example.com"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_force_authn.0.value", "true"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_snowflake_issuer_url.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_snowflake_acs_url.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_snowflake_metadata.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_digest_methods_used.0.value"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.saml2_signature_methods_used.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.allowed_user_domains.0.value", "[example.com]"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.allowed_email_patterns.0.value", "[^(.+dev)@example.com$]"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.comment.0.value", "foo"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/saml2/optionals_unset"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "SAML2"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_Scim(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)
	configVariables := config.Variables{
		"name":           config.StringVariable(id.Name()),
		"comment":        config.StringVariable(comment),
		"network_policy": config.StringVariable(networkPolicy.ID().Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/optionals_set"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "SCIM - GENERIC"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.network_policy.0.value", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.run_as_role.0.value", "GENERIC_SCIM_PROVISIONER"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.sync_password.0.value", "true"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.comment.0.value", comment),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/optionals_unset"),
				ConfigVariables: configVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "SCIM - GENERIC"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_Filtering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	commonVariables := config.Variables{
		"name_1": config.StringVariable(idOne.Name()),
		"name_2": config.StringVariable(idTwo.Name()),
		"name_3": config.StringVariable(idThree.Name()),
	}

	likeConfig := config.Variables{
		"like": config.StringVariable(idOne.Name()),
	}
	maps.Copy(likeConfig, commonVariables)

	likeConfig2 := config.Variables{
		"like": config.StringVariable(prefix + "%"),
	}
	maps.Copy(likeConfig2, commonVariables)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/like"),
				ConfigVariables: likeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/like"),
				ConfigVariables: likeConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "2"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_SecurityIntegrationNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one security integration"),
			},
		},
	})
}
