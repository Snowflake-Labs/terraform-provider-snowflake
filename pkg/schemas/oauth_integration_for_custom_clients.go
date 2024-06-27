package schemas

import (
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeOauthIntegrationForCustomClients = map[string]*schema.Schema{
	"oauth_client_type":                     DescribePropertyListSchema,
	"oauth_redirect_uri":                    DescribePropertyListSchema,
	"enabled":                               DescribePropertyListSchema,
	"oauth_allow_non_tls_redirect_uri":      DescribePropertyListSchema,
	"oauth_enforce_pkce":                    DescribePropertyListSchema,
	"oauth_use_secondary_roles":             DescribePropertyListSchema,
	"pre_authorized_roles_list":             DescribePropertyListSchema,
	"blocked_roles_list":                    DescribePropertyListSchema,
	"oauth_issue_refresh_tokens":            DescribePropertyListSchema,
	"oauth_refresh_token_validity":          DescribePropertyListSchema,
	"network_policy":                        DescribePropertyListSchema,
	"oauth_client_rsa_public_key_fp":        DescribePropertyListSchema,
	"oauth_client_rsa_public_key_2_fp":      DescribePropertyListSchema,
	"comment":                               DescribePropertyListSchema,
	"oauth_client_id":                       DescribePropertyListSchema,
	"oauth_authorization_endpoint":          DescribePropertyListSchema,
	"oauth_token_endpoint":                  DescribePropertyListSchema,
	"oauth_allowed_authorization_endpoints": DescribePropertyListSchema,
	"oauth_allowed_token_endpoints":         DescribePropertyListSchema,
}

func DescribeOauthIntegrationForCustomClientsToSchema(integrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	propsSchema := make(map[string]any)
	for _, property := range integrationProperties {
		switch property.Name {
		case "OAUTH_CLIENT_TYPE",
			"OAUTH_REDIRECT_URI",
			"ENABLED",
			"OAUTH_ALLOW_NON_TLS_REDIRECT_URI",
			"OAUTH_ENFORCE_PKCE",
			"OAUTH_USE_SECONDARY_ROLES",
			"PRE_AUTHORIZED_ROLES_LIST",
			"BLOCKED_ROLES_LIST",
			"OAUTH_ISSUE_REFRESH_TOKENS",
			"OAUTH_REFRESH_TOKEN_VALIDITY",
			"NETWORK_POLICY",
			"OAUTH_CLIENT_RSA_PUBLIC_KEY_FP",
			"OAUTH_CLIENT_RSA_PUBLIC_KEY_2_FP",
			"COMMENT",
			"OAUTH_CLIENT_ID",
			"OAUTH_AUTHORIZATION_ENDPOINT",
			"OAUTH_TOKEN_ENDPOINT",
			"OAUTH_ALLOWED_AUTHORIZATION_ENDPOINTS",
			"OAUTH_ALLOWED_TOKEN_ENDPOINTS":
			propsSchema[strings.ToLower(property.Name)] = []map[string]any{
				{
					"name":    property.Name,
					"type":    property.Type,
					"value":   property.Value,
					"default": property.Default,
				},
			}
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", property.Name)
		}
	}
	return propsSchema
}
