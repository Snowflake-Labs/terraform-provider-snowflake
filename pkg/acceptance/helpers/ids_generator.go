package helpers

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type IdsGenerator struct {
	context *TestClientContext
}

func NewIdsGenerator(context *TestClientContext) *IdsGenerator {
	return &IdsGenerator{
		context: context,
	}
}

func (c *IdsGenerator) databaseId() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.context.database)
}

func (c *IdsGenerator) schemaId() sdk.DatabaseObjectIdentifier {
	return sdk.NewDatabaseObjectIdentifier(c.context.database, c.context.schema)
}

func (c *IdsGenerator) warehouseId() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.context.warehouse)
}

func (c *IdsGenerator) newSchemaObjectIdentifier(name string) sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifier(c.context.database, c.context.schema, name)
}

func (c *IdsGenerator) RandomAccountObjectIdentifier() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.Alpha())
}

func (c *IdsGenerator) RandomSchemaObjectIdentifier() sdk.SchemaObjectIdentifier {
	return c.newSchemaObjectIdentifier(c.Alpha())
}

func (c *IdsGenerator) Alpha() string {
	return c.AlphaN(6)
}

func (c *IdsGenerator) AlphaN(n int) string {
	return c.withTestObjectSuffix(strings.ToUpper(random.AlphaN(n)))
}

func (c *IdsGenerator) AlphaContaining(part string) string {
	return c.withTestObjectSuffix(c.Alpha() + part)
}

func (c *IdsGenerator) AlphaWithPrefix(part string) string {
	return c.withTestObjectSuffix(part + c.Alpha())
}

func (c *IdsGenerator) withTestObjectSuffix(text string) string {
	return text + c.context.testObjectSuffix
}
