package resources_test

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func MustParseInt(input string) int64 {
	i, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func extractList(in map[string]string, name string) ([]string, error) {
	out := make([]string, 0)
	r, err := regexp.Compile(fmt.Sprintf(`^%s.\d+(.+)$`, name))
	if err != nil {
		return out, err
	}
	for k, v := range in {
		if r.MatchString(k) {
			log.Printf("[DEBUG] matched %s %s", k, v)
			out = append(out, v)
		} else {
			log.Printf("[DEBUG] no match %s", k)
		}
	}
	return out, nil
}

func listEquiv(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if strings.ToUpper(a[i]) != strings.ToUpper(b[i]) {
			return false
		}
	}
	return true
}

func testCheckRolesAndUsers(path string, roles, users []string) func(state *terraform.State) error {

	return func(state *terraform.State) error {
		is := state.RootModule().Resources[path].Primary

		if c, ok := is.Attributes["roles.#"]; !ok || MustParseInt(c) != int64(len(roles)) {
			return fmt.Errorf("expected roles.# to equal %d but got %s", len(roles), c)
		}
		r, err := extractList(is.Attributes, "roles")
		if err != nil {
			return err
		}

		if !listEquiv(roles, r) {
			return fmt.Errorf("expected roles %#v but got %#v", roles, r)
		}

		if c, ok := is.Attributes["users.#"]; !ok || MustParseInt(c) != int64(len(users)) {
			return fmt.Errorf("expected users.# to equal %d but got %s", len(users), c)
		}
		u, err := extractList(is.Attributes, "users")
		if err != nil {
			return err
		}

		if !listEquiv(users, u) {
			return fmt.Errorf("expected users %#v but got %#v", users, u)
		}

		return nil
	}
}

func TestAccGrantRole(t *testing.T) {
	role1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	role2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: rgConfig(role1, role2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.r", "name", strings.ToUpper(role1)),
					resource.TestCheckResourceAttr("snowflake_role.r2", "name", strings.ToUpper(role2)),
					resource.TestCheckResourceAttr("snowflake_role_grants.w", "name", strings.ToUpper(role1)),
					testCheckRolesAndUsers("snowflake_role_grants.w", []string{role2}, []string{}),
				),
			},
			// 			// CHANGE PROPERTIES
			// 			{
			// 				Config: uConfig2(prefix2),
			// 				Check: resource.ComposeTestCheckFunc(
			// 					resource.TestCheckResourceAttr("snowflake_grant_role.w", "name", strings.ToUpper(prefix2)),
			// 					resource.TestCheckResourceAttr("snowflake_grant_role.w", "comment", "test comment 2"),
			// 					resource.TestCheckResourceAttr("snowflake_grant_role.w", "password", "best password"),
			// 				),
			// 			},
			// 			// IMPORT
			// 			{
			// 				ResourceName:            "snowflake_grant_role.w",
			// 				ImportState:             true,
			// 				ImportStateVerify:       true,
			// 			},
		},
	})
}

func rgConfig(prefix, prefix2 string) string {
	s := `
resource "snowflake_role" "r" {
	name = "%s"
}
resource "snowflake_role" "r2" {
	name = "%s"
}
resource "snowflake_role_grants" "w" {
	name = "${snowflake_role.r.name}"
	role_name = "${snowflake_role.r.name}"
	roles = ["${snowflake_role.r2.name}"]
}
`
	return fmt.Sprintf(s, prefix, prefix2)
}

func rgConfig2(prefix string) string {
	s := `
resource "snowflake_grant_role" "w" {
	name = "%s"
	comment = "test comment 2"
	password = "best password"
}
`
	return fmt.Sprintf(s, prefix)
}
