package resources

import (
	"context"
	"strconv"
)

func v0_94_1_ViewStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	if v, ok := rawState["is_secure"]; ok {
		rawState["is_secure"] = strconv.FormatBool(v.(bool))
	}

	delete(rawState, "tag")

	return rawState, nil
}
