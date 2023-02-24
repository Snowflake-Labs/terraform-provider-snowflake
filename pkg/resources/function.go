package resources

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var languages = []string{"javascript", "java", "sql", "python"}

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
					DiffSuppressFunc: DiffTypes,
					Description:      "The argument name",
				},
				"type": {
					Type:     schema.TypeString,
					Required: true,
					// Suppress the diff shown if the values are equal when both compared in lower case.
					DiffSuppressFunc: DiffTypes,
					Description:      "The argument type",
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
		DiffSuppressFunc: DiffTypes,
		Required:         true,
		ForceNew:         true,
	},
	"statement": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the javascript / java / sql / python code used to create the function.",
		ForceNew:         true,
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"language": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice(languages, false),
		Description:  "The language of the statement",
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

// Function returns a pointer to the resource representing a stored function.
func Function() *schema.Resource {
	return &schema.Resource{
		Create: CreateFunction,
		Read:   ReadFunction,
		Update: UpdateFunction,
		Delete: DeleteFunction,

		Schema: functionSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateFunction implements schema.CreateFunc.
func CreateFunction(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	s := d.Get("statement").(string)
	ret := d.Get("return_type").(string)

	builder := snowflake.NewFunctionBuilder(database, schema, name, []string{}).WithStatement(s).WithReturnType(ret)

	// Set optionals, args
	if _, ok := d.GetOk("arguments"); ok {
		args := []map[string]string{}
		for _, arg := range d.Get("arguments").([]interface{}) {
			argDef := map[string]string{}
			for key, val := range arg.(map[string]interface{}) {
				argDef[key] = val.(string)
			}
			args = append(args, argDef)
		}
		builder.WithArgs(args)
	}

	// Set optionals, default is false
	if v, ok := d.GetOk("return_behavior"); ok {
		builder.WithReturnBehavior(v.(string))
	}

	// Set optionals, default is false
	if v, ok := d.GetOk("null_input_behavior"); ok {
		builder.WithNullInputBehavior(v.(string))
	}

	// Set optionals, default is OWNER
	if v, ok := d.GetOk("language"); ok {
		builder.WithLanguage(v.(string))
	}

	// Set optionals, runtime version for Python
	if v, ok := d.GetOk("runtime_version"); ok {
		builder.WithRuntimeVersion(v.(string))
	}

	// Set optionals, default is false
	if v, ok := d.GetOk("is_secure"); ok && v.(bool) {
		builder.WithSecure()
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	// Set optionals, packages for Java / Python
	if _, ok := d.GetOk("packages"); ok {
		packages := []string{}
		for _, pack := range d.Get("packages").([]interface{}) {
			packages = append(packages, pack.(string))
		}
		builder.WithPackages(packages)
	}

	// Set optionals, imports for Java / Python
	if _, ok := d.GetOk("imports"); ok {
		imports := []string{}
		for _, imp := range d.Get("imports").([]interface{}) {
			imports = append(imports, imp.(string))
		}
		builder.WithImports(imports)
	}

	// handler for Java / Python
	if v, ok := d.GetOk("handler"); ok {
		builder.WithHandler(v.(string))
	}

	// target path for Java / Python
	if v, ok := d.GetOk("target_path"); ok {
		builder.WithTargetPath(v.(string))
	}

	q, err := builder.Create()
	if err != nil {
		return err
	}
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error creating function %v err = %w", name, err)
	}

	functionID := &functionID{
		DatabaseName: database,
		SchemaName:   schema,
		FunctionName: name,
		ArgTypes:     builder.ArgTypes(),
	}

	d.SetId(functionID.String())

	return ReadFunction(d, meta)
}

// ReadFunction implements schema.ReadFunc.
func ReadFunction(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	functionID, err := splitFunctionID(d.Id())
	if err != nil {
		return err
	}
	funct := snowflake.NewFunctionBuilder(
		functionID.DatabaseName,
		functionID.SchemaName,
		functionID.FunctionName,
		functionID.ArgTypes,
	)

	// some attributes can be retrieved only by Describe and some only by Show
	stmt, err := funct.Describe()
	if err != nil {
		return err
	}
	rows, err := snowflake.Query(db, stmt)
	if err != nil && snowflake.IsResourceNotExistOrNotAuthorized(err.Error(), "Function") {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] function (%s) not found or we are not authorized.Err:\n%s", d.Id(), err.Error())
		d.SetId("")
		return nil
	}
	defer rows.Close()
	descPropValues, err := snowflake.ScanFunctionDescription(rows)
	if err != nil {
		return err
	}
	for _, desc := range descPropValues {
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

				if err := d.Set("arguments", args); err != nil {
					return err
				}
			}
		case "null handling":
			if err := d.Set("null_input_behavior", desc.Value.String); err != nil {
				return err
			}
		case "volatility":
			if err := d.Set("return_behavior", desc.Value.String); err != nil {
				return err
			}
		case "body":
			if err := d.Set("statement", desc.Value.String); err != nil {
				return err
			}
		case "returns":
			// Format in Snowflake DB is returnType(<some number>)
			re := regexp.MustCompile(`^(.*)\([0-9]*\)$`)
			match := re.FindStringSubmatch(desc.Value.String)
			rt := desc.Value.String
			if match != nil {
				rt = match[1]
			}
			if err := d.Set("return_type", rt); err != nil {
				return err
			}
		case "language":
			if snowflake.Contains(languages, desc.Value.String) {
				if err := d.Set("language", desc.Value.String); err != nil {
					return err
				}
			}
		case "packages":
			packagesString := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(desc.Value.String, "[", ""), "]", ""), "'", "")
			if packagesString != "" { // Do nothing for Java / Python functions without packages
				packages := strings.Split(packagesString, ",")
				if err := d.Set("packages", packages); err != nil {
					return err
				}
			}
		case "imports":
			importsString := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(desc.Value.String, "[", ""), "]", ""), "'", "")
			if importsString != "" { // Do nothing for Java functions without imports
				imports := strings.Split(importsString, ",")
				if err := d.Set("imports", imports); err != nil {
					return err
				}
			}
		case "handler":
			if err := d.Set("handler", desc.Value.String); err != nil {
				return err
			}
		case "target_path":
			if err := d.Set("target_path", desc.Value.String); err != nil {
				return err
			}
		case "runtime_version":
			if err := d.Set("runtime_version", desc.Value.String); err != nil {
				return err
			}
		default:
			log.Printf("[WARN] unexpected function property %v returned from Snowflake", desc.Property.String)
		}
	}

	q := funct.Show()
	showRows, err := snowflake.Query(db, q)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] function (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	defer showRows.Close()

	foundFunctions, err := snowflake.ScanFunctions(showRows)
	if err != nil {
		return err
	}
	// function names can be overloaded with different argument types so we
	// iterate over and find the correct one
	argSig, _ := funct.ArgumentsSignature()

	functionIsSecure := map[string]bool{
		"Y": true,
		"N": false,
	}

	for _, v := range foundFunctions {
		if v.Arguments.String == argSig {
			if err := d.Set("comment", v.Comment.String); err != nil {
				return err
			}
			if err = d.Set("is_secure", functionIsSecure[v.IsSecure.String]); err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateFunction implements schema.UpdateFunction.
func UpdateFunction(d *schema.ResourceData, meta interface{}) error {
	pID, err := splitFunctionID(d.Id())
	if err != nil {
		return err
	}
	builder := snowflake.NewFunctionBuilder(
		pID.DatabaseName,
		pID.SchemaName,
		pID.FunctionName,
		pID.ArgTypes,
	)

	db := meta.(*sql.DB)
	if d.HasChange("name") {
		name := d.Get("name")
		q, err := builder.Rename(name.(string))
		if err != nil {
			return err
		}
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error renaming function %v", d.Id())
		}
		newID := &functionID{
			DatabaseName: pID.DatabaseName,
			SchemaName:   pID.SchemaName,
			FunctionName: name.(string),
			ArgTypes:     pID.ArgTypes,
		}
		d.SetId(newID.String())
	}
	if d.HasChange("is_secure") {
		secure := d.Get("is_secure")

		if secure.(bool) {
			q, err := builder.Secure()
			if err != nil {
				return err
			}
			if err = snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error setting secure for function %v", d.Id())
			}
		} else {
			q, err := builder.Unsecure()
			if err != nil {
				return err
			}
			if err = snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error unsetting secure for function %v", d.Id())
			}
		}
	}

	if d.HasChange("comment") {
		comment := d.Get("comment")

		if c := comment.(string); c == "" {
			q, err := builder.RemoveComment()
			if err != nil {
				return err
			}
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error unsetting comment for function %v err = %w", d.Id(), err)
			}
		} else {
			q, err := builder.ChangeComment(c)
			if err != nil {
				return err
			}
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error updating comment for function %v err = %w", d.Id(), err)
			}
		}
	}

	return ReadFunction(d, meta)
}

