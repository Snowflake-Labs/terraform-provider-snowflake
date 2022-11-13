package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/pkg/errors"
	"github.com/snowflakedb/gosnowflake"
)

var (
	ErrNoRecord = errors.New("record not found")
)

type Config struct {
	Account   string
	User      string
	Password  string
	Region    string
	Role      string
	Host      string
	Warehouse string
}

func DefaultConfig() *Config {
	config := &Config{
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
	if config.Region == "us-west-2" {
		config.Region = ""
	}
	return config
}

type Client struct {
	conn *sql.DB

	Users           Users
	Roles           Roles
	Warehouses      Warehouses
	Databases       Databases
	Schemas         Schemas
	Tables          Tables
	NetworkPolicies NetworkPolicies
}

func NewClient(cfg *Config) (*Client, error) {
	config := DefaultConfig()
	if cfg != nil {
		if cfg.Account != "" {
			config.Account = cfg.Account
		}
		if cfg.User != "" {
			config.User = cfg.User
		}
		if cfg.Password != "" {
			config.Password = cfg.Password
		}
		// us-west-2 is Snowflake's default region, but if you actually specify that it won't trigger the default code
		//  https://github.com/snowflakedb/gosnowflake/blob/52137ce8c32eaf93b0bd22fc5c7297beff339812/dsn.go#L61
		if cfg.Region != "" && cfg.Region != "us-west-2" {
			config.Region = cfg.Region
		}
		if cfg.Role != "" {
			config.Role = cfg.Role
		}
		if cfg.Host != "" {
			config.Host = cfg.Host
			// if host is set trust it and do not use the region
			config.Region = ""
		}
		if cfg.Warehouse != "" {
			config.Warehouse = cfg.Warehouse
		}
	}

	dsn, err := gosnowflake.DSN(&gosnowflake.Config{
		Account:   config.Account,
		User:      config.User,
		Password:  config.Password,
		Region:    config.Region,
		Role:      config.Role,
		Warehouse: config.Warehouse,
	})
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
	conn, err := sql.Open("snowflake-instrumented", dsn)
	if err != nil {
		return nil, fmt.Errorf("open snowflake connection: %w", err)
	}
	client := &Client{
		conn: conn,
	}

	client.Users = &users{client: client}
	client.Roles = &roles{client: client}
	client.Warehouses = &warehouses{client: client}
	client.Databases = &databases{client: client}
	client.Schemas = &schemas{client: client}
	client.Tables = &tables{client: client}
	client.NetworkPolicies = &networkPolicies{client: client}

	return client, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) exec(ctx context.Context, sql string) (sql.Result, error) {
	return c.conn.ExecContext(ctx, sql)
}

func (c *Client) query(ctx context.Context, sql string) (*sqlx.Rows, error) {
	return sqlx.NewDb(c.conn, "snowflake-instrumented").Unsafe().QueryxContext(ctx, sql)
}

// drop a resource
func (c *Client) drop(ctx context.Context, resource string, name string) error {
	sql := fmt.Sprintf("DROP %s %s", resource, name)
	if _, err := c.exec(ctx, sql); err != nil {
		return fmt.Errorf("db exec: %w", err)
	}
	return nil
}

// rename a resource
func (c *Client) rename(ctx context.Context, resource string, old string, new string) error {
	sql := fmt.Sprintf("ALTER %s %s RENAME TO %s", resource, old, new)
	if _, err := c.exec(ctx, sql); err != nil {
		return fmt.Errorf("db exec: %w", err)
	}
	return nil
}

// clone a resource
func (c *Client) clone(ctx context.Context, resource string, source string, dest string) error {
	sql := fmt.Sprintf("CREATE %s %s CLONE %s", resource, dest, source)
	if _, err := c.exec(ctx, sql); err != nil {
		return fmt.Errorf("db exec: %w", err)
	}
	return nil
}

// use a resource
func (c *Client) use(ctx context.Context, resource string, name string) error {
	sql := fmt.Sprintf("USE %s %s", resource, name)
	if _, err := c.exec(ctx, sql); err != nil {
		return fmt.Errorf("db exec: %w", err)
	}
	return nil
}

// read a resource
func (c *Client) read(ctx context.Context, resources string, name string, v interface{}) error {
	sql := fmt.Sprintf("SHOW %s LIKE '%s'", resources, name)
	rows, err := c.query(ctx, sql)
	if err != nil {
		return fmt.Errorf("do query: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrNoRecord
	}
	if err := rows.StructScan(v); err != nil {
		return fmt.Errorf("rows scan: %w", err)
	}
	return nil
}
