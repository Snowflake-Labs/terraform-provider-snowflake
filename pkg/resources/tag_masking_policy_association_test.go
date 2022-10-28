package resources_test

import (
	"database/sql"
	"testing"

	"regexp"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestTagMaskingPolicyAssociation(t *testing.T) {
	r := require.New(t)
	err := resources.Tag().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestTagMaskingPolicyAssociationCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"tag_id":            "tag_db|tag_schema|tag_name",
		"masking_policy_id": "mp_db|mp_schema|mp_name",
	}
	d := schema.TestResourceDataRaw(t, resources.TagMaskingPolicyAssociation().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^ALTER TAG "tag_db"."tag_schema"."tag_name" SET MASKING POLICY "mp_db"."mp_schema"."mp_name"$`).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadTestTagMaskingPolicyAssociation(mock)
		err := resources.CreateTagMaskingPolicyAssociation(d, db)
		r.NoError(err)
	})
}

func TestTagMaskingPolicyAssociationDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"tag_id":            "tag_db|tag_schema|tag_name",
		"masking_policy_id": "mp_db|mp_schema|mp_name",
	}

	d := schema.TestResourceDataRaw(t, resources.TagMaskingPolicyAssociation().Schema, in)
	d.SetId("tag_db|tag_schema|tag_name|mp_db|mp_schema|mp_name")
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^ALTER TAG "tag_db"."tag_schema"."tag_name" UNSET MASKING POLICY "mp_db"."mp_schema"."mp_name"$`).WillReturnResult(sqlmock.NewResult(1, 1))

		err := resources.DeleteTagMaskingPolicyAssociation(d, db)

		r.NoError(err)
	})
}

func TestTagMaskingPolicyAssociationRead(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"tag_id":            "tag_db|tag_schema|tag_name",
		"masking_policy_id": "mp_db|mp_schema|mp_name",
	}

	d := schema.TestResourceDataRaw(t, resources.TagMaskingPolicyAssociation().Schema, in)
	d.SetId("tag_db|tag_schema|tag_name|mp_db|mp_schema|mp_name")

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		mP := snowflake.MaskingPolicy("mp_name", "mp_db", "mp_schema")
		q := snowflake.Tag("tag_name").WithDB("tag_db").WithSchema("tag_schema").WithMaskingPolicy(mP).ShowAttachedPolicy()
		mock.ExpectQuery(regexp.QuoteMeta(q)).WillReturnError(sql.ErrNoRows)
		err := resources.ReadTagMaskingPolicyAssociation(d, db)

		r.Empty(d.State())
		r.Nil(err)
	})
}

func expectReadTestTagMaskingPolicyAssociation(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"POLICY_DB", "POLICY_SCHEMA", "POLICY_NAME", "POLICY_KIND", "REF_DATABASE_NAME", "REF_SCHEMA_NAME", "REF_ENTITY_NAME", "REF_ENTITY_DOMAIN"},
	).AddRow("mp_db", "mp_schema", "mp_name", "MASKING", "tag_db", "tag_schema", "tag_name", "TAG")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * from table ("tag_db".information_schema.policy_references(ref_entity_name => '"tag_db"."tag_schema"."tag_name"', ref_entity_domain => 'TAG')) where policy_db='mp_db' and policy_schema='mp_schema' and policy_name='mp_name'`)).WillReturnRows(rows)
}
