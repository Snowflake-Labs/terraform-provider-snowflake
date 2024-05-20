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

func (c *IdsGenerator) DatabaseId() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.context.database)
}

func (c *IdsGenerator) SchemaId() sdk.DatabaseObjectIdentifier {
	return sdk.NewDatabaseObjectIdentifier(c.context.database, c.context.schema)
}

func (c *IdsGenerator) WarehouseId() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.context.warehouse)
}

func (c *IdsGenerator) AccountIdentifierWithLocator() sdk.AccountIdentifier {
	return sdk.NewAccountIdentifierFromAccountLocator(c.context.client.GetAccountLocator())
}

func (c *IdsGenerator) NewSchemaObjectIdentifier(name string) sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifierInSchema(c.SchemaId(), name)
}

func (c *IdsGenerator) RandomAccountObjectIdentifier() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.Alpha())
}

func (c *IdsGenerator) RandomAccountObjectIdentifierWithPrefix(prefix string) sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.AlphaWithPrefix(prefix))
}

func (c *IdsGenerator) RandomAccountObjectIdentifierContaining(part string) sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(c.AlphaContaining(part))
}

func (c *IdsGenerator) RandomDatabaseObjectIdentifier() sdk.DatabaseObjectIdentifier {
	return sdk.NewDatabaseObjectIdentifier(c.DatabaseId().Name(), c.Alpha())
}

func (c *IdsGenerator) RandomSchemaObjectIdentifier() sdk.SchemaObjectIdentifier {
	return c.RandomSchemaObjectIdentifierInSchema(c.SchemaId())
}

func (c *IdsGenerator) RandomSchemaObjectIdentifierInSchema(schemaId sdk.DatabaseObjectIdentifier) sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifierInSchema(schemaId, c.Alpha())
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

func (c *IdsGenerator) AlphaWithPrefix(prefix string) string {
	return prefix + c.Alpha()
}

func (c *IdsGenerator) withTestObjectSuffix(text string) string {
	return text + c.context.testObjectSuffix
}
