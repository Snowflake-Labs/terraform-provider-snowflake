// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type FunctionPythonResourceAssert struct {
	*assert.ResourceAssert
}

func FunctionPythonResource(t *testing.T, name string) *FunctionPythonResourceAssert {
	t.Helper()

	return &FunctionPythonResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedFunctionPythonResource(t *testing.T, id string) *FunctionPythonResourceAssert {
	t.Helper()

	return &FunctionPythonResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (f *FunctionPythonResourceAssert) HasArgumentsString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("arguments", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasCommentString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("comment", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasDatabaseString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("database", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasEnableConsoleOutputString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("enable_console_output", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasExternalAccessIntegrationsString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("external_access_integrations", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasFullyQualifiedNameString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasFunctionDefinitionString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("function_definition", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasFunctionLanguageString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("function_language", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasHandlerString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("handler", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasImportsString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("imports", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasIsAggregateString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("is_aggregate", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasIsSecureString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("is_secure", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasLogLevelString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("log_level", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasMetricLevelString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("metric_level", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasNameString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("name", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasNullInputBehaviorString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("null_input_behavior", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasPackagesString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("packages", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasReturnBehaviorString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("return_behavior", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasReturnTypeString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("return_type", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasRuntimeVersionString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("runtime_version", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasSchemaString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("schema", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasSecretsString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("secrets", expected))
	return f
}

func (f *FunctionPythonResourceAssert) HasTraceLevelString(expected string) *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueSet("trace_level", expected))
	return f
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (f *FunctionPythonResourceAssert) HasNoArguments() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("arguments"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoComment() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("comment"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoDatabase() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("database"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoEnableConsoleOutput() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("enable_console_output"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoExternalAccessIntegrations() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("external_access_integrations"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoFullyQualifiedName() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoFunctionDefinition() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("function_definition"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoFunctionLanguage() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("function_language"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoHandler() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("handler"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoImports() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("imports"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoIsAggregate() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("is_aggregate"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoIsSecure() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("is_secure"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoLogLevel() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("log_level"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoMetricLevel() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("metric_level"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoName() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("name"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoNullInputBehavior() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("null_input_behavior"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoPackages() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("packages"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoReturnBehavior() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("return_behavior"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoReturnType() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("return_type"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoRuntimeVersion() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("runtime_version"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoSchema() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("schema"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoSecrets() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("secrets"))
	return f
}

func (f *FunctionPythonResourceAssert) HasNoTraceLevel() *FunctionPythonResourceAssert {
	f.AddAssertion(assert.ValueNotSet("trace_level"))
	return f
}