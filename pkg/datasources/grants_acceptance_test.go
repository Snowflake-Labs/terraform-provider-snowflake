package datasources_test

import (
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func checkAtLeastOneGrantPresent() resource.TestCheckFunc {
	datasourceName := "data.snowflake_grants.test"
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(datasourceName, "grants.#"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.created_on"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.privilege"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.granted_on"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.name"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.granted_to"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.grantee_name"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.grant_option"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.granted_by"),
	)
}

// TODO: tests (examples from the correct ones):
// + on - account
// + on - account object
// + on - db object
// + on - schema object
// + on - invalid config - no attribute
// + on - invalid config - missing object type or name
// - to - application
// - to - application role
// - to - role
// - to - user
// - to - share
// - to - share with application package
// - to - invalid config - no attribute
// - to - invalid config - share name missing
// - of - role
// - of - application role
// - of - share
// - of - invalid config - no attribute
// - future in - database
// - future in - schema (both db and sc present)
// - future in - schema (only sc present)
// - future in - invalid config - no attribute
// - future in - invalid config - schema with no schema name
// - future to - role
// - future to - database role
// - future to - invalid config - no attribute
// - future to - invalid config - database role id invalid
func TestAcc_Grants_On_Account(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/Account"),
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_AccountObject(t *testing.T) {
	configVariables := config.Variables{
		"database": config.StringVariable(acc.TestDatabaseName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/AccountObject"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_DatabaseObject(t *testing.T) {
	configVariables := config.Variables{
		"database": config.StringVariable(acc.TestDatabaseName),
		"schema":   config.StringVariable(acc.TestSchemaName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/DatabaseObject"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_SchemaObject(t *testing.T) {
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"database": config.StringVariable(acc.TestDatabaseName),
		"schema":   config.StringVariable(acc.TestSchemaName),
		"table":    config.StringVariable(tableName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/SchemaObject"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_On_Invalid_MissingObjectType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/Invalid/MissingObjectType"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Missing required argument"),
			},
		},
	})
}
