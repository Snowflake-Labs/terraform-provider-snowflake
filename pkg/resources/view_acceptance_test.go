package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_View(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
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
		Providers:    acc.TestAccProviders(),
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
		Providers:    acc.TestAccProviders(),
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

// Checks that copy_grants changes don't trigger a drop
func TestAcc_ViewChangeCopyGrants(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	var createdOn string

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, false, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "created_on", func(value string) error {
						createdOn = value
						return nil
					}),
					checkBool("snowflake_view.test", "is_secure", true),
				),
			},
			{
				Config: viewConfig(accName, true, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("snowflake_view.test", "created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("View was recreated")
						}
						return nil
					}),
					checkBool("snowflake_view.test", "is_secure", true),
				),
			},
		},
	})
}

func TestAcc_ViewChangeCopyGrantsReversed(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	var createdOn string

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, true, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "true"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "created_on", func(value string) error {
						createdOn = value
						return nil
					}),
					checkBool("snowflake_view.test", "is_secure", true),
				),
			},
			{
				Config: viewConfig(accName, false, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("snowflake_view.test", "created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("View was recreated")
						}
						return nil
					}),
					checkBool("snowflake_view.test", "is_secure", true),
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
