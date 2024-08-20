package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-cty/cty"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var apiAuthAuthorizationCodeGrantSchema = func() map[string]*schema.Schema {
	apiAuthAuthorizationCodeGrant := map[string]*schema.Schema{
		"oauth_authorization_endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the URL for authenticating to the external service. If removed from the config, the resource is recreated.",
		},
		"oauth_allowed_scopes": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "Specifies a list of scopes to use when making a request from the OAuth by a role with USAGE on the integration during the OAuth client credentials flow.",
		},
	}
	return helpers.MergeMaps(apiAuthCommonSchema, apiAuthAuthorizationCodeGrant)
}()

func ApiAuthenticationIntegrationWithAuthorizationCodeGrant() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextApiAuthenticationIntegrationWithAuthorizationCodeGrant,
		ReadContext:   ReadContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(true),
		UpdateContext: UpdateContextApiAuthenticationIntegrationWithAuthorizationCodeGrant,
		DeleteContext: DeleteContextApiAuthenticationIntegrationWithAuthorizationCodeGrant,
		Description:   "Resource used to manage api authentication security integration objects with authorization code grant. For more information, check [security integrations documentation](https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth).",

		CustomizeDiff: customdiff.All(
			ForceNewIfChangeToEmptyString("oauth_token_endpoint"),
			ForceNewIfChangeToEmptyString("oauth_authorization_endpoint"),
			ForceNewIfChangeToEmptyString("oauth_client_auth_method"),
			ComputedIfAnyAttributeChanged(ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(DescribeOutputAttributeName, "enabled", "comment", "oauth_access_token_validity", "oauth_refresh_token_validity",
				"oauth_client_id", "oauth_client_auth_method", "oauth_authorization_endpoint",
				"oauth_token_endpoint", "oauth_allowed_scopes"),
		),
		Schema: apiAuthAuthorizationCodeGrantSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportApiAuthenticationWithAuthorizationCodeGrant,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName,
			},
		},
	}
}

func ImportApiAuthenticationWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting api auth integration with authorization code grant import")
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	properties, err := client.SecurityIntegrations.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := handleApiAuthImport(d, integration, properties); err != nil {
		return nil, err
	}
	oauthAuthorizationEndpoint, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_AUTHORIZATION_ENDPOINT"
	})
	if err == nil {
		if err = d.Set("oauth_authorization_endpoint", oauthAuthorizationEndpoint.Value); err != nil {
			return nil, err
		}
	}
	oauthAllowedScopes, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_ALLOWED_SCOPES" })
	if err == nil {
		if err = d.Set("oauth_allowed_scopes", sdk.ParseCommaSeparatedStringArray(oauthAllowedScopes.Value, false)); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	commonCreate, err := handleApiAuthCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := sdk.ParseAccountObjectIdentifier(commonCreate.name)
	if err != nil {
		return diag.FromErr(err)
	}

	req := sdk.NewCreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest(id, commonCreate.enabled, commonCreate.oauthClientId, commonCreate.oauthClientSecret)
	req.WithOauthGrantAuthorizationCode(true)
	req.Comment = commonCreate.comment
	req.OauthAccessTokenValidity = commonCreate.oauthAccessTokenValidity
	req.OauthRefreshTokenValidity = commonCreate.oauthRefreshTokenValidity
	req.OauthTokenEndpoint = commonCreate.oauthTokenEndpoint
	req.OauthClientAuthMethod = commonCreate.oauthClientAuthMethod

	if v, ok := d.GetOk("oauth_authorization_endpoint"); ok {
		req.WithOauthAuthorizationEndpoint(v.(string))
	}

	if v, ok := d.GetOk("oauth_allowed_scopes"); ok {
		elems := expandStringList(v.(*schema.Set).List())
		allowedScopes := make([]sdk.AllowedScope, len(elems))
		for i := range elems {
			allowedScopes[i] = sdk.AllowedScope{Scope: elems[i]}
		}
		req.WithOauthAllowedScopes(allowedScopes)
	}

	if err := client.SecurityIntegrations.CreateApiAuthenticationWithAuthorizationCodeGrantFlow(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(false)(ctx, d, meta)
}

func ReadContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query security integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Security integration name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}
		properties, err := client.SecurityIntegrations.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if c := integration.Category; c != sdk.SecurityIntegrationCategory {
			return diag.FromErr(fmt.Errorf("expected %v to be a %s integration, got %v", id, sdk.SecurityIntegrationCategory, c))
		}
		oauthAuthorizationEndpoint, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool {
			return property.Name == "OAUTH_AUTHORIZATION_ENDPOINT"
		})
		if err != nil {
			return diag.FromErr(err)
		}

		oauthAllowedScopes, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_ALLOWED_SCOPES" })
		if err != nil {
			return diag.FromErr(err)
		}

		if err := handleApiAuthRead(d, id, integration, properties, withExternalChangesMarking, []describeMapping{
			{"oauth_authorization_endpoint", "oauth_authorization_endpoint", oauthAuthorizationEndpoint.Value, oauthAuthorizationEndpoint.Value, nil},
			{"oauth_allowed_scopes", "oauth_allowed_scopes", oauthAllowedScopes.Value, sdk.ParseCommaSeparatedStringArray(oauthAllowedScopes.Value, false), nil},
		}); err != nil {
			return diag.FromErr(err)
		}
		if err := setStateToValuesFromConfig(d, apiAuthAuthorizationCodeGrantSchema, []string{
			"oauth_authorization_endpoint",
			"oauth_allowed_scopes",
			"oauth_client_auth_method",
			"oauth_token_endpoint",
		}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
}

func UpdateContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	commonSet, commonUnset, err := handleApiAuthUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}
	set := &sdk.ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSetRequest{
		Enabled:                   commonSet.enabled,
		OauthTokenEndpoint:        commonSet.oauthTokenEndpoint,
		OauthClientAuthMethod:     commonSet.oauthClientAuthMethod,
		OauthClientId:             commonSet.oauthClientId,
		OauthClientSecret:         commonSet.oauthClientSecret,
		OauthAccessTokenValidity:  commonSet.oauthAccessTokenValidity,
		OauthRefreshTokenValidity: commonSet.oauthRefreshTokenValidity,
		Comment:                   commonSet.comment,
	}
	unset := &sdk.ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationUnsetRequest{
		Comment: commonUnset.comment,
	}
	if d.HasChange("oauth_authorization_endpoint") {
		set.WithOauthAuthorizationEndpoint(d.Get("oauth_authorization_endpoint").(string))
	}

	if d.HasChange("oauth_allowed_scopes") {
		elems := expandStringList(d.Get("oauth_allowed_scopes").(*schema.Set).List())
		allowedScopes := make([]sdk.AllowedScope, len(elems))
		for i := range elems {
			allowedScopes[i] = sdk.AllowedScope{Scope: elems[i]}
		}
		set.WithOauthAllowedScopes(allowedScopes)
	}

	if !reflect.DeepEqual(*set, sdk.ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSetRequest{}) {
		if err := client.SecurityIntegrations.AlterApiAuthenticationWithAuthorizationCodeGrantFlow(ctx, sdk.NewAlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unset, sdk.ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterApiAuthenticationWithAuthorizationCodeGrantFlow(ctx, sdk.NewAlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(false)(ctx, d, meta)
}

func DeleteContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(id).WithIfExists(true))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting integration",
				Detail:   fmt.Sprintf("id %v err = %v", id.Name(), err),
			},
		}
	}

	d.SetId("")
	return nil
}
