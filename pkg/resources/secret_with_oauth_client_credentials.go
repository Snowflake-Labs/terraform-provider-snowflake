package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"reflect"

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
			Description:      "Specifies the name value of the Snowflake security integration that connects Snowflake to an external service.",
			DiffSuppressFunc: suppressIdentifierQuoting,
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
		Description:   "Resource used to manage secret objects with OAuth Client Credentials. For more information, check [secret documentation](https://docs.snowflake.com/en/sql-reference/sql/create-secret).",

		Schema: secretClientCredentialsSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportSecretWithClientCredentials,
		},
	}
}

func ImportSecretWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting secret with client credentials import")
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	if err := handleSecretImport(d); err != nil {
		return nil, err
	}
	secretDescription, err := client.Secrets.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := d.Set("api_authentication", secretDescription.IntegrationName); err != nil {
		return nil, err
	}

	if err := d.Set("oauth_scopes", secretDescription.OauthScopes); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
func CreateContextSecretWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName, schemaName, name := handleSecretCreate(d)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	apiIntegrationString := d.Get("api_authentication").(string)
	apiIntegration, err := sdk.ParseAccountObjectIdentifier(apiIntegrationString)
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, apiIntegration)

	stringScopes := expandStringList(d.Get("oauth_scopes").(*schema.Set).List())
	oauthScopes := make([]sdk.ApiIntegrationScope, len(stringScopes))
	for i, scope := range stringScopes {
		oauthScopes[i] = sdk.ApiIntegrationScope{Scope: scope}
	}
	request.WithOauthScopes(sdk.OauthScopesListRequest{OauthScopesList: oauthScopes})

	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	err = client.Secrets.CreateWithOAuthClientCredentialsFlow(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSecretWithClientCredentials(ctx, d, meta)
}

func ReadContextSecretWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
					Summary:  "Failed to retrieve secret with client credentials. Target object not found. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Secret with client credentials name: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve secret with client credentials.",
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

func UpdateContextSecretWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	commonSet, commonUnset := handleSecretUpdate(d)
	set := &sdk.SecretSetRequest{
		Comment: commonSet.comment,
	}
	unset := &sdk.SecretUnsetRequest{
		Comment: commonUnset.comment,
	}

	if d.HasChange("oauth_scopes") {
		stringScopes := expandStringList(d.Get("oauth_scopes").(*schema.Set).List())
		oauthScopes := make([]sdk.ApiIntegrationScope, len(stringScopes))
		for i, scope := range stringScopes {
			oauthScopes[i] = sdk.ApiIntegrationScope{Scope: scope}
		}
		req := sdk.NewSetForOAuthClientCredentialsRequest().WithOauthScopes(*sdk.NewOauthScopesListRequest(oauthScopes))
		set.WithSetForFlow(*sdk.NewSetForFlowRequest().WithSetForOAuthClientCredentials(*req))
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

	return ReadContextSecretWithClientCredentials(ctx, d, meta)
}

func DeleteContextSecretWithClientCredentials(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
