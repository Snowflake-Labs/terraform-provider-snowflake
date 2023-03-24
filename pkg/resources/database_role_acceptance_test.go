package resources_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type (
	TestAccDatabaseRoleSettings struct {
		WarehouseName string
		DatabaseName  string
		DatabaseRole  *DatabaseRoleSettings
	}

	DatabaseRoleSettings struct {
		Name    string
		Comment string
	}
)

var (
	resourceName = "snowflake_database_role.test_db_role"
	databaseName = "db_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	dbRoleName   = "db_role_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	dbRoleInitialState = &TestAccDatabaseRoleSettings{ //nolint
		DatabaseName: databaseName,
		DatabaseRole: &DatabaseRoleSettings{
			Name:    dbRoleName,
			Comment: "dummy",
		},
	}

	dbRoleStepOne = &TestAccDatabaseRoleSettings{ //nolint
		DatabaseName: databaseName,
		DatabaseRole: &DatabaseRoleSettings{
			Name:    dbRoleName,
			Comment: "test",
		},
	}

	dbRoleStepTwo = &TestAccDatabaseRoleSettings{ //nolint
		DatabaseName: databaseName,
		DatabaseRole: &DatabaseRoleSettings{
			Name:    dbRoleName,
			Comment: "text",
		},
	}
)

func TestAcc_DatabaseRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: databaseRoleConfig(dbRoleInitialState),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dbRoleName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
				),
			},
			{
				Config: databaseRoleConfig(dbRoleStepOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dbRoleName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "comment", dbRoleStepOne.DatabaseRole.Comment),
				),
			},
			{
				Config: databaseRoleConfig(dbRoleStepTwo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dbRoleName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "comment", dbRoleStepTwo.DatabaseRole.Comment),
				),
			},
			{
				Config: databaseRoleConfig(dbRoleInitialState),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dbRoleName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "comment", dbRoleInitialState.DatabaseRole.Comment),
				),
			},
		},
	})
}

func databaseRoleConfig(settings *TestAccDatabaseRoleSettings) string { //nolint
	config, err := template.New("db_role_acceptance_test_config").Parse(`
resource "snowflake_database" "test_db" {
	name = "{{ .DatabaseName }}"
}
resource "snowflake_database_role" "test_db_role" {
	name     	  = "{{ .DatabaseRole.Name }}"
	database  	  = snowflake_database.test_db.name
	comment       = "{{ .DatabaseRole.Comment }}"
}
	`)
	if err != nil {
		fmt.Println(err)
	}

	var result bytes.Buffer
	err = config.Execute(&result, settings) //nolint
	if err != nil {
		fmt.Println(err)
	}
	return result.String()
}
