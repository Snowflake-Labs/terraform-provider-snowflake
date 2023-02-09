package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// FailoverGroupBuilder abstracts the creation of SQL queries for a Snowflake file format.
type FailoverGroupBuilder struct {
	name                              string
	objectTypes                       []string
	allowedDatabases                  []string
	allowedShares                     []string
	allowedIntegrationTypes           []string
	allowedAccounts                   []string
	ignoreEditionCheck                bool
	replicationScheduleInterval       int
	replicationScheduleCronExpression string
	replicationScheduleTimeZone       string
}

func (b *FailoverGroupBuilder) WithName(name string) *FailoverGroupBuilder {
	b.name = name
	return b
}

func (b *FailoverGroupBuilder) WithObjectTypes(objectTypes []string) *FailoverGroupBuilder {
	b.objectTypes = objectTypes
	return b
}

func (b *FailoverGroupBuilder) WithAllowedDatabases(allowedDatabases []string) *FailoverGroupBuilder {
	b.allowedDatabases = allowedDatabases
	return b
}

func (b *FailoverGroupBuilder) WithAllowedShares(allowedShares []string) *FailoverGroupBuilder {
	b.allowedShares = allowedShares
	return b
}

func (b *FailoverGroupBuilder) WithAllowedIntegrationTypes(allowedIntegrationTypes []string) *FailoverGroupBuilder {
	b.allowedIntegrationTypes = allowedIntegrationTypes
	return b
}

func (b *FailoverGroupBuilder) WithAllowedAccounts(allowedAccounts []string) *FailoverGroupBuilder {
	b.allowedAccounts = allowedAccounts
	return b
}

func (b *FailoverGroupBuilder) WithIgnoreEditionCheck(ignoreEditionCheck bool) *FailoverGroupBuilder {
	b.ignoreEditionCheck = ignoreEditionCheck
	return b
}

func (b *FailoverGroupBuilder) WithReplicationScheduleInterval(replicationScheduleInterval int) *FailoverGroupBuilder {
	b.replicationScheduleInterval = replicationScheduleInterval
	return b
}

func (b *FailoverGroupBuilder) WithReplicationScheduleCronExpression(replicationScheduleCronExpression string) *FailoverGroupBuilder {
	b.replicationScheduleCronExpression = replicationScheduleCronExpression
	return b
}

func (b *FailoverGroupBuilder) WithReplicationScheduleTimeZone(replicationScheduleTimeZone string) *FailoverGroupBuilder {
	b.replicationScheduleTimeZone = replicationScheduleTimeZone
	return b
}

// CreateFailoverGroup returns a pointer to a Builder that abstracts the DDL operations for a failover group.
func NewFailoverGroupBuilder(name string) *FailoverGroupBuilder {
	return &FailoverGroupBuilder{
		name: name,
	}
}

func (b *FailoverGroupBuilder) CreateFromReplica(name string) string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(` CREATE FAILOVER GROUP %v`, b.name))
	q.WriteString(fmt.Sprintf(` AS REPLICA OF %v`, name))
	return q.String()
}

// Create returns the SQL query that will create a new failover group.
func (b *FailoverGroupBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` FAILOVER GROUP %v`, b.name))
	q.WriteString(fmt.Sprintf(` OBJECT_TYPES = %v`, strings.Join(b.objectTypes, ",")))

	if len(b.allowedDatabases) > 0 {
		allowedDatabasesWithQuotes := make([]string, len(b.allowedDatabases))
		for i, v := range b.allowedDatabases {
			allowedDatabasesWithQuotes[i] = "\"" + v + "\""
		}
		q.WriteString(fmt.Sprintf(` ALLOWED_DATABASES = %v`, strings.Join(allowedDatabasesWithQuotes, ",")))
	}

	if len(b.allowedShares) > 0 {
		q.WriteString(fmt.Sprintf(` ALLOWED_SHARES = %v`, strings.Join(b.allowedShares, ",")))
	}

	if len(b.allowedIntegrationTypes) > 0 {
		q.WriteString(fmt.Sprintf(` ALLOWED_INTEGRATION_TYPES = %v`, strings.Join(b.allowedIntegrationTypes, ",")))
	}

	if len(b.allowedAccounts) > 0 {
		q.WriteString(fmt.Sprintf(` ALLOWED_ACCOUNTS = %v`, strings.Join(b.allowedAccounts, ",")))
	}

	if b.ignoreEditionCheck {
		q.WriteString(` IGNORE_EDITION_CHECK`)
	}

	if b.replicationScheduleCronExpression != "" {
		q.WriteString(fmt.Sprintf(" REPLICATION_SCHEDULE = 'USING CRON %v", b.replicationScheduleCronExpression))
		if b.replicationScheduleTimeZone != "" {
			q.WriteString(fmt.Sprintf(" %v", b.replicationScheduleTimeZone))
		}
		q.WriteString("'")
	}
	if b.replicationScheduleInterval > 0 {
		q.WriteString(fmt.Sprintf(" REPLICATION_SCHEDULE = '%v MINUTE'", b.replicationScheduleInterval))
	}

	return q.String()
}

