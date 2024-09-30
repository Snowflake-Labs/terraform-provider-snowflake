package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

var secretAuthorizationCodeSchema = func() map[string]*schema.Schema {
	secretAuthorizationCode := map[string]*schema.Schema{
		"oauth_refresh_token": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the token as a string that is used to obtain a new access token from the OAuth authorization server when the access token expires.",
		},
		"oauth_refresh_token_expiry_time": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the timestamp as a string when the OAuth refresh token expires. Accepted formats: YYYY-MM-DD, YYYY-MM-DD HH:MI:SS",
		},
		"api_authentication": {
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			Required:         true,
			Description:      "Specifies the name value of the Snowflake security integration that connects Snowflake to an external service when setting Type to OAUTH2.",
		},
	}
	return helpers.MergeMaps(secretCommonSchema, secretAuthorizationCode)
}()

func SecretWithAuthorizationCode() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextSecretWithAuthorizationCode,
		ReadContext:   ReadContextSecretWithAuthorizationCode,
		UpdateContext: UpdateContextSecretWithAuthorizationCode,
		DeleteContext: DeleteContextSecretWithAuthorizationCode,

		Schema: secretAuthorizationCodeSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextSecretWithAuthorizationCode(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	commonCreate := handleSecretCreate(d)

	id := sdk.NewSchemaObjectIdentifier(commonCreate.database, commonCreate.schema, commonCreate.name)

	apiIntegrationString := d.Get("api_authentication").(string)
	apiIntegration, err := sdk.ParseAccountObjectIdentifier(apiIntegrationString)
	if err != nil {
		return diag.FromErr(err)
	}

	refreshToken := d.Get("oauth_refresh_token").(string)
	refreshTokenExpiryTime := d.Get("oauth_refresh_token_expiry_time").(string)

	request := sdk.NewCreateWithOAuthAuthorizationCodeFlowSecretRequest(id, refreshToken, refreshTokenExpiryTime, apiIntegration)
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	err = client.Secrets.CreateWithOAuthAuthorizationCodeFlow(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSecretWithAuthorizationCode(ctx, d, meta)
}

func ReadContextSecretWithAuthorizationCode(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := client.Secrets.ShowByID(ctx, id)
	if secret == nil || err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve secret. Target object not found. Marking the resource as removed.",
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve secret.",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}
	secretDescription, err := client.Secrets.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := handleSecretRead(d, id, secret, secretDescription); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("oauth_refresh_token_expiry_time", secretDescription.OauthRefreshTokenExpiryTime.In(time.UTC).Format(time.DateOnly)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("api_authentication", secretDescription.IntegrationName); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateContextSecretWithAuthorizationCode(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if request := handleSecretUpdate(id, d); request != nil {
		if err := client.Secrets.Alter(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("oauth_refresh_token") {
		refreshToken := d.Get("oauth_refresh_token").(string)

		request := sdk.NewAlterSecretRequest(id)
		setRequest := sdk.NewSetForOAuthAuthorizationFlowRequest().WithOauthRefreshToken(refreshToken)
		request.WithSet(*sdk.NewSecretSetRequest().WithSetForOAuthAuthorizationFlow(*setRequest))

		if err := client.Secrets.Alter(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("oauth_refresh_token_expiry_time") {
		refreshTokenExpiryTime := d.Get("oauth_refresh_token_expiry_time").(string)

		request := sdk.NewAlterSecretRequest(id)
		setRequest := sdk.NewSetForOAuthAuthorizationFlowRequest().WithOauthRefreshTokenExpiryTime(refreshTokenExpiryTime)
		request.WithSet(*sdk.NewSecretSetRequest().WithSetForOAuthAuthorizationFlow(*setRequest))

		if err := client.Secrets.Alter(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSecretWithAuthorizationCode(ctx, d, meta)
}

func DeleteContextSecretWithAuthorizationCode(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Secrets.Drop(ctx, sdk.NewDropSecretRequest(id).WithIfExists(*sdk.Bool(true))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
