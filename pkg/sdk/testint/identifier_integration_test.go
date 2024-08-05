package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_IdentifiersForOnePartIdentifierAsNameAndReference(t *testing.T) {
	testCases := []struct {
		Name     sdk.AccountObjectIdentifier
		ShowName string
		Error    string
	}{
		// special cases
		{Name: sdk.NewAccountObjectIdentifier(``), Error: "invalid object identifier"},
		{Name: sdk.NewAccountObjectIdentifier(`"`), Error: "invalid object identifier"},
		// This is a valid identifier, but because in NewXIdentifier functions we're trimming double quotes it won't work
		{Name: sdk.NewAccountObjectIdentifier(`""`), Error: "invalid object identifier"},
		// This is a valid identifier, but because in NewXIdentifier functions we're trimming double quotes it won't work
		{Name: sdk.NewAccountObjectIdentifier(`""""`), Error: "invalid object identifier"},
		{Name: sdk.NewAccountObjectIdentifier(`"."`), ShowName: `.`},

		// lower case
		{Name: sdk.NewAccountObjectIdentifier(`abc`), ShowName: `abc`},
		{Name: sdk.NewAccountObjectIdentifier(`ab.c`), ShowName: `ab.c`},
		{Name: sdk.NewAccountObjectIdentifier(`a"bc`), Error: `unexpected '"`},
		{Name: sdk.NewAccountObjectIdentifier(`"a""bc"`), ShowName: `a"bc`},

		// upper case
		{Name: sdk.NewAccountObjectIdentifier(`ABC`), ShowName: `ABC`},
		{Name: sdk.NewAccountObjectIdentifier(`AB.C`), ShowName: `AB.C`},
		{Name: sdk.NewAccountObjectIdentifier(`A"BC`), Error: `unexpected '"`},
		{Name: sdk.NewAccountObjectIdentifier(`"A""BC"`), ShowName: `A"BC`},

		// mixed case
		{Name: sdk.NewAccountObjectIdentifier(`AbC`), ShowName: `AbC`},
		{Name: sdk.NewAccountObjectIdentifier(`Ab.C`), ShowName: `Ab.C`},
		{Name: sdk.NewAccountObjectIdentifier(`A"bC`), Error: `unexpected '"`},
		{Name: sdk.NewAccountObjectIdentifier(`"A""bC"`), ShowName: `A"bC`},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("one part identifier name and reference for input: %s", testCase.Name.FullyQualifiedName()), func(t *testing.T) {
			ctx := context.Background()

			err := testClient(t).ResourceMonitors.Create(ctx, testCase.Name, new(sdk.CreateResourceMonitorOptions))
			if testCase.Error != "" {
				require.ErrorContains(t, err, testCase.Error)
			} else {
				t.Cleanup(testClientHelper().ResourceMonitor.DropResourceMonitorFunc(t, testCase.Name))
			}

			err = testClient(t).Warehouses.Create(ctx, testCase.Name, &sdk.CreateWarehouseOptions{
				ResourceMonitor: &testCase.Name,
			})
			if testCase.Error != "" {
				require.ErrorContains(t, err, testCase.Error)
			} else {
				require.NoError(t, err)
				t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, testCase.Name))
				var result struct {
					Name            string `db:"name"`
					ResourceMonitor string `db:"resource_monitor"`
				}
				err = testClient(t).QueryOneForTests(ctx, &result, fmt.Sprintf("SHOW WAREHOUSES LIKE '%s'", testCase.ShowName))
				require.NoError(t, err)

				// For one part identifiers, we expect Snowflake to return unescaped identifiers (just like the ones we used for SHOW)
				assert.Equal(t, testCase.ShowName, result.Name)
				assert.Equal(t, testCase.ShowName, result.ResourceMonitor)
			}
		})
	}
}

