// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type FunctionScalaModel struct {
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
	IsSecure                   tfconfig.Variable `json:"is_secure,omitempty"`
	LogLevel                   tfconfig.Variable `json:"log_level,omitempty"`
	MetricLevel                tfconfig.Variable `json:"metric_level,omitempty"`
	Name                       tfconfig.Variable `json:"name,omitempty"`
	NullInputBehavior          tfconfig.Variable `json:"null_input_behavior,omitempty"`
	Packages                   tfconfig.Variable `json:"packages,omitempty"`
	ReturnResultsBehavior      tfconfig.Variable `json:"return_results_behavior,omitempty"`
	ReturnType                 tfconfig.Variable `json:"return_type,omitempty"`
	RuntimeVersion             tfconfig.Variable `json:"runtime_version,omitempty"`
	Schema                     tfconfig.Variable `json:"schema,omitempty"`
	Secrets                    tfconfig.Variable `json:"secrets,omitempty"`
	TargetPath                 tfconfig.Variable `json:"target_path,omitempty"`
	TraceLevel                 tfconfig.Variable `json:"trace_level,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func FunctionScala(
	resourceName string,
	database string,
	handler string,
	name string,
	returnType string,
	runtimeVersion string,
	schema string,
) *FunctionScalaModel {
	f := &FunctionScalaModel{ResourceModelMeta: config.Meta(resourceName, resources.FunctionScala)}
	f.WithDatabase(database)
	f.WithHandler(handler)
	f.WithName(name)
	f.WithReturnType(returnType)
	f.WithRuntimeVersion(runtimeVersion)
	f.WithSchema(schema)
	return f
}

func FunctionScalaWithDefaultMeta(
	database string,
	handler string,
	name string,
	returnType string,
	runtimeVersion string,
	schema string,
) *FunctionScalaModel {
	f := &FunctionScalaModel{ResourceModelMeta: config.DefaultMeta(resources.FunctionScala)}
	f.WithDatabase(database)
	f.WithHandler(handler)
	f.WithName(name)
	f.WithReturnType(returnType)
	f.WithRuntimeVersion(runtimeVersion)
	f.WithSchema(schema)
	return f
}

///////////////////////////////////////////////////////
// set proper json marshalling and handle depends on //
///////////////////////////////////////////////////////

func (f *FunctionScalaModel) MarshalJSON() ([]byte, error) {
	type Alias FunctionScalaModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}

func (f *FunctionScalaModel) WithDependsOn(values ...string) *FunctionScalaModel {
	f.SetDependsOn(values...)
	return f
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

// arguments attribute type is not yet supported, so WithArguments can't be generated

func (f *FunctionScalaModel) WithComment(comment string) *FunctionScalaModel {
	f.Comment = tfconfig.StringVariable(comment)
	return f
}

func (f *FunctionScalaModel) WithDatabase(database string) *FunctionScalaModel {
	f.Database = tfconfig.StringVariable(database)
	return f
}

func (f *FunctionScalaModel) WithEnableConsoleOutput(enableConsoleOutput bool) *FunctionScalaModel {
	f.EnableConsoleOutput = tfconfig.BoolVariable(enableConsoleOutput)
	return f
}

// external_access_integrations attribute type is not yet supported, so WithExternalAccessIntegrations can't be generated

func (f *FunctionScalaModel) WithFullyQualifiedName(fullyQualifiedName string) *FunctionScalaModel {
	f.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return f
}

func (f *FunctionScalaModel) WithFunctionDefinition(functionDefinition string) *FunctionScalaModel {
	f.FunctionDefinition = config.MultilineWrapperVariable(functionDefinition)
	return f
}

func (f *FunctionScalaModel) WithFunctionLanguage(functionLanguage string) *FunctionScalaModel {
	f.FunctionLanguage = tfconfig.StringVariable(functionLanguage)
	return f
}

func (f *FunctionScalaModel) WithHandler(handler string) *FunctionScalaModel {
	f.Handler = tfconfig.StringVariable(handler)
	return f
}

// imports attribute type is not yet supported, so WithImports can't be generated

func (f *FunctionScalaModel) WithIsSecure(isSecure string) *FunctionScalaModel {
	f.IsSecure = tfconfig.StringVariable(isSecure)
	return f
}

func (f *FunctionScalaModel) WithLogLevel(logLevel string) *FunctionScalaModel {
	f.LogLevel = tfconfig.StringVariable(logLevel)
	return f
}

func (f *FunctionScalaModel) WithMetricLevel(metricLevel string) *FunctionScalaModel {
	f.MetricLevel = tfconfig.StringVariable(metricLevel)
	return f
}

func (f *FunctionScalaModel) WithName(name string) *FunctionScalaModel {
	f.Name = tfconfig.StringVariable(name)
	return f
}

func (f *FunctionScalaModel) WithNullInputBehavior(nullInputBehavior string) *FunctionScalaModel {
	f.NullInputBehavior = tfconfig.StringVariable(nullInputBehavior)
	return f
}

// packages attribute type is not yet supported, so WithPackages can't be generated

func (f *FunctionScalaModel) WithReturnResultsBehavior(returnResultsBehavior string) *FunctionScalaModel {
	f.ReturnResultsBehavior = tfconfig.StringVariable(returnResultsBehavior)
	return f
}

func (f *FunctionScalaModel) WithReturnType(returnType string) *FunctionScalaModel {
	f.ReturnType = tfconfig.StringVariable(returnType)
	return f
}

func (f *FunctionScalaModel) WithRuntimeVersion(runtimeVersion string) *FunctionScalaModel {
	f.RuntimeVersion = tfconfig.StringVariable(runtimeVersion)
	return f
}

func (f *FunctionScalaModel) WithSchema(schema string) *FunctionScalaModel {
	f.Schema = tfconfig.StringVariable(schema)
	return f
}

// secrets attribute type is not yet supported, so WithSecrets can't be generated

// target_path attribute type is not yet supported, so WithTargetPath can't be generated

func (f *FunctionScalaModel) WithTraceLevel(traceLevel string) *FunctionScalaModel {
	f.TraceLevel = tfconfig.StringVariable(traceLevel)
	return f
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (f *FunctionScalaModel) WithArgumentsValue(value tfconfig.Variable) *FunctionScalaModel {
	f.Arguments = value
	return f
}

func (f *FunctionScalaModel) WithCommentValue(value tfconfig.Variable) *FunctionScalaModel {
	f.Comment = value
	return f
}

func (f *FunctionScalaModel) WithDatabaseValue(value tfconfig.Variable) *FunctionScalaModel {
	f.Database = value
	return f
}

func (f *FunctionScalaModel) WithEnableConsoleOutputValue(value tfconfig.Variable) *FunctionScalaModel {
	f.EnableConsoleOutput = value
	return f
}

func (f *FunctionScalaModel) WithExternalAccessIntegrationsValue(value tfconfig.Variable) *FunctionScalaModel {
	f.ExternalAccessIntegrations = value
	return f
}

func (f *FunctionScalaModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *FunctionScalaModel {
	f.FullyQualifiedName = value
	return f
}

func (f *FunctionScalaModel) WithFunctionDefinitionValue(value tfconfig.Variable) *FunctionScalaModel {
	f.FunctionDefinition = value
	return f
}

func (f *FunctionScalaModel) WithFunctionLanguageValue(value tfconfig.Variable) *FunctionScalaModel {
	f.FunctionLanguage = value
	return f
}

func (f *FunctionScalaModel) WithHandlerValue(value tfconfig.Variable) *FunctionScalaModel {
	f.Handler = value
	return f
}

func (f *FunctionScalaModel) WithImportsValue(value tfconfig.Variable) *FunctionScalaModel {
	f.Imports = value
	return f
}

func (f *FunctionScalaModel) WithIsSecureValue(value tfconfig.Variable) *FunctionScalaModel {
	f.IsSecure = value
	return f
}

func (f *FunctionScalaModel) WithLogLevelValue(value tfconfig.Variable) *FunctionScalaModel {
	f.LogLevel = value
	return f
}

func (f *FunctionScalaModel) WithMetricLevelValue(value tfconfig.Variable) *FunctionScalaModel {
	f.MetricLevel = value
	return f
}

func (f *FunctionScalaModel) WithNameValue(value tfconfig.Variable) *FunctionScalaModel {
	f.Name = value
	return f
}

func (f *FunctionScalaModel) WithNullInputBehaviorValue(value tfconfig.Variable) *FunctionScalaModel {
	f.NullInputBehavior = value
	return f
}

func (f *FunctionScalaModel) WithPackagesValue(value tfconfig.Variable) *FunctionScalaModel {
	f.Packages = value
	return f
}

func (f *FunctionScalaModel) WithReturnResultsBehaviorValue(value tfconfig.Variable) *FunctionScalaModel {
	f.ReturnResultsBehavior = value
	return f
}

func (f *FunctionScalaModel) WithReturnTypeValue(value tfconfig.Variable) *FunctionScalaModel {
	f.ReturnType = value
	return f
}

func (f *FunctionScalaModel) WithRuntimeVersionValue(value tfconfig.Variable) *FunctionScalaModel {
	f.RuntimeVersion = value
	return f
}

func (f *FunctionScalaModel) WithSchemaValue(value tfconfig.Variable) *FunctionScalaModel {
	f.Schema = value
	return f
}

func (f *FunctionScalaModel) WithSecretsValue(value tfconfig.Variable) *FunctionScalaModel {
	f.Secrets = value
	return f
}

func (f *FunctionScalaModel) WithTargetPathValue(value tfconfig.Variable) *FunctionScalaModel {
	f.TargetPath = value
	return f
}

func (f *FunctionScalaModel) WithTraceLevelValue(value tfconfig.Variable) *FunctionScalaModel {
	f.TraceLevel = value
	return f
}
