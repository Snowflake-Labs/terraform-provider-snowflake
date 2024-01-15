package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/snowflakedb/gosnowflake"
)

var (
	instrumentedSQL         bool
	gosnowflakeLoggingLevel string
)

func init() {
	instrumentedSQL = os.Getenv("SF_TF_NO_INSTRUMENTED_SQL") == ""
	gosnowflakeLoggingLevel = os.Getenv("SF_TF_GOSNOWFLAKE_LOG_LEVEL")
}

type Client struct {
	config         *gosnowflake.Config
	db             *sqlx.DB
	sessionID      string
	accountLocator string
	dryRun         bool
	traceLogs      []string

	// System-Defined Functions
	ContextFunctions     ContextFunctions
	ConversionFunctions  ConversionFunctions
	SystemFunctions      SystemFunctions
	ReplicationFunctions ReplicationFunctions

	// DDL Commands
	Accounts         Accounts
	Alerts           Alerts
	ApplicationRoles ApplicationRoles
	Applications     Applications
	Comments         Comments
	DatabaseRoles    DatabaseRoles
	Databases        Databases
	DynamicTables    DynamicTables
	ExternalTables   ExternalTables
	EventTables      EventTables
	FailoverGroups   FailoverGroups
	FileFormats      FileFormats
	Functions        Functions
	Grants           Grants
	MaskingPolicies  MaskingPolicies
	NetworkPolicies  NetworkPolicies
	Parameters       Parameters
	PasswordPolicies PasswordPolicies
	Pipes            Pipes
	Procedures       Procedures
	ResourceMonitors ResourceMonitors
	Roles            Roles
	Schemas          Schemas
	SessionPolicies  SessionPolicies
	Sessions         Sessions
	Shares           Shares
	Stages           Stages
	Streams          Streams
	Tables           Tables
	Tags             Tags
	Tasks            Tasks
	Users            Users
	Views            Views
	Warehouses       Warehouses
}

func (c *Client) GetAccountLocator() string {
	return c.accountLocator
}

func (c *Client) GetConfig() *gosnowflake.Config {
	return c.config
}

func (c *Client) GetConn() *sqlx.DB {
	return c.db
}

func NewDefaultClient() (*Client, error) {
	return NewClient(nil)
}

func NewDryRunClient() *Client {
	client := &Client{
		dryRun:    true,
		traceLogs: []string{},
	}
	client.initialize()
	return client
}

func NewClient(cfg *gosnowflake.Config) (*Client, error) {
	var err error
	if cfg == nil {
		log.Printf("[DEBUG] Searching for default config in credentials chain...\n")
		cfg = DefaultConfig()
	}

	var client *Client
	// register the snowflake driver if it hasn't been registered yet

	driverName := "snowflake"
	if instrumentedSQL {
		if !slices.Contains(sql.Drivers(), "snowflake-instrumented") {
			log.Println("[DEBUG] Registering snowflake-instrumented driver")
			logger := instrumentedsql.LoggerFunc(func(ctx context.Context, s string, kv ...interface{}) {
				switch s {
				case "sql-conn-query", "sql-conn-exec":
					log.Printf("[DEBUG] %s: %v (%s)\n", s, kv, ctx.Value(snowflakeAccountLocatorContextKey))
				default:
					return
				}
			})
			sql.Register("snowflake-instrumented", instrumentedsql.WrapDriver(new(gosnowflake.SnowflakeDriver), instrumentedsql.WithLogger(logger)))
		}
		driverName = "snowflake-instrumented"
	}

	if gosnowflakeLoggingLevel != "" {
		cfg.Tracing = gosnowflakeLoggingLevel
	}

	dsn, err := gosnowflake.DSN(cfg)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open snowflake connection: %w", err)
	}

	client = &Client{
		// snowflake does not adhere to the normal sql driver interface, so we have to use unsafe
		db:     db.Unsafe(),
		config: cfg,
	}
	client.initialize()

	err = client.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping snowflake: %w", err)
	}
	ctx := context.Background()
	currentAccount, err := client.ContextFunctions.CurrentAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current account: %w", err)
	}
	client.accountLocator = currentAccount

	sessionID, err := client.ContextFunctions.CurrentSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current session: %w", err)
	}
	client.sessionID = sessionID
	log.Printf("[DEBUG] connection success! Account: %s, Session identifier: %s\n", currentAccount, sessionID)

	return client, nil
}

