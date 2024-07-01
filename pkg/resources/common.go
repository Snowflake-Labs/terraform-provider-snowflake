package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DiffSuppressStatement will suppress diffs between statements if they differ in only case or in
// runs of whitespace (\s+ = \s). This is needed because the snowflake api does not faithfully
// round-trip queries, so we cannot do a simple character-wise comparison to detect changes.
//
// Warnings: We will have false positives in cases where a change in case or run of whitespace is
// semantically significant.
//
// If we can find a sql parser that can handle the snowflake dialect then we should switch to parsing
// queries and either comparing ASTs or emitting a canonical serialization for comparison. I couldn't
// find such a library.
func DiffSuppressStatement(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(normalizeQuery(old), normalizeQuery(new))
}

func normalizeQuery(str string) string {
	return strings.TrimSpace(space.ReplaceAllString(str, " "))
}

// TODO [SNOW-999049]: address during identifiers rework
func suppressIdentifierQuoting(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	if oldValue == "" || newValue == "" {
		return false
	} else {
		oldId, err := helpers.DecodeSnowflakeParameterID(oldValue)
		if err != nil {
			return false
		}
		newId, err := helpers.DecodeSnowflakeParameterID(newValue)
		if err != nil {
			return false
		}
		return oldId.FullyQualifiedName() == newId.FullyQualifiedName()
	}
}

// TODO [SNOW-1325214]: address during stage resource rework
func suppressQuoting(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	if oldValue == "" || newValue == "" {
		return false
	} else {
		oldWithoutQuotes := strings.ReplaceAll(oldValue, "'", "")
		newWithoutQuotes := strings.ReplaceAll(newValue, "'", "")
		return oldWithoutQuotes == newWithoutQuotes
	}
}

var apiAuthCommonSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier (i.e. name) for the integration. This value must be unique in your account.",
	},
	"enabled": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Specifies whether this security integration is enabled or disabled.",
	},
	"oauth_token_endpoint": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "unknown",
		Description: "Specifies the token endpoint used by the client to obtain an access token by presenting its authorization grant or refresh token. The token endpoint is used with every authorization grant except for the implicit grant type (since an access token is issued directly).",
	},
	"oauth_client_auth_method": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption),
		Default:          "unknown",
		Description:      fmt.Sprintf("Specifies that POST is used as the authentication method to the external service. Valid options are: %v", sdk.AsStringList(sdk.AllApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption)),
	},
	"oauth_client_id": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the client ID for the OAuth application in the external service.",
	},
	"oauth_client_secret": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the client secret for the OAuth application in the ServiceNow instance from the previous step. The connector uses this to request an access token from the ServiceNow instance.",
	},
	"oauth_access_token_validity": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateFunc:     validation.IntAtLeast(0),
		Default:          -1,
		Description:      "Specifies the default lifetime of the OAuth access token (in seconds) issued by an OAuth server.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("oauth_access_token_validity"),
	},
	"oauth_refresh_token_validity": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateFunc:     validation.IntAtLeast(1),
		Default:          -1,
		Description:      "Specifies the default lifetime of the OAuth access token (in seconds) issued by an OAuth server.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("oauth_refresh_token_validity"),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the integration.",
	},
	showOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SECURITY INTEGRATIONS` for the given security integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSecurityIntegrationSchema,
		},
	},
	describeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SECURITY INTEGRATIONS` for the given security integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeApiAuthSecurityIntegrationSchema,
		},
	},
}

func listValueToSlice(value string, trimBrackets bool, trimQuotes bool) []string {
	if trimBrackets {
		value = strings.TrimLeft(value, "[")
		value = strings.TrimRight(value, "]")
	}
	if value == "" {
		return nil
	}
	elems := strings.Split(value, ",")
	for i := range elems {
		if trimQuotes {
			elems[i] = strings.Trim(elems[i], " '")
		}
	}
	return elems
}

func ctyValToSliceString(val cty.Value) []string {
	valueElems := val.AsValueSlice()
	elems := make([]string, len(valueElems))
	for i, v := range valueElems {
		elems[i] = v.AsString()
	}
	return elems
}
