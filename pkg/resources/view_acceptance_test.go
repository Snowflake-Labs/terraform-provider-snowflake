package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_View(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, false, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

func TestAcc_View2(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, false, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES where ROLE_OWNER like 'foo%%';"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

func TestAcc_ViewWithCopyGrants(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, true, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "true"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

func TestAcc_ViewWithChangeTracking(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewConfigCreateViewWithChangeTracking(accName, fmt.Sprintf("SELECT * FROM %s.%s", accName, accName)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", fmt.Sprintf("%s_vw", accName)),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					resource.TestCheckResourceAttr("snowflake_view.test", "change_tracking", "true"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

func viewConfig(n string, copyGrants bool, q string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%v"
}

resource "snowflake_view" "test" {
	name        = "%v"
	comment     = "Terraform test resource"
	database    = snowflake_database.test.name
	schema      = "PUBLIC"
	is_secure   = true
	or_replace  = %t
	copy_grants = %t
	statement   = "%s"
}
`, n, n, copyGrants, copyGrants, q)
}

func viewConfigCreateViewWithChangeTracking(n string, q string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%v"
}

resource "snowflake_schema" "test" {
	name     = "%v"
	database = snowflake_database.test.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test" {
	database = snowflake_database.test.name
	schema   = snowflake_schema.test.name
	name     = "%v"
	comment  = "Terraform acceptance test"
	change_tracking = true
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
		nullable = false
	}
	column {
		name = "column3"
		type = "FLOAT"
		nullable = false
	}
	primary_key {
		name = "new_name"
		keys = ["column2","column3"]
	}
}

resource "snowflake_view" "test" {
	name        = "%v_vw"
	comment     = "Terraform test resource"
	database    = snowflake_database.test.name
	schema      = snowflake_schema.test.name
	is_secure   = true
	or_replace  = false
	copy_grants = false
	change_tracking = true
	statement   = "%s"

	depends_on = [snowflake_table.test]
}
`, n, n, n, n, q)
}
