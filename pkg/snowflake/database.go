package snowflake

import (
	"database/sql"
	"fmt"
	"log"

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

// DatabaseReplicaBuilder is a basic builder that just creates databases from an avilable replication source
type DatabaseReplicaBuilder struct {
	name    string
	replica string
}

// DatabaseFromReplica returns a pointer to a builder that can create a database from an avilable replication source
func DatabaseFromReplica(name, replica string) *DatabaseReplicaBuilder {
	return &DatabaseReplicaBuilder{
		name:    name,
		replica: replica,
	}
}

// Create returns the SQL statement required to create a database from an avilable replication source
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
