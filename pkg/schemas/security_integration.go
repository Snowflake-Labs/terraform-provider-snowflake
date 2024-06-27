package schemas

import (
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowSecurityIntegrationSchema represents output of SHOW query for the single SecurityIntegration.
var ShowSecurityIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"integration_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"category": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowSecurityIntegrationSchema

func SecurityIntegrationToSchema(securityIntegration *sdk.SecurityIntegration) map[string]any {
	securityIntegrationSchema := make(map[string]any)
	securityIntegrationSchema["name"] = securityIntegration.Name
	securityIntegrationSchema["integration_type"] = securityIntegration.IntegrationType
	securityIntegrationSchema["category"] = securityIntegration.Category
	securityIntegrationSchema["enabled"] = securityIntegration.Enabled
	securityIntegrationSchema["comment"] = securityIntegration.Comment
	securityIntegrationSchema["created_on"] = securityIntegration.CreatedOn.String()
	return securityIntegrationSchema
}

var _ = SecurityIntegrationToSchema

var DescribeSaml2IntegrationSchema = map[string]*schema.Schema{
	"saml2_issuer": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_sso_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_provider": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_x509_cert": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_sp_initiated_login_page_label": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_enable_sp_initiated": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"saml2_snowflake_x509_cert": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_sign_request": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"saml2_requested_nameid_format": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_post_logout_redirect_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_force_authn": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"saml2_snowflake_issuer_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_snowflake_acs_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_snowflake_metadata": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_digest_methods_used": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"saml2_signature_methods_used": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"allowed_user_domains": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Computed: true,
	},
	"allowed_email_patterns": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func DescribeSaml2IntegrationToSchema(props []sdk.SecurityIntegrationProperty) map[string]any {
	propsSchema := make(map[string]any)
	for _, property := range props {
		name := property.Name
		value := property.Value
		switch name {
		case "ENABLED", "COMMENT", "SAML2_ISSUER", "SAML2_SSO_URL", "SAML2_PROVIDER", "SAML2_X509_CERT", "SAML2_SP_INITIATED_LOGIN_PAGE_LABEL", "SAML2_SNOWFLAKE_X509_CERT",
			"SAML2_REQUESTED_NAMEID_FORMAT", "SAML2_POST_LOGOUT_REDIRECT_URL", "SAML2_SNOWFLAKE_ISSUER_URL", "SAML2_SNOWFLAKE_ACS_URL", "SAML2_SNOWFLAKE_METADATA",
			"SAML2_DIGEST_METHODS_USED", "SAML2_SIGNATURE_METHODS_USED":
			propsSchema[strings.ToLower(name)] = value
		case "SAML2_ENABLE_SP_INITIATED", "SAML2_SIGN_REQUEST", "SAML2_FORCE_AUTHN":
			propsSchema[strings.ToLower(name)] = helpers.StringToBool(value)
		case "ALLOWED_USER_DOMAINS", "ALLOWED_EMAIL_PATTERNS":
			// TODO: extract
			value = strings.TrimLeft(value, "[")
			value = strings.TrimRight(value, "]")
			elems := strings.Split(value, ",")
			if value == "" {
				elems = nil
			}
			propsSchema[strings.ToLower(name)] = elems
		default:
			log.Printf("[WARN] unexpected property %v returned from Snowflake", name)
		}
	}
	return propsSchema
}
