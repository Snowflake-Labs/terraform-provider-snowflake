// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type FunctionSqlResourceAssert struct {
	*assert.ResourceAssert
}

func FunctionSqlResource(t *testing.T, name string) *FunctionSqlResourceAssert {
	t.Helper()

	return &FunctionSqlResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedFunctionSqlResource(t *testing.T, id string) *FunctionSqlResourceAssert {
	t.Helper()

	return &FunctionSqlResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (f *FunctionSqlResourceAssert) HasArgumentsString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("arguments", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasCommentString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("comment", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasDatabaseString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("database", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasEnableConsoleOutputString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("enable_console_output", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasFullyQualifiedNameString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasFunctionDefinitionString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("function_definition", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasFunctionLanguageString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("function_language", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasIsSecureString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("is_secure", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasLogLevelString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("log_level", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasMetricLevelString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("metric_level", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasNameString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("name", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasReturnResultsBehaviorString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("return_results_behavior", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasReturnTypeString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("return_type", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasSchemaString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("schema", expected))
	return f
}

func (f *FunctionSqlResourceAssert) HasTraceLevelString(expected string) *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("trace_level", expected))
	return f
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (f *FunctionSqlResourceAssert) HasNoArguments() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("arguments.#", "0"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoComment() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("comment"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoDatabase() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("database"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoEnableConsoleOutput() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("enable_console_output"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoFullyQualifiedName() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoFunctionDefinition() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("function_definition"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoFunctionLanguage() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("function_language"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoIsSecure() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("is_secure"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoLogLevel() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("log_level"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoMetricLevel() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("metric_level"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoName() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("name"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoReturnResultsBehavior() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("return_results_behavior"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoReturnType() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("return_type"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoSchema() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("schema"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNoTraceLevel() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueNotSet("trace_level"))
	return f
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (f *FunctionSqlResourceAssert) HasCommentEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("comment", ""))
	return f
}

func (f *FunctionSqlResourceAssert) HasFullyQualifiedNameEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("fully_qualified_name", ""))
	return f
}

func (f *FunctionSqlResourceAssert) HasFunctionLanguageEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("function_language", ""))
	return f
}

func (f *FunctionSqlResourceAssert) HasIsSecureEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("is_secure", ""))
	return f
}

func (f *FunctionSqlResourceAssert) HasLogLevelEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("log_level", ""))
	return f
}

func (f *FunctionSqlResourceAssert) HasMetricLevelEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("metric_level", ""))
	return f
}

func (f *FunctionSqlResourceAssert) HasReturnResultsBehaviorEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("return_results_behavior", ""))
	return f
}

func (f *FunctionSqlResourceAssert) HasTraceLevelEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValueSet("trace_level", ""))
	return f
}

///////////////////////////////
// Attribute presence checks //
///////////////////////////////

func (f *FunctionSqlResourceAssert) HasArgumentsNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("arguments"))
	return f
}

func (f *FunctionSqlResourceAssert) HasCommentNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("comment"))
	return f
}

func (f *FunctionSqlResourceAssert) HasDatabaseNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("database"))
	return f
}

func (f *FunctionSqlResourceAssert) HasEnableConsoleOutputNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("enable_console_output"))
	return f
}

func (f *FunctionSqlResourceAssert) HasFullyQualifiedNameNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("fully_qualified_name"))
	return f
}

func (f *FunctionSqlResourceAssert) HasFunctionDefinitionNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("function_definition"))
	return f
}

func (f *FunctionSqlResourceAssert) HasFunctionLanguageNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("function_language"))
	return f
}

func (f *FunctionSqlResourceAssert) HasIsSecureNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("is_secure"))
	return f
}

func (f *FunctionSqlResourceAssert) HasLogLevelNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("log_level"))
	return f
}

func (f *FunctionSqlResourceAssert) HasMetricLevelNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("metric_level"))
	return f
}

func (f *FunctionSqlResourceAssert) HasNameNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("name"))
	return f
}

func (f *FunctionSqlResourceAssert) HasReturnResultsBehaviorNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("return_results_behavior"))
	return f
}

func (f *FunctionSqlResourceAssert) HasReturnTypeNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("return_type"))
	return f
}

func (f *FunctionSqlResourceAssert) HasSchemaNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("schema"))
	return f
}

func (f *FunctionSqlResourceAssert) HasTraceLevelNotEmpty() *FunctionSqlResourceAssert {
	f.AddAssertion(assert.ValuePresent("trace_level"))
	return f
}
