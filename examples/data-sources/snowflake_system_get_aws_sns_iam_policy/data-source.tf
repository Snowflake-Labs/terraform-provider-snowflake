data "snowflake_system_get_aws_sns_iam_policy" "snowflake_policy" {
  aws_sns_topic_arn = "<aws_sns_topic_arn>"
}
