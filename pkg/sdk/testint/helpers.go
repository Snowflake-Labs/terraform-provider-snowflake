package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func createDatabaseFromShare(t *testing.T) (*sdk.Database, func()) {
	t.Helper()
	client := testClient(t)
	secondaryClient := testSecondaryClient(t)
	ctx := testContext(t)

	shareTest, shareCleanup := secondaryTestClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	sharedDatabase, sharedDatabaseCleanup := secondaryTestClientHelper().Database.CreateDatabase(t)
	t.Cleanup(sharedDatabaseCleanup)

	databaseId := sharedDatabase.ID()

	err := secondaryClient.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
		Database: sharedDatabase.ID(),
	}, shareTest.ID())
	require.NoError(t, err)
	t.Cleanup(func() {
		err := secondaryClient.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: sharedDatabase.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})

	err = secondaryClient.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
		IfExists: sdk.Bool(true),
		Set: &sdk.ShareSet{
			Accounts: []sdk.AccountIdentifier{
				testClientHelper().Account.GetAccountIdentifier(t),
			},
		},
	})
	require.NoError(t, err)

	err = client.Databases.CreateShared(ctx, databaseId, shareTest.ExternalID(), &sdk.CreateSharedDatabaseOptions{})
	require.NoError(t, err)

	database, err := client.Databases.ShowByID(ctx, databaseId)
	require.NoError(t, err)

	return database, testClientHelper().Database.DropDatabaseFunc(t, database.ID())
}

func createDatabaseReplica(t *testing.T) (*sdk.Database, func()) {
	t.Helper()
	client := testClient(t)
	secondaryClient := testSecondaryClient(t)
	ctx := testContext(t)

	sharedDatabase, sharedDatabaseCleanup := secondaryTestClientHelper().Database.CreateDatabase(t)
	t.Cleanup(sharedDatabaseCleanup)

	err := secondaryClient.Databases.AlterReplication(ctx, sharedDatabase.ID(), &sdk.AlterDatabaseReplicationOptions{
		EnableReplication: &sdk.EnableReplication{
			ToAccounts: []sdk.AccountIdentifier{
				testClientHelper().Account.GetAccountIdentifier(t),
			},
			IgnoreEditionCheck: sdk.Bool(true),
		},
	})
	require.NoError(t, err)

	externalDatabaseId := sdk.NewExternalObjectIdentifier(secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(), sharedDatabase.ID())
	err = client.Databases.CreateSecondary(ctx, sharedDatabase.ID(), externalDatabaseId, &sdk.CreateSecondaryDatabaseOptions{})
	require.NoError(t, err)

	database, err := client.Databases.ShowByID(ctx, sharedDatabase.ID())
	require.NoError(t, err)

	return database, testClientHelper().Database.DropDatabaseFunc(t, sharedDatabase.ID())
}

func createApplicationPackage(t *testing.T) (*sdk.ApplicationPackage, func()) {
	t.Helper()

	stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(cleanupStage)

	testClientHelper().Stage.PutOnStageWithContent(t, stage.ID(), "manifest.yml", "")
	testClientHelper().Stage.PutOnStageWithContent(t, stage.ID(), "setup.sql", "CREATE APPLICATION ROLE IF NOT EXISTS APP_HELLO_SNOWFLAKE;")

	applicationPackage, cleanupApplicationPackage := testClientHelper().ApplicationPackage.CreateApplicationPackage(t)
	t.Cleanup(cleanupApplicationPackage)

	testClientHelper().ApplicationPackage.AddApplicationPackageVersion(t, applicationPackage.ID(), stage.ID(), "V01")

	return applicationPackage, cleanupApplicationPackage
}

func createShare(t *testing.T, ctx context.Context, client *sdk.Client) (*sdk.Share, func()) {
	t.Helper()
	object, objectCleanup := testClientHelper().Share.CreateShare(t)
	err := client.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
		Database: testClientHelper().Ids.DatabaseId(),
	}, object.ID())
	require.NoError(t, err)
	cleanup := func() {
		err = client.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: testClientHelper().Ids.DatabaseId(),
		}, object.ID())
		require.NoError(t, err)
		objectCleanup()
	}
	return object, cleanup
}

func createPipe(t *testing.T) (*sdk.Pipe, func()) {
	t.Helper()
	table, tableCleanup := testClientHelper().Table.Create(t)

	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)

	object, objectCleanup := testClientHelper().Pipe.CreatePipe(t, fmt.Sprintf("COPY INTO %s\nFROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName()))
	cleanup := func() {
		objectCleanup()
		stageCleanup()
		tableCleanup()
	}
	return object, cleanup
}

func createMaterializedView(t *testing.T) (*sdk.MaterializedView, func()) {
	t.Helper()
	table, tableCleanup := testClientHelper().Table.Create(t)
	query := fmt.Sprintf(`SELECT * FROM %s`, table.ID().FullyQualifiedName())
	object, objectCleanup := testClientHelper().MaterializedView.CreateMaterializedView(t, query, false)
	cleanup := func() {
		objectCleanup()
		tableCleanup()
	}
	return object, cleanup
}

func createStream(t *testing.T) (*sdk.Stream, func()) {
	t.Helper()
	table, tableCleanup := testClientHelper().Table.CreateInSchema(t, testClientHelper().Ids.SchemaId())
	t.Cleanup(tableCleanup)

	object, objectCleanup := testClientHelper().Stream.CreateOnTable(t, table.ID())
	t.Cleanup(objectCleanup)
	cleanup := func() {
		objectCleanup()
		tableCleanup()
	}
	return object, cleanup
}

func createExternalTable(t *testing.T) (*sdk.ExternalTable, func()) {
	t.Helper()
	stageID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
	stageLocation := fmt.Sprintf("@%s", stageID.FullyQualifiedName())
	_, stageCleanup := testClientHelper().Stage.CreateStageWithURL(t, stageID)
	object, objectCleanup := testClientHelper().ExternalTable.CreateWithLocation(t, stageLocation)
	cleanup := func() {
		objectCleanup()
		stageCleanup()
	}

	return object, cleanup
}
