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
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
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
