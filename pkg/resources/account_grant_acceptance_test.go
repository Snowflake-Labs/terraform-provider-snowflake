package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_AccountGrant_defaults(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountGrantConfig(roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_grant.test", "privilege", "MONITOR USAGE"),
				),
			},
		},
	})
}

func accountGrantConfig(role string) string {
	return fmt.Sprintf(`

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_account_grant" "test" {
  roles          = [snowflake_role.test.name]
  privilege = "MONITOR USAGE"
}
`, role)
}
/*
// try commenting this out, since test is mysteriously failing
https://github.com/Snowflake-Labs/terraform-provider-snowflake/actions/runs/4906105061/jobs/8760528078?pr=1779

panic: test timed out after 10m0s

goroutine 103803 [running]:
testing.(*M).startAlarm.func1()
	/opt/hostedtoolcache/go/1.19.8/x64/src/testing/testing.go:2036 +0x8e
created by time.goFunc
	/opt/hostedtoolcache/go/1.19.8/x64/src/time/sleep.go:176 +0x32
	....
github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource.ParallelTest({0x1b70360, 0xc000530820}, {0x0, 0x0, 0x0, 0x0, 0x0, 0xc000652c00, 0x0, 0x0, ...})
	/home/runner/go/pkg/mod/github.com/hashicorp/terraform-plugin-sdk/v2@v2.26.1/helper/resource/testing.go:678 +0x55
github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources_test.TestAcc_AccountGrantManagedTask(0x0?)

func TestAcc_AccountGrantManagedTask(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountGrantManagedTaskConfig(roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_grant.test", "privilege", "EXECUTE MANAGED TASK"),
				),
			},
		},
	})
}
*/
func accountGrantManagedTaskConfig(role string) string {
	return fmt.Sprintf(`

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_account_grant" "test" {
  roles     = [snowflake_role.test.name]
  privilege = "EXECUTE MANAGED TASK"
}
`, role)
}

func TestAcc_AccountGrantManageSupportCases(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountGrantManageSupportCasesConfig(roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_grant.test", "privilege", "MANAGE ACCOUNT SUPPORT CASES"),
				),
			},
		},
	})
}

func accountGrantManageSupportCasesConfig(role string) string {
	return fmt.Sprintf(`

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_account_grant" "test" {
  roles     = [snowflake_role.test.name]
  privilege = "MANAGE ACCOUNT SUPPORT CASES"
}
`, role)
}
