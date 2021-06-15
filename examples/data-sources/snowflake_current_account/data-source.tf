data "snowflake_current_account" "this" {}

resource "aws_ssm_parameter" "snowflake_account_url" {
  name  = "/snowflake/account_url"
  type  = "String"
  value = data.snowflake_current_account.this.url
}
