package auth
/*
import (
	"database/sql"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/db"
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

func SQLConnection(config Config) (*sql.DB, error) {
	// us-west-2 is Snowflake's default region, but if you actually specify that it won't trigger the default code
	//  https://github.com/snowflakedb/gosnowflake/blob/52137ce8c32eaf93b0bd22fc5c7297beff339812/dsn.go#L61
	if config.Region == "us-west-2" {
		config.Region = ""
	}

	snowflakeConfig := gosnowflake.Config{
		Account:   config.Account,
		User:      config.User,
		Password:  config.Password,
		Region:    config.Region,
		Role:      config.Role,
		Warehouse: config.Warehouse,
	}

	dsn := gosnowflake.DSN(&snowflakeConfig)
	db, err := db.Open(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open snowflake database.")
	}
}
*/
