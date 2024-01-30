package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseGrantPrivilegesToShareId(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier string
		Expected   GrantPrivilegesToShareId
		Error      string
	}{
		{
			Name:       "grant privileges on database to share",
			Identifier: `account-name."share-name"|REFERENCE_USAGE|OnDatabase|"on-database-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"REFERENCE_USAGE"},
				Kind:       OnDatabaseShareGrantKind,
				Identifier: sdk.NewAccountObjectIdentifier("on-database-name"),
			},
		},
		{
			Name:       "grant privileges on schema to share",
			Identifier: `account-name."share-name"|USAGE|OnSchema|"on-database-name"."on-schema-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"USAGE"},
				Kind:       OnSchemaShareGrantKind,
				Identifier: sdk.NewDatabaseObjectIdentifier("on-database-name", "on-schema-name"),
			},
		},
		// TODO(SNOW-1021686): This is wrong and should be fixed (function's last part of identifier cannot be enclosed with quotes like that)
		//{
		//	Name:       "grant privileges on function to share",
		//	Identifier: `account-name."share-name"|USAGE|OnFunction|"on-database-name"."on-schema-name".on-function-name(INT, VARCHAR)`,
		//	Expected: GrantPrivilegesToShareId{
		//		ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
		//		Privileges: []string{"USAGE"},
		//		Kind:       OnFunctionShareGrantKind,
		//		Identifier: sdk.NewSchemaObjectIdentifier("on-database-name", "on-schema-name", "on-function-name(INT, VARCHAR)"),
		//	},
		//},
		{
			Name:       "grant privileges on table to share",
			Identifier: `account-name."share-name"|EVOLVE SCHEMA|OnTable|"on-database-name"."on-schema-name"."on-table-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"EVOLVE SCHEMA"},
				Kind:       OnTableShareGrantKind,
				Identifier: sdk.NewSchemaObjectIdentifier("on-database-name", "on-schema-name", "on-table-name"),
			},
		},
		{
			Name:       "grant privileges on all tables in schema to share",
			Identifier: `account-name."share-name"|EVOLVE SCHEMA,SELECT|OnAllTablesInSchema|"on-database-name"."on-schema-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"EVOLVE SCHEMA", "SELECT"},
				Kind:       OnAllTablesInSchemaShareGrantKind,
				Identifier: sdk.NewDatabaseObjectIdentifier("on-database-name", "on-schema-name"),
			},
		},
		{
			Name:       "grant privileges on tag to share",
			Identifier: `account-name."share-name"|READ|OnTag|"on-tag-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"READ"},
				Kind:       OnTagShareGrantKind,
				Identifier: sdk.NewAccountObjectIdentifier("on-tag-name"),
			},
		},
		{
			Name:       "grant privileges on view to share",
			Identifier: `account-name."share-name"|READ|OnView|"on-database-name"."on-schema-name"."on-view-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"READ"},
				Kind:       OnViewShareGrantKind,
				Identifier: sdk.NewSchemaObjectIdentifier("on-database-name", "on-schema-name", "on-view-name"),
			},
		},
		{
			Name:       "validation: not enough parts",
			Identifier: `account-name."share-name"|SELECT|OnDatabase`,
			Error:      `snowflake_grant_privileges_to_share id is composed out of 4 parts "<account_name>.<share_name>|<privileges>|<grant_on_type>|<grant_on_identifier>", but got 3 parts: [account-name."share-name" SELECT OnDatabase]`,
		},
		{
			Name:       "validation: empty privileges",
			Identifier: `account-name."share-name"||OnDatabase|"database-name"`,
			Error:      `invalid Privileges value: [], should be comma separated list of privileges`,
		},
		{
			Name:       "validation: unsupported kind",
			Identifier: `account-name."share-name"|SELECT|OnSomething|"object-name"`,
			Error:      `unexpected share grant kind: OnSomething`,
		},
		{
			Name:       "validation: invalid identifier",
			Identifier: `account-name."share-name"|SELECT|OnDatabase|one.two.three.four.five.six.seven.eight.nine.ten`,
			Error:      `unable to classify identifier: one.two.three.four.five.six.seven.eight.nine.ten`,
		},
		{
			Name:       "validation: invalid account object identifier",
			Identifier: `account-name."share-name"|SELECT|OnDatabase|one.two`,
			Error:      `invalid identifier, expected fully qualified name of account object: <name>, but instead got: <database_name>.<name>`,
		},
		{
			Name:       "validation: invalid database object identifier",
			Identifier: `account-name."share-name"|SELECT|OnSchema|one.two.three`,
			Error:      `invalid identifier, expected fully qualified name of database object: <database_name>.<name>, but instead got: <database_name>.<schema_name>.<name>`,
		},
		{
			Name:       "validation: invalid schema object identifier",
			Identifier: `account-name."share-name"|SELECT|OnTable|one`,
			Error:      `invalid identifier, expected fully qualified name of schema object: <database_name>.<schema_name>.<name>, but instead got: <name>`,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			id, err := ParseGrantPrivilegesToShareId(tt.Identifier)
			if tt.Error == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.Expected, id)
			} else {
				assert.ErrorContains(t, err, tt.Error)
			}
		})
	}
}

