package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// DatabaseBuilder abstracts the creation of SQL queries for a Snowflake database.
type DatabaseBuilder struct {
	name                 string
	comment              string
	transient            bool
	cloneDatabase        string
	setDataRetentionDays bool
	dataRetentionDays    int
	tags                 []TagValue
}

func (db *DatabaseBuilder) QualifiedName() string {
	return fmt.Sprintf(`"%v"`, db.name)
}

// Clone adds CLONE to the DatabaseBuilder to create a clone of another database.
func (db *DatabaseBuilder) Clone(database string) *DatabaseBuilder {
	db.cloneDatabase = database
	return db
}

// Transient adds the TRANSIENT flag to the DatabaseBuilder.
func (db *DatabaseBuilder) Transient() *DatabaseBuilder {
	db.transient = true
	return db
}

// WithComment adds a comment to the DatabaseBuilder.
func (db *DatabaseBuilder) WithComment(c string) *DatabaseBuilder {
	db.comment = c
	return db
}

// WithDataRetentionDays adds the days to retain data to the DatabaseBuilder (must
// be 0-1 for standard edition, 0-90 for enterprise edition).
func (db *DatabaseBuilder) WithDataRetentionDays(d int) *DatabaseBuilder {
	db.setDataRetentionDays = true
	db.dataRetentionDays = d
	return db
}

// WithTags sets the tags on the DatabaseBuilder.
func (db *DatabaseBuilder) WithTags(tags []TagValue) *DatabaseBuilder {
	db.tags = tags
	return db
}

// AddTag returns the SQL query that will add a new tag to the database.
func (db *DatabaseBuilder) AddTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER DATABASE %s SET TAG "%v"."%v"."%v" = "%v"`, db.QualifiedName(), tag.Database, tag.Schema, tag.Name, tag.Value)
}

// ChangeTag returns the SQL query that will alter a tag on the database.
func (db *DatabaseBuilder) ChangeTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER DATABASE %s SET TAG "%v"."%v"."%v" = "%v"`, db.QualifiedName(), tag.Database, tag.Schema, tag.Name, tag.Value)
}

// UnsetTag returns the SQL query that will unset a tag on the database.
func (db *DatabaseBuilder) UnsetTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER DATABASE %s UNSET TAG "%v"."%v"."%v"`, db.QualifiedName(), tag.Database, tag.Schema, tag.Name)
}

// Database returns a pointer to a Builder that abstracts the DDL operations for a database.
//
// Supported DDL operations are:
//   - CREATE DATABASE
//   - ALTER DATABASE
//   - DROP DATABASE
//   - UNDROP DATABASE
//   - USE DATABASE
//   - SHOW DATABASE
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/sql-reference/ddl-database.html#database-management)
func NewDatabaseBuilder(name string) *DatabaseBuilder {
	return &DatabaseBuilder{
		name: name,
	}
}

// Create returns the SQL query that will create a new database.
func (db *DatabaseBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	if db.transient {
		q.WriteString(` TRANSIENT`)
	}

	q.WriteString(fmt.Sprintf(` DATABASE %v`, db.QualifiedName()))

	if db.cloneDatabase != "" {
		q.WriteString(fmt.Sprintf(` CLONE "%v"`, db.cloneDatabase))
	}

	if db.setDataRetentionDays {
		q.WriteString(fmt.Sprintf(` DATA_RETENTION_TIME_IN_DAYS = %d`, db.dataRetentionDays))
	}

	if db.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(db.comment)))
	}

	return q.String()
}

// Rename returns the SQL query that will rename the database.
func (db *DatabaseBuilder) Rename(newName string) string {
	oldName := db.QualifiedName()
	db.name = newName
	return fmt.Sprintf(`ALTER DATABASE %v RENAME TO %v`, oldName, db.QualifiedName())
}

// Swap returns the SQL query that Swaps all objects (tables, views, etc.) and
// metadata, including identifiers, between the two specified databases.
func (db *DatabaseBuilder) Swap(targetDatabase string) string {
	sourceDatabase := db.QualifiedName()
	db.name = targetDatabase
	return fmt.Sprintf(`ALTER DATABASE %v SWAP WITH %v`, sourceDatabase, db.QualifiedName())
}

// ChangeComment returns the SQL query that will update the comment on the database.
func (db *DatabaseBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER DATABASE %v SET COMMENT = '%v'`, db.QualifiedName(), EscapeString(c))
}

