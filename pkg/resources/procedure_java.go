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

func ProcedureJava() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ProcedureJava, CreateContextProcedureJava),
		ReadContext:   TrackingReadWrapper(resources.ProcedureJava, ReadContextProcedureJava),
		UpdateContext: TrackingUpdateWrapper(resources.ProcedureJava, UpdateContextProcedureJava),
		DeleteContext: TrackingDeleteWrapper(resources.ProcedureJava, DeleteProcedure),
		Description:   "Resource used to manage java procedure objects. For more information, check [procedure documentation](https://docs.snowflake.com/en/sql-reference/sql/create-procedure).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ProcedureJava, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(javaProcedureSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(javaProcedureSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(procedureParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllProcedureParameters), strings.ToLower)...),
			procedureParametersCustomDiff,
			// TODO[SNOW-1348103]: recreate when type changed externally
		)),

		Schema: collections.MergeMaps(javaProcedureSchema, procedureParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextProcedureJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadContextProcedureJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateContextProcedureJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
