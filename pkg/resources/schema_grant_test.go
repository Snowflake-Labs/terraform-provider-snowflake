package resources_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestSchemaGrant(t *testing.T) {
	r := require.New(t)
	err := resources.SchemaGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestSchemaGrantCreate(t *testing.T) {
	r := require.New(t)

	for _, testPriv := range []string{"USAGE", "MODIFY", "CREATE TAG"} {
		testPriv := testPriv
		in := map[string]interface{}{
			"schema_name":       "test-schema",
			"database_name":     "test-db",
			"privilege":         testPriv,
			"roles":             []interface{}{"test-role-1", "test-role-2"},
			"shares":            []interface{}{"test-share-1", "test-share-2"},
			"with_grant_option": true,
		}
		d := schema.TestResourceDataRaw(t, resources.SchemaGrant().Resource.Schema, in)
		r.NotNil(d)

		WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
			mock.ExpectExec(
				fmt.Sprintf(`^GRANT %s ON SCHEMA "test-db"."test-schema" TO ROLE "test-role-1" WITH GRANT OPTION$`, testPriv),
			).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(
				fmt.Sprintf(`^GRANT %s ON SCHEMA "test-db"."test-schema" TO ROLE "test-role-2" WITH GRANT OPTION$`, testPriv),
			).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(
				fmt.Sprintf(`^GRANT %s ON SCHEMA "test-db"."test-schema" TO SHARE "test-share-1" WITH GRANT OPTION$`, testPriv),
			).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(
				fmt.Sprintf(`^GRANT %s ON SCHEMA "test-db"."test-schema" TO SHARE "test-share-2" WITH GRANT OPTION$`, testPriv),
			).WillReturnResult(sqlmock.NewResult(1, 1))
			expectReadSchemaGrant(mock, testPriv)
			err := resources.CreateSchemaGrant(d, db)
			r.NoError(err)
		})
	}
}

func TestSchemaGrantRead(t *testing.T) {
	r := require.New(t)

	d := schemaGrant(t, "test-db|test-schema||USAGE||false", map[string]interface{}{
		"schema_name":       "test-schema",
		"database_name":     "test-db",
		"privilege":         "USAGE",
		"roles":             []interface{}{},
		"shares":            []interface{}{},
		"with_grant_option": false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadSchemaGrant(mock, "USAGE")
		err := resources.ReadSchemaGrant(d, db)
		r.NoError(err)
	})
	roles := d.Get("roles").(*schema.Set)
	r.True(roles.Contains("test-role-1"))
	r.True(roles.Contains("test-role-2"))
	r.Equal(2, roles.Len())

	shares := d.Get("shares").(*schema.Set)
	r.True(shares.Contains("test-share-1"))
	r.True(shares.Contains("test-share-2"))
	r.Equal(2, shares.Len())
}

func expectReadSchemaGrant(mock sqlmock.Sqlmock, testPriv string) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), testPriv, "SCHEMA", "test-schema", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), testPriv, "SCHEMA", "test-schema", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), testPriv, "SCHEMA", "test-schema", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), testPriv, "SCHEMA", "test-schema", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON SCHEMA "test-db"."test-schema"$`).WillReturnRows(rows)
}

func TestFutureSchemaGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"on_future":         true,
		"database_name":     "test-db",
		"privilege":         "USAGE",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.SchemaGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE SCHEMAS IN DATABASE "test-db" TO ROLE "test-role-1" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE SCHEMAS IN DATABASE "test-db" TO ROLE "test-role-2" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureSchemaGrant(mock)
		err := resources.CreateSchemaGrant(d, db)
		r.NoError(err)
	})
}

func expectReadFutureSchemaGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "SCHEMA", "test-db.<SCHEMA>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "SCHEMA", "test-db.<SCHEMA>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN DATABASE "test-db"$`).WillReturnRows(rows)
}

func TestParseSchemaGrantID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseSchemaGrantID("test-db|test-schema|USAGE|false|role1,role2|share1,share2")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("test-schema", grantID.SchemaName)
	r.Equal("USAGE", grantID.Privilege)
	r.Equal(false, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
	r.Equal(2, len(grantID.Shares))
	r.Equal("share1", grantID.Shares[0])
	r.Equal("share2", grantID.Shares[1])
}

func TestParseSchemaGrantEmojiID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseSchemaGrantID("test-db❄️test-schema❄️USAGE❄️false❄️role1,role2❄️share1,share2")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("test-schema", grantID.SchemaName)
	r.Equal("USAGE", grantID.Privilege)
	r.Equal(false, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
	r.Equal(2, len(grantID.Shares))
	r.Equal("share1", grantID.Shares[0])
	r.Equal("share2", grantID.Shares[1])
}

func TestParseSchemaGrantOldID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseSchemaGrantID("test-db|test-schema||USAGE|role1,role2|false")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("test-schema", grantID.SchemaName)
	r.Equal("USAGE", grantID.Privilege)
	r.Equal(false, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
	r.Equal(0, len(grantID.Shares))
}

func TestParseSchemaGrantReallyOldID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseSchemaGrantID("test-db|test-schema||CREATE TABLE|false")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("test-schema", grantID.SchemaName)
	r.Equal("CREATE TABLE", grantID.Privilege)
	r.Equal(false, grantID.WithGrantOption)
	r.Equal(0, len(grantID.Roles))
	r.Equal(0, len(grantID.Shares))
}
