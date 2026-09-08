package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HasFromDefinition asserts the whole `from` block for a connector created from a definition. The generator
// only produces a presence check for a nested list, and the point of asserting the block is that the source
// not in use is empty.
func (o *OpenflowConnectorResourceAssert) HasFromDefinition(definition string) *OpenflowConnectorResourceAssert {
	o.ValueSet("from.#", "1")
	o.ValueSet("from.0.definition", definition)
	o.ValueSet("from.0.stage", "")
	o.ValueSet("from.0.path", "")
	return o
}

// HasFromEmpty asserts the block is absent, which is what an imported connector has: SHOW cannot tell which
// source created it. The generator stopped emitting this once `from` became required.
func (o *OpenflowConnectorResourceAssert) HasFromEmpty() *OpenflowConnectorResourceAssert {
	o.ValueSet("from.#", "0")
	return o
}

// HasFromStage asserts the whole `from` block for a connector created from a bundle on a stage. The definition
// stays empty even though Snowflake resolves one, because the block is not read back.
func (o *OpenflowConnectorResourceAssert) HasFromStage(stageId sdk.SchemaObjectIdentifier, path string) *OpenflowConnectorResourceAssert {
	o.ValueSet("from.#", "1")
	o.ValueSet("from.0.definition", "")
	o.ValueSet("from.0.stage", stageId.FullyQualifiedName())
	o.ValueSet("from.0.path", path)
	return o
}
