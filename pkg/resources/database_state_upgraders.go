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

// resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.#", "1"),
// resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.0.ignore_edition_check", "true"),
// resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.0.accounts.#", "1"),
// resource.TestCheckResourceAttr("snowflake_database.test", "replication_configuration.0.accounts.0", secondaryAccountIdentifier),

// resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.ignore_edition_check", "true"),
// resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.#", "1"),
// resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.account_identifier", secondaryAccountIdentifier),
// resource.TestCheckResourceAttr("snowflake_database.test", "replication.0.enable_to_account.0.with_failover", "false"),
