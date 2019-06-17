package snowflake

// ManagedAccount returns a pointer to a Builder that abstracts the DDL
// operations for a reader account.
//
// Supported DDL operations are:
//   - CREATE MANAGED ACCOUNT
//   - DROP MANAGED ACCOUNT
//   - SHOW MANAGED ACCOUNTS
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/user-guide/data-sharing-reader-create.html)
func ManagedAccount(name string) *Builder {
	return &Builder{
		entityType: ManagedAccountType,
		name:       name,
	}
}
