// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type FunctionJavascriptModel struct {
	Arguments             tfconfig.Variable `json:"arguments,omitempty"`
	Comment               tfconfig.Variable `json:"comment,omitempty"`
	Database              tfconfig.Variable `json:"database,omitempty"`
	EnableConsoleOutput   tfconfig.Variable `json:"enable_console_output,omitempty"`
	FullyQualifiedName    tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	FunctionDefinition    tfconfig.Variable `json:"function_definition,omitempty"`
	FunctionLanguage      tfconfig.Variable `json:"function_language,omitempty"`
	IsSecure              tfconfig.Variable `json:"is_secure,omitempty"`
	LogLevel              tfconfig.Variable `json:"log_level,omitempty"`
	MetricLevel           tfconfig.Variable `json:"metric_level,omitempty"`
	Name                  tfconfig.Variable `json:"name,omitempty"`
	NullInputBehavior     tfconfig.Variable `json:"null_input_behavior,omitempty"`
	ReturnResultsBehavior tfconfig.Variable `json:"return_results_behavior,omitempty"`
	ReturnType            tfconfig.Variable `json:"return_type,omitempty"`
	Schema                tfconfig.Variable `json:"schema,omitempty"`
	TraceLevel            tfconfig.Variable `json:"trace_level,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func FunctionJavascript(
	resourceName string,
	database string,
	functionDefinition string,
	name string,
	returnType string,
	schema string,
) *FunctionJavascriptModel {
	f := &FunctionJavascriptModel{ResourceModelMeta: config.Meta(resourceName, resources.FunctionJavascript)}
	f.WithDatabase(database)
	f.WithFunctionDefinition(functionDefinition)
	f.WithName(name)
	f.WithReturnType(returnType)
	f.WithSchema(schema)
	return f
}

func FunctionJavascriptWithDefaultMeta(
	database string,
	functionDefinition string,
	name string,
	returnType string,
	schema string,
) *FunctionJavascriptModel {
	f := &FunctionJavascriptModel{ResourceModelMeta: config.DefaultMeta(resources.FunctionJavascript)}
	f.WithDatabase(database)
	f.WithFunctionDefinition(functionDefinition)
	f.WithName(name)
	f.WithReturnType(returnType)
	f.WithSchema(schema)
	return f
}

///////////////////////////////////////////////////////
// set proper json marshalling and handle depends on //
///////////////////////////////////////////////////////

func (f *FunctionJavascriptModel) MarshalJSON() ([]byte, error) {
	type Alias FunctionJavascriptModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}

func (f *FunctionJavascriptModel) WithDependsOn(values ...string) *FunctionJavascriptModel {
	f.SetDependsOn(values...)
	return f
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

// arguments attribute type is not yet supported, so WithArguments can't be generated

func (f *FunctionJavascriptModel) WithComment(comment string) *FunctionJavascriptModel {
	f.Comment = tfconfig.StringVariable(comment)
	return f
}

func (f *FunctionJavascriptModel) WithDatabase(database string) *FunctionJavascriptModel {
	f.Database = tfconfig.StringVariable(database)
	return f
}

func (f *FunctionJavascriptModel) WithEnableConsoleOutput(enableConsoleOutput bool) *FunctionJavascriptModel {
	f.EnableConsoleOutput = tfconfig.BoolVariable(enableConsoleOutput)
	return f
}

func (f *FunctionJavascriptModel) WithFullyQualifiedName(fullyQualifiedName string) *FunctionJavascriptModel {
	f.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return f
}

func (f *FunctionJavascriptModel) WithFunctionDefinition(functionDefinition string) *FunctionJavascriptModel {
	f.FunctionDefinition = config.MultilineWrapperVariable(functionDefinition)
	return f
}

func (f *FunctionJavascriptModel) WithFunctionLanguage(functionLanguage string) *FunctionJavascriptModel {
	f.FunctionLanguage = tfconfig.StringVariable(functionLanguage)
	return f
}

func (f *FunctionJavascriptModel) WithIsSecure(isSecure string) *FunctionJavascriptModel {
	f.IsSecure = tfconfig.StringVariable(isSecure)
	return f
}

func (f *FunctionJavascriptModel) WithLogLevel(logLevel string) *FunctionJavascriptModel {
	f.LogLevel = tfconfig.StringVariable(logLevel)
	return f
}

func (f *FunctionJavascriptModel) WithMetricLevel(metricLevel string) *FunctionJavascriptModel {
	f.MetricLevel = tfconfig.StringVariable(metricLevel)
	return f
}

func (f *FunctionJavascriptModel) WithName(name string) *FunctionJavascriptModel {
	f.Name = tfconfig.StringVariable(name)
	return f
}

func (f *FunctionJavascriptModel) WithNullInputBehavior(nullInputBehavior string) *FunctionJavascriptModel {
	f.NullInputBehavior = tfconfig.StringVariable(nullInputBehavior)
	return f
}

func (f *FunctionJavascriptModel) WithReturnResultsBehavior(returnResultsBehavior string) *FunctionJavascriptModel {
	f.ReturnResultsBehavior = tfconfig.StringVariable(returnResultsBehavior)
	return f
}

func (f *FunctionJavascriptModel) WithReturnType(returnType string) *FunctionJavascriptModel {
	f.ReturnType = tfconfig.StringVariable(returnType)
	return f
}

func (f *FunctionJavascriptModel) WithSchema(schema string) *FunctionJavascriptModel {
	f.Schema = tfconfig.StringVariable(schema)
	return f
}

func (f *FunctionJavascriptModel) WithTraceLevel(traceLevel string) *FunctionJavascriptModel {
	f.TraceLevel = tfconfig.StringVariable(traceLevel)
	return f
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (f *FunctionJavascriptModel) WithArgumentsValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.Arguments = value
	return f
}

func (f *FunctionJavascriptModel) WithCommentValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.Comment = value
	return f
}

func (f *FunctionJavascriptModel) WithDatabaseValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.Database = value
	return f
}

func (f *FunctionJavascriptModel) WithEnableConsoleOutputValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.EnableConsoleOutput = value
	return f
}

func (f *FunctionJavascriptModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.FullyQualifiedName = value
	return f
}

func (f *FunctionJavascriptModel) WithFunctionDefinitionValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.FunctionDefinition = value
	return f
}

func (f *FunctionJavascriptModel) WithFunctionLanguageValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.FunctionLanguage = value
	return f
}

func (f *FunctionJavascriptModel) WithIsSecureValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.IsSecure = value
	return f
}

func (f *FunctionJavascriptModel) WithLogLevelValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.LogLevel = value
	return f
}

func (f *FunctionJavascriptModel) WithMetricLevelValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.MetricLevel = value
	return f
}

func (f *FunctionJavascriptModel) WithNameValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.Name = value
	return f
}

func (f *FunctionJavascriptModel) WithNullInputBehaviorValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.NullInputBehavior = value
	return f
}

func (f *FunctionJavascriptModel) WithReturnResultsBehaviorValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.ReturnResultsBehavior = value
	return f
}

func (f *FunctionJavascriptModel) WithReturnTypeValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.ReturnType = value
	return f
}

func (f *FunctionJavascriptModel) WithSchemaValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.Schema = value
	return f
}

func (f *FunctionJavascriptModel) WithTraceLevelValue(value tfconfig.Variable) *FunctionJavascriptModel {
	f.TraceLevel = value
	return f
}
