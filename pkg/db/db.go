package db

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/ExpansiveWorlds/instrumentedsql"
	"github.com/snowflakedb/gosnowflake"
)

func init() {
	re := regexp.MustCompile(`\r?\n`)

	logger := instrumentedsql.LoggerFunc(func(ctx context.Context, msg string, keyvals ...interface{}) {
		s := fmt.Sprintf("[DEBUG] %s %v\n", msg, keyvals)
		fmt.Println(re.ReplaceAllString(s, " "))
	})

	sql.Register("snowflake-instrumented", instrumentedsql.WrapDriver(&gosnowflake.SnowflakeDriver{}, instrumentedsql.WithLogger(logger)))
}

func Open(dsn string) (*sql.DB, error) {
	return sql.Open("snowflake-instrumented", dsn)
}
