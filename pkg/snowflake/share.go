package snowflake

import "strings"

// Share returns a pointer to a Builder that abstracts the DDL operations for a share.
//
// Supported DDL operations are:
//   - CREATE SHARE
//   - ALTER SHARE
//   - DROP SHARE
//   - SHOW SHARES
//   - DESCRIBE SHARE
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/sql-reference/ddl-database.html#share-management)
func Share(name string) *Builder {
	return &Builder{
		entityType: ShareType,
		name:       strings.ToUpper(name),
	}
}
