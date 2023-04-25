package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

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

	PasswordPolicies PasswordPolicies
}

func NewDefaultClient() (*Client, error) {
	return NewClient(nil)
}

func NewClient(cfg *gosnowflake.Config) (*Client, error) {
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

	dsn, err := gosnowflake.DSN(config)
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

func (c *Client) execContext(ctx context.Context, sql string) (sql.Result, error) {
	if !c.dryRun {
		return c.db.ExecContext(ctx, sql)
	}
	return nil, nil
}

func (c *Client) selectContext(ctx context.Context, dest interface{}, query string) error {
	if !c.dryRun {
		return c.db.SelectContext(ctx, dest, query)
	}
	return nil
}

/*
func (c *Client) getContext(ctx context.Context, dest interface{}, query string) error {
	if !c.dryRun {
		return c.db.GetContext(ctx, dest, query)
	}
	return nil
}*/

type SQLOperation string

const (
	sqlOperationCreate   SQLOperation = "CREATE"
	sqlOperationAlter    SQLOperation = "ALTER"
	sqlOperationDrop     SQLOperation = "DROP"
	sqlOperationShow     SQLOperation = "SHOW"
	sqlOperationDescribe SQLOperation = "DESCRIBE"
)

func (v SQLOperation) String() string {
	return string(v)
}

func (c *Client) sql(sqlOperation SQLOperation, clause ...ddlClause) string {
	sb := strings.Builder{}
	sb.WriteString(sqlOperation.String())
	for _, c := range clause {
		sb.WriteString(fmt.Sprintf(" %s", c.String()))
	}
	return sb.String()
}

// misc helper functions.
func pSlice[T any](a []T) []*T {
	b := make([]*T, len(a))
	for i := range a {
		b = append(b, &a[i])
	}
	return b
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func getUnexportedClause(field reflect.StructField, value reflect.Value) ddlClause {
	if field.Tag.Get("ddl") == "" {
		return nil
	}
	tagParts := strings.Split(field.Tag.Get("ddl"), ",")
	ddlType := tagParts[0]
	dbTag := field.Tag.Get("db")
	switch ddlType {
	case "object_type":
		return ddlClauseObjectType(ObjectType(value.String()))
	case "name":
		return ddlClauseName(value.String())
	case "keyword":
		if value.Kind() == reflect.Bool {
			useKeyword := value.Bool()
			if useKeyword {
				return ddlClauseKeyword(dbTag)
			}
		}
		return ddlClauseKeyword(value.String())
	}

	return nil
}

func getExportedClause(field reflect.StructField, value reflect.Value) ddlClause {
	if field.Tag.Get("ddl") == "" {
		return nil
	}
	if field.Tag.Get("db") == "" {
		return nil
	}

	ddlTag := strings.Split(field.Tag.Get("ddl"), ",")[0]
	dbTag := field.Tag.Get("db")

	var clause ddlClause
	switch ddlTag {
	case "keyword":
		if value.Kind() == reflect.Bool {
			useKeyword := value.Interface().(bool)
			if useKeyword {
				clause = ddlClauseKeyword(dbTag)
			}
		} else {
			clause = ddlClauseKeyword(value.Interface().(string))
		}
	case "command_param":
		clause = ddlClauseCommandParameter{
			key:   dbTag,
			value: value.Interface().(string),
			qt:    getQuoteTypeFromTag(field.Tag.Get("ddl")),
		}
	case "param":
		clause = ddlClauseParameter{
			key:   dbTag,
			value: value.Interface(),
			qt:    getQuoteTypeFromTag(field.Tag.Get("ddl")),
		}
	case "list":
		clause = ddlClauseList{
			key:   dbTag,
			value: value.Interface().([]string),
		}
	default:
		return nil
	}
	return clause
}

func ddlClausesForObject(s interface{}) ([]ddlClause, error) {
	clauses := make([]ddlClause, 0)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", v.Kind())
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// unexported fields need to be handled separately.
		if !value.CanInterface() {
			clause := getUnexportedClause(field, value)
			if clause != nil {
				clauses = append(clauses, clause)
			}
			continue
		}

		if value.Kind() == reflect.Ptr {
			// skip nil pointers for attributes, since they are not set.
			if value.IsNil() {
				continue
			}
			value = value.Elem()
		}

		if value.Kind() == reflect.Struct {
			// check if there is any keyword on the struct
			// if there is, then we need to add it to the clause
			// if there is not, then we need to recurse into the struct
			// and get the clauses from there
			ddlTag := field.Tag.Get("ddl")
			if ddlTag != "" {
				ddlTagParts := strings.Split(ddlTag, ",")
				if ddlTagParts[0] == "keyword" {
					clauses = append(clauses, ddlClauseKeyword(field.Tag.Get("db")))
				}
			}
			innerClauses, err := ddlClausesForObject(value.Interface())
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, innerClauses...)
			continue
		}
		clause := getExportedClause(field, value)
		if clause != nil {
			clauses = append(clauses, clause)
		}
	}
	return clauses, nil
}

func getQuoteTypeFromTag(t string) quoteType {
	parts := strings.Split(t, ",")
	for _, part := range parts {
		if strings.Contains(part, "quote") {
			return quoteType(part)
		}
	}
	return NoQuotes
}

type ddlClause interface {
	String() string
}

type ddlClauseKeyword string

func (v ddlClauseKeyword) String() string {
	return string(v)
}

type ddlClauseObjectType ObjectType

func (v ddlClauseObjectType) String() string {
	return string(v)
}

type ddlClauseName string

func (v ddlClauseName) String() string {
	return string(v)
}

type ddlClauseParameter struct {
	key   string
	value interface{} // string list, string, string literal, bool, int
	qt    quoteType
}

func (v ddlClauseParameter) String() string {
	vType := reflect.TypeOf(v.value)
	var result string
	if v.key != "" {
		result = fmt.Sprintf("%s = ", v.key)
	}
	if vType.Kind() == reflect.String {
		result += fmt.Sprintf("%s%s%s", v.qt.String(), v.value, v.qt.String())
	} else {
		result += fmt.Sprintf("%v", v.value)
	}

	return result
}

type ddlClauseCommandParameter struct {
	key   string
	value string
	qt    quoteType
}

func (v ddlClauseCommandParameter) String() string {
	return fmt.Sprintf("%s %s%s%s", v.key, v.qt.String(), v.value, v.qt.String())
}

type ddlClauseList struct {
	key   string
	value []string
}

func (v ddlClauseList) String() string {
	return fmt.Sprintf("%s = %s", v.key, strings.Join(v.value, ","))
}

type quoteType string

const (
	NoQuotes     quoteType = "no_quotes"
	DoubleQuotes quoteType = "double_quotes"
	SingleQuotes quoteType = "single_quotes"
)

func (v quoteType) String() string {
	switch v {
	case NoQuotes:
		return ""
	case DoubleQuotes:
		return "\""
	case SingleQuotes:
		return "'"
	}
	return ""
}
