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

type tagAttachment struct {
	ColumnId       sql.NullString `db:"COLUMN_ID"`
	ColumnName     sql.NullString `db:"COLUMN_NAME"`
	Domain         sql.NullString `db:"DOMAIN"`
	ObjectDatabase sql.NullString `db:"OBJECT_DATABASE"`
	ObjectDeleted  sql.NullString `db:"OBJECT_DELETED"`
	ObjectId       sql.NullString `db:"OBJECT_ID"`
	ObjectName     sql.NullString `db:"OBJECT_NAME"`
	ObjectSchema   sql.NullString `db:"OBJECT_SCHEMA"`
	TagDatabase    sql.NullString `db:"TAG_DATABASE"`
	TagId          sql.NullString `db:"TAG_ID"`
	TagName        sql.NullString `db:"TAG_NAME"`
	TagSchema      sql.NullString `db:"TAG_SCHEMA"`
	TagValue       sql.NullString `db:"TAG_VALUE"`
}

// WithResourceId adds the name of the schema to the TagAttachmentBuilder
func (tb *TagAttachmentBuilder) WithResourceId(resourceId string) *TagAttachmentBuilder {
	tb.resourceId = resourceId
	return tb
}

// WithObjectType adds the object type of the resource to add tag attachement to the TagAttachmentBuilder
func (tb *TagAttachmentBuilder) WithObjectType(objectType string) *TagAttachmentBuilder {
	tb.objectType = objectType
	return tb
}

// WithTagValue adds the name of the tag value to the TagAttachmentBuilder
func (tb *TagAttachmentBuilder) WithTagValue(tagValue string) *TagAttachmentBuilder {
	tb.tagValue = tagValue
	return tb
}

// GetTagDatabase returns the value of the tag database of TagAttachmentBuilder
func (tb *TagAttachmentBuilder) GetTagDatabase() string {
	return tb.databaseName
}

// GetTagName returns the value of the tag name of TagAttachmentBuilder
func (tb *TagAttachmentBuilder) GetTagName() string {
	return tb.schemaName
}

// GetTagSchema returns the value of the tag schema of TagAttachmentBuilder
func (tb *TagAttachmentBuilder) GetTagSchema() string {
	return tb.schemaName
}

// TagAttachment returns a pointer to a Builder that abstracts the DDL operations for a tag attachment.
//
// Supported DDL operations are:
//   - CREATE TAG
//   - DROP TAG
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/object-tagging.html)
func TagAttachment(tagId string) *TagAttachmentBuilder {
	parsedString := strings.Replace(tagId, "\"", "", -1)

	var s []string
	if strings.Contains(parsedString, "|") {
		s = strings.Split(strings.ToUpper(parsedString), "|")
	} else {
		s = strings.Split(strings.ToUpper(parsedString), ".")
	}
	return &TagAttachmentBuilder{
		databaseName: s[0],
		schemaName:   s[1],
		tagName:      s[2],
	}
}

// Create returns the SQL query that will set the tag on a resource.
func (tb *TagAttachmentBuilder) Create() string {
	//return fmt.Sprintf(`USE WAREHOUSE %v ALTER %v %v SET TAG %v = "%v"`, tb.warehouse, tb.objectType, tb.resourceId, tb.tagName, tb.tagValue)
	return fmt.Sprintf(`ALTER %v %v SET TAG "%v"."%v"."%v" = '%v'`, tb.objectType, tb.resourceId, tb.databaseName, tb.schemaName, tb.tagName, tb.tagValue)
}

// Drop returns the SQL query that will remove a tag from a resource.
func (tb *TagAttachmentBuilder) Drop() string {
	return fmt.Sprintf(`ALTER %v %v UNSET TAG "%v"."%v"."%v"`, tb.objectType, tb.resourceId, tb.databaseName, tb.schemaName, tb.tagName)
}

func ListTagAttachments(tb *TagAttachmentBuilder, db *sql.DB) ([]tagAttachment, error) {
	stmt := fmt.Sprintf(`SELECT * FROM SNOWFLAKE.ACCOUNT_USAGE.TAG_REFERENCES WHERE TAG_NAME = '%v' AND DOMAIN = '%v' AND TAG_VALUE = '%v'`, tb.tagName, tb.objectType, tb.tagValue)
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tagAttachments := []tagAttachment{}
	err = sqlx.StructScan(rows, &tagAttachments)
	log.Printf("[DEBUG] tagAttachments is %v", tagAttachments)

	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no tag attachments found")
		return nil, err
	}
	return tagAttachments, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
