package provider_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Provider_configHierarchy(t *testing.T) {
	user := os.Getenv(testenvs.User)
	pass := os.Getenv(testenvs.Password)
	account := os.Getenv(testenvs.Account)
	role := os.Getenv(testenvs.Role)
	host := os.Getenv(testenvs.Host)
	if user == "" || pass == "" || account == "" || role == "" || host == "" {
		t.Skip("Skipping TestAcc_Provider_configHierarchy")
	}

	nonExistingUser := "non-existing-user"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// make sure that we fail for incorrect profile
			{
				Config:      providerConfig(testprofiles.IncorrectUserAndPassword),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// incorrect user in provider config should not be rewritten by profile and cause error
			{
				Config:      providerConfigWithUser(nonExistingUser, testprofiles.Default),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// correct user and password in provider's config should not be rewritten by a faulty config
			{
				Config: providerConfigWithUserAndPassword(user, pass, testprofiles.IncorrectUserAndPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// incorrect user in env variable should not be rewritten by profile and cause error
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, nonExistingUser)
				},
				Config:      providerConfig(testprofiles.Default),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// correct user and password in env should not be rewritten by a faulty config
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, user)
					t.Setenv(snowflakeenvs.Password, pass)
				},
				Config: providerConfig(testprofiles.IncorrectUserAndPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// user on provider level wins (it's incorrect - env and profile ones are)
			{
				Config:      providerConfigWithUser(nonExistingUser, testprofiles.Default),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// there is no config (by setting the dir to something different than .snowflake/config)
			{
				PreConfig: func() {
					dir, err := os.UserHomeDir()
					require.NoError(t, err)
					t.Setenv(snowflakeenvs.ConfigPath, dir)
				},
				Config:      providerConfigWithUserAndPassword(user, pass, testprofiles.Default),
				ExpectError: regexp.MustCompile("account is empty"),
			},
			// provider's config should not be rewritten by env when there is no profile (incorrect user in config versus correct one in env) - proves #2242
			{
				PreConfig: func() {
					require.NotEmpty(t, os.Getenv(snowflakeenvs.ConfigPath))
					t.Setenv(snowflakeenvs.User, user)
					t.Setenv(snowflakeenvs.Password, pass)
					t.Setenv(snowflakeenvs.Account, account)
					t.Setenv(snowflakeenvs.Role, role)
					t.Setenv(snowflakeenvs.Host, host)
				},
				Config:      providerConfigWithUser(nonExistingUser, testprofiles.Default),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// make sure the teardown is fine by using a correct env config at the end
			{
				PreConfig: func() {
					require.NotEmpty(t, os.Getenv(snowflakeenvs.ConfigPath))
					require.NotEmpty(t, os.Getenv(snowflakeenvs.User))
					require.NotEmpty(t, os.Getenv(snowflakeenvs.Password))
					require.NotEmpty(t, os.Getenv(snowflakeenvs.Account))
					require.NotEmpty(t, os.Getenv(snowflakeenvs.Role))
					require.NotEmpty(t, os.Getenv(snowflakeenvs.Host))
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
