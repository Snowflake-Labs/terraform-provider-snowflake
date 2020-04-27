package resources

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func Test_grantToRole(t *testing.T) {
	r := require.New(t)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`GRANT ROLE "foo" TO ROLE "bar"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := grantRoleToRole(db, "foo", "bar")
		r.NoError(err)
	})
}

func Test_grantToUser(t *testing.T) {
	r := require.New(t)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`GRANT ROLE "foo" TO USER "bar"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := grantRoleToUser(db, "foo", "bar")
		r.NoError(err)
	})
}

func Test_readGrants(t *testing.T) {
	r := require.New(t)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"created_on", "role", "granted_to", "grantee_name", "granted_by"}).AddRow("_", "foo", "ROLE", "bam", "")
		mock.ExpectQuery(`SHOW GRANTS OF ROLE "foo"`).WillReturnRows(rows)
		read, err := readGrants(db, "foo")
		r.NoError(err)
		r.Len(read, 1)
		g := read[0]
		r.Equal("ROLE", g.GrantedTo.String)
		r.Equal("bam", g.GranteeName.String)
	})
}

func Test_revokeRoleFromRole(t *testing.T) {
	r := require.New(t)
	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`REVOKE ROLE "foo" FROM ROLE "bar"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := revokeRoleFromRole(db, "foo", "bar")
		r.NoError(err)

	})

}
func Test_revokeRoleFromUser(t *testing.T) {
	r := require.New(t)
	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`REVOKE ROLE "foo" FROM USER "bar"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := revokeRoleFromUser(db, "foo", "bar")
		r.NoError(err)

	})

}
