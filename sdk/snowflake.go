package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/pkg/errors"
	"github.com/snowflakedb/gosnowflake"
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

	Users Users
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
		// If host is set trust it and do not use the region value
		if cfg.Host != "" {
			config.Region = ""
			config.Host = cfg.Host
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
		return nil, errors.Wrap(err, "could not build dsn for snowflake connection")
	}

	logger := instrumentedsql.LoggerFunc(func(ctx context.Context, msg string, keyvals ...interface{}) {
		s := fmt.Sprintf("[DEBUG] %s %v\n", msg, keyvals)
		re := regexp.MustCompile(`\r?\n`)
		log.Println(re.ReplaceAllString(s, " "))
	})

	sql.Register("snowflake-instrumented", instrumentedsql.WrapDriver(&gosnowflake.SnowflakeDriver{}, instrumentedsql.WithLogger(logger)))

	conn, err := sql.Open("snowflake-instrumented", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open snowflake database.")
	}

	client := &Client{
		conn: conn,
	}
	client.Users = &users{client: client}
	return client, nil

}
func (c *Client) Exec(query string, args ...interface{}) (sql.Result, error) {
	log.Print("[DEBUG] exec stmt ", query)
	return c.conn.Exec(query, args...)
}

func (c *Client) ExecMulti(queries []string) error {
	log.Print("[DEBUG] exec stmts ", queries)

	tx, err := c.conn.Begin()
	if err != nil {
		return err
	}

	for _, query := range queries {
		_, err = tx.Exec(query)
		if err != nil {
			return tx.Rollback()
		}
	}
	return tx.Commit()
}

func (c *Client) QueryRow(stmt string) *sqlx.Row {
	log.Print("[DEBUG] query stmt ", stmt)
	sdb := sqlx.NewDb(c.conn, "snowflake").Unsafe()
	return sdb.QueryRowx(stmt)
}

func (c *Client) Query(stmt string) (*sqlx.Rows, error) {
	log.Print("[DEBUG] query stmt ", stmt)
	sdb := sqlx.NewDb(c.conn, "snowflake").Unsafe()
	return sdb.Queryx(stmt)
}
