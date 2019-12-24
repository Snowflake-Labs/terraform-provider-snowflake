package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStageGrant(t *testing.T) {
	r := require.New(t)
	err := resources.StageGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestStageGrantCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"stage_name":    "test-stage",
		"schema_name":   "PUBLIC",
		"database_name": "test-db",
		"privilege":     "USAGE",
		"roles":         []interface{}{"test-role-1", "test-role-2"},
		"shares":        []interface{}{"test-share-1", "test-share-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.StageGrant().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT USAGE ON STAGE "test-db"."PUBLIC"."test-stage" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON STAGE "test-db"."PUBLIC"."test-stage" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON STAGE "test-db"."PUBLIC"."test-stage" TO SHARE "test-share-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON STAGE "test-db"."PUBLIC"."test-stage" TO SHARE "test-share-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadStageGrant(mock)
		err := resources.CreateStageGrant(d, db)
		a.NoError(err)
	})
}

func expectReadStageGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "STAGE", "test-stage", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "STAGE", "test-stage", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "STAGE", "test-stage", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "STAGE", "test-stage", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON STAGE "test-db"."PUBLIC"."test-stage"$`).WillReturnRows(rows)
}
