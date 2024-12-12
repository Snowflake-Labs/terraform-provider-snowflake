package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFunctionOrProcedureArguments(d *schema.ResourceData, args []sdk.NormalizedArgument) error {
	if len(args) == 0 {
		// TODO [SNOW-1348103]: handle empty list
		return nil
	}
	// We do it the unusual way because the default values are not returned by SF.
	// We update what we have - leaving the defaults unchanged.
	if currentArgs, ok := d.Get("arguments").([]map[string]any); !ok {
		return fmt.Errorf("arguments must be a list")
	} else {
		for i, arg := range args {
			currentArgs[i]["arg_name"] = arg.Name
			currentArgs[i]["arg_data_type"] = arg.DataType.ToSql()
		}
		return d.Set("arguments", currentArgs)
	}
}

func readFunctionOrProcedureImports(d *schema.ResourceData, imports []sdk.NormalizedPath) error {
	if len(imports) == 0 {
		// don't do anything if imports not present
		return nil
	}
	imps := collections.Map(imports, func(imp sdk.NormalizedPath) map[string]any {
		return map[string]any{
			"stage_location": imp.StageLocation,
			"path_on_stage":  imp.PathOnStage,
		}
	})
	return d.Set("imports", imps)
}

func readFunctionOrProcedureTargetPath(d *schema.ResourceData, normalizedPath *sdk.NormalizedPath) error {
	if normalizedPath == nil {
		// don't do anything if imports not present
		return nil
	}
	tp := make([]map[string]any, 1)
	tp[0] = map[string]any{
		"stage_location": normalizedPath.StageLocation,
		"path_on_stage":  normalizedPath.PathOnStage,
	}
	return d.Set("target_path", tp)
}
