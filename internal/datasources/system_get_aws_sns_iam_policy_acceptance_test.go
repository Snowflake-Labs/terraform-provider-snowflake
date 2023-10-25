// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_SystemGetAWSSNSIAMPolicy_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: policyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_system_get_aws_sns_iam_policy.p", "aws_sns_topic_arn", "arn:aws:sns:us-east-1:1234567890123456:mytopic"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_aws_sns_iam_policy.p", "aws_sns_topic_policy_json"),
				),
			},
		},
	})
}

func policyConfig() string {
	s := `
	data snowflake_system_get_aws_sns_iam_policy p {
		aws_sns_topic_arn = "arn:aws:sns:us-east-1:1234567890123456:mytopic"
	}
	`
	return s
}
