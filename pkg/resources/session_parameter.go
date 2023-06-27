package resources

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
		ValidateFunc: validation.StringInSlice(maps.Keys(snowflake.GetParameterDefaults(snowflake.ParameterTypeSession)), false),
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of session parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation.",
	},
	"on_account": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If true, the session parameter will be set on the account level.",
	},
	"user": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The user to set the session parameter for. Required if on_account is false",
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
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	parameterDefault := snowflake.GetParameterDefaults(snowflake.ParameterTypeSession)[key]
	if parameterDefault.Validate != nil {
		if err := parameterDefault.Validate(value); err != nil {
			return err
		}
	}

	onAccount := d.Get("on_account").(bool)
	user := d.Get("user").(string)
	parameter := sdk.SessionParameter(key)

	var err error
	if onAccount {
		err := client.Parameters.SetSessionParameterForAccount(ctx, parameter, value)
		if err != nil {
			return err
		}
	} else {
		if user == "" {
			return fmt.Errorf("user is required if on_account is false")
		}
		// add quotes to value if it is a string
		typeString := reflect.TypeOf("")
		if reflect.TypeOf(parameterDefault.DefaultValue) == typeString {
			value = fmt.Sprintf("'%s'", snowflake.EscapeString(value))
		}
		builder := snowflake.NewSessionParameter(key, value, db)
		builder.SetUser(user)
		err = builder.SetParameter()
		if err != nil {
			return fmt.Errorf("error creating session parameter err = %w", err)
		}
	}

	d.SetId(key)

	return ReadSessionParameter(d, meta)
}

// ReadSessionParameter implements schema.ReadFunc.
func ReadSessionParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	parameter := d.Id()

	onAccount := d.Get("on_account").(bool)
	var err error
	var p *sdk.Parameter
	if onAccount {
		p, err = client.Sessions.ShowAccountParameter(ctx, sdk.AccountParameter(parameter))
	} else {
		user := d.Get("user").(string)
		userId := sdk.NewAccountObjectIdentifier(user)
		p, err = client.Sessions.ShowUserParameter(ctx, sdk.UserParameter(parameter), userId)
	}
	if err != nil {
		return fmt.Errorf("error reading session parameter err = %w", err)
	}
	err = d.Set("value", p.Value)
	if err != nil {
		return fmt.Errorf("error setting session parameter err = %w", err)
	}
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
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	onAccount := d.Get("on_account").(bool)
	parameter := sdk.SessionParameter(key)

	var err error
	if onAccount {
		defaultParameter, err := client.Sessions.ShowAccountParameter(ctx, sdk.AccountParameter(key))
		if err != nil {
			return err
		}
		defaultValue := defaultParameter.Default
		err = client.Parameters.SetSessionParameterForAccount(ctx, parameter, defaultValue)
		if err != nil {
			return fmt.Errorf("error creating session parameter err = %w", err)
		}
	} else {
		user := d.Get("user").(string)
		if user == "" {
			return fmt.Errorf("user is required if on_account is false")
		}
		parameterDefault := snowflake.GetParameterDefaults(snowflake.ParameterTypeSession)[key]
		defaultValue := parameterDefault.DefaultValue
		value := fmt.Sprintf("%v", defaultValue)
		typeString := reflect.TypeOf("")
		if reflect.TypeOf(parameterDefault.DefaultValue) == typeString {
			value = fmt.Sprintf("'%s'", value)
		}
		builder := snowflake.NewSessionParameter(key, value, db)
		builder.SetUser(user)
		err = builder.SetParameter()
		if err != nil {
			return fmt.Errorf("error deleting session parameter err = %w", err)
		}
	}

	d.SetId(key)
	return nil
}
