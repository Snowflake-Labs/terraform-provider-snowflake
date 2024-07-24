package resources

import (
	"context"
)

func v093SchemaStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	rawState["with_managed_access"] = rawState["is_managed"]
	delete(rawState, "is_managed")

	return rawState, nil
}
