package resources_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestStageGrant(t *testing.T) {
	r := require.New(t)
	err := resources.StageGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestStageGrantCreate(t *testing.T) {
	r := require.New(t)

	for _, test_priv := range []string{"USAGE", "READ"} {
		in := map[string]interface{}{
			"stage_name":        "test-stage",
			"schema_name":       "test-schema",
			"database_name":     "test-db",
			"privilege":         test_priv,
			"roles":             []interface{}{"test-role-1", "test-role-2"},
			"shares":            []interface{}{"test-share-1", "test-share-2"},
			"with_grant_option": true,
		}
		d := schema.TestResourceDataRaw(t, resources.StageGrant().Resource.Schema, in)
		r.NotNil(d)

		WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
			mock.ExpectExec(fmt.Sprintf(`^GRANT %s ON STAGE "test-db"."test-schema"."test-stage" TO ROLE "test-role-1" WITH GRANT OPTION$`, test_priv)).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(fmt.Sprintf(`^GRANT %s ON STAGE "test-db"."test-schema"."test-stage" TO ROLE "test-role-2" WITH GRANT OPTION$`, test_priv)).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(fmt.Sprintf(`^GRANT %s ON STAGE "test-db"."test-schema"."test-stage" TO SHARE "test-share-1" WITH GRANT OPTION$`, test_priv)).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(fmt.Sprintf(`^GRANT %s ON STAGE "test-db"."test-schema"."test-stage" TO SHARE "test-share-2" WITH GRANT OPTION$`, test_priv)).WillReturnResult(sqlmock.NewResult(1, 1))
			expectReadStageGrant(mock, test_priv)
			err := resources.CreateStageGrant(d, db)
			r.NoError(err)
		})
	}
}

func TestStageGrantRead(t *testing.T) {
	r := require.New(t)

	d := stageGrant(t, "test-db|test-schema|test-stage|USAGE|false", map[string]interface{}{
		"stage_name":        "test-stage",
		"schema_name":       "test-schema",
		"database_name":     "test-db",
		"privilege":         "USAGE",
		"roles":             []interface{}{},
		"shares":            []interface{}{},
		"with_grant_option": false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadStageGrant(mock, "USAGE")
		err := resources.ReadStageGrant(d, db)
		r.NoError(err)
	})

	roles := d.Get("roles").(*schema.Set)
	r.True(roles.Contains("test-role-1"))
	r.True(roles.Contains("test-role-2"))
	r.Equal(roles.Len(), 2)

	shares := d.Get("shares").(*schema.Set)
	r.True(shares.Contains("test-share-1"))
	r.True(shares.Contains("test-share-2"))
	r.Equal(shares.Len(), 2)
}

func expectReadStageGrant(mock sqlmock.Sqlmock, test_priv string) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "STAGE", "test-stage", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "STAGE", "test-stage", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "STAGE", "test-stage", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "STAGE", "test-stage", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON STAGE "test-db"."test-schema"."test-stage"$`).WillReturnRows(rows)
}
