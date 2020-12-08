package datasources_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FIXME refactor to testhelpers
func providers() map[string]*schema.Provider {
	p := provider.Provider()
	return map[string]*schema.Provider{
		"snowflake": p,
	}
}

func TestAccSystemGetAWSSNSIAMPolicy_basic(t *testing.T) {
	// r := require.New(t)

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: policyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_system_get_aws_sns_iam_policy.p", "aws_sns_topic_arn", "arn:aws:sns:us-east-1:1234567890123456:mytopic"),
					resource.TestCheckResourceAttr("data.snowflake_system_get_aws_sns_iam_policy.p", "aws_sns_topic_policy_json",
						`{"Version":"2012-10-17","Statement":[{"Sid":"1","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::494544507972:user/3gmi-s-ssca3411"},"Action":["sns:Subscribe"],"Resource":["arn:aws:sns:us-east-1:1234567890123456:mytopic"]}]}`),
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
