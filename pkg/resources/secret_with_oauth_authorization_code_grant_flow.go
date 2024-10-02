package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var secretAuthorizationCodeGrantSchema = func() map[string]*schema.Schema {
	secretAuthorizationCodeGrant := map[string]*schema.Schema{
		"oauth_refresh_token": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the token as a string that is used to obtain a new access token from the OAuth authorization server when the access token expires.",
		},
		"oauth_refresh_token_expiry_time": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("oauth_refresh_token_expiry_time"),
			Description:      "Specifies the timestamp as a string when the OAuth refresh token expires. Accepted string formats: YYYY-MM-DD, YYYY-MM-DD HH:MI, YYYY-MM-DD HH:MI:SS, YYYY-MM-DD HH:MI PDT",
		},
		"api_authentication": {
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			Required:         true,
			Description:      "Specifies the name value of the Snowflake security integration that connects Snowflake to an external service when setting Type to OAUTH2.",
		},
	}
	return helpers.MergeMaps(secretCommonSchema, secretAuthorizationCodeGrant)
}()

func SecretWithAuthorizationCodeGrant() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextSecretWithAuthorizationCodeGrant,
		ReadContext:   ReadContextSecretWithAuthorizationCodeGrant,
		UpdateContext: UpdateContextSecretWithAuthorizationCodeGrant,
		DeleteContext: DeleteContextSecretWithAuthorizationCodeGrant,

		Schema: secretAuthorizationCodeGrantSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportSecretWithAuthorizationCodeGrant,
		},
		Description: "Secret with OAuth authorization code grant where Secret's Type attribute is set to 'OAUTH2'.'",

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(secretAuthorizationCodeGrantSchema, DescribeOutputAttributeName, "oauth_refresh_token_expiry_time"),
		),
	}
}

func ImportSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting secret with authorization code import")
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	secretDescription, err := client.Secrets.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	// cannot import oauth_refresh_token because it is not present both in SHOW or DESCRIBE

	if err := d.Set("oauth_refresh_token_expiry_time", secretDescription.OauthRefreshTokenExpiryTime.String()); err != nil {
		return nil, err
	}
	if err := d.Set("api_authentication", secretDescription.IntegrationName); err != nil {
		return nil, err
	}

	secret, err := client.Secrets.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := handleSecretRead(d, id, secret, secretDescription); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateContextSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return ReadContextSecretWithAuthorizationCodeGrant(ctx, d, meta)
}

func ReadContextSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if err = setStateToValuesFromConfig(d, secretAuthorizationCodeGrantSchema, []string{"oauth_refresh_token_expiry_time"}); err != nil {
		return diag.FromErr(err)
	}
	/*
		// Possible limitation
		// Accepted formats are: YYYY-MM-DD; YYYY-MM-DD HH:MI:SS
		// But snowflake holds this value as timestamp, so with this code below we can parse it and keep in state only with one of time.DateOnly or time.DateTime
		if err = d.Set("oauth_refresh_token_expiry_time", secretDescription.OauthRefreshTokenExpiryTime.In(time.UTC).Format(time.DateOnly)); err != nil {
			return diag.FromErr(err)
		}
	*/
	if err = d.Set("api_authentication", secretDescription.IntegrationName); err != nil {
		return diag.FromErr(err)
	}
	if err := handleSecretRead(d, id, secret, secretDescription); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateContextSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return ReadContextSecretWithAuthorizationCodeGrant(ctx, d, meta)
}

func DeleteContextSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
