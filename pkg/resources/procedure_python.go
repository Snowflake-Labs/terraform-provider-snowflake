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

func ProcedurePython() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ProcedurePython, CreateContextProcedurePython),
		ReadContext:   TrackingReadWrapper(resources.ProcedurePython, ReadContextProcedurePython),
		UpdateContext: TrackingUpdateWrapper(resources.ProcedurePython, UpdateContextProcedurePython),
		DeleteContext: TrackingDeleteWrapper(resources.ProcedurePython, DeleteProcedure),
		Description:   "Resource used to manage python procedure objects. For more information, check [procedure documentation](https://docs.snowflake.com/en/sql-reference/sql/create-procedure).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ProcedurePython, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(pythonProcedureSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(pythonProcedureSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(procedureParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllProcedureParameters), strings.ToLower)...),
			procedureParametersCustomDiff,
			// TODO[SNOW-1348103]: recreate when type changed externally
		)),

		Schema: collections.MergeMaps(pythonProcedureSchema, procedureParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextProcedurePython(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadContextProcedurePython(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateContextProcedurePython(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
