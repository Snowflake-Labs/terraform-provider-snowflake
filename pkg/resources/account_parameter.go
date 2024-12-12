package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToAccountParameter),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToAccountParameter),
		Description:      fmt.Sprintf("Name of account parameter. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AsStringList(sdk.AllAccountParameters))),
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of account parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation. The parameter values are validated in Snowflake.",
	},
}

func AccountParameter() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.AccountParameter, CreateAccountParameter),
		ReadContext:   TrackingReadWrapper(resources.AccountParameter, ReadAccountParameter),
		UpdateContext: TrackingUpdateWrapper(resources.AccountParameter, UpdateAccountParameter),
		DeleteContext: TrackingDeleteWrapper(resources.AccountParameter, DeleteAccountParameter),

		Description: "Resource used to manage current account parameters. For more information, check [parameters documentation](https://docs.snowflake.com/en/sql-reference/parameters).",

		Schema: accountParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateAccountParameter implements schema.CreateFunc.
func CreateAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	parameter, err := sdk.ToAccountParameter(key)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.Parameters.SetAccountParameter(ctx, parameter, value)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(string(parameter)))
	return ReadAccountParameter(ctx, d, meta)
}

// ReadAccountParameter implements schema.ReadFunc.
func ReadAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	parameterNameRaw := d.Id()
	parameterName, err := sdk.ToAccountParameter(parameterNameRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	parameter, err := client.Parameters.ShowAccountParameter(ctx, parameterName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading account parameter: %w", err))
	}
	errs := errors.Join(
		d.Set("value", parameter.Value),
		d.Set("key", parameter.Key),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	return nil
}

// UpdateAccountParameter implements schema.UpdateFunc.
func UpdateAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return CreateAccountParameter(ctx, d, meta)
}

// DeleteAccountParameter implements schema.DeleteFunc.
func DeleteAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	key := d.Get("key").(string)
	parameter := sdk.AccountParameter(key)

	err := client.Parameters.UnsetAccountParameter(ctx, parameter)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unsetting account parameter: %w", err))
	}

	d.SetId("")
	return nil
}
