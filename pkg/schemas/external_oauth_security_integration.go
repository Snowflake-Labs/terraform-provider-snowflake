package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeExternalOauthSecurityIntegrationSchema represents output of DESCRIBE query for the single SecurityIntegration.
var DescribeExternalOauthSecurityIntegrationSchema = map[string]*schema.Schema{
	"enabled":                                         DescribePropertyListSchema,
	"external_oauth_issuer":                           DescribePropertyListSchema,
	"external_oauth_jws_keys_url":                     DescribePropertyListSchema,
	"external_oauth_any_role_mode":                    DescribePropertyListSchema,
	"external_oauth_rsa_public_key":                   DescribePropertyListSchema,
	"external_oauth_rsa_public_key_2":                 DescribePropertyListSchema,
	"external_oauth_blocked_roles_list":               DescribePropertyListSchema,
	"external_oauth_allowed_roles_list":               DescribePropertyListSchema,
	"external_oauth_audience_list":                    DescribePropertyListSchema,
	"external_oauth_token_user_mapping_claim":         DescribePropertyListSchema,
	"external_oauth_snowflake_user_mapping_attribute": DescribePropertyListSchema,
	"external_oauth_scope_delimiter":                  DescribePropertyListSchema,
	"comment":                                         DescribePropertyListSchema,
}

var _ = DescribeExternalOauthSecurityIntegrationSchema

func ExternalOauthSecurityIntegrationPropertiesToSchema(securityIntegrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	securityIntegrationSchema := make(map[string]any)
	for _, securityIntegrationProperty := range securityIntegrationProperties {
		securityIntegrationProperty := securityIntegrationProperty
		switch securityIntegrationProperty.Name {
		case "ENABLED",
			"EXTERNAL_OAUTH_ISSUER",
			"EXTERNAL_OAUTH_JWS_KEYS_URL",
			"EXTERNAL_OAUTH_ANY_ROLE_MODE",
			"EXTERNAL_OAUTH_RSA_PUBLIC_KEY",
			"EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2",
			"EXTERNAL_OAUTH_BLOCKED_ROLES_LIST",
			"EXTERNAL_OAUTH_ALLOWED_ROLES_LIST",
			"EXTERNAL_OAUTH_AUDIENCE_LIST",
			"EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM",
			"EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE",
			"EXTERNAL_OAUTH_SCOPE_DELIMITER",
			"COMMENT":
			securityIntegrationSchema[strings.ToLower(securityIntegrationProperty.Name)] = []map[string]any{SecurityIntegrationPropertyToSchema(&securityIntegrationProperty)}
		}
	}
	return securityIntegrationSchema
}

var _ = ExternalOauthSecurityIntegrationPropertiesToSchema
