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

func ProcedureSql() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ProcedureSql, CreateContextProcedureSql),
		ReadContext:   TrackingReadWrapper(resources.ProcedureSql, ReadContextProcedureSql),
		UpdateContext: TrackingUpdateWrapper(resources.ProcedureSql, UpdateContextProcedureSql),
		DeleteContext: TrackingDeleteWrapper(resources.ProcedureSql, DeleteContextProcedureSql),
		Description:   "Resource used to manage sql procedure objects. For more information, check [procedure documentation](https://docs.snowflake.com/en/sql-reference/sql/create-procedure).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ProcedureSql, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(sqlProcedureSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(sqlProcedureSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(procedureParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllProcedureParameters), strings.ToLower)...),
			procedureParametersCustomDiff,
			// TODO[SNOW-1348103]: recreate when type changed externally
		)),

		Schema: collections.MergeMaps(sqlProcedureSchema, procedureParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextProcedureSql(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadContextProcedureSql(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateContextProcedureSql(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func DeleteContextProcedureSql(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
