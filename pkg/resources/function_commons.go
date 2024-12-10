package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {

}

// TODO [next PR]: currently all database.schema.name are ForceNew but based on the docs it is possible to rename with moving to different db/schema
// TODO [next PR]: copyGrants and orReplace logic omitted for now, will be added to the limitations docs
// TODO [next PR]: temporary is not supported because it creates a per-session object; add to limitations/design decisions
var functionBaseSchema = map[string]*schema.Schema{
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the function."),
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the function."),
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The name of the function; the identifier does not need to be unique for the schema in which the function is created because UDFs are identified and resolved by the combination of the name and argument types. Check the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#all-languages)."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"is_secure": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("is_secure"),
		Description:      booleanStringFieldDescription("Specifies that the function is secure. By design, the Snowflake's `SHOW FUNCTIONS` command does not provide information about secure views (consult [function docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#id1) and [Protecting Sensitive Information with Secure UDFs and Stored Procedures](https://docs.snowflake.com/en/developer-guide/secure-udf-procedure)) which is essential to manage/import function with Terraform. Use the role owning the function while managing secure functions."),
	},
	"is_aggregate": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("is_aggregate"),
		Description:      booleanStringFieldDescription("Specifies that the function is an aggregate function. For more information about user-defined aggregate functions, see [Python user-defined aggregate functions](https://docs.snowflake.com/en/developer-guide/udf/python/udf-python-aggregate-functions)."),
	},
	"arguments": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"arg_name": {
					Type:     schema.TypeString,
					Required: true,
					// TODO [SNOW-1348103]: adjust diff suppression accordingly.
					Description: "The argument name.",
				},
				"arg_data_type": {
					Type:     schema.TypeString,
					Required: true,
					// TODO [SNOW-1348103]: adjust diff suppression accordingly.
					Description: "The argument type.",
				},
			},
		},
		Optional:    true,
		ForceNew:    true,
		Description: "List of the arguments for the function. Consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#all-languages) for more details.",
	},
	// TODO [SNOW-1348103]: for now, the proposal is to leave return type as string, add TABLE to data types, and here always parse (easier handling and diff suppression)
	"return_type": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the results returned by the UDF, which determines the UDF type. Use `<result_data_type>` to create a scalar UDF that returns a single value with the specified data type. Use `TABLE (col_name col_data_type, ...)` to creates a table UDF that returns tabular results with the specified table column(s) and column type(s). For the details, consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#all-languages).",
		// TODO [SNOW-1348103]: adjust DiffSuppressFunc
	},
	"null_input_behavior": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToNullInputBehavior),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToNullInputBehavior), IgnoreChangeToCurrentSnowflakeValueInShow("null_input_behavior")),
		Description:      fmt.Sprintf("Specifies the behavior of the function when called with null inputs. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllAllowedNullInputBehaviors)),
	},
	"return_behavior": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToReturnResultsBehavior),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToReturnResultsBehavior), IgnoreChangeToCurrentSnowflakeValueInShow("return_behavior")),
		Description:      fmt.Sprintf("Specifies the behavior of the function when returning results. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllAllowedReturnResultsBehaviors)),
	},
	// TODO [this PR]: set optional/required correctly
	"runtime_version": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		// TODO [this PR]: may be optional for java without consequence because if it is not set, the describe is not returning any version.
		//Description: "Specifies the Java JDK runtime version to use. The supported versions of Java are 11.x and 17.x. If RUNTIME_VERSION is not set, Java JDK 11 is used.",
		//Description: "Specifies the Python version to use. The supported versions of Python are: 3.9, 3.10, and 3.11.",
		//Description: "Specifies the Scala runtime version to use. The supported versions of Scala are: 2.12.",
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
		// TODO [SNOW-1348103]: handle dynamic comment - this is a workaround for now
		Default:     "user-defined function",
		Description: "Specifies a comment for the function.",
	},
	// TODO [SNOW-1348103]: because of https://docs.snowflake.com/en/sql-reference/sql/create-function#id6, maybe it will be better to split into stage_name + target_path
	"imports": {
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Optional: true,
		ForceNew: true,
		//Description: "The location (stage), path, and name of the file(s) to import. A file can be a JAR file or another type of file. If the file is a JAR file, it can contain one or more .class files and zero or more resource files. JNI (Java Native Interface) is not supported. Snowflake prohibits loading libraries that contain native code (as opposed to Java bytecode). Java UDFs can also read non-JAR files. For an example, see [Reading a file specified statically in IMPORTS](https://docs.snowflake.com/en/developer-guide/udf/java/udf-java-cookbook.html#label-reading-file-from-java-udf-imports). Consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#java).",
		//Description: "The location (stage), path, and name of the file(s) to import. A file can be a `.py` file or another type of file. Python UDFs can also read non-Python files, such as text files. For an example, see [Reading a file](https://docs.snowflake.com/en/developer-guide/udf/python/udf-python-examples.html#label-udf-python-read-files). Consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#python).",
		//Description: "The location (stage), path, and name of the file(s) to import, such as a JAR or other kind of file. The JAR file might contain handler dependency libraries. It can contain one or more .class files and zero or more resource files. JNI (Java Native Interface) is not supported. Snowflake prohibits loading libraries that contain native code (as opposed to Java bytecode). A non-JAR file might a file read by handler code. For an example, see [Reading a file specified statically in IMPORTS](https://docs.snowflake.com/en/developer-guide/udf/java/udf-java-cookbook.html#label-reading-file-from-java-udf-imports). Consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-function#scala).",
	},
	"packages": {
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Optional: true,
		ForceNew: true,
		// TODO [SNOW-1348103]: what do we do with the version "latest".
		//Description: "The name and version number of Snowflake system packages required as dependencies. The value should be of the form `package_name:version_number`, where `package_name` is `snowflake_domain:package`.",
		//Description: "The name and version number of packages required as dependencies. The value should be of the form `package_name==version_number`.",
		//Description: "The name and version number of Snowflake system packages required as dependencies. The value should be of the form `package_name:version_number`, where `package_name` is `snowflake_domain:package`.",
	},
	"handler": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		// Description: "The name of the handler method or class. If the handler is for a scalar UDF, returning a non-tabular value, the HANDLER value should be a method name, as in the following form: `MyClass.myMethod`. If the handler is for a tabular UDF, the HANDLER value should be the name of a handler class.",
		// Description: "The name of the handler function or class. If the handler is for a scalar UDF, returning a non-tabular value, the HANDLER value should be a function name. If the handler code is in-line with the CREATE FUNCTION statement, you can use the function name alone. When the handler code is referenced at a stage, this value should be qualified with the module name, as in the following form: `my_module.my_function`. If the handler is for a tabular UDF, the HANDLER value should be the name of a handler class.",
		// Description: "The name of the handler method or class. If the handler is for a scalar UDF, returning a non-tabular value, the HANDLER value should be a method name, as in the following form: `MyClass.myMethod`.",
	},
	"external_access_integrations": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		Optional:    true,
		ForceNew:    true,
		Description: "The names of [external access integrations](https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration) needed in order for this function’s handler code to access external networks. An external access integration specifies [network rules](https://docs.snowflake.com/en/sql-reference/sql/create-network-rule) and [secrets](https://docs.snowflake.com/en/sql-reference/sql/create-secret) that specify external locations and credentials (if any) allowed for use by handler code when making requests of an external network, such as an external REST API.",
	},
	"secrets": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"secret_variable_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The variable that will be used in handler code when retrieving information from the secret.",
				},
				"secret_id": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Fully qualified name of the allowed secret. You will receive an error if you specify a SECRETS value whose secret isn’t also included in an integration specified by the EXTERNAL_ACCESS_INTEGRATIONS parameter.",
					DiffSuppressFunc: suppressIdentifierQuoting,
				},
			},
		},
		Description: "Assigns the names of secrets to variables so that you can use the variables to reference the secrets when retrieving information from secrets in handler code. Secrets you specify here must be allowed by the [external access integration](https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration) specified as a value of this CREATE FUNCTION command’s EXTERNAL_ACCESS_INTEGRATIONS parameter.",
	},
	"target_path": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		//Description: "The TARGET_PATH clause specifies the location to which Snowflake should write the compiled code (JAR file) after compiling the source code specified in the `function_definition`. If this clause is included, the user should manually remove the JAR file when it is no longer needed (typically when the Java UDF is dropped). If this clause is omitted, Snowflake re-compiles the source code each time the code is needed. The JAR file is not stored permanently, and the user does not need to clean up the JAR file. Snowflake returns an error if the TARGET_PATH matches an existing file; you cannot use TARGET_PATH to overwrite an existing file.",
		//Description: "The TARGET_PATH clause specifies the location to which Snowflake should write the compiled code (JAR file) after compiling the source code specified in the `function_definition`. If this clause is included, you should manually remove the JAR file when it is no longer needed (typically when the UDF is dropped). If this clause is omitted, Snowflake re-compiles the source code each time the code is needed. The JAR file is not stored permanently, and you do not need to clean up the JAR file. Snowflake returns an error if the TARGET_PATH matches an existing file; you cannot use TARGET_PATH to overwrite an existing file.",
	},
	"function_definition": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: DiffSuppressStatement,
		// TODO [this PR]: generalize description for function_definition
	},
	"function_language": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies language for the user. Used to detect external changes.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW FUNCTION` for the given function.",
		Elem: &schema.Resource{
			Schema: schemas.ShowFunctionSchema,
		},
	},
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN FUNCTION` for the given function.",
		Elem: &schema.Resource{
			Schema: functionParametersSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}
