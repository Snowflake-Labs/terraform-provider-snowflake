package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DatabasesCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("minimal", func(t *testing.T) {
		databaseID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Databases.Create(ctx, databaseID, &sdk.CreateDatabaseOptions{
			OrReplace: sdk.Bool(true),
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, databaseID))

		database, err := client.Databases.ShowByID(ctx, databaseID)
		require.NoError(t, err)
		assert.Equal(t, databaseID.Name(), database.Name)
	})

	t.Run("as clone", func(t *testing.T) {
		cloneDatabase, cloneDatabaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cloneDatabaseCleanup)

		databaseID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Databases.Create(ctx, databaseID, &sdk.CreateDatabaseOptions{
			Clone: &sdk.Clone{
				SourceObject: cloneDatabase.ID(),
				At: &sdk.TimeTravel{
					Offset: sdk.Int(0),
				},
			},
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, databaseID))

		database, err := client.Databases.ShowByID(ctx, databaseID)
		require.NoError(t, err)
		assert.Equal(t, databaseID.Name(), database.Name)
	})

	t.Run("complete", func(t *testing.T) {
		databaseID := testClientHelper().Ids.RandomAccountObjectIdentifier()

		// new database and schema created on purpose
		databaseTest, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		schemaTest, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, databaseTest.ID())
		t.Cleanup(schemaCleanup)

		tagTest, tagCleanup := testClientHelper().Tag.CreateTagInSchema(t, schemaTest.ID())
		t.Cleanup(tagCleanup)

		tag2Test, tag2Cleanup := testClientHelper().Tag.CreateTagInSchema(t, schemaTest.ID())
		t.Cleanup(tag2Cleanup)

		externalVolume, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(externalVolumeCleanup)

		catalog, catalogCleanup := testClientHelper().CatalogIntegration.Create(t)
		t.Cleanup(catalogCleanup)

		comment := random.Comment()
		err := client.Databases.Create(ctx, databaseID, &sdk.CreateDatabaseOptions{
			Transient:                  sdk.Bool(true),
			IfNotExists:                sdk.Bool(true),
			DataRetentionTimeInDays:    sdk.Int(1),
			MaxDataExtensionTimeInDays: sdk.Int(1),
			ExternalVolume:             &externalVolume,
			Catalog:                    &catalog,
			DefaultDDLCollation:        sdk.String("en_US"),
			LogLevel:                   sdk.Pointer(sdk.LogLevelInfo),
			TraceLevel:                 sdk.Pointer(sdk.TraceLevelOnEvent),
			Comment:                    sdk.String(comment),
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
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, databaseID))

		database, err := client.Databases.ShowByID(ctx, databaseID)
		require.NoError(t, err)
		assert.Equal(t, databaseID.Name(), database.Name)
		assert.Equal(t, comment, database.Comment)
		assert.Equal(t, 1, database.RetentionTime)

		param, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterMaxDataExtensionTimeInDays, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: databaseID})
		assert.NoError(t, err)
		assert.Equal(t, "1", param.Value)

		externalVolumeParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterExternalVolume, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: databaseID})
		assert.NoError(t, err)
		assert.Equal(t, externalVolume.Name(), externalVolumeParam.Value)

		catalogParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterCatalog, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: databaseID})
		assert.NoError(t, err)
		assert.Equal(t, catalog.Name(), catalogParam.Value)

		logLevelParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterLogLevel, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: databaseID})
		assert.NoError(t, err)
		assert.Equal(t, string(sdk.LogLevelInfo), logLevelParam.Value)

		traceLevelParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterTraceLevel, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: databaseID})
		assert.NoError(t, err)
		assert.Equal(t, string(sdk.TraceLevelOnEvent), traceLevelParam.Value)

		tag1Value, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), database.ID(), sdk.ObjectTypeDatabase)
		require.NoError(t, err)
		assert.Equal(t, "v1", tag1Value)

		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2Test.ID(), database.ID(), sdk.ObjectTypeDatabase)
		require.NoError(t, err)
		assert.Equal(t, "v2", tag2Value)
	})
}

