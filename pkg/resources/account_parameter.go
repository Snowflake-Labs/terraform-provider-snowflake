package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Name of account parameter. Valid values are those in [account parameters](https://docs.snowflake.com/en/sql-reference/parameters.html#account-parameters).",
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of account parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation.",
	},
}

func AccountParameter() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.AccountParameter, CreateAccountParameter),
		ReadContext:   TrackingReadWrapper(resources.AccountParameter, ReadAccountParameter),
		UpdateContext: TrackingUpdateWrapper(resources.AccountParameter, UpdateAccountParameter),
		DeleteContext: TrackingDeleteWrapper(resources.AccountParameter, DeleteAccountParameter),

		Schema: accountParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateAccountParameter implements schema.CreateFunc.
func CreateAccountParameter(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	parameter := sdk.AccountParameter(key)
	err := client.Parameters.SetAccountParameter(ctx, parameter, value)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(key)
	return ReadAccountParameter(ctx, d, meta)
}

// ReadAccountParameter implements schema.ReadFunc.
func ReadAccountParameter(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	parameterName := d.Id()
	parameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(parameterName))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading account parameter err = %w", err))
	}
	err = d.Set("value", parameter.Value)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting account parameter err = %w", err))
	}
	err = d.Set("key", parameter.Key)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting account parameter err = %w", err))
	}
	return nil
}

// UpdateAccountParameter implements schema.UpdateFunc.
func UpdateAccountParameter(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return CreateAccountParameter(ctx, d, meta)
}

// DeleteAccountParameter implements schema.DeleteFunc.
func DeleteAccountParameter(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	key := d.Get("key").(string)
	parameter := sdk.AccountParameter(key)
	defaultParameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(key))
	if err != nil {
		return diag.FromErr(err)
	}
	defaultValue := defaultParameter.Default
	err = client.Parameters.SetAccountParameter(ctx, parameter, defaultValue)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error resetting account parameter err = %w", err))
	}

	d.SetId("")
	return nil
}