// Rename returns the SQL query that will rename a failover group.
func (b *FailoverGroupBuilder) Rename(name string) string {
	s := fmt.Sprintf(`ALTER FAILOVER GROUP %v RENAME TO %v`, b.name, name)
	b.name = name
	return s
}

// ChangeObjectTypes returns the SQL query that will change the object types of a failover group.
func (b *FailoverGroupBuilder) ChangeObjectTypes(objectTypes []string) string {
	s := fmt.Sprintf(`ALTER FAILOVER GROUP %v SET OBJECT_TYPES = %v`, b.name, strings.Join(objectTypes, ","))
	b.objectTypes = objectTypes
	return s
}

// ChangeReplicationCronSchedule returns the SQL query that will change the replication schedule of a failover group.
func (b *FailoverGroupBuilder) ChangeReplicationCronSchedule(replicationScheduleCronExpression string, replicationScheduleTimeZone string) string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`ALTER FAILOVER GROUP %v SET REPLICATION_SCHEDULE = 'USING CRON %v`, b.name, replicationScheduleCronExpression))

	if replicationScheduleTimeZone != "" {
		q.WriteString(fmt.Sprintf(` %v`, replicationScheduleTimeZone))
	}
	q.WriteString(`'`)
	b.replicationScheduleCronExpression = replicationScheduleCronExpression
	b.replicationScheduleTimeZone = replicationScheduleTimeZone
	return q.String()
}

// ChangeReplicationIntervalSchedule returns the SQL query that will change the replication schedule of a failover group.
func (b *FailoverGroupBuilder) ChangeReplicationIntervalSchedule(replicationScheduleInterval int) string {
	s := fmt.Sprintf(`ALTER FAILOVER GROUP %v SET REPLICATION_SCHEDULE = %v MINUTE`, b.name, replicationScheduleInterval)
	b.replicationScheduleInterval = replicationScheduleInterval
	return s
}

// ChangeAllowedIntegrationTypes returns the SQL query that will change the allowed integration types of a failover group.
func (b *FailoverGroupBuilder) ChangeAllowedIntegrationTypes(allowedIntegrationTypes []string) string {
	s := fmt.Sprintf(`ALTER FAILOVER GROUP %v SET ALLOWED_INTEGRATION_TYPES = %v`, b.name, strings.Join(allowedIntegrationTypes, ","))
	b.allowedIntegrationTypes = allowedIntegrationTypes
	return s
}

// AddAllowedDatabases returns the SQL query that will change the allowed databases of a failover group.
func (b *FailoverGroupBuilder) AddAllowedDatabases(allowedDatabases []string) string {
	return fmt.Sprintf(`ALTER FAILOVER GROUP %v ADD %v TO ALLOWED_DATABASES`, b.name, strings.Join(allowedDatabases, ","))
}

// RemoveAllowedDatabases returns the SQL query that will change the allowed databases of a failover group.
func (b *FailoverGroupBuilder) RemoveAllowedDatabases(allowedDatabases []string) string {
	return fmt.Sprintf(`ALTER FAILOVER GROUP %v REMOVE %v FROM ALLOWED_DATABASES`, b.name, strings.Join(allowedDatabases, ","))
}

// AddAllowedShares returns the SQL query that will change the allowed shares of a failover group.
func (b *FailoverGroupBuilder) AddAllowedShares(allowedShares []string) string {
	return fmt.Sprintf(`ALTER FAILOVER GROUP %v ADD %v TO ALLOWED_SHARES`, b.name, strings.Join(allowedShares, ","))
}

// RemoveAllowedShares returns the SQL query that will change the allowed shares of a failover group.
func (b *FailoverGroupBuilder) RemoveAllowedShares(allowedShares []string) string {
	return fmt.Sprintf(`ALTER FAILOVER GROUP %v REMOVE %v FROM ALLOWED_SHARES`, b.name, strings.Join(allowedShares, ","))
}

