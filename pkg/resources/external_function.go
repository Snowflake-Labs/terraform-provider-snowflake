package resources

import (
	"context"
	"database/sql"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		Description: "Indicates whether the function can return NULL values (true) or must return only NON-NULL values (false).",
		Default:     true,
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
	"request_translator": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "This specifies the name of the request translator function",
	},
	"response_translator": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "This specifies the name of the response translator function.",
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
		Description: "A description of the external function.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the external function was created.",
	},
}

// ExternalFunction returns a pointer to the resource representing an external function.
func ExternalFunction() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: CreateContextExternalFunction,
		ReadContext:   ReadContextExternalFunction,
		UpdateContext: UpdateContextExternalFunction,
		DeleteContext: DeleteContextExternalFunction,

		Schema: externalFunctionSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v085ExternalFunctionStateUpgrader,
			},
		},
	}
}

func CreateContextExternalFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	database := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	returnType := d.Get("return_type").(string)
	resultDataType, err := sdk.ToDataType(returnType)
	if err != nil {
		return diag.FromErr(err)
	}
	apiIntegration := sdk.NewAccountObjectIdentifier(d.Get("api_integration").(string))
	urlOfProxyAndResource := d.Get("url_of_proxy_and_resource").(string)
	req := sdk.NewCreateExternalFunctionRequest(id, resultDataType, &apiIntegration, urlOfProxyAndResource)

	// Set optionals
	args := make([]sdk.ExternalFunctionArgumentRequest, 0)
	if v, ok := d.GetOk("arg"); ok {
		for _, arg := range v.([]interface{}) {
			argName := arg.(map[string]interface{})["name"].(string)
			argType := arg.(map[string]interface{})["type"].(string)
			argDataType, err := sdk.ToDataType(argType)
			if err != nil {
				return diag.FromErr(err)
			}
			args = append(args, sdk.ExternalFunctionArgumentRequest{ArgName: argName, ArgDataType: argDataType})
		}
	}
	if len(args) > 0 {
		req.WithArguments(args)
	}

	if v, ok := d.GetOk("return_null_allowed"); ok {
		if v.(bool) {
			req.WithReturnNullValues(&sdk.ReturnNullValuesNull)
		} else {
			req.WithReturnNullValues(&sdk.ReturnNullValuesNotNull)
		}
	}

	if v, ok := d.GetOk("return_behavior"); ok {
		if v.(string) == "VOLATILE" {
			req.WithReturnResultsBehavior(&sdk.ReturnResultsBehaviorVolatile)
		} else {
			req.WithReturnResultsBehavior(&sdk.ReturnResultsBehaviorImmutable)
		}
	}

	if v, ok := d.GetOk("null_input_behavior"); ok {
		switch {
		case v.(string) == "CALLED ON NULL INPUT":
			req.WithNullInputBehavior(sdk.Pointer(sdk.NullInputBehaviorCalledOnNullInput))
		case v.(string) == "RETURNS NULL ON NULL INPUT":
			req.WithNullInputBehavior(sdk.Pointer(sdk.NullInputBehaviorReturnNullInput))
		default:
			req.WithNullInputBehavior(sdk.Pointer(sdk.NullInputBehaviorStrict))
		}
	}

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(sdk.String(v.(string)))
	}

	if _, ok := d.GetOk("header"); ok {
		headers := make([]sdk.ExternalFunctionHeaderRequest, 0)
		for _, header := range d.Get("header").(*schema.Set).List() {
			m := header.(map[string]interface{})
			headerName := m["name"].(string)
			headerValue := m["value"].(string)
			headers = append(headers, sdk.ExternalFunctionHeaderRequest{
				Name:  headerName,
				Value: headerValue,
			})
		}
		req.WithHeaders(headers)
	}

	if v, ok := d.GetOk("context_headers"); ok {
		contextHeadersList := expandStringList(v.([]interface{}))
		contextHeaders := make([]sdk.ExternalFunctionContextHeaderRequest, 0)
		for _, header := range contextHeadersList {
			contextHeaders = append(contextHeaders, sdk.ExternalFunctionContextHeaderRequest{
				ContextFunction: header,
			})
		}
		req.WithContextHeaders(contextHeaders)
	}

	if v, ok := d.GetOk("max_batch_rows"); ok {
		req.WithMaxBatchRows(sdk.Int(v.(int)))
	}

	if v, ok := d.GetOk("compression"); ok {
		req.WithCompression(sdk.String(v.(string)))
	}

	if v, ok := d.GetOk("request_translator"); ok {
		req.WithRequestTranslator(sdk.Pointer(sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(v.(string))))
	}

	if v, ok := d.GetOk("response_translator"); ok {
		req.WithResponseTranslator(sdk.Pointer(sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(v.(string))))
	}

	if err := client.ExternalFunctions.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	argTypes := make([]sdk.DataType, 0, len(args))
	for _, item := range args {
		argTypes = append(argTypes, item.ArgDataType)
	}
	sid := sdk.NewSchemaObjectIdentifierWithArguments(database, schemaName, name, argTypes)
	d.SetId(sid.FullyQualifiedName())
	return ReadContextExternalFunction(ctx, d, meta)
}

func ReadContextExternalFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(d.Id())
	externalFunction, err := client.ExternalFunctions.ShowByID(ctx, id.WithoutArguments(), id.Arguments())
	if err != nil {
		d.SetId("")
		return nil
	}

	// Some properties can come from the SHOW EXTERNAL FUNCTION call
	if err := d.Set("name", externalFunction.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("schema", strings.Trim(externalFunction.SchemaName, "\"")); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database", strings.Trim(externalFunction.CatalogName, "\"")); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", externalFunction.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_on", externalFunction.CreatedOn); err != nil {
		return diag.FromErr(err)
	}

	// Some properties come from the DESCRIBE FUNCTION call
	externalFunctionPropertyRows, err := client.ExternalFunctions.Describe(ctx, sdk.NewDescribeExternalFunctionRequest(id.WithoutArguments(), id.Arguments()))
	if err != nil {
		d.SetId("")
		return nil
	}

	for _, row := range externalFunctionPropertyRows {
		switch row.Property {
		case "signature":
			// Format in Snowflake DB is: (argName argType, argName argType, ...)
			args := strings.ReplaceAll(strings.ReplaceAll(row.Value, "(", ""), ")", "")

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

				if err := d.Set("arg", args); err != nil {
					return diag.Errorf("error setting arg: %v", err)
				}
			}
		case "returns":
			returnType := row.Value
			// We first check for VARIANT or OBJECT
			if returnType == "VARIANT" || returnType == "OBJECT" {
				if err := d.Set("return_type", returnType); err != nil {
					return diag.Errorf("error setting return_type: %v", err)
				}
				break
			}

			// otherwise, format in Snowflake DB is returnType(<some number>)
			re := regexp.MustCompile(`^(\w+)\([0-9]*\)$`)
			match := re.FindStringSubmatch(row.Value)
			if len(match) < 2 {
				return diag.Errorf("return_type %s not recognized", returnType)
			}
			if err := d.Set("return_type", match[1]); err != nil {
				return diag.Errorf("error setting return_type: %v", err)
			}

		case "null handling":
			if err := d.Set("null_input_behavior", row.Value); err != nil {
				return diag.Errorf("error setting null_input_behavior: %v", err)
			}
		case "volatility":
			if err := d.Set("return_behavior", row.Value); err != nil {
				return diag.Errorf("error setting return_behavior: %v", err)
			}
		case "headers":
			if row.Value != "" && row.Value != "null" {
				// Format in Snowflake DB is: {"head1":"val1","head2":"val2"}
				headerPairs := strings.Split(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(row.Value, "{", ""), "}", ""), "\"", ""), ",")
				headers := []interface{}{}

				for _, headerPair := range headerPairs {
					headerItem := strings.Split(headerPair, ":")

					header := map[string]interface{}{}
					header["name"] = headerItem[0]
					header["value"] = headerItem[1]
					headers = append(headers, header)
				}

				if err := d.Set("header", headers); err != nil {
					return diag.Errorf("error setting return_behavior: %v", err)
				}
			}
		case "context_headers":
			if row.Value != "" && row.Value != "null" {
				// Format in Snowflake DB is: ["CONTEXT_FUNCTION_1","CONTEXT_FUNCTION_2"]
				contextHeaders := strings.Split(strings.Trim(row.Value, "[]"), ",")
				for i, v := range contextHeaders {
					contextHeaders[i] = strings.Trim(v, "\"")
				}
				if err := d.Set("context_headers", contextHeaders); err != nil {
					return diag.Errorf("error setting context_headers: %v", err)
				}
			}
		case "max_batch_rows":
			if row.Value != "not set" {
				maxBatchRows, err := strconv.ParseInt(row.Value, 10, 64)
				if err != nil {
					return diag.Errorf("error parsing max_batch_rows: %v", err)
				}
				if err := d.Set("max_batch_rows", maxBatchRows); err != nil {
					return diag.Errorf("error setting max_batch_rows: %v", err)
				}
			}
		case "compression":
			if err := d.Set("compression", row.Value); err != nil {
				return diag.Errorf("error setting compression: %v", err)
			}
		case "body":
			if err := d.Set("url_of_proxy_and_resource", row.Value); err != nil {
				return diag.Errorf("error setting url_of_proxy_and_resource: %v", err)
			}
		case "language":
			// To ignore
		default:
			log.Printf("[WARN] unexpected external function property %v returned from Snowflake", row.Property)
		}
	}

	return nil
}

func UpdateContextExternalFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(d.Id())
	req := sdk.NewAlterFunctionRequest(id.WithoutArguments(), id.Arguments())
	if d.HasChange("comment") {
		_, new := d.GetChange("comment")
		if new == "" {
			req.UnsetComment = sdk.Bool(true)
		} else {
			req.SetComment = sdk.String(new.(string))
		}
		err := client.Functions.Alter(ctx, req)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadContextExternalFunction(ctx, d, meta)
}

func DeleteContextExternalFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(d.Id())
	req := sdk.NewDropFunctionRequest(id.WithoutArguments(), id.Arguments())
	if err := client.Functions.Drop(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
