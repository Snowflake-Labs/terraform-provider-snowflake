// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package datasources_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Accounts(t *testing.T) {
	if _, ok := os.LookupEnv("SNOWFLAKE_TEST_ACCOUNTS_SHOW"); !ok {
		t.Skip("Skipping TestInt_Accounts")
	}
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_accounts.accounts", "accounts.#"),
				),
			},
		},
	})
}

func accountsConfig() string {
	return `data "snowflake_accounts" "accounts" {}`
}
