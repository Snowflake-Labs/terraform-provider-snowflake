package resources

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"golang.org/x/exp/maps"
)

var accountParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Name of account parameter. Valid values are those in [account parameters](https://docs.snowflake.com/en/sql-reference/parameters.html#account-parameters).",
		ValidateFunc: validation.StringInSlice(maps.Keys(sdk.GetParameterDefaults(sdk.ParameterTypeAccount)), false),
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
	db := meta.(*sql.DB)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	parameter := sdk.AccountParameter(key)

	parameterDefault := sdk.GetParameterDefaults(sdk.ParameterTypeSession)[key]
	if parameterDefault.Validate != nil {
		if err := parameterDefault.Validate(value); err != nil {
			return err
		}
	}

	err := client.Parameters.SetAccountParameter(ctx, parameter, value)
	if err != nil {
		return err
	}
	d.SetId(key)
	return ReadAccountParameter(d, meta)
}

// ReadAccountParameter implements schema.ReadFunc.
func ReadAccountParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
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
	db := meta.(*sql.DB)
	key := d.Get("key").(string)
	client := sdk.NewClientFromDB(db)
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
