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
)

var secretClientCredentialsSchema = func() map[string]*schema.Schema {
	secretClientCredentials := map[string]*schema.Schema{
		"api_authentication": {
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			Required:         true,
			Description:      "Specifies the name value of the Snowflake security integration that connects Snowflake to an external service when setting Type to OAUTH2.",
		},
		"oauth_scopes": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Required:    true,
			Description: "Specifies a list of scopes to use when making a request from the OAuth server by a role with USAGE on the integration during the OAuth client credentials flow.",
		},
	}
	return helpers.MergeMaps(secretCommonSchema, secretClientCredentials)
}()

func SecretWithClientCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextSecretWithClientCredentials,
		ReadContext:   ReadContextSecretWithClientCredentials,
		UpdateContext: UpdateContextSecretWithClientCredentials,
		DeleteContext: DeleteContextSecretWithClientCredentials,

		Schema: secretClientCredentialsSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextSecretWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	commonCreate := handleSecretCreate(d)

	apiIntegrationString := d.Get("api_authentication").(string)
	apiIntegration, err := sdk.ParseAccountObjectIdentifier(apiIntegrationString)
	if err != nil {
		return diag.FromErr(err)
	}

	id := sdk.NewSchemaObjectIdentifier(commonCreate.database, commonCreate.schema, commonCreate.name)

	stringScopes := expandStringList(d.Get("oauth_scopes").(*schema.Set).List())
	oauthScopes := make([]sdk.ApiIntegrationScope, len(stringScopes))
	for i, scope := range stringScopes {
		oauthScopes[i] = sdk.ApiIntegrationScope{Scope: scope}
	}

	request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, apiIntegration, oauthScopes)

	if commonCreate.comment != nil {
		request.WithComment(*commonCreate.comment)
	}

	err = client.Secrets.CreateWithOAuthClientCredentialsFlow(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSecretWithClientCredentials(ctx, d, meta)
}

func ReadContextSecretWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if err = d.Set("api_authentication", secretDescription.IntegrationName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("oauth_scopes", secretDescription.OauthScopes); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateContextSecretWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("oauth_scopes") {
		stringScopes := expandStringList(d.Get("oauth_scopes").(*schema.Set).List())
		oauthScopes := make([]sdk.ApiIntegrationScope, len(stringScopes))
		for i, scope := range stringScopes {
			oauthScopes[i] = sdk.ApiIntegrationScope{Scope: scope}
		}

		request := sdk.NewAlterSecretRequest(id)
		setRequest := sdk.NewSetForOAuthClientCredentialsFlowRequest(oauthScopes)
		request.WithSet(*sdk.NewSecretSetRequest().WithSetForOAuthClientCredentialsFlow(*setRequest))

		if err := client.Secrets.Alter(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}

	if request := handleSecretUpdate(id, d); request != nil {
		if err := client.Secrets.Alter(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSecretWithClientCredentials(ctx, d, meta)
}

func DeleteContextSecretWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
