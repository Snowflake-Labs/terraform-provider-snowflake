package resources

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var procedureLanguages = []string{"javascript", "java", "scala", "SQL", "python"}

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
		ValidateFunc: validation.StringInSlice(procedureLanguages, true),
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

func DiffTypes(_, o, n string, _ *schema.ResourceData) bool {
	return strings.EqualFold(strings.ToUpper(o), strings.ToUpper(n))
}

// Procedure returns a pointer to the resource representing a stored procedure.
func Procedure() *schema.Resource {
	return &schema.Resource{
		Create: CreateProcedure,
		Read:   ReadProcedure,
		Update: UpdateProcedure,
		Delete: DeleteProcedure,

		Schema: procedureSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateProcedure implements schema.CreateFunc.
func CreateProcedure(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	schema := d.Get("schema").(string)
	database := d.Get("database").(string)
	s := d.Get("statement").(string)
	ret := d.Get("return_type").(string)

	builder := snowflake.NewProcedureBuilder(database, schema, name, []string{}).WithStatement(s).WithReturnType(ret)

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
	if v, ok := d.GetOk("execute_as"); ok {
		builder.WithExecuteAs(v.(string))
	}

	// Set optionals, default is SQL
	if v, ok := d.GetOk("language"); ok {
		builder.WithLanguage(strings.ToUpper(v.(string)))
	}

	// Set optionals, runtime version for Python
	if v, ok := d.GetOk("runtime_version"); ok {
		builder.WithRuntimeVersion(v.(string))
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

	q, err := builder.Create()
	if err != nil {
		return err
	}

	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error creating procedure %v err = %w", name, err)
	}

	procedureID := &procedureID{
		DatabaseName:  database,
		SchemaName:    schema,
		ProcedureName: name,
		ArgTypes:      builder.ArgTypes(),
	}

	d.SetId(procedureID.String())

	return ReadProcedure(d, meta)
}

// ReadProcedure implements schema.ReadFunc.
func ReadProcedure(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	procedureID, err := splitProcedureID(d.Id())
	if err != nil {
		return err
	}
	proc := snowflake.NewProcedureBuilder(
		procedureID.DatabaseName,
		procedureID.SchemaName,
		procedureID.ProcedureName,
		procedureID.ArgTypes,
	)

	// some attributes can be retrieved only by Describe and some only by Show
	stmt, err := proc.Describe()
	if err != nil {
		return err
	}
	rows, err := snowflake.Query(db, stmt)
	if err != nil {
		return err
	}
	defer rows.Close()
	descPropValues, err := snowflake.ScanProcedureDescription(rows)
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
		case "execute as":
			if err := d.Set("execute_as", desc.Value.String); err != nil {
				return err
			}
		case "returns":
			if err := d.Set("return_type", desc.Value.String); err != nil {
				return err
			}
		case "language":
			if err := d.Set("language", desc.Value.String); err != nil {
				return err
			}
		case "runtime_version":
			if err := d.Set("runtime_version", desc.Value.String); err != nil {
				return err
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
			importsString := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(desc.Value.String, "[", ""), "]", ""), "'", ""), " ", "")
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

		default:
			log.Printf("[WARN] unexpected procedure property %v returned from Snowflake", desc.Property.String)
		}
	}

	q := proc.Show()
	showRows, err := snowflake.Query(db, q)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] procedure (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	defer showRows.Close()

	foundProcedures, err := snowflake.ScanProcedures(showRows)
	if err != nil {
		return err
	}
	// procedure names can be overloaded with different argument types so we
	// iterate over and find the correct one
	argSig, _ := proc.ArgumentsSignature()

	for _, v := range foundProcedures {
		showArgs := strings.Split(v.Arguments.String, " RETURN ")
		if showArgs[0] == argSig {
			if err := d.Set("name", v.Name.String); err != nil {
				return err
			}
			database := strings.Trim(v.DatabaseName.String, "\"")
			if err := d.Set("database", database); err != nil {
				return err
			}
			schema := strings.Trim(v.SchemaName.String, "\"")
			if err := d.Set("schema", schema); err != nil {
				return err
			}
			if err := d.Set("comment", v.Comment.String); err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateProcedure implements schema.UpdateProcedure.
func UpdateProcedure(d *schema.ResourceData, meta interface{}) error {
	pID, err := splitProcedureID(d.Id())
	if err != nil {
		return err
	}
	builder := snowflake.NewProcedureBuilder(
		pID.DatabaseName,
		pID.SchemaName,
		pID.ProcedureName,
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
			return fmt.Errorf("error renaming procedure %v", d.Id())
		}
		newID := &procedureID{
			DatabaseName:  pID.DatabaseName,
			SchemaName:    pID.SchemaName,
			ProcedureName: name.(string),
			ArgTypes:      pID.ArgTypes,
		}
		d.SetId(newID.String())
	}

	if d.HasChange("comment") {
		comment := d.Get("comment")

		if c := comment.(string); c == "" {
			q, err := builder.RemoveComment()
			if err != nil {
				return err
			}
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error unsetting comment for procedure %v", d.Id())
			}
		} else {
			q, err := builder.ChangeComment(c)
			if err != nil {
				return err
			}
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error updating comment for procedure %v", d.Id())
			}
		}
	}
	if d.HasChange("execute_as") {
		executeAs := d.Get("execute_as")

		q, err := builder.ChangeExecuteAs(executeAs.(string))
		if err != nil {
			return err
		}
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error changing execute as for procedure %v", d.Id())
		}
	}

	return ReadProcedure(d, meta)
}

// DeleteProcedure implements schema.DeleteFunc.
func DeleteProcedure(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	pID, err := splitProcedureID(d.Id())
	if err != nil {
		return err
	}
	builder := snowflake.NewProcedureBuilder(
		pID.DatabaseName,
		pID.SchemaName,
		pID.ProcedureName,
		pID.ArgTypes,
	)

	q, err := builder.Drop()
	if err != nil {
		return err
	}

	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting procedure %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}

type procedureID struct {
	DatabaseName  string
	SchemaName    string
	ProcedureName string
	ArgTypes      []string
}

// splitProcedureID takes the <database_name>|<schema_name>|<view_name>|<argtypes> ID and returns
// the procedureID struct, for example MYDB|PUBLIC|PROC1|VARCHAR-DATE-VARCHAR
// returns struct
//
//	DatabaseName: MYDB
//	SchemaName: PUBLIC
//	ProcedureName: PROC1
//	ArgTypes: [VARCHAR, DATE, VARCHAR]
func splitProcedureID(v string) (*procedureID, error) {
	arr := strings.Split(v, "|")
	if len(arr) != 4 {
		return nil, fmt.Errorf("ID %v is invalid", v)
	}

	return &procedureID{
		DatabaseName:  arr[0],
		SchemaName:    arr[1],
		ProcedureName: arr[2],
		ArgTypes:      strings.Split(arr[3], "-"),
	}, nil
}

// the opposite of splitProcedureID.
func (pi *procedureID) String() string {
	return fmt.Sprintf("%v|%v|%v|%v",
		pi.DatabaseName,
		pi.SchemaName,
		pi.ProcedureName,
		strings.Join(pi.ArgTypes, "-"))
}
