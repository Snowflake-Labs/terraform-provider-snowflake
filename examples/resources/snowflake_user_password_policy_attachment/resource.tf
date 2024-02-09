resource "snowflake_user" "user" {
  name = "USER_NAME"
}
resource "snowflake_password_policy" "pp" {
  database = "prod"
  schema   = "security"
  name     = "default_policy"
}

resource "snowflake_user_password_policy_attachment" "ppa" {
  password_policy_database = snowflake_password_policy.pp.database
  password_policy_schema   = snowflake_password_policy.pp.schema
  password_policy_name     = snowflake_password_policy.pp.name
  user_name                = snowflake_user.user.name
}
