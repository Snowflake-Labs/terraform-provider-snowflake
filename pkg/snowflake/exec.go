package snowflake

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

func Exec(db *sql.DB, query string) error {
	log.Print("[DEBUG] stmt ", query)

	_, err := db.Exec(query)
	return err
}

// ExecAndGetId runs the query without returning the result, like sql.DB.Exec
// However, it then fetches the query id of the Exec call, and returns it
// This is needed for cases that need to use RESULT_SCAN to post-process the result of a SQL command like SHOW or DESC
func ExecAndGetId(db *sql.DB, query string) (string, error) {
	// Execute the query like a normal Exec call
	err := Exec(db, query)
	if err != nil {
		return "", err
	}

	// Fetch and return the latest query id from this Snowflake session, which is the Exec query
	var queryId string
	row := QueryRow(db, "SELECT LAST_QUERY_ID()::STRING AS QueryId;")
	row.Scan(&queryId)
	return queryId, err
}

// QueryRow will run stmt against the db and return the row. We use
// [DB.Unsafe](https://godoc.org/github.com/jmoiron/sqlx#DB.Unsafe) so that we can scan to structs
// without worrying about newly introduced columns
func QueryRow(db *sql.DB, stmt string) *sqlx.Row {
	log.Print("[DEBUG] stmt ", stmt)
	sdb := sqlx.NewDb(db, "snowflake").Unsafe()
	return sdb.QueryRowx(stmt)
}

// Query will run stmt against the db and return the rows. We use
// [DB.Unsafe](https://godoc.org/github.com/jmoiron/sqlx#DB.Unsafe) so that we can scan to structs
// without worrying about newly introduced columns
func Query(db *sql.DB, stmt string) (*sqlx.Rows, error) {
	sdb := sqlx.NewDb(db, "snowflake").Unsafe()
	return sdb.Queryx(stmt)
}
