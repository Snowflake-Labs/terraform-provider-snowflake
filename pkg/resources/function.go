package resources

import (
	"context"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"regexp"
	"strings"
)

var languages = []string{"javascript", "scala", "java", "sql", "python"}

var functionSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the function; does not have to be unique for the schema in which the function is created. Don't use the | character.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the function. Don't use the | character.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the function. Don't use the | character.",
		ForceNew:    true,
	},
	"arguments": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
					// Suppress the diff shown if the values are equal when both compared in lower case.
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return strings.EqualFold(old, new)
					},
					Description: "The argument name",
				},
				"type": {
					Type:     schema.TypeString,
					Required: true,
					// Suppress the diff shown if the values are equal when both compared in lower case.
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return strings.EqualFold(old, new)
					},
					Description: "The argument type",
				},
			},
		},
		Optional:    true,
		Description: "List of the arguments for the function",
		ForceNew:    true,
	},
	"return_type": {
		Type:        schema.TypeString,
		Description: "The return type of the function",
		// Suppress the diff shown if the values are equal when both compared in lower case.
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.EqualFold(old, new)
		},
		Required: true,
		ForceNew: true,
	},
	"statement": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the javascript / java / scala / sql / python code used to create the function.",
		ForceNew:         true,
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"language": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "SQL",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.EqualFold(old, new)
		},
		ValidateFunc: validation.StringInSlice(languages, true),
		Description:  "Specifies the language of the stored function code.",
	},
	"null_input_behavior": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "CALLED ON NULL INPUT",
		ForceNew: true,
		// We do not use STRICT, because Snowflake then in the Read phase returns RETURNS NULL ON NULL INPUT
		ValidateFunc: validation.StringInSlice([]string{"CALLED ON NULL INPUT", "RETURNS NULL ON NULL INPUT"}, false),
		Description:  "Specifies the behavior of the function when called with null inputs.",
	},
	"return_behavior": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "VOLATILE",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"VOLATILE", "IMMUTABLE"}, false),
		Description:  "Specifies the behavior of the function when returning results",
	},
	"is_secure": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies that the function is secure.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "user-defined function",
		Description: "Specifies a comment for the function.",
	},
	"runtime_version": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Required for Python functions. Specifies Python runtime version.",
	},
	"packages": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		ForceNew:    true,
		Description: "List of package imports to use for Java / Python functions. For Java, package imports should be of the form: package_name:version_number, where package_name is snowflake_domain:package. For Python use it should be: ('numpy','pandas','xgboost==1.5.0').",
	},
	"imports": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		ForceNew:    true,
		Description: "Imports for Java / Python functions. For Java this a list of jar files, for Python this is a list of Python files.",
	},
	"handler": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "The handler method for Java / Python function.",
	},
	"target_path": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "The target path for the Java / Python functions. For Java, it is the path of compiled jar files and for the Python it is the path of the Python files.",
	},
}

func Function() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: CreateContextFunction,
		ReadContext:   ReadContextFunction,
		UpdateContext: UpdateContextFunction,
		DeleteContext: DeleteContextFunction,

		Schema: functionSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v085FunctionIdStateUpgrader,
			},
		},
	}
}

func CreateContextFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lang := strings.ToUpper(d.Get("language").(string))
	switch lang {
	case "JAVA":
		return createJavaFunction(ctx, d, meta)
	case "JAVASCRIPT":
		return createJavascriptFunction(ctx, d, meta)
	case "PYTHON":
		return createPythonFunction(ctx, d, meta)
	case "SCALA":
		return createScalaFunction(ctx, d, meta)
	case "", "SQL": // SQL if language is not set
		return createSQLFunction(ctx, d, meta)
	default:
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid language",
				Detail:   fmt.Sprintf("Language %s is not supported", lang),
			},
		}
	}
}

func createJavaFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schema, name)

	// Set required
	returns, diags := parseFunctionReturnsRequest(d.Get("return_type").(string))
	if diags != nil {
		return diags
	}
	handler := d.Get("handler").(string)
	// create request with required
	request := sdk.NewCreateForJavaFunctionRequest(id, *returns, handler)
	functionDefinition := d.Get("statement").(string)
	request.WithFunctionDefinition(functionDefinition)

	// Set optionals
	if v, ok := d.GetOk("is_secure"); ok {
		request.WithSecure(v.(bool))
	}
	arguments, diags := parseFunctionArguments(d)
	if diags != nil {
		return diags
	}
	if len(arguments) > 0 {
		request.WithArguments(arguments)
	}
	if v, ok := d.GetOk("null_input_behavior"); ok {
		request.WithNullInputBehavior(sdk.NullInputBehavior(v.(string)))
	}
	if v, ok := d.GetOk("return_behavior"); ok {
		request.WithReturnResultsBehavior(sdk.ReturnResultsBehavior(v.(string)))
	}
	if v, ok := d.GetOk("runtime_version"); ok {
		request.WithRuntimeVersion(v.(string))
	}
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}
	if _, ok := d.GetOk("imports"); ok {
		imports := []sdk.FunctionImportRequest{}
		for _, item := range d.Get("imports").([]interface{}) {
			imports = append(imports, *sdk.NewFunctionImportRequest().WithImport(item.(string)))
		}
		request.WithImports(imports)
	}
	if _, ok := d.GetOk("packages"); ok {
		packages := []sdk.FunctionPackageRequest{}
		for _, item := range d.Get("packages").([]interface{}) {
			packages = append(packages, *sdk.NewFunctionPackageRequest().WithPackage(item.(string)))
		}
		request.WithPackages(packages)
	}
	if v, ok := d.GetOk("target_path"); ok {
		request.WithTargetPath(v.(string))
	}

	if err := client.Functions.CreateForJava(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	argumentTypes := make([]sdk.DataType, 0, len(arguments))
	for _, item := range arguments {
		argumentTypes = append(argumentTypes, item.ArgDataType)
	}
	nid := sdk.NewSchemaObjectIdentifierWithArguments(database, schema, name, argumentTypes...)
	d.SetId(nid.FullyQualifiedName())
	return ReadContextFunction(ctx, d, meta)
}

func createScalaFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schema, name)

	// Set required
	returnType := d.Get("return_type").(string)
	returnDataType, diags := convertFunctionDataType(returnType)
	if diags != nil {
		return diags
	}
	functionDefinition := d.Get("statement").(string)
	handler := d.Get("handler").(string)
	// create request with required
	request := sdk.NewCreateForScalaFunctionRequest(id, returnDataType, handler)
	request.WithFunctionDefinition(functionDefinition)

	// Set optionals
	if v, ok := d.GetOk("is_secure"); ok {
		request.WithSecure(v.(bool))
	}
	arguments, diags := parseFunctionArguments(d)
	if diags != nil {
		return diags
	}
	if len(arguments) > 0 {
		request.WithArguments(arguments)
	}
	if v, ok := d.GetOk("null_input_behavior"); ok {
		request.WithNullInputBehavior(sdk.NullInputBehavior(v.(string)))
	}
	if v, ok := d.GetOk("return_behavior"); ok {
		request.WithReturnResultsBehavior(sdk.ReturnResultsBehavior(v.(string)))
	}
	if v, ok := d.GetOk("runtime_version"); ok {
		request.WithRuntimeVersion(v.(string))
	}
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}
	if _, ok := d.GetOk("imports"); ok {
		imports := []sdk.FunctionImportRequest{}
		for _, item := range d.Get("imports").([]interface{}) {
			imports = append(imports, *sdk.NewFunctionImportRequest().WithImport(item.(string)))
		}
		request.WithImports(imports)
	}
	if _, ok := d.GetOk("packages"); ok {
		packages := []sdk.FunctionPackageRequest{}
		for _, item := range d.Get("packages").([]interface{}) {
			packages = append(packages, *sdk.NewFunctionPackageRequest().WithPackage(item.(string)))
		}
		request.WithPackages(packages)
	}
	if v, ok := d.GetOk("target_path"); ok {
		request.WithTargetPath(v.(string))
	}

	if err := client.Functions.CreateForScala(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	argumentTypes := make([]sdk.DataType, 0, len(arguments))
	for _, item := range arguments {
		argumentTypes = append(argumentTypes, item.ArgDataType)
	}
	nid := sdk.NewSchemaObjectIdentifierWithArguments(database, schema, name, argumentTypes...)
	d.SetId(nid.FullyQualifiedName())
	return ReadContextFunction(ctx, d, meta)
}

func createSQLFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schema, name)

	// Set required
	returns, diags := parseFunctionReturnsRequest(d.Get("return_type").(string))
	if diags != nil {
		return diags
	}
	functionDefinition := d.Get("statement").(string)
	// create request with required
	request := sdk.NewCreateForSQLFunctionRequest(id, *returns, functionDefinition)

	// Set optionals
	if v, ok := d.GetOk("is_secure"); ok {
		request.WithSecure(v.(bool))
	}
	arguments, diags := parseFunctionArguments(d)
	if diags != nil {
		return diags
	}
	if len(arguments) > 0 {
		request.WithArguments(arguments)
	}
	if v, ok := d.GetOk("return_behavior"); ok {
		request.WithReturnResultsBehavior(sdk.ReturnResultsBehavior(v.(string)))
	}
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	if err := client.Functions.CreateForSQL(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	argumentTypes := make([]sdk.DataType, 0, len(arguments))
	for _, item := range arguments {
		argumentTypes = append(argumentTypes, item.ArgDataType)
	}
	nid := sdk.NewSchemaObjectIdentifierWithArguments(database, schema, name, argumentTypes...)
	d.SetId(nid.FullyQualifiedName())
	return ReadContextFunction(ctx, d, meta)
}

func createPythonFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schema, name)

	// Set required
	returns, diags := parseFunctionReturnsRequest(d.Get("return_type").(string))
	if diags != nil {
		return diags
	}
	functionDefinition := d.Get("statement").(string)
	version := d.Get("runtime_version").(string)
	handler := d.Get("handler").(string)
	// create request with required
	request := sdk.NewCreateForPythonFunctionRequest(id, *returns, version, handler)
	request.WithFunctionDefinition(functionDefinition)

	// Set optionals
	if v, ok := d.GetOk("is_secure"); ok {
		request.WithSecure(v.(bool))
	}
	arguments, diags := parseFunctionArguments(d)
	if diags != nil {
		return diags
	}
	if len(arguments) > 0 {
		request.WithArguments(arguments)
	}
	if v, ok := d.GetOk("null_input_behavior"); ok {
		request.WithNullInputBehavior(sdk.NullInputBehavior(v.(string)))
	}
	if v, ok := d.GetOk("return_behavior"); ok {
		request.WithReturnResultsBehavior(sdk.ReturnResultsBehavior(v.(string)))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}
	if _, ok := d.GetOk("imports"); ok {
		imports := []sdk.FunctionImportRequest{}
		for _, item := range d.Get("imports").([]interface{}) {
			imports = append(imports, *sdk.NewFunctionImportRequest().WithImport(item.(string)))
		}
		request.WithImports(imports)
	}
	if _, ok := d.GetOk("packages"); ok {
		packages := []sdk.FunctionPackageRequest{}
		for _, item := range d.Get("packages").([]interface{}) {
			packages = append(packages, *sdk.NewFunctionPackageRequest().WithPackage(item.(string)))
		}
		request.WithPackages(packages)
	}

	if err := client.Functions.CreateForPython(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	argumentTypes := make([]sdk.DataType, 0, len(arguments))
	for _, item := range arguments {
		argumentTypes = append(argumentTypes, item.ArgDataType)
	}
	nid := sdk.NewSchemaObjectIdentifierWithArguments(database, schema, name, argumentTypes...)
	d.SetId(nid.FullyQualifiedName())
	return ReadContextFunction(ctx, d, meta)
}

func createJavascriptFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schema, name)

	// Set required
	returns, diags := parseFunctionReturnsRequest(d.Get("return_type").(string))
	if diags != nil {
		return diags
	}
	functionDefinition := d.Get("statement").(string)
	// create request with required
	request := sdk.NewCreateForJavascriptFunctionRequest(id, *returns, functionDefinition)

	// Set optionals
	if v, ok := d.GetOk("is_secure"); ok {
		request.WithSecure(v.(bool))
	}
	arguments, diags := parseFunctionArguments(d)
	if diags != nil {
		return diags
	}
	if len(arguments) > 0 {
		request.WithArguments(arguments)
	}
	if v, ok := d.GetOk("null_input_behavior"); ok {
		request.WithNullInputBehavior(sdk.NullInputBehavior(v.(string)))
	}
	if v, ok := d.GetOk("return_behavior"); ok {
		request.WithReturnResultsBehavior(sdk.ReturnResultsBehavior(v.(string)))
	}
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	if err := client.Functions.CreateForJavascript(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	argumentTypes := make([]sdk.DataType, 0, len(arguments))
	for _, item := range arguments {
		argumentTypes = append(argumentTypes, item.ArgDataType)
	}
	nid := sdk.NewSchemaObjectIdentifierWithArguments(database, schema, name, argumentTypes...)
	// TODO: Create upgrader for id migration
	d.SetId(nid.FullyQualifiedName())
	return ReadContextFunction(ctx, d, meta)
}

func ReadContextFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*provider.Context).Client

	id, err := sdk.NewSchemaObjectIdentifierWithArgumentsFromFullyQualifiedName(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", id.Name()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("database", id.DatabaseName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", id.SchemaName()); err != nil {
		return diag.FromErr(err)
	}

	arguments := d.Get("arguments").([]interface{})
	argumentTypes := make([]string, len(arguments))
	for i, arg := range arguments {
		argumentTypes[i] = arg.(map[string]interface{})["type"].(string)
	}
	functionDetails, err := client.Functions.Describe(ctx, id)
	if err != nil {
		// if function is not found then mark resource to be removed from state file during apply or refresh
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Describe function failed.",
				Detail:   "See our document on design decisions for functions: <LINK (coming soon)>",
			},
		}
	}
	for _, desc := range functionDetails {
		switch desc.Property {
		case "signature":
			// Format in Snowflake DB is: (argName argType, argName argType, ...)
			value := strings.ReplaceAll(strings.ReplaceAll(desc.Value, "(", ""), ")", "")
			if value != "" { // Do nothing for functions without arguments
				pairs := strings.Split(value, ", ")

				arguments := []interface{}{}
				for _, pair := range pairs {
					item := strings.Split(pair, " ")
					argument := map[string]interface{}{}
					argument["name"] = item[0]
					argument["type"] = item[1]
					arguments = append(arguments, argument)
				}
				if err := d.Set("arguments", arguments); err != nil {
					diag.FromErr(err)
				}
			}
		case "null handling":
			if err := d.Set("null_input_behavior", desc.Value); err != nil {
				diag.FromErr(err)
			}
		case "volatility":
			if err := d.Set("return_behavior", desc.Value); err != nil {
				diag.FromErr(err)
			}
		case "body":
			if err := d.Set("statement", desc.Value); err != nil {
				diag.FromErr(err)
			}
		case "returns":
			// Format in Snowflake DB is returnType(<some number>)
			re := regexp.MustCompile(`^(.*)\([0-9]*\)$`)
			match := re.FindStringSubmatch(desc.Value)
			rt := desc.Value
			if match != nil {
				rt = match[1]
			}
			if err := d.Set("return_type", rt); err != nil {
				diag.FromErr(err)
			}
		case "language":
			if snowflake.Contains(languages, strings.ToLower(desc.Value)) {
				if err := d.Set("language", desc.Value); err != nil {
					diag.FromErr(err)
				}
			} else {
				log.Printf("[INFO] Unexpected language for function %v returned from Snowflake", desc.Value)
			}
		case "packages":
			value := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(desc.Value, "[", ""), "]", ""), "'", "")
			if value != "" { // Do nothing for Java / Python functions without packages
				packages := strings.Split(value, ",")
				if err := d.Set("packages", packages); err != nil {
					diag.FromErr(err)
				}
			}
		case "imports":
			value := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(desc.Value, "[", ""), "]", ""), "'", "")
			if value != "" { // Do nothing for Java functions without imports
				imports := strings.Split(value, ",")
				if err := d.Set("imports", imports); err != nil {
					diag.FromErr(err)
				}
			}
		case "handler":
			if err := d.Set("handler", desc.Value); err != nil {
				diag.FromErr(err)
			}
		case "target_path":
			if err := d.Set("target_path", desc.Value); err != nil {
				diag.FromErr(err)
			}
		case "runtime_version":
			if err := d.Set("runtime_version", desc.Value); err != nil {
				diag.FromErr(err)
			}
		default:
			log.Printf("[INFO] Unexpected function property %v returned from Snowflake with value %v", desc.Property, desc.Value)
		}
	}

	function, err := client.Functions.ShowByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_secure", function.IsSecure); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", function.Description); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func UpdateContextFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.NewSchemaObjectIdentifierWithArgumentsFromFullyQualifiedName(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		newId := sdk.NewSchemaObjectIdentifierWithArguments(id.DatabaseName(), id.SchemaName(), name, id.ArgumentDataTypes()...)

		if err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithRenameTo(newId.SchemaObjectId())); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(newId.FullyQualifiedName())
		id = newId
	}

	if d.HasChange("is_secure") {
		secure := d.Get("is_secure")
		if secure.(bool) {
			if err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetSecure(true)); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetSecure(true)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("comment") {
		comment := d.Get("comment")
		if comment != "" {
			if err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetComment(comment.(string))); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetComment(true)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return ReadContextFunction(ctx, d, meta)
}

func DeleteContextFunction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.NewSchemaObjectIdentifierWithArgumentsFromFullyQualifiedName(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func parseFunctionArguments(d *schema.ResourceData) ([]sdk.FunctionArgumentRequest, diag.Diagnostics) {
	args := make([]sdk.FunctionArgumentRequest, 0)
	if v, ok := d.GetOk("arguments"); ok {
		for _, arg := range v.([]interface{}) {
			argName := arg.(map[string]interface{})["name"].(string)
			argType := arg.(map[string]interface{})["type"].(string)
			argDataType, diags := convertFunctionDataType(argType)
			if diags != nil {
				return nil, diags
			}
			args = append(args, sdk.FunctionArgumentRequest{ArgName: argName, ArgDataType: argDataType})
		}
	}
	return args, nil
}

func convertFunctionDataType(s string) (sdk.DataType, diag.Diagnostics) {
	dataType, err := sdk.ToDataType(s)
	if err != nil {
		return dataType, diag.FromErr(err)
	}
	return dataType, nil
}

func convertFunctionColumns(s string) ([]sdk.FunctionColumn, diag.Diagnostics) {
	pattern := regexp.MustCompile(`(\w+)\s+(\w+)`)
	matches := pattern.FindAllStringSubmatch(s, -1)
	var columns []sdk.FunctionColumn
	for _, match := range matches {
		if len(match) == 3 {
			dataType, err := sdk.ToDataType(match[2])
			if err != nil {
				return nil, diag.FromErr(err)
			}
			columns = append(columns, sdk.FunctionColumn{
				ColumnName:     match[1],
				ColumnDataType: dataType,
			})
		}
	}
	return columns, nil
}

func parseFunctionReturnsRequest(s string) (*sdk.FunctionReturnsRequest, diag.Diagnostics) {
	returns := sdk.NewFunctionReturnsRequest()
	if strings.HasPrefix(strings.ToLower(s), "table") {
		columns, diags := convertFunctionColumns(s)
		if diags != nil {
			return nil, diags
		}
		var cr []sdk.FunctionColumnRequest
		for _, item := range columns {
			cr = append(cr, *sdk.NewFunctionColumnRequest(item.ColumnName, item.ColumnDataType))
		}
		returns.WithTable(*sdk.NewFunctionReturnsTableRequest().WithColumns(cr))
	} else {
		returnDataType, diags := convertFunctionDataType(s)
		if diags != nil {
			return nil, diags
		}
		returns.WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(returnDataType))
	}
	return returns, nil
}
