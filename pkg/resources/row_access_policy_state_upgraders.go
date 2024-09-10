package resources

import (
	"context"
)

func v0_95_0_RowAccessPolicyStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	rawState["body"] = rawState["row_access_expression"]
	delete(rawState, "row_access_expression")

	signature := rawState["signature"].(map[string]any)
	args := make([]map[string]any, 0)
	for k, v := range signature {
		args = append(args, map[string]any{
			"name": k,
			"type": v,
		})
	}
	rawState["argument"] = args
	delete(rawState, "signature")

	return migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName(ctx, rawState, meta)
}
