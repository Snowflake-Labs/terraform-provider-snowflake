package schemas

import (
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