func TestInt_DatabasesCreateShared(t *testing.T) {
	client := testClient(t)
	secondaryClient := testSecondaryClient(t)
	ctx := testContext(t)

	databaseTest, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, databaseTest.ID())
	t.Cleanup(schemaCleanup)

	testTag, testTagCleanup := testClientHelper().Tag.CreateTagInSchema(t, schemaTest.ID())
	t.Cleanup(testTagCleanup)

	externalVolume, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalog, catalogCleanup := testClientHelper().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	// prepare a database on the secondary account
	shareTest, shareCleanup := secondaryTestClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	sharedDatabase, sharedDatabaseCleanup := secondaryTestClientHelper().Database.CreateDatabase(t)
	t.Cleanup(sharedDatabaseCleanup)

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

	comment := random.Comment()
	err = client.Databases.CreateShared(ctx, sharedDatabase.ID(), shareTest.ExternalID(), &sdk.CreateSharedDatabaseOptions{
		Transient:           sdk.Bool(true),
		IfNotExists:         sdk.Bool(true),
		ExternalVolume:      &externalVolume,
		Catalog:             &catalog,
		DefaultDDLCollation: sdk.String("en_US"),
		LogLevel:            sdk.Pointer(sdk.LogLevelDebug),
		TraceLevel:          sdk.Pointer(sdk.TraceLevelAlways),
		Comment:             sdk.String(comment),
		Tag: []sdk.TagAssociation{
			{
				Name:  testTag.ID(),
				Value: "v1",
			},
		},
	})
	require.NoError(t, err)
	t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, sharedDatabase.ID()))

	database, err := client.Databases.ShowByID(ctx, sharedDatabase.ID())
	require.NoError(t, err)

	assert.Equal(t, sharedDatabase.ID().Name(), database.Name)
	assert.Equal(t, comment, database.Comment)

	externalVolumeParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterExternalVolume, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: sharedDatabase.ID()})
	assert.NoError(t, err)
	assert.Equal(t, externalVolume.Name(), externalVolumeParam.Value)

	catalogParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterCatalog, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: sharedDatabase.ID()})
	assert.NoError(t, err)
	assert.Equal(t, catalog.Name(), catalogParam.Value)

	logLevelParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterLogLevel, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: sharedDatabase.ID()})
	assert.NoError(t, err)
	assert.Equal(t, string(sdk.LogLevelDebug), logLevelParam.Value)

	traceLevelParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterTraceLevel, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: sharedDatabase.ID()})
	assert.NoError(t, err)
	assert.Equal(t, string(sdk.TraceLevelAlways), traceLevelParam.Value)

	tag1Value, err := client.SystemFunctions.GetTag(ctx, testTag.ID(), database.ID(), sdk.ObjectTypeDatabase)
	require.NoError(t, err)
	assert.Equal(t, "v1", tag1Value)
}

func TestInt_DatabasesCreateSecondary(t *testing.T) {
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

	externalVolume, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalog, catalogCleanup := testClientHelper().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	externalDatabaseId := sdk.NewExternalObjectIdentifier(secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(), sharedDatabase.ID())
	comment := random.Comment()
	err = client.Databases.CreateSecondary(ctx, sharedDatabase.ID(), externalDatabaseId, &sdk.CreateSecondaryDatabaseOptions{
		IfNotExists:                sdk.Bool(true),
		DataRetentionTimeInDays:    sdk.Int(1),
		MaxDataExtensionTimeInDays: sdk.Int(10),
		ExternalVolume:             &externalVolume,
		Catalog:                    &catalog,
		DefaultDDLCollation:        sdk.String("en_US"),
		LogLevel:                   sdk.Pointer(sdk.LogLevelDebug),
		TraceLevel:                 sdk.Pointer(sdk.TraceLevelAlways),
		Comment:                    sdk.String(comment),
	})
	require.NoError(t, err)
	t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, sharedDatabase.ID()))

	database, err := client.Databases.ShowByID(ctx, sharedDatabase.ID())
	require.NoError(t, err)

	assert.Equal(t, sharedDatabase.ID().Name(), database.Name)
	assert.Equal(t, 1, database.RetentionTime)
	assert.Equal(t, comment, database.Comment)

	param, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterMaxDataExtensionTimeInDays, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: sharedDatabase.ID()})
	assert.NoError(t, err)
	assert.Equal(t, "10", param.Value)

	externalVolumeParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterExternalVolume, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: sharedDatabase.ID()})
	assert.NoError(t, err)
	assert.Equal(t, externalVolume.Name(), externalVolumeParam.Value)

	catalogParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterCatalog, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: sharedDatabase.ID()})
	assert.NoError(t, err)
	assert.Equal(t, catalog.Name(), catalogParam.Value)

	logLevelParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterLogLevel, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: sharedDatabase.ID()})
	assert.NoError(t, err)
	assert.Equal(t, string(sdk.LogLevelDebug), logLevelParam.Value)

	traceLevelParam, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterTraceLevel, sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: sharedDatabase.ID()})
	assert.NoError(t, err)
	assert.Equal(t, string(sdk.TraceLevelAlways), traceLevelParam.Value)
}

