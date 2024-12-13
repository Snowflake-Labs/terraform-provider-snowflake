// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type ProcedureJavascriptResourceAssert struct {
	*assert.ResourceAssert
}

func ProcedureJavascriptResource(t *testing.T, name string) *ProcedureJavascriptResourceAssert {
	t.Helper()

	return &ProcedureJavascriptResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedProcedureJavascriptResource(t *testing.T, id string) *ProcedureJavascriptResourceAssert {
	t.Helper()

	return &ProcedureJavascriptResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (p *ProcedureJavascriptResourceAssert) HasArgumentsString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("arguments", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasCommentString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("comment", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasDatabaseString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("database", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasEnableConsoleOutputString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("enable_console_output", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasExecuteAsString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("execute_as", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasFullyQualifiedNameString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasIsSecureString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("is_secure", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasLogLevelString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("log_level", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasMetricLevelString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("metric_level", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNameString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("name", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNullInputBehaviorString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("null_input_behavior", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasProcedureDefinitionString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("procedure_definition", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasProcedureLanguageString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("procedure_language", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasReturnTypeString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("return_type", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasSchemaString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("schema", expected))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasTraceLevelString(expected string) *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("trace_level", expected))
	return p
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (p *ProcedureJavascriptResourceAssert) HasNoArguments() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("arguments.#", "0"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoComment() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("comment"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoDatabase() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("database"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoEnableConsoleOutput() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("enable_console_output"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoExecuteAs() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("execute_as"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoFullyQualifiedName() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoIsSecure() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("is_secure"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoLogLevel() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("log_level"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoMetricLevel() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("metric_level"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoName() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("name"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoNullInputBehavior() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("null_input_behavior"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoProcedureDefinition() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("procedure_definition"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoProcedureLanguage() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("procedure_language"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoReturnType() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("return_type"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoSchema() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("schema"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNoTraceLevel() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueNotSet("trace_level"))
	return p
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (p *ProcedureJavascriptResourceAssert) HasCommentEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("comment", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasDatabaseEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("database", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasExecuteAsEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("execute_as", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasFullyQualifiedNameEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("fully_qualified_name", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasIsSecureEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("is_secure", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasLogLevelEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("log_level", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasMetricLevelEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("metric_level", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasNameEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("name", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasNullInputBehaviorEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("null_input_behavior", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasProcedureDefinitionEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("procedure_definition", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasProcedureLanguageEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("procedure_language", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasReturnTypeEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("return_type", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasSchemaEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("schema", ""))
	return p
}
func (p *ProcedureJavascriptResourceAssert) HasTraceLevelEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValueSet("trace_level", ""))
	return p
}

///////////////////////////////
// Attribute presence checks //
///////////////////////////////

func (p *ProcedureJavascriptResourceAssert) HasArgumentsNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("arguments"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasCommentNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("comment"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasDatabaseNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("database"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasEnableConsoleOutputNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("enable_console_output"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasExecuteAsNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("execute_as"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasFullyQualifiedNameNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("fully_qualified_name"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasIsSecureNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("is_secure"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasLogLevelNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("log_level"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasMetricLevelNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("metric_level"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNameNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("name"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasNullInputBehaviorNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("null_input_behavior"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasProcedureDefinitionNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("procedure_definition"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasProcedureLanguageNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("procedure_language"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasReturnTypeNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("return_type"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasSchemaNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("schema"))
	return p
}

func (p *ProcedureJavascriptResourceAssert) HasTraceLevelNotEmpty() *ProcedureJavascriptResourceAssert {
	p.AddAssertion(assert.ValuePresent("trace_level"))
	return p
}
