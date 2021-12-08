package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"
)

const (
	externalFunctionIDDelimiter = '|'
)

var externalFunctionSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the external function. The identifier can contain the schema name and database name, as well as the function name. The function's signature (name and argument data types) must be unique within the schema.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the external function.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the external function.",
	},
	"arg": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the arguments/inputs for the external function. These should correspond to the arguments that the remote service expects.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
					// Suppress the diff shown if the values are equal when both compared in lower case.
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return strings.EqualFold(strings.ToLower(old), strings.ToLower(new))
					},
					Description: "Argument name",
				},
				"type": {
					Type:     schema.TypeString,
					Required: true,
					// Suppress the diff shown if the values are equal when both compared in lower case.
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return strings.EqualFold(strings.ToLower(old), strings.ToLower(new))
					},
					Description: "Argument type, e.g. VARCHAR",
				},
			},
		},
	},
	"null_input_behavior": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "CALLED ON NULL INPUT",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"CALLED ON NULL INPUT", "RETURNS NULL ON NULL INPUT", "STRICT"}, false),
		Description:  "Specifies the behavior of the external function when called with null inputs.",
	},
	"return_type": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		// Suppress the diff shown if the values are equal when both compared in lower case.
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.EqualFold(strings.ToLower(old), strings.ToLower(new))
		},
		Description: "Specifies the data type returned by the external function.",
	},
	"return_null_allowed": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Description: "Indicates whether the function can return NULL values or must return only NON-NULL values.",
	},
	"return_behavior": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"VOLATILE", "IMMUTABLE"}, false),
		Description:  "Specifies the behavior of the function when returning results",
	},
	"api_integration": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The name of the API integration object that should be used to authenticate the call to the proxy service.",
	},
	"header": {
		Type:        schema.TypeSet,
		Optional:    true,
		ForceNew:    true,
		Description: "Allows users to specify key-value metadata that is sent with every request as HTTP headers.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Header name",
				},
				"value": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Header value",
				},
			},
		},
	},
	"context_headers": {
		Type:     schema.TypeList,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Optional: true,
		ForceNew: true,
		// Suppress the diff shown if the values are equal when both compared in lower case.
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.EqualFold(strings.ToLower(old), strings.ToLower(new))
		},
		Description: "Binds Snowflake context function results to HTTP headers.",
	},
	"max_batch_rows": {
		Type:        schema.TypeInt,
		Optional:    true,
		ForceNew:    true,
		Description: "This specifies the maximum number of rows in each batch sent to the proxy service.",
	},
	"compression": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "AUTO",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"NONE", "AUTO", "GZIP", "DEFLATE"}, false),
		Description:  "If specified, the JSON payload is compressed when sent from Snowflake to the proxy service, and when sent back from the proxy service to Snowflake.",
	},
	"url_of_proxy_and_resource": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "This is the invocation URL of the proxy service and resource through which Snowflake calls the remote service.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "user-defined function",
		ForceNew:    true,
		Description: "A description of the external function.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the external function was created.",
	},
}

