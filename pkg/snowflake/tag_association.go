package snowflake

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// TagAssociationBuilder abstracts the creation of SQL queries for a Snowflake tag.
type TagAssociationBuilder struct {
	databaseName string
	objectName   string
	objectType   string
	schemaName   string
	tagName      string
	tagValue     string
}

type tagAssociation struct {
	ColumnID       sql.NullString `db:"COLUMN_ID"`
	ColumnName     sql.NullString `db:"COLUMN_NAME"`
	Domain         sql.NullString `db:"DOMAIN"`
	ObjectDatabase sql.NullString `db:"OBJECT_DATABASE"`
	ObjectDeleted  sql.NullString `db:"OBJECT_DELETED"`
	ObjectID       sql.NullString `db:"OBJECT_ID"`
	ObjectName     sql.NullString `db:"OBJECT_NAME"`
	ObjectSchema   sql.NullString `db:"OBJECT_SCHEMA"`
	TagDatabase    sql.NullString `db:"TAG_DATABASE"`
	TagID          sql.NullString `db:"TAG_ID"`
	TagName        sql.NullString `db:"TAG_NAME"`
	TagSchema      sql.NullString `db:"TAG_SCHEMA"`
	TagValue       sql.NullString `db:"TAG_VALUE"`
}

// WithObjectId adds the name of the schema to the TagAssociationBuilder.
func (tb *TagAssociationBuilder) WithObjectName(objectName string) *TagAssociationBuilder {
	tb.objectName = objectName
	return tb
}

// WithObjectType adds the object type of the resource to add tag attachement to the TagAssociationBuilder.
func (tb *TagAssociationBuilder) WithObjectType(objectType string) *TagAssociationBuilder {
	tb.objectType = objectType
	return tb
}

// WithTagValue adds the name of the tag value to the TagAssociationBuilder.
func (tb *TagAssociationBuilder) WithTagValue(tagValue string) *TagAssociationBuilder {
	tb.tagValue = tagValue
	return tb
}

// GetTagDatabase returns the value of the tag database of TagAssociationBuilder.
func (tb *TagAssociationBuilder) GetTagDatabase() string {
	return tb.databaseName
}

// GetTagName returns the value of the tag name of TagAssociationBuilder.
func (tb *TagAssociationBuilder) GetTagName() string {
	return tb.schemaName
}

// GetTagSchema returns the value of the tag schema of TagAssociationBuilder.
func (tb *TagAssociationBuilder) GetTagSchema() string {
	return tb.schemaName
}

// TagAssociation returns a pointer to a Builder that abstracts the DDL operations for a tag sssociation.
//
// Supported DDL operations are:
//   - CREATE TAG
//   - DROP TAG
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/object-tagging.html)
func TagAssociation(tagID string) *TagAssociationBuilder {
	databaseName, schemaName, tagName := validation.ParseFullyQualifiedObjectID(tagID)
	return &TagAssociationBuilder{
		databaseName: databaseName,
		schemaName:   schemaName,
		tagName:      tagName,
	}
}

// Create returns the SQL query that will set the tag on an object.
func (tb *TagAssociationBuilder) Create() string {
	return fmt.Sprintf(`ALTER %v "%v" SET TAG "%v"."%v"."%v" = '%v'`, tb.objectType, tb.objectName, tb.databaseName, tb.schemaName, tb.tagName, tb.tagValue)
}

// Drop returns the SQL query that will remove a tag from an object.
func (tb *TagAssociationBuilder) Drop() string {
	return fmt.Sprintf(`ALTER %v "%v" UNSET TAG "%v"."%v"."%v"`, tb.objectType, tb.objectName, tb.databaseName, tb.schemaName, tb.tagName)
}

func ListTagAssociations(tb *TagAssociationBuilder, db *sql.DB) ([]tagAssociation, error) {
	stmt := fmt.Sprintf(`SELECT * FROM SNOWFLAKE.ACCOUNT_USAGE.TAG_REFERENCES WHERE TAG_NAME = '%v' AND DOMAIN = '%v' AND TAG_VALUE = '%v'`, tb.tagName, tb.objectType, tb.tagValue)
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tagAssociations := []tagAssociation{}
	err = sqlx.StructScan(rows, &tagAssociations)
	log.Printf("[DEBUG] tagAssociations is %v", tagAssociations)

	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no tag associations found for tag %s", tb.tagName)
		return nil, err
	}
	return tagAssociations, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
