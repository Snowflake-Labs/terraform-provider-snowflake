package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// TagAttachmentBuilder abstracts the creation of SQL queries for a Snowflake tag
type TagAttachmentBuilder struct {
	databaseName string
	resourceId   string
	objectType   string
	schemaName   string
	tagName      string
	tagValue     string
}

// WithResourceId adds the name of the schema to the TagAttachmentBuilder
func (tb *TagAttachmentBuilder) WithResourceId(resourceId string) *TagAttachmentBuilder {
	tb.resourceId = resourceId
	return tb
}

// WithComment adds a comment to the TagAttachmentBuilder
func (tb *TagAttachmentBuilder) WithObjectType(objectType string) *TagAttachmentBuilder {
	tb.objectType = objectType
	return tb
}

// WithTagValue adds the name of the tag value to the TagAttachmentBuilder
func (tb *TagAttachmentBuilder) WithTagValue(tagValue string) *TagAttachmentBuilder {
	tb.tagValue = tagValue
	return tb
}

// TagAttachment returns a pointer to a Builder that abstracts the DDL operations for a tag attachment.
//
// Supported DDL operations are:
//   - CREATE TAG
//   - DROP TAG
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/object-tagging.html)
func TagAttachment(tag string) *TagAttachmentBuilder {
	s := strings.Split(strings.ToUpper(tag), ".")
	return &TagAttachmentBuilder{
		databaseName: s[0],
		schemaName:   s[1],
		tagName:      s[2],
	}
}

// Create returns the SQL query that will set the tag on a resource.
func (tb *TagAttachmentBuilder) Create() string {
	return fmt.Sprintf(`ALTER %v "%v" SET TAG %v ="%v"`, tb.objectType, tb.resourceId, tb.tagName, tb.tagValue)
}

// Drop returns the SQL query that will remove a tag from a resource.
func (tb *TagAttachmentBuilder) Drop() string {
	return fmt.Sprintf(`ALTER %v UNSET TAG %v`, tb.objectType, tb.tagName)
}

// ListTagAttachments returns a list of tags in a database or schema
func ListTagAttachments(tb *TagAttachmentBuilder, db *sql.DB) ([]tag, error) {
	var stmt strings.Builder
	stmt.WriteString(fmt.Sprintf(`SELECT * FROM SNOWFLAKE.ACCOUNT_USAGE.TAG_REFERENCES`))

	// add filters for the tag
	stmt.WriteString(fmt.Sprintf(`WHERE TAG_NAME = '%v' AND DOMAIN = '%v' AND TAG_VALUE = '%v'`, tb.tagName, tb.objectType, tb.tagValue))
	rows, err := Query(db, stmt.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []tag
	err = sqlx.StructScan(rows, &tags)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no tags found")
		return nil, nil
	}
	return tags, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
