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

var secretGenericStringSchema = func() map[string]*schema.Schema {
	secretGenericString := map[string]*schema.Schema{
		"secret_string": {
			Type:        schema.TypeString,
			Required:    true,
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

		Schema: secretGenericStringSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Secret with Generic string where Secrets Type attribute is set to GENERIC_STRING.",
	}
}

func CreateContextSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	commonCreate := handleSecretCreate(d)

	id := sdk.NewSchemaObjectIdentifier(commonCreate.database, commonCreate.schema, commonCreate.name)

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

func ReadContextSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	return nil
}

func UpdateContextSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if d.HasChange("secret_string") {
		secretString := d.Get("secret_string").(string)

		request := sdk.NewAlterSecretRequest(id)
		setRequest := sdk.NewSetForGenericStringRequest().WithSecretString(secretString)
		request.WithSet(*sdk.NewSecretSetRequest().WithSetForGenericString(*setRequest))

		if err := client.Secrets.Alter(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSecretWithGenericString(ctx, d, meta)
}

func DeleteContextSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