func TestInt_DatabasesAlter(t *testing.T) {
	client := testClient(t)
	secondaryClient := testSecondaryClient(t)
	ctx := testContext(t)

	queryParameterValueForDatabase := func(t *testing.T, id sdk.AccountObjectIdentifier, parameter sdk.ObjectParameter) string {
		t.Helper()
		value, err := client.Parameters.ShowObjectParameter(ctx, parameter, sdk.Object{
			ObjectType: sdk.ObjectTypeDatabase,
			Name:       id,
		})
		require.NoError(t, err)
		return value.Value
	}

	testCases := []struct {
		DatabaseType string
		CreateFn     func(t *testing.T) (*sdk.Database, func())
	}{
		{
			DatabaseType: "Normal",
			CreateFn: func(t *testing.T) (*sdk.Database, func()) {
				t.Helper()
				return testClientHelper().Database.CreateDatabase(t)
			},
		},
		{
			DatabaseType: "From Share",
			CreateFn: func(t *testing.T) (*sdk.Database, func()) {
				t.Helper()

				shareTest, shareCleanup := secondaryTestClientHelper().Share.CreateShare(t)
				t.Cleanup(shareCleanup)

				sharedDatabase, sharedDatabaseCleanup := secondaryTestClientHelper().Database.CreateDatabase(t)
				t.Cleanup(sharedDatabaseCleanup)

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

				err = client.Databases.CreateShared(ctx, sharedDatabase.ID(), shareTest.ExternalID(), &sdk.CreateSharedDatabaseOptions{})
				require.NoError(t, err)

				database, err := client.Databases.ShowByID(ctx, sharedDatabase.ID())
				require.NoError(t, err)

				return database, testClientHelper().Database.DropDatabaseFunc(t, sharedDatabase.ID())
			},
		},
		{
			DatabaseType: "Replica",
			CreateFn: func(t *testing.T) (*sdk.Database, func()) {
				t.Helper()

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
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Database: %s - renaming", testCase.DatabaseType), func(t *testing.T) {
			databaseTest, databaseTestCleanup := testCase.CreateFn(t)
			t.Cleanup(databaseTestCleanup)
			newName := testClientHelper().Ids.RandomAccountObjectIdentifier()

			err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				NewName: &newName,
			})
			require.NoError(t, err)
			t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, newName))

			database, err := client.Databases.ShowByID(ctx, newName)
			require.NoError(t, err)
			assert.Equal(t, newName.Name(), database.Name)
		})

		t.Run(fmt.Sprintf("Database: %s - setting and unsetting log_level and trace_level", testCase.DatabaseType), func(t *testing.T) {
			if testCase.DatabaseType == "From Share" {
				t.Skipf("Skipping database test because from share is not supported")
			}

			databaseTest, databaseTestCleanup := testCase.CreateFn(t)
			t.Cleanup(databaseTestCleanup)

			err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Set: &sdk.DatabaseSet{
					LogLevel:   sdk.Pointer(sdk.LogLevelInfo),
					TraceLevel: sdk.Pointer(sdk.TraceLevelOnEvent),
				},
			})
			require.NoError(t, err)

			require.Equal(t, string(sdk.LogLevelInfo), queryParameterValueForDatabase(t, databaseTest.ID(), sdk.ObjectParameterLogLevel))
			require.Equal(t, string(sdk.TraceLevelOnEvent), queryParameterValueForDatabase(t, databaseTest.ID(), sdk.ObjectParameterTraceLevel))

			err = client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Unset: &sdk.DatabaseUnset{
					LogLevel:   sdk.Bool(true),
					TraceLevel: sdk.Bool(true),
				},
			})
			require.NoError(t, err)

			require.Equal(t, string(sdk.LogLevelOff), queryParameterValueForDatabase(t, databaseTest.ID(), sdk.ObjectParameterLogLevel))
			require.Equal(t, string(sdk.TraceLevelOff), queryParameterValueForDatabase(t, databaseTest.ID(), sdk.ObjectParameterTraceLevel))
		})

		t.Run(fmt.Sprintf("Database: %s - setting and unsetting external volume and catalog", testCase.DatabaseType), func(t *testing.T) {
			if testCase.DatabaseType == "From Share" {
				t.Skipf("Skipping database test because from share is not supported")
			}

			databaseTest, databaseTestCleanup := testCase.CreateFn(t)
			t.Cleanup(databaseTestCleanup)

			externalVolumeTest, externalVolumeTestCleanup := testClientHelper().ExternalVolume.Create(t)
			t.Cleanup(externalVolumeTestCleanup)

			catalogIntegrationTest, catalogIntegrationTestCleanup := testClientHelper().CatalogIntegration.Create(t)
			t.Cleanup(catalogIntegrationTestCleanup)

			err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Set: &sdk.DatabaseSet{
					ExternalVolume: &externalVolumeTest,
					Catalog:        &catalogIntegrationTest,
				},
			})
			require.NoError(t, err)
			require.Equal(t, externalVolumeTest.Name(), queryParameterValueForDatabase(t, databaseTest.ID(), sdk.ObjectParameterExternalVolume))
			require.Equal(t, catalogIntegrationTest.Name(), queryParameterValueForDatabase(t, databaseTest.ID(), sdk.ObjectParameterCatalog))

			err = client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Unset: &sdk.DatabaseUnset{
					ExternalVolume: sdk.Bool(true),
					Catalog:        sdk.Bool(true),
				},
			})
			require.NoError(t, err)
			require.Empty(t, queryParameterValueForDatabase(t, databaseTest.ID(), sdk.ObjectParameterExternalVolume))
			require.Empty(t, queryParameterValueForDatabase(t, databaseTest.ID(), sdk.ObjectParameterCatalog))
		})

		t.Run(fmt.Sprintf("Database: %s - setting and unsetting retention time", testCase.DatabaseType), func(t *testing.T) {
			if testCase.DatabaseType == "From Share" {
				t.Skipf("Skipping database test because from share is not supported")
			}

			databaseTest, databaseTestCleanup := testCase.CreateFn(t)
			t.Cleanup(databaseTestCleanup)

			err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Set: &sdk.DatabaseSet{
					DataRetentionTimeInDays: sdk.Int(42),
				},
			})
			require.NoError(t, err)

			database, err := client.Databases.ShowByID(ctx, databaseTest.ID())
			require.NoError(t, err)
			assert.Equal(t, 42, database.RetentionTime)

			err = client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Unset: &sdk.DatabaseUnset{
					DataRetentionTimeInDays: sdk.Bool(true),
				},
			})
			require.NoError(t, err)

			database, err = client.Databases.ShowByID(ctx, databaseTest.ID())
			require.NoError(t, err)
			assert.NotEqual(t, 42, database.RetentionTime)
		})

		t.Run(fmt.Sprintf("Database: %s - setting and unsetting comment", testCase.DatabaseType), func(t *testing.T) {
			databaseTest, databaseTestCleanup := testCase.CreateFn(t)
			t.Cleanup(databaseTestCleanup)

			err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Set: &sdk.DatabaseSet{
					Comment: sdk.String("test comment"),
				},
			})
			require.NoError(t, err)

			database, err := client.Databases.ShowByID(ctx, databaseTest.ID())
			require.NoError(t, err)

			assert.Equal(t, "test comment", database.Comment)

			err = client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Unset: &sdk.DatabaseUnset{
					Comment: sdk.Bool(true),
				},
			})
			require.NoError(t, err)

			database, err = client.Databases.ShowByID(ctx, databaseTest.ID())
			require.NoError(t, err)
			assert.Equal(t, "", database.Comment)
		})

		t.Run(fmt.Sprintf("Database: %s - swap with another database", testCase.DatabaseType), func(t *testing.T) {
			databaseTest, databaseTestCleanup := testCase.CreateFn(t)
			t.Cleanup(databaseTestCleanup)

			databaseTest2, databaseCleanup2 := testClientHelper().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup2)

			err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				SwapWith: sdk.Pointer(databaseTest2.ID()),
			})
			require.NoError(t, err)
		})
	}
}

