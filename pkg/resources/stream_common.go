package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func handleStreamTimeTravel(d *schema.ResourceData, req *sdk.CreateOnTableStreamRequest) {
	if v := d.Get(AtAttributeName).([]any); len(v) > 0 {
		req.WithOn(*sdk.NewOnStreamRequest().WithAt(true).WithStatement(handleStreamTimeTravelStatement(v[0].(map[string]any))))
	}
	if v := d.Get(BeforeAttributeName).([]any); len(v) > 0 {
		req.WithOn(*sdk.NewOnStreamRequest().WithBefore(true).WithStatement(handleStreamTimeTravelStatement(v[0].(map[string]any))))
	}
}

func handleStreamTimeTravelStatement(timeTravelConfig map[string]any) sdk.OnStreamStatementRequest {
	statement := sdk.OnStreamStatementRequest{}
	if v := timeTravelConfig["timestamp"].(string); len(v) > 0 {
		statement.WithTimestamp(v)
	}
	if v := timeTravelConfig["offset"].(string); len(v) > 0 {
		statement.WithOffset(v)
	}
	if v := timeTravelConfig["statement"].(string); len(v) > 0 {
		statement.WithStatement(v)
	}
	if v := timeTravelConfig["stream"].(string); len(v) > 0 {
		statement.WithStream(v)
	}
	return statement
}

func DeleteStreamContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id).WithIfExists(true))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting stream",
				Detail:   fmt.Sprintf("id %v err = %v", id.Name(), err),
			},
		}
	}

	d.SetId("")
	return nil
}
