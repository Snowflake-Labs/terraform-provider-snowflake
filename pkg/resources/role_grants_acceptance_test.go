package resources_test

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func mustParseInt(t *testing.T, input string) int64 {
	t.Helper()
	i, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		t.Error(err)
	}
	return i
}

func extractList(in map[string]string, name string) ([]string, error) {
	out := make([]string, 0)
	r, err := regexp.Compile(fmt.Sprintf(`^%s.\d+$`, name))
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

func listSetEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func testCheckRolesAndUsers(t *testing.T, path string, roles, users []string) func(state *terraform.State) error {
	t.Helper()
	return func(state *terraform.State) error {
		is := state.RootModule().Resources[path].Primary
		if c, ok := is.Attributes["roles.#"]; !ok || mustParseInt(t, c) != int64(len(roles)) {
			return fmt.Errorf("expected roles.# to equal %d but got %s", len(roles), c)
		}
		r, err := extractList(is.Attributes, "roles")
		if err != nil {
			return err
		}

		// TODO no longer case sensitive
		if !listSetEqual(roles, r) {
			return fmt.Errorf("expected roles %#v but got %#v", roles, r)
		}

		if c, ok := is.Attributes["users.#"]; !ok || mustParseInt(t, c) != int64(len(users)) {
			return fmt.Errorf("expected users.# to equal %d but got %s", len(users), c)
		}
		u, err := extractList(is.Attributes, "users")
		if err != nil {
			return err
		}

		if !listSetEqual(users, u) {
			return fmt.Errorf("expected users %#v but got %#v", users, u)
		}

		return nil
	}
}

func TestAcc_RoleGrant(t *testing.T) {
	role1 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	role2 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	role3 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	user1 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	user2 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	basicChecks := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("snowflake_role.r", "name", role1),
		resource.TestCheckResourceAttr("snowflake_role.r2", "name", role2),
		resource.TestCheckResourceAttr("snowflake_role_grants.w", "role_name", role1),
	)

	baselineStep := resource.TestStep{
		Config:       rgConfig(role1, role2, role3, user1, user2),
		ResourceName: "snowflake_role_grants.w",
		Check: resource.ComposeTestCheckFunc(
			basicChecks,
			testCheckRolesAndUsers(t, "snowflake_role_grants.w", []string{role2, role3}, []string{user1, user2}),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// test setup + removing a role
			baselineStep,
			{
				Config:       rgConfig2(role1, role2, role3, user1, user2),
				ResourceName: "snowflake_role_grants.w",
				Check: resource.ComposeTestCheckFunc(
					basicChecks,
					testCheckRolesAndUsers(t, "snowflake_role_grants.w", []string{role2}, []string{user1, user2})),
			},
			// back to baseline, which means adding a role
			baselineStep,
			// then remove a user
			{
				Config:       rgConfig3(role1, role2, role3, user1, user2),
				ResourceName: "snowflake_role_grants.w",

				Check: resource.ComposeTestCheckFunc(
					basicChecks,
					testCheckRolesAndUsers(t, "snowflake_role_grants.w", []string{role2, role3}, []string{user1})),
			},
			// add the user back to get back to baseline
			baselineStep,
			// now try reordering and ensure there is no diff
			{
				Config:             rgConfig4(role1, role2, role3, user1, user2),
				ResourceName:       "snowflake_role_grants.w",
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,

				Check: resource.ComposeTestCheckFunc(
					basicChecks,
					func(state *terraform.State) error {
						return nil
					},
				),
			},
			baselineStep,
			// IMPORT
			{
				ResourceName:            "snowflake_role_grants.w",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_multiple_grants"},
			},
		},
	})
}

func rolesAndUser(role1, role2, role3, user1, user2 string) string {
	s := `
resource "snowflake_role" "r" {
	name = "%s"
}
resource "snowflake_role" "r2" {
	name = "%s"
}
resource "snowflake_role" "r3" {
	name = "%s"
}
resource "snowflake_user" "u" {
	name = "%s"
}
resource "snowflake_user" "u2" {
	name = "%s"
}
`
	return fmt.Sprintf(s, role1, role2, role3, user1, user2)
}

func rgConfig(role1, role2, role3, user1, user2 string) string {
	s := `
	%s

	resource "snowflake_role_grants" "w" {
		role_name = "${snowflake_role.r.name}"
		roles = ["${snowflake_role.r2.name}", "${snowflake_role.r3.name}"]
		users = ["${snowflake_user.u.name}", "${snowflake_user.u2.name}"]
	}
	`
	return fmt.Sprintf(s, rolesAndUser(role1, role2, role3, user1, user2))
}

func rgConfig2(role1, role2, role3, user1, user2 string) string {
	s := `

	%s

	resource "snowflake_role_grants" "w" {
		role_name = "${snowflake_role.r.name}"
		roles = ["${snowflake_role.r2.name}"]
		users = ["${snowflake_user.u.name}", "${snowflake_user.u2.name}"]
	}
	`

	return fmt.Sprintf(s, rolesAndUser(role1, role2, role3, user1, user2))
}

func rgConfig3(role1, role2, role3, user1, user2 string) string {
	s := `

	%s

	resource "snowflake_role_grants" "w" {
		role_name = "${snowflake_role.r.name}"
		roles = ["${snowflake_role.r2.name}", "${snowflake_role.r3.name}"]
		users = ["${snowflake_user.u.name}"]
	}
	`

	return fmt.Sprintf(s, rolesAndUser(role1, role2, role3, user1, user2))
}

func rgConfig4(role1, role2, role3, user1, user2 string) string {
	s := `

%s
	resource "snowflake_role_grants" "w" {
	role_name = "${snowflake_role.r.name}"
	roles = ["${snowflake_role.r3.name}", "${snowflake_role.r2.name}"]
	users = ["${snowflake_user.u2.name}", "${snowflake_user.u.name}"]
}
`
	return fmt.Sprintf(s, rolesAndUser(role1, role2, role3, user1, user2))
}
