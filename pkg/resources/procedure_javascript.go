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

func ProcedureJavascript() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ProcedureJavascript, CreateContextProcedureJavascript),
		ReadContext:   TrackingReadWrapper(resources.ProcedureJavascript, ReadContextProcedureJavascript),
		UpdateContext: TrackingUpdateWrapper(resources.ProcedureJavascript, UpdateContextProcedureJavascript),
		DeleteContext: TrackingDeleteWrapper(resources.ProcedureJavascript, DeleteProcedure),
		Description:   "Resource used to manage javascript procedure objects. For more information, check [procedure documentation](https://docs.snowflake.com/en/sql-reference/sql/create-procedure).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ProcedureJavascript, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(javascriptProcedureSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(javascriptProcedureSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(procedureParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllProcedureParameters), strings.ToLower)...),
			procedureParametersCustomDiff,
			// TODO[SNOW-1348103]: recreate when type changed externally
		)),

		Schema: collections.MergeMaps(javascriptProcedureSchema, procedureParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextProcedureJavascript(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadContextProcedureJavascript(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateContextProcedureJavascript(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
