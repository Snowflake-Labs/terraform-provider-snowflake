package resources_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TableConstraint_fk(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConstraintFKConfig(name, acc.TestDatabaseName, acc.TestSchemaName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "type", "FOREIGN KEY"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "enforced", "false"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "deferrable", "false"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "comment", "hello fk"),
				),
			},
		},
	})
}

func tableConstraintFKConfig(n string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	name     = "%s"
	database = "%s"
	schema   = "%s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table" "fk_t" {
	name     = "fk_%s"
	database = "%s"
	schema   = "%s"
	column {
		name     = "fk_col1"
		type     = "text"
		nullable = false
	  }
}

resource "snowflake_table_constraint" "fk" {
	name="%s"
	type= "FOREIGN KEY"
	table_id = snowflake_table.t.qualified_name
	columns = ["col1"]
	foreign_key_properties {
	  references {
		table_id = snowflake_table.fk_t.qualified_name
		columns = ["fk_col1"]
	  }
	}
	enforced = false
	deferrable = false
	initially = "IMMEDIATE"
	comment = "hello fk"
}

`, n, databaseName, schemaName, n, databaseName, schemaName, n)
}

// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2674
// It is connected with https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2629.
// Provider defaults will be reworked during resources redesign.
func TestAcc_TableConstraint_pk(t *testing.T) {
	tableName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	constraintName := fmt.Sprintf("%s_pk", tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConstraintPKConfig(acc.TestDatabaseName, acc.TestSchemaName, tableName, constraintName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.pk", "type", "PRIMARY KEY"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.pk", "comment", "hello pk"),
					checkPrimaryKeyExists(sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, tableName), constraintName),
				),
			},
		},
	})
}

func tableConstraintPKConfig(databaseName string, schemaName string, tableName string, constraintName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table_constraint" "pk" {
	name = "%[4]s"
	type = "PRIMARY KEY"
	table_id = snowflake_table.t.qualified_name
	columns = ["col1"]
	enable = false
	deferrable = false
	comment = "hello pk"
}
`, databaseName, schemaName, tableName, constraintName)
}

type PrimaryKeys struct {
	ConstraintName string `db:"constraint_name"`
}

func checkPrimaryKeyExists(tableId sdk.SchemaObjectIdentifier, constraintName string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		var keys []PrimaryKeys
		err := client.QueryForTests(ctx, &keys, fmt.Sprintf("show primary keys in %s", tableId.FullyQualifiedName()))
		if err != nil {
			return err
		}

		var found bool
		for _, pk := range keys {
			if pk.ConstraintName == strings.ToUpper(constraintName) {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("unable to find primary key %s on table %s, found: %v", constraintName, tableId.FullyQualifiedName(), keys)
		}

		return nil
	}
}

func TestAcc_TableConstraint_unique(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConstraintUniqueConfig(name, acc.TestDatabaseName, acc.TestSchemaName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "type", "UNIQUE"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "enforced", "true"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "deferrable", "false"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "comment", "hello unique"),
				),
			},
		},
	})
}

func tableConstraintUniqueConfig(n string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	name     = "%s"
	database = "%s"
	schema   = "%s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table_constraint" "unique" {
	name="%s"
	type= "UNIQUE"
	table_id = snowflake_table.t.qualified_name
	columns = ["col1"]
	enforced = true
	deferrable = false
	initially = "IMMEDIATE"
	comment = "hello unique"
}

`, n, databaseName, schemaName, n)
}

// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2535
func TestAcc_Table_issue2535_newConstraint(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.86.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config:      tableConstraintUniqueConfigUsingTableId(accName, acc.TestDatabaseName, acc.TestSchemaName),
				ExpectError: regexp.MustCompile(`.*table id is incorrect.*`),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConstraintUniqueConfigUsingTableId(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "type", "UNIQUE"),
				),
			},
		},
	})
}

// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2535
func TestAcc_Table_issue2535_existingTable(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// reference done by table.id in 0.85.0
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.85.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: tableConstraintUniqueConfigUsingTableId(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "type", "UNIQUE"),
				),
			},
			// switched to qualified_name in 0.86.0
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.86.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config:      tableConstraintUniqueConfigUsingFullyQualifiedName(accName, acc.TestDatabaseName, acc.TestSchemaName),
				ExpectError: regexp.MustCompile(`.*table id is incorrect.*`),
			},
			// fixed in the current version
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConstraintUniqueConfigUsingFullyQualifiedName(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "type", "UNIQUE"),
				),
			},
		},
	})
}

func tableConstraintUniqueConfigUsingFullyQualifiedName(n string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	name     = "%s"
	database = "%s"
	schema   = "%s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table_constraint" "unique" {
	name     = "%s"
	type     = "UNIQUE"
	table_id = snowflake_table.t.qualified_name
	columns  = ["col1"]
}
`, n, databaseName, schemaName, n)
}

func tableConstraintUniqueConfigUsingTableId(n string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	name     = "%s"
	database = "%s"
	schema   = "%s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table_constraint" "unique" {
	name     = "%s"
	type     = "UNIQUE"
	table_id = snowflake_table.t.id
	columns  = ["col1"]
}
`, n, databaseName, schemaName, n)
}
