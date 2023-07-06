resource "snowflake_password_policy" "default" {
  database = "prod"
  schema   = "security"
  name     = "default_policy"
}

resource "snowflake_account_password_policy_attachment" "attachment" {
  password_policy = snowflake_password_policy.default.qualified_name
}
