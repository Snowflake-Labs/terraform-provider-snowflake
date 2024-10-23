package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"

	assertions "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
)

func TestInt_Connections(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("Create minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Connection.DropFunc(t, id))

		orgName, err := client.ContextFunctions.CurrentOrganizationName(ctx)
		require.NoError(t, err)

		accountName, err := client.ContextFunctions.CurrentAccountName(ctx)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.Connection(t, id).
			HasName(id.Name()).
			HasNoComment().
			HasAccountLocator(client.GetAccountLocator()).
			HasFailoverAllowedToAccounts(
				[]string{
					fmt.Sprintf("%s.%s", orgName, accountName),
				},
			).
			HasIsPrimary(true),
		)
	})

	t.Run("CreateReplicated", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("AlterFailover", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Alter", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Drop", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Show", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("ShowByID", func(t *testing.T) {
		// TODO: fill me
	})
}
