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

func QueryRow(db *sql.DB, stmt string) *sqlx.Row {
	log.Print("[DEBUG] stmt ", stmt)
	sdb := sqlx.NewDb(db, "snowflake")
	return sdb.QueryRowx(stmt)
}

func Query(db *sql.DB, stmt string) (*sqlx.Rows, err) {
	sdb := sqlx.NewDb(db, "snowflake")
	return sdb.Queryx(stmt)
}
