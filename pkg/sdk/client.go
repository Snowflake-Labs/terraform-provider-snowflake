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
	config *gosnowflake.Config
	db     *sqlx.DB
	dryRun bool

	ContextFunctions ContextFunctions
	Databases        Databases
	Grants           Grants
	MaskingPolicies  MaskingPolicies
	PasswordPolicies PasswordPolicies
	Sessions         Sessions
	Shares           Shares
	SystemFunctions  SystemFunctions
	Warehouses       Warehouses
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

	// register the snowflake driver if it hasn't been registered yet
	if !slices.Contains(sql.Drivers(), "snowflake-instrumented") {
		logger := instrumentedsql.LoggerFunc(func(ctx context.Context, s string, kv ...interface{}) {
			switch s {
			case "sql-conn-query", "sql-conn-exec":
				log.Printf("[DEBUG] %s: %v\n", s, kv)
			default:
				return
			}
		})
		sql.Register("snowflake-instrumented", instrumentedsql.WrapDriver(gosnowflake.SnowflakeDriver{}, instrumentedsql.WithLogger(logger)))
	}

	dsn, err := gosnowflake.DSN(cfg)
	if err != nil {
		return nil, decodeDriverError(err)
	}

	db, err := sqlx.Connect("snowflake-instrumented", dsn)
	if err != nil {
		return nil, fmt.Errorf("open snowflake connection: %w", err)
	}

	client := &Client{
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
	sessionID, err := client.ContextFunctions.CurrentSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current session: %w", err)
	}
	log.Printf("[DEBUG] connection success! Session identifier: %s\n", sessionID)

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
	b := &sqlBuilder{}
	c.ContextFunctions = &contextFunctions{client: c, builder: b}
	c.Databases = &databases{client: c, builder: b}
	c.Grants = &grants{client: c, builder: b}
	c.MaskingPolicies = &maskingPolicies{client: c, builder: b}
	c.PasswordPolicies = &passwordPolicies{client: c, builder: b}
	c.Sessions = &sessions{client: c, builder: b}
	c.Shares = &shares{client: c, builder: b}
	c.SystemFunctions = &systemFunctions{client: c, builder: b}
	c.Warehouses = &warehouses{client: c, builder: b}
}

func (c *Client) SetDryRun(dryRun bool) {
	c.dryRun = dryRun
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

// Exec executes a query that does not return rows.
func (c *Client) exec(ctx context.Context, sql string) (sql.Result, error) {
	if !c.dryRun {
		return c.db.ExecContext(ctx, sql)
	}
	return nil, nil
}

// query runs a query and returns the rows. dest is expected to be a slice of structs.
func (c *Client) query(ctx context.Context, dest interface{}, sql string) error {
	if !c.dryRun {
		return c.db.SelectContext(ctx, dest, sql)
	}
	return nil
}

// queryOne runs a query and returns one row. dest is expected to be a pointer to a struct.
func (c *Client) queryOne(ctx context.Context, dest interface{}, sql string) error {
	if !c.dryRun {
		return c.db.GetContext(ctx, dest, sql)
	}
	return nil
}
