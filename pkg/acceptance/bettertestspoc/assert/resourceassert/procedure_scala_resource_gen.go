// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type ProcedureScalaResourceAssert struct {
	*assert.ResourceAssert
}

func ProcedureScalaResource(t *testing.T, name string) *ProcedureScalaResourceAssert {
	t.Helper()

	return &ProcedureScalaResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedProcedureScalaResource(t *testing.T, id string) *ProcedureScalaResourceAssert {
	t.Helper()

	return &ProcedureScalaResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (p *ProcedureScalaResourceAssert) HasArgumentsString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("arguments", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasCommentString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("comment", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasDatabaseString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("database", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasEnableConsoleOutputString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("enable_console_output", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasExecuteAsString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("execute_as", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasExternalAccessIntegrationsString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("external_access_integrations", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasFullyQualifiedNameString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasHandlerString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("handler", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasImportsString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("imports", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasIsSecureString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("is_secure", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasLogLevelString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("log_level", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasMetricLevelString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("metric_level", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNameString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("name", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNullInputBehaviorString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("null_input_behavior", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasPackagesString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("packages", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasProcedureDefinitionString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("procedure_definition", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasProcedureLanguageString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("procedure_language", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasReturnTypeString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("return_type", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasRuntimeVersionString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("runtime_version", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasSchemaString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("schema", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasSecretsString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("secrets", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasSnowparkPackageString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("snowpark_package", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasTargetPathString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("target_path", expected))
	return p
}

func (p *ProcedureScalaResourceAssert) HasTraceLevelString(expected string) *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueSet("trace_level", expected))
	return p
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (p *ProcedureScalaResourceAssert) HasNoArguments() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("arguments"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoComment() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("comment"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoDatabase() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("database"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoEnableConsoleOutput() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("enable_console_output"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoExecuteAs() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("execute_as"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoExternalAccessIntegrations() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("external_access_integrations"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoFullyQualifiedName() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoHandler() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("handler"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoImports() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("imports"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoIsSecure() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("is_secure"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoLogLevel() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("log_level"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoMetricLevel() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("metric_level"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoName() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("name"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoNullInputBehavior() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("null_input_behavior"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoPackages() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("packages"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoProcedureDefinition() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("procedure_definition"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoProcedureLanguage() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("procedure_language"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoReturnType() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("return_type"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoRuntimeVersion() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("runtime_version"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoSchema() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("schema"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoSecrets() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("secrets"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoSnowparkPackage() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("snowpark_package"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoTargetPath() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("target_path"))
	return p
}

func (p *ProcedureScalaResourceAssert) HasNoTraceLevel() *ProcedureScalaResourceAssert {
	p.AddAssertion(assert.ValueNotSet("trace_level"))
	return p
}
