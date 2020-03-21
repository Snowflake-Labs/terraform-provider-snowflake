package snowflake

// StorageIntegration returns a pointer to a Builder that abstracts the DDL operations for a storage integration.
//
// Supported DDL operations are:
//   - CREATE STORAGE INTEGRATION
//   - ALTER STORAGE INTEGRATION
//   - DROP INTEGRATION
//   - SHOW INTEGRATIONS
//   - DESCRIBE INTEGRATION
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/sql-reference/ddl-user-security.html#storage-integrations)
func StorageIntegration(name string) *Builder {
	return &Builder{
		entityType: StorageIntegrationType,
		name:       name,
	}
}
