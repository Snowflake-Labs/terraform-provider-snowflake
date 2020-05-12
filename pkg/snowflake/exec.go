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
