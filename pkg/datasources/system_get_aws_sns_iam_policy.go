package datasources

import (
	"database/sql"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
func ReadSystemGetAWSSNSIAMPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	awsSNSTopicArn := d.Get("aws_sns_topic_arn").(string)

	sel := snowflake.SystemGetAWSSNSIAMPolicy(awsSNSTopicArn).Select()
	row := snowflake.QueryRow(db, sel)
	policy, err := snowflake.ScanAWSSNSIAMPolicy(row)
	if err == sql.ErrNoRows {
		log.Printf("[WARN] system_get_aws_sns_iam_policy (%s) not found, removing from state file", data.Id())
		data.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	d.SetId(awsSNSTopicArn)
	return d.Set("aws_sns_topic_policy_json", policy.Policy)
}
