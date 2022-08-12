package snowflake

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// TagAttachmentBuilder abstracts the creation of SQL queries for a Snowflake tag
type TagAttachmentBuilder struct {
	resourceId string
	objectType string
	comment    string
	tag        interface{}
}

// WithResourceId adds the name of the schema to the TagAttachmentBuilder
func (tb *TagAttachmentBuilder) WithResourceId(resourceId string) *TagAttachmentBuilder {
	tb.resourceId = resourceId
	return tb
}

// WithComment adds a comment to the TagAttachmentBuilder
func (tb *TagAttachmentBuilder) WithObjectType(object_type string) *TagAttachmentBuilder {
	tb.objectType = object_type
	return tb
}

// TagAttachment returns a pointer to a Builder that abstracts the DDL operations for a tag.
//
// Supported DDL operations are:
//   - CREATE TAG
//   - ALTER TAG
//   - DROP TAG
//   - UNDROP TAG
//   - SHOW TAGS
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/object-tagging.html)
func TagAttachment(tag interface{}) *TagAttachmentBuilder {
	return &TagAttachmentBuilder{
		tag: tag,
	}
}

// Create returns the SQL query that will set the tag on a resource.
func (tb *TagAttachmentBuilder) Create() string {
	tagKey := tb.tag.([]interface{})[0].(map[string]interface{})["name"]
	tagValue := tb.tag.([]interface{})[0].(map[string]interface{})["value"]
	return fmt.Sprintf(`ALTER %v SET TAG %v ='%v'`, tb.objectType, tagKey, tagValue)
}

// Drop returns the SQL query that will remove a tag from a resource.
func (tb *TagAttachmentBuilder) Drop() string {
	tagKey := tb.tag.([]interface{})[0].(map[string]interface{})["name"]
	tagValue := tb.tag.([]interface{})[0].(map[string]interface{})["value"]
	return fmt.Sprintf(`ALTER %v UNSET TAG %v ='%v'`, tb.objectType, tagKey, tagValue)
}

// ListTagsForResource returns a list of tags for a resource
func ListTagsForResource(databaseName, schemaName string, db *sql.DB) ([]tag, error) {
	stmt := fmt.Sprintf(`SHOW TAGS IN %v %v`, object_type, resource_id)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []tag{}
	err = sqlx.StructScan(rows, &tags)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no tags found")
		return nil, nil
	}
	return tags, errors.Wrapf(err, "unable to scan row for %s", stmt)
}

//// RemoveTagFromResource returns a list of tags for a resource
//func RemoveTagFromResource(databaseName, schemaName string, db *sql.DB) ([]tag, error) {
//	stmt := fmt.Sprintf(`SHOW TAGS IN %v %v`, object_type, resource_id)
//	rows, err := Query(db, stmt)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	tags := []tag{}
//	err = sqlx.StructScan(rows, &tags)
//	if err == sql.ErrNoRows {
//		log.Printf("[DEBUG] no tags found")
//		return nil, nil
//	}
//	return tags, errors.Wrapf(err, "unable to scan row for %s", stmt)
//}
