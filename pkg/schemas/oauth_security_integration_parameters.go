package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

var ParametersOauthIntegrationForPartnerApplicationsSchema = map[string]*schema.Schema{
	strings.ToLower(string(sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList)): ParameterListSchema,
}

func OauthIntegrationForPartnerApplicationsParametersSchema(parameters []*sdk.Parameter) map[string]any {
	oauthIntegrationForPartnerApplicationsParameters := make(map[string]any)
	for _, param := range parameters {
		switch key := strings.ToUpper(param.Key); key {
		case string(sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList):
			oauthIntegrationForPartnerApplicationsParameters[strings.ToLower(key)] = []map[string]any{ParameterToSchema(param)}
		}
	}
	return oauthIntegrationForPartnerApplicationsParameters
}