func TestInt_IdentifiersForTwoPartIdentifierAsReference(t *testing.T) {
	type RawGrantOutput struct {
		Name      string `db:"name"`
		Privilege string `db:"privilege"`
	}

	testCases := []struct {
		Name                            sdk.DatabaseObjectIdentifier
		OverrideExpectedSnowflakeOutput string
		Error                           string
	}{
		// special cases
		{Name: sdk.NewDatabaseObjectIdentifier(``, ``), Error: "invalid object identifier"},
		{Name: sdk.NewDatabaseObjectIdentifier(`"`, `"`), Error: "invalid object identifier"},
		// This is a valid identifier, but because in NewXIdentifier functions we're trimming double quotes it won't work
		{Name: sdk.NewDatabaseObjectIdentifier(`""`, `""`), Error: "invalid object identifier"},
		// This is a valid identifier, but because in NewXIdentifier functions we're trimming double quotes it won't work
		{Name: sdk.NewDatabaseObjectIdentifier(`""""`, `""""`), Error: "invalid object identifier"},
		{Name: sdk.NewDatabaseObjectIdentifier(`"."`, `"."`)},

		// lower case
		{Name: sdk.NewDatabaseObjectIdentifier(`abc`, `abc`)},
		{Name: sdk.NewDatabaseObjectIdentifier(`ab.c`, `ab.c`)},
		{Name: sdk.NewDatabaseObjectIdentifier(`a"bc`, `a"bc`), Error: `unexpected '"`},
		{Name: sdk.NewDatabaseObjectIdentifier(`"a""bc"`, `"a""bc"`)},

		// upper case
		{Name: sdk.NewDatabaseObjectIdentifier(`ABC`, `ABC`), OverrideExpectedSnowflakeOutput: `ABC.ABC`},
		{Name: sdk.NewDatabaseObjectIdentifier(`AB.C`, `AB.C`)},
		{Name: sdk.NewDatabaseObjectIdentifier(`A"BC`, `A"BC`), Error: `unexpected '"`},
		{Name: sdk.NewDatabaseObjectIdentifier(`"A""BC"`, `"A""BC"`)},

		// mixed case
		{Name: sdk.NewDatabaseObjectIdentifier(`AbC`, `AbC`)},
		{Name: sdk.NewDatabaseObjectIdentifier(`Ab.C`, `Ab.C`)},
		{Name: sdk.NewDatabaseObjectIdentifier(`A"bC`, `A"bC`), Error: `unexpected '"`},
		{Name: sdk.NewDatabaseObjectIdentifier(`"A""bC"`, `"A""bC"`)},
	}

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("two part identifier reference for input: %s", testCase.Name.FullyQualifiedName()), func(t *testing.T) {
			ctx := context.Background()

			err := testClient(t).Databases.Create(ctx, testCase.Name.DatabaseId(), new(sdk.CreateDatabaseOptions))
			if testCase.Error != "" {
				require.ErrorContains(t, err, testCase.Error)
			} else {
				t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, testCase.Name.DatabaseId()))
			}

			err = testClient(t).Schemas.Create(ctx, testCase.Name, new(sdk.CreateSchemaOptions))
			if testCase.Error != "" {
				require.ErrorContains(t, err, testCase.Error)
			} else {
				require.NoError(t, err)
				t.Cleanup(testClientHelper().Schema.DropSchemaFunc(t, testCase.Name))

				testClientHelper().Grant.GrantOnSchemaToAccountRole(t, testCase.Name, role.ID(), sdk.SchemaPrivilegeCreateTable)

				var grants []RawGrantOutput
				err = testClient(t).QueryForTests(ctx, &grants, fmt.Sprintf("SHOW GRANTS ON SCHEMA %s", testCase.Name.FullyQualifiedName()))
				require.NoError(t, err)

				createTableGrant, err := collections.FindOne(grants, func(output RawGrantOutput) bool { return output.Privilege == sdk.SchemaPrivilegeCreateTable.String() })
				require.NoError(t, err)

				// For two part identifiers, we expect Snowflake to return escaped identifiers with exception
				// to identifiers that don't have any lowercase character and special symbol in it.
				if testCase.OverrideExpectedSnowflakeOutput != "" {
					assert.Equal(t, testCase.OverrideExpectedSnowflakeOutput, createTableGrant.Name)
				} else {
					assert.Equal(t, testCase.Name.FullyQualifiedName(), createTableGrant.Name)
				}
			}
		})
	}
}
