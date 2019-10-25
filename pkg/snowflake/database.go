package snowflake

import (
	"fmt"
)

// Database returns a pointer to a Builder for a database
func Database(name string) *Builder {
	return &Builder{
		name:       name,
		entityType: DatabaseType,
	}
}

// DatabaseShareBuilder is a basic builder that just creates databases from shares
type DatabaseShareBuilder struct {
	name     string
	provider string
	share    string
}

// DatabaseFromShare returns a pointer to a builder that can create a database from a share
func DatabaseFromShare(name, provider, share string) *DatabaseShareBuilder {
	return &DatabaseShareBuilder{
		name:     name,
		provider: provider,
		share:    share,
	}
}

// Create returns the SQL statement required to create a database from a share
func (dsb *DatabaseShareBuilder) Create() string {
	return fmt.Sprintf(`CREATE DATABASE "%v" FROM SHARE "%v"."%v"`, dsb.name, dsb.provider, dsb.share)
}

// DatabaseCloneBuilder is a basic builder that just creates databases from a source database
type DatabaseCloneBuilder struct {
	name     string
	database string
}

// DatabaseFromDatabase returns a pointer to a builder that can create a database from a source database
func DatabaseFromDatabase(name, database string) *DatabaseCloneBuilder {
	return &DatabaseCloneBuilder{
		name:     name,
		database: database,
	}
}

// Create returns the SQL statement required to create a database from a source database
func (dsb *DatabaseCloneBuilder) Create() string {
	return fmt.Sprintf(`CREATE DATABASE "%v" CLONE "%v"`, dsb.name, dsb.database)
}
