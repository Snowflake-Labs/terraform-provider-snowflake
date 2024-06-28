package resources

import (
	"context"
	"strconv"
)

func v092ScimIntegrationStateUpgrader(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		return rawState, nil
	}

	rawState["run_as_role"] = rawState["provisioner_role"]
	delete(rawState, "provisioner_role")

	if v, ok := rawState["enabled"]; ok {
		rawState["enabled"] = strconv.FormatBool(v.(bool))
	}

	return rawState, nil
}
