package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
)

func TestTagAttachment(t *testing.T) {
	r := require.New(t)
	err := resources.Tag().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestTagAttachmentCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"resource_id": "good_name",
		"object_type": "USER",
		"tag_id":      "testDb.testSchema.testTagName",
		"tag_value":   "test_tag_value",
	}
	d := schema.TestResourceDataRaw(t, resources.TagAttachment().Schema, in)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^ALTER USER good_name SET TAG "TESTDB"."TESTSCHEMA"."TESTTAGNAME" = 'test_tag_value'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		rows := sqlmock.NewRows([]string{
			"tag_database", "tag_id", "tag_name", "tag_schema", "tag_value", "object_database", "object_deleted", "object_id", "object_name", "object_schema",
			"domain", "column_id", "column_name",
		})
		mock.ExpectQuery(
			`^[SELECT * FROM SNOWFLAKE.ACCOUNT_USAGE.TAG_REFERENCES WHERE TAG_NAME = 'TESTTAGNAME' AND DOMAIN = 'USER' AND TAG_VALUE = 'test_tag_value'$`,
		).WillReturnRows(rows)
		//expectReadTagAttachment(mock)
		err := resources.CreateTagAttachment(d, db)
		r.NoError(err)
	})
}

//func TestTagAttachmentDelete(t *testing.T) {
//	r := require.New(t)
//
//	in := map[string]interface{}{
//		"resource_id": "good_name",
//		"object_type": "USER",
//		"tag_id":      "testDb.testSchema.testTagName",
//		"tag_value":   "test_tag_value",
//	}
//
//	//d := tagAttachment(t, "test_tag_attachment", in)
//	//r.NotNil(d)
//	d := schema.TestResourceDataRaw(t, resources.TagAttachment().Schema, in)
//	r.NotNil(d)
//
//	testhelpers.WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
//		mock.ExpectExec(
//			`^ALTER USER good_name UNSET TAG "TESTDB"."TESTSCHEMA"."TESTTAGNAME" = 'test_tag_value'$`,
//		).WillReturnResult(sqlmock.NewResult(1, 1))
//
//		err := resources.DeleteTagAttachment(d, db)
//		r.NoError(err)
//	})
//}

//func TestTagAttachmentRead(t *testing.T) {
//	r := require.New(t)
//
//	in := map[string]interface{}{
//		"resource_id": "good_name",
//		"object_type": "USER",
//		"tag_id":      "testDb.testSchema.testTagName",
//		"tag_value":   "test_tag_value",
//	}
//
//	d := schema.TestResourceDataRaw(t, resources.TagAttachment().Schema, in)
//	d.SetId("test_tag_attachment")
//
//	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
//		// Test when resource is not found, checking if state will be empty
//		r.NotEmpty(d.State())
//		q := snowflake.Tag("good_name").WithDB("test_db").WithSchema("test_schema").Show()
//		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
//		err := resources.ReadTag(d, db)
//		r.Empty(d.State())
//		r.Nil(err)
//	})
//}

func expectReadTagAttachment(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"tag_database", "tag_id", "tag_name", "tag_schema", "tag_value", "object_database", "object_deleted", "object_id", "object_name", "object_schema",
		"domain", "column_id", "column_name",
	})
	mock.ExpectQuery(`^SHOW TAGS LIKE 'good_name' IN SCHEMA "test_db"."test_schema"$`).WillReturnRows(rows)
}
