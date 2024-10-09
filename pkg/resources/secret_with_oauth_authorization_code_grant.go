package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

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
			Sensitive:   true,
			Description: externalChangesNotDetectedFieldDescription("Specifies the token as a string that is used to obtain a new access token from the OAuth authorization server when the access token expires."),
		},
		"oauth_refresh_token_expiry_time": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("oauth_refresh_token_expiry_time"),
			Description:      "Specifies the timestamp as a string when the OAuth refresh token expires. Accepted string formats: YYYY-MM-DD, YYYY-MM-DD HH:MI, YYYY-MM-DD HH:MI:SS, YYYY-MM-DD HH:MI <timezone>",
		},
		"api_authentication": {
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			Required:         true,
			Description:      "Specifies the name value of the Snowflake security integration that connects Snowflake to an external service.",
			DiffSuppressFunc: suppressIdentifierQuoting,
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
		Description:   "Resource used to manage secret objects with OAuth Authorization Code Grant. For more information, check [secret documentation](https://docs.snowflake.com/en/sql-reference/sql/create-secret).",

		Schema: secretAuthorizationCodeGrantSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportSecretWithAuthorizationCodeGrant,
		},

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(secretAuthorizationCodeGrantSchema, DescribeOutputAttributeName, "oauth_refresh_token_expiry_time", "api_authentication"),
		),
	}
}

func ImportSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting secret with authorization code import")
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err := handleSecretImport(d); err != nil {
		return nil, err
	}

	secretDescription, err := client.Secrets.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = d.Set("oauth_refresh_token_expiry_time", secretDescription.OauthRefreshTokenExpiryTime.String()); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName, schemaName, name := handleSecretCreate(d)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

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

func ReadContextSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := client.Secrets.ShowByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve secret with authorization code grant. Target object not found. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Secret with authorization code grant name: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve secret with authorization code grant.",
				Detail:   fmt.Sprintf("Secret with authorization code grant name: %s, Err: %s", id.FullyQualifiedName(), err),
			},
		}
	}
	secretDescription, err := client.Secrets.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("api_authentication", secretDescription.IntegrationName); err != nil {
		return diag.FromErr(err)
	}

	if err = setStateToValuesFromConfig(d, secretAuthorizationCodeGrantSchema, []string{"oauth_refresh_token_expiry_time"}); err != nil {
		return diag.FromErr(err)
	}

	if err := handleSecretRead(d, id, secret, secretDescription); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateContextSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	commonSet, commonUnset := handleSecretUpdate(d)
	set := &sdk.SecretSetRequest{
		Comment: commonSet.comment,
		SetForFlow: &sdk.SetForFlowRequest{
			SetForOAuthAuthorization: &sdk.SetForOAuthAuthorizationRequest{},
		},
	}

	unset := &sdk.SecretUnsetRequest{
		Comment: commonUnset.comment,
	}

	if d.HasChange("oauth_refresh_token") {
		refreshToken := d.Get("oauth_refresh_token").(string)
		set.SetForFlow.SetForOAuthAuthorization.WithOauthRefreshToken(refreshToken)
	}

	if d.HasChange("oauth_refresh_token_expiry_time") {
		refreshTokenExpiryTime := d.Get("oauth_refresh_token_expiry_time").(string)
		set.SetForFlow.SetForOAuthAuthorization.WithOauthRefreshTokenExpiryTime(refreshTokenExpiryTime)
	}

	if !reflect.DeepEqual(*set, sdk.SecretSetRequest{}) {
		if err := client.Secrets.Alter(ctx, sdk.NewAlterSecretRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if !reflect.DeepEqual(*unset, sdk.SecretUnsetRequest{}) {
		if err := client.Secrets.Alter(ctx, sdk.NewAlterSecretRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSecretWithAuthorizationCodeGrant(ctx, d, meta)
}

func DeleteContextSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Secrets.Drop(ctx, sdk.NewDropSecretRequest(id).WithIfExists(true)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
