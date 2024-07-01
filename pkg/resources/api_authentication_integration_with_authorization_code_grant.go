package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var apiAuthAuthorizationCodeGrantSchema = func() map[string]*schema.Schema {
	uniq := map[string]*schema.Schema{
		"oauth_authorization_endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "unknown",
			Description: "Specifies the URL for authenticating to the external service.",
		},
		"oauth_allowed_scopes": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "Specifies a list of scopes to use when making a request from the OAuth by a role with USAGE on the integration during the OAuth client credentials flow.",
		},
		"oauth_grant": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"unknown", "AUTHORIZATION_CODE"}, true),
			Description:  "Specifies the type of OAuth flow.",
			Default:      "unknown",
		},
	}
	return MergeMaps(apiAuthCommonSchema, uniq)
}()

func ApiAuthenticationIntegrationWithAuthorizationCodeGrant() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextApiAuthenticationIntegrationWithAuthorizationCodeGrant,
		ReadContext:   ReadContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(true),
		UpdateContext: UpdateContextApiAuthenticationIntegrationWithAuthorizationCodeGrant,
		DeleteContext: DeleteContextApiAuthenticationIntegrationWithAuthorizationCodeGrant,
		CustomizeDiff: customdiff.All(
			ForceNewIfChangeToDefaultString("oauth_token_endpoint"),
			ForceNewIfChangeToDefaultString("oauth_authorization_endpoint"),
			ForceNewIfChangeToDefaultString("oauth_client_auth_method"),
			ForceNewIfChangeToDefaultString("oauth_grant"),
			ComputedIfAnyAttributeChanged(showOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(describeOutputAttributeName, "enabled", "comment", "oauth_access_token_validity", "oauth_refresh_token_validity",
				"oauth_client_id", "oauth_client_auth_method", "oauth_authorization_endpoint",
				"oauth_token_endpoint", "oauth_allowed_scopes", "oauth_grant"),
		),
		Schema: apiAuthAuthorizationCodeGrantSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportApiAuthenticationWithAuthorizationCodeGrant,
		},
	}
}

func ImportApiAuthenticationWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting api auth integration with authorization code grant import")
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	properties, err := client.SecurityIntegrations.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = d.Set("name", integration.Name); err != nil {
		return nil, err
	}
	if err = d.Set("enabled", integration.Enabled); err != nil {
		return nil, err
	}
	if err = d.Set("comment", integration.Comment); err != nil {
		return nil, err
	}

	oauthAccessTokenValidity, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_ACCESS_TOKEN_VALIDITY"
	})
	if err == nil {
		value, err := strconv.Atoi(oauthAccessTokenValidity.Value)
		if err != nil {
			return nil, err
		}
		if err = d.Set("oauth_access_token_validity", value); err != nil {
			return nil, err
		}
	}
	oauthRefreshTokenValidity, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_REFRESH_TOKEN_VALIDITY"
	})
	if err == nil {
		value, err := strconv.Atoi(oauthRefreshTokenValidity.Value)
		if err != nil {
			return nil, err
		}
		if err = d.Set("oauth_refresh_token_validity", value); err != nil {
			return nil, err
		}
	}
	oauthClientId, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_CLIENT_ID" })
	if err == nil {
		if err = d.Set("oauth_client_id", oauthClientId.Value); err != nil {
			return nil, err
		}
	}
	oauthClientAuthMethod, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_CLIENT_AUTH_METHOD"
	})
	if err == nil {
		if err = d.Set("oauth_client_auth_method", oauthClientAuthMethod.Value); err != nil {
			return nil, err
		}
	}
	oauthAuthorizationEndpoint, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_AUTHORIZATION_ENDPOINT"
	})
	if err == nil {
		if err = d.Set("oauth_authorization_endpoint", oauthAuthorizationEndpoint.Value); err != nil {
			return nil, err
		}
	}
	oauthTokenEndpoint, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_TOKEN_ENDPOINT" })
	if err == nil {
		if err = d.Set("oauth_token_endpoint", oauthTokenEndpoint.Value); err != nil {
			return nil, err
		}
	}
	oauthAllowedScopes, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_ALLOWED_SCOPES" })
	if err == nil {
		if err = d.Set("oauth_allowed_scopes", listValueToSlice(oauthAllowedScopes.Value, true, false)); err != nil {
			return nil, err
		}
	}
	oauthGrant, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_GRANT" })
	if err == nil {
		if err = d.Set("oauth_grant", oauthGrant.Value); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	enabled := d.Get("enabled").(bool)
	name := d.Get("name").(string)
	oauthClientId := d.Get("oauth_client_id").(string)
	oauthClientSecret := d.Get("oauth_client_secret").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	req := sdk.NewCreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest(id, enabled, oauthClientId, oauthClientSecret)

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if v := d.Get("oauth_access_token_validity").(int); v != -1 {
		req.WithOauthAccessTokenValidity(v)
	}

	if v := d.Get("oauth_authorization_endpoint").(string); v != "unknown" {
		req.WithOauthAuthorizationEndpoint(v)
	}

	if v := d.Get("oauth_client_auth_method").(string); v != "unknown" {
		value, err := sdk.ToApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithOauthClientAuthMethod(value)
	}

	if v := d.Get("oauth_refresh_token_validity").(int); v != -1 {
		req.WithOauthRefreshTokenValidity(v)
	}

	if v := d.Get("oauth_grant").(string); v != "unknown" {
		if v == "AUTHORIZATION_CODE" {
			req.WithOauthGrantAuthorizationCode(true)
		}
	}

	if v := d.Get("oauth_token_endpoint").(string); v != "unknown" {
		req.WithOauthTokenEndpoint(v)
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

	d.SetId(name)

	return ReadContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(false)(ctx, d, meta)
}

func ReadContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

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
		if err := d.Set("name", integration.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("comment", integration.Comment); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("enabled", integration.Enabled); err != nil {
			return diag.FromErr(err)
		}
		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"comment", "comment", integration.Comment, integration.Comment, nil},
				showMapping{"enabled", "enabled", integration.Enabled, integration.Enabled, nil},
			); err != nil {
				return diag.FromErr(err)
			}

			enabled, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "ENABLED" })
			if err != nil {
				return diag.FromErr(err)
			}

			oauthAccessTokenValidity, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_ACCESS_TOKEN_VALIDITY"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			oauthRefreshTokenValidity, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_REFRESH_TOKEN_VALIDITY"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			oauthClientId, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_CLIENT_ID" })
			if err != nil {
				return diag.FromErr(err)
			}

			oauthClientAuthMethod, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_CLIENT_AUTH_METHOD"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			oauthAuthorizationEndpoint, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_AUTHORIZATION_ENDPOINT"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			oauthTokenEndpoint, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_TOKEN_ENDPOINT" })
			if err != nil {
				return diag.FromErr(err)
			}

			oauthAllowedScopes, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_ALLOWED_SCOPES" })
			if err != nil {
				return diag.FromErr(err)
			}

			oauthGrant, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_GRANT" })
			if err != nil {
				return diag.FromErr(err)
			}
			oauthAccessTokenValidityInt, err := strconv.Atoi(oauthAccessTokenValidity.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			oauthRefreshTokenValidityInt, err := strconv.Atoi(oauthRefreshTokenValidity.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err = handleExternalChangesToObjectInDescribe(d,
				describeMapping{"enabled", "enabled", enabled.Value, enabled.Value, nil},
				describeMapping{"oauth_access_token_validity", "oauth_access_token_validity", oauthAccessTokenValidityInt, oauthAccessTokenValidityInt, stringToIntNormalizer},
				describeMapping{"oauth_refresh_token_validity", "oauth_refresh_token_validity", oauthRefreshTokenValidityInt, oauthRefreshTokenValidityInt, stringToIntNormalizer},
				describeMapping{"oauth_client_id", "oauth_client_id", oauthClientId.Value, oauthClientId.Value, nil},
				describeMapping{"oauth_client_auth_method", "oauth_client_auth_method", oauthClientAuthMethod.Value, oauthClientAuthMethod.Value, nil},
				describeMapping{"oauth_authorization_endpoint", "oauth_authorization_endpoint", oauthAuthorizationEndpoint.Value, oauthAuthorizationEndpoint.Value, nil},
				describeMapping{"oauth_token_endpoint", "oauth_token_endpoint", oauthTokenEndpoint.Value, oauthTokenEndpoint.Value, nil},
				describeMapping{"oauth_allowed_scopes", "oauth_allowed_scopes", oauthAllowedScopes.Value, listValueToSlice(oauthAllowedScopes.Value, true, false), nil},
				describeMapping{"oauth_grant", "oauth_grant", oauthGrant.Value, oauthGrant.Value, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}
		if !d.GetRawConfig().IsNull() {
			if v := d.GetRawConfig().AsValueMap()["enabled"]; !v.IsNull() {
				if err = d.Set("enabled", v.True()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["oauth_access_token_validity"]; !v.IsNull() {
				intVal, _ := v.AsBigFloat().Int64()
				if err = d.Set("oauth_access_token_validity", intVal); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["oauth_refresh_token_validity"]; !v.IsNull() {
				intVal, _ := v.AsBigFloat().Int64()
				if err = d.Set("oauth_refresh_token_validity", intVal); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["oauth_client_id"]; !v.IsNull() {
				if err = d.Set("oauth_client_id", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["oauth_client_auth_method"]; !v.IsNull() {
				if err = d.Set("oauth_client_auth_method", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["oauth_authorization_endpoint"]; !v.IsNull() {
				if err = d.Set("oauth_authorization_endpoint", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["oauth_token_endpoint"]; !v.IsNull() {
				if err = d.Set("oauth_token_endpoint", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["oauth_allowed_scopes"]; !v.IsNull() {
				if err = d.Set("oauth_allowed_scopes", ctyValToSliceString(v)); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["oauth_grant"]; !v.IsNull() {
				if err = d.Set("oauth_grant", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["comment"]; !v.IsNull() {
				if err = d.Set("comment", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if err = d.Set(showOutputAttributeName, []map[string]any{schemas.SecurityIntegrationToSchema(integration)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(describeOutputAttributeName, []map[string]any{schemas.ApiAuthSecurityIntegrationPropertiesToSchema(properties)}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateContextApiAuthenticationIntegrationWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	set, unset := sdk.NewApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSetRequest(), sdk.NewApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationUnsetRequest()

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.WithComment(v.(string))
		} else {
			unset.WithComment(true)
		}
	}

	if d.HasChange("enabled") {
		if v := d.Get("comment").(string); v != "unknown" {
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithEnabled(parsed)
		} else {
			unset.WithEnabled(true)
		}
	}

	if d.HasChange("oauth_access_token_validity") {
		if v := d.Get("oauth_access_token_validity").(int); v != -1 {
			set.WithOauthAccessTokenValidity(v)
		} else {
			// TODO: use UNSET
			set.WithOauthAccessTokenValidity(0)
		}
	}

	if d.HasChange("oauth_authorization_endpoint") {
		set.WithOauthAuthorizationEndpoint(d.Get("oauth_authorization_endpoint").(string))
	}

	if d.HasChange("oauth_client_auth_method") {
		v := d.Get("oauth_client_auth_method").(string)
		if len(v) > 0 {
			value, err := sdk.ToApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithOauthClientAuthMethod(value)
		}
	}

	if d.HasChange("oauth_client_id") {
		set.WithOauthClientId(d.Get("oauth_client_id").(string))
	}
	if d.HasChange("oauth_grant") {
		if v := d.Get("oauth_grant").(string); v == "AUTHORIZATION_CODE" {
			set.WithOauthGrantAuthorizationCode(true)
		}
	}

	if d.HasChange("oauth_client_secret") {
		set.WithOauthClientSecret(d.Get("oauth_client_secret").(string))
	}

	if d.HasChange("oauth_refresh_token_validity") {
		if v := d.Get("oauth_refresh_token_validity").(int); v != -1 {
			set.WithOauthRefreshTokenValidity(v)
		} else {
			// TODO: use UNSET
			set.WithOauthRefreshTokenValidity(7776000)
		}
	}

	if d.HasChange("oauth_token_endpoint") {
		set.WithOauthTokenEndpoint(d.Get("oauth_token_endpoint").(string))
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
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	client := meta.(*provider.Context).Client

	err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(sdk.NewAccountObjectIdentifier(id.Name())).WithIfExists(true))
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
