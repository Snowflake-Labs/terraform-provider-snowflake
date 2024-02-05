package resources

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var procedureSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the procedure; does not have to be unique for the schema in which the procedure is created. Don't use the | character.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the procedure. Don't use the | character.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the procedure. Don't use the | character.",
		ForceNew:    true,
	},
	"secure": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies that the procedure is secure. For more information about secure procedures, see Protecting Sensitive Information with Secure UDFs and Stored Procedures.",
		Default:     false,
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
					// todo: add validation that this is a valid data type
					Description: "The argument type",
				},
			},
		},
		Optional:    true,
		Description: "List of the arguments for the procedure",
		ForceNew:    true,
	},
	"return_type": {
		Type:        schema.TypeString,
		Description: "The return type of the procedure",
		// Suppress the diff shown if the values are equal when both compared in lower case.
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			if strings.EqualFold(old, new) {
				return true
			}

			varcharType := []string{"VARCHAR(16777216)", "VARCHAR", "text", "string", "NVARCHAR", "NVARCHAR2", "CHAR VARYING", "NCHAR VARYING"}
			if slices.Contains(varcharType, strings.ToUpper(old)) && slices.Contains(varcharType, strings.ToUpper(new)) {
				return true
			}

			// all these types are equivalent https://docs.snowflake.com/en/sql-reference/data-types-numeric.html#int-integer-bigint-smallint-tinyint-byteint
			integerTypes := []string{"INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT", "BYTEINT", "NUMBER(38,0)"}
			if slices.Contains(integerTypes, strings.ToUpper(old)) && slices.Contains(integerTypes, strings.ToUpper(new)) {
				return true
			}
			return false
		},
		// todo: add validation that this is a valid data type
		Required: true,
		ForceNew: true,
	},
	"statement": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the code used to create the procedure.",
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
		ValidateFunc: validation.StringInSlice([]string{"javascript", "java", "scala", "SQL", "python"}, true),
		Description:  "Specifies the language of the stored procedure code.",
	},
	"execute_as": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "OWNER",
		Description: "Sets execute context - see caller's rights and owner's rights",
	},
	"null_input_behavior": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "CALLED ON NULL INPUT",
		ForceNew: true,
		// We do not use STRICT, because Snowflake then in the Read phase returns RETURNS NULL ON NULL INPUT
		ValidateFunc: validation.StringInSlice([]string{"CALLED ON NULL INPUT", "RETURNS NULL ON NULL INPUT"}, false),
		Description:  "Specifies the behavior of the procedure when called with null inputs.",
	},
	"return_behavior": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "VOLATILE",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"VOLATILE", "IMMUTABLE"}, false),
		Description:  "Specifies the behavior of the function when returning results",
		Deprecated:   "These keywords are deprecated for stored procedures. These keywords are not intended to apply to stored procedures. In a future release, these keywords will be removed from the documentation.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "user-defined procedure",
		Description: "Specifies a comment for the procedure.",
	},
	"runtime_version": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Required for Python procedures. Specifies Python runtime version.",
	},
	"packages": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		ForceNew:    true,
		Description: "List of package imports to use for Java / Python procedures. For Java, package imports should be of the form: package_name:version_number, where package_name is snowflake_domain:package. For Python use it should be: ('numpy','pandas','xgboost==1.5.0').",
	},
	"imports": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		ForceNew:    true,
		Description: "Imports for Java / Python procedures. For Java this a list of jar files, for Python this is a list of Python files.",
	},
	"handler": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "The handler method for Java / Python procedures.",
	},
}