func TestGrantPrivilegesToShareIdString(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier GrantPrivilegesToShareId
		Expected   string
		Error      string
	}{
		{
			Name: "grant privileges on database to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"REFERENCE_USAGE"},
				Kind:       OnDatabaseShareGrantKind,
				Identifier: sdk.NewAccountObjectIdentifier("database-name"),
			},
			Expected: `account-name."share-name"|REFERENCE_USAGE|OnDatabase|"database-name"`,
		},
		{
			Name: "grant privileges on schema to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"USAGE"},
				Kind:       OnSchemaShareGrantKind,
				Identifier: sdk.NewDatabaseObjectIdentifier("database-name", "schema-name"),
			},
			Expected: `account-name."share-name"|USAGE|OnSchema|"database-name"."schema-name"`,
		},
		// TODO(SNOW-1021686): This is wrong and should be fixed (function's last part of identifier cannot be enclosed with quotes like that)
		//{
		//	Name: "grant privileges on function to share",
		//	Identifier: GrantPrivilegesToShareId{
		//		ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
		//		Privileges: []string{"USAGE"},
		//		Kind:       OnFunctionShareGrantKind,
		//		Identifier: sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "function-name(INT, VARCHAR)"),
		//	},
		//	Expected: `account-name."share-name"|USAGE|OnFunction|"database-name"."schema-name".\"function-name(INT, VARCHAR)\"`,
		//},
		{
			Name: "grant privileges on table to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"EVOLVE SCHEMA", "SELECT"},
				Kind:       OnTableShareGrantKind,
				Identifier: sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "table-name"),
			},
			Expected: `account-name."share-name"|EVOLVE SCHEMA,SELECT|OnTable|"database-name"."schema-name"."table-name"`,
		},
		{
			Name: "grant privileges on all tables in schema to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"EVOLVE SCHEMA", "SELECT"},
				Kind:       OnAllTablesInSchemaShareGrantKind,
				Identifier: sdk.NewDatabaseObjectIdentifier("database-name", "schema-name"),
			},
			Expected: `account-name."share-name"|EVOLVE SCHEMA,SELECT|OnAllTablesInSchema|"database-name"."schema-name"`,
		},
		{
			Name: "grant privileges on tag to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"READ"},
				Kind:       OnTagShareGrantKind,
				Identifier: sdk.NewAccountObjectIdentifier("tag-name"),
			},
			Expected: `account-name."share-name"|READ|OnTag|"tag-name"`,
		},
		{
			Name: "grant privileges on view to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewExternalObjectIdentifierFromFullyQualifiedName("account-name.share-name"),
				Privileges: []string{"SELECT"},
				Kind:       OnViewShareGrantKind,
				Identifier: sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "view-name"),
			},
			Expected: `account-name."share-name"|SELECT|OnView|"database-name"."schema-name"."view-name"`,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Identifier.String())
		})
	}
}
