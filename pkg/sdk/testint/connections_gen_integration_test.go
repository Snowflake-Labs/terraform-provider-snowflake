package testint

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"

	assertions "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
)

const ConnectionFailoverToAccountInSameRegionErrorMessage = "The connection cannot be failed over to an account in the same region"

func TestInt_Connections(t *testing.T) {
	client := testClient(t)
	secondaryClient := testSecondaryClient(t)
	ctx := testContext(t)

	sessionDetails, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	require.NoError(t, err)

	t.Run("Create minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		require.NoError(t, err)

		err = client.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Connection.DropFunc(t, id))

		assertions.AssertThatObject(t, objectassert.Connection(t, id).
			HasSnowflakeRegion(sessionDetails.Region).
			HasAccountName(sessionDetails.AccountName).
			HasName(id.Name()).
			HasNoComment().
			HasIsPrimary(true).
			HasPrimary(fmt.Sprintf("%s.%s.%s", sessionDetails.OrganizationName, sessionDetails.AccountName, id.Name())).
			HasFailoverAllowedToAccounts(
				[]string{
					fmt.Sprintf("%s.%s", sessionDetails.OrganizationName, sessionDetails.AccountName),
				},
			).
			HasOrganizationName(sessionDetails.OrganizationName).
			HasAccountLocator(client.GetAccountLocator()).
			HasConnectionUrl(
				strings.ToLower(
					fmt.Sprintf("%s-%s.snowflakecomputing.com", sessionDetails.OrganizationName, id.Name()),
				),
			),
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
			HasSnowflakeRegion(sessionDetails.Region).
			HasAccountName(sessionDetails.AccountName).
			HasName(id.Name()).
			HasComment("test comment for connection").
			HasIsPrimary(true).
			HasPrimary(fmt.Sprintf("%s.%s.%s", sessionDetails.OrganizationName, sessionDetails.AccountName, id.Name())).
			HasFailoverAllowedToAccounts(
				[]string{
					fmt.Sprintf("%s.%s", sessionDetails.OrganizationName, sessionDetails.AccountName),
				},
			).
			HasOrganizationName(sessionDetails.OrganizationName).
			HasAccountLocator(client.GetAccountLocator()).
			HasConnectionUrl(
				strings.ToLower(
					fmt.Sprintf("%s-%s.snowflakecomputing.com", sessionDetails.OrganizationName, id.Name()),
				),
			),
		)
	})

	t.Run("Alter enable failover", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		secondaryAccountId := secondaryTestClientHelper().Ids.AccountIdentifierWithLocator()

		_, connectionCleanup := testClientHelper().Connection.Create(t, id)
		t.Cleanup(connectionCleanup)

		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithEnableConnectionFailover(
				*sdk.NewEnableConnectionFailoverRequest().WithToAccounts(
					[]sdk.AccountIdentifier{
						secondaryAccountId,
					},
				),
			),
		)
		require.ErrorContains(t, err, ConnectionFailoverToAccountInSameRegionErrorMessage)

		// TODO: [SNOW-1763442]
		/*
		   require.NoError(t, err)
		   assertions.AssertThatObject(t, objectassert.Connection(t, id).
		       HasSnowflakeRegion(sessionDetails.Region).
		       HasAccountName(sessionDetails.AccountName).
		       HasName(id.Name()).
		       HasNoComment().
		       HasIsPrimary(true).
		       HasPrimary(fmt.Sprintf("%s.%s.%s", sessionDetails.OrganizationName, sessionDetails.AccountName, id.Name())).
		       HasFailoverAllowedToAccounts(
		           []string{
		               fmt.Sprintf("%s.%s", sessionDetails.OrganizationName, sessionDetails.AccountName),
		               fmt.Sprintf("%s.%s", sessionDetails.OrganizationName, secondaryAccountId.Name()),
		           },
		       ).
		       HasOrganizationName(sessionDetails.OrganizationName).
		       HasAccountLocator(client.GetAccountLocator()),
		       HasConnectionUrl(
		           strings.ToLower(
		               fmt.Sprintf("%s-%s.snowflakecomputing.com", sessionDetails.OrganizationName, id.Name()),
		           ),
		       ),
		   )
		*/
	})

	t.Run("Create as replica of", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		accountId := testClientHelper().Ids.AccountIdentifierWithLocator()
		_ = id
		_ = accountId
		secondaryAccountId := secondaryTestClientHelper().Ids.AccountIdentifierWithLocator()

		primaryConn, connectionCleanup := testClientHelper().Connection.Create(t, testClientHelper().Ids.RandomAccountObjectIdentifier())
		t.Cleanup(connectionCleanup)

		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(primaryConn.ID()).
			WithEnableConnectionFailover(
				*sdk.NewEnableConnectionFailoverRequest().WithToAccounts(
					[]sdk.AccountIdentifier{
						secondaryAccountId,
					},
				),
			),
		)
		require.ErrorContains(t, err, ConnectionFailoverToAccountInSameRegionErrorMessage)
		// TODO: [SNOW-1763442]
		//
		// require.NoError(t, err)

		/*
		   // create replica on secondary account
		   err = secondaryClient.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id).
		       WithAsReplicaOf(sdk.AsReplicaOfRequest{
		           AsReplicaOf: sdk.NewExternalObjectIdentifier(accountId, id.Name()),
		       }))
		   t.Cleanup(testClientHelper().Connection.DropFunc(t, id))
		   require.NoError(t, err)

		   assertions.AssertThatObject(t, objectassert.Connection(t, id).
		       HasSnowflakeRegion(sessionDetails.Region).
		       HasAccountName(sessionDetails.AccountName).
		       HasName(id.Name()).
		       HasNoComment().
		       HasIsPrimary(false).
		       HasPrimary(fmt.Sprintf("%s.%s.%s", sessionDetails.OrganizationName, sessionDetails.AccountName, id.Name())).
		       HasFailoverAllowedToAccounts(
		           []string{
		               fmt.Sprintf("%s.%s", sessionDetails.OrganizationName, sessionDetails.AccountName),
		               fmt.Sprintf("%s.%s", sessionDetails.OrganizationName, secondaryAccountId.Name()),
		           },
		       ).
		       HasOrganizationName(sessionDetails.OrganizationName).
		       HasAccountLocator(client.GetAccountLocator()).
		       HasConnectionUrl(
		           strings.ToLower(
		               fmt.Sprintf("%s-%s.snowflakecomputing.com", sessionDetails.OrganizationName, id.Name()),
		           ),
		       ),
		   )
		*/
	})

	t.Run("Alter disable failover", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		accountId := testClientHelper().Ids.AccountIdentifierWithLocator()
		secondaryAccountId := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

		primaryConn, connectionCleanup := testClientHelper().Connection.Create(t, id)
		t.Cleanup(connectionCleanup)

		// Add secondary account to failover list
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithEnableConnectionFailover(
				*sdk.NewEnableConnectionFailoverRequest().WithToAccounts(
					[]sdk.AccountIdentifier{
						secondaryAccountId,
					},
				),
			),
		)
		require.ErrorContains(t, err, ConnectionFailoverToAccountInSameRegionErrorMessage)
		// TODO: [SNOW-1763442]
		//
		// require.NoError(t, err)

		// Disable promotion of this connection
		err = client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithDisableConnectionFailover(*sdk.NewDisableConnectionFailoverRequest()))
		require.NoError(t, err)

		// Assert that promotion for other account has been disabled
		assertions.AssertThatObject(t, objectassert.Connection(t, primaryConn.ID()).
			HasPrimary(fmt.Sprintf("%s.%s.%s", sessionDetails.OrganizationName, sessionDetails.AccountName, id.Name())).
			HasFailoverAllowedToAccounts(
				[]string{
					fmt.Sprintf("%s.%s", sessionDetails.OrganizationName, sessionDetails.AccountName),
				},
			),
		)

		// Try to create repllication on secondary account
		err = secondaryClient.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id).
			WithAsReplicaOf(sdk.AsReplicaOfRequest{
				AsReplicaOf: sdk.NewExternalObjectIdentifier(accountId, id),
			}))
		require.ErrorContains(t, err, "This account is not authorized to create a secondary connection of this primary connection")
	})

	t.Run("Alter", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		_, connectionCleanup := testClientHelper().Connection.Create(t, id)
		t.Cleanup(connectionCleanup)

		// Set
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithSet(*sdk.NewSetRequest().
				WithComment("new integration test comment")))
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.Connection(t, id).
			HasName(id.Name()).
			HasComment("new integration test comment"),
		)

		// Unset
		err = client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithUnset(*sdk.NewUnsetRequest().
				WithComment(true)))
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.Connection(t, id).
			HasName(id.Name()).
			HasNoComment(),
		)
	})

	t.Run("Drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		_, connectionCleanup := testClientHelper().Connection.Create(t, id)
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

		connection1, connectionCleanup1 := testClientHelper().Connection.Create(t, id1)
		t.Cleanup(connectionCleanup1)

		connection2, connectionCleanup2 := testClientHelper().Connection.Create(t, id2)
		t.Cleanup(connectionCleanup2)

		returnedConnections, err := client.Connections.Show(ctx, sdk.NewShowConnectionRequest())
		require.NoError(t, err)
		require.Contains(t, returnedConnections, *connection1)
		require.Contains(t, returnedConnections, *connection2)
	})

	t.Run("Show with Like", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		id2 := testClientHelper().Ids.RandomAccountObjectIdentifier()

		connection1, connectionCleanup1 := testClientHelper().Connection.Create(t, id1)
		t.Cleanup(connectionCleanup1)

		connection2, connectionCleanup2 := testClientHelper().Connection.Create(t, id2)
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

		_, connectionCleanup := testClientHelper().Connection.Create(t, id)
		t.Cleanup(connectionCleanup)

		returnedConnection, err := client.Connections.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, returnedConnection.ID())
		require.Equal(t, id.Name(), returnedConnection.Name)
	})
}
