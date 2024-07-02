package schemas

import (
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeApiAuthSecurityIntegrationSchema represents output of DESCRIBE query for the single SecurityIntegration.
var DescribeApiAuthSecurityIntegrationSchema = map[string]*schema.Schema{
	"enabled":                      DescribePropertyListSchema,
	"oauth_access_token_validity":  DescribePropertyListSchema,
	"oauth_refresh_token_validity": DescribePropertyListSchema,
	"oauth_client_id":              DescribePropertyListSchema,
	"oauth_client_auth_method":     DescribePropertyListSchema,
	"oauth_authorization_endpoint": DescribePropertyListSchema,
	"oauth_token_endpoint":         DescribePropertyListSchema,
	"oauth_allowed_scopes":         DescribePropertyListSchema,
	"oauth_grant":                  DescribePropertyListSchema,
	"parent_integration":           DescribePropertyListSchema,
	"auth_type":                    DescribePropertyListSchema,
	"comment":                      DescribePropertyListSchema,
}

var _ = DescribeApiAuthSecurityIntegrationSchema

func ApiAuthSecurityIntegrationPropertiesToSchema(securityIntegrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	securityIntegrationSchema := make(map[string]any)
	for _, securityIntegrationProperty := range securityIntegrationProperties {
		switch securityIntegrationProperty.Name {
		case "ENABLED",
			"OAUTH_ACCESS_TOKEN_VALIDITY",
			"OAUTH_REFRESH_TOKEN_VALIDITY",
			"OAUTH_CLIENT_ID",
			"OAUTH_CLIENT_AUTH_METHOD",
			"OAUTH_AUTHORIZATION_ENDPOINT",
			"OAUTH_TOKEN_ENDPOINT",
			"OAUTH_ALLOWED_SCOPES",
			"OAUTH_GRANT",
			"PARENT_INTEGRATION",
			"AUTH_TYPE",
			"COMMENT":
			securityIntegrationSchema[strings.ToLower(securityIntegrationProperty.Name)] = []map[string]any{SecurityIntegrationPropertyToSchema(&securityIntegrationProperty)}
		default:
			log.Printf("unknown field from DESCRIBE: %v", securityIntegrationProperty.Name)
		}
	}
	return securityIntegrationSchema
}

var _ = ApiAuthSecurityIntegrationPropertiesToSchema
