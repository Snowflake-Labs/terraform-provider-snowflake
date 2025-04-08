package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func Exec(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	return err
}

func ExecMulti(db *sql.DB, queries []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, query := range queries {
		_, err = tx.Exec(query)
		if err != nil {
			return tx.Rollback()
		}
	}
	return tx.Commit()
}

// QueryRow will run stmt against the db and return the row. We use
// [DB.Unsafe](https://godoc.org/github.com/jmoiron/sqlx#DB.Unsafe) so that we can scan to structs
// without worrying about newly introduced columns.
func QueryRow(db *sql.DB, stmt string) *sqlx.Row {
	return sqlx.NewDb(db, "snowflake").Unsafe().QueryRowx(stmt)
}

// Query will run stmt against the db and return the rows. We use
// [DB.Unsafe](https://godoc.org/github.com/jmoiron/sqlx#DB.Unsafe) so that we can scan to structs
// without worrying about newly introduced columns.
func Query(db *sql.DB, stmt string) (*sqlx.Rows, error) {
	return sqlx.NewDb(db, "snowflake").Unsafe().Queryx(stmt)
}