// Procedure returns a pointer to the resource representing a stored procedure.
func Procedure() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextProcedure,
		ReadContext:   ReadContextProcedure,
		UpdateContext: UpdateContextProcedure,
		DeleteContext: DeleteContextProcedure,

		Schema: procedureSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextProcedure(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lang := strings.ToUpper(d.Get("language").(string))
	switch lang {
	case "JAVA":
		return createJavaProcedure(ctx, d, meta)
	case "JAVASCRIPT":
		return createJavaScriptProcedure(ctx, d, meta)
	case "PYTHON":
		return createPythonProcedure(ctx, d, meta)
	case "SCALA":
		return createScalaProcedure(ctx, d, meta)
	case "SQL":
		return createSQLProcedure(ctx, d, meta)
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

func createJavaProcedure(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	schemaObjectId := sdk.NewSchemaObjectIdentifier(database, schema, name)
	returnType := d.Get("return_type").(string)
	returnDataType, diags := convertProcedureDataType(returnType)
	if diags != nil {
		return diags
	}
	procedureDefinition := d.Get("statement").(string)

	returns := sdk.NewProcedureReturnsRequest().WithResultDataType(sdk.NewProcedureReturnsResultDataTypeRequest(returnDataType))
	runtimeVersion := d.Get("runtime_version").(string)

	packages := []sdk.ProcedurePackageRequest{}
	pkgs := d.Get("packages").([]interface{})
	for _, pkg := range pkgs {
		packages = append(packages, *sdk.NewProcedurePackageRequest(pkg.(string)))
	}
	handler := d.Get("handler").(string)
	req := sdk.NewCreateForJavaProcedureRequest(schemaObjectId, *returns, runtimeVersion, packages, handler)
	req.WithProcedureDefinition(sdk.String(procedureDefinition))
	args, diags := getProcedureArguments(d)
	if diags != nil {
		return diags
	}
	if len(args) > 0 {
		req.WithArguments(args)
	}

	// read optional params
	if v, ok := d.GetOk("execute_as"); ok {
		if strings.ToUpper(v.(string)) == "OWNER" {
			req.WithExecuteAs(sdk.Pointer(sdk.ExecuteAsOwner))
		} else if strings.ToUpper(v.(string)) == "CALLER" {
			req.WithExecuteAs(sdk.Pointer(sdk.ExecuteAsCaller))
		}
	}
	if v, ok := d.GetOk("null_input_behavior"); ok {
		req.WithNullInputBehavior(sdk.Pointer(sdk.NullInputBehavior(v.(string))))
	}
	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(sdk.String(v.(string)))
	}
	if v, ok := d.GetOk("secure"); ok {
		req.WithSecure(sdk.Bool(v.(bool)))
	}
	if _, ok := d.GetOk("imports"); ok {
		imports := []sdk.ProcedureImportRequest{}
		for _, imp := range d.Get("imports").([]interface{}) {
			imports = append(imports, *sdk.NewProcedureImportRequest(imp.(string)))
		}
		req.WithImports(imports)
	}

	if err := client.Procedures.CreateForJava(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	argTypes := make([]string, len(args))
	for i, arg := range args {
		argTypes[i] = string(arg.ArgDataType)
	}
	id := NewProcedureID(database, schema, name, argTypes)
	d.SetId(id.String())
	return ReadContextProcedure(ctx, d, meta)
}

func createJavaScriptProcedure(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	schemaObjectId := sdk.NewSchemaObjectIdentifier(database, schema, name)
	returnType := d.Get("return_type").(string)
	returnDataType, diags := convertProcedureDataType(returnType)
	if diags != nil {
		return diags
	}
	procedureDefinition := d.Get("statement").(string)

	req := sdk.NewCreateForJavaScriptProcedureRequest(schemaObjectId, returnDataType, procedureDefinition)
	args, diags := getProcedureArguments(d)
	if diags != nil {
		return diags
	}
	if len(args) > 0 {
		req.WithArguments(args)
	}

	// read optional params
	if v, ok := d.GetOk("execute_as"); ok {
		if strings.ToUpper(v.(string)) == "OWNER" {
			req.WithExecuteAs(sdk.Pointer(sdk.ExecuteAsOwner))
		} else if strings.ToUpper(v.(string)) == "CALLER" {
			req.WithExecuteAs(sdk.Pointer(sdk.ExecuteAsCaller))
		}
	}
	if v, ok := d.GetOk("null_input_behavior"); ok {
		req.WithNullInputBehavior(sdk.Pointer(sdk.NullInputBehavior(v.(string))))
	}
	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(sdk.String(v.(string)))
	}
	if v, ok := d.GetOk("secure"); ok {
		req.WithSecure(sdk.Bool(v.(bool)))
	}

	if err := client.Procedures.CreateForJavaScript(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	argTypes := make([]string, len(args))
	for i, arg := range args {
		argTypes[i] = string(arg.ArgDataType)
	}
	id := NewProcedureID(database, schema, name, argTypes)
	d.SetId(id.String())
	return ReadContextProcedure(ctx, d, meta)
}

func createScalaProcedure(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	schemaObjectId := sdk.NewSchemaObjectIdentifier(database, schema, name)
	returnType := d.Get("return_type").(string)
	returnDataType, diags := convertProcedureDataType(returnType)
	if diags != nil {
		return diags
	}
	procedureDefinition := d.Get("statement").(string)

	returns := sdk.NewProcedureReturnsRequest().WithResultDataType(sdk.NewProcedureReturnsResultDataTypeRequest(returnDataType))
	runtimeVersion := d.Get("runtime_version").(string)

	packages := []sdk.ProcedurePackageRequest{}
	pkgs := d.Get("packages").([]interface{})
	for _, pkg := range pkgs {
		packages = append(packages, *sdk.NewProcedurePackageRequest(pkg.(string)))
	}
	handler := d.Get("handler").(string)
	req := sdk.NewCreateForScalaProcedureRequest(schemaObjectId, *returns, runtimeVersion, packages, handler)
	req.WithProcedureDefinition(sdk.String(procedureDefinition))
	args, diags := getProcedureArguments(d)
	if diags != nil {
		return diags
	}
	if len(args) > 0 {
		req.WithArguments(args)
	}

	// read optional params
	if v, ok := d.GetOk("execute_as"); ok {
		if strings.ToUpper(v.(string)) == "OWNER" {
			req.WithExecuteAs(sdk.Pointer(sdk.ExecuteAsOwner))
		} else if strings.ToUpper(v.(string)) == "CALLER" {
			req.WithExecuteAs(sdk.Pointer(sdk.ExecuteAsCaller))
		}
	}
	if v, ok := d.GetOk("null_input_behavior"); ok {
		req.WithNullInputBehavior(sdk.Pointer(sdk.NullInputBehavior(v.(string))))
	}
	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(sdk.String(v.(string)))
	}
	if v, ok := d.GetOk("secure"); ok {
		req.WithSecure(sdk.Bool(v.(bool)))
	}
	if _, ok := d.GetOk("imports"); ok {
		imports := []sdk.ProcedureImportRequest{}
		for _, imp := range d.Get("imports").([]interface{}) {
			imports = append(imports, *sdk.NewProcedureImportRequest(imp.(string)))
		}
		req.WithImports(imports)
	}

	if err := client.Procedures.CreateForScala(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	argTypes := make([]string, len(args))
	for i, arg := range args {
		argTypes[i] = string(arg.ArgDataType)
	}
	id := NewProcedureID(database, schema, name, argTypes)
	d.SetId(id.String())
	return ReadContextProcedure(ctx, d, meta)
}

func createSQLProcedure(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	schemaObjectId := sdk.NewSchemaObjectIdentifier(database, schema, name)
	returnType := d.Get("return_type").(string)
	returnDataType, diags := convertProcedureDataType(returnType)
	if diags != nil {
		return diags
	}
	procedureDefinition := d.Get("statement").(string)

	returns := sdk.NewProcedureSQLReturnsRequest().WithResultDataType(sdk.NewProcedureReturnsResultDataTypeRequest(returnDataType))
	req := sdk.NewCreateForSQLProcedureRequest(schemaObjectId, *returns, procedureDefinition)
	args, diags := getProcedureArguments(d)
	if diags != nil {
		return diags
	}
	if len(args) > 0 {
		req.WithArguments(args)
	}

	// read optional params
	if v, ok := d.GetOk("execute_as"); ok {
		if strings.ToUpper(v.(string)) == "OWNER" {
			req.WithExecuteAs(sdk.Pointer(sdk.ExecuteAsOwner))
		} else if strings.ToUpper(v.(string)) == "CALLER" {
			req.WithExecuteAs(sdk.Pointer(sdk.ExecuteAsCaller))
		}
	}
	if v, ok := d.GetOk("null_input_behavior"); ok {
		req.WithNullInputBehavior(sdk.Pointer(sdk.NullInputBehavior(v.(string))))
	}
	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(sdk.String(v.(string)))
	}
	if v, ok := d.GetOk("secure"); ok {
		req.WithSecure(sdk.Bool(v.(bool)))
	}

	if err := client.Procedures.CreateForSQL(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	argTypes := make([]string, len(args))
	for i, arg := range args {
		argTypes[i] = string(arg.ArgDataType)
	}
	id := NewProcedureID(database, schema, name, argTypes)
	d.SetId(id.String())
	return ReadContextProcedure(ctx, d, meta)
}

func createPythonProcedure(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	schemaObjectId := sdk.NewSchemaObjectIdentifier(database, schema, name)
	returns := sdk.NewProcedureReturnsRequest()

	returnType := d.Get("return_type").(string)
	if strings.HasPrefix(strings.ToLower(returnType), "table") {
		columns, diags := convertProcedureColumns(returnType)
		if diags != nil {
			return diags
		}
		var cr []sdk.ProcedureColumnRequest
		for _, item := range columns {
			cr = append(cr, *sdk.NewProcedureColumnRequest(item.ColumnName, item.ColumnDataType))
		}
		returns.WithTable(sdk.NewProcedureReturnsTableRequest().WithColumns(cr))
	} else {
		returnDataType, diags := convertProcedureDataType(returnType)
		if diags != nil {
			return diags
		}
		returns.WithResultDataType(sdk.NewProcedureReturnsResultDataTypeRequest(returnDataType))
	}

	procedureDefinition := d.Get("statement").(string)
	runtimeVersion := d.Get("runtime_version").(string)
	packages := []sdk.ProcedurePackageRequest{}
	pkgs := d.Get("packages").([]interface{})
	for _, pkg := range pkgs {
		packages = append(packages, *sdk.NewProcedurePackageRequest(pkg.(string)))
	}
	handler := d.Get("handler").(string)
	req := sdk.NewCreateForPythonProcedureRequest(schemaObjectId, *returns, runtimeVersion, packages, handler)
	req.WithProcedureDefinition(sdk.String(procedureDefinition))
	args, diags := getProcedureArguments(d)
	if diags != nil {
		return diags
	}
	if len(args) > 0 {
		req.WithArguments(args)
	}

	// read optional params
	if v, ok := d.GetOk("execute_as"); ok {
		if strings.ToUpper(v.(string)) == "OWNER" {
			req.WithExecuteAs(sdk.Pointer(sdk.ExecuteAsOwner))
		} else if strings.ToUpper(v.(string)) == "CALLER" {
			req.WithExecuteAs(sdk.Pointer(sdk.ExecuteAsCaller))
		}
	}

	// TODO: [ { CALLED ON NULL INPUT | { RETURNS NULL ON NULL INPUT | STRICT } } ] not works for python
	// if v, ok := d.GetOk("null_input_behavior"); ok {
	// 	req.WithNullInputBehavior(sdk.Pointer(sdk.NullInputBehavior(v.(string))))
	// }

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(sdk.String(v.(string)))
	}
	if v, ok := d.GetOk("secure"); ok {
		req.WithSecure(sdk.Bool(v.(bool)))
	}
	if _, ok := d.GetOk("imports"); ok {
		imports := []sdk.ProcedureImportRequest{}
		for _, imp := range d.Get("imports").([]interface{}) {
			imports = append(imports, *sdk.NewProcedureImportRequest(imp.(string)))
		}
		req.WithImports(imports)
	}

	if err := client.Procedures.CreateForPython(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	argTypes := make([]string, len(args))
	for i, arg := range args {
		argTypes[i] = string(arg.ArgDataType)
	}
	id := NewProcedureID(database, schema, name, argTypes)
	d.SetId(id.String())
	return ReadContextProcedure(ctx, d, meta)
}

func ReadContextProcedure(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	id, err := DecodeProcedureID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", id.name.Name()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("database", id.DatabaseName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", id.SchemaName()); err != nil {
		return diag.FromErr(err)
	}
	args := d.Get("arguments").([]interface{})
	argTypes := make([]string, len(args))
	for i, arg := range args {
		argTypes[i] = arg.(map[string]interface{})["type"].(string)
	}
	procedureDetails, err := client.Procedures.Describe(ctx, sdk.NewDescribeProcedureRequest(id.name, id.argTypes))
	if err != nil {
		// if procedure is not found then mark resource to be removed from state file during apply or refresh
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Describe procedure failed.",
				// TODO: link to the design decisions doc
				Detail: "See our document on design decisions for procedures: <LINK (coming soon)>",
			},
		}
	}
	for _, desc := range procedureDetails {
		switch desc.Property {
		case "signature":
			// Format in Snowflake DB is: (argName argType, argName argType, ...)
			args := strings.ReplaceAll(strings.ReplaceAll(desc.Value, "(", ""), ")", "")

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

				if err := d.Set("arguments", args); err != nil {
					return diag.FromErr(err)
				}
			}
		case "null handling":
			if err := d.Set("null_input_behavior", desc.Value); err != nil {
				return diag.FromErr(err)
			}
		case "volatility":
			if err := d.Set("return_behavior", desc.Value); err != nil {
				return diag.FromErr(err)
			}
		case "body":
			if err := d.Set("statement", desc.Value); err != nil {
				return diag.FromErr(err)
			}
		case "execute as":
			if err := d.Set("execute_as", desc.Value); err != nil {
				return diag.FromErr(err)
			}
		case "returns":
			if err := d.Set("return_type", desc.Value); err != nil {
				return diag.FromErr(err)
			}
		case "language":
			if err := d.Set("language", desc.Value); err != nil {
				return diag.FromErr(err)
			}
		case "runtime_version":
			if err := d.Set("runtime_version", desc.Value); err != nil {
				return diag.FromErr(err)
			}
		case "packages":
			packagesString := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(desc.Value, "[", ""), "]", ""), "'", "")
			if packagesString != "" { // Do nothing for Java / Python functions without packages
				packages := strings.Split(packagesString, ",")
				if err := d.Set("packages", packages); err != nil {
					return diag.FromErr(err)
				}
			}
		case "imports":
			importsString := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(desc.Value, "[", ""), "]", ""), "'", ""), " ", "")
			if importsString != "" { // Do nothing for Java functions without imports
				imports := strings.Split(importsString, ",")
				if err := d.Set("imports", imports); err != nil {
					return diag.FromErr(err)
				}
			}
		case "handler":
			if err := d.Set("handler", desc.Value); err != nil {
				return diag.FromErr(err)
			}
		default:
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Unexpected procedure property returned from Snowflake",
				Detail:   fmt.Sprintf("Unexpected procedure property %v returned from Snowflake", desc.Property),
			})
		}
	}

	request := sdk.NewShowProcedureRequest().WithIn(&sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())}).WithLike(&sdk.Like{Pattern: sdk.String(id.name.Name())})

	procedures, err := client.Procedures.Show(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}
	// procedure names can be overloaded with different argument types so we iterate over and find the correct one
	// the ShowByID function should probably be updated to also require the list of arg types, like describe procedure
	for _, procedure := range procedures {
		argumentSignature := strings.Split(procedure.Arguments, " RETURN ")[0]
		argumentSignature = strings.ReplaceAll(argumentSignature, " ", "")
		if argumentSignature == id.ArgumentSignature() {
			if err := d.Set("secure", procedure.IsSecure); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("comment", procedure.Description); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return diags
}

func UpdateContextProcedure(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	id, err := DecodeProcedureID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("name") {
		_, newName := d.GetChange("name")
		err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id.name, id.argTypes).WithRenameTo(sdk.Pointer(sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), newName.(string)))))
		if err != nil {
			return diag.FromErr(err)
		}
		id.name = sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), newName.(string))
		if err := d.Set("name", newName); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("comment") {
		_, newComment := d.GetChange("comment")
		if newComment != "" {
			err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id.name, id.argTypes).WithSetComment(sdk.String(newComment.(string))))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id.name, id.argTypes).WithUnsetComment(sdk.Bool(true)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("execute_as") {
		_, newExecuteAs := d.GetChange("execute_as")
		err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id.name, id.argTypes).WithExecuteAs(sdk.Pointer(sdk.ExecuteAs(newExecuteAs.(string)))))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextProcedure(ctx, d, meta)
}

