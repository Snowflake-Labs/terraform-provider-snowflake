package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var (
	SecurityIntegrationDescribeSchema     = helpers.MergeMaps(DescribeSaml2IntegrationSchema, DescribeScimSecurityIntegrationSchema)
	allSecurityIntegrationPropertiesNames = append(Saml2PropertiesNames, ScimPropertiesNames...)
)

func SecurityIntegrationsDescriptionsToSchema(descriptions []sdk.SecurityIntegrationProperty) map[string]any {
	securityIntegrationProperties := make(map[string]any)
	for _, desc := range descriptions {
		desc := desc
		if slices.Contains(allSecurityIntegrationPropertiesNames, desc.Name) {
			securityIntegrationProperties[strings.ToLower(desc.Name)] = []map[string]any{SecurityIntegrationPropertyToSchema(&desc)}
		}
	}
	return securityIntegrationProperties
}
