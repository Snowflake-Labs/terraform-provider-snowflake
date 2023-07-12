package resources

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
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
		ValidateFunc: validation.StringInSlice(maps.Keys(snowflake.GetParameterDefaults(snowflake.ParameterTypeAccount)), false),
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

	parameterDefault := snowflake.GetParameterDefaults(snowflake.ParameterTypeAccount)[key]
	if parameterDefault.Validate != nil {
		if err := parameterDefault.Validate(value); err != nil {
			return err
		}
	}

	// add quotes to value if it is a string
	typeString := reflect.TypeOf("")
	if reflect.TypeOf(parameterDefault.DefaultValue) == typeString {
		value = fmt.Sprintf("'%s'", snowflake.EscapeString(value))
	}

	builder := snowflake.NewAccountParameter(key, value, db)
	err := builder.SetParameter()
	if err != nil {
		return fmt.Errorf("error creating account parameter err = %w", err)
	}

	d.SetId(key)
	p, err := snowflake.ShowAccountParameter(db, key)
	if err != nil {
		return fmt.Errorf("error reading account parameter err = %w", err)
	}
	err = d.Set("value", p.Value.String)
	if err != nil {
		return fmt.Errorf("error setting account parameter value err = %w", err)
	}
	return nil
}

// ReadAccountParameter implements schema.ReadFunc.
func ReadAccountParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	key := d.Id()
	p, err := snowflake.ShowAccountParameter(db, key)
	if err != nil {
		return fmt.Errorf("error reading account parameter err = %w", err)
	}
	err = d.Set("value", p.Value.String)
	if err != nil {
		return fmt.Errorf("error setting account parameter value err = %w", err)
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

	parameterDefault := snowflake.GetParameterDefaults(snowflake.ParameterTypeAccount)[key]
	defaultValue := parameterDefault.DefaultValue
	value := fmt.Sprintf("%v", defaultValue)

	// add quotes to value if it is a string
	typeString := reflect.TypeOf("")
	if reflect.TypeOf(parameterDefault.DefaultValue) == typeString {
		value = fmt.Sprintf("'%s'", value)
	}
	builder := snowflake.NewAccountParameter(key, value, db)
	err := builder.SetParameter()
	if err != nil {
		return fmt.Errorf("error creating account parameter err = %w", err)
	}
	_, err = snowflake.ShowAccountParameter(db, key)
	if err != nil {
		return fmt.Errorf("error reading account parameter err = %w", err)
	}

	d.SetId("")
	return nil
}
