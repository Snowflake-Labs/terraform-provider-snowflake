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

func FunctionSql() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.FunctionSql, CreateContextFunctionSql),
		ReadContext:   TrackingReadWrapper(resources.FunctionSql, ReadContextFunctionSql),
		UpdateContext: TrackingUpdateWrapper(resources.FunctionSql, UpdateContextFunctionSql),
		DeleteContext: TrackingDeleteWrapper(resources.FunctionSql, DeleteContextFunctionSql),
		Description:   "Resource used to manage sql function objects. For more information, check [function documentation](https://docs.snowflake.com/en/sql-reference/sql/create-function).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FunctionSql, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(sqlFunctionSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(sqlFunctionSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(functionParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllFunctionParameters), strings.ToLower)...),
			functionParametersCustomDiff,
			// TODO[SNOW-1348103]: recreate when type changed externally
		)),

		Schema: collections.MergeMaps(sqlFunctionSchema, functionParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextFunctionSql(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadContextFunctionSql(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateContextFunctionSql(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func DeleteContextFunctionSql(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
