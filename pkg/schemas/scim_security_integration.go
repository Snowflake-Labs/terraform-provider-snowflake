package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeScimSecurityIntegrationSchema represents output of DESCRIBE query for the single SecurityIntegration.
var DescribeScimSecurityIntegrationSchema = map[string]*schema.Schema{
	"enabled":        DescribePropertyListSchema,
	"network_policy": DescribePropertyListSchema,
	"run_as_role":    DescribePropertyListSchema,
	"sync_password":  DescribePropertyListSchema,
	"comment":        DescribePropertyListSchema,
}

var _ = DescribeScimSecurityIntegrationSchema

func ScimSecurityIntegrationPropertiesToSchema(securityIntegrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	securityIntegrationSchema := make(map[string]any)
	for _, securityIntegrationProperty := range securityIntegrationProperties {
		securityIntegrationProperty := securityIntegrationProperty
		switch securityIntegrationProperty.Name {
		case "ENABLED",
			"NETWORK_POLICY",
			"RUN_AS_ROLE",
			"SYNC_PASSWORD",
			"COMMENT":
			securityIntegrationSchema[strings.ToLower(securityIntegrationProperty.Name)] = []map[string]any{SecurityIntegrationPropertyToSchema(&securityIntegrationProperty)}
		}
	}
	return securityIntegrationSchema
}

var _ = ScimSecurityIntegrationPropertiesToSchema
