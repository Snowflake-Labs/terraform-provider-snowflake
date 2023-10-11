package testint

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DatabasesCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("minimal", func(t *testing.T) {
		databaseID := randomAccountObjectIdentifier(t)
		err := client.Databases.Create(ctx, databaseID, nil)
		require.NoError(t, err)
		database, err := client.Databases.ShowByID(ctx, databaseID)
		require.NoError(t, err)
		assert.Equal(t, databaseID.Name(), database.Name)
		t.Cleanup(func() {
			err = client.Databases.Drop(ctx, databaseID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("as clone", func(t *testing.T) {
		cloneDatabase, cloneDatabaseCleanup := createDatabase(t, client)
		t.Cleanup(cloneDatabaseCleanup)
		databaseID := randomAccountObjectIdentifier(t)
		opts := &sdk.CreateDatabaseOptions{
			Clone: &sdk.Clone{
				SourceObject: cloneDatabase.ID(),
				At: &sdk.TimeTravel{
					Offset: sdk.Int(0),
				},
			},
		}
		err := client.Databases.Create(ctx, databaseID, opts)
		require.NoError(t, err)
		database, err := client.Databases.ShowByID(ctx, databaseID)
		require.NoError(t, err)
		assert.Equal(t, databaseID.Name(), database.Name)
		t.Cleanup(func() {
			err = client.Databases.Drop(ctx, databaseID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("complete", func(t *testing.T) {
		databaseID := randomAccountObjectIdentifier(t)

		databaseTest, databaseCleanup := createDatabase(t, client)
		t.Cleanup(databaseCleanup)
		schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
		t.Cleanup(schemaCleanup)
		tagTest, tagCleanup := createTag(t, client, databaseTest, schemaTest)
		t.Cleanup(tagCleanup)
		tag2Test, tag2Cleanup := createTag(t, client, databaseTest, schemaTest)
		t.Cleanup(tag2Cleanup)

		comment := sdk.RandomComment(t)
		opts := &sdk.CreateDatabaseOptions{
			OrReplace:                  sdk.Bool(true),
			Transient:                  sdk.Bool(true),
			Comment:                    sdk.String(comment),
			DataRetentionTimeInDays:    sdk.Int(1),
			MaxDataExtensionTimeInDays: sdk.Int(1),
			Tag: []sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
				{
					Name:  tag2Test.ID(),
					Value: "v2",
				},
			},
		}
		err := client.Databases.Create(ctx, databaseID, opts)
		require.NoError(t, err)
		database, err := client.Databases.ShowByID(ctx, databaseID)
		require.NoError(t, err)
		assert.Equal(t, databaseID.Name(), database.Name)
		assert.Equal(t, comment, database.Comment)
		assert.Equal(t, 1, database.RetentionTime)
		// MAX_DATA_EXTENSION_IN_DAYS is an object parameter, not in Database object
		param, err := client.Parameters.ShowObjectParameter(ctx, "MAX_DATA_EXTENSION_TIME_IN_DAYS", sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: databaseID})
		assert.NoError(t, err)
		assert.Equal(t, "1", param.Value)

		// verify tags
		tag1Value, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), database.ID(), sdk.ObjectTypeDatabase)
		require.NoError(t, err)
		assert.Equal(t, "v1", tag1Value)
		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2Test.ID(), database.ID(), sdk.ObjectTypeDatabase)
		require.NoError(t, err)
		assert.Equal(t, "v2", tag2Value)

		t.Cleanup(func() {
			err = client.Databases.Drop(ctx, databaseID, nil)
			require.NoError(t, err)
		})
	})
}

func TestInt_CreateShared(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	shareTest, _ := createShare(t, client)
	// t.Cleanup(shareCleanup)
	err := client.Grants.GrantPrivilegeToShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.GrantPrivilegeToShareOn{
		Database: databaseTest.ID(),
	}, shareTest.ID())
	require.NoError(t, err)
	t.Cleanup(func() {
		err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.ObjectPrivilegeUsage, &sdk.RevokePrivilegeFromShareOn{
			Database: databaseTest.ID(),
		}, shareTest.ID())
	})
	require.NoError(t, err)
	secondaryClient := testSecondaryClient(t)
	accountsToSet := []sdk.AccountIdentifier{
		getAccountIdentifier(t, secondaryClient),
	}
	// first add the account.
	err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
		IfExists: sdk.Bool(true),
		Set: &sdk.ShareSet{
			Accounts: accountsToSet,
		},
	})

	databaseID := randomAccountObjectIdentifier(t)
	err = secondaryClient.Databases.CreateShared(ctx, databaseID, shareTest.ExternalID(), nil)
	require.NoError(t, err)
	database, err := secondaryClient.Databases.ShowByID(ctx, databaseID)
	require.NoError(t, err)
	assert.Equal(t, databaseID.Name(), database.Name)
	t.Cleanup(func() {
		err = secondaryClient.Databases.Drop(ctx, databaseID, nil)
		require.NoError(t, err)
	})
}

