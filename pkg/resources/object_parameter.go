package resources

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	snowflakeValidation "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"
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
	"object_identifier": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		Description: "Specifies the object identifier for the object parameter.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Name of the object to set the parameter for.",
				},
				"database": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the database that the object was created in.",
				},
				"schema": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the schema that the object was created in.",
				},
			},
		},
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
	objectType := snowflake.ObjectType(d.Get("object_type").(string))
	objectDatabase, objectSchema, objectName := expandObjectIdentifier(d.Get("object_identifier"))
	fullyQualifierObjectIdentifier := snowflakeValidation.FormatFullyQualifiedObjectID(objectDatabase, objectSchema, objectName)

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
	builder.WithObjectIdentifier(fullyQualifierObjectIdentifier)
	builder.WithObjectType(objectType)
	err := builder.SetParameter()
	if err != nil {
		return fmt.Errorf("error creating object parameter err = %w", err)
	}
	id := fmt.Sprintf("%v❄️%v❄️%v", key, objectType, fullyQualifierObjectIdentifier)
	d.SetId(id)
	p, err := snowflake.ShowObjectParameter(db, key, objectType, fullyQualifierObjectIdentifier)
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
		return fmt.Errorf("unexpected format of ID (%v), expected key❄️object_type❄️object_identifier", id)
	}
	key := parts[0]
	objectType := snowflake.ObjectType(parts[1])
	objectIdentifier := parts[2]
	p, err := snowflake.ShowObjectParameter(db, key, objectType, objectIdentifier)
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
	objectDatabase, objectSchema, objectName := expandObjectIdentifier(d.Get("object_identifier"))
	fullyQualifierObjectIdentifier := snowflakeValidation.FormatFullyQualifiedObjectID(objectDatabase, objectSchema, objectName)

	parameterDefault := snowflake.GetParameterDefaults(snowflake.ParameterTypeObject)[key]
	defaultValue := parameterDefault.DefaultValue
	value := fmt.Sprintf("%v", defaultValue)

	// add quotes to value if it is a string
	typeString := reflect.TypeOf("")
	if reflect.TypeOf(parameterDefault.DefaultValue) == typeString {
		value = fmt.Sprintf("'%s'", value)
	}
	builder := snowflake.NewParameter(key, value, snowflake.ParameterTypeObject, db)
	builder.WithObjectIdentifier(fullyQualifierObjectIdentifier)
	builder.WithObjectType(objectType)
	err := builder.SetParameter()
	if err != nil {
		return fmt.Errorf("error restoring default for object parameter err = %w", err)
	}
	_, err = snowflake.ShowObjectParameter(db, key, objectType, fullyQualifierObjectIdentifier)
	if err != nil {
		return fmt.Errorf("error reading object parameter err = %w", err)
	}

	d.SetId("")
	return nil
}
