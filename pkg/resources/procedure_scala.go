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

func ProcedureScala() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ProcedureScala, CreateContextProcedureScala),
		ReadContext:   TrackingReadWrapper(resources.ProcedureScala, ReadContextProcedureScala),
		UpdateContext: TrackingUpdateWrapper(resources.ProcedureScala, UpdateContextProcedureScala),
		DeleteContext: TrackingDeleteWrapper(resources.ProcedureScala, DeleteProcedure),
		Description:   "Resource used to manage scala procedure objects. For more information, check [procedure documentation](https://docs.snowflake.com/en/sql-reference/sql/create-procedure).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ProcedureScala, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(scalaProcedureSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(scalaProcedureSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(procedureParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllProcedureParameters), strings.ToLower)...),
			procedureParametersCustomDiff,
			// TODO[SNOW-1348103]: recreate when type changed externally
		)),

		Schema: collections.MergeMaps(scalaProcedureSchema, procedureParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextProcedureScala(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadContextProcedureScala(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateContextProcedureScala(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
