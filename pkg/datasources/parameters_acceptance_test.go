package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_ParametersOnAccount(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: parametersConfigOnAccount(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "pattern", "AUTOCOMMIT"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "parameters.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "parameters.0.key", "AUTOCOMMIT"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.value"),
				),
			},
		},
	})
}

func TestAcc_ParametersOnSession(t *testing.T) {
	userName := "TEST_USER_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: parametersConfigOnSession(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.key"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "user", userName),
				),
			},
		},
	})
}

func TestAcc_ParametersOnObject(t *testing.T) {
	dbName := "TEST_DB_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: parametersConfigOnObject(dbName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.key"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "object_type", "DATABASE"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "object_name", dbName),
				),
			},
		},
	})
}

func parametersConfigOnAccount() string {
	return `data "snowflake_parameters" "p" {
		parameter_type = "ACCOUNT"
		pattern = "AUTOCOMMIT"
	}`
}

func parametersConfigOnSession(user string) string {
	s := `
	resource "snowflake_user" "u" {
		name = "%s"
	}

	data "snowflake_parameters" "p" {
		parameter_type = "SESSION"
		user = snowflake_user.u.name
	}`
	return fmt.Sprintf(s, user)
}

func parametersConfigOnObject(name string) string {
	stmt := `
	resource "snowflake_database" "d" {
		name = "%s"
	}
	data "snowflake_parameters" "p" {
		parameter_type = "OBJECT"
		object_type = "DATABASE"
		object_name = snowflake_database.d.name
	}`
	return fmt.Sprintf(stmt, name)
}
