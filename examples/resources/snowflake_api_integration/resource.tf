resource "snowflake_api_integration" "api_integration" {
  name = "aws_integration"
  api_provider = "aws_api_gateway"
  api_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
  api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
  enabled = true
}