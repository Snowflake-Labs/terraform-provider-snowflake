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

func FunctionScala() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.FunctionScala, CreateContextFunctionScala),
		ReadContext:   TrackingReadWrapper(resources.FunctionScala, ReadContextFunctionScala),
		UpdateContext: TrackingUpdateWrapper(resources.FunctionScala, UpdateContextFunctionScala),
		DeleteContext: TrackingDeleteWrapper(resources.FunctionScala, DeleteContextFunctionScala),
		Description:   "Resource used to manage scala function objects. For more information, check [function documentation](https://docs.snowflake.com/en/sql-reference/sql/create-function).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FunctionScala, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(scalaFunctionSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(scalaFunctionSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(functionParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllFunctionParameters), strings.ToLower)...),
			functionParametersCustomDiff,
			// TODO[SNOW-1348103]: recreate when type changed externally
		)),

		Schema: collections.MergeMaps(scalaFunctionSchema, functionParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextFunctionScala(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadContextFunctionScala(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateContextFunctionScala(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func DeleteContextFunctionScala(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
