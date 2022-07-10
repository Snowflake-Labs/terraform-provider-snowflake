package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	name             string
	provider         string //Deprecated: use organization_name and account_name instead
	share            string
	organizationName string
	accountName      string
	comment          string
}

// DatabaseFromShare returns a pointer to a builder that can create a database from a share
func DatabaseFromShare(name, share string) *DatabaseShareBuilder {
	return &DatabaseShareBuilder{
		name:  name,
		share: share,
	}
}

// WithProvider adds a provider to the DatabaseShareBuilder
// Deprecated: use WithOrg instead
func (dsb *DatabaseShareBuilder) WithProvider(provider string) *DatabaseShareBuilder {
	dsb.provider = provider
	return dsb
}

// WithOrg adds a organizationName and accountName to the DatabaseShareBuilder
func (dsb *DatabaseShareBuilder) WithOrg(organizationName, accountName string) *DatabaseShareBuilder {
	dsb.organizationName = organizationName
	dsb.accountName = accountName
	return dsb
}

// WithComment adds a comment to the DatabaseShareBuilder
func (dsb *DatabaseShareBuilder) WithComment(comment string) *DatabaseShareBuilder {
	dsb.comment = comment
	return dsb
}

// Create returns the SQL statement required to create a database from a share
func (dsb *DatabaseShareBuilder) Create() string {
	var stmt strings.Builder
	stmt.WriteString(fmt.Sprintf(`CREATE DATABASE "%v" FROM SHARE`, dsb.name))
	if dsb.provider != "" {
		stmt.WriteString(fmt.Sprintf(` "%v"."%v"`, dsb.provider, dsb.share))
	} else {
		stmt.WriteString(fmt.Sprintf(` "%v"."%v"."%v"`, dsb.organizationName, dsb.accountName, dsb.share))
	}

	if dsb.comment != "" {
		stmt.WriteString(fmt.Sprintf(` COMMENT = '%v'`, dsb.comment))
	}

	return stmt.String()
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

// DatabaseReplicaBuilder is a basic builder that just creates databases from an available replication source
type DatabaseReplicaBuilder struct {
	name    string
	replica string
}

// DatabaseFromReplica returns a pointer to a builder that can create a database from an available replication source
func DatabaseFromReplica(name, replica string) *DatabaseReplicaBuilder {
	return &DatabaseReplicaBuilder{
		name:    name,
		replica: replica,
	}
}

// Create returns the SQL statement required to create a database from an available replication source
func (dsb *DatabaseReplicaBuilder) Create() string {
	return fmt.Sprintf(`CREATE DATABASE "%v" AS REPLICA OF "%v"`, dsb.name, dsb.replica)
}

type database struct {
	CreatedOn     sql.NullString `db:"created_on"`
	DBName        sql.NullString `db:"name"`
	IsDefault     sql.NullString `db:"is_default"`
	IsCurrent     sql.NullString `db:"is_current"`
	Origin        sql.NullString `db:"origin"`
	Owner         sql.NullString `db:"owner"`
	Comment       sql.NullString `db:"comment"`
	Options       sql.NullString `db:"options"`
	RetentionTime sql.NullString `db:"retention_time"`
}

func ScanDatabase(row *sqlx.Row) (*database, error) {
	d := &database{}
	e := row.StructScan(d)
	return d, e
}

func ListDatabases(sdb *sqlx.DB) ([]database, error) {
	stmt := "SHOW DATABASES"
	rows, err := sdb.Queryx(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []database{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no databases found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}

func ListDatabase(sdb *sqlx.DB, databaseName string) (*database, error) {
	stmt := fmt.Sprintf("SHOW DATABASES LIKE '%s'", databaseName)
	rows, err := sdb.Queryx(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []database{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows || len(dbs) == 0 {
		log.Printf("[DEBUG] no databases found")
		return nil, nil
	}
	db := &database{}
	for _, d := range dbs {
		if d.DBName.String == databaseName {
			db = &d
			break
		}
	}
	return db, errors.Wrapf(err, "unable to scan row for %s", stmt)
}

// EnableReplicationAccounts returns the SQL query that will enable replication to provided accounts
func (db *Builder) EnableReplicationAccounts(dbName string, accounts string) string {
	return fmt.Sprintf(`ALTER DATABASE "%v" ENABLE REPLICATION TO ACCOUNTS %v`, dbName, accounts)
}

// DisableReplicationAccounts returns the SQL query that will disable replication to provided accounts
func (db *Builder) DisableReplicationAccounts(dbName string, accounts string) string {
	return fmt.Sprintf(`ALTER DATABASE "%v" DISABLE REPLICATION TO ACCOUNTS %v`, dbName, accounts)
}

// GetRemovedAccountsFromReplicationConfiguration compares two old and new configurations and returns any values that
// were deleted from the old configuration.
func (db *Builder) GetRemovedAccountsFromReplicationConfiguration(oldAcc []interface{}, newAcc []interface{}) []interface{} {
	accountMap := make(map[string]bool)
	var removedAccounts []interface{}
	// insert all values from new configuration into mapping
	for _, v := range newAcc {
		accountMap[v.(string)] = true
	}
	for _, v := range oldAcc {
		if !accountMap[v.(string)] {
			removedAccounts = append(removedAccounts, v.(string))
		}
	}
	return removedAccounts
}
