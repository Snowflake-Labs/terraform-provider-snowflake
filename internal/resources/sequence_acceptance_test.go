// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Sequence(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accRename := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: sequenceConfig(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", accName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", fmt.Sprintf(`"%v"."%v".%v`, acc.TestDatabaseName, acc.TestSchemaName, accName)),
				),
			},
			// Set comment and rename
			{
				Config: sequenceConfigWithComment(accRename, "look at me I am a comment", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", accRename),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "comment", "look at me I am a comment"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", fmt.Sprintf(`"%v"."%v".%v`, acc.TestDatabaseName, acc.TestSchemaName, accRename)),
				),
			},
			// Unset comment and set increment
			{
				Config: sequenceConfigWithIncrement(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", accName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "increment", "32"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", fmt.Sprintf(`"%v"."%v".%v`, acc.TestDatabaseName, acc.TestSchemaName, accName)),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_sequence.test_sequence",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func sequenceConfigWithIncrement(sequenceName string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_sequence" "test_sequence" {
	name       = "%s"
	database   = "%s"
	schema     = "%s"
    increment = 32
}
`
	return fmt.Sprintf(s, sequenceName, databaseName, schemaName)
}

func sequenceConfig(sequenceName string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_sequence" "test_sequence" {
	name     = "%s"
	database   = "%s"
	schema     = "%s"
}
`
	return fmt.Sprintf(s, sequenceName, databaseName, schemaName)
}

func sequenceConfigWithComment(sequenceName, comment string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_sequence" "test_sequence" {
	name     = "%s"
	database   = "%s"
	schema     = "%s"
    comment  = "%s"
}
`
	return fmt.Sprintf(s, sequenceName, databaseName, schemaName, comment)
}
