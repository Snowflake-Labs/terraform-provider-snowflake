package resources

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// FIXME de-dupe
func withMockDb(t *testing.T, f func(*sql.DB, sqlmock.Sqlmock)) {
	a := assert.New(t)
	db, mock, err := sqlmock.New()
	defer db.Close()
	a.NoError(err)

	f(db, mock)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func Test_grantToRole(t *testing.T) {
	a := assert.New(t)

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("GRANT ROLE foo TO ROLE bar").WillReturnResult(sqlmock.NewResult(1, 1))
		err := grantRoleToRole(db, "foo", "bar")
		a.NoError(err)
	})
}

func Test_grantRoletoRoles(t *testing.T) {
	a := assert.New(t)

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`GRANT ROLE foo TO ROLE bar`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`GRANT ROLE foo TO ROLE bam`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := grantRoleToRoles(db, "foo", []string{"bar", "bam"})
		a.NoError(err)
	})

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`GRANT ROLE foo TO ROLE bar`).WillReturnError(errors.New("uh oh"))
		err := grantRoleToRoles(db, "foo", []string{"bar", "bam"})
		a.Error(err)
	})
}

func Test_grantToUser(t *testing.T) {
	a := assert.New(t)

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("GRANT ROLE foo TO USER bar").WillReturnResult(sqlmock.NewResult(1, 1))
		err := grantRoleToUser(db, "foo", "bar")
		a.NoError(err)
	})
}

func Test_grantRoletoUsers(t *testing.T) {
	a := assert.New(t)

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`GRANT ROLE foo TO USER bar`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`GRANT ROLE foo TO USER bam`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := grantRoleToUsers(db, "foo", []string{"bar", "bam"})
		a.NoError(err)
	})

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`GRANT ROLE foo TO USER bar`).WillReturnError(errors.New("uh oh"))
		err := grantRoleToUsers(db, "foo", []string{"bar", "bam"})
		a.Error(err)
	})
}

func Test_readGrants(t *testing.T) {
	a := assert.New(t)

	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"created_on", "role", "granted_to", "grantee_name", "granted_by"}).AddRow("_", "foo", "ROLE", "bam", "")
		mock.ExpectQuery(`SHOW GRANTS OF ROLE foo`).WillReturnRows(rows)
		r, err := readGrants(db, "foo")
		a.NoError(err)
		a.Len(r, 1)
		g := r[0]
		a.Equal("ROLE", g.GrantedTo.String)
		a.Equal("bam", g.GranteeName.String)
	})
}

func Test_revokeRoleFromRole(t *testing.T) {
	a := assert.New(t)
	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`REVOKE ROLE foo FROM ROLE bar`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := revokeRoleFromRole(db, "foo", "bar")
		a.NoError(err)

	})

}
func Test_revokeRoleFromUser(t *testing.T) {
	a := assert.New(t)
	withMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`REVOKE ROLE foo FROM USER bar`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := revokeRoleFromUser(db, "foo", "bar")
		a.NoError(err)

	})

}
