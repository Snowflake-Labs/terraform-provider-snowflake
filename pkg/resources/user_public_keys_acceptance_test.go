package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [SNOW-1348101 - this PR]: change description of user public keys resource (should be used only if user is not managed by terraform)
func TestAcc_UserPublicKeys(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	key1, _ := random.GenerateRSAPublicKey(t)
	key2, _ := random.GenerateRSAPublicKey(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, userCleanup := acc.TestClient().User.CreateUserWithOptions(t, userId, nil)
					t.Cleanup(userCleanup)
				},
				Config: uPublicKeysConfig(userId, key1, key2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_public_keys.foobar", "rsa_public_key", key1),
					resource.TestCheckResourceAttr("snowflake_user_public_keys.foobar", "rsa_public_key_2", key2),
					resource.TestCheckNoResourceAttr("snowflake_user_public_keys.foobar", "has_rsa_public_key"),
				),
			},
		},
	})
}

func uPublicKeysConfig(userId sdk.AccountObjectIdentifier, key1 string, key2 string) string {
	return fmt.Sprintf(`
resource "snowflake_user_public_keys" "foobar" {
	name = %s
	rsa_public_key = <<KEY
%s
	KEY

	rsa_public_key_2 = <<KEY
%s
	KEY
}
`, userId.FullyQualifiedName(), key1, key2)
}