func TestInt_DatabasesAlterReplication(t *testing.T) {
	t.Run("enable and disable replication", func(t *testing.T) {
		ctx := testContext(t)

		database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		err := testClient(t).Databases.AlterReplication(ctx, database.ID(), &sdk.AlterDatabaseReplicationOptions{
			EnableReplication: &sdk.EnableReplication{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
				IgnoreEditionCheck: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		err = testClient(t).Databases.AlterReplication(ctx, database.ID(), &sdk.AlterDatabaseReplicationOptions{
			DisableReplication: &sdk.DisableReplication{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("refresh replicated database", func(t *testing.T) {
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

		externalVolume, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(externalVolumeCleanup)

		catalog, catalogCleanup := testClientHelper().CatalogIntegration.Create(t)
		t.Cleanup(catalogCleanup)

		externalDatabaseId := sdk.NewExternalObjectIdentifier(secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(), sharedDatabase.ID())
		comment := random.Comment()
		err = client.Databases.CreateSecondary(ctx, sharedDatabase.ID(), externalDatabaseId, &sdk.CreateSecondaryDatabaseOptions{
			IfNotExists:                sdk.Bool(true),
			DataRetentionTimeInDays:    sdk.Int(1),
			MaxDataExtensionTimeInDays: sdk.Int(10),
			ExternalVolume:             &externalVolume,
			Catalog:                    &catalog,
			DefaultDDLCollation:        sdk.String("en_US"),
			LogLevel:                   sdk.Pointer(sdk.LogLevelDebug),
			TraceLevel:                 sdk.Pointer(sdk.TraceLevelAlways),
			Comment:                    sdk.String(comment),
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, sharedDatabase.ID()))

		err = secondaryClient.Databases.Alter(ctx, sharedDatabase.ID(), &sdk.AlterDatabaseOptions{
			Set: &sdk.DatabaseSet{
				Comment: sdk.String("some comment"),
			},
		})
		require.NoError(t, err)

		database, err := client.Databases.ShowByID(ctx, sharedDatabase.ID())
		require.NoError(t, err)

		assert.Equal(t, sharedDatabase.ID().Name(), database.Name)
		assert.Equal(t, 1, database.RetentionTime)
		assert.Equal(t, comment, database.Comment)

		err = client.Databases.AlterReplication(ctx, sharedDatabase.ID(), &sdk.AlterDatabaseReplicationOptions{
			Refresh: sdk.Bool(true),
		})
		require.NoError(t, err)

		database, err = client.Databases.ShowByID(ctx, sharedDatabase.ID())
		require.NoError(t, err)

		assert.Equal(t, sharedDatabase.ID().Name(), database.Name)
		assert.Equal(t, 1, database.RetentionTime)
		assert.Equal(t, comment, database.Comment)
	})
}

func TestInt_DatabasesAlterFailover(t *testing.T) {
	t.Run("enable and disable failover", func(t *testing.T) {
		ctx := testContext(t)

		database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		err := testClient(t).Databases.AlterReplication(ctx, database.ID(), &sdk.AlterDatabaseReplicationOptions{
			EnableReplication: &sdk.EnableReplication{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
				IgnoreEditionCheck: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(ctx, database.ID(), &sdk.AlterDatabaseFailoverOptions{
			EnableFailover: &sdk.EnableFailover{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
			},
		})
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(ctx, database.ID(), &sdk.AlterDatabaseFailoverOptions{
			DisableFailover: &sdk.DisableFailover{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("promote to primary", func(t *testing.T) {
		t.Skipf("Can be unskipped after [SNOW-1002023]. CI Snowflake Edition doesn't support this feature")

		ctx := testContext(t)

		database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		err := testClient(t).Databases.AlterReplication(ctx, database.ID(), &sdk.AlterDatabaseReplicationOptions{
			EnableReplication: &sdk.EnableReplication{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
				IgnoreEditionCheck: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(ctx, database.ID(), &sdk.AlterDatabaseFailoverOptions{
			EnableFailover: &sdk.EnableFailover{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
			},
		})
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(ctx, database.ID(), &sdk.AlterDatabaseFailoverOptions{
			Primary: sdk.Bool(true),
		})
		require.NoError(t, err)
	})
}

func TestInt_DatabasesDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("drop with nil options", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseTestCleanup)

		err := client.Databases.Drop(ctx, databaseTest.ID(), nil)
		require.NoError(t, err)
	})

	t.Run("drop if exists", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		databaseTestCleanup()

		err := client.Databases.Drop(ctx, databaseTest.ID(), &sdk.DropDatabaseOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	})

	t.Run("drop with cascade", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseTestCleanup)

		err := client.Databases.Drop(ctx, databaseTest.ID(), &sdk.DropDatabaseOptions{
			IfExists: sdk.Bool(true),
			Cascade:  sdk.Bool(true),
		})
		require.NoError(t, err)
	})

	t.Run("drop with restrict", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseTestCleanup)

		err := client.Databases.Drop(ctx, databaseTest.ID(), &sdk.DropDatabaseOptions{
			IfExists: sdk.Bool(true),
			Restrict: sdk.Bool(true),
		})
		require.NoError(t, err)
	})
}

func TestInt_DatabasesUndrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
	databaseCleanup()

	_, err := client.Databases.ShowByID(ctx, databaseTest.ID())
	require.Error(t, err)

	err = client.Databases.Undrop(ctx, databaseTest.ID())
	require.NoError(t, err)

	database, err := client.Databases.ShowByID(ctx, databaseTest.ID())
	require.NoError(t, err)

	assert.Equal(t, databaseTest.Name, database.Name)
}

func TestInt_DatabasesShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	databaseTest2, databaseCleanup2 := testClientHelper().Database.CreateDatabase(t)
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
		assert.Equal(t, "ROLE", databases[0].OwnerRoleType)
	})

	t.Run("with terse", func(t *testing.T) {
		databases, err := client.Databases.Show(ctx, &sdk.ShowDatabasesOptions{
			Terse: sdk.Bool(true),
			Like: &sdk.Like{
				Pattern: sdk.String(databaseTest.Name),
			},
		})
		require.NoError(t, err)

		database, err := collections.FindOne(databases, func(database sdk.Database) bool { return database.Name == databaseTest.Name })
		require.NoError(t, err)

		assert.Equal(t, databaseTest.Name, database.Name)
		assert.NotEmpty(t, database.CreatedOn)
		assert.Empty(t, database.DroppedOn)
		assert.Empty(t, database.Owner)
	})

	t.Run("with history", func(t *testing.T) {
		databaseTest3, databaseCleanup3 := testClientHelper().Database.CreateDatabase(t)
		databaseCleanup3()

		databases, err := client.Databases.Show(ctx, &sdk.ShowDatabasesOptions{
			History: sdk.Bool(true),
			Like: &sdk.Like{
				Pattern: sdk.String(databaseTest3.Name),
			},
		})
		require.NoError(t, err)

		droppedDatabase, err := collections.FindOne(databases, func(database sdk.Database) bool { return database.Name == databaseTest3.Name })
		require.NoError(t, err)

		assert.Equal(t, databaseTest3.Name, droppedDatabase.Name)
		assert.NotEmpty(t, droppedDatabase.DroppedOn)
	})

	t.Run("with like starts with", func(t *testing.T) {
		databases, err := client.Databases.Show(ctx, &sdk.ShowDatabasesOptions{
			StartsWith: sdk.String(databaseTest.Name),
			LimitFrom: &sdk.LimitFrom{
				Rows: sdk.Int(1),
			},
		})
		require.NoError(t, err)

		database, err := collections.FindOne(databases, func(database sdk.Database) bool { return database.Name == databaseTest.Name })
		require.NoError(t, err)

		assert.Equal(t, databaseTest.Name, database.Name)
	})

	t.Run("when searching a non-existent database", func(t *testing.T) {
		databases, err := client.Databases.Show(ctx, &sdk.ShowDatabasesOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		})
		require.NoError(t, err)

		assert.Equal(t, 0, len(databases))
	})
}

func TestInt_DatabasesDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertContainsSchema := func(details *sdk.DatabaseDetails, schemaName string) {
		_, err := collections.FindOne(details.Rows, func(row sdk.DatabaseDetailsRow) bool { return row.Kind == "SCHEMA" && row.Name == schemaName })
		assert.NoError(t, err)
	}

	schemaTest, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	databaseDetails, err := client.Databases.Describe(ctx, schemaTest.ID().DatabaseId())
	require.NoError(t, err)

	assertContainsSchema(databaseDetails, schemaTest.ID().Name())
	assertContainsSchema(databaseDetails, "INFORMATION_SCHEMA")
	assertContainsSchema(databaseDetails, "PUBLIC")
}
