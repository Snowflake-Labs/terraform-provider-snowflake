package resources

import (
	"context"
	"fmt"

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
		Create: CreateAccountParameter,
		Read:   ReadAccountParameter,
		Update: UpdateAccountParameter,
		Delete: DeleteAccountParameter,

		Schema: accountParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateAccountParameter implements schema.CreateFunc.
func CreateAccountParameter(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	ctx := context.Background()
	parameter := sdk.AccountParameter(key)
	err := client.Parameters.SetAccountParameter(ctx, parameter, value)
	if err != nil {
		return err
	}
	d.SetId(key)
	return ReadAccountParameter(d, meta)
}

// ReadAccountParameter implements schema.ReadFunc.
func ReadAccountParameter(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	parameterName := d.Id()
	parameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(parameterName))
	if err != nil {
		return fmt.Errorf("error reading account parameter err = %w", err)
	}
	err = d.Set("value", parameter.Value)
	if err != nil {
		return fmt.Errorf("error setting account parameter err = %w", err)
	}
	return nil
}

// UpdateAccountParameter implements schema.UpdateFunc.
func UpdateAccountParameter(d *schema.ResourceData, meta interface{}) error {
	return CreateAccountParameter(d, meta)
}

// DeleteAccountParameter implements schema.DeleteFunc.
func DeleteAccountParameter(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	key := d.Get("key").(string)
	ctx := context.Background()
	parameter := sdk.AccountParameter(key)
	defaultParameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(key))
	if err != nil {
		return err
	}
	defaultValue := defaultParameter.Default
	err = client.Parameters.SetAccountParameter(ctx, parameter, defaultValue)
	if err != nil {
		return fmt.Errorf("error resetting account parameter err = %w", err)
	}

	d.SetId("")
	return nil
}
