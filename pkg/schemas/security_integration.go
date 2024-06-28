package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-1348100]: multiple PRs touching the security integrations are in progress, this should be filled by all the possible properties (the mapping method below should be too)
var SecurityIntegrationDescribeSchema = map[string]*schema.Schema{
	"todo": DescribePropertyListSchema,
}

func SecurityIntegrationsDescriptionsToSchema(descriptions []sdk.SecurityIntegrationProperty) map[string]any {
	securityIntegrationProperties := make(map[string]any)
	for _, desc := range descriptions {
		desc := desc
		propertySchema := SecurityIntegrationPropertyToSchema(&desc)
		securityIntegrationProperties["todo"] = []map[string]any{propertySchema}
	}
	return securityIntegrationProperties
}
