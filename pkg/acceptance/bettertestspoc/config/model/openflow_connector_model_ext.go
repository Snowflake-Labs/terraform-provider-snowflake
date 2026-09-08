package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

// WithFrom creates the connector from a Snowflake-managed definition. `from` is a required nested list, so the
// generated constructor passes it positionally and calls this; the generator cannot build setters for nested
// lists itself. The parameter is the definition rather than a stage because no single type carries both
// sources, and the definition is the common case. Use WithFromStage for the other one.
func (o *OpenflowConnectorModel) WithFrom(definition string) *OpenflowConnectorModel {
	return o.WithFromValue(
		tfconfig.ListVariable(
			tfconfig.MapVariable(map[string]tfconfig.Variable{
				"definition": tfconfig.StringVariable(definition),
			}),
		),
	)
}

// WithFromStage creates the connector from a bundle on a stage. An empty path is left out rather than written
// as an empty string, so the configuration matches what a user would write.
func (o *OpenflowConnectorModel) WithFromStage(stageId sdk.SchemaObjectIdentifier, path string) *OpenflowConnectorModel {
	from := map[string]tfconfig.Variable{
		"stage": tfconfig.StringVariable(stageId.FullyQualifiedName()),
	}
	if path != "" {
		from["path"] = tfconfig.StringVariable(path)
	}
	return o.WithFromValue(tfconfig.ListVariable(tfconfig.MapVariable(from)))
}
