package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"github.com/luna-duclos/instrumentedsql"
	"github.com/snowflakedb/gosnowflake"
)

func init() {
	re := regexp.MustCompile(`\r?\n`)

	logger := instrumentedsql.LoggerFunc(func(ctx context.Context, msg string, keyvals ...interface{}) {
		s := fmt.Sprintf("[DEBUG] %s %v\n", msg, keyvals)
		log.Println(re.ReplaceAllString(s, " "))
	})

	sql.Register("snowflake-instrumented", instrumentedsql.WrapDriver(&gosnowflake.SnowflakeDriver{}, instrumentedsql.WithLogger(logger)))
}

func Open(dsn string) (*sql.DB, error) {
	return sql.Open("snowflake-instrumented", dsn)
}
