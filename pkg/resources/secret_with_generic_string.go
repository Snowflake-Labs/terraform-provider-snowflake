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

var secretGenericStringSchema = func() map[string]*schema.Schema {
	secretGenericString := map[string]*schema.Schema{
		"secret_string": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Specifies the string to store in the secret. The string can be an API token or a string of sensitive value that can be used in the handler code of a UDF or stored procedure. For details, see [Creating and using an external access integration](https://docs.snowflake.com/en/developer-guide/external-network-access/creating-using-external-network-access). You should not use this property to store any kind of OAuth token; use one of the other secret types for your OAuth use cases.",
		},
	}
	return helpers.MergeMaps(secretCommonSchema, secretGenericString)
}()

func SecretWithGenericString() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextSecretWithGenericString,
		ReadContext:   ReadContextSecretWithGenericString,
		UpdateContext: UpdateContextSecretWithGenericString,
		DeleteContext: DeleteContextSecretWithGenericString,
		Description:   "Resource used to manage secret objects with Generic String. For more information, check [secret documentation](https://docs.snowflake.com/en/sql-reference/sql/create-secret).",

		Schema: secretGenericStringSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportSecretWithGenericString,
		},
	}
}

func ImportSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting secret with generic string import")

	_, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	if err := handleSecretImport(d); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := sdk.NewSchemaObjectIdentifier(handleSecretCreate(d))

	secretSting := d.Get("secret_string").(string)

	request := sdk.NewCreateWithGenericStringSecretRequest(id, secretSting)
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	err := client.Secrets.CreateWithGenericString(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSecretWithGenericString(ctx, d, meta)
}

func ReadContextSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	if err := handleSecretRead(d, id, secret, secretDescription); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateContextSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	if d.HasChange("secret_string") {
		secretString := d.Get("secret_string").(string)
		req := sdk.NewSetForFlowRequest().WithSetForGenericString(*sdk.NewSetForGenericStringRequest().WithSecretString(secretString))
		set.WithSetForFlow(*req)
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

	return ReadContextSecretWithGenericString(ctx, d, meta)
}

func DeleteContextSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
