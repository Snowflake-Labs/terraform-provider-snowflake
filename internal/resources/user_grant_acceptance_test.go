// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserGrant(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_USER_GRANT_TESTS"); ok {
		t.Skip("Skipping TestAccUserGrant")
	}
	wName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: userGrantConfig(wName, roleName, "MONITOR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_grant.test", "user_name", wName),
					resource.TestCheckResourceAttr("snowflake_user_grant.test", "privilege", "MONITOR"),
				),
			},
			// UPDATE
			{
				Config: userGrantConfig(wName, roleName, "ALL PRIVILEGES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_grant.test", "user_name", wName),
					resource.TestCheckResourceAttr("snowflake_user_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_user_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func userGrantConfig(n, role, privilege string) string {
	return fmt.Sprintf(`

resource "snowflake_user" "test" {
  name      = "%v"
}

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_user_grant" "test" {
  user_name = snowflake_user.test.name
  roles     = [snowflake_role.test.name]
  privilege = "%s"
}
`, n, role, privilege)
}
