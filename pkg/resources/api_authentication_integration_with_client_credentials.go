package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var apiAuthClientCredentialsSchema = func() map[string]*schema.Schema {
	apiAuthClientCredentials := map[string]*schema.Schema{
		"oauth_allowed_scopes": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "Specifies a list of scopes to use when making a request from the OAuth by a role with USAGE on the integration during the OAuth client credentials flow.",
		},
		"oauth_grant": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{sdk.ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials}, true),
			Description:  "Specifies the type of OAuth flow.",
		},
	}
	return MergeMaps(apiAuthCommonSchema, apiAuthClientCredentials)
}()

func ApiAuthenticationIntegrationWithClientCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextApiAuthenticationIntegrationWithClientCredentials,
		ReadContext:   ReadContextApiAuthenticationIntegrationWithClientCredentials(true),
		UpdateContext: UpdateContextApiAuthenticationIntegrationWithClientCredentials,
		DeleteContext: DeleteContextApiAuthenticationIntegrationWithClientCredentials,
		Schema:        apiAuthClientCredentialsSchema,
		CustomizeDiff: customdiff.All(
			ForceNewIfChangeToEmptyString("oauth_token_endpoint"),
			ForceNewIfChangeToEmptyString("oauth_client_auth_method"),
			ForceNewIfChangeToEmptyString("oauth_grant"),
			ComputedIfAnyAttributeChanged(ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(DescribeOutputAttributeName, "enabled", "comment", "oauth_access_token_validity", "oauth_refresh_token_validity",
				"oauth_client_id", "oauth_client_auth_method", "oauth_token_endpoint", "oauth_allowed_scopes", "oauth_grant"),
		),
		Importer: &schema.ResourceImporter{
			StateContext: ImportApiAuthenticationWithClientCredentials,
		},
	}
}

func ImportApiAuthenticationWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting api auth integration with client credentials import")
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
	if err := handleApiAuthImport(ctx, d, integration, properties); err != nil {
		return nil, err
	}
	oauthAllowedScopes, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_ALLOWED_SCOPES" })
	if err == nil {
		if err = d.Set("oauth_allowed_scopes", listValueToSlice(oauthAllowedScopes.Value, false)); err != nil {
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

func CreateContextApiAuthenticationIntegrationWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	commonCreate, err := handleApiAuthCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}
	id := sdk.NewAccountObjectIdentifier(commonCreate.name)
	req := sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(id, commonCreate.enabled, commonCreate.oauthClientId, commonCreate.oauthClientSecret)
	req.Comment = commonCreate.comment
	req.OauthAccessTokenValidity = commonCreate.oauthAccessTokenValidity
	req.OauthRefreshTokenValidity = commonCreate.oauthRefreshTokenValidity
	req.OauthTokenEndpoint = commonCreate.oauthTokenEndpoint
	req.OauthClientAuthMethod = commonCreate.oauthClientAuthMethod

	if v, ok := d.GetOk("oauth_grant"); ok {
		if v.(string) == sdk.ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials {
			req.WithOauthGrantClientCredentials(true)
		}
	}

	if v, ok := d.GetOk("oauth_allowed_scopes"); ok {
		elems := expandStringList(v.(*schema.Set).List())
		allowedScopes := make([]sdk.AllowedScope, len(elems))
		for i := range elems {
			allowedScopes[i] = sdk.AllowedScope{Scope: elems[i]}
		}
		req.WithOauthAllowedScopes(allowedScopes)
	}

	if err := client.SecurityIntegrations.CreateApiAuthenticationWithClientCredentialsFlow(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))
	return ReadContextApiAuthenticationIntegrationWithClientCredentials(false)(ctx, d, meta)
}

func ReadContextApiAuthenticationIntegrationWithClientCredentials(withExternalChangesMarking bool) schema.ReadContextFunc {
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
		oauthAllowedScopes, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_ALLOWED_SCOPES" })
		if err != nil {
			return diag.FromErr(err)
		}

		oauthGrant, err := collections.FindOne(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_GRANT" })
		if err != nil {
			return diag.FromErr(err)
		}

		if err := handleApiAuthRead(d, integration, properties, withExternalChangesMarking, []describeMapping{
			{"oauth_allowed_scopes", "oauth_allowed_scopes", oauthAllowedScopes.Value, listValueToSlice(oauthAllowedScopes.Value, false), nil},
			{"oauth_grant", "oauth_grant", oauthGrant.Value, oauthGrant.Value, nil},
		}); err != nil {
			return diag.FromErr(err)
		}
		if !d.GetRawConfig().IsNull() {
			if v := d.GetRawConfig().AsValueMap()["oauth_allowed_scopes"]; !v.IsNull() {
				if err := d.Set("oauth_allowed_scopes", ctyValToSliceString(v.AsValueSlice())); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["oauth_grant"]; !v.IsNull() {
				if err := d.Set("oauth_grant", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		return nil
	}
}

func UpdateContextApiAuthenticationIntegrationWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	commonSet, commonUnset, err := handleApiAuthUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}
	set := &sdk.ApiAuthenticationWithClientCredentialsFlowIntegrationSetRequest{
		Enabled:                   commonSet.enabled,
		OauthTokenEndpoint:        commonSet.oauthTokenEndpoint,
		OauthClientAuthMethod:     commonSet.oauthClientAuthMethod,
		OauthClientId:             commonSet.oauthClientId,
		OauthClientSecret:         commonSet.oauthClientSecret,
		OauthAccessTokenValidity:  commonSet.oauthAccessTokenValidity,
		OauthRefreshTokenValidity: commonSet.oauthRefreshTokenValidity,
		Comment:                   commonSet.comment,
	}
	unset := &sdk.ApiAuthenticationWithClientCredentialsFlowIntegrationUnsetRequest{
		Comment: commonUnset.comment,
	}
	if d.HasChange("oauth_grant") {
		if v := d.Get("oauth_grant").(string); v == sdk.ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials {
			set.WithOauthGrantClientCredentials(true)
		}
		// else: force new
	}

	if d.HasChange("oauth_allowed_scopes") {
		elems := expandStringList(d.Get("oauth_allowed_scopes").(*schema.Set).List())
		allowedScopes := make([]sdk.AllowedScope, len(elems))
		for i := range elems {
			allowedScopes[i] = sdk.AllowedScope{Scope: elems[i]}
		}
		set.WithOauthAllowedScopes(allowedScopes)
	}

	if !reflect.DeepEqual(*set, sdk.ApiAuthenticationWithClientCredentialsFlowIntegrationSetRequest{}) {
		if err := client.SecurityIntegrations.AlterApiAuthenticationWithClientCredentialsFlow(ctx, sdk.NewAlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unset, sdk.ApiAuthenticationWithClientCredentialsFlowIntegrationUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterApiAuthenticationWithClientCredentialsFlow(ctx, sdk.NewAlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextApiAuthenticationIntegrationWithClientCredentials(false)(ctx, d, meta)
}

func DeleteContextApiAuthenticationIntegrationWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
