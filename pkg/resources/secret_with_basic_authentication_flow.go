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

var secretBasicAuthenticationSchema = func() map[string]*schema.Schema {
	secretAuthorizationCode := map[string]*schema.Schema{
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the username value to store in the secret when setting the TYPE value to PASSWORD.",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the password value to store in the secret when setting the TYPE value to PASSWORD.",
		},
	}
	return helpers.MergeMaps(secretCommonSchema, secretAuthorizationCode)
}()

func SecretWithBasicAuthentication() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextSecretWithBasicAuthentication,
		ReadContext:   ReadContextSecretWithBasicAuthentication,
		UpdateContext: UpdateContextSecretWithBasicAuthentication,
		DeleteContext: DeleteContextSecretWithBasicAuthentication,

		Schema: secretBasicAuthenticationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextSecretWithBasicAuthentication(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	commonCreate := handleSecretCreate(d)

	id := sdk.NewSchemaObjectIdentifier(commonCreate.database, commonCreate.schema, commonCreate.name)

	usernameString := d.Get("username").(string)
	passwordString := d.Get("password").(string)

	request := sdk.NewCreateWithBasicAuthenticationSecretRequest(id, usernameString, passwordString)
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	err := client.Secrets.CreateWithBasicAuthentication(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSecretWithBasicAuthentication(ctx, d, meta)
}

func ReadContextSecretWithBasicAuthentication(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if err = d.Set("username", secretDescription.Username); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateContextSecretWithBasicAuthentication(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if d.HasChange("username") {
		username := d.Get("username").(string)

		request := sdk.NewAlterSecretRequest(id)
		setRequest := sdk.NewSetForBasicAuthenticationRequest().WithUsername(username)
		request.WithSet(*sdk.NewSecretSetRequest().WithSetForBasicAuthentication(*setRequest))

		if err := client.Secrets.Alter(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)

		request := sdk.NewAlterSecretRequest(id)
		setRequest := sdk.NewSetForBasicAuthenticationRequest().WithPassword(password)
		request.WithSet(*sdk.NewSecretSetRequest().WithSetForBasicAuthentication(*setRequest))

		if err := client.Secrets.Alter(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSecretWithBasicAuthentication(ctx, d, meta)
}

func DeleteContextSecretWithBasicAuthentication(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
