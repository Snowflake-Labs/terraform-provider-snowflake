package resources

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var objectParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Name of object parameter. Valid values are those in [object parameters](https://docs.snowflake.com/en/sql-reference/parameters.html#object-parameters).",
		ValidateFunc: validation.StringInSlice(maps.Keys(snowflake.GetParameterDefaults(snowflake.ParameterTypeObject)), false),
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of object parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation.",
	},
	"object_type": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Type of object to which the parameter applies. Valid values are those in [object types](https://docs.snowflake.com/en/sql-reference/parameters.html#object-types).",
		ValidateFunc: validation.StringInSlice(snowflake.GetParameterObjectTypeSetAsStrings(), false),
	},
	"object_name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Name of object to which the parameter applies.",
	},
}

func ObjectParameter() *schema.Resource {
	return &schema.Resource{
		Create: CreateObjectParameter,
		Read:   ReadObjectParameter,
		Update: UpdateObjectParameter,
		Delete: DeleteObjectParameter,

		Schema: objectParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateObjectParameter implements schema.CreateFunc.
func CreateObjectParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	objectName := d.Get("object_name").(string)
	objectType := snowflake.ObjectType(d.Get("object_type").(string))
	parameterDefault := snowflake.GetParameterDefaults(snowflake.ParameterTypeObject)[key]
	if parameterDefault.Validate != nil {
		if err := parameterDefault.Validate(value); err != nil {
			return err
		}
	}
	if ok := slices.Contains(parameterDefault.AllowedObjectTypes, objectType); !ok {
		return fmt.Errorf("object_type '%v' is not allowed for parameter '%v'", objectType, key)
	}

	// add quotes to value if it is a string
	typeString := reflect.TypeOf("")
	if reflect.TypeOf(parameterDefault.DefaultValue) == typeString {
		value = fmt.Sprintf("'%s'", snowflake.EscapeString(value))
	}

	builder := snowflake.NewParameter(key, value, snowflake.ParameterTypeObject, db)
	builder.WithObjectName(objectName)
	builder.WithObjectType(objectType)
	err := builder.SetParameter()
	if err != nil {
		return fmt.Errorf("error creating object parameter err = %w", err)
	}
	id := fmt.Sprintf("%v❄️%v❄️%v", key, objectType, objectName)
	d.SetId(id)
	p, err := snowflake.ShowObjectParameter(db, key, objectType, objectName)
	if err != nil {
		return fmt.Errorf("error reading object parameter err = %w", err)
	}
	err = d.Set("value", p.Value.String)
	if err != nil {
		return err
	}
	return nil
}

// ReadObjectParameter implements schema.ReadFunc.
func ReadObjectParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()
	parts := strings.Split(id, "❄️")
	if len(parts) != 3 {
		return fmt.Errorf("unexpected format of ID (%v), expected key❄️object_type❄️object_name", id)
	}
	key := parts[0]
	objectType := snowflake.ObjectType(parts[1])
	objectName := parts[2]
	p, err := snowflake.ShowObjectParameter(db, key, objectType, objectName)
	if err != nil {
		return fmt.Errorf("error reading object parameter err = %w", err)
	}
	err = d.Set("value", p.Value.String)
	if err != nil {
		return err
	}
	return nil
}

// UpdateObjectParameter implements schema.UpdateFunc.
func UpdateObjectParameter(d *schema.ResourceData, meta interface{}) error {
	return CreateObjectParameter(d, meta)
}

// DeleteObjectParameter implements schema.DeleteFunc.
func DeleteObjectParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	key := d.Get("key").(string)
	objectType := snowflake.ObjectType(d.Get("object_type").(string))
	objectName := d.Get("object_name").(string)

	parameterDefault := snowflake.GetParameterDefaults(snowflake.ParameterTypeObject)[key]
	defaultValue := parameterDefault.DefaultValue
	value := fmt.Sprintf("%v", defaultValue)

	// add quotes to value if it is a string
	typeString := reflect.TypeOf("")
	if reflect.TypeOf(parameterDefault.DefaultValue) == typeString {
		value = fmt.Sprintf("'%s'", value)
	}
	builder := snowflake.NewParameter(key, value, snowflake.ParameterTypeObject, db)
	builder.WithObjectName(objectName)
	builder.WithObjectType(objectType)
	err := builder.SetParameter()
	if err != nil {
		return fmt.Errorf("error restoring default for object parameter err = %w", err)
	}
	_, err = snowflake.ShowObjectParameter(db, key, objectType, objectName)
	if err != nil {
		return fmt.Errorf("error reading object parameter err = %w", err)
	}

	d.SetId("")
	return nil
}
