package snowflake

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/pkg/errors"
	"github.com/snowflakedb/gosnowflake"
)

var ErrNoRecord = errors.New("record not found")

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

	dryRun           bool
	PasswordPolicies PasswordPolicies
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

	/*client.Users = &users{client: client}
	client.Roles = &roles{client: client}
	client.Warehouses = &warehouses{client: client}
	client.Databases = &databases{client: client}
	client.Schemas = &schemas{client: client}
	client.Tables = &tables{client: client}
	client.NetworkPolicies = &networkPolicies{client: client}*/
	client.PasswordPolicies = &passwordPolicies{client: client}

	return client, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) exec(ctx context.Context, sql string) (sql.Result, error) {
	if !c.dryRun {
		return c.conn.ExecContext(ctx, sql)
	}
	return nil, nil
}

func (c *Client) query(ctx context.Context, sql string) (*sqlx.Rows, error) {
	return sqlx.NewDb(c.conn, "snowflake-instrumented").Unsafe().QueryxContext(ctx, sql)
}

type createOptions struct {
	objectType ObjectType
	name       string

	// generic create DDL options
	// orReplace   bool
	// ifNotExists bool
	// comment string

	// object specific DDL options
	// ddlPreKeyword []ddlKeyword
	ddlProperties []ddlProperty
}

type ddlProperty interface {
	String() string
}

type parameterDDLProperty struct {
	key   string
	value interface{} // string list, string, string literal, bool, int
}

func (p *parameterDDLProperty) String() string {
	// case if parameter is actually a keyword
	if p.value == nil {
		return p.key
	}
	return fmt.Sprintf("%s = %v", p.key, p.value)
}

type keywordPosition int

const (
	keywordPositionBefore keywordPosition = iota
	keywordPositionAfter  keywordPosition = iota
)

type SQLOperation string

const (
	sqlOperationCreate SQLOperation = "CREATE"
	sqlOperationAlter  SQLOperation = "ALTER"
	sqlOperationDrop   SQLOperation = "DROP"
)

func (v SQLOperation) String() string {
	return string(v)
}

func (c *Client) sql(sqlOperation SQLOperation, properties []ddlProperty) string {
	sb := strings.Builder{}
	sb.WriteString(sqlOperation.String())
	for _, p := range properties {
		sb.WriteString(fmt.Sprintf(" %s", p.String()))
	}
	return sb.String()
}
