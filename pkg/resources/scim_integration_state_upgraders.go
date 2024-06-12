package resources

import "context"

func v091ScimIntegrationStateUpgrader(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		return rawState, nil
	}

	rawState["run_as_role"] = rawState["provisioner_role"]
	delete(rawState, "provisioner_role")
	return rawState, nil
}
