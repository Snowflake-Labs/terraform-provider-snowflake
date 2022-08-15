package snowflake

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// TagAttachmentBuilder abstracts the creation of SQL queries for a Snowflake tag
type TagAttachmentBuilder struct {
	resourceId   string
	objectType   string
	databaseName string
	schemaName   string
	tagName      string
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
func TagAttachment(tag string) *TagAttachmentBuilder {
	s := strings.Split(tag, ".")
	return &TagAttachmentBuilder{
		databaseName: s[0],
		schemaName:   s[1],
		tagName:      s[2],
	}
}

// Create returns the SQL query that will set the tag on a resource.
func (tb *TagAttachmentBuilder) Create() string {
	return fmt.Sprintf(`ALTER %v SET TAG %v ='%v'`, tb.objectType, tb.tagName, tagValue)
}

// Drop returns the SQL query that will show a tag.
func (tb *TagAttachmentBuilder) Show() string {
	return fmt.Sprintf(`SHOW TAGS LIKE '%v' IN SCHEMA "%v"."%v"`, tb.tagName, tb.databaseName, tb.schemaName)
}

// Show returns the SQL query that will show a pipe.
func (pb *PipeBuilder) Show() string {
	return fmt.Sprintf(`SHOW PIPES LIKE '%v' IN SCHEMA "%v"."%v"`, pb.name, pb.db, pb.schema)
}

// Drop returns the SQL query that will remove a tag from a resource.
func (tb *TagAttachmentBuilder) Drop() string {
	return fmt.Sprintf(`ALTER %v UNSET TAG %v ='%v'`, tb.objectType, tb.tagName, tagValue)
}

// ScanTagAttachment returns a list of tags for a resource
func ScanTagAttachment(tb *TagAttachmentBuilder) string {
	p := &pipe{}
	e := row.StructScan(p)
	return p, e
}

func ScanPipe(row *sqlx.Row) (*pipe, error) {
	p := &pipe{}
	e := row.StructScan(p)
	return p, e
}
