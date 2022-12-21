package datasources

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var parametersSchema = map[string]*schema.Schema{
	"parameter_type": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "SESSION",
		Description:  "The type of parameter to filter by. Valid values are: \"ACCOUNT\", \"SESSION\", \"OBJECT\".",
		ValidateFunc: validation.StringInSlice([]string{"ACCOUNT", "SESSION", "OBJECT"}, true),
	},
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Allows limiting the list of parameters by name using LIKE clause. Refer to [Limiting the List of Parameters by Name](https://docs.snowflake.com/en/sql-reference/parameters.html#limiting-the-list-of-parameters-by-name)",
	},
	"object_type": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "If parameter_type is set to \"OBJECT\" then object_type is the type of object to display object parameters for. Valid values are any object supported by the IN clause of the [SHOW PARAMETERS](https://docs.snowflake.com/en/sql-reference/sql/show-parameters.html#parameters) statement, including: WAREHOUSE | DATABASE | SCHEMA | TASK | TABLE",
		ValidateFunc: validation.StringInSlice(snowflake.GetParameterObjectTypeSetAsStrings(), false),
	},
	"object_name": {
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "If parameter_type is set to \"OBJECT\" then object_name is the name of the object to display object parameters for.",
	},
	"parameters": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The pipes in the schema",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The name of the parameter",
				},
				"value": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The value of the parameter",
				},
				"default": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The default value of the parameter",
				},
				"level": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The level of the parameter",
				},
				"description": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The description of the parameter",
				},
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The type of the parameter",
				},
			},
		},
	},
}

func Parameters() *schema.Resource {
	return &schema.Resource{
		Read:   ReadParameters,
		Schema: parametersSchema,
	}
}

func ReadParameters(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	p, ok := d.GetOk("pattern")
	pattern := ""
	if ok {
		pattern = p.(string)
	}
	parameterType := snowflake.ParameterType(strings.ToUpper(d.Get("parameter_type").(string)))
	if parameterType == snowflake.ParameterTypeObject {
		oType := d.Get("object_type").(string)
		objectType := snowflake.ObjectType(oType)
		objectName := d.Get("object_name").(string)
		parameters, err := snowflake.ListObjectParameters(db, objectType, objectName, pattern)
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[DEBUG] parameters not found")
			d.SetId("")
			return nil
		} else if err != nil {
			log.Printf("[DEBUG] error occured during read: %v", err.Error())
			return err
		}
		d.SetId("parameters")
		params := []map[string]interface{}{}
		for _, param := range parameters {
			paramMap := map[string]interface{}{}

			paramMap["key"] = param.Key.String
			paramMap["value"] = param.Value.String
			paramMap["default"] = param.Default.String
			paramMap["level"] = param.Level.String
			paramMap["description"] = param.Description.String
			paramMap["type"] = param.PType.String

			params = append(params, paramMap)
		}
		return d.Set("parameters", params)
	} else {
		parameters, err := snowflake.ListParameters(db, parameterType, pattern)
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[DEBUG] parameters not found")
			d.SetId("")
			return nil
		} else if err != nil {
			log.Printf("[DEBUG] error occured during read: %v", err.Error())
			return err
		}
		d.SetId("parameters")

		params := []map[string]interface{}{}
		for _, param := range parameters {
			paramMap := map[string]interface{}{}

			paramMap["key"] = param.Key.String
			paramMap["value"] = param.Value.String
			paramMap["default"] = param.Default.String
			paramMap["level"] = param.Level.String
			paramMap["description"] = param.Description.String
			paramMap["type"] = param.PType.String

			params = append(params, paramMap)
		}
		return d.Set("parameters", params)
	}
}
