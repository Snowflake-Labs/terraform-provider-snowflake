package schemas

import (
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var (
	SecurityIntegrationDescribeSchema = collections.MergeMaps(
		DescribeApiAuthSecurityIntegrationSchema,
		DescribeExternalOauthSecurityIntegrationSchema,
		DescribeOauthIntegrationForCustomClients,
		DescribeOauthIntegrationForPartnerApplications,
		DescribeSaml2IntegrationSchema,
		DescribeScimSecurityIntegrationSchema,
	)
	allSecurityIntegrationPropertiesNames = helpers.ConcatSlices(
		ApiAuthenticationPropertiesNames,
		ExternalOauthPropertiesNames,
		OauthIntegrationForCustomClientsPropertiesNames,
		OauthIntegrationForPartnerApplicationsPropertiesNames,
		Saml2PropertiesNames,
		ScimPropertiesNames,
	)
)

func SecurityIntegrationsDescriptionsToSchema(integrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	securityIntegrationProperties := make(map[string]any)
	for _, desc := range integrationProperties {
		desc := desc
		if slices.Contains(allSecurityIntegrationPropertiesNames, desc.Name) {
			securityIntegrationProperties[strings.ToLower(desc.Name)] = []map[string]any{SecurityIntegrationPropertyToSchema(&desc)}
		} else {
			log.Printf("[WARN] unexpected property %v in security integration returned from Snowflake", desc.Name)
		}
	}
	return securityIntegrationProperties
}
