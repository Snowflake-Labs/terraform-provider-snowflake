package snowflake

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// SystemGetAWSSNSIAMPolicyBuilder abstracts calling the SYSTEM$GET_AWS_SNS_IAM_POLICY system function
type SystemGetAWSSNSIAMPolicyBuilder struct {
	awsSnsTopicArn string
}

// SystemGetAWSSNSIAMPolicy returns a pointer to a builder that abstracts calling the the SYSTEM$GET_AWS_SNS_IAM_POLICY system function
func SystemGetAWSSNSIAMPolicy(awsSnsTopicArn string) *SystemGetAWSSNSIAMPolicyBuilder {
	return &SystemGetAWSSNSIAMPolicyBuilder{
		awsSnsTopicArn: awsSnsTopicArn,
	}
}

// Select generates the select statement for obtaining the aws sns iam policy
func (pb *SystemGetAWSSNSIAMPolicyBuilder) Select() string {
	return fmt.Sprintf(`SELECT SYSTEM$GET_AWS_SNS_IAM_POLICY('%v') AS "policy"`, pb.awsSnsTopicArn)
}

type awsSNSIAMPolicy struct {
	Policy string `db:"policy"`
}

// ScanAWSSNSIAMPolicy convert a result into a
func ScanAWSSNSIAMPolicy(row *sqlx.Row) (*awsSNSIAMPolicy, error) {
	p := &awsSNSIAMPolicy{}
	e := row.StructScan(p)
	return p, e
}
