package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
)

func TestDatabaseGrant(t *testing.T) {
	r := require.New(t)
	err := resources.DatabaseGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestDatabaseGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"database_name":     "test-database",
		"privilege":         "USAGE",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"shares":            []interface{}{"test-share-1", "test-share-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.DatabaseGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO ROLE "test-role-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO ROLE "test-role-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO SHARE "test-share-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO SHARE "test-share-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadDatabaseGrant(mock)
		err := resources.CreateDatabaseGrant(d, db)
		r.NoError(err)
	})
}

func TestDatabaseGrantRead(t *testing.T) {
	r := require.New(t)

	d := databaseGrant(t, "test-database|||USAGE||false", map[string]interface{}{
		"database_name":     "test-database",
		"privilege":         "USAGE",
		"roles":             []interface{}{},
		"shares":            []interface{}{},
		"with_grant_option": false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadDatabaseGrant(mock)
		err := resources.ReadDatabaseGrant(d, db)
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

func expectReadDatabaseGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "DATABASE", "test-database", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "DATABASE", "test-database", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "DATABASE", "test-database", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "DATABASE", "test-database", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON DATABASE "test-database"$`).WillReturnRows(rows)
}

func TestParseDatabaseGrantID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseDatabaseGrantID("test-database|USAGE|true|role1,role2|share1,share2")
	r.NoError(err)
	r.Equal("test-database", grantID.DatabaseName)
	r.Equal("USAGE", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
	r.Equal(2, len(grantID.Shares))
	r.Equal("share1", grantID.Shares[0])
	r.Equal("share2", grantID.Shares[1])
}

func TestParseDatabaseGrantEmojiID(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseDatabaseGrantID("test-database❄️USAGE❄️true❄️role1,role2❄️share1,share2")
	r.NoError(err)
	r.Equal("test-database", grantID.DatabaseName)
	r.Equal("USAGE", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
	r.Equal(2, len(grantID.Shares))
	r.Equal("share1", grantID.Shares[0])
	r.Equal("share2", grantID.Shares[1])
}

func TestParseDatabaseGrantOldID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseDatabaseGrantID("test-database|||USAGE|role1,role2|true")
	r.NoError(err)
	r.Equal("test-database", grantID.DatabaseName)
	r.Equal("USAGE", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}
