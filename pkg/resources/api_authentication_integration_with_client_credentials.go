package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var apiAuthClientCredentialsSchema = func() map[string]*schema.Schema {
	uniq := map[string]*schema.Schema{
		"oauth_allowed_scopes": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "Specifies a list of scopes to use when making a request from the OAuth by a role with USAGE on the integration during the OAuth client credentials flow.",
		},
	}
	return MergeMaps(apiAuthCommonSchema, uniq)
}()

func ApiAuthenticationIntegrationWithClientCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextApiAuthenticationIntegrationWithClientCredentials,
		ReadContext:   ReadContextApiAuthenticationIntegrationWithClientCredentials,
		UpdateContext: UpdateContextApiAuthenticationIntegrationWithClientCredentials,
		DeleteContext: DeleteContextApiAuthenticationIntegrationWithClientCredentials,
		Schema:        apiAuthClientCredentialsSchema,
		CustomizeDiff: customdiff.All(
			ForceNewIfChangeToEmptyString("oauth_token_endpoint"),
			ForceNewIfChangeToEmptyString("oauth_client_auth_method"),
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextApiAuthenticationIntegrationWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	enabled := d.Get("enabled").(bool)
	name := d.Get("name").(string)
	oauthClientId := d.Get("oauth_client_id").(string)
	oauthClientSecret := d.Get("oauth_client_secret").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	req := sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(id, enabled, oauthClientId, oauthClientSecret)

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if v, ok := d.GetOk("oauth_access_token_validity"); ok {
		req.WithOauthAccessTokenValidity(v.(int))
	}

	if v, ok := d.GetOk("oauth_allowed_scopes"); ok {
		elems := expandStringList(v.(*schema.Set).List())
		allowedScopes := make([]sdk.AllowedScope, len(elems))
		for i := range elems {
			allowedScopes[i] = sdk.AllowedScope{Scope: elems[i]}
		}
		req.WithOauthAllowedScopes(allowedScopes)
	}

	if v, ok := d.GetOk("oauth_client_auth_method"); ok {
		valueRaw := v.(string)
		value, err := sdk.ToApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption(valueRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithOauthClientAuthMethod(value)
	}

	if v, ok := d.GetOk("oauth_token_endpoint"); ok {
		req.WithOauthTokenEndpoint(v.(string))
	}

	if err := client.SecurityIntegrations.CreateApiAuthenticationWithClientCredentialsFlow(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return ReadContextApiAuthenticationIntegrationWithClientCredentials(ctx, d, meta)
}

func ReadContextApiAuthenticationIntegrationWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if c := integration.Category; c != sdk.SecurityIntegrationCategory {
		return diag.FromErr(fmt.Errorf("expected %v to be a %s integration, got %v", id, sdk.SecurityIntegrationCategory, c))
	}

	if err := d.Set("name", integration.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", integration.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_on", integration.CreatedOn.String()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", integration.Enabled); err != nil {
		return diag.FromErr(err)
	}

	properties, err := client.SecurityIntegrations.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, property := range properties {
		name := property.Name
		value := property.Value
		switch name {
		case "COMMENT":
			if err := d.Set("comment", value); err != nil {
				return diag.FromErr(err)
			}
		case "CREATED_ON":
			if err := d.Set("created_on", value); err != nil {
				return diag.FromErr(err)
			}
		case "ENABLED":
			if err := d.Set("enabled", helpers.StringToBool(value)); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_ACCESS_TOKEN_VALIDITY":
			valueInt, err := strconv.Atoi(value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("oauth_access_token_validity", valueInt); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_ALLOWED_SCOPES":
			value = strings.TrimLeft(value, "[")
			value = strings.TrimRight(value, "]")
			elems := strings.Split(value, ",")
			if value == "" {
				elems = nil
			}
			if err := d.Set("oauth_allowed_scopes", elems); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_CLIENT_AUTH_METHOD":
			if err := d.Set("oauth_client_auth_method", value); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_CLIENT_ID":
			if err := d.Set("oauth_client_id", value); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_CLIENT_SECRET":
			if err := d.Set("oauth_client_secret", value); err != nil {
				return diag.FromErr(err)
			}
		case "OAUTH_TOKEN_ENDPOINT":
			if err := d.Set("oauth_token_endpoint", value); err != nil {
				return diag.FromErr(err)
			}
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", name)
		}
	}

	return nil
}

func UpdateContextApiAuthenticationIntegrationWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	set, unset := sdk.NewApiAuthenticationWithClientCredentialsFlowIntegrationSetRequest(), sdk.NewApiAuthenticationWithClientCredentialsFlowIntegrationUnsetRequest()

	if d.HasChange("comment") {
		set.WithComment(d.Get("comment").(string))
	}

	if d.HasChange("enabled") {
		set.WithEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("oauth_access_token_validity") {
		set.WithOauthAccessTokenValidity(d.Get("oauth_access_token_validity").(int))
	}

	if d.HasChange("oauth_allowed_scopes") {
		elems := expandStringList(d.Get("oauth_allowed_scopes").(*schema.Set).List())
		allowedScopes := make([]sdk.AllowedScope, len(elems))
		for i := range elems {
			allowedScopes[i] = sdk.AllowedScope{Scope: elems[i]}
		}
		set.WithOauthAllowedScopes(allowedScopes)
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

	if d.HasChange("oauth_client_secret") {
		set.WithOauthClientSecret(d.Get("oauth_client_secret").(string))
	}

	if d.HasChange("oauth_token_endpoint") {
		set.WithOauthTokenEndpoint(d.Get("oauth_token_endpoint").(string))
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
	return ReadContextApiAuthenticationIntegrationWithClientCredentials(ctx, d, meta)
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