// RemoveComment returns the SQL query that will remove the comment on the database.
func (db *DatabaseBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER DATABASE %v UNSET COMMENT`, db.QualifiedName())
}

// ChangeDataRetentionDays returns the SQL query that will update the data retention days on the database.
func (db *DatabaseBuilder) ChangeDataRetentionDays(d int) string {
	return fmt.Sprintf(`ALTER DATABASE %v SET DATA_RETENTION_TIME_IN_DAYS = %d`, db.QualifiedName(), d)
}

// RemoveDataRetentionDays returns the SQL query that will remove the data retention days on the database.
func (db *DatabaseBuilder) RemoveDataRetentionDays() string {
	return fmt.Sprintf(`ALTER DATABASE %v UNSET DATA_RETENTION_TIME_IN_DAYS`, db.QualifiedName())
}

// Drop returns the SQL query that will drop a database.
func (db *DatabaseBuilder) Drop() string {
	return fmt.Sprintf(`DROP DATABASE %v`, db.QualifiedName())
}

// Undrop returns the SQL query that will undrop a database.
func (db *DatabaseBuilder) Undrop() string {
	return fmt.Sprintf(`UNDROP DATABASE %v`, db.QualifiedName())
}

// Use returns the SQL query that will use a database.
func (db *DatabaseBuilder) Use() string {
	return fmt.Sprintf(`USE DATABASE %v`, db.QualifiedName())
}

// Show returns the SQL query that will show a database.
func (db *DatabaseBuilder) Show() string {
	return fmt.Sprintf(`SHOW DATABASES LIKE '%v'`, db.name)
}

// EnableReplicationAccounts returns the SQL query that will enable replication to provided accounts.
func (db *DatabaseBuilder) EnableReplicationAccounts(dbName string, accounts string) string {
	return fmt.Sprintf(`ALTER DATABASE "%v" ENABLE REPLICATION TO ACCOUNTS %v`, dbName, accounts)
}

// DisableReplicationAccounts returns the SQL query that will disable replication to provided accounts.
func (db *DatabaseBuilder) DisableReplicationAccounts(dbName string, accounts string) string {
	return fmt.Sprintf(`ALTER DATABASE "%v" DISABLE REPLICATION TO ACCOUNTS %v`, dbName, accounts)
}

// DatabaseShareBuilder is a basic builder that just creates databases from shares.
type DatabaseShareBuilder struct {
	name     string
	provider string
	share    string
	comment  string
}

// DatabaseFromShare returns a pointer to a builder that can create a database from a share.
func DatabaseFromShare(name, provider, share string) *DatabaseShareBuilder {
	return &DatabaseShareBuilder{
		name:     name,
		provider: provider,
		share:    share,
	}
}

// WithComment adds a comment to the DatabaseShareBuilder.
func (dsb *DatabaseShareBuilder) WithComment(comment string) *DatabaseShareBuilder {
	dsb.comment = comment
	return dsb
}

// Create returns the SQL statement required to create a database from a share.
func (dsb *DatabaseShareBuilder) Create() string {
	var q strings.Builder
	q.WriteString(fmt.Sprintf(`CREATE DATABASE "%v" FROM SHARE "%v"."%v"`, dsb.name, dsb.provider, dsb.share))

	if dsb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, dsb.comment))
	}

	return q.String()
}

// DatabaseReplicaBuilder is a basic builder that just creates databases from an available replication source.
type DatabaseReplicaBuilder struct {
	name    string
	replica string
}

// DatabaseFromReplica returns a pointer to a builder that can create a database from an available replication source.
func DatabaseFromReplica(name, replica string) *DatabaseReplicaBuilder {
	return &DatabaseReplicaBuilder{
		name:    name,
		replica: replica,
	}
}

// Create returns the SQL statement required to create a database from an available replication source.
func (dsb *DatabaseReplicaBuilder) Create() string {
	return fmt.Sprintf(`CREATE DATABASE "%v" AS REPLICA OF "%v"`, dsb.name, dsb.replica)
}

// GetRemovedAccountsFromReplicationConfiguration compares two old and new configurations and returns any values that
// were deleted from the old configuration.
func (db *DatabaseBuilder) GetRemovedAccountsFromReplicationConfiguration(oldAcc []interface{}, newAcc []interface{}) []interface{} {
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

type Database struct {
	CreatedOn     sql.NullString `db:"created_on"`
	DBName        sql.NullString `db:"name"`
	IsDefault     sql.NullString `db:"is_default"`
	IsCurrent     sql.NullString `db:"is_current"`
	Origin        sql.NullString `db:"origin"`
	Owner         sql.NullString `db:"owner"`
	Comment       sql.NullString `db:"comment"`
	Options       sql.NullString `db:"options"`
	RetentionTime sql.NullString `db:"retention_time"`
	Kind          sql.NullString `db:"kind"`
	Budget        sql.NullString `db:"budget"`
}

func ScanDatabase(row *sqlx.Row) (*Database, error) {
	d := &Database{}
	e := row.StructScan(d)
	return d, e
}

func ListDatabases(sdb *sqlx.DB) ([]Database, error) {
	stmt := "SHOW DATABASES"
	rows, err := sdb.Queryx(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []Database{}
	if err := sqlx.StructScan(rows, &dbs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no databases found")
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return dbs, nil
}

func ListDatabase(sdb *sqlx.DB, databaseName string) (*Database, error) {
	stmt := fmt.Sprintf("SHOW DATABASES LIKE '%s'", databaseName)
	rows, err := sdb.Queryx(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []Database{}
	if err := sqlx.StructScan(rows, &dbs); err != nil {
		if errors.Is(err, sql.ErrNoRows) || len(dbs) == 0 {
			log.Println("[DEBUG] no databases found")
			return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	for _, d := range dbs {
		d := d
		if d.DBName.String == databaseName {
			return &d, nil
		}
	}
	return nil, errors.New("database not found")
}