// AddAllowedAccounts returns the SQL query that will change the allowed accounts of a failover group.
func (b *FailoverGroupBuilder) AddAllowedAccounts(allowedAccounts []string) string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`ALTER FAILOVER GROUP %v ADD %v TO ALLOWED_ACCOUNTS`, b.name, strings.Join(allowedAccounts, ",")))
	if b.ignoreEditionCheck {
		q.WriteString(` IGNORE EDITION CHECK`)
	}
	return q.String()
}

// RemoveAllowedAccounts returns the SQL query that will change the allowed accounts of a failover group.
func (b *FailoverGroupBuilder) RemoveAllowedAccounts(allowedAccounts []string) string {
	return fmt.Sprintf(`ALTER FAILOVER GROUP %v REMOVE %v FROM ALLOWED_ACCOUNTS`, b.name, strings.Join(allowedAccounts, ","))
}

// Drop returns the SQL query that will drop a failover group.
func (b *FailoverGroupBuilder) Drop() string {
	return fmt.Sprintf(`DROP FAILOVER GROUP %v`, b.name)
}

// Show returns the SQL query that will show a failover group.
func (b *FailoverGroupBuilder) Show() string {
	return "SHOW FAILOVER GROUPS"
}

// ListFailoverGroups returns a list of all failover groups in the account.
func ListFailoverGroups(db *sql.DB, accountLocator string) ([]FailoverGroup, error) {
	stmt := fmt.Sprintf("SHOW FAILOVER GROUPS IN ACCOUNT %s", accountLocator)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	v := []FailoverGroup{}
	err = sqlx.StructScan(rows, &v)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("[DEBUG] no failover groups found")
		return nil, nil
	}

	return v, nil
}

type FailoverGroup struct {
	SnowflakeRegion         sql.NullString `db:"snowflake_region"`
	CreatedOn               sql.NullString `db:"created_on"`
	AccountName             sql.NullString `db:"account_name"`
	Name                    sql.NullString `db:"name"`
	IsPrimary               sql.NullString `db:"is_primary"`
	Primary                 sql.NullString `db:"primary"`
	ObjectTypes             sql.NullString `db:"object_types"`
	AllowedIntegrationTypes sql.NullString `db:"allowed_integration_types"`
	AllowedAccounts         sql.NullString `db:"allowed_accounts"`
	OrganizationName        sql.NullString `db:"organization_name"`
	AccountLocator          sql.NullString `db:"account_locator"`
	ReplicationSchedule     sql.NullString `db:"replication_schedule"`
	SecondaryState          sql.NullString `db:"secondary_state"`
}

type failoverGroupAllowedDatabase struct {
	Name sql.NullString `db:"name"`
}

func ShowDatabasesInFailoverGroup(name string, db *sql.DB) ([]string, error) {
	stmt := fmt.Sprintf(`SHOW DATABASES IN FAILOVER GROUP %v`, name)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, fmt.Errorf("error listing allowed databases for failover group %v err = %w", name, err)
	}
	defer rows.Close()

	failoverGroupAllowedDatabase := []failoverGroupAllowedDatabase{}
	err = sqlx.StructScan(rows, &failoverGroupAllowedDatabase)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("[DEBUG] no failover group databases found")
		return nil, nil
	}

	result := make([]string, 0, len(failoverGroupAllowedDatabase))
	for _, v := range failoverGroupAllowedDatabase {
		result = append(result, v.Name.String)
	}
	return result, nil
}

type failoverGroupAllowedShare struct {
	Name sql.NullString `db:"name"`
}

func ShowSharesInFailoverGroup(name string, db *sql.DB) ([]string, error) {
	stmt := fmt.Sprintf(`SHOW SHARES IN FAILOVER GROUP %v`, name)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, fmt.Errorf("error listing allowed shares for failover group %v err = %w", name, err)
	}

	defer rows.Close()

	failoverGroupAllowedShares := []failoverGroupAllowedShare{}
	if err := sqlx.StructScan(rows, &failoverGroupAllowedShares); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no failover group shares found")
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan row for %s err = %w", stmt, err)
	}

	result := make([]string, 0, len(failoverGroupAllowedShares))
	for _, v := range failoverGroupAllowedShares {
		result = append(result, v.Name.String)
	}
	return result, nil
}
