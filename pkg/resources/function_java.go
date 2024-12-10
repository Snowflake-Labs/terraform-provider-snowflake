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

func FunctionJava() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.FunctionJava, CreateContextFunctionJava),
		ReadContext:   TrackingReadWrapper(resources.FunctionJava, ReadContextFunctionJava),
		UpdateContext: TrackingUpdateWrapper(resources.FunctionJava, UpdateContextFunctionJava),
		DeleteContext: TrackingDeleteWrapper(resources.FunctionJava, DeleteFunction),
		Description:   "Resource used to manage java function objects. For more information, check [function documentation](https://docs.snowflake.com/en/sql-reference/sql/create-function).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FunctionJava, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(javaFunctionSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(javaFunctionSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(functionParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllFunctionParameters), strings.ToLower)...),
			functionParametersCustomDiff,
			// TODO[SNOW-1348103]: recreate when type changed externally
		)),

		Schema: collections.MergeMaps(javaFunctionSchema, functionParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextFunctionJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadContextFunctionJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateContextFunctionJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
