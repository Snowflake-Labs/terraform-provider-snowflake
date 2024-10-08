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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
)

var secretBasicAuthenticationSchema = func() map[string]*schema.Schema {
	secretBasicAuthentication := map[string]*schema.Schema{
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the username value to store in the secret.",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: externalChangesNotDetectedFieldDescription("Specifies the password value to store in the secret."),
		},
	}
	return helpers.MergeMaps(secretCommonSchema, secretBasicAuthentication)
}()

func SecretWithBasicAuthentication() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextSecretWithBasicAuthentication,
		ReadContext:   ReadContextSecretWithBasicAuthentication,
		UpdateContext: UpdateContextSecretWithBasicAuthentication,
		DeleteContext: DeleteContextSecretWithBasicAuthentication,
		Description:   "Resource used to manage secret objects with Basic Authentication. For more information, check [secret documentation](https://docs.snowflake.com/en/sql-reference/sql/create-secret).",

		Schema: secretBasicAuthenticationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportSecretWithBasicAuthentication,
		},
	}
}

func ImportSecretWithBasicAuthentication(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting secret with basic authentication import")
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

	if err := d.Set("username", secretDescription.Username); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextSecretWithBasicAuthentication(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName, schemaName, name := handleSecretCreate(d)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

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

func ReadContextSecretWithBasicAuthentication(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	if err = handleSecretRead(d, id, secret, secretDescription); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("username", secretDescription.Username); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateContextSecretWithBasicAuthentication(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	setForFlow := &sdk.SetForFlowRequest{
		SetForBasicAuthentication: &sdk.SetForBasicAuthenticationRequest{},
	}

	if d.HasChange("username") {
		username := d.Get("username").(string)
		setForFlow.SetForBasicAuthentication.WithUsername(username)
		set.WithSetForFlow(*setForFlow)
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		setForFlow.SetForBasicAuthentication.WithPassword(password)
		set.WithSetForFlow(*setForFlow)
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

	return ReadContextSecretWithBasicAuthentication(ctx, d, meta)
}

func DeleteContextSecretWithBasicAuthentication(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
