// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type FunctionPythonModel struct {
	Arguments                  tfconfig.Variable `json:"arguments,omitempty"`
	Comment                    tfconfig.Variable `json:"comment,omitempty"`
	Database                   tfconfig.Variable `json:"database,omitempty"`
	EnableConsoleOutput        tfconfig.Variable `json:"enable_console_output,omitempty"`
	ExternalAccessIntegrations tfconfig.Variable `json:"external_access_integrations,omitempty"`
	FullyQualifiedName         tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	FunctionDefinition         tfconfig.Variable `json:"function_definition,omitempty"`
	FunctionLanguage           tfconfig.Variable `json:"function_language,omitempty"`
	Handler                    tfconfig.Variable `json:"handler,omitempty"`
	Imports                    tfconfig.Variable `json:"imports,omitempty"`
	IsAggregate                tfconfig.Variable `json:"is_aggregate,omitempty"`
	IsSecure                   tfconfig.Variable `json:"is_secure,omitempty"`
	LogLevel                   tfconfig.Variable `json:"log_level,omitempty"`
	MetricLevel                tfconfig.Variable `json:"metric_level,omitempty"`
	Name                       tfconfig.Variable `json:"name,omitempty"`
	NullInputBehavior          tfconfig.Variable `json:"null_input_behavior,omitempty"`
	Packages                   tfconfig.Variable `json:"packages,omitempty"`
	ReturnBehavior             tfconfig.Variable `json:"return_behavior,omitempty"`
	ReturnType                 tfconfig.Variable `json:"return_type,omitempty"`
	RuntimeVersion             tfconfig.Variable `json:"runtime_version,omitempty"`
	Schema                     tfconfig.Variable `json:"schema,omitempty"`
	Secrets                    tfconfig.Variable `json:"secrets,omitempty"`
	TraceLevel                 tfconfig.Variable `json:"trace_level,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func FunctionPython(
	resourceName string,
	database string,
	functionDefinition string,
	handler string,
	name string,
	returnType string,
	runtimeVersion string,
	schema string,
) *FunctionPythonModel {
	f := &FunctionPythonModel{ResourceModelMeta: config.Meta(resourceName, resources.FunctionPython)}
	f.WithDatabase(database)
	f.WithFunctionDefinition(functionDefinition)
	f.WithHandler(handler)
	f.WithName(name)
	f.WithReturnType(returnType)
	f.WithRuntimeVersion(runtimeVersion)
	f.WithSchema(schema)
	return f
}

func FunctionPythonWithDefaultMeta(
	database string,
	functionDefinition string,
	handler string,
	name string,
	returnType string,
	runtimeVersion string,
	schema string,
) *FunctionPythonModel {
	f := &FunctionPythonModel{ResourceModelMeta: config.DefaultMeta(resources.FunctionPython)}
	f.WithDatabase(database)
	f.WithFunctionDefinition(functionDefinition)
	f.WithHandler(handler)
	f.WithName(name)
	f.WithReturnType(returnType)
	f.WithRuntimeVersion(runtimeVersion)
	f.WithSchema(schema)
	return f
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

// arguments attribute type is not yet supported, so WithArguments can't be generated

func (f *FunctionPythonModel) WithComment(comment string) *FunctionPythonModel {
	f.Comment = tfconfig.StringVariable(comment)
	return f
}

func (f *FunctionPythonModel) WithDatabase(database string) *FunctionPythonModel {
	f.Database = tfconfig.StringVariable(database)
	return f
}

func (f *FunctionPythonModel) WithEnableConsoleOutput(enableConsoleOutput bool) *FunctionPythonModel {
	f.EnableConsoleOutput = tfconfig.BoolVariable(enableConsoleOutput)
	return f
}

// external_access_integrations attribute type is not yet supported, so WithExternalAccessIntegrations can't be generated

func (f *FunctionPythonModel) WithFullyQualifiedName(fullyQualifiedName string) *FunctionPythonModel {
	f.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return f
}

func (f *FunctionPythonModel) WithFunctionDefinition(functionDefinition string) *FunctionPythonModel {
	f.FunctionDefinition = tfconfig.StringVariable(functionDefinition)
	return f
}

func (f *FunctionPythonModel) WithFunctionLanguage(functionLanguage string) *FunctionPythonModel {
	f.FunctionLanguage = tfconfig.StringVariable(functionLanguage)
	return f
}

func (f *FunctionPythonModel) WithHandler(handler string) *FunctionPythonModel {
	f.Handler = tfconfig.StringVariable(handler)
	return f
}

// imports attribute type is not yet supported, so WithImports can't be generated

func (f *FunctionPythonModel) WithIsAggregate(isAggregate string) *FunctionPythonModel {
	f.IsAggregate = tfconfig.StringVariable(isAggregate)
	return f
}

func (f *FunctionPythonModel) WithIsSecure(isSecure string) *FunctionPythonModel {
	f.IsSecure = tfconfig.StringVariable(isSecure)
	return f
}

func (f *FunctionPythonModel) WithLogLevel(logLevel string) *FunctionPythonModel {
	f.LogLevel = tfconfig.StringVariable(logLevel)
	return f
}

func (f *FunctionPythonModel) WithMetricLevel(metricLevel string) *FunctionPythonModel {
	f.MetricLevel = tfconfig.StringVariable(metricLevel)
	return f
}

func (f *FunctionPythonModel) WithName(name string) *FunctionPythonModel {
	f.Name = tfconfig.StringVariable(name)
	return f
}

func (f *FunctionPythonModel) WithNullInputBehavior(nullInputBehavior string) *FunctionPythonModel {
	f.NullInputBehavior = tfconfig.StringVariable(nullInputBehavior)
	return f
}

// packages attribute type is not yet supported, so WithPackages can't be generated

func (f *FunctionPythonModel) WithReturnBehavior(returnBehavior string) *FunctionPythonModel {
	f.ReturnBehavior = tfconfig.StringVariable(returnBehavior)
	return f
}

func (f *FunctionPythonModel) WithReturnType(returnType string) *FunctionPythonModel {
	f.ReturnType = tfconfig.StringVariable(returnType)
	return f
}

func (f *FunctionPythonModel) WithRuntimeVersion(runtimeVersion string) *FunctionPythonModel {
	f.RuntimeVersion = tfconfig.StringVariable(runtimeVersion)
	return f
}

func (f *FunctionPythonModel) WithSchema(schema string) *FunctionPythonModel {
	f.Schema = tfconfig.StringVariable(schema)
	return f
}

// secrets attribute type is not yet supported, so WithSecrets can't be generated

func (f *FunctionPythonModel) WithTraceLevel(traceLevel string) *FunctionPythonModel {
	f.TraceLevel = tfconfig.StringVariable(traceLevel)
	return f
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (f *FunctionPythonModel) WithArgumentsValue(value tfconfig.Variable) *FunctionPythonModel {
	f.Arguments = value
	return f
}

func (f *FunctionPythonModel) WithCommentValue(value tfconfig.Variable) *FunctionPythonModel {
	f.Comment = value
	return f
}

func (f *FunctionPythonModel) WithDatabaseValue(value tfconfig.Variable) *FunctionPythonModel {
	f.Database = value
	return f
}

func (f *FunctionPythonModel) WithEnableConsoleOutputValue(value tfconfig.Variable) *FunctionPythonModel {
	f.EnableConsoleOutput = value
	return f
}

func (f *FunctionPythonModel) WithExternalAccessIntegrationsValue(value tfconfig.Variable) *FunctionPythonModel {
	f.ExternalAccessIntegrations = value
	return f
}

func (f *FunctionPythonModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *FunctionPythonModel {
	f.FullyQualifiedName = value
	return f
}

func (f *FunctionPythonModel) WithFunctionDefinitionValue(value tfconfig.Variable) *FunctionPythonModel {
	f.FunctionDefinition = value
	return f
}

func (f *FunctionPythonModel) WithFunctionLanguageValue(value tfconfig.Variable) *FunctionPythonModel {
	f.FunctionLanguage = value
	return f
}

func (f *FunctionPythonModel) WithHandlerValue(value tfconfig.Variable) *FunctionPythonModel {
	f.Handler = value
	return f
}

func (f *FunctionPythonModel) WithImportsValue(value tfconfig.Variable) *FunctionPythonModel {
	f.Imports = value
	return f
}

func (f *FunctionPythonModel) WithIsAggregateValue(value tfconfig.Variable) *FunctionPythonModel {
	f.IsAggregate = value
	return f
}

func (f *FunctionPythonModel) WithIsSecureValue(value tfconfig.Variable) *FunctionPythonModel {
	f.IsSecure = value
	return f
}

func (f *FunctionPythonModel) WithLogLevelValue(value tfconfig.Variable) *FunctionPythonModel {
	f.LogLevel = value
	return f
}

func (f *FunctionPythonModel) WithMetricLevelValue(value tfconfig.Variable) *FunctionPythonModel {
	f.MetricLevel = value
	return f
}

func (f *FunctionPythonModel) WithNameValue(value tfconfig.Variable) *FunctionPythonModel {
	f.Name = value
	return f
}

func (f *FunctionPythonModel) WithNullInputBehaviorValue(value tfconfig.Variable) *FunctionPythonModel {
	f.NullInputBehavior = value
	return f
}

func (f *FunctionPythonModel) WithPackagesValue(value tfconfig.Variable) *FunctionPythonModel {
	f.Packages = value
	return f
}

func (f *FunctionPythonModel) WithReturnBehaviorValue(value tfconfig.Variable) *FunctionPythonModel {
	f.ReturnBehavior = value
	return f
}

func (f *FunctionPythonModel) WithReturnTypeValue(value tfconfig.Variable) *FunctionPythonModel {
	f.ReturnType = value
	return f
}

func (f *FunctionPythonModel) WithRuntimeVersionValue(value tfconfig.Variable) *FunctionPythonModel {
	f.RuntimeVersion = value
	return f
}

func (f *FunctionPythonModel) WithSchemaValue(value tfconfig.Variable) *FunctionPythonModel {
	f.Schema = value
	return f
}

func (f *FunctionPythonModel) WithSecretsValue(value tfconfig.Variable) *FunctionPythonModel {
	f.Secrets = value
	return f
}

func (f *FunctionPythonModel) WithTraceLevelValue(value tfconfig.Variable) *FunctionPythonModel {
	f.TraceLevel = value
	return f
}