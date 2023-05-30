package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
	"golang.org/x/exp/slices"

	"github.com/snowflakedb/gosnowflake"
)

type Client struct {
	config         *gosnowflake.Config
	db             *sqlx.DB
	sessionID      string
	accountLocator string

	Accounts             Accounts
	Comments             Comments
	ContextFunctions     ContextFunctions
	Databases            Databases
	FailoverGroups       FailoverGroups
	Grants               Grants
	MaskingPolicies      MaskingPolicies
	PasswordPolicies     PasswordPolicies
	ReplicationFunctions ReplicationFunctions
	ResourceMonitors     ResourceMonitors
	Roles                Roles
	SessionPolicies      SessionPolicies
	Sessions             Sessions
	Shares               Shares
	SystemFunctions      SystemFunctions
	Warehouses           Warehouses
}

func NewDefaultClient() (*Client, error) {
	return NewClient(nil)
}

func NewClient(cfg *gosnowflake.Config) (*Client, error) {
	var err error
	if cfg == nil {
		log.Printf("[DEBUG] Searching for default config in credentials chain...\n")
		cfg = DefaultConfig()
	}

	var client *Client
	// register the snowflake driver if it hasn't been registered yet
	if !slices.Contains(sql.Drivers(), "snowflake-instrumented") {
		logger := instrumentedsql.LoggerFunc(func(ctx context.Context, s string, kv ...interface{}) {
			switch s {
			case "sql-conn-query", "sql-conn-exec":
				log.Printf("[DEBUG] %s: %v (%s)\n", s, kv, ctx.Value(snowflakeAccountLocatorContextKey))
			default:
				return
			}
		})
		sql.Register("snowflake-instrumented", instrumentedsql.WrapDriver(gosnowflake.SnowflakeDriver{}, instrumentedsql.WithLogger(logger)))
	}

	dsn, err := gosnowflake.DSN(cfg)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("snowflake-instrumented", dsn)
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
	c.Comments = &comments{client: c}
	c.ContextFunctions = &contextFunctions{client: c}
	c.Databases = &databases{client: c}
	c.FailoverGroups = &failoverGroups{client: c}
	c.Grants = &grants{client: c}
	c.MaskingPolicies = &maskingPolicies{client: c}
	c.PasswordPolicies = &passwordPolicies{client: c}
	c.ReplicationFunctions = &replicationFunctions{client: c}
	c.ResourceMonitors = &resourceMonitors{client: c}
	c.Roles = &roles{client: c}
	c.SessionPolicies = &sessionPolicies{client: c}
	c.Sessions = &sessions{client: c}
	c.Shares = &shares{client: c}
	c.SystemFunctions = &systemFunctions{client: c}
	c.Warehouses = &warehouses{client: c}
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
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	result, err := c.db.ExecContext(ctx, sql)
	return result, decodeDriverError(err)
}

// query runs a query and returns the rows. dest is expected to be a slice of structs.
func (c *Client) query(ctx context.Context, dest interface{}, sql string) error {
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	return decodeDriverError(c.db.SelectContext(ctx, dest, sql))
}

// queryOne runs a query and returns one row. dest is expected to be a pointer to a struct.
func (c *Client) queryOne(ctx context.Context, dest interface{}, sql string) error {
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	return decodeDriverError(c.db.GetContext(ctx, dest, sql))
}
