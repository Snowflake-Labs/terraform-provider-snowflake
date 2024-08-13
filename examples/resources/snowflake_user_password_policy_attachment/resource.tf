resource "snowflake_user" "user" {
  name = "USER_NAME"
}
resource "snowflake_password_policy" "pp" {
  database = "prod"
  schema   = "security"
  name     = "default_policy"
}

resource "snowflake_user_password_policy_attachment" "ppa" {
  password_policy_name = snowflake_password_policy.pp.fully_qualified_name
  user_name            = snowflake_user.user.name
}