func DeleteContextProcedure(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	id, err := DecodeProcedureID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.Procedures.Drop(ctx, sdk.NewDropProcedureRequest(id.name, id.argTypes))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func NewProcedureID(database, schema, name string, argTypes []string) *procedureID {
	argDataTypes := make([]sdk.DataType, len(argTypes))
	for i, argType := range argTypes {
		argDataTypes[i], _ = sdk.ToDataType(argType)
	}
	return &procedureID{
		name:     sdk.NewSchemaObjectIdentifier(database, schema, name),
		argTypes: argDataTypes,
	}
}

func DecodeProcedureID(id string) (*procedureID, error) {
	parts := strings.Split(id, "|")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid procedure id %v", id)
	}
	argTypes := strings.Split(parts[3], ",")
	return NewProcedureID(parts[0], parts[1], parts[2], argTypes), nil
}

type procedureID struct {
	name     sdk.SchemaObjectIdentifier
	argTypes []sdk.DataType
}

func (i procedureID) DatabaseName() string {
	return i.name.DatabaseName()
}

func (i procedureID) SchemaName() string {
	return i.name.SchemaName()
}

func (i procedureID) ProcedureName() string {
	return i.name.Name()
}

func (i procedureID) ArgTypes() []sdk.DataType {
	return i.argTypes
}

