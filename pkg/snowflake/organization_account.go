package snowflake

import (
	"database/sql"
    "log"

	"github.com/jmoiron/sqlx"
    "github.com/pkg/errors"
)

func ShowOrganizationAccounts() string {
	return `SHOW ORGANIZATION ACCOUNTS";`
}

type orgAccount struct {
	RegionGroup 	  sql.NullString `db:"region_group"`
	SnowflakeRegion   sql.NullString `db:"snowflake_region"`
	Name  			  sql.NullString `db:"name"`
	Edition  		  sql.NullString `db:"edition"`
	CreatedOn 		  sql.NullString `db:"created_on"`
	AccountUrl  	  sql.NullString `db:"account_url"`
	Comment  		  sql.NullString `db:"comment"`
	ManagedAccounts   sql.NullInt32  `db:"managed_accounts"`
	AccountLocatorUrl sql.NullString `db:"account_locator_url"`
}

func ScanOrganizationAccount(row *sqlx.Row) (*orgAccount, error) {
	acc := &orgAccount{}
	err := row.StructScan(acc)
	return acc, err
}

func ListOrganizationAccounts(db *sql.DB) ([]orgAccount, error) {
    stmt := ShowOrganizationAccounts()
    rows, err := Query(db, stmt)
    if err != nil {
		return nil, err
	}
	defer rows.Close()

    dbs := []orgAccount{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no organization accounts found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)

}