// ExternalFunction returns a pointer to the resource representing an external function
func ExternalFunction() *schema.Resource {
	return &schema.Resource{
		Create: CreateExternalFunction,
		Read:   ReadExternalFunction,
		Delete: DeleteExternalFunction,

		Schema: externalFunctionSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

type externalFunctionID struct {
	DatabaseName             string
	SchemaName               string
	ExternalFunctionName     string
	ExternalFunctionArgTypes string
}

func (si *externalFunctionID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = externalFunctionIDDelimiter
	err := csvWriter.WriteAll([][]string{{si.DatabaseName, si.SchemaName, si.ExternalFunctionName, si.ExternalFunctionArgTypes}})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func externalFunctionIDFromString(stringID string) (*externalFunctionID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = externalFunctionIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line at a time")
	}
	if len(lines[0]) != 4 {
		return nil, fmt.Errorf("4 fields allowed")
	}

	return &externalFunctionID{
		DatabaseName:             lines[0][0],
		SchemaName:               lines[0][1],
		ExternalFunctionName:     lines[0][2],
		ExternalFunctionArgTypes: lines[0][3],
	}, nil
}

// CreateExternalFunction implements schema.CreateFunc
func CreateExternalFunction(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	dbSchema := d.Get("schema").(string)
	name := d.Get("name").(string)
	var argtypes string

	builder := snowflake.ExternalFunction(name, database, dbSchema)
	builder.WithReturnType(d.Get("return_type").(string))
	builder.WithReturnBehavior(d.Get("return_behavior").(string))
	builder.WithAPIIntegration(d.Get("api_integration").(string))
	builder.WithURLOfProxyAndResource(d.Get("url_of_proxy_and_resource").(string))

	// Set optionals
	if _, ok := d.GetOk("arg"); ok {
		var types []string
		args := []map[string]string{}
		for _, arg := range d.Get("arg").([]interface{}) {
			argDef := map[string]string{}
			for key, val := range arg.(map[string]interface{}) {
				argDef[key] = val.(string)

				if key == "type" {
					// Also store arg types in distinct array as list of types is required for some Snowflake commands (DESC, DROP)
					types = append(types, argDef[key])
				}
			}
			args = append(args, argDef)
		}

		// Use '-' as a separator between arg types as the result will end in the Terraform resource id
		argtypes = strings.Join(types, "-")

		builder.WithArgs(args)
		builder.WithArgTypes(argtypes)
	}

	if v, ok := d.GetOk("return_null_allowed"); ok {
		builder.WithReturnNullAllowed(v.(bool))
	}

	if v, ok := d.GetOk("null_input_behavior"); ok {
		builder.WithNullInputBehavior(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if _, ok := d.GetOk("header"); ok {
		headers := []map[string]string{}
		for _, header := range d.Get("header").(*schema.Set).List() {
			headerDef := map[string]string{}
			for key, val := range header.(map[string]interface{}) {
				headerDef[key] = val.(string)
			}
			headers = append(headers, headerDef)
		}

		builder.WithHeaders(headers)
	}

	if v, ok := d.GetOk("context_headers"); ok {
		contextHeaders := expandStringList(v.([]interface{}))
		builder.WithContextHeaders(contextHeaders)
	}

	if v, ok := d.GetOk("max_batch_rows"); ok {
		builder.WithMaxBatchRows(v.(int))
	}

	if v, ok := d.GetOk("compression"); ok {
		builder.WithCompression(v.(string))
	}

	stmt := builder.Create()
	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating external function %v", name)
	}

	externalFunctionID := &externalFunctionID{
		DatabaseName:             database,
		SchemaName:               dbSchema,
		ExternalFunctionName:     name,
		ExternalFunctionArgTypes: argtypes,
	}
	dataIDInput, err := externalFunctionID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadExternalFunction(d, meta)
}

// ReadExternalFunction implements schema.ReadFunc
func ReadExternalFunction(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	externalFunctionID, err := externalFunctionIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := externalFunctionID.DatabaseName
	dbSchema := externalFunctionID.SchemaName
	name := externalFunctionID.ExternalFunctionName
	argtypes := externalFunctionID.ExternalFunctionArgTypes

	// Some properties can come from the SHOW EXTERNAL FUNCTION call
	stmt := snowflake.ExternalFunction(name, dbName, dbSchema).Show()
	row := snowflake.QueryRow(db, stmt)
	externalFunction, err := snowflake.ScanExternalFunction(row)
	if err != nil {
		return err
	}

	// Note: 'language' must be EXTERNAL and 'is_external_function' set to Y
	if externalFunction.Language.String != "EXTERNAL" || externalFunction.IsExternalFunction.String != "Y" {
		return fmt.Errorf("Expected %v to be an external function, got 'language=%v' and 'is_external_function=%v'", d.Id(), externalFunction.Language.String, externalFunction.IsExternalFunction.String)
	}

	if err := d.Set("name", externalFunction.ExternalFunctionName.String); err != nil {
		return err
	}

	if err := d.Set("schema", externalFunction.SchemaName.String); err != nil {
		return err
	}

	if err := d.Set("database", externalFunction.DatabaseName.String); err != nil {
		return err
	}

	if err := d.Set("comment", externalFunction.Comment.String); err != nil {
		return err
	}

	if err := d.Set("created_on", externalFunction.CreatedOn.String); err != nil {
		return err
	}

	// Some properties come from the DESCRIBE FUNCTION call
	stmt = snowflake.ExternalFunction(name, dbName, dbSchema).WithArgTypes(argtypes).Describe()
	externalFunctionDescriptionRows, err := snowflake.Query(db, stmt)
	if err != nil {
		return err
	}

	externalFunctionDescription, err := snowflake.ScanExternalFunctionDescription(externalFunctionDescriptionRows)
	if err != nil {
		return err
	}

	for _, desc := range externalFunctionDescription {
		switch desc.Property.String {
		case "signature":
			// Format in Snowflake DB is: (argName argType, argName argType, ...)
			args := strings.ReplaceAll(strings.ReplaceAll(desc.Value.String, "(", ""), ")", "")

			if args != "" { // Do nothing for functions without arguments
				argPairs := strings.Split(args, ", ")
				args := []interface{}{}

				for _, argPair := range argPairs {
					argItem := strings.Split(argPair, " ")

					arg := map[string]interface{}{}
					arg["name"] = argItem[0]
					arg["type"] = argItem[1]
					args = append(args, arg)
				}

				if err = d.Set("arg", args); err != nil {
					return err
				}
			}
		case "returns":
			returnType := desc.Value.String
			// We first check for VARIANT
			if returnType == "VARIANT" {
				if err = d.Set("return_type", returnType); err != nil {
					return err
				}
				break
			}

			// otherwise, format in Snowflake DB is returnType(<some number>)
			re := regexp.MustCompile(`^(\w+)\([0-9]*\)$`)
			match := re.FindStringSubmatch(desc.Value.String)
			if len(match) < 2 {
				return errors.Errorf("return_type %s not recognized", returnType)
			}
			if err = d.Set("return_type", match[1]); err != nil {
				return err
			}

		case "null handling":
			if err = d.Set("null_input_behavior", desc.Value.String); err != nil {
				return err
			}
		case "volatility":
			if err = d.Set("return_behavior", desc.Value.String); err != nil {
				return err
			}
		case "headers":
			if desc.Value.Valid && desc.Value.String != "null" {
				// Format in Snowflake DB is: {"head1":"val1","head2":"val2"}
				headerPairs := strings.Split(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(desc.Value.String, "{", ""), "}", ""), "\"", ""), ",")
				headers := []interface{}{}

				for _, headerPair := range headerPairs {
					headerItem := strings.Split(headerPair, ":")

					header := map[string]interface{}{}
					header["name"] = headerItem[0]
					header["value"] = headerItem[1]
					headers = append(headers, header)
				}

				if err = d.Set("header", headers); err != nil {
					return err
				}
			}
		case "context_headers":
			if desc.Value.Valid && desc.Value.String != "null" {
				// Format in Snowflake DB is: ["CONTEXT_FUNCTION_1","CONTEXT_FUNCTION_2"]
				contextHeaders := strings.Split(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(desc.Value.String, "[", ""), "]", ""), "\"", ""), ",")

				if err = d.Set("context_headers", contextHeaders); err != nil {
					return err
				}
			}
		case "max_batch_rows":
			if desc.Value.String != "not set" {
				i, err := strconv.ParseInt(desc.Value.String, 10, 64)
				if err != nil {
					return err
				}

				if err = d.Set("max_batch_rows", i); err != nil {
					return err
				}
			}
		case "compression":
			if err = d.Set("compression", desc.Value.String); err != nil {
				return err
			}
		case "body":
			if err = d.Set("url_of_proxy_and_resource", desc.Value.String); err != nil {
				return err
			}
		case "language":
			// To ignore
		default:
			log.Printf("[WARN] unexpected external function property %v returned from Snowflake", desc.Property.String)
		}
	}

	return nil
}

// DeleteExternalFunction implements schema.DeleteFunc
func DeleteExternalFunction(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	externalFunctionID, err := externalFunctionIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := externalFunctionID.DatabaseName
	dbSchema := externalFunctionID.SchemaName
	name := externalFunctionID.ExternalFunctionName
	argtypes := externalFunctionID.ExternalFunctionArgTypes

	q := snowflake.ExternalFunction(name, dbName, dbSchema).WithArgTypes(argtypes).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting external function %v", d.Id())
	}

	d.SetId("")
	return nil
}
