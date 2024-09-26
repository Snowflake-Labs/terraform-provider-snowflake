package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var secretClientCredentialsSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("String that specifies the identifier (i.e. name) for the secret, must be unique in your schema."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the secret"),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the secret."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
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
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the secret.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func SecretWithClientCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextSecret,
		ReadContext:   ReadContextSecret,
		UpdateContext: UpdateContextSecret,
		DeleteContext: DeleteContextSecret,

		Schema: secretClientCredentialsSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextSecret(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	apiIntegrationString := d.Get("api_authentication").(string)
	apiIntegration, err := sdk.ParseAccountObjectIdentifier(apiIntegrationString)
	if err != nil {
		return diag.FromErr(err)
	}

	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	stringScopes := expandStringList(d.Get("oauth_scopes").(*schema.Set).List())
	oauthScopes := make([]sdk.ApiIntegrationScope, len(stringScopes))
	for i, scope := range stringScopes {
		oauthScopes[i] = sdk.ApiIntegrationScope{Scope: scope}
	}

	request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, apiIntegration, oauthScopes)

	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	err = client.Secrets.CreateWithOAuthClientCredentialsFlow(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSecret(ctx, d, meta)
}

func ReadContextSecret(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if err = d.Set("name", secretDescription.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("database", secretDescription.DatabaseName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("schema", secretDescription.SchemaName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("api_authentication", secretDescription.IntegrationName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("oauth_scopes", secretDescription.OauthScopes); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("comment", secretDescription.Comment); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateContextSecret(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		request := sdk.NewAlterSecretRequest(id)
		if len(comment) == 0 {
			unsetRequest := sdk.NewSecretUnsetRequest().WithComment(*sdk.Bool(true))
			request.WithUnset(*unsetRequest)
		} else {
			setRequest := sdk.NewSecretSetRequest().WithComment(comment)
			request.WithSet(*setRequest)
		}
		if err := client.Secrets.Alter(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSecret(ctx, d, meta)
}

func DeleteContextSecret(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