func NewClientFromDB(db *sql.DB) *Client {
	dbx := sqlx.NewDb(db, "snowflake")
	client := &Client{
		db: dbx.Unsafe(),
	}
	client.initialize()
	return client
}

func (c *Client) initialize() {
	c.Accounts = &accounts{client: c}
	c.Alerts = &alerts{client: c}
	c.ApplicationRoles = &applicationRoles{client: c}
	c.Applications = &applications{client: c}
	c.Comments = &comments{client: c}
	c.ContextFunctions = &contextFunctions{client: c}
	c.ConversionFunctions = &conversionFunctions{client: c}
	c.DatabaseRoles = &databaseRoles{client: c}
	c.Databases = &databases{client: c}
	c.DynamicTables = &dynamicTables{client: c}
	c.ExternalTables = &externalTables{client: c}
	c.EventTables = &eventTables{client: c}
	c.FailoverGroups = &failoverGroups{client: c}
	c.FileFormats = &fileFormats{client: c}
	c.Functions = &functions{client: c}
	c.Grants = &grants{client: c}
	c.MaskingPolicies = &maskingPolicies{client: c}
	c.NetworkPolicies = &networkPolicies{client: c}
	c.Parameters = &parameters{client: c}
	c.PasswordPolicies = &passwordPolicies{client: c}
	c.Pipes = &pipes{client: c}
	c.Procedures = &procedures{client: c}
	c.ReplicationFunctions = &replicationFunctions{client: c}
	c.ResourceMonitors = &resourceMonitors{client: c}
	c.Roles = &roles{client: c}
	c.Schemas = &schemas{client: c}
	c.SessionPolicies = &sessionPolicies{client: c}
	c.Sessions = &sessions{client: c}
	c.Shares = &shares{client: c}
	c.Stages = &stages{client: c}
	c.Streams = &streams{client: c}
	c.SystemFunctions = &systemFunctions{client: c}
	c.Tables = &tables{client: c}
	c.Tags = &tags{client: c}
	c.Tasks = &tasks{client: c}
	c.Users = &users{client: c}
	c.Views = &views{client: c}
	c.Warehouses = &warehouses{client: c}
}

func (c *Client) TraceLogs() []string {
	return c.traceLogs
}

func (c *Client) Ping() error {
	return c.db.Ping()
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

type snowflakeAccountLocatorContext string

const (
	snowflakeAccountLocatorContextKey snowflakeAccountLocatorContext = "snowflake_account_locator"
)

// Exec executes a query that does not return rows.
func (c *Client) exec(ctx context.Context, sql string) (sql.Result, error) {
	if c.dryRun {
		c.traceLogs = append(c.traceLogs, sql)
		log.Printf("[DEBUG] sql-conn-exec-dry: %v\n", sql)
		return nil, nil
	}
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	result, err := c.db.ExecContext(ctx, sql)
	return result, decodeDriverError(err)
}

// query runs a query and returns the rows. dest is expected to be a slice of structs.
func (c *Client) query(ctx context.Context, dest interface{}, sql string) error {
	if c.dryRun {
		c.traceLogs = append(c.traceLogs, sql)
		log.Printf("[DEBUG] sql-conn-query-dry: %v\n", sql)
		return nil
	}
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	return decodeDriverError(c.db.SelectContext(ctx, dest, sql))
}

// queryOne runs a query and returns one row. dest is expected to be a pointer to a struct.
func (c *Client) queryOne(ctx context.Context, dest interface{}, sql string) error {
	if c.dryRun {
		c.traceLogs = append(c.traceLogs, sql)
		log.Printf("[DEBUG] sql-conn-query-one-dry: %v\n", sql)
		return nil
	}
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	return decodeDriverError(c.db.GetContext(ctx, dest, sql))
}
