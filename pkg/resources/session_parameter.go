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

var sessionParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Name of session parameter. Valid values are those in [session parameters](https://docs.snowflake.com/en/sql-reference/parameters.html#session-parameters).",
		ValidateFunc: validation.StringInSlice(maps.Keys(snowflake.GetParameters(snowflake.ParameterTypeSession)), false),
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of session parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation.",
	},
}

func SessionParameter() *schema.Resource {
	return &schema.Resource{
		Create: CreateSessionParameter,
		Read:   ReadSessionParameter,
		Update: UpdateSessionParameter,
		Delete: DeleteSessionParameter,

		Schema: sessionParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateSessionParameter implements schema.CreateFunc.
func CreateSessionParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	parameterDefault := snowflake.GetParameters(snowflake.ParameterTypeSession)[key]
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

	builder := snowflake.NewParameter(key, value, snowflake.ParameterTypeSession, db)
	err := builder.SetParameter()
	if err != nil {
		return fmt.Errorf("error creating session parameter err = %v", err)
	}

	d.SetId(key)
	p, err := snowflake.ShowParameter(db, key, snowflake.ParameterTypeSession)
	if err != nil {
		return fmt.Errorf("error reading session parameter err = %v", err)
	}
	d.Set("value", p.Value.String)
	return nil
}

// ReadSessionParameter implements schema.ReadFunc.
func ReadSessionParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	key := d.Id()
	p, err := snowflake.ShowParameter(db, key, snowflake.ParameterTypeSession)
	if err != nil {
		return fmt.Errorf("error reading session parameter err = %v", err)
	}
	d.Set("value", p.Value.String)
	return nil
}

// UpdateSessionParameter implements schema.UpdateFunc.
func UpdateSessionParameter(d *schema.ResourceData, meta interface{}) error {
	return CreateSessionParameter(d, meta)
}

// DeleteSessionParameter implements schema.DeleteFunc.
func DeleteSessionParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	key := d.Get("key").(string)

	parameterDefault := snowflake.GetParameters(snowflake.ParameterTypeSession)[key]
	defaultValue := parameterDefault.DefaultValue
	value := fmt.Sprintf("%v", defaultValue)

	// add quotes to value if it is a string
	typeString := reflect.TypeOf("")
	if reflect.TypeOf(parameterDefault.DefaultValue) == typeString {
		value = fmt.Sprintf("'%s'", value)
	}
	builder := snowflake.NewParameter(key, value, snowflake.ParameterTypeSession, db)
	err := builder.SetParameter()
	if err != nil {
		return fmt.Errorf("error creating account parameter err = %v", err)
	}
	_, err = snowflake.ShowParameter(db, key, snowflake.ParameterTypeSession)
	if err != nil {
		return fmt.Errorf("error reading a parameter err = %v", err)
	}

	d.SetId("")
	return nil
}
