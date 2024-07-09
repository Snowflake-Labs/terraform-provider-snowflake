package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var streamlitsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC STREAMLIT for each streamlit returned by SHOW STREAMLITS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"streamlits": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all streamlits details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW STREAMLITS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowStreamlitSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE STREAMLITS.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeStreamlitSchema,
					},
				},
			},
		},
	},
}

func Streamlits() *schema.Resource {
	return &schema.Resource{
		ReadContext: ReadStreamlits,
		Schema:      streamlitsSchema,
		Description: "Datasource used to get details of filtered streamlits. Filtering is aligned with the current possibilities for [SHOW STREAMLITS](https://docs.snowflake.com/en/sql-reference/sql/show-integrations) query (only `like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection `security_integrations`.",
	}
}

func ReadStreamlits(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	showRequest := sdk.NewShowStreamlitRequest()

	if likePattern, ok := d.GetOk("like"); ok {
		showRequest.WithLike(sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		})
	}

	streamlits, err := client.Streamlits.Show(ctx, showRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("security_integrations_read")

	flattenedStreamlits := make([]map[string]any, len(streamlits))

	for i, streamlit := range streamlits {
		streamlit := streamlit
		var streamlitDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			descriptions, err := client.Streamlits.Describe(ctx, streamlit.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			streamlitDescriptions = make([]map[string]any, 1)
			streamlitDescriptions[0], err = schemas.StreamlitPropertiesToSchema(*descriptions)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		flattenedStreamlits[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.StreamlitToSchema(&streamlit)},
			resources.DescribeOutputAttributeName: streamlitDescriptions,
		}
	}

	err = d.Set("streamlits", flattenedStreamlits)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
