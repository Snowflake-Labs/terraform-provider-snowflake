package resources

import (
	"context"
	"fmt"
	"strings"
)

func v0_95_0_MaskingPolicyStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	rawState["body"] = rawState["masking_expression"]

	signature := rawState["signature"].([]any)
	if len(signature) != 1 {
		return nil, fmt.Errorf("corrupted signature in state: expected list of length 1, got %d", len(signature))
	}
	columns := signature[0].(map[string]any)["column"].([]any)
	args := make([]map[string]any, 0)
	for _, v := range columns {
		column := v.(map[string]any)
		args = append(args, map[string]any{
			"name": strings.ToUpper(column["name"].(string)),
			"type": column["type"].(string),
		})
	}
	rawState["argument"] = args

	return migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName(ctx, rawState, meta)
}
