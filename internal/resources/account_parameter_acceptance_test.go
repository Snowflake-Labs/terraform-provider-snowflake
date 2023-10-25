// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_AccountParameter(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountParameterBasic("ALLOW_ID_TOKEN", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "key", "ALLOW_ID_TOKEN"),
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "value", "true"),
				),
			},
		},
	})
}

func accountParameterBasic(key, value string) string {
	s := `
resource "snowflake_account_parameter" "p" {
	key = "%s"
	value = "%s"
}
`
	return fmt.Sprintf(s, key, value)
}

func TestAcc_AccountParameter_PREVENT_LOAD_FROM_INLINE_URL(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountParameterBasic("PREVENT_LOAD_FROM_INLINE_URL", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "key", "PREVENT_LOAD_FROM_INLINE_URL"),
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "value", "true"),
				),
			},
		},
	})
}

func TestAcc_AccountParameter_REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountParameterBasic("REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "key", "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"),
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "value", "true"),
				),
			},
		},
	})
}