func (i procedureID) ArgumentSignature() string {
	argTypes := make([]string, len(i.argTypes))
	for i, argType := range i.argTypes {
		argTypes[i] = string(argType)
	}
	return fmt.Sprintf("%v(%v)", i.name.Name(), strings.Join(argTypes, ","))
}

func (i procedureID) String() string {
	argTypes := make([]string, len(i.argTypes))
	for i, argType := range i.argTypes {
		argTypes[i] = string(argType)
	}
	return fmt.Sprintf("%v|%v|%v|%v", i.DatabaseName(), i.SchemaName(), i.ProcedureName(), strings.Join(argTypes, ","))
}

func getProcedureArguments(d *schema.ResourceData) ([]sdk.ProcedureArgumentRequest, diag.Diagnostics) {
	args := make([]sdk.ProcedureArgumentRequest, 0)
	if v, ok := d.GetOk("arguments"); ok {
		for _, arg := range v.([]interface{}) {
			argName := arg.(map[string]interface{})["name"].(string)
			argType := arg.(map[string]interface{})["type"].(string)
			argDataType, diags := convertProcedureDataType(argType)
			if diags != nil {
				return nil, diags
			}
			args = append(args, sdk.ProcedureArgumentRequest{ArgName: argName, ArgDataType: argDataType})
		}
	}
	return args, nil
}

func convertProcedureDataType(s string) (sdk.DataType, diag.Diagnostics) {
	dataType, err := sdk.ToDataType(s)
	if err != nil {
		return dataType, diag.FromErr(err)
	}
	return dataType, nil
}

func convertProcedureColumns(s string) ([]sdk.ProcedureColumn, diag.Diagnostics) {
	pattern := regexp.MustCompile(`(\w+)\s+(\w+)`)
	matches := pattern.FindAllStringSubmatch(s, -1)
	var columns []sdk.ProcedureColumn
	for _, match := range matches {
		if len(match) == 3 {
			dataType, err := sdk.ToDataType(match[2])
			if err != nil {
				return nil, diag.FromErr(err)
			}
			columns = append(columns, sdk.ProcedureColumn{
				ColumnName:     match[1],
				ColumnDataType: dataType,
			})
		}
	}
	return columns, nil
}
