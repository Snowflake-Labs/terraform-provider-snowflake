package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type AccountEdition string

const (
	AccountEditionStandard         AccountEdition = "STANDARD"
	AccountEditionEnterprise       AccountEdition = "ENTERPRISE"
	AccountEditionBusinessCritical AccountEdition = "BUSINESS_CRITICAL"
)

// AccountBuilder abstracts the creation of SQL queries for a Snowflake account.
type AccountBuilder struct {
	db                 *sql.DB
	name               string
	adminName          string
	adminPassword      string
	adminRSAPublicKey  string
	firstName          string
	lastName           string
	email              string
	mustChangePassword bool
	edition            AccountEdition
	regionGroup        string
	region             string
	comment            string
}

// NewCreateAccountBuilder creates a new builder for a Snowflake account.
func NewCreateAccountBuilder(name, adminName, email string, edition AccountEdition, db *sql.DB) *AccountBuilder {
	return &AccountBuilder{
		name:      name,
		adminName: adminName,
		email:     email,
		edition:   edition,
		db:        db,
	}
}

// NewAlterAccountBuilder creates a new builder for a Snowflake account.
func NewAlterAccountBuilder(name string, db *sql.DB) *AccountBuilder {
	return &AccountBuilder{
		name: name,
		db:   db,
	}
}

// WithAdminPassword sets adminPassword.
func (b *AccountBuilder) WithAdminPassword(adminPassword string) *AccountBuilder {
	b.adminPassword = adminPassword
	return b
}

// WithAdminRSAPublicKey sets adminRSAPublicKey.
func (b *AccountBuilder) WithAdminRSAPublicKey(adminRSAPublicKey string) *AccountBuilder {
	b.adminRSAPublicKey = adminRSAPublicKey
	return b
}

// WithFirstName sets firstName.
func (b *AccountBuilder) WithFirstName(firstName string) *AccountBuilder {
	b.firstName = firstName
	return b
}

// WithLastName sets lastName.
func (b *AccountBuilder) WithLastName(lastName string) *AccountBuilder {
	b.lastName = lastName
	return b
}

// WithMustChangePassword sets mustChangePassword.
func (b *AccountBuilder) WithMustChangePassword(mustChangePassword bool) *AccountBuilder {
	b.mustChangePassword = mustChangePassword
	return b
}

// WithRegionGroup sets regionGroup.
func (b *AccountBuilder) WithRegionGroup(regionGroup string) *AccountBuilder {
	b.regionGroup = regionGroup
	return b
}

// WithRegion sets region.
func (b *AccountBuilder) WithRegion(region string) *AccountBuilder {
	b.region = region
	return b
}

// WithComment sets comment.
func (b *AccountBuilder) WithComment(comment string) *AccountBuilder {
	b.comment = comment
	return b
}

// Create returns the SQL query that will create a new account.
func (b *AccountBuilder) Create() (*Account, error) {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`CREATE ACCOUNT %s `, b.name))
	q.WriteString(fmt.Sprintf(` ADMIN_NAME = %s`, b.adminName))
	if (b.adminPassword == "") == (b.adminRSAPublicKey == "") {
		return nil, errors.New("either adminPassword or adminRSAPublicKey must be set, but not both")
	}
	if b.adminPassword != "" {
		q.WriteString(fmt.Sprintf(` ADMIN_PASSWORD = '%s'`, EscapeString(b.adminPassword)))
	}
	if b.adminRSAPublicKey != "" {
		q.WriteString(fmt.Sprintf(` ADMIN_RSA_PUBLIC_KEY = '%s'`, EscapeString(b.adminRSAPublicKey)))
	}
	if b.firstName != "" {
		q.WriteString(fmt.Sprintf(` FIRST_NAME = '%s'`, EscapeString(b.firstName)))
	}
	if b.lastName != "" {
		q.WriteString(fmt.Sprintf(` LAST_NAME ='%s'`, EscapeString(b.lastName)))
	}
	q.WriteString(fmt.Sprintf(` EMAIL = '%s'`, b.email))
	if b.mustChangePassword {
		q.WriteString(fmt.Sprintf(` MUST_CHANGE_PASSWORD = %t`, b.mustChangePassword))
	}

	q.WriteString(fmt.Sprintf(` EDITION = %s`, b.edition))

	if b.regionGroup != "" {
		q.WriteString(fmt.Sprintf(` REGION_GROUP = %s`, EscapeString(b.regionGroup)))
	}
	if b.region != "" {
		q.WriteString(fmt.Sprintf(` REGION = %s`, EscapeString(b.region)))
	}
	if b.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%s'`, EscapeString(b.comment)))
	}

	_, err := b.db.Exec(q.String())
	if err != nil {
		return nil, err
	}

	return ShowAccount(b.db, b.name)
}

// Rename returns the SQL query that will rename the account.
func (b *AccountBuilder) Rename(newName string) error {
	stmt := fmt.Sprintf(`ALTER ACCOUNT %s RENAME TO %s`, b.name, EscapeString(newName))
	_, err := b.db.Exec(stmt)
	return err
}

// SetComment returns the SQL query that will set the comment on the account.
func (b *AccountBuilder) SetComment(comment string) error {
	stmt := fmt.Sprintf(`COMMENT ON ACCOUNT %s IS '%s'`, b.name, EscapeString(comment))
	_, err := b.db.Exec(stmt)
	return err
}

type Account struct {
	OrganizationName                     sql.NullString `db:"organization_name"`
	AccountName                          sql.NullString `db:"account_name"`
	RegionGroup                          sql.NullString `db:"region_group"`
	SnowflakeRegion                      sql.NullString `db:"snowflake_region"`
	Edition                              sql.NullString `db:"edition"`
	AccountURL                           sql.NullString `db:"account_url"`
	CreatedOn                            sql.NullString `db:"created_on"`
	Comment                              sql.NullString `db:"comment"`
	AccountLocator                       sql.NullString `db:"account_locator"`
	AccountLocatorURL                    sql.NullString `db:"account_locator_url"`
	ManagedAccounts                      sql.NullString `db:"managed_accounts"`
	ConsumptionBillingEntityName         sql.NullString `db:"consumption_billing_entity_name"`
	MarketplaceConsumerBillingEntityName sql.NullString `db:"marketplace_consumer_billing_entity_name"`
	MarketplaceProviderBillingEntityName sql.NullString `db:"marketplace_provider_billing_entity_name"`
	OldAccountURL                        sql.NullString `db:"old_account_url"`
	IsOrgAdmin                           sql.NullBool   `db:"is_org_admin"`
}

// Show returns the SQL query that will show a specific account by pattern.
func ShowAccount(db *sql.DB, pattern string) (*Account, error) {
	stmt := fmt.Sprintf("SHOW ORGANIZATION ACCOUNTS LIKE '%s'", pattern)
	rows, err := db.Query(stmt, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	accounts := []Account{}
	log.Printf("[DEBUG] accounts: %v", accounts)
	if err := sqlx.StructScan(rows, &accounts); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[DEBUG] no accounts found for pattern: %s", pattern)
			return nil, err
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return &accounts[0], nil
}
