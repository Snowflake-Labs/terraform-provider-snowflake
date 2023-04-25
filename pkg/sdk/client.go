package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"

	"github.com/snowflakedb/gosnowflake"
)

// ObjectType is the type of object.
type ObjectType string

func (o ObjectType) String() string {
	return string(o)
}

func DefaultConfig() *gosnowflake.Config {
	cfg := &gosnowflake.Config{
		Account:   os.Getenv("SNOWFLAKE_ACCOUNT"),
		User:      os.Getenv("SNOWFLAKE_USER"),
		Password:  os.Getenv("SNOWFLAKE_PASSWORD"),
		Region:    os.Getenv("SNOWFLAKE_REGION"),
		Role:      os.Getenv("SNOWFLAKE_ROLE"),
		Host:      os.Getenv("SNOWFLAKE_HOST"),
		Warehouse: os.Getenv("SNOWFLAKE_WAREHOUSE"),
	}
	// us-west-2 is Snowflake's default region, but if you actually specify that it won't trigger the default code
	//  https://github.com/snowflakedb/gosnowflake/blob/52137ce8c32eaf93b0bd22fc5c7297beff339812/dsn.go#L61
	if cfg.Region == "us-west-2" {
		cfg.Region = ""
	}
	return cfg
}

type Client struct {
	db     *sqlx.DB
	dryRun bool
	sqlBuilder

	PasswordPolicies PasswordPolicies
}

func NewDefaultClient() (*Client, error) {
	return NewClient(nil)
}

func NewClient(cfg *gosnowflake.Config) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	dsn, err := gosnowflake.DSN(cfg)
	if err != nil {
		return nil, fmt.Errorf("build dsn for snowflake connection: %w", err)
	}

	logger := instrumentedsql.LoggerFunc(func(ctx context.Context, fn string, kv ...interface{}) {
		switch fn {
		case "sql-conn-query", "sql-conn-exec":
			log.Printf("[DEBUG] %s: %v", fn, kv)
		default:
			return
		}
	})
	sql.Register("snowflake-instrumented", instrumentedsql.WrapDriver(&gosnowflake.SnowflakeDriver{}, instrumentedsql.WithLogger(logger)))
	db, err := sqlx.Connect("snowflake-instrumented", dsn)
	if err != nil {
		return nil, fmt.Errorf("open snowflake connection: %w", err)
	}
	client := &Client{
		// snowflake does not adhere to the normal sql driver interface, so we have to use unsafe
		db: db.Unsafe(),
	}
	client.initialize()

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
	c.PasswordPolicies = &passwordPolicies{client: c}
}

func (c *Client) SetDryRun(dryRun bool) {
	c.dryRun = dryRun
}

func (c *Client) Ping() error {
	return c.db.Ping()
}

func (c *Client) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

// Exec executes a query that does not return rows.W
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
