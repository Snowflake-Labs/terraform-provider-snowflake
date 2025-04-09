package helpers

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"

type RandomValuesGenerator struct {
	context *TestClientContext
}

func NewRandomValuesGenerator(context *TestClientContext) *RandomValuesGenerator {
	return &RandomValuesGenerator{
		context: context,
	}
}

func (c *RandomValuesGenerator) Secret() string {
	return c.context.generatedRandomSecret + random.String()
}