func TestInt_DatabasesCreateSecondary(t *testing.T) {
	// todo: once ReplicationGroups are supported.
}

func TestInt_DatabasesDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, _ := createDatabase(t, client)
	databaseID := databaseTest.ID()
	t.Run("drop with nil options", func(t *testing.T) {
		err := client.Databases.Drop(ctx, databaseID, nil)
		require.NoError(t, err)
	})
}

/*
this test keeps failing need to fix.

	func TestInt_DatabasesUndrop(t *testing.T) {
		client := testClient(t)
		ctx := testContext(t)
		databaseTest, databaseCleanup := createDatabase(t, client)
		t.Cleanup(databaseCleanup)
		databaseID := databaseTest.ID()
		err := client.Databases.Drop(ctx, databaseID, nil)
		require.NoError(t, err)
		_, err = client.Databases.ShowByID(ctx, databaseID)
		require.Error(t, err)
		err = client.Databases.Undrop(ctx, databaseID)
		require.NoError(t, err)
		database, err := client.Databases.ShowByID(ctx, databaseID)
		require.NoError(t, err)
		assert.Equal(t, databaseID.Name(), database.Name)
	}
*/
func TestInt_DatabasesDescribe(t *testing.T) {
	client := testClient(t)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)
	ctx := testContext(t)
	databaseDetails, err := client.Databases.Describe(ctx, databaseTest.ID())
	require.NoError(t, err)
	rows := databaseDetails.Rows
	found := false
	for _, row := range rows {
		if row.Name == schemaTest.ID().Name() {
			found = true
		}
	}
	assert.True(t, found)
}

func TestInt_DatabasesAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("renaming", func(t *testing.T) {
		databaseTest, _ := createDatabase(t, client)
		newName := randomAccountObjectIdentifier(t)
		err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
			NewName: newName,
		})
		require.NoError(t, err)
		database, err := client.Databases.ShowByID(ctx, newName)
		assert.Equal(t, newName.Name(), database.Name)
		t.Cleanup(func() {
			err = client.Databases.Drop(ctx, newName, nil)
			require.NoError(t, err)
		})
	})

	t.Run("swap with another database", func(t *testing.T) {
		databaseTest, databaseCleanup := createDatabase(t, client)
		t.Cleanup(databaseCleanup)
		databaseTest2, databaseCleanup2 := createDatabase(t, client)
		t.Cleanup(databaseCleanup2)
		err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
			SwapWith: databaseTest2.ID(),
		})
		require.NoError(t, err)
	})

	t.Run("setting and unsetting retention time + comment ", func(t *testing.T) {
		databaseTest, _ := createDatabase(t, client)
		err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
			Set: &sdk.DatabaseSet{
				DataRetentionTimeInDays: sdk.Int(42),
				Comment:                 sdk.String("test comment"),
			},
		})
		require.NoError(t, err)
		database, err := client.Databases.ShowByID(ctx, databaseTest.ID())
		require.NoError(t, err)
		assert.Equal(t, 42, database.RetentionTime)
		assert.Equal(t, "test comment", database.Comment)
		err = client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
			Unset: &sdk.DatabaseUnset{
				DataRetentionTimeInDays: sdk.Bool(true),
				Comment:                 sdk.Bool(true),
			},
		})
		require.NoError(t, err)
		database, err = client.Databases.ShowByID(ctx, databaseTest.ID())
		require.NoError(t, err)
		assert.NotEqual(t, 42, database.RetentionTime)
		assert.Equal(t, "", database.Comment)
	})
}

