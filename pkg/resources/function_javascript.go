package resources

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FunctionJavascript() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.FunctionJavascript, CreateContextFunctionJavascript),
		ReadContext:   TrackingReadWrapper(resources.FunctionJavascript, ReadContextFunctionJavascript),
		UpdateContext: TrackingUpdateWrapper(resources.FunctionJavascript, UpdateContextFunctionJavascript),
		DeleteContext: TrackingDeleteWrapper(resources.FunctionJavascript, DeleteContextFunctionJavascript),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FunctionJavascript, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(javascriptFunctionSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(javascriptFunctionSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(functionParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllFunctionParameters), strings.ToLower)...),
			functionParametersCustomDiff,
			// TODO[SNOW-1348103]: recreate when type changed externally
		)),

		Schema: collections.MergeMaps(javascriptFunctionSchema, functionParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextFunctionJavascript(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadContextFunctionJavascript(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateContextFunctionJavascript(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func DeleteContextFunctionJavascript(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
