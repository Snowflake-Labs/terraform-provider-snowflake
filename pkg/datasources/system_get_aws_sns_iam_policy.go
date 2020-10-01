package datasources

import (
	"database/sql"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var systemGetAWSSNSIAMPolicySchema = map[string]*schema.Schema{
	"aws_sns_topic_arn": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Amazon Resource Name (ARN) of the SNS topic for your S3 bucket",
	},

	"aws_sns_topic_policy_json": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "IAM policy for Snowflake’s SQS queue to subscribe to this topic",
	},
}

func SystemGetAWSSNSIAMPolicy() *schema.Resource {
	return &schema.Resource{
		Read:   ReadSystemGetAWSSNSIAMPolicy,
		Schema: systemGetAWSSNSIAMPolicySchema,
	}
}

// ReadSystemGetAWSSNSIAMPolicy implements schema.ReadFunc
func ReadSystemGetAWSSNSIAMPolicy(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	awsSNSTopicArn := data.Get("aws_sns_topic_arn").(string)

	sel := snowflake.SystemGetAWSSNSIAMPolicy(awsSNSTopicArn).Select()
	row := snowflake.QueryRow(db, sel)
	policy, err := snowflake.ScanAWSSNSIAMPolicy(row)
	if err != nil {
		return err
	}

	data.SetId(awsSNSTopicArn)
	data.Set("aws_sns_topic_policy_json", policy.Policy)
	return nil
}
