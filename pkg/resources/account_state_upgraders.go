package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func v0_99_0_AccountStateUpgrader(ctx context.Context, state map[string]any, meta any) (map[string]any, error) {
	client := meta.(*provider.Context).Client
	state["must_change_password"] = booleanStringFromBool(state["must_change_password"].(bool))
	state["is_org_admin"] = booleanStringFromBool(state["is_org_admin"].(bool))
	account, err := client.Accounts.ShowByID(ctx, sdk.NewAccountObjectIdentifier(state["name"].(string)))
	if err != nil {
		return nil, err
	}

	state["id"] = helpers.EncodeResourceIdentifier(sdk.NewAccountIdentifier(account.OrganizationName, account.AccountName))

	return state, nil
}