// DeleteFunction implements schema.DeleteFunc.
func DeleteFunction(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	pID, err := splitFunctionID(d.Id())
	if err != nil {
		return err
	}
	builder := snowflake.NewFunctionBuilder(
		pID.DatabaseName,
		pID.SchemaName,
		pID.FunctionName,
		pID.ArgTypes,
	)

	q, err := builder.Drop()
	if err != nil {
		return err
	}
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting function %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}

type functionID struct {
	DatabaseName string
	SchemaName   string
	FunctionName string
	ArgTypes     []string
}

// splitFunctionID takes the <database_name>|<schema_name>|<view_name>|<argtypes> ID and returns
// the functionID struct, for example MYDB|PUBLIC|FUNC1|VARCHAR-DATE-VARCHAR
// returns struct
//
//	DatabaseName: MYDB
//	SchemaName: PUBLIC
//	FunctionName: FUNC1
//	ArgTypes: [VARCHAR, DATE, VARCHAR]
func splitFunctionID(v string) (*functionID, error) {
	arr := strings.Split(v, "|")
	if len(arr) != 4 {
		return nil, fmt.Errorf("ID %v is invalid", v)
	}

	return &functionID{
		DatabaseName: arr[0],
		SchemaName:   arr[1],
		FunctionName: arr[2],
		ArgTypes:     strings.Split(arr[3], "-"),
	}, nil
}

// the opposite of splitFunctionID.
func (pi *functionID) String() string {
	return fmt.Sprintf("%v|%v|%v|%v",
		pi.DatabaseName,
		pi.SchemaName,
		pi.FunctionName,
		strings.Join(pi.ArgTypes, "-"))
}
