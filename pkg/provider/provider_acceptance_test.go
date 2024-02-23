package provider_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Provider_configHierarchy(t *testing.T) {
	user := os.Getenv("TEST_SF_TF_USER")
	pass := os.Getenv("TEST_SF_TF_PASSWORD")
	account := os.Getenv("TEST_SF_TF_ACCOUNT")
	role := os.Getenv("TEST_SF_TF_ROLE")
	host := os.Getenv("TEST_SF_TF_HOST")
	if user == "" || pass == "" || account == "" || role == "" || host == "" {
		t.Skip("Skipping TestAcc_Provider_configHierarchy")
	}

	nonExistingUser := "non-existing-user"
	profileWithIncorrectUserAndPassword := "incorrect_test_profile"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// make sure that we fail for incorrect profile
			{
				Config:      providerConfig(profileWithIncorrectUserAndPassword),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// incorrect user in provider config should not be rewritten by profile and cause error
			{
				Config:      providerConfigWithUser(nonExistingUser, "default"),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// correct user and password in provider's config should not be rewritten by a faulty config
			{
				Config: providerConfigWithUserAndPassword(user, pass, profileWithIncorrectUserAndPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// incorrect user in env variable should not be rewritten by profile and cause error
			{
				PreConfig: func() {
					t.Setenv("SNOWFLAKE_USER", nonExistingUser)
				},
				Config:      providerConfig("default"),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// correct user and password in env should not be rewritten by a faulty config
			{
				PreConfig: func() {
					t.Setenv("SNOWFLAKE_USER", user)
					t.Setenv("SNOWFLAKE_PASSWORD", pass)
				},
				Config: providerConfig(profileWithIncorrectUserAndPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// user on provider level wins (it's incorrect - env and profile ones are)
			{
				Config:      providerConfigWithUser(nonExistingUser, "default"),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// there is no config (by setting the dir to something different than .snowflake/config)
			{
				PreConfig: func() {
					dir, err := os.UserHomeDir()
					require.NoError(t, err)
					t.Setenv("SNOWFLAKE_CONFIG_PATH", dir)
				},
				Config:      providerConfigWithUserAndPassword(user, pass, "default"),
				ExpectError: regexp.MustCompile("account is empty"),
			},
			// provider's config should not be rewritten by env when there is no profile (incorrect user in config versus correct one in env) - proves #2242
			{
				PreConfig: func() {
					require.NotEmpty(t, os.Getenv("SNOWFLAKE_CONFIG_PATH"))
					t.Setenv("SNOWFLAKE_USER", user)
					t.Setenv("SNOWFLAKE_PASSWORD", pass)
					t.Setenv("SNOWFLAKE_ACCOUNT", account)
					t.Setenv("SNOWFLAKE_ROLE", role)
					t.Setenv("SNOWFLAKE_HOST", host)
				},
				Config:      providerConfigWithUser(nonExistingUser, "default"),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// make sure the teardown is fine by using a correct env config at the end
			{
				PreConfig: func() {
					require.NotEmpty(t, os.Getenv("SNOWFLAKE_CONFIG_PATH"))
					require.NotEmpty(t, os.Getenv("SNOWFLAKE_USER"))
					require.NotEmpty(t, os.Getenv("SNOWFLAKE_PASSWORD"))
					require.NotEmpty(t, os.Getenv("SNOWFLAKE_ACCOUNT"))
					require.NotEmpty(t, os.Getenv("SNOWFLAKE_ROLE"))
					require.NotEmpty(t, os.Getenv("SNOWFLAKE_HOST"))
				},
				Config: emptyProviderConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
		},
	})
}

func emptyProviderConfig() string {
	return `
provider "snowflake" {
}` + datasourceConfig()
}

func providerConfig(profile string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
}
`, profile) + datasourceConfig()
}

func providerConfigWithUser(user string, profile string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	user = "%[1]s"
	profile = "%[2]s"
}
`, user, profile) + datasourceConfig()
}

func providerConfigWithUserAndPassword(user string, pass string, profile string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	user = "%[1]s"
	password = "%[2]s"
	profile = "%[3]s"
}
`, user, pass, profile) + datasourceConfig()
}

func datasourceConfig() string {
	return fmt.Sprintf(`
data snowflake_database "t" {
	name = "%s"
}`, acc.TestDatabaseName)
}