func TestInt_AlterReplication(t *testing.T) {
	// todo: once ReplicationGroups are supported.
}

func TestInt_AlterFailover(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	secondaryClient := testSecondaryClient(t)

	toAccounts := []sdk.AccountIdentifier{
		getAccountIdentifier(t, secondaryClient),
	}
	t.Run("enable and disable failover", func(t *testing.T) {
		opts := &sdk.AlterDatabaseFailoverOptions{
			EnableFailover: &sdk.EnableFailover{
				ToAccounts: toAccounts,
			},
		}
		err := client.Databases.AlterFailover(ctx, databaseTest.ID(), opts)
		if strings.Contains(err.Error(), "Accounts enabled for failover must also be enabled for replication. Enable replication to account") {
			t.Skip("Skipping test because secondary account not enabled for replication")
		}
		require.NoError(t, err)
		opts = &sdk.AlterDatabaseFailoverOptions{
			DisableFailover: &sdk.DisableFailover{
				ToAccounts: toAccounts,
			},
		}
		err = client.Databases.AlterFailover(ctx, databaseTest.ID(), opts)
		require.NoError(t, err)
		opts = &sdk.AlterDatabaseFailoverOptions{
			Primary: sdk.Bool(true),
		}
		err = client.Databases.AlterFailover(ctx, databaseTest.ID(), opts)
		require.NoError(t, err)
	})
}

func TestInt_DatabasesShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	databaseTest2, databaseCleanup2 := createDatabase(t, client)
	t.Cleanup(databaseCleanup2)
	t.Run("without show options", func(t *testing.T) {
		databases, err := client.Databases.Show(ctx, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(databases), 2)
		databaseIDs := make([]sdk.AccountObjectIdentifier, len(databases))
		for i, database := range databases {
			databaseIDs[i] = database.ID()
		}
		assert.Contains(t, databaseIDs, databaseTest.ID())
		assert.Contains(t, databaseIDs, databaseTest2.ID())
	})

	t.Run("with terse", func(t *testing.T) {
		showOptions := &sdk.ShowDatabasesOptions{
			Terse: sdk.Bool(true),
			Like: &sdk.Like{
				Pattern: sdk.String(databaseTest.Name),
			},
		}
		databases, err := client.Databases.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(databases))
		database := databases[0]
		assert.Equal(t, databaseTest.Name, database.Name)
		assert.NotEmpty(t, database.CreatedOn)
		assert.Empty(t, database.DroppedOn)
		assert.Empty(t, database.Owner)
	})
	/*
	   this test keeps failing, need to fix
	   	t.Run("with history", func(t *testing.T) {
	   		// need to drop a database to test if the "dropped_on" column is populated
	   		databaseCleanup2()
	   		showOptions := &ShowDatabasesOptions{
	   			History: Bool(true),
	   			Like: &Like{
	   				Pattern: String(databaseTest2.Name),
	   			},
	   		}
	   		databases, err := client.Databases.Show(ctx, showOptions)
	   		require.NoError(t, err)
	   		assert.Equal(t, 1, len(databases))
	   		database := databases[0]
	   		assert.Equal(t, databaseTest2.Name, database.Name)
	   		assert.NotEmpty(t, database.DroppedOn)
	   	})
	*/
	t.Run("with like starts with", func(t *testing.T) {
		showOptions := &sdk.ShowDatabasesOptions{
			StartsWith: sdk.String(databaseTest.Name),
			LimitFrom: &sdk.LimitFrom{
				Rows: sdk.Int(1),
			},
		}
		databases, err := client.Databases.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(databases))
		database := databases[0]
		assert.Equal(t, databaseTest.Name, database.Name)
	})

	t.Run("when searching a non-existent database", func(t *testing.T) {
		showOptions := &sdk.ShowDatabasesOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		databases, err := client.Databases.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(databases))
	})
}
