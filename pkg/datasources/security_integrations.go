package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var securityIntegrationsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC SECURITY INTEGRATION for each security integration returned by SHOW SECURITY INTEGRATIONS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"security_integrations": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all security integrations details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"show_output": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW SECURITY INTEGRATIONS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowSecurityIntegrationSchema,
					},
				},
				"describe_output": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE SECURITY INTEGRATIONS.",
					Elem: &schema.Resource{
						Schema: schemas.SecurityIntegrationDescribeSchema,
					},
				},
			},
		},
	},
}

func SecurityIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: ReadSecurityIntegrations,
		Schema:      securityIntegrationsSchema,
		Description: "Datasource used to get details of filtered security integrations. Filtering is aligned with the current possibilities for [SHOW SECURITY INTEGRATIONS](https://docs.snowflake.com/en/sql-reference/sql/show-integrations) query (only `like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection `security_integrations`.",
	}
}

func ReadSecurityIntegrations(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	showRequest := sdk.NewShowSecurityIntegrationRequest()

	if likePattern, ok := d.GetOk("like"); ok {
		showRequest.WithLike(sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		})
	}

	securityIntegrations, err := client.SecurityIntegrations.Show(ctx, showRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("security_integrations_read")

	flattenedSecurityIntegrations := make([]map[string]any, len(securityIntegrations))

	for i, securityIntegration := range securityIntegrations {
		securityIntegration := securityIntegration
		var securityIntegrationDescription map[string]any
		if d.Get("with_describe").(bool) {
			descriptions, err := client.SecurityIntegrations.Describe(ctx, securityIntegration.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			securityIntegrationDescription = schemas.SecurityIntegrationsDescriptionsToSchema(descriptions)
		}

		flattenedSecurityIntegrations[i] = map[string]any{
			"show_output":     []map[string]any{schemas.SecurityIntegrationToSchema(&securityIntegration)},
			"describe_output": []map[string]any{securityIntegrationDescription},
		}
	}

	err = d.Set("security_integrations", flattenedSecurityIntegrations)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
