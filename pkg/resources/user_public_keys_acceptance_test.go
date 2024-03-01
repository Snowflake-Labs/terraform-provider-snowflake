package resources_test

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_UserPublicKeys(t *testing.T) {
	r := require.New(t)
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	sshkey1, err := testhelpers.Fixture("userkey1")
	r.NoError(err)
	sshkey2, err := testhelpers.Fixture("userkey2")
	r.NoError(err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: uPublicKeysConfig(r, PublicKeyData{
					Prefix:     prefix,
					PublicKey1: sshkey1,
					PublicKey2: sshkey2,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix),

					resource.TestCheckResourceAttr("snowflake_user_public_keys.foobar", "rsa_public_key", sshkey1),
					resource.TestCheckResourceAttr("snowflake_user_public_keys.foobar", "rsa_public_key_2", sshkey2),
					resource.TestCheckNoResourceAttr("snowflake_user_public_keys.foobar", "has_rsa_public_key"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_user.w",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignoring because keys are currently altered outside of snowflake_user resource (in snowflake_user_public_keys).
				ImportStateVerifyIgnore: []string{"password", "rsa_public_key", "rsa_public_key_2", "has_rsa_public_key", "must_change_password"},
			},
		},
	})
}

type PublicKeyData struct {
	Prefix     string
	PublicKey1 string
	PublicKey2 string
}

func uPublicKeysConfig(r *require.Assertions, data PublicKeyData) string {
	t := `
resource "snowflake_user" "w" {
	name = "{{.Prefix}}"
	comment = "test comment"
	login_name = "{{.Prefix}}_login"
	display_name = "Display Name"
	first_name = "Marcin"
	last_name = "Zukowski"
	email = "fake@email.com"
	disabled = false
	default_warehouse="foo"
	default_role="foo"
	default_namespace="foo"
}

resource "snowflake_user_public_keys" "foobar" {
	name = snowflake_user.w.name
	rsa_public_key = <<KEY
{{ .PublicKey1 }}
	KEY

	rsa_public_key_2 = <<KEY
{{ .PublicKey2 }}
	KEY
}
`
	conf := bytes.NewBuffer(nil)
	err := template.Must(template.New("user").Parse(t)).Execute(conf, data)
	r.NoError(err)
	return conf.String()
}
