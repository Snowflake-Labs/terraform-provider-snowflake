package resources

import (
	"context"
)

func v092DatabaseStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	if replicationConfigurations, ok := rawState["replication_configuration"]; ok && len(replicationConfigurations.([]any)) == 1 {
		replicationConfiguration := replicationConfigurations.([]any)[0].(map[string]any)
		replication := make(map[string]any)
		replication["ignore_edition_check"] = replicationConfiguration["ignore_edition_check"]

		accounts := replicationConfiguration["accounts"].([]any)
		enableForAccounts := make([]map[string]any, len(accounts))
		for i, account := range accounts {
			enableForAccounts[i] = map[string]any{
				"account_identifier": account,
			}
		}

		rawState["replication"] = []any{replication}
	}

	return rawState, nil
}
