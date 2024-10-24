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

	orgName, err := client.ContextFunctions.CurrentOrganizationName(ctx)
	require.NoError(t, err)

	accountName, err := client.ContextFunctions.CurrentAccountName(ctx)
	require.NoError(t, err)

	region, err := client.ContextFunctions.CurrentRegion(ctx)
	require.NoError(t, err)

	t.Run("Create minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Connection.DropFunc(t, id))

		assertions.AssertThatObject(t, objectassert.Connection(t, id).
			HasSnowflakeRegion(region).
			HasAccountName(accountName).
			HasName(id.Name()).
			HasNoComment().
			HasIsPrimary(true).
			HasPrimary(fmt.Sprintf("%s.%s.%s", orgName, accountName, id.Name())).
			HasFailoverAllowedToAccounts(
				[]string{
					fmt.Sprintf("%s.%s", orgName, accountName),
				},
			).
			HasOrganizationName(orgName).
			HasAccountLocator(client.GetAccountLocator()),
		)
	})

	t.Run("Create all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id).
			WithIfNotExists(true).
			WithComment("test comment for connection"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Connection.DropFunc(t, id))

		assertions.AssertThatObject(t, objectassert.Connection(t, id).
			HasSnowflakeRegion(region).
			HasAccountName(accountName).
			HasName(id.Name()).
			HasComment("test comment for connection").
			HasIsPrimary(true).
			HasPrimary(fmt.Sprintf("%s.%s.%s", orgName, accountName, id.Name())).
			HasFailoverAllowedToAccounts(
				[]string{
					fmt.Sprintf("%s.%s", orgName, accountName),
				},
			).
			HasOrganizationName(orgName).
			HasAccountLocator(client.GetAccountLocator()),
		)
	})

	// TODO: uncomment when able to change accounts to different regions
	// Snowflake error: The connection cannot be failed over to an account in the same region
	/*
		t.Run("AlterFailover EnableFailover", func(t *testing.T) {
			id := testClientHelper().Ids.RandomAccountObjectIdentifier()
			secondaryAccountId := secondaryTestClientHelper().Ids.AccountIdentifierWithLocator()

			primaryConn, connectionCleanup := testClientHelper().Connection.CreateConnection(t, id)
			t.Cleanup(connectionCleanup)

			err := client.Connections.AlterFailover(ctx, sdk.NewAlterFailoverConnectionRequest(id).
				WithEnableConnectionFailover(
					*sdk.NewEnableConnectionFailoverRequest().WithToAccounts(
						[]sdk.AccountIdentifier{
							secondaryAccountId,
						},
					),
				),
			)
			require.NoError(t, err)

			assertions.AssertThatObject(t, objectassert.Connection(t, id).
				HasSnowflakeRegion(region).
				HasAccountName(accountName).
				HasName(id.Name()).
				HasNoComment().
				HasIsPrimary(true).
				HasPrimary(fmt.Sprintf("%s.%s.%s", orgName, accountName, id.Name())).
				HasFailoverAllowedToAccounts(
					[]string{
						fmt.Sprintf("%s.%s", orgName, accountName),
						fmt.Sprintf("%s.%s", orgName, secondaryAccountId.Name()),
					},
				).
				HasOrganizationName(orgName).
				HasAccountLocator(client.GetAccountLocator()),
			)
		})

		t.Run("AlterFailover EnableFailover With Ignore Edittion Check", func(t *testing.T) {
			id := testClientHelper().Ids.RandomAccountObjectIdentifier()
			secondaryAccountId := secondaryTestClientHelper().Ids.AccountIdentifierWithLocator()

			primaryConn, connectionCleanup := testClientHelper().Connection.CreateConnection(t, id)
			t.Cleanup(connectionCleanup)

			err := client.Connections.AlterFailover(ctx, sdk.NewAlterFailoverConnectionRequest(primaryConn.ID()).
				WithEnableConnectionFailover(
					*sdk.NewEnableConnectionFailoverRequest().WithToAccounts(
						[]sdk.AccountIdentifier{
							secondaryAccountId,
						},
					).WithIgnoreEditionCheck(true),
				),
			)
			require.NoError(t, err)

			assertions.AssertThatObject(t, objectassert.Connection(t, id).
				HasPrimary(fmt.Sprintf("%s.%s.%s", orgName, accountName, id.Name())).
				HasFailoverAllowedToAccounts(
					[]string{
						fmt.Sprintf("%s.%s", orgName, accountName),
						fmt.Sprintf("%s.%s", orgName, secondaryAccountId.Name()),
					},
				),
			)

			// try to alter enable failover to accoutns list
			err = client.Connections.AlterFailover(ctx, sdk.NewAlterFailoverConnectionRequest(id).
				WithEnableConnectionFailover(
					*sdk.NewEnableConnectionFailoverRequest().WithToAccounts(
						[]sdk.AccountIdentifier{},
					),
				),
			)
			require.NoError(t, err)

			// assert that list has not been changed
			assertions.AssertThatObject(t, objectassert.Connection(t, id).
				HasPrimary(fmt.Sprintf("%s.%s.%s", orgName, accountName, id.Name())).
				HasFailoverAllowedToAccounts(
					[]string{
						fmt.Sprintf("%s.%s", orgName, accountName),
						fmt.Sprintf("%s.%s", orgName, secondaryAccountId.Name()),
					},
				),
			)
		})

		t.Run("AlterFailover DisableFailover", func(t *testing.T) {
			id := testClientHelper().Ids.RandomAccountObjectIdentifier()
			accountId := testClientHelper().Ids.AccountIdentifierWithLocator()
			secondaryAccountId := secondaryTestClientHelper().Ids.AccountIdentifierWithLocator()

			primaryConn, connectionCleanup := testClientHelper().Connection.CreateConnection(t, testClientHelper().Ids.RandomAccountObjectIdentifier())
			t.Cleanup(connectionCleanup)

			// add secondary account to failover list
			err := client.Connections.AlterFailover(ctx, sdk.NewAlterFailoverConnectionRequest(primaryConn.ID()).
				WithEnableConnectionFailover(
					*sdk.NewEnableConnectionFailoverRequest().WithToAccounts(
						[]sdk.AccountIdentifier{
							secondaryAccountId,
						},
					),
				),
			)
			require.NoError(t, err)

			assertions.AssertThatObject(t, objectassert.Connection(t, primaryConn.ID()).
				HasPrimary(fmt.Sprintf("%s.%s.%s", orgName, accountName, primaryConn.ID().Name())).
				HasFailoverAllowedToAccounts(
					[]string{
						fmt.Sprintf("%s.%s", orgName, accountName),
						fmt.Sprintf("%s.%s", orgName, secondaryAccountId.Name()),
					},
				),
			)

			// create repllication for secondary account
			err = client.Connections.CreateReplicated(ctx, sdk.NewCreateReplicatedConnectionRequest(id, sdk.NewExternalObjectIdentifier(accountId, primaryConn.ID())))

			// assert that it is not a primary connection
			assertions.AssertThatObject(t, objectassert.Connection(t, id).
				HasPrimary(fmt.Sprintf("%s.%s.%s", orgName, secondaryAccountId.Name(), id.Name())).
				HasIsPrimary(false),
			)

			// Promote to primary
			err = client.Connections.AlterFailover(ctx, sdk.NewAlterFailoverConnectionRequest(id).
				WithPrimary(true))

			assertions.AssertThatObject(t, objectassert.Connection(t, id).
				HasPrimary(fmt.Sprintf("%s.%s.%s", orgName, secondaryAccountId.Name(), id.Name())).
				HasIsPrimary(true),
			)
		})

		t.Run("CreateReplicated", func(t *testing.T) {
			id := testClientHelper().Ids.RandomAccountObjectIdentifier()
			accountId := testClientHelper().Ids.AccountIdentifierWithLocator()
			secondaryAccountId := secondaryTestClientHelper().Ids.AccountIdentifierWithLocator()

			primaryConn, connectionCleanup := testClientHelper().Connection.CreateConnection(t, testClientHelper().Ids.RandomAccountObjectIdentifier())
			t.Cleanup(connectionCleanup)

			err := client.Connections.AlterFailover(ctx, sdk.NewAlterFailoverConnectionRequest(primaryConn.ID()).
				WithEnableConnectionFailover(
					*sdk.NewEnableConnectionFailoverRequest().WithToAccounts(
						[]sdk.AccountIdentifier{
							secondaryAccountId,
						},
					),
				),
			)
			require.NoError(t, err)

			err = client.Connections.CreateReplicated(ctx, sdk.NewCreateReplicatedConnectionRequest(id, sdk.NewExternalObjectIdentifier(accountId, primaryConn.ID())))
			require.NoError(t, err)
			t.Cleanup(testClientHelper().Connection.DropFunc(t, id))

			assertions.AssertThatObject(t, objectassert.Connection(t, id).
				HasSnowflakeRegion(region).
				HasAccountName(accountName).
				HasName(id.Name()).
				HasNoComment().
				HasIsPrimary(false).
				HasPrimary(fmt.Sprintf("%s.%s.%s", orgName, accountName, primaryConn.ID().Name())).
				HasFailoverAllowedToAccounts(
					[]string{
						fmt.Sprintf("%s.%s", orgName, accountName),
						fmt.Sprintf("%s.%s", orgName, secondaryAccountId.Name()),
					},
				).
				HasOrganizationName(orgName).
				HasAccountLocator(client.GetAccountLocator()),
			)
		})
	*/

	t.Run("AlterConnection", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		_, connectionCleanup := testClientHelper().Connection.CreateConnection(t, id)
		t.Cleanup(connectionCleanup)

		// Set
		client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithSet(*sdk.NewSetRequest().
				WithComment("new integration test comment")))

		assertions.AssertThatObject(t, objectassert.Connection(t, id).
			HasName(id.Name()).
			HasComment("new integration test comment"),
		)

		// Unset
		client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithUnset(*sdk.NewUnsetRequest().
				WithComment(true)))

		assertions.AssertThatObject(t, objectassert.Connection(t, id).
			HasName(id.Name()).
			HasNoComment(),
		)
	})

	t.Run("Drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		_, connectionCleanup := testClientHelper().Connection.CreateConnection(t, id)
		t.Cleanup(connectionCleanup)

		connection, err := client.Connections.ShowByID(ctx, id)
		require.NoError(t, err)
		require.NotNil(t, connection)

		err = client.Connections.Drop(ctx, sdk.NewDropConnectionRequest(id))
		require.NoError(t, err)

		connection, err = client.Connections.ShowByID(ctx, id)
		require.Nil(t, connection)
		require.Error(t, err)
	})

	t.Run("Show", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		id2 := testClientHelper().Ids.RandomAccountObjectIdentifier()

		connection1, connectionCleanup1 := testClientHelper().Connection.CreateConnection(t, id1)
		t.Cleanup(connectionCleanup1)

		connection2, connectionCleanup2 := testClientHelper().Connection.CreateConnection(t, id2)
		t.Cleanup(connectionCleanup2)

		returnedConnections, err := client.Connections.Show(ctx, sdk.NewShowConnectionRequest())
		require.NoError(t, err)
		require.Contains(t, returnedConnections, *connection1)
		require.Contains(t, returnedConnections, *connection2)
	})

	t.Run("Show with Like", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		id2 := testClientHelper().Ids.RandomAccountObjectIdentifier()

		connection1, connectionCleanup1 := testClientHelper().Connection.CreateConnection(t, id1)
		t.Cleanup(connectionCleanup1)

		connection2, connectionCleanup2 := testClientHelper().Connection.CreateConnection(t, id2)
		t.Cleanup(connectionCleanup2)

		returnedConnections, err := client.Connections.Show(ctx, sdk.NewShowConnectionRequest().
			WithLike(sdk.Like{
				Pattern: sdk.String(id1.Name()),
			}))
		require.NoError(t, err)
		require.Contains(t, returnedConnections, *connection1)
		require.NotContains(t, returnedConnections, *connection2)
	})

	t.Run("ShowByID", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		_, connectionCleanup := testClientHelper().Connection.CreateConnection(t, id)
		t.Cleanup(connectionCleanup)

		returnedConnection, err := client.Connections.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, returnedConnection.ID())
		require.Equal(t, id.Name(), returnedConnection.Name)
	})
}
