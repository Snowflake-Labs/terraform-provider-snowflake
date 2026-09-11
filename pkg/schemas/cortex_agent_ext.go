package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	profile := cortexAgentProfileSchema()
	ShowCortexAgentSchema["profile"] = profile
	DescribeCortexAgentDetailsSchema["profile"] = profile
}

func cortexAgentProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"avatar": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"color": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func cortexAgentProfileToSchema(profile sdk.CortexAgentProfile) []map[string]any {
	return []map[string]any{
		{
			"display_name": profile.DisplayName,
			"avatar":       profile.Avatar,
			"color":        profile.Color,
		},
	}
}

func CortexAgentToSchemaWithProfile(cortexAgent *sdk.CortexAgent) map[string]any {
	cortexAgentSchema := CortexAgentToSchema(cortexAgent)
	cortexAgentSchema["profile"] = cortexAgentProfileToSchema(cortexAgent.Profile)
	return cortexAgentSchema
}

func CortexAgentDetailsToSchemaWithProfile(cortexAgentDetails *sdk.CortexAgentDetails) map[string]any {
	cortexAgentDetailsSchema := CortexAgentDetailsToSchema(cortexAgentDetails)
	cortexAgentDetailsSchema["profile"] = cortexAgentProfileToSchema(cortexAgentDetails.Profile)
	return cortexAgentDetailsSchema
}
