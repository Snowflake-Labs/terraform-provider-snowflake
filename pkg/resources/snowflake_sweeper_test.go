package resources_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jmoiron/sqlx"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func getWarehousesSweeper(name string) *resource.Sweeper {
	return &resource.Sweeper{
		Name: name,
		F: func(ununsed string) error {
			db, err := provider.GetDatabaseHandleFromEnv()
			if err != nil {
				return fmt.Errorf("Error getting db handle: %w", err)
			}

			warehouses, err := snowflake.ListWarehouses(db)
			if err != nil {
				return fmt.Errorf("Error listing warehouses: %w", err)
			}

			for _, wh := range warehouses {
				log.Printf("[DEBUG] Testing if warehouse %s starts with tst-terraform", wh.Name)
				if strings.HasPrefix(wh.Name, "tst-terraform") {
					log.Printf("[DEBUG] deleting warehouse %s", wh.Name)
					whBuilder := snowflake.Warehouse(name).Builder
					stmt := whBuilder.Drop()
					err = snowflake.Exec(db, stmt)
					if err != nil {
						return fmt.Errorf("Error deleting warehouse %q %w", wh.Name, err)
					}
				}
			}
			return nil
		},
	}
}

func getDatabaseSweepers(name string) *resource.Sweeper {
	return &resource.Sweeper{
		Name: name,
		F: func(ununsed string) error {
			db, err := provider.GetDatabaseHandleFromEnv()
			if err != nil {
				return fmt.Errorf("Error getting db handle: %w", err)
			}
			dbx := sqlx.NewDb(db, "snowflake")
			databases, err := snowflake.ListDatabases(dbx)
			if err != nil {
				return fmt.Errorf("Error listing databases: %w", err)
			}

			for _, database := range databases {
				log.Printf("[DEBUG] Testing if database %s starts with tst-terraform", database.DBName.String)
				if strings.HasPrefix(database.DBName.String, "tst-terraform") {
					log.Printf("[DEBUG] deleting database %s", database.DBName.String)
					stmt := snowflake.Database(database.DBName.String).Drop()
					err = snowflake.Exec(db, stmt)
					if err != nil {
						return fmt.Errorf("Error deleting database %q %w", database.DBName.String, err)
					}
				}
			}
			return nil
		},
	}
}

func getRolesSweeper(name string) *resource.Sweeper {
	return &resource.Sweeper{
		Name: name,
		F: func(ununsed string) error {
			db, err := provider.GetDatabaseHandleFromEnv()
			if err != nil {
				return fmt.Errorf("Error getting db handle: %w", err)
			}

			roles, err := snowflake.ListRoles(db)
			if err != nil {
				return fmt.Errorf("Error listing roles: %w", err)
			}

			for _, role := range roles {
				log.Printf("[DEBUG] Testing if role %s starts with tst-terraform", role.Name.String)
				if strings.HasPrefix(role.Name.String, "tst-terraform") {
					log.Printf("[DEBUG] deleting role %s", role.Name.String)
					stmt := snowflake.Role(role.Name.String).Drop()
					err = snowflake.Exec(db, stmt)
					if err != nil {
						return fmt.Errorf("Error deleting role %q %w", role.Name.String, err)
					}
				}
			}
			return nil
		},
	}
}

func getUsersSweeper(name string) *resource.Sweeper {
	return &resource.Sweeper{
		Name: name,
		F: func(ununsed string) error {
			db, err := provider.GetDatabaseHandleFromEnv()
			if err != nil {
				return fmt.Errorf("Error getting db handle: %w", err)
			}

			users, err := snowflake.ListUsers("*", db)
			if err != nil {
				return fmt.Errorf("Error listing users: %w", err)
			}

			for _, user := range users {
				log.Printf("[DEBUG] Testing if user %s starts with tst-terraform", user.Name.String)
				if strings.HasPrefix(user.Name.String, "tst-terraform") {
					log.Printf("[DEBUG] deleting user %s", user.Name.String)
					stmt := snowflake.User(user.Name.String).Drop()
					err = snowflake.Exec(db, stmt)
					if err != nil {
						return fmt.Errorf("Error deleting user %q %w", user.Name.String, err)
					}
				}
			}
			return nil
		},
	}
}

// Sweepers usually go along with the tests. In TF[CE]'s case everything depends on the organization,
// which means that if we delete it then all the other entities will  be deleted automatically.
func init() {
	resource.AddTestSweepers("wh_sweeper", getWarehousesSweeper("wh_sweeper"))
	resource.AddTestSweepers("db_sweeper", getDatabaseSweepers("db_sweeper"))
	resource.AddTestSweepers("role_sweeper", getRolesSweeper("role_sweeper"))
	resource.AddTestSweepers("user_sweeper", getUsersSweeper("user_sweeper"))
}
