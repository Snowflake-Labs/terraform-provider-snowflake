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

func FunctionPython() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.FunctionPython, CreateContextFunctionPython),
		ReadContext:   TrackingReadWrapper(resources.FunctionPython, ReadContextFunctionPython),
		UpdateContext: TrackingUpdateWrapper(resources.FunctionPython, UpdateContextFunctionPython),
		DeleteContext: TrackingDeleteWrapper(resources.FunctionPython, DeleteContextFunctionPython),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FunctionPython, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(pythonFunctionSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(pythonFunctionSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(functionParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllFunctionParameters), strings.ToLower)...),
			functionParametersCustomDiff,
			// TODO[SNOW-1348103]: recreate when type changed externally
		)),

		Schema: collections.MergeMaps(pythonFunctionSchema, functionParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextFunctionPython(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadContextFunctionPython(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateContextFunctionPython(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func DeleteContextFunctionPython(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
