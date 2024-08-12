package resources

import (
	"context"
)

func v0_94_1_TableStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	rawState[FullyQualifiedNameAttributeName] = rawState["qualified_name"]
	delete(rawState, "qualified_name")

	return rawState, nil
}